package user

import "time"

// ============================================================================
// AUTHENTICATION MODELS
// ============================================================================

type AuthProvider string

const (
	ProviderGoogle   AuthProvider = "google"
	ProviderFacebook AuthProvider = "facebook"
	ProviderApple    AuthProvider = "apple"
)

// ============================================================================
// USER CREDENTIALS
// ============================================================================

// UserCredential stores authentication credentials
type UserCredential struct {
	ID           int     `json:"id" db:"id"`
	UserID       int     `json:"user_id" db:"user_id"`
	PasswordHash *string `json:"-" db:"password_hash"`

	// Two-factor authentication
	TwoFactorEnabled bool    `json:"two_factor_enabled" db:"two_factor_enabled"`
	TwoFactorSecret  *string `json:"-" db:"two_factor_secret"` // Encrypted

	AuditFields
}

// Has2FAEnabled checks if 2FA is enabled
func (c *UserCredential) Has2FAEnabled() bool {
	return c.TwoFactorEnabled && c.TwoFactorSecret != nil
}

// ============================================================================
// PASSWORD MANAGEMENT
// ============================================================================

// PasswordHistory represents password history (prevent reuse)
type PasswordHistory struct {
	ID           int       `json:"id" db:"id"`
	UserID       int       `json:"user_id" db:"user_id"`
	PasswordHash string    `json:"-" db:"password_hash"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
}

// PasswordResetToken represents password reset token
type PasswordResetToken struct {
	ID        int        `json:"id" db:"id"`
	UserID    int        `json:"user_id" db:"user_id"`
	Token     string     `json:"token" db:"token"`
	ExpiresAt time.Time  `json:"expires_at" db:"expires_at"`
	UsedAt    *time.Time `json:"used_at,omitempty" db:"used_at"`

	// Security audit
	IPAddress       *string `json:"-" db:"ip_address"`
	UserAgent       *string `json:"-" db:"user_agent"`
	UsedBySessionID *int    `json:"-" db:"used_by_session_id"`

	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

// IsExpired checks if token is expired
func (t *PasswordResetToken) IsExpired() bool {
	return time.Now().After(t.ExpiresAt)
}

// IsUsed checks if token has been used
func (t *PasswordResetToken) IsUsed() bool {
	return t.UsedAt != nil
}

// IsValid checks if token is valid
func (t *PasswordResetToken) IsValid() bool {
	return !t.IsExpired() && !t.IsUsed()
}

// MarkAsUsed marks token as used
func (t *PasswordResetToken) MarkAsUsed(sessionID *int) {
	now := time.Now()
	t.UsedAt = &now
	t.UsedBySessionID = sessionID
}

// ============================================================================
// EMAIL VERIFICATION
// ============================================================================

// EmailVerificationToken represents email verification token
type EmailVerificationToken struct {
	ID        int        `json:"id" db:"id"`
	UserID    int        `json:"user_id" db:"user_id"`
	Token     string     `json:"token" db:"token"`
	ExpiresAt time.Time  `json:"expires_at" db:"expires_at"`
	UsedAt    *time.Time `json:"used_at,omitempty" db:"used_at"`

	// Security audit
	IPAddress       *string `json:"-" db:"ip_address"`
	UserAgent       *string `json:"-" db:"user_agent"`
	UsedBySessionID *int    `json:"-" db:"used_by_session_id"`

	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

// IsExpired checks if token is expired
func (t *EmailVerificationToken) IsExpired() bool {
	return time.Now().After(t.ExpiresAt)
}

// IsUsed checks if token has been used
func (t *EmailVerificationToken) IsUsed() bool {
	return t.UsedAt != nil
}

// IsValid checks if token is valid
func (t *EmailVerificationToken) IsValid() bool {
	return !t.IsExpired() && !t.IsUsed()
}

// MarkAsUsed marks token as used
func (t *EmailVerificationToken) MarkAsUsed(sessionID *int) {
	now := time.Now()
	t.UsedAt = &now
	t.UsedBySessionID = sessionID
}

// ============================================================================
// SOCIAL AUTHENTICATION
// ============================================================================

// UserSocialAuth represents social authentication
type UserSocialAuth struct {
	ID             int          `json:"id" db:"id"`
	UserID         int          `json:"user_id" db:"user_id"`
	Provider       AuthProvider `json:"provider" db:"provider"`
	ProviderUserID string       `json:"provider_user_id" db:"provider_user_id"`

	// OAuth tokens (encrypted in application layer)
	AccessToken    *string    `json:"-" db:"access_token"`
	RefreshToken   *string    `json:"-" db:"refresh_token"`
	TokenExpiresAt *time.Time `json:"-" db:"token_expires_at"`

	// Cached profile
	ProviderEmail     *string `json:"provider_email" db:"provider_email"`
	ProviderName      *string `json:"provider_name" db:"provider_name"`
	ProviderAvatarURL *string `json:"provider_avatar_url" db:"provider_avatar_url"`

	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

// IsTokenExpired checks if OAuth token is expired
func (s *UserSocialAuth) IsTokenExpired() bool {
	if s.TokenExpiresAt == nil {
		return false
	}
	return time.Now().After(*s.TokenExpiresAt)
}

// ============================================================================
// TWO-FACTOR AUTHENTICATION
// ============================================================================

// TwoFactorBackupCode represents 2FA backup codes
type TwoFactorBackupCode struct {
	ID        int        `json:"id" db:"id"`
	UserID    int        `json:"user_id" db:"user_id"`
	Code      string     `json:"-" db:"code"` // Hashed
	UsedAt    *time.Time `json:"used_at,omitempty" db:"used_at"`
	CreatedAt time.Time  `json:"created_at" db:"created_at"`
}

// IsUsed checks if backup code has been used
func (b *TwoFactorBackupCode) IsUsed() bool {
	return b.UsedAt != nil
}

// MarkAsUsed marks backup code as used
func (b *TwoFactorBackupCode) MarkAsUsed() {
	now := time.Now()
	b.UsedAt = &now
}
