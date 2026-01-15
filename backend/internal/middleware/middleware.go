package middleware

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/LePhuocVuTien/SurvivalPro-Backend/infrastructure/ratelimit"
	models "github.com/LePhuocVuTien/SurvivalPro-Backend/internal/domain/user"
	"github.com/gofiber/fiber/v2"
)

// ============================================================================
// CONTEXT KEYS
// ============================================================================

const (
	UserContextKey = "user"
	UserIDKey      = "user_id"
	UserRoleKey    = "user_role"
)

// ============================================================================
// CONTEXT HELPERS
// ============================================================================

// SetUserInContext stores user in context
func SetUserInContext(c *fiber.Ctx, user *models.User) {
	c.Locals(UserContextKey, user)
	c.Locals(UserIDKey, user.ID)
	c.Locals(UserRoleKey, user.Role)
}

// GetUserFromContext retrieves user from context
func GetUserFromContext(c *fiber.Ctx) *models.User {
	user, ok := c.Locals(UserContextKey).(*models.User)
	if !ok {
		return nil
	}
	return user
}

// GetUserIDFromContext retrieves user ID from context
func GetUserIDFromContext(c *fiber.Ctx) (int, error) {
	userID, ok := c.Locals(UserIDKey).(int)
	if !ok {
		return 0, errors.New("user ID not found in context")
	}
	return userID, nil
}

// GetUserRoleFromContext retrieves user role from context
func GetUserRoleFromContext(c *fiber.Ctx) (models.UserRole, error) {
	role, ok := c.Locals(UserRoleKey).(models.UserRole)
	if !ok {
		return "", errors.New("user role not found in context")
	}
	return role, nil
}

// ============================================================================
// ROLE-BASED ACCESS CONTROL MIDDLEWARE
// ============================================================================

// RequireRole middleware checks if user has the required role
func RequireRole(roles ...models.UserRole) fiber.Handler {
	return func(c *fiber.Ctx) error {
		user := GetUserFromContext(c)
		if user == nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Unauthorized - user not found in context",
			})
		}

		// Check if user has any of the required roles
		hasRole := false
		for _, role := range roles {
			if user.Role == role {
				hasRole = true
				break
			}
		}

		if !hasRole {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"error":    "Forbidden - insufficient permissions",
				"required": roles,
				"actual":   user.Role,
			})
		}

		return c.Next()
	}
}

// RequireAdmin middleware - only admins can access
func RequireAdmin() fiber.Handler {
	return RequireRole(models.UserRoleAdmin)
}

// RequireLeaderOrAdmin middleware - leaders and admins can access
func RequireLeaderOrAdmin() fiber.Handler {
	return RequireRole(models.UserRoleAdmin, models.UserRoleLeader)
}

// RequirePermission middleware - checks permission hierarchy
func RequirePermission(minimumRole models.UserRole) fiber.Handler {
	return func(c *fiber.Ctx) error {
		user := GetUserFromContext(c)
		if user == nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Unauthorized - user not found in context",
			})
		}

		if !user.HasPermission(minimumRole) {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"error":          "Forbidden - insufficient permissions",
				"required_level": minimumRole,
				"your_level":     user.Role,
			})
		}

		return c.Next()
	}
}

// RequireActiveAccount middleware - only active accounts can access
func RequireActiveAccount() fiber.Handler {
	return func(c *fiber.Ctx) error {
		user := GetUserFromContext(c)
		if user == nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Unauthorized - user not found in context",
			})
		}

		if err := user.CanLogin(); err != nil {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"error":  "Account is not active",
				"status": user.AccountStatus,
			})
		}

		return c.Next()
	}
}

// RequireVerifiedEmail middleware - only verified emails can access
func RequireVerifiedEmail() fiber.Handler {
	return func(c *fiber.Ctx) error {
		user := GetUserFromContext(c)
		if user == nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Unauthorized - user not found in context",
			})
		}

		if !user.EmailVerified {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"error": "Email not verified",
			})
		}

		return c.Next()
	}
}

// RequireOwnershipOrAdmin checks if user owns the resource or is admin
func RequireOwnershipOrAdmin(userIDParam string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		user := GetUserFromContext(c)
		if user == nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Unauthorized - user not found in context",
			})
		}

		// Admins can access everything
		if user.IsAdmin() {
			return c.Next()
		}

		// Get resource owner ID from URL params
		resourceUserID, err := c.ParamsInt(userIDParam)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid user ID parameter",
			})
		}

		// Check ownership
		if user.ID != resourceUserID {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"error": "Forbidden - you can only access your own resources",
			})
		}

		return c.Next()
	}
}

// RequireOwnershipOrLeader checks if user owns the resource or is leader/admin
func RequireOwnershipOrLeader(userIDParam string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		user := GetUserFromContext(c)
		if user == nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Unauthorized - user not found in context",
			})
		}

		// Admins and leaders can access everything
		if user.IsAdmin() || user.IsLeader() {
			return c.Next()
		}

		// Get resource owner ID from URL params
		resourceUserID, err := c.ParamsInt(userIDParam)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid user ID parameter",
			})
		}

		// Check ownership
		if user.ID != resourceUserID {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"error": "Forbidden - you can only access your own resources",
			})
		}

		return c.Next()
	}
}

// ============================================================================
// REDIS RATE LIMIT MIDDLEWARE (FIXED)
// ============================================================================

// RedisRateLimitMiddleware creates rate limiting middleware using Redis
func RedisRateLimitMiddleware(limiter *ratelimit.RedisLimiter, action string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// ✅ FIX 1: Use c.Context() instead of context.Background()
		ctx := c.Context()

		// Get identifier (IP address by default)
		identifier := normalizeIdentifier(c.IP())

		// Record attempt and get status
		status, err := limiter.RecordAttempt(ctx, identifier, action)
		if err != nil {
			// On Redis error, log but allow request (fail open)
			// In production, log this error
			return c.Next()
		}

		// Set rate limit headers
		setRateLimitHeaders(c, status)

		// If blocked, return 429
		if status.Blocked {
			return c.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{
				"error":         "Rate limit exceeded",
				"message":       getRateLimitMessage(action, status),
				"retry_after":   getRetryAfter(status),
				"blocked_until": status.BlockedUntil,
			})
		}

		return c.Next()
	}
}

// RedisRateLimitByUserID creates rate limiting middleware using user ID
func RedisRateLimitByUserID(limiter *ratelimit.RedisLimiter, action string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Get user ID from context (set by JWT middleware)
		userID, err := GetUserIDFromContext(c)
		if err != nil {
			// Fallback to IP if user ID not available
			return RedisRateLimitMiddleware(limiter, action)(c)
		}

		// ✅ FIX 1: Use c.Context() instead of context.Background()
		ctx := c.Context()

		// Use user ID as identifier
		identifier := fmt.Sprintf("user:%d", userID)

		// Record attempt
		status, err := limiter.RecordAttempt(ctx, identifier, action)
		if err != nil {
			return c.Next() // Fail open
		}

		// Set headers
		setRateLimitHeaders(c, status)

		if status.Blocked {
			return c.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{
				"error":         "Rate limit exceeded",
				"message":       getRateLimitMessage(action, status),
				"retry_after":   getRetryAfter(status),
				"blocked_until": status.BlockedUntil,
			})
		}

		return c.Next()
	}
}

// RedisRateLimitByEmail creates rate limiting middleware using email
// ✅ FIX 2: Only parse email field, not entire body
func RedisRateLimitByEmail(limiter *ratelimit.RedisLimiter, action string, emailField string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// ✅ FIX 2: Parse only the email field to avoid body consumption
		type EmailPayload struct {
			Email string `json:"email"`
		}

		var payload EmailPayload
		if err := c.BodyParser(&payload); err != nil {
			// If can't parse body, fallback to IP
			return RedisRateLimitMiddleware(limiter, action)(c)
		}

		if payload.Email == "" {
			// If no email, fallback to IP
			return RedisRateLimitMiddleware(limiter, action)(c)
		}

		// ✅ FIX 1: Use c.Context() instead of context.Background()
		ctx := c.Context()

		// ✅ FIX 4: Normalize email (lowercase, trim spaces)
		email := normalizeEmail(payload.Email)
		identifier := fmt.Sprintf("email:%s", email)

		// Record attempt
		status, err := limiter.RecordAttempt(ctx, identifier, action)
		if err != nil {
			return c.Next() // Fail open
		}

		// Set headers
		setRateLimitHeaders(c, status)

		if status.Blocked {
			return c.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{
				"error":         "Rate limit exceeded",
				"message":       getRateLimitMessage(action, status),
				"retry_after":   getRetryAfter(status),
				"blocked_until": status.BlockedUntil,
			})
		}

		return c.Next()
	}
}

// CombinedRedisRateLimit combines IP and user-based rate limiting
func CombinedRedisRateLimit(limiter *ratelimit.RedisLimiter, action string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// ✅ FIX 1: Use c.Context() instead of context.Background()
		ctx := c.Context()

		// Check IP-based rate limit first
		ip := normalizeIdentifier(c.IP())
		status, err := limiter.RecordAttempt(ctx, ip, action)
		if err == nil && status.Blocked {
			setRateLimitHeaders(c, status)
			return c.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{
				"error":       "Rate limit exceeded (IP)",
				"message":     getRateLimitMessage(action, status),
				"retry_after": getRetryAfter(status),
			})
		}

		// Check user-based rate limit if authenticated
		userID, err := GetUserIDFromContext(c)
		if err == nil {
			identifier := fmt.Sprintf("user:%d", userID)
			status, err := limiter.RecordAttempt(ctx, identifier, action)
			if err == nil && status.Blocked {
				setRateLimitHeaders(c, status)
				return c.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{
					"error":       "Rate limit exceeded (User)",
					"message":     getRateLimitMessage(action, status),
					"retry_after": getRetryAfter(status),
				})
			}
		}

		return c.Next()
	}
}

// ============================================================================
// HELPER FUNCTIONS (FIXED)
// ============================================================================

// ✅ FIX 4: Normalize identifier (trim spaces, lowercase)
func normalizeIdentifier(identifier string) string {
	return strings.ToLower(strings.TrimSpace(identifier))
}

// ✅ FIX 4: Normalize email (trim spaces, lowercase)
func normalizeEmail(email string) string {
	return strings.ToLower(strings.TrimSpace(email))
}

// setRateLimitHeaders sets standard rate limit response headers
func setRateLimitHeaders(c *fiber.Ctx, status *ratelimit.RateLimitStatus) {
	c.Set("X-RateLimit-Limit", fmt.Sprintf("%d", status.MaxAttempts))
	c.Set("X-RateLimit-Remaining", fmt.Sprintf("%d", status.RemainingTries))

	if !status.WindowEnd.IsZero() {
		c.Set("X-RateLimit-Reset", fmt.Sprintf("%d", status.WindowEnd.Unix()))
	}

	if status.Blocked && status.BlockedUntil != nil {
		retryAfter := int(time.Until(*status.BlockedUntil).Seconds())
		if retryAfter < 0 {
			retryAfter = 0
		}
		c.Set("Retry-After", fmt.Sprintf("%d", retryAfter))
	}
}

// getRateLimitMessage returns user-friendly message based on action
func getRateLimitMessage(action string, status *ratelimit.RateLimitStatus) string {
	messages := map[string]string{
		"login":          "Too many login attempts. Please try again later.",
		"register":       "Too many registration attempts. Please try again later.",
		"password_reset": "Too many password reset attempts. Please try again later.",
		"email_verify":   "Too many email verification attempts. Please try again later.",
		"api":            "Too many API requests. Please slow down.",
		"upload":         "Too many file uploads. Please try again later.",
	}

	msg, ok := messages[action]
	if !ok {
		msg = "Rate limit exceeded. Please try again later."
	}

	// Add retry time if blocked
	if status.Blocked && status.BlockedUntil != nil {
		retrySeconds := int(time.Until(*status.BlockedUntil).Seconds())
		if retrySeconds > 0 {
			if retrySeconds < 60 {
				msg += fmt.Sprintf(" You can try again in %d seconds.", retrySeconds)
			} else {
				retryMinutes := retrySeconds / 60
				msg += fmt.Sprintf(" You can try again in %d minutes.", retryMinutes)
			}
		}
	}

	return msg
}

// getRetryAfter returns retry after duration in seconds
func getRetryAfter(status *ratelimit.RateLimitStatus) int {
	if status.BlockedUntil == nil {
		return 0
	}

	retryAfter := int(time.Until(*status.BlockedUntil).Seconds())
	if retryAfter < 0 {
		return 0
	}

	return retryAfter
}

// ============================================================================
// PERMISSION CHECKER HELPERS
// ============================================================================

// CanModifyUser checks if current user can modify target user
func CanModifyUser(currentUser *models.User, targetUserID int) bool {
	if currentUser == nil {
		return false
	}

	// Admin can modify anyone
	if currentUser.IsAdmin() {
		return true
	}

	// Users can only modify themselves
	return currentUser.ID == targetUserID
}

// CanChangeRole checks if current user can change role
func CanChangeRole(currentUser *models.User, targetRole models.UserRole) bool {
	if currentUser == nil {
		return false
	}

	// Only admins can change roles
	return currentUser.IsAdmin()
}

// CanChangeStatus checks if current user can change account status
func CanChangeStatus(currentUser *models.User) bool {
	if currentUser == nil {
		return false
	}

	// Only admins can change status
	return currentUser.IsAdmin()
}

// CanDeleteUser checks if current user can delete target user
func CanDeleteUser(currentUser *models.User, targetUserID int) bool {
	if currentUser == nil {
		return false
	}

	// Admin can delete anyone except themselves
	if currentUser.IsAdmin() && currentUser.ID != targetUserID {
		return true
	}

	// Users can delete themselves
	return currentUser.ID == targetUserID
}

// CanViewUser checks if current user can view target user
func CanViewUser(currentUser *models.User, targetUserID int) bool {
	if currentUser == nil {
		return false
	}

	// Admin and leader can view anyone
	if currentUser.IsAdmin() || currentUser.IsLeader() {
		return true
	}

	// Users can view themselves
	return currentUser.ID == targetUserID
}

// CanListUsers checks if current user can list users
func CanListUsers(currentUser *models.User) bool {
	if currentUser == nil {
		return false
	}

	// Admin and leader can list users
	return currentUser.IsAdmin() || currentUser.IsLeader()
}

// ============================================================================
// ERROR RESPONSES
// ============================================================================

// UnauthorizedResponse returns unauthorized error
func UnauthorizedResponse(c *fiber.Ctx) error {
	return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
		"error": "Unauthorized",
	})
}

// ForbiddenResponse returns forbidden error
func ForbiddenResponse(c *fiber.Ctx, message string) error {
	if message == "" {
		message = "Forbidden - insufficient permissions"
	}
	return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
		"error": message,
	})
}
