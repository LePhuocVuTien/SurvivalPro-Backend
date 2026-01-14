package ratelimit

import (
	"context"
	"fmt"
	"time"
)

// ============================================================================
// RATE LIMITER INTERFACE
// ============================================================================
// This defines the contract for rate limiting implementations.
// Implementations can use Redis, memory, database, or any other storage.

// Limiter defines the interface for rate limiting operations.
// All business logic lives in the implementations (MemoryLimiter, RedisLimiter),
// not in the models.
type Limiter interface {
	// Check checks current rate limit status WITHOUT recording an attempt.
	// This is a read-only operation with no side effects.
	//
	// identifier can be: user_id, ip_address, email, api_key, etc.
	// Format is implementation-defined (e.g., "192.168.1.1", "user:123").
	//
	// Returns true if the action is allowed (not rate limited).
	// Returns false if the action is blocked or rate limit exceeded.
	//
	// Use this when you want to check status without affecting the counter.
	// For most use cases, use RecordAttempt() instead which combines check + record.
	Check(ctx context.Context, identifier, action string) (bool, error)

	// RecordAttempt records an attempt and returns the current rate limit status.
	// This is the primary method for rate limiting - it both checks AND records.
	//
	// The decision logic (should block? increment counter? start new window?)
	// lives in the implementation, not in the models.
	//
	// Returns RateLimitStatus with detailed information including:
	// - Whether the action is allowed
	// - Remaining attempts
	// - When the limit resets
	// - Block information if blocked
	//
	// Example usage:
	//   status, err := limiter.RecordAttempt(ctx, clientIP, ActionLogin)
	//   if err != nil {
	//       return err
	//   }
	//   if !status.IsAllowed() {
	//       return RateLimitError{Status: status}
	//   }
	RecordAttempt(ctx context.Context, identifier, action string) (*RateLimitStatus, error)

	// GetStatus gets current rate limit status without recording an attempt.
	// Similar to Check(), but returns detailed status information.
	//
	// Use this when you need detailed information (like remaining attempts,
	// time until reset) without affecting the counter.
	GetStatus(ctx context.Context, identifier, action string) (*RateLimitStatus, error)

	// Reset resets rate limit for identifier and action.
	// This removes all rate limit state for the identifier-action pair.
	//
	// Common use cases:
	// - Reset after successful authentication (clear failed login attempts)
	// - Admin action to unblock a user
	// - Testing/development purposes
	Reset(ctx context.Context, identifier, action string) error

	// Block manually blocks an identifier for an action.
	//
	// If duration is nil or zero, blocks indefinitely (until manual Unblock).
	// If duration is provided, blocks for that specific duration.
	//
	// This is typically used by admins to manually block abusive users.
	//
	// Example:
	//   // Block for 24 hours
	//   limiter.Block(ctx, "192.168.1.1", ActionLogin, 24*time.Hour)
	//
	//   // Block indefinitely
	//   limiter.Block(ctx, "user:123", ActionAPICall, 0)
	Block(ctx context.Context, identifier, action string, duration time.Duration) error

	// Unblock manually unblocks an identifier for an action.
	// This removes any manual or automatic blocks.
	//
	// Typically used by admins to unblock users.
	Unblock(ctx context.Context, identifier, action string) error
}

// ============================================================================
// STORAGE REPOSITORY INTERFACE
// ============================================================================

// Repository defines the interface for rate limit persistence.
// This is optional - implementations like MemoryLimiter don't need it.
// Implementations using database storage can implement this interface.
type Repository interface {
	// GetLog gets rate limit log for identifier and action
	GetLog(ctx context.Context, identifier, action string) (*RateLimitLog, error)

	// CreateLog creates a new rate limit log
	CreateLog(ctx context.Context, log *RateLimitLog) error

	// UpdateLog updates existing rate limit log
	UpdateLog(ctx context.Context, log *RateLimitLog) error

	// DeleteLog deletes rate limit log
	DeleteLog(ctx context.Context, identifier, action string) error

	// GetRule gets rate limit rule for action
	GetRule(ctx context.Context, action string) (*RateLimitRule, error)

	// ListRules lists all active rate limit rules
	ListRules(ctx context.Context) ([]*RateLimitRule, error)

	// SaveRule saves or updates rate limit rule
	SaveRule(ctx context.Context, rule *RateLimitRule) error

	// CleanExpired removes expired logs (for maintenance)
	CleanExpired(ctx context.Context, before time.Time) error
}

// ============================================================================
// RESULT TYPES
// ============================================================================

// CheckResult provides rich information about a rate limit check.
// This type is designed for middleware and handler layers that need
// detailed information to construct proper HTTP responses.
//
// Note: The Limiter interface methods return simpler types (bool, *RateLimitStatus)
// for easier use in business logic. This type is for presentation layer use.
//
// Example usage in HTTP middleware:
//
//	result := &CheckResult{
//	    Allowed: status.IsAllowed(),
//	    Status: status,
//	    RetryAfter: status.TimeUntilReset(),
//	    ResetAt: &status.WindowEnd,
//	}
//	if !result.Allowed {
//	    w.Header().Set("Retry-After", fmt.Sprintf("%d", int(result.RetryAfter.Seconds())))
//	    w.Header().Set("X-RateLimit-Reset", result.ResetAt.Format(time.RFC3339))
//	    http.Error(w, "Rate limit exceeded", http.StatusTooManyRequests)
//	}
type CheckResult struct {
	// Allowed indicates if the action is allowed
	Allowed bool

	// Status contains detailed rate limit information
	Status *RateLimitStatus

	// RetryAfter indicates how long to wait before retrying (for HTTP Retry-After header)
	// This is nil if Allowed is true
	RetryAfter *time.Duration

	// ResetAt indicates when the rate limit resets (for HTTP X-RateLimit-Reset header)
	// This is nil if no active rate limit
	ResetAt *time.Time
}

// ShouldRetry indicates if the client should retry later
func (r *CheckResult) ShouldRetry() bool {
	return !r.Allowed && r.RetryAfter != nil
}

// HTTPHeaders returns suggested HTTP headers for rate limit response
// This is a convenience method for HTTP handlers
func (r *CheckResult) HTTPHeaders() map[string]string {
	headers := make(map[string]string)

	if r.Status != nil {
		// Standard rate limit headers (draft RFC)
		headers["X-RateLimit-Limit"] = fmt.Sprintf("%d", r.Status.MaxAttempts)
		headers["X-RateLimit-Remaining"] = fmt.Sprintf("%d", r.Status.RemainingTries)

		if r.ResetAt != nil {
			headers["X-RateLimit-Reset"] = fmt.Sprintf("%d", r.ResetAt.Unix())
		}

		// Retry-After header for blocked requests (HTTP standard)
		if !r.Allowed && r.RetryAfter != nil {
			headers["Retry-After"] = fmt.Sprintf("%d", int(r.RetryAfter.Seconds()))
		}
	}

	return headers
}

// ============================================================================
// HELPER FUNCTIONS
// ============================================================================

// ToCheckResult converts a RateLimitStatus to CheckResult
// This is a convenience function for middleware/handler layers
func ToCheckResult(status *RateLimitStatus) *CheckResult {
	if status == nil {
		return &CheckResult{Allowed: true}
	}

	result := &CheckResult{
		Allowed: status.IsAllowed(),
		Status:  status,
	}

	if !result.Allowed {
		retryAfter := status.TimeUntilReset()
		result.RetryAfter = &retryAfter
		result.ResetAt = &status.WindowEnd

		if status.BlockedUntil != nil {
			result.ResetAt = status.BlockedUntil
		}
	}

	return result
}

// ============================================================================
// IDENTIFIER TYPES (for documentation and type safety)
// ============================================================================

// IdentifierType represents the type of identifier used for rate limiting
type IdentifierType string

const (
	// IdentifierIP represents IP address identifier (e.g., "192.168.1.1")
	// Use for preventing brute force attacks per IP
	IdentifierIP IdentifierType = "ip"

	// IdentifierUserID represents user ID identifier (e.g., "user:123")
	// Use for per-user rate limiting
	IdentifierUserID IdentifierType = "user_id"

	// IdentifierEmail represents email identifier (e.g., "email:user@example.com")
	// Use for email-based rate limiting (registration, password reset)
	IdentifierEmail IdentifierType = "email"

	// IdentifierAPIKey represents API key identifier (e.g., "apikey:abc123")
	// Use for API rate limiting
	IdentifierAPIKey IdentifierType = "api_key"

	// IdentifierDevice represents device identifier (e.g., "device:uuid")
	// Use for per-device rate limiting
	IdentifierDevice IdentifierType = "device"

	// IdentifierSession represents session identifier (e.g., "session:xyz")
	// Use for per-session rate limiting
	IdentifierSession IdentifierType = "session"
)

// FormatIdentifier formats an identifier with type prefix for clarity
// This is optional but recommended for better debugging and monitoring
//
// Example:
//
//	FormatIdentifier(IdentifierIP, "192.168.1.1") → "ip:192.168.1.1"
//	FormatIdentifier(IdentifierUserID, "123") → "user_id:123"
func FormatIdentifier(typ IdentifierType, value string) string {
	return string(typ) + ":" + value
}

// ============================================================================
// RATE LIMIT ERROR
// ============================================================================

// Error represents a rate limit error with detailed information
// This can be used in application layer to return structured errors
type Error struct {
	Action     string
	Identifier string
	Status     *RateLimitStatus
	Message    string
}

// Error implements error interface
func (e *Error) Error() string {
	if e.Message != "" {
		return e.Message
	}

	if e.Status != nil && e.Status.Blocked {
		return fmt.Sprintf("rate limit exceeded for %s, blocked until %s",
			e.Action, e.Status.BlockedUntil.Format(time.RFC3339))
	}

	return fmt.Sprintf("rate limit exceeded for %s, %d attempts remaining",
		e.Action, e.Status.RemainingTries)
}

// IsRateLimitError checks if an error is a rate limit error
func IsRateLimitError(err error) bool {
	_, ok := err.(*Error)
	return ok
}
