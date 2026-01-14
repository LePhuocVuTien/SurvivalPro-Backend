package user

import "time"

// ============================================================================
// USER DTOs
// ============================================================================
// DTOs define API contracts - they can be different from domain models

// UserCreateRequest represents user creation request
type UserCreateRequest struct {
	Email    string   `json:"email"`
	Password string   `json:"password"`
	Name     string   `json:"name"`
	Phone    *string  `json:"phone,omitempty"`
	Role     UserRole `json:"role,omitempty"`
}

// UserUpdateRequest represents user update request
type UserUpdateRequest struct {
	Name                  *string `json:"name,omitempty"`
	Phone                 *string `json:"phone,omitempty"`
	AvatarURL             *string `json:"avatar_url,omitempty"`
	BloodType             *string `json:"blood_type,omitempty"`
	Allergies             *string `json:"allergies,omitempty"`
	EmergencyContactName  *string `json:"emergency_contact_name,omitempty"`
	EmergencyContactPhone *string `json:"emergency_contact_phone,omitempty"`
}

// UserResponse represents user response (without sensitive data)
type UserResponse struct {
	ID                    int           `json:"id"`
	Email                 string        `json:"email"`
	Name                  string        `json:"name"`
	Role                  UserRole      `json:"role"`
	AccountStatus         AccountStatus `json:"account_status"`
	Phone                 *string       `json:"phone,omitempty"`
	AvatarURL             *string       `json:"avatar_url,omitempty"`
	BloodType             *string       `json:"blood_type,omitempty"`
	Allergies             *string       `json:"allergies,omitempty"`
	EmergencyContactName  *string       `json:"emergency_contact_name,omitempty"`
	EmergencyContactPhone *string       `json:"emergency_contact_phone,omitempty"`
	EmailVerified         bool          `json:"email_verified"`
	PhoneVerified         bool          `json:"phone_verified"`
	TwoFactorEnabled      bool          `json:"two_factor_enabled"`
	CreatedAt             time.Time     `json:"created_at"`
	LastLoginAt           *time.Time    `json:"last_login_at,omitempty"`
}

// ToResponse converts User to UserResponse
func (u *User) ToResponse(securityInfo *UserSecurityInfo, credential *UserCredential) *UserResponse {
	resp := &UserResponse{
		ID:                    u.ID,
		Email:                 u.Email,
		Name:                  u.Name,
		Role:                  u.Role,
		AccountStatus:         u.AccountStatus,
		Phone:                 u.Phone,
		AvatarURL:             u.AvatarURL,
		BloodType:             u.BloodType,
		Allergies:             u.Allergies,
		EmergencyContactName:  u.EmergencyContactName,
		EmergencyContactPhone: u.EmergencyContactPhone,
		EmailVerified:         u.EmailVerified,
		PhoneVerified:         u.PhoneVerified,
		CreatedAt:             u.CreatedAt,
	}

	if securityInfo != nil {
		resp.LastLoginAt = securityInfo.LastLoginAt
	}

	if credential != nil {
		resp.TwoFactorEnabled = credential.TwoFactorEnabled
	}

	return resp
}

// UserListResponse represents paginated user list
type UserListResponse struct {
	Users      []*UserResponse `json:"users"`
	Total      int64           `json:"total"`
	Page       int             `json:"page"`
	PageSize   int             `json:"page_size"`
	TotalPages int             `json:"total_pages"`
}

// ============================================================================
// AUTHENTICATION DTOs
// ============================================================================

// LoginRequest represents login request
type LoginRequest struct {
	Email         string  `json:"email"`
	Password      string  `json:"password"`
	DeviceID      *string `json:"device_id,omitempty"`
	DeviceName    *string `json:"device_name,omitempty"`
	Platform      *string `json:"platform,omitempty"`
	TwoFactorCode *string `json:"two_factor_code,omitempty"`
}

// LoginResponse represents login response
type LoginResponse struct {
	User         *UserResponse `json:"user"`
	AccessToken  string        `json:"access_token"`
	RefreshToken string        `json:"refresh_token"`
	SessionID    int           `json:"session_id"` // For debugging/revoke
	ExpiresAt    time.Time     `json:"expires_at"`
}

// RefreshTokenRequest represents refresh token request
type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token"`
}

// RefreshTokenResponse represents refresh token response
type RefreshTokenResponse struct {
	AccessToken  string    `json:"access_token"`
	RefreshToken string    `json:"refresh_token"`
	SessionID    int       `json:"session_id"`
	ExpiresAt    time.Time `json:"expires_at"`
}

// LogoutRequest represents logout request
type LogoutRequest struct {
	RefreshToken string `json:"refresh_token"`
}

// LogoutAllRequest represents logout all devices request
type LogoutAllRequest struct {
	Password string `json:"password"`
}

// SocialLoginRequest represents social login request
type SocialLoginRequest struct {
	Provider       AuthProvider `json:"provider"`
	AccessToken    string       `json:"access_token"`
	ProviderUserID string       `json:"provider_user_id"`
	Email          *string      `json:"email,omitempty"`
	Name           *string      `json:"name,omitempty"`
	AvatarURL      *string      `json:"avatar_url,omitempty"`
	DeviceID       *string      `json:"device_id,omitempty"`
	DeviceName     *string      `json:"device_name,omitempty"`
	Platform       *string      `json:"platform,omitempty"`
}

// ============================================================================
// PASSWORD MANAGEMENT DTOs
// ============================================================================

// ChangePasswordRequest represents change password request
type ChangePasswordRequest struct {
	OldPassword string `json:"old_password"`
	NewPassword string `json:"new_password"`
}

// ForgotPasswordRequest represents forgot password request
type ForgotPasswordRequest struct {
	Email string `json:"email"`
}

// ResetPasswordRequest represents reset password request
type ResetPasswordRequest struct {
	Token       string `json:"token"`
	NewPassword string `json:"new_password"`
}

// ============================================================================
// EMAIL VERIFICATION DTOs
// ============================================================================

// VerifyEmailRequest represents verify email request
type VerifyEmailRequest struct {
	Token string `json:"token"`
}

// ResendVerificationRequest represents resend verification email request
type ResendVerificationRequest struct {
	Email string `json:"email"`
}

// ============================================================================
// SESSION MANAGEMENT DTOs
// ============================================================================

// SessionResponse represents session list item
type SessionResponse struct {
	ID         int       `json:"id"`
	DeviceName string    `json:"device_name"`
	Platform   string    `json:"platform"`
	Location   string    `json:"location"`
	IPAddress  string    `json:"ip_address"`
	LastUsedAt time.Time `json:"last_used_at"`
	CreatedAt  time.Time `json:"created_at"`
	IsCurrent  bool      `json:"is_current"`
}

// RevokeSessionRequest represents revoke session request
type RevokeSessionRequest struct {
	SessionID int `json:"session_id"`
}

// ============================================================================
// ACCOUNT MANAGEMENT DTOs (Admin)
// ============================================================================

// UpdateAccountStatusRequest represents update account status request
type UpdateAccountStatusRequest struct {
	UserID     int           `json:"user_id"`
	Status     AccountStatus `json:"status"`
	Reason     string        `json:"reason"`
	NotifyUser bool          `json:"notify_user"`
}

// ============================================================================
// TWO-FACTOR AUTHENTICATION DTOs
// ============================================================================

// Enable2FARequest represents enable 2FA request
type Enable2FARequest struct {
	Password string `json:"password"`
}

// Enable2FAResponse represents enable 2FA response
type Enable2FAResponse struct {
	Secret      string   `json:"secret"`
	QRCodeURL   string   `json:"qr_code_url"`
	BackupCodes []string `json:"backup_codes"`
}

// Verify2FARequest represents verify 2FA request
type Verify2FARequest struct {
	Code string `json:"code"`
}

// Disable2FARequest represents disable 2FA request
type Disable2FARequest struct {
	Password string `json:"password"`
	Code     string `json:"code"`
}

// ============================================================================
// QUERY & FILTER DTOs
// ============================================================================

// UserFilter represents user query filters
type UserFilter struct {
	Email         *string        `json:"email,omitempty"`
	Role          *UserRole      `json:"role,omitempty"`
	AccountStatus *AccountStatus `json:"account_status,omitempty"`
	EmailVerified *bool          `json:"email_verified,omitempty"`
	Search        *string        `json:"search,omitempty"`
	Page          int            `json:"page"`
	PageSize      int            `json:"page_size"`
	SortBy        string         `json:"sort_by,omitempty"`
	SortOrder     string         `json:"sort_order,omitempty"`
}

// ============================================================================
// PUSH TOKEN DTOs
// ============================================================================

// RegisterPushTokenRequest represents register push token request
type RegisterPushTokenRequest struct {
	Token    string `json:"token"`
	DeviceID string `json:"device_id"`
	Platform string `json:"platform"`
}

// UpdatePushTokenRequest represents update push token request
type UpdatePushTokenRequest struct {
	OldToken string `json:"old_token"`
	NewToken string `json:"new_token"`
}
