package ratelimit

import "time"

// ============================================================================
// RATE LIMIT MODELS (Data + Minimal Helpers)
// ============================================================================

// RateLimitLog represents rate limiting record in storage
// This is a DATA MODEL - business logic belongs in the limiter implementation
type RateLimitLog struct {
	ID           int        `json:"id" db:"id"`
	Identifier   string     `json:"identifier" db:"identifier"` // IP, user_id, email
	Action       string     `json:"action" db:"action"`
	Count        int        `json:"count" db:"count"`
	WindowStart  time.Time  `json:"window_start" db:"window_start"`
	WindowEnd    time.Time  `json:"window_end" db:"window_end"`
	Blocked      bool       `json:"blocked" db:"blocked"`
	BlockedUntil *time.Time `json:"blocked_until,omitempty" db:"blocked_until"`
	CreatedAt    time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at" db:"updated_at"`
}

// IsCurrentlyBlocked is a simple query helper (not business logic)
// Just checks current state, doesn't make decisions
func (r *RateLimitLog) IsCurrentlyBlocked() bool {
	if !r.Blocked || r.BlockedUntil == nil {
		return false
	}
	return time.Now().Before(*r.BlockedUntil)
}

// IsWindowActive is a simple query helper
func (r *RateLimitLog) IsWindowActive() bool {
	now := time.Now()
	return now.After(r.WindowStart) && now.Before(r.WindowEnd)
}

// ============================================================================
// RATE LIMIT RULE (Configuration)
// ============================================================================

// RateLimitRule represents rate limiting configuration
// This is pure configuration - no business logic
type RateLimitRule struct {
	ID                   int           `json:"id" db:"id"`
	Action               string        `json:"action" db:"action"`
	MaxAttempts          int           `json:"max_attempts" db:"max_attempts"`
	WindowSize           time.Duration `json:"-"` // Not stored, computed from seconds
	BlockDuration        time.Duration `json:"-"` // Not stored, computed from seconds
	WindowSizeSeconds    int           `json:"window_size_seconds" db:"window_size_seconds"`
	BlockDurationSeconds int           `json:"block_duration_seconds" db:"block_duration_seconds"`
	IsActive             bool          `json:"is_active" db:"is_active"`
	CreatedAt            time.Time     `json:"created_at" db:"created_at"`
	UpdatedAt            time.Time     `json:"updated_at" db:"updated_at"`
}

// LoadDurations converts stored seconds to time.Duration for easier use
func (r *RateLimitRule) LoadDurations() {
	r.WindowSize = time.Duration(r.WindowSizeSeconds) * time.Second
	r.BlockDuration = time.Duration(r.BlockDurationSeconds) * time.Second
}

// ============================================================================
// RATE LIMIT STATUS (Query Result)
// ============================================================================

// RateLimitStatus represents current rate limit status
// This is a READ MODEL - just data for clients
type RateLimitStatus struct {
	Identifier     string
	Action         string
	Count          int
	MaxAttempts    int
	RemainingTries int
	WindowEnd      time.Time
	Blocked        bool
	BlockedUntil   *time.Time
}

// IsAllowed is a simple query helper
func (s *RateLimitStatus) IsAllowed() bool {
	if s.Blocked && s.BlockedUntil != nil {
		return time.Now().After(*s.BlockedUntil)
	}
	return s.RemainingTries > 0
}

// TimeUntilReset returns time until rate limit resets
func (s *RateLimitStatus) TimeUntilReset() time.Duration {
	if s.Blocked && s.BlockedUntil != nil {
		return time.Until(*s.BlockedUntil)
	}
	return time.Until(s.WindowEnd)
}

// ============================================================================
// COMMON ACTIONS (Constants)
// ============================================================================

// Common rate limit actions
const (
	ActionLogin         = "login"
	ActionPasswordReset = "password_reset"
	ActionEmailVerify   = "email_verify"
	ActionResendEmail   = "resend_email"
	ActionAPICall       = "api_call"
	ActionRegistration  = "registration"
	ActionOTPRequest    = "otp_request"
)
