package ratelimit

import (
	"context"
	"crypto/sha1"
	"fmt"
	"strconv"
	"sync"
	"time"

	"github.com/redis/go-redis/v9"
)

// ============================================================================
// REDIS LUA SCRIPTS
// ============================================================================
// Using Lua scripts ensures atomic operations and better performance

const (
	// luaIncrWithExpire atomically increments counter and sets expiry
	// This prevents race condition between INCR and EXPIRE
	luaIncrWithExpire = `
		local key = KEYS[1]
		local window = tonumber(ARGV[1])
		
		local count = redis.call("INCR", key)
		
		if count == 1 then
			redis.call("EXPIRE", key, window)
		end
		
		return count
	`

	// luaGetMulti gets counter, TTL, and block status in one call
	// Reduces Redis round-trips from 3 to 1
	luaGetMulti = `
		local counterKey = KEYS[1]
		local blockKey = KEYS[2]
		
		local count = redis.call("GET", counterKey)
		local ttl = redis.call("TTL", counterKey)
		local blockVal = redis.call("GET", blockKey)
		
		return {count or "0", tostring(ttl), blockVal or ""}
	`
)

// ============================================================================
// REDIS RATE LIMITER IMPLEMENTATION
// ============================================================================

// RedisLimiter implements Limiter using Redis
// Thread-safe and production-ready with atomic operations
type RedisLimiter struct {
	client *redis.Client

	// Use sync.RWMutex for thread-safe rule access
	mu    sync.RWMutex
	rules map[string]*RateLimitRule

	// Lua scripts (preloaded for better performance)
	scriptIncr *redis.Script
	scriptGet  *redis.Script
}

// NewRedisLimiter creates a new Redis-backed rate limiter
func NewRedisLimiter(client *redis.Client, rules []*RateLimitRule) *RedisLimiter {
	limiter := &RedisLimiter{
		client:     client,
		rules:      make(map[string]*RateLimitRule),
		scriptIncr: redis.NewScript(luaIncrWithExpire),
		scriptGet:  redis.NewScript(luaGetMulti),
	}

	// Load rules (thread-safe)
	for _, rule := range rules {
		rule.LoadDurations()
		limiter.rules[rule.Action] = rule
	}

	return limiter
}

// Check checks if action is allowed WITHOUT recording
// This is a read-only operation
func (l *RedisLimiter) Check(ctx context.Context, identifier, action string) (bool, error) {
	status, err := l.GetStatus(ctx, identifier, action)
	if err != nil {
		return false, err
	}

	return status.IsAllowed(), nil
}

// RecordAttempt records an attempt using Redis with atomic operations
func (l *RedisLimiter) RecordAttempt(ctx context.Context, identifier, action string) (*RateLimitStatus, error) {
	// Get rule (thread-safe read)
	l.mu.RLock()
	rule, exists := l.rules[action]
	l.mu.RUnlock()

	if !exists || !rule.IsActive {
		// No rate limiting for this action
		return &RateLimitStatus{
			Identifier:     identifier,
			Action:         action,
			RemainingTries: -1,
			Blocked:        false,
		}, nil
	}

	key := l.makeKey(identifier, action)
	blockKey := l.makeBlockKey(identifier, action)
	now := time.Now()

	// Check if blocked (using timestamp instead of JSON for efficiency)
	blockedUntilStr, err := l.client.Get(ctx, blockKey).Result()
	if err == nil && blockedUntilStr != "" {
		blockedUntilUnix, err := strconv.ParseInt(blockedUntilStr, 10, 64)
		if err == nil {
			blockedUntil := time.Unix(blockedUntilUnix, 0)
			if now.Before(blockedUntil) {
				// Still blocked
				return &RateLimitStatus{
					Identifier:     identifier,
					Action:         action,
					Count:          rule.MaxAttempts,
					MaxAttempts:    rule.MaxAttempts,
					RemainingTries: 0,
					WindowEnd:      blockedUntil,
					Blocked:        true,
					BlockedUntil:   &blockedUntil,
				}, nil
			}
		}
	}

	// Atomically increment counter and set expiry using Lua script
	// This prevents race condition between INCR and EXPIRE
	countResult, err := l.scriptIncr.Run(ctx, l.client,
		[]string{key},
		int(rule.WindowSize.Seconds()),
	).Result()
	if err != nil {
		return nil, fmt.Errorf("redis incr script failed: %w", err)
	}

	count, ok := countResult.(int64)
	if !ok {
		return nil, fmt.Errorf("unexpected script result type: %T", countResult)
	}

	// Get TTL for window end (with safe fallback)
	ttl, err := l.client.TTL(ctx, key).Result()
	if err != nil || ttl <= 0 {
		// Fallback to window size if TTL query fails or returns invalid value
		// -1 = no expiry, -2 = key doesn't exist
		ttl = rule.WindowSize
	}
	windowEnd := now.Add(ttl)

	// Check if should block
	shouldBlock := int(count) >= rule.MaxAttempts
	var blockedUntil *time.Time

	if shouldBlock {
		blockedTime := now.Add(rule.BlockDuration)
		blockedUntil = &blockedTime

		// Store block as Unix timestamp (more efficient than JSON)
		blockedUntilUnix := blockedTime.Unix()
		err := l.client.Set(ctx, blockKey,
			strconv.FormatInt(blockedUntilUnix, 10),
			rule.BlockDuration,
		).Err()
		if err != nil {
			// Log but don't fail - block will be enforced on next attempt
			// In production, you should log this error
		}
	}

	// Calculate remaining tries (ensure non-negative)
	remaining := rule.MaxAttempts - int(count)
	if remaining < 0 {
		remaining = 0
	}

	status := &RateLimitStatus{
		Identifier:     identifier,
		Action:         action,
		Count:          int(count),
		MaxAttempts:    rule.MaxAttempts,
		RemainingTries: remaining,
		WindowEnd:      windowEnd,
		Blocked:        shouldBlock,
		BlockedUntil:   blockedUntil,
	}

	return status, nil
}

// GetStatus gets current status from Redis using optimized multi-get
func (l *RedisLimiter) GetStatus(ctx context.Context, identifier, action string) (*RateLimitStatus, error) {
	// Get rule (thread-safe read)
	l.mu.RLock()
	rule, exists := l.rules[action]
	l.mu.RUnlock()

	if !exists || !rule.IsActive {
		return &RateLimitStatus{
			Identifier:     identifier,
			Action:         action,
			RemainingTries: -1,
			Blocked:        false,
		}, nil
	}

	key := l.makeKey(identifier, action)
	blockKey := l.makeBlockKey(identifier, action)
	now := time.Now()

	// Use Lua script to get all values in one round-trip
	result, err := l.scriptGet.Run(ctx, l.client,
		[]string{key, blockKey},
	).Result()

	if err != nil {
		return nil, fmt.Errorf("redis get script failed: %w", err)
	}

	// Parse result array
	values, ok := result.([]interface{})
	if !ok || len(values) != 3 {
		return nil, fmt.Errorf("unexpected script result format")
	}

	// Parse count
	countStr, _ := values[0].(string)
	count, _ := strconv.Atoi(countStr)

	// Parse TTL
	ttlStr, _ := values[1].(string)
	ttlSeconds, _ := strconv.ParseInt(ttlStr, 10, 64)
	ttl := time.Duration(ttlSeconds) * time.Second
	if ttl <= 0 {
		ttl = rule.WindowSize
	}
	windowEnd := now.Add(ttl)

	// Parse block status
	blockStr, _ := values[2].(string)
	var blocked bool
	var blockedUntil *time.Time

	if blockStr != "" {
		blockedUntilUnix, err := strconv.ParseInt(blockStr, 10, 64)
		if err == nil {
			blockedTime := time.Unix(blockedUntilUnix, 0)
			if now.Before(blockedTime) {
				blocked = true
				blockedUntil = &blockedTime
			}
		}
	}

	// Calculate remaining tries (ensure non-negative)
	remaining := rule.MaxAttempts - count
	if remaining < 0 {
		remaining = 0
	}

	status := &RateLimitStatus{
		Identifier:     identifier,
		Action:         action,
		Count:          count,
		MaxAttempts:    rule.MaxAttempts,
		RemainingTries: remaining,
		WindowEnd:      windowEnd,
		Blocked:        blocked,
		BlockedUntil:   blockedUntil,
	}

	return status, nil
}

// Reset resets rate limit in Redis using pipeline for efficiency
func (l *RedisLimiter) Reset(ctx context.Context, identifier, action string) error {
	key := l.makeKey(identifier, action)
	blockKey := l.makeBlockKey(identifier, action)

	// Use pipeline to delete both keys in one round-trip
	pipe := l.client.Pipeline()
	pipe.Del(ctx, key)
	pipe.Del(ctx, blockKey)
	_, err := pipe.Exec(ctx)

	return err
}

// Block manually blocks an identifier in Redis
// Duration of 0 means indefinite block (until manual unblock)
func (l *RedisLimiter) Block(ctx context.Context, identifier, action string, duration time.Duration) error {
	if duration == 0 {
		// Indefinite block - use very long duration
		duration = 365 * 24 * time.Hour // 1 year
	}

	blockKey := l.makeBlockKey(identifier, action)
	blockedUntil := time.Now().Add(duration)

	// Store as Unix timestamp (efficient)
	blockedUntilUnix := blockedUntil.Unix()
	return l.client.Set(ctx, blockKey,
		strconv.FormatInt(blockedUntilUnix, 10),
		duration,
	).Err()
}

// Unblock manually unblocks an identifier in Redis
func (l *RedisLimiter) Unblock(ctx context.Context, identifier, action string) error {
	key := l.makeKey(identifier, action)
	blockKey := l.makeBlockKey(identifier, action)

	// Delete both counter and block keys
	pipe := l.client.Pipeline()
	pipe.Del(ctx, key)
	pipe.Del(ctx, blockKey)
	_, err := pipe.Exec(ctx)

	return err
}

// ============================================================================
// HELPER METHODS
// ============================================================================

// makeKey generates Redis key for rate limit counter
// Uses SHA1 hash of identifier to keep keys short and consistent
func (l *RedisLimiter) makeKey(identifier, action string) string {
	// Hash long identifiers (emails, UUIDs) to keep keys short
	hash := hashIdentifier(identifier)
	return fmt.Sprintf("rl:c:%s:%s", action, hash)
}

// makeBlockKey generates Redis key for block status
func (l *RedisLimiter) makeBlockKey(identifier, action string) string {
	hash := hashIdentifier(identifier)
	return fmt.Sprintf("rl:b:%s:%s", action, hash)
}

// hashIdentifier hashes identifier to keep Redis keys short
// Long identifiers (emails, UUIDs) → 40 character SHA1 hash
func hashIdentifier(identifier string) string {
	// For short identifiers (IPs), no need to hash
	if len(identifier) <= 20 {
		return identifier
	}

	// Hash long identifiers
	h := sha1.New()
	h.Write([]byte(identifier))
	return fmt.Sprintf("%x", h.Sum(nil))
}

// ============================================================================
// RULE MANAGEMENT (Thread-safe)
// ============================================================================

// AddRule adds or updates a rate limit rule (thread-safe)
// Rules are immutable after initialization in production,
// but this method is useful for dynamic rule updates
func (l *RedisLimiter) AddRule(rule *RateLimitRule) {
	rule.LoadDurations()

	l.mu.Lock()
	l.rules[rule.Action] = rule
	l.mu.Unlock()
}

// RemoveRule removes a rate limit rule (thread-safe)
func (l *RedisLimiter) RemoveRule(action string) {
	l.mu.Lock()
	delete(l.rules, action)
	l.mu.Unlock()
}

// GetRule gets a rate limit rule (thread-safe)
func (l *RedisLimiter) GetRule(action string) (*RateLimitRule, bool) {
	l.mu.RLock()
	defer l.mu.RUnlock()

	rule, exists := l.rules[action]
	return rule, exists
}

// ListRules returns all configured rules (thread-safe)
func (l *RedisLimiter) ListRules() []*RateLimitRule {
	l.mu.RLock()
	defer l.mu.RUnlock()

	rules := make([]*RateLimitRule, 0, len(l.rules))
	for _, rule := range l.rules {
		rules = append(rules, rule)
	}

	return rules
}

// ============================================================================
// CONNECTION MANAGEMENT
// ============================================================================

// Close closes the Redis connection
func (l *RedisLimiter) Close() error {
	return l.client.Close()
}

// Ping checks Redis connectivity
func (l *RedisLimiter) Ping(ctx context.Context) error {
	return l.client.Ping(ctx).Err()
}

// ============================================================================
// STATISTICS & MONITORING
// ============================================================================

// GetStats returns statistics about rate limiting
// Useful for monitoring and debugging
func (l *RedisLimiter) GetStats(ctx context.Context) (map[string]interface{}, error) {
	// Get total number of rate limit keys
	counterKeys, err := l.client.Keys(ctx, "rl:c:*").Result()
	if err != nil {
		return nil, err
	}

	blockKeys, err := l.client.Keys(ctx, "rl:b:*").Result()
	if err != nil {
		return nil, err
	}

	l.mu.RLock()
	ruleCount := len(l.rules)
	l.mu.RUnlock()

	stats := map[string]interface{}{
		"rules_configured": ruleCount,
		"active_counters":  len(counterKeys),
		"active_blocks":    len(blockKeys),
		"store_type":       "redis",
		"redis_connected":  l.client.Ping(ctx).Err() == nil,
	}

	return stats, nil
}

// CleanExpired removes expired rate limit keys
// Redis does this automatically with TTL, but this can be used
// to force cleanup if needed (e.g., after changing rules)
func (l *RedisLimiter) CleanExpired(ctx context.Context) error {
	// Redis automatically removes expired keys with TTL
	// This is just a placeholder for interface compatibility
	// In practice, Redis handles this automatically
	return nil
}
