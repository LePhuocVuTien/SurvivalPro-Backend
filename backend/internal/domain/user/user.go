package user

import "time"

// ============================================================================
// ENUMS
// ============================================================================

type UserRole string

const (
	UserRoleAdmin  UserRole = "admin"
	UserRoleLeader UserRole = "leader"
	UserRoleUser   UserRole = "user"
)

type AccountStatus string

const (
	AccountPending   AccountStatus = "pending"
	AccountActive    AccountStatus = "active"
	AccountSuspended AccountStatus = "suspended"
	AccountBanned    AccountStatus = "banned"
	AccountClosed    AccountStatus = "closed"
)

// ============================================================================
// AUDIT FIELDS (Embedded)
// ============================================================================

// AuditFields contains common audit fields
type AuditFields struct {
	CreatedAt time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt time.Time  `json:"updated_at" db:"updated_at"`
	CreatedBy *int       `json:"-" db:"created_by"`
	UpdatedBy *int       `json:"-" db:"updated_by"`
	DeletedAt *time.Time `json:"-" db:"deleted_at"`
	DeletedBy *int       `json:"-" db:"deleted_by"`
}

// IsDeleted checks if entity is soft deleted
func (a *AuditFields) IsDeleted() bool {
	return a.DeletedAt != nil
}

// ============================================================================
// USER ENTITY (Domain Model)
// ============================================================================

// User represents the core user entity
type User struct {
	ID    int      `json:"id" db:"id"`
	Email string   `json:"email" db:"email"`
	Name  string   `json:"name" db:"name"`
	Role  UserRole `json:"role" db:"role"`

	// Account status
	AccountStatus   AccountStatus `json:"account_status" db:"account_status"`
	StatusChangedAt *time.Time    `json:"status_changed_at,omitempty" db:"status_changed_at"`
	StatusChangedBy *int          `json:"-" db:"status_changed_by"`
	StatusReason    *string       `json:"-" db:"status_reason"`

	// Profile
	Phone     *string `json:"phone,omitempty" db:"phone"`
	AvatarURL *string `json:"avatar_url,omitempty" db:"avatar_url"`

	// Emergency info
	BloodType             *string `json:"blood_type,omitempty" db:"blood_type"`
	Allergies             *string `json:"allergies,omitempty" db:"allergies"`
	EmergencyContactName  *string `json:"emergency_contact_name,omitempty" db:"emergency_contact_name"`
	EmergencyContactPhone *string `json:"emergency_contact_phone,omitempty" db:"emergency_contact_phone"`

	// Verification status
	EmailVerified bool `json:"email_verified" db:"email_verified"`
	PhoneVerified bool `json:"phone_verified" db:"phone_verified"`

	// Audit fields
	AuditFields
}

// ============================================================================
// DOMAIN METHODS (State queries)
// ============================================================================

// IsActive checks if account is active and can be used
func (u *User) IsActive() bool {
	return u.AccountStatus == AccountActive && !u.IsDeleted()
}

// IsVerified checks if user has verified their email
func (u *User) IsVerified() bool {
	return u.EmailVerified
}

// HasRole checks if user has specific role
func (u *User) HasRole(role UserRole) bool {
	return u.Role == role
}

// IsAdmin checks if user is admin
func (u *User) IsAdmin() bool {
	return u.Role == UserRoleAdmin
}

// IsLeader checks if user is doctor
func (u *User) IsLeader() bool {
	return u.Role == UserRoleLeader
}

// IsUser checks if user has user role
func (u *User) IsUser() bool {
	return u.Role == UserRoleUser
}

// IsSuspended checks if account is suspended
func (u *User) IsSuspended() bool {
	return u.AccountStatus == AccountSuspended
}

// IsBanned checks if account is banned
func (u *User) IsBanned() bool {
	return u.AccountStatus == AccountBanned
}

// IsClosed checks if account is closed
func (u *User) IsClosed() bool {
	return u.AccountStatus == AccountClosed
}

// IsPending checks if account is pending
func (u *User) IsPending() bool {
	return u.AccountStatus == AccountPending
}

// CanTransitionTo checks if user can transition to new status
// This is a convenience method that wraps the policy function
func (u *User) CanTransitionTo(newStatus AccountStatus) bool {
	return isValidStatusTransition(u.AccountStatus, newStatus)
}

// ============================================================================
// ACCOUNT STATUS CHANGE LOG
// ============================================================================

// AccountStatusChange represents account status change audit log
type AccountStatusChange struct {
	ID        int           `json:"id" db:"id"`
	UserID    int           `json:"user_id" db:"user_id"`
	OldStatus AccountStatus `json:"old_status" db:"old_status"`
	NewStatus AccountStatus `json:"new_status" db:"new_status"`
	Reason    *string       `json:"reason,omitempty" db:"reason"`
	ChangedBy int           `json:"changed_by" db:"changed_by"`
	CreatedAt time.Time     `json:"created_at" db:"created_at"`
}
