package ratelimit

import (
	"context"
	"sync"
	"time"
)

// ============================================================================
// IN-MEMORY RATE LIMITER IMPLEMENTATION
// ============================================================================
// This is where BUSINESS LOGIC lives - not in the models
// Suitable for single-server deployments and development

// MemoryLimiter implements Limiter using in-memory storage
type MemoryLimiter struct {
	mu    sync.RWMutex
	logs  map[string]*RateLimitLog  // key: "identifier:action"
	rules map[string]*RateLimitRule // key: action

	// Cleanup configuration
	cleanupInterval time.Duration
	stopCleanup     chan struct{}
}

// NewMemoryLimiter creates a new in-memory rate limiter
func NewMemoryLimiter(rules []*RateLimitRule) *MemoryLimiter {
	limiter := &MemoryLimiter{
		logs:            make(map[string]*RateLimitLog),
		rules:           make(map[string]*RateLimitRule),
		cleanupInterval: 5 * time.Minute,
		stopCleanup:     make(chan struct{}),
	}

	// Load rules
	for _, rule := range rules {
		rule.LoadDurations()
		limiter.rules[rule.Action] = rule
	}

	// Start cleanup goroutine
	go limiter.cleanup()

	return limiter
}

// Check checks if action is allowed WITHOUT recording
// This is a read-only operation
func (l *MemoryLimiter) Check(ctx context.Context, identifier, action string) (bool, error) {
	status, err := l.GetStatus(ctx, identifier, action)
	if err != nil {
		return false, err
	}

	return status.IsAllowed(), nil
}

// RecordAttempt records an attempt and returns status
// THIS IS WHERE THE BUSINESS LOGIC IS
func (l *MemoryLimiter) RecordAttempt(ctx context.Context, identifier, action string) (*RateLimitStatus, error) {
	l.mu.Lock()
	defer l.mu.Unlock()

	// Get rule for action
	rule, exists := l.rules[action]
	if !exists || !rule.IsActive {
		// No rate limiting for this action
		return &RateLimitStatus{
			Identifier:     identifier,
			Action:         action,
			Count:          0,
			MaxAttempts:    0,
			RemainingTries: -1, // Unlimited
			Blocked:        false,
		}, nil
	}

	key := l.makeKey(identifier, action)
	log, exists := l.logs[key]
	now := time.Now()

	// Early return if already blocked (consistent with Redis implementation)
	// This prevents Count from increasing while blocked
	if exists && log.IsCurrentlyBlocked() {
		return &RateLimitStatus{
			Identifier:     identifier,
			Action:         action,
			Count:          log.Count,
			MaxAttempts:    rule.MaxAttempts,
			RemainingTries: 0, // Already blocked = no remaining tries
			WindowEnd:      log.WindowEnd,
			Blocked:        true,
			BlockedUntil:   log.BlockedUntil,
		}, nil
	}

	if !exists || !log.IsWindowActive() {
		// Create new window
		log = &RateLimitLog{
			Identifier:  identifier,
			Action:      action,
			Count:       1,
			WindowStart: now,
			WindowEnd:   now.Add(rule.WindowSize),
			Blocked:     false,
			CreatedAt:   now,
			UpdatedAt:   now,
		}
		l.logs[key] = log
	} else {
		// Increment counter in existing window
		log.Count++
		log.UpdatedAt = now

		// Check if should block
		if log.Count >= rule.MaxAttempts && !log.Blocked {
			log.Blocked = true
			blockedUntil := now.Add(rule.BlockDuration)
			log.BlockedUntil = &blockedUntil
		}
	}

	// Build status response
	// Calculate remaining tries (ensure non-negative)
	var remaining int
	if log.Blocked {
		// Already blocked = no remaining tries
		remaining = 0
	} else {
		remaining = rule.MaxAttempts - log.Count
		if remaining < 0 {
			remaining = 0
		}
	}

	status := &RateLimitStatus{
		Identifier:     identifier,
		Action:         action,
		Count:          log.Count,
		MaxAttempts:    rule.MaxAttempts,
		RemainingTries: remaining,
		WindowEnd:      log.WindowEnd,
		Blocked:        log.Blocked,
		BlockedUntil:   log.BlockedUntil,
	}

	return status, nil
}

// GetStatus gets current status without recording
func (l *MemoryLimiter) GetStatus(ctx context.Context, identifier, action string) (*RateLimitStatus, error) {
	l.mu.RLock()
	defer l.mu.RUnlock()

	rule, exists := l.rules[action]
	if !exists || !rule.IsActive {
		return &RateLimitStatus{
			Identifier:     identifier,
			Action:         action,
			Count:          0,
			MaxAttempts:    0,
			RemainingTries: -1,
			Blocked:        false,
		}, nil
	}

	key := l.makeKey(identifier, action)
	log, exists := l.logs[key]

	if !exists || !log.IsWindowActive() {
		// No active rate limit
		// WindowEnd is estimated when no active window exists
		return &RateLimitStatus{
			Identifier:     identifier,
			Action:         action,
			Count:          0,
			MaxAttempts:    rule.MaxAttempts,
			RemainingTries: rule.MaxAttempts,
			WindowEnd:      time.Now().Add(rule.WindowSize),
			Blocked:        false,
		}, nil
	}

	// Check if still blocked
	blocked := log.IsCurrentlyBlocked()

	// Calculate remaining tries (ensure non-negative)
	var remaining int
	if blocked {
		// Blocked = no remaining tries
		remaining = 0
	} else {
		remaining = rule.MaxAttempts - log.Count
		if remaining < 0 {
			remaining = 0
		}
	}

	status := &RateLimitStatus{
		Identifier:     identifier,
		Action:         action,
		Count:          log.Count,
		MaxAttempts:    rule.MaxAttempts,
		RemainingTries: remaining,
		WindowEnd:      log.WindowEnd,
		Blocked:        blocked,
		BlockedUntil:   log.BlockedUntil,
	}

	return status, nil
}

// Reset resets rate limit for identifier and action
func (l *MemoryLimiter) Reset(ctx context.Context, identifier, action string) error {
	l.mu.Lock()
	defer l.mu.Unlock()

	key := l.makeKey(identifier, action)
	delete(l.logs, key)

	return nil
}

// Block manually blocks an identifier
// Duration of 0 means indefinite block
func (l *MemoryLimiter) Block(ctx context.Context, identifier, action string, duration time.Duration) error {
	l.mu.Lock()
	defer l.mu.Unlock()

	if duration == 0 {
		// Indefinite block - use very long duration
		duration = 365 * 24 * time.Hour // 1 year
	}

	key := l.makeKey(identifier, action)
	now := time.Now()
	blockedUntil := now.Add(duration)

	log, exists := l.logs[key]
	if !exists {
		rule, ruleExists := l.rules[action]
		windowEnd := now.Add(1 * time.Hour)
		if ruleExists {
			windowEnd = now.Add(rule.WindowSize)
		}

		log = &RateLimitLog{
			Identifier:   identifier,
			Action:       action,
			Count:        999, // High count to indicate manual block
			WindowStart:  now,
			WindowEnd:    windowEnd,
			Blocked:      true,
			BlockedUntil: &blockedUntil,
			CreatedAt:    now,
			UpdatedAt:    now,
		}
		l.logs[key] = log
	} else {
		log.Blocked = true
		log.BlockedUntil = &blockedUntil
		log.UpdatedAt = now
	}

	return nil
}

// Unblock manually unblocks an identifier
func (l *MemoryLimiter) Unblock(ctx context.Context, identifier, action string) error {
	l.mu.Lock()
	defer l.mu.Unlock()

	key := l.makeKey(identifier, action)
	log, exists := l.logs[key]
	if exists {
		log.Blocked = false
		log.BlockedUntil = nil
		log.UpdatedAt = time.Now()
	}

	return nil
}

// ============================================================================
// RULE MANAGEMENT (Thread-safe)
// ============================================================================

// AddRule adds or updates a rate limit rule (thread-safe)
func (l *MemoryLimiter) AddRule(rule *RateLimitRule) {
	l.mu.Lock()
	defer l.mu.Unlock()

	rule.LoadDurations()
	l.rules[rule.Action] = rule
}

// RemoveRule removes a rate limit rule (thread-safe)
func (l *MemoryLimiter) RemoveRule(action string) {
	l.mu.Lock()
	defer l.mu.Unlock()

	delete(l.rules, action)
}

// GetRule gets a rate limit rule (thread-safe)
func (l *MemoryLimiter) GetRule(action string) (*RateLimitRule, bool) {
	l.mu.RLock()
	defer l.mu.RUnlock()

	rule, exists := l.rules[action]
	return rule, exists
}

// ListRules returns all configured rules (thread-safe)
func (l *MemoryLimiter) ListRules() []*RateLimitRule {
	l.mu.RLock()
	defer l.mu.RUnlock()

	rules := make([]*RateLimitRule, 0, len(l.rules))
	for _, rule := range l.rules {
		rules = append(rules, rule)
	}

	return rules
}

// ============================================================================
// HELPER METHODS
// ============================================================================

// makeKey creates a key for the rate limit log
// Optimized: uses string concatenation instead of fmt.Sprintf
func (l *MemoryLimiter) makeKey(identifier, action string) string {
	return identifier + ":" + action
}

// cleanup periodically removes expired logs
func (l *MemoryLimiter) cleanup() {
	ticker := time.NewTicker(l.cleanupInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			l.mu.Lock()

			for key, log := range l.logs {
				// Remove if window expired and not blocked, or if block expired
				if !log.IsWindowActive() && !log.IsCurrentlyBlocked() {
					delete(l.logs, key)
				}
			}

			l.mu.Unlock()

		case <-l.stopCleanup:
			return
		}
	}
}

// ============================================================================
// LIFECYCLE MANAGEMENT
// ============================================================================

// Close stops the cleanup goroutine
// Call this when shutting down to prevent goroutine leak
func (l *MemoryLimiter) Close() error {
	close(l.stopCleanup)
	return nil
}

// SetCleanupInterval changes the cleanup interval
// Useful for testing or tuning performance
func (l *MemoryLimiter) SetCleanupInterval(interval time.Duration) {
	l.mu.Lock()
	defer l.mu.Unlock()

	l.cleanupInterval = interval
}

// ============================================================================
// STATISTICS & MONITORING
// ============================================================================

// GetStats returns current limiter statistics
// Useful for monitoring and debugging
func (l *MemoryLimiter) GetStats() map[string]interface{} {
	l.mu.RLock()
	defer l.mu.RUnlock()

	blocked := 0
	active := 0

	for _, log := range l.logs {
		if log.IsCurrentlyBlocked() {
			blocked++
		}
		if log.IsWindowActive() {
			active++
		}
	}

	return map[string]interface{}{
		"total_logs":       len(l.logs),
		"active_windows":   active,
		"blocked_ids":      blocked,
		"configured_rules": len(l.rules),
		"store_type":       "memory",
		"cleanup_interval": l.cleanupInterval.String(),
	}
}

// CleanExpired forces immediate cleanup of expired entries
// Normally cleanup runs automatically, but this can be useful for testing
func (l *MemoryLimiter) CleanExpired(ctx context.Context) error {
	l.mu.Lock()
	defer l.mu.Unlock()

	for key, log := range l.logs {
		if !log.IsWindowActive() && !log.IsCurrentlyBlocked() {
			delete(l.logs, key)
		}
	}

	return nil
}

// GetLog returns a copy of the rate limit log for debugging
// Returns nil if log doesn't exist
func (l *MemoryLimiter) GetLog(identifier, action string) *RateLimitLog {
	l.mu.RLock()
	defer l.mu.RUnlock()

	key := l.makeKey(identifier, action)
	log, exists := l.logs[key]
	if !exists {
		return nil
	}

	// Return a copy to prevent external modification
	logCopy := *log
	return &logCopy
}

// Clear removes all rate limit logs (useful for testing)
func (l *MemoryLimiter) Clear() {
	l.mu.Lock()
	defer l.mu.Unlock()

	l.logs = make(map[string]*RateLimitLog)
}
