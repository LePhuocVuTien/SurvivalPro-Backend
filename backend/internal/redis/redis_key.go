package redis

import (
	"fmt"
	"time"
)

// CacheKey represents a cache key with TTL
type CacheKey struct {
	Key string
	TTL time.Duration
}

// Cache key prefixes for organization
const (
	PrefixUser         = "user"
	PrefixSession      = "session"
	PrefixToken        = "token"
	PrefixVerification = "verification"
	PrefixRateLimit    = "ratelimit"
	PrefixOTP          = "otp"
	PrefixAppointment  = "appointment"
	PrefixDoctor       = "doctor"
	PrefixPatient      = "patient"
	PrefixNotification = "notification"
	PrefixCache        = "cache"
)

// Default TTL values
const (
	TTLSession      = 24 * time.Hour
	TTLToken        = 15 * time.Minute
	TTLRefreshToken = 7 * 24 * time.Hour
	TTLVerification = 1 * time.Hour
	TTLRateLimit    = 1 * time.Hour
	TTLOTP          = 5 * time.Minute
	TTLShortCache   = 5 * time.Minute
	TTLMediumCache  = 30 * time.Minute
	TTLLongCache    = 2 * time.Hour
	TTLDayCache     = 24 * time.Hour
)

// ============================================================================
// User Cache Keys
// ============================================================================

// UserKey returns cache key for user data
func UserKey(userID int) string {
	return fmt.Sprintf("%s:%d", PrefixUser, userID)
}

// UserProfileKey returns cache key for user profile
func UserProfileKey(userID int) string {
	return fmt.Sprintf("%s:profile:%d", PrefixUser, userID)
}

// UserByEmailKey returns cache key for user lookup by email
func UserByEmailKey(email string) string {
	return fmt.Sprintf("%s:email:%s", PrefixUser, email)
}

// UserPreferencesKey returns cache key for user preferences
func UserPreferencesKey(userID int) string {
	return fmt.Sprintf("%s:preferences:%d", PrefixUser, userID)
}

// ============================================================================
// Session Cache Keys
// ============================================================================

// SessionKey returns cache key for session
func SessionKey(sessionID string) string {
	return fmt.Sprintf("%s:%s", PrefixSession, sessionID)
}

// UserSessionsKey returns cache key for user's sessions list
func UserSessionsKey(userID int) string {
	return fmt.Sprintf("%s:user:%d", PrefixSession, userID)
}

// DeviceSessionKey returns cache key for device session
func DeviceSessionKey(deviceID string) string {
	return fmt.Sprintf("%s:device:%s", PrefixSession, deviceID)
}

// ============================================================================
// Token Cache Keys
// ============================================================================

// AccessTokenKey returns cache key for access token
func AccessTokenKey(tokenID string) string {
	return fmt.Sprintf("%s:access:%s", PrefixToken, tokenID)
}

// RefreshTokenKey returns cache key for refresh token
func RefreshTokenKey(tokenID string) string {
	return fmt.Sprintf("%s:refresh:%s", PrefixToken, tokenID)
}

// ResetTokenKey returns cache key for password reset token
func ResetTokenKey(token string) string {
	return fmt.Sprintf("%s:reset:%s", PrefixToken, token)
}

// RevokedTokenKey returns cache key for revoked token
func RevokedTokenKey(tokenID string) string {
	return fmt.Sprintf("%s:revoked:%s", PrefixToken, tokenID)
}

// ============================================================================
// Verification Cache Keys
// ============================================================================

// EmailVerificationKey returns cache key for email verification
func EmailVerificationKey(token string) string {
	return fmt.Sprintf("%s:email:%s", PrefixVerification, token)
}

// PhoneVerificationKey returns cache key for phone verification
func PhoneVerificationKey(phone string) string {
	return fmt.Sprintf("%s:phone:%s", PrefixVerification, phone)
}

// ============================================================================
// Rate Limiting Cache Keys
// ============================================================================

// RateLimitKey returns cache key for rate limiting
func RateLimitKey(identifier, action string) string {
	return fmt.Sprintf("%s:%s:%s", PrefixRateLimit, action, identifier)
}

// LoginAttemptsKey returns cache key for login attempts tracking
func LoginAttemptsKey(identifier string) string {
	return fmt.Sprintf("%s:login:%s", PrefixRateLimit, identifier)
}

// PasswordResetAttemptsKey returns cache key for password reset attempts
func PasswordResetAttemptsKey(identifier string) string {
	return fmt.Sprintf("%s:password_reset:%s", PrefixRateLimit, identifier)
}

// APIRateLimitKey returns cache key for API rate limiting
func APIRateLimitKey(userID int, endpoint string) string {
	return fmt.Sprintf("%s:api:%s:%d", PrefixRateLimit, endpoint, userID)
}

// ============================================================================
// OTP Cache Keys
// ============================================================================

// OTPKey returns cache key for OTP
func OTPKey(identifier, purpose string) string {
	return fmt.Sprintf("%s:%s:%s", PrefixOTP, purpose, identifier)
}

// TwoFactorCodeKey returns cache key for 2FA code
func TwoFactorCodeKey(userID int) string {
	return fmt.Sprintf("%s:2fa:%d", PrefixOTP, userID)
}

// ============================================================================
// Appointment Cache Keys
// ============================================================================

// AppointmentKey returns cache key for appointment
func AppointmentKey(appointmentID int) string {
	return fmt.Sprintf("%s:%d", PrefixAppointment, appointmentID)
}

// UserAppointmentsKey returns cache key for user's appointments
func UserAppointmentsKey(userID int) string {
	return fmt.Sprintf("%s:user:%d", PrefixAppointment, userID)
}

// DoctorAppointmentsKey returns cache key for doctor's appointments
func DoctorAppointmentsKey(doctorID int, date string) string {
	return fmt.Sprintf("%s:doctor:%d:%s", PrefixAppointment, doctorID, date)
}

// AppointmentsByDateKey returns cache key for appointments by date
func AppointmentsByDateKey(date string) string {
	return fmt.Sprintf("%s:date:%s", PrefixAppointment, date)
}

// ============================================================================
// Doctor Cache Keys
// ============================================================================

// DoctorKey returns cache key for doctor
func DoctorKey(doctorID int) string {
	return fmt.Sprintf("%s:%d", PrefixDoctor, doctorID)
}

// DoctorScheduleKey returns cache key for doctor schedule
func DoctorScheduleKey(doctorID int, date string) string {
	return fmt.Sprintf("%s:schedule:%d:%s", PrefixDoctor, doctorID, date)
}

// DoctorAvailabilityKey returns cache key for doctor availability
func DoctorAvailabilityKey(doctorID int) string {
	return fmt.Sprintf("%s:availability:%d", PrefixDoctor, doctorID)
}

// DoctorsBySpecialtyKey returns cache key for doctors by specialty
func DoctorsBySpecialtyKey(specialty string) string {
	return fmt.Sprintf("%s:specialty:%s", PrefixDoctor, specialty)
}

// ============================================================================
// Patient Cache Keys
// ============================================================================

// PatientKey returns cache key for patient
func PatientKey(patientID int) string {
	return fmt.Sprintf("%s:%d", PrefixPatient, patientID)
}

// PatientMedicalHistoryKey returns cache key for patient medical history
func PatientMedicalHistoryKey(patientID int) string {
	return fmt.Sprintf("%s:history:%d", PrefixPatient, patientID)
}

// ============================================================================
// Notification Cache Keys
// ============================================================================

// NotificationKey returns cache key for notification
func NotificationKey(notificationID int) string {
	return fmt.Sprintf("%s:%d", PrefixNotification, notificationID)
}

// UserNotificationsKey returns cache key for user's notifications
func UserNotificationsKey(userID int) string {
	return fmt.Sprintf("%s:user:%d", PrefixNotification, userID)
}

// UnreadNotificationsCountKey returns cache key for unread notifications count
func UnreadNotificationsCountKey(userID int) string {
	return fmt.Sprintf("%s:unread:%d", PrefixNotification, userID)
}

// ============================================================================
// General Cache Keys
// ============================================================================

// CacheKeyWithPrefix returns a cache key with custom prefix
func CacheKeyWithPrefix(prefix string, identifier string) string {
	return fmt.Sprintf("%s:%s:%s", PrefixCache, prefix, identifier)
}

// ListCacheKey returns cache key for list/collection
func ListCacheKey(entity string, filters ...string) string {
	key := fmt.Sprintf("%s:list:%s", PrefixCache, entity)
	for _, filter := range filters {
		key += ":" + filter
	}
	return key
}

// CountCacheKey returns cache key for count
func CountCacheKey(entity string) string {
	return fmt.Sprintf("%s:count:%s", PrefixCache, entity)
}

// ============================================================================
// Pattern Helpers
// ============================================================================

// UserPattern returns pattern for all user keys
func UserPattern() string {
	return fmt.Sprintf("%s:*", PrefixUser)
}

// SessionPattern returns pattern for all session keys
func SessionPattern() string {
	return fmt.Sprintf("%s:*", PrefixSession)
}

// UserSessionPattern returns pattern for specific user's session keys
func UserSessionPattern(userID int) string {
	return fmt.Sprintf("%s:user:%d:*", PrefixSession, userID)
}

// RateLimitPattern returns pattern for all rate limit keys
func RateLimitPattern() string {
	return fmt.Sprintf("%s:*", PrefixRateLimit)
}

// OTPPattern returns pattern for all OTP keys
func OTPPattern() string {
	return fmt.Sprintf("%s:*", PrefixOTP)
}

// ============================================================================
// Cache Invalidation Helpers
// ============================================================================

// InvalidateUser invalidates all user-related cache
func InvalidateUser(userID int) error {
	keys := []string{
		UserKey(userID),
		UserProfileKey(userID),
		UserPreferencesKey(userID),
		UserAppointmentsKey(userID),
		UserNotificationsKey(userID),
		UnreadNotificationsCountKey(userID),
	}
	return DeleteMultiple(keys...)
}

// InvalidateSession invalidates session cache
func InvalidateSession(sessionID string) error {
	return Delete(SessionKey(sessionID))
}

// InvalidateUserSessions invalidates all sessions for a user
func InvalidateUserSessions(userID int) error {
	return DeletePattern(UserSessionPattern(userID))
}

// InvalidateDoctor invalidates doctor-related cache
func InvalidateDoctor(doctorID int) error {
	// Delete specific keys
	keys := []string{
		DoctorKey(doctorID),
		DoctorAvailabilityKey(doctorID),
	}
	if err := DeleteMultiple(keys...); err != nil {
		return err
	}

	// Delete schedule patterns
	pattern := fmt.Sprintf("%s:schedule:%d:*", PrefixDoctor, doctorID)
	return DeletePattern(pattern)
}

// InvalidateAppointment invalidates appointment-related cache
func InvalidateAppointment(appointmentID int) error {
	return Delete(AppointmentKey(appointmentID))
}

// ============================================================================
// Helper Functions with TTL
// ============================================================================

// SetWithDefaultTTL sets a value with default TTL based on prefix
func SetWithDefaultTTL(key string, value interface{}) error {
	ttl := TTLMediumCache // Default

	// Determine TTL based on key prefix
	switch {
	case len(key) > len(PrefixSession) && key[:len(PrefixSession)] == PrefixSession:
		ttl = TTLSession
	case len(key) > len(PrefixToken) && key[:len(PrefixToken)] == PrefixToken:
		ttl = TTLToken
	case len(key) > len(PrefixVerification) && key[:len(PrefixVerification)] == PrefixVerification:
		ttl = TTLVerification
	case len(key) > len(PrefixRateLimit) && key[:len(PrefixRateLimit)] == PrefixRateLimit:
		ttl = TTLRateLimit
	case len(key) > len(PrefixOTP) && key[:len(PrefixOTP)] == PrefixOTP:
		ttl = TTLOTP
	}

	return SetJSON(key, value, ttl)
}
