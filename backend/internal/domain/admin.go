package models

import (
	"fmt"
	"time"
)

// AdminAuditLogResponse represents audit log with admin details
type AdminAuditLogResponse struct {
	ID           int              `json:"id"`
	AdminID      *int             `json:"admin_id,omitempty"`
	AdminName    *string          `json:"admin_name,omitempty"`
	AdminEmail   *string          `json:"admin_email,omitempty"`
	Action       AdminAction      `json:"action"`
	TargetType   *AuditTargetType `json:"target_type,omitempty"`
	TargetID     *int             `json:"target_id,omitempty"`
	Reason       *string          `json:"reason,omitempty"`
	Metadata     map[string]any   `json:"metadata,omitempty"` // Parsed from JSONB
	Success      bool             `json:"success"`
	ErrorMessage *string          `json:"error_message,omitempty"`
	IPAddress    *string          `json:"ip_address,omitempty"`
	CreatedAt    time.Time        `json:"created_at"`
}

// AdminAuditLogFilter represents filter options for audit log queries
type AdminAuditLogFilter struct {
	AdminID    *int             `json:"admin_id,omitempty" form:"admin_id" binding:"omitempty,min=1"`
	Action     *AdminAction     `json:"action,omitempty" form:"action"`
	TargetType *AuditTargetType `json:"target_type,omitempty" form:"target_type"`
	TargetID   *int             `json:"target_id,omitempty" form:"target_id" binding:"omitempty,min=1"`
	Success    *bool            `json:"success,omitempty" form:"success"`
	StartDate  *time.Time       `json:"start_date,omitempty" form:"start_date"`
	EndDate    *time.Time       `json:"end_date,omitempty" form:"end_date"`
	Limit      int              `json:"limit,omitempty" form:"limit" binding:"omitempty,min=1,max=100"`
	Offset     int              `json:"offset,omitempty" form:"offset" binding:"omitempty,min=0"`
	SortBy     *string          `json:"sort_by,omitempty" form:"sort_by"`       // created_at, action
	SortOrder  *string          `json:"sort_order,omitempty" form:"sort_order"` // asc, desc
}

// Validate validates the filter
func (f *AdminAuditLogFilter) Validate() error {
	if f.StartDate != nil && f.EndDate != nil {
		if f.EndDate.Before(*f.StartDate) {
			return fmt.Errorf("end_date must be after start_date")
		}
	}

	if f.Action != nil && !f.Action.IsValid() {
		return fmt.Errorf("invalid action filter")
	}

	if f.TargetType != nil && !f.TargetType.IsValid() {
		return fmt.Errorf("invalid target type filter")
	}

	if f.SortBy != nil {
		validSortFields := []string{"created_at", "action", "success"}
		valid := false
		for _, field := range validSortFields {
			if *f.SortBy == field {
				valid = true
				break
			}
		}
		if !valid {
			return fmt.Errorf("invalid sort_by field: must be one of created_at, action, success")
		}
	}

	if f.SortOrder != nil {
		if *f.SortOrder != "asc" && *f.SortOrder != "desc" {
			return fmt.Errorf("sort_order must be 'asc' or 'desc'")
		}
	}

	return nil
}

// AdminDashboardStats represents admin dashboard statistics
type AdminDashboardStats struct {
	// User stats
	TotalUsers        int `json:"total_users"`
	ActiveUsers       int `json:"active_users"`
	InactiveUsers     int `json:"inactive_users"`
	BannedUsers       int `json:"banned_users"`
	NewUsersToday     int `json:"new_users_today"`
	NewUsersThisWeek  int `json:"new_users_this_week"`
	NewUsersThisMonth int `json:"new_users_this_month"`

	// Group stats
	TotalGroups  int `json:"total_groups"`
	ActiveGroups int `json:"active_groups"`

	// Subscription stats
	TotalSubscriptions   int `json:"total_subscriptions"`
	ActiveSubscriptions  int `json:"active_subscriptions"`
	TrialSubscriptions   int `json:"trial_subscriptions"`
	ExpiredSubscriptions int `json:"expired_subscriptions"`

	// SOS stats
	TotalSOSEvents    int `json:"total_sos_events"`
	ActiveSOSEvents   int `json:"active_sos_events"`
	ResolvedSOSEvents int `json:"resolved_sos_events"`

	// Danger zone stats
	TotalDangerZones  int `json:"total_danger_zones"`
	ActiveDangerZones int `json:"active_danger_zones"`
	UsersInDanger     int `json:"users_in_danger"`

	// Revenue stats
	MonthlyRevenue float64 `json:"monthly_revenue"`
	YearlyRevenue  float64 `json:"yearly_revenue"`
	TotalRevenue   float64 `json:"total_revenue"`

	// System stats
	AverageSessionDurationMinutes int `json:"average_session_duration_minutes"` // Rõ đơn vị
	TotalAuditLogs                int `json:"total_audit_logs"`
	FailedActionsLast24h          int `json:"failed_actions_last_24h"`
	ActiveSessionsNow             int `json:"active_sessions_now"`
}

// UserManagementAction represents admin actions on users
type UserManagementAction struct {
	Action AdminAction `json:"action" binding:"required"`
	Reason *string     `json:"reason,omitempty" binding:"omitempty,max=1000"`
}

// Validate validates the user management action
func (a *UserManagementAction) Validate() error {
	validActions := []AdminAction{
		AdminActionUserActivate,
		AdminActionUserDeactivate,
		AdminActionUserPromote,
		AdminActionUserDemote,
		AdminActionUserDelete,
		AdminActionUserBan,
		AdminActionUserUnban,
	}

	valid := false
	for _, validAction := range validActions {
		if a.Action == validAction {
			valid = true
			break
		}
	}

	if !valid {
		return fmt.Errorf("invalid action for user management: %s", a.Action)
	}

	return nil
}

// GroupManagementAction represents admin actions on groups
type GroupManagementAction struct {
	Action      AdminAction `json:"action" binding:"required"`
	Reason      *string     `json:"reason,omitempty" binding:"omitempty,max=1000"`
	NewLeaderID *int        `json:"new_leader_id,omitempty" binding:"omitempty,min=1"` // Required for transfer_leader
}

// Validate validates the group management action
func (a *GroupManagementAction) Validate() error {
	validActions := []AdminAction{
		AdminActionGroupActivate,
		AdminActionGroupDeactivate,
		AdminActionGroupDelete,
		AdminActionGroupTransferLeader,
	}

	valid := false
	for _, validAction := range validActions {
		if a.Action == validAction {
			valid = true
			break
		}
	}

	if !valid {
		return fmt.Errorf("invalid action for group management: %s", a.Action)
	}

	// transfer_leader requires NewLeaderID
	if a.Action == AdminActionGroupTransferLeader && a.NewLeaderID == nil {
		return fmt.Errorf("new_leader_id is required for transfer_leader action")
	}

	return nil
}

// SystemHealth represents system health status
type SystemHealth struct {
	Status            HealthStatus `json:"status"`
	DatabaseConnected bool         `json:"database_connected"`
	DatabaseLatency   int          `json:"database_latency_ms"`
	CacheConnected    *bool        `json:"cache_connected,omitempty"`  // Optional
	CacheLatency      *int         `json:"cache_latency_ms,omitempty"` // Optional
	ActiveConnections int          `json:"active_connections"`
	QueueSize         *int         `json:"queue_size,omitempty"` // Optional
	ErrorRate         float64      `json:"error_rate_percent"`
	LastHealthCheck   time.Time    `json:"last_health_check"`
	UptimeSeconds     int64        `json:"uptime_seconds"`
	Version           string       `json:"version"`
	Environment       string       `json:"environment"` // production, staging, development
}

// IsHealthy checks if system is healthy
func (s *SystemHealth) IsHealthy() bool {
	return s.Status == HealthStatusHealthy &&
		s.DatabaseConnected &&
		s.ErrorRate < 5.0
}

// IsDegraded checks if system is degraded
func (s *SystemHealth) IsDegraded() bool {
	return s.Status == HealthStatusDegraded
}

// IsDown checks if system is down
func (s *SystemHealth) IsDown() bool {
	return s.Status == HealthStatusDown || !s.DatabaseConnected
}

// BulkActionError represents an error in bulk operation
type BulkActionError struct {
	ID      int    `json:"id"`
	Message string `json:"message"`
}

// BulkActionResult represents result of a bulk operation
type BulkActionResult struct {
	TotalRequested int                `json:"total_requested"`
	Successful     int                `json:"successful"`
	Failed         int                `json:"failed"`
	Errors         []*BulkActionError `json:"errors,omitempty"`        // Chi tiết hơn với ID
	ProcessedIDs   []int              `json:"processed_ids,omitempty"` // IDs thành công
	FailedIDs      []int              `json:"failed_ids,omitempty"`    // IDs thất bại
}

// IsFullSuccess checks if all operations succeeded
func (r *BulkActionResult) IsFullSuccess() bool {
	return r.Failed == 0 && r.Successful == r.TotalRequested
}

// IsPartialSuccess checks if some operations succeeded
func (r *BulkActionResult) IsPartialSuccess() bool {
	return r.Successful > 0 && r.Failed > 0
}

// IsFullFailure checks if all operations failed
func (r *BulkActionResult) IsFullFailure() bool {
	return r.Successful == 0 && r.Failed == r.TotalRequested
}

// SuccessRate returns the success rate percentage
func (r *BulkActionResult) SuccessRate() float64 {
	if r.TotalRequested == 0 {
		return 0
	}
	return float64(r.Successful) / float64(r.TotalRequested) * 100
}

// BulkUserAction represents bulk action on multiple users
type BulkUserAction struct {
	UserIDs []int       `json:"user_ids" binding:"required,min=1,max=100,dive,min=1"`
	Action  AdminAction `json:"action" binding:"required"`
	Reason  *string     `json:"reason,omitempty" binding:"omitempty,max=1000"`
}

// Validate validates the bulk user action
func (a *BulkUserAction) Validate() error {
	if len(a.UserIDs) == 0 {
		return fmt.Errorf("user_ids cannot be empty")
	}

	if len(a.UserIDs) > 100 {
		return fmt.Errorf("cannot process more than 100 users at once")
	}

	// Check for duplicates
	seen := make(map[int]bool)
	for _, id := range a.UserIDs {
		if seen[id] {
			return fmt.Errorf("duplicate user_id: %d", id)
		}
		seen[id] = true
	}

	validActions := []AdminAction{
		AdminActionUserActivate,
		AdminActionUserDeactivate,
		AdminActionUserDelete,
		AdminActionUserBan,
		AdminActionUserUnban,
	}

	valid := false
	for _, validAction := range validActions {
		if a.Action == validAction {
			valid = true
			break
		}
	}

	if !valid {
		return fmt.Errorf("invalid action for bulk user management: %s", a.Action)
	}

	return nil
}

// BulkGroupAction represents bulk action on multiple groups
type BulkGroupAction struct {
	GroupIDs []int       `json:"group_ids" binding:"required,min=1,max=100,dive,min=1"`
	Action   AdminAction `json:"action" binding:"required"`
	Reason   *string     `json:"reason,omitempty" binding:"omitempty,max=1000"`
}

// Validate validates the bulk group action
func (a *BulkGroupAction) Validate() error {
	if len(a.GroupIDs) == 0 {
		return fmt.Errorf("group_ids cannot be empty")
	}

	if len(a.GroupIDs) > 100 {
		return fmt.Errorf("cannot process more than 100 groups at once")
	}

	// Check for duplicates
	seen := make(map[int]bool)
	for _, id := range a.GroupIDs {
		if seen[id] {
			return fmt.Errorf("duplicate group_id: %d", id)
		}
		seen[id] = true
	}

	validActions := []AdminAction{
		AdminActionGroupActivate,
		AdminActionGroupDeactivate,
		AdminActionGroupDelete,
	}

	valid := false
	for _, validAction := range validActions {
		if a.Action == validAction {
			valid = true
			break
		}
	}

	if !valid {
		return fmt.Errorf("invalid action for bulk group management: %s", a.Action)
	}

	return nil
}

// AdminStatsFilter represents filter for admin statistics
type AdminStatsFilter struct {
	StartDate *time.Time `json:"start_date,omitempty" form:"start_date"`
	EndDate   *time.Time `json:"end_date,omitempty" form:"end_date"`
}

// Validate validates the stats filter
func (f *AdminStatsFilter) Validate() error {
	if f.StartDate != nil && f.EndDate != nil {
		if f.EndDate.Before(*f.StartDate) {
			return fmt.Errorf("end_date must be after start_date")
		}
	}
	return nil
}

// Helper methods for AdminAuditLog
func (a *AdminAuditLog) IsSuccessful() bool {
	return a.Success
}

func (a *AdminAuditLog) HasError() bool {
	return a.ErrorMessage != nil && *a.ErrorMessage != ""
}

func (a *AdminAuditLog) IsUserAction() bool {
	return a.TargetType != nil && *a.TargetType == AuditTargetUser
}

func (a *AdminAuditLog) IsGroupAction() bool {
	return a.TargetType != nil && *a.TargetType == AuditTargetGroup
}
