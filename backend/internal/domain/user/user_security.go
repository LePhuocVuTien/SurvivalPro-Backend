package user

import "time"

// ============================================================================
// USER SESSION
// ============================================================================

// UserSession represents user session with device tracking
type UserSession struct {
	ID           int    `json:"id" db:"id"`
	UserID       int    `json:"user_id" db:"user_id"`
	RefreshToken string `json:"-" db:"refresh_token"` // Hashed

	// Device info
	DeviceID   *string `json:"device_id,omitempty" db:"device_id"`
	DeviceName *string `json:"device_name,omitempty" db:"device_name"`
	Platform   *string `json:"platform,omitempty" db:"platform"`
	AppVersion *string `json:"app_version,omitempty" db:"app_version"`

	// Security
	IPAddress *string `json:"ip_address,omitempty" db:"ip_address"`
	UserAgent *string `json:"-" db:"user_agent"`
	Location  *string `json:"location,omitempty" db:"location"`

	// Lifecycle
	ExpiresAt  time.Time  `json:"expires_at" db:"expires_at"`
	LastUsedAt time.Time  `json:"last_used_at" db:"last_used_at"`
	RevokedAt  *time.Time `json:"revoked_at,omitempty" db:"revoked_at"`
	RevokedBy  *int       `json:"-" db:"revoked_by"`

	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

// IsExpired checks if session is expired
func (s *UserSession) IsExpired() bool {
	return time.Now().After(s.ExpiresAt)
}

// IsRevoked checks if session is revoked
func (s *UserSession) IsRevoked() bool {
	return s.RevokedAt != nil
}

// IsValid checks if session is valid
func (s *UserSession) IsValid() bool {
	return !s.IsExpired() && !s.IsRevoked()
}

// UpdateLastUsed updates the last used timestamp
func (s *UserSession) UpdateLastUsed() {
	s.LastUsedAt = time.Now()
}

// Revoke revokes the session
func (s *UserSession) Revoke(revokedBy *int) {
	now := time.Now()
	s.RevokedAt = &now
	s.RevokedBy = revokedBy
}

// ============================================================================
// PUSH NOTIFICATIONS
// ============================================================================

// UserPushToken represents push notification token
// One user can have multiple devices/tokens
type UserPushToken struct {
	ID       int    `json:"id" db:"id"`
	UserID   int    `json:"user_id" db:"user_id"`
	Token    string `json:"token" db:"token"`
	DeviceID string `json:"device_id" db:"device_id"`
	Platform string `json:"platform" db:"platform"`

	// Token metadata
	IsActive      bool       `json:"is_active" db:"is_active"`
	LastUsedAt    time.Time  `json:"last_used_at" db:"last_used_at"`
	DeactivatedAt *time.Time `json:"deactivated_at,omitempty" db:"deactivated_at"`

	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

// Deactivate marks the token as inactive
func (t *UserPushToken) Deactivate() {
	t.IsActive = false
	now := time.Now()
	t.DeactivatedAt = &now
}

// Activate marks the token as active
func (t *UserPushToken) Activate() {
	t.IsActive = true
	t.DeactivatedAt = nil
	t.LastUsedAt = time.Now()
}

// ============================================================================
// SECURITY INFO
// ============================================================================

// UserSecurityInfo stores security-related information
type UserSecurityInfo struct {
	ID                  int        `json:"id" db:"id"`
	UserID              int        `json:"user_id" db:"user_id"`
	FailedLoginAttempts int        `json:"failed_login_attempts" db:"failed_login_attempts"`
	LastFailedLoginAt   *time.Time `json:"last_failed_login_at,omitempty" db:"last_failed_login_at"`
	LockedUntil         *time.Time `json:"locked_until,omitempty" db:"locked_until"`
	LastPasswordChange  *time.Time `json:"last_password_change,omitempty" db:"last_password_change"`
	LastLoginAt         *time.Time `json:"last_login_at,omitempty" db:"last_login_at"`

	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

// IsLocked checks if user account is locked
func (s *UserSecurityInfo) IsLocked() bool {
	if s.LockedUntil == nil {
		return false
	}
	return time.Now().Before(*s.LockedUntil)
}

// IncrementFailedAttempts increments failed login attempts
func (s *UserSecurityInfo) IncrementFailedAttempts() {
	s.FailedLoginAttempts++
	now := time.Now()
	s.LastFailedLoginAt = &now

	// Lock account after 5 failed attempts for 30 minutes
	if s.FailedLoginAttempts >= 5 {
		lockUntil := now.Add(30 * time.Minute)
		s.LockedUntil = &lockUntil
	}
}

// ResetFailedAttempts resets failed login attempts counter
func (s *UserSecurityInfo) ResetFailedAttempts() {
	s.FailedLoginAttempts = 0
	s.LastFailedLoginAt = nil
	s.LockedUntil = nil
}

// UpdateLastLogin updates last login timestamp
func (s *UserSecurityInfo) UpdateLastLogin() {
	now := time.Now()
	s.LastLoginAt = &now
}

// UpdatePasswordChange updates last password change timestamp
func (s *UserSecurityInfo) UpdatePasswordChange() {
	now := time.Now()
	s.LastPasswordChange = &now
}
