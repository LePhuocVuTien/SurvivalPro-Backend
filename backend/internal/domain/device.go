package models

import "time"

// UserDevice represents a user's device
type UserDevice struct {
	ID                  int          `json:"id" db:"id"`
	UserID              int          `json:"user_id" db:"user_id"`
	DeviceUUID          string       `json:"device_uuid" db:"device_uuid"` // UUID/FCM/APNs token
	DeviceName          *string      `json:"device_name,omitempty" db:"device_name"`
	DeviceType          *string      `json:"device_type,omitempty" db:"device_type"`
	Platform            PlatformType `json:"platform_type" db:"plat_form"`
	OSVersion           *string      `json:"os_version,omitempty" db:"os_version"`
	AppVersion          *string      `json:"app_version,omitempty" db:"app_version"`
	IsActive            bool         `json:"is_active" db:"is_active"`
	LastActiveAt        *time.Time   `json:"last_active_at,omitempty" db:"last_active_at"`     // Nullable cho device mới
	LastIPAddress       *string      `json:"last_ip_address,omitempty" db:"last_ip_address"`   // Validate ở service layer
	LastLocationCountry *string      `json:"last_location_country" db:"last_location_country"` // ISO country code (VN, US, etc)
	CreatedAt           time.Time    `json:"created_at" db:"created_at"`
	UpdatedAt           time.Time    `json:"updated_at" db:"updated_at"`
	DeletedAt           *time.Time   `json:"-" db:"deleted_at"` // Soft delete
}

// UserDeviceCreate represents device registration request
type UserDeviceCreate struct {
	DeviceUUID string       `json:"device_uuid" binding:"required,max=255"`
	DeviceName *string      `json:"device_name,omitempty" binding:"omitempty,max=100"`
	DeviceType *string      `json:"device_type,omitempty" binding:"omitempty,max=50"`
	Platform   PlatformType `json:"plat_form" binding:"required,oneof=ios android web"`
	OSVersion  *string      `json:"os_version,omitempty" binding:"omitempty,max=50"`
	AppVersion *string      `json:"app_version,omitempty"`
}

// UserSession represents a user session
type UserSession struct {
	ID               int        `json:"id" db:"id"`
	UserID           int        `json:"user_id" db:"user_id"`
	UserDeviceID     int        `json:"user_device_id" db:"user_device_id"`
	SessionToken     string     `json:"-" db:"session_token"`
	SessionTokenHash string     `json:"-" db:"session_token_hash"`
	RefreshToken     *string    `json:"-" db:"refresh_token"`
	RefreshTokenHash *string    `json:"-" db:"refresh_token_hash"`
	IPAddress        *string    `json:"ip_address,omitempty" db:"ip_address"`
	UserAgent        *string    `json:"user_agent,omitempty" db:"user_agent"`
	ExpiresAt        time.Time  `json:"expires_at" db:"expires_at"`
	CreatedAt        time.Time  `json:"created_at" db:"created_at"`
	LastActivityAt   time.Time  `json:"last_activity_at" db:"last_activity_at"`
	RevokedAt        *time.Time `json:"revoked_at,omitempty" db:"revoked_at"`
	RevokedReason    *string    `json:"revoked_reason,omitempty" db:"revoked_reason"`
}

// UserSessionCreate represents session creation request
type UserSessionCreate struct {
	UserDeviceID int     `json:"user_devive_id" binding:"request"`
	IPAddress    *string `json:"ip_address,omitempty" binding:"omitempty,ip"` // Validate IP format
	UserAgent    *string `json:"user_agent,omitempty" binding:"omitempty,max=500"`
}

// UserSessionResponse represents session with device info
type UserSessionResponse struct {
	ID              int          `json:"id"`
	UserID          int          `json:"user_id"`
	UserDeviceID    int          `json:"user_device_id"`
	DeviceName      *string      `json:"device_name,omitempty"`
	DeviceType      *string      `json:"device_type,omitempty"`
	Platform        PlatformType `json:"platform"`
	IPAddress       *string      `json:"ip_address,omitempty"`
	LocationCountry *string      `json:"location_country,omitempty"`
	IsCurrent       bool         `json:"is_current"`
	ExpiresAt       time.Time    `json:"expires_at"`
	CreatedAt       time.Time    `json:"created_at"`
	LastActivityAt  time.Time    `json:"last_activity_at"`
}

// UserSessionListResponse represents paginated session list
type UserSessionListResponse struct {
	Session []UserSessionResponse `json:"sessions"`
	Total   int                   `json:"total"`
}

// RefreshTokenRequest represents token refresh request
type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required,min=32"`
}

// RefreshTokenReponse represents token refresh response
type RefreshTokenReponse struct {
	AccessToken  string    `json:"access_token"`
	RefreshToken string    `json:"refresh_token"`
	ExpiresAt    time.Time `json:"expires_at"`
}

func (s *UserSession) IsExpired() bool {
	return time.Now().After(s.ExpiresAt)
}

func (s *UserSession) IsRevoked() bool {
	return s.RevokedAt != nil
}

func (d *UserDevice) IsRecentlyActive(days int) bool {
	if d.LastActiveAt == nil {
		return false
	}
	return time.Since(*d.LastActiveAt) < time.Duration(days)*24*time.Hour
}

func (d *UserDevice) IsDeleted() bool {
	return d.DeletedAt != nil
}
