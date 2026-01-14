package user

import "time"

// ============================================================================
// AUDIT LOGS
// ============================================================================
// Audit logs are domain events that happened
// They belong to domain layer, not infrastructure

// LoginActivity represents login activity log
type LoginActivity struct {
	ID        int       `json:"id" db:"id"`
	UserID    *int      `json:"user_id,omitempty" db:"user_id"` // NULL for non-existent users
	Email     string    `json:"email" db:"email"`
	Success   bool      `json:"success" db:"success"`
	IPAddress string    `json:"ip_address" db:"ip_address"`
	UserAgent string    `json:"-" db:"user_agent"`
	Location  *string   `json:"location,omitempty" db:"location"`
	Reason    *string   `json:"reason,omitempty" db:"reason"`
	SessionID *int      `json:"session_id,omitempty" db:"session_id"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

// UserActivityLog represents general user activity audit log
type UserActivityLog struct {
	ID        int       `json:"id" db:"id"`
	UserID    int       `json:"user_id" db:"user_id"`
	Action    string    `json:"action" db:"action"` // "profile_update", "password_change"
	Entity    string    `json:"entity" db:"entity"` // "user", "appointment"
	EntityID  *int      `json:"entity_id,omitempty" db:"entity_id"`
	Changes   *string   `json:"changes,omitempty" db:"changes"` // JSON
	IPAddress string    `json:"ip_address" db:"ip_address"`
	UserAgent *string   `json:"-" db:"user_agent"`
	SessionID *int      `json:"session_id,omitempty" db:"session_id"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

// SecurityEvent represents security-related events
type SecurityEvent struct {
	ID          int        `json:"id" db:"id"`
	UserID      *int       `json:"user_id,omitempty" db:"user_id"`
	EventType   string     `json:"event_type" db:"event_type"` // "suspicious_login", "2fa_enabled"
	Severity    string     `json:"severity" db:"severity"`     // "low", "medium", "high", "critical"
	Description string     `json:"description" db:"description"`
	IPAddress   string     `json:"ip_address" db:"ip_address"`
	UserAgent   *string    `json:"-" db:"user_agent"`
	Metadata    *string    `json:"metadata,omitempty" db:"metadata"` // JSON
	Resolved    bool       `json:"resolved" db:"resolved"`
	ResolvedAt  *time.Time `json:"resolved_at,omitempty" db:"resolved_at"`
	ResolvedBy  *int       `json:"resolved_by,omitempty" db:"resolved_by"`
	CreatedAt   time.Time  `json:"created_at" db:"created_at"`
}

// MarkAsResolved marks security event as resolved
func (e *SecurityEvent) MarkAsResolved(resolvedBy int) {
	e.Resolved = true
	now := time.Now()
	e.ResolvedAt = &now
	e.ResolvedBy = &resolvedBy
}
