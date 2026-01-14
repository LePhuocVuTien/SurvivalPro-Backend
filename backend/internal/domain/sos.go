package models

import "time"

// SOSEvent represents an SOS event
type SOSEvent struct {
	ID                int        `json:"id" db:"id"`
	UserID            int        `json:"user_id" db:"user_id"`
	Location          *Point     `json:"location,omitempty" db:"location"`
	Status            SOSStatus  `json:"status" db:"status"`
	TriggerSource     SOSTrigger `json:"trigger_source" db:"trigger_source"`
	RespondedByUserID *int       `json:"responded_by_user_id,omitempty" db:"responded_by_user_id"`
	RespondedAt       *time.Time `json:"responded_at,omitempty" db:"responded_at"`
	ResponseNotes     *string    `json:"response_notes,omitempty" db:"response_notes"`
	ActivatedAt       time.Time  `json:"activated_at" db:"activated_at"` // DB default NOW()
	ResolvedAt        *time.Time `json:"resolved_at,omitempty" db:"resolved_at"`
	CancelledAt       *time.Time `json:"cancelled_at,omitempty" db:"cancelled_at"`
	BatteryLevel      *int       `json:"battery_level,omitempty" db:"battery_level"`
	DeviceInfo        *JSONB     `json:"device_info,omitempty" db:"device_info"`
	CreatedAt         time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt         time.Time  `json:"updated_at" db:"updated_at"`
}

// SOSEventCreate represents SOS activation request
type SOSEventCreate struct {
	Location      *Point         `json:"location,omitempty" binding:"required"`
	TriggerSource SOSTrigger     `json:"trigger_source" binding:"required"`
	BatteryLevel  *int           `json:"battery_level,omitempty" binding:"omitempty,min=0,max=100"`
	DeviceInfo    map[string]any `json:"device_info,omitempty"`
}

// SOSEventResponse represents SOS event with user details
type SOSEventResponse struct {
	ID                int                 `json:"id"`
	UserID            int                 `json:"user_id"`
	UserName          string              `json:"user_name"`
	UserPhone         *string             `json:"user_phone,omitempty"`
	UserAvatarURL     *string             `json:"user_avatar_url,omitempty"`
	UserBloodType     *string             `json:"user_blood_type,omitempty"`
	UserAllergies     *string             `json:"user_allergies,omitempty"`
	Location          *Point              `json:"location,omitempty"`
	Status            SOSStatus           `json:"status"`
	TriggerSource     SOSTrigger          `json:"trigger_source"`
	RespondedByUserID *int                `json:"responded_by_user_id,omitempty"`
	ResponderName     *string             `json:"responder_name,omitempty"`
	RespondedAt       *time.Time          `json:"responded_at,omitempty"`
	ResponseNotes     *string             `json:"response_notes,omitempty"`
	ActivatedAt       time.Time           `json:"activated_at"`
	ResolvedAt        *time.Time          `json:"resolved_at,omitempty"`
	CancelledAt       *time.Time          `json:"cancelled_at,omitempty"`
	BatteryLevel      *int                `json:"battery_level,omitempty"`
	DurationMinutes   *float64            `json:"duration_minutes,omitempty"`
	EmergencyContacts []*EmergencyContact `json:"emergency_contacts,omitempty"`
}

// CalculateDuration computes duration in minutes
func (s *SOSEventResponse) CalculateDuration() *float64 {
	if s.ResolvedAt == nil {
		return nil
	}
	duration := s.ResolvedAt.Sub(s.ActivatedAt).Minutes()
	return &duration
}

// SOSResolveRequest represents SOS resolution request
type SOSResolveRequest struct {
	Status        SOSStatus `json:"status" binding:"required,oneof=resolved false_alarm"`
	ResponseNotes *string   `json:"response_notes,omitempty" binding:"omitempty,max=1000"`
}

// SOSCancelRequest represents SOS cancellation request by user
type SOSCancelRequest struct {
	Reason *string `json:"reason,omitempty" binding:"omitempty,max=500"`
}

// EmergencyContact represents an emergency contact
type EmergencyContact struct {
	ID            int        `json:"id" db:"id"`
	UserID        int        `json:"user_id" db:"user_id"`
	Name          string     `json:"name" db:"name"`
	Phone         string     `json:"phone" db:"phone"`
	Relationship  *string    `json:"relationship,omitempty" db:"relationship"`
	Email         *string    `json:"email,omitempty" db:"email"`
	PriorityOrder int        `json:"priority_order" db:"priority_order"`
	NotifyOnSOS   bool       `json:"notify_on_sos" db:"notify_on_sos"`
	CreatedAt     time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at" db:"updated_at"`
	DeletedAt     *time.Time `json:"-" db:"deleted_at"`
}

// EmergencyContactCreate represents emergency contact creation request
type EmergencyContactCreate struct {
	Name          string  `json:"name" binding:"required,max=255"`
	Phone         string  `json:"phone" binding:"required,e164"`
	Relationship  *string `json:"relationship,omitempty" binding:"omitempty,max=100"`
	Email         *string `json:"email,omitempty" binding:"omitempty,email"`
	PriorityOrder int     `json:"priority_order" binding:"required,min=1,max=10"`
	NotifyOnSOS   bool    `json:"notify_on_sos"`
}

// EmergencyContactUpdate represents emergency contact update request
type EmergencyContactUpdate struct {
	Name          *string `json:"name,omitempty" binding:"omitempty,max=255"`
	Phone         *string `json:"phone,omitempty" binding:"omitempty,e164"`
	Relationship  *string `json:"relationship,omitempty" binding:"omitempty,max=100"`
	Email         *string `json:"email,omitempty" binding:"omitempty,email"`
	PriorityOrder *int    `json:"priority_order,omitempty" binding:"omitempty,min=1,max=10"`
	NotifyOnSOS   *bool   `json:"notify_on_sos,omitempty"`
}

// SOSEventStatusHistory tracks status changes for audit
type SOSEventStatusHistory struct {
	ID        int       `json:"id" db:"id"`
	EventID   int       `json:"event_id" db:"event_id"`
	Status    SOSStatus `json:"status" db:"status"`
	ChangedBy *int      `json:"changed_by,omitempty" db:"changed_by"`
	Notes     *string   `json:"notes,omitempty" db:"notes"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

// SOSStats represents SOS statistics
type SOSStats struct {
	TotalSOSEvents             int                 `json:"total_sos_events"`
	ActiveSOSEvents            int                 `json:"active_sos_events"`
	ResolvedSOSEvents          int                 `json:"resolved_sos_events"`
	FalseAlarmEvents           int                 `json:"false_alarm_events"`
	CancelledEvents            int                 `json:"cancelled_events"`
	AverageResponseTimeMinutes float64             `json:"average_response_time_minutes"`
	ByStatus                   map[SOSStatus]int   `json:"by_status"`
	ByTriggerSource            map[SOSTrigger]int  `json:"by_trigger_source"`
	RecentEvents               []*SOSEventResponse `json:"recent_events,omitempty"`
}

// Helper methods for SOSEvent
func (s *SOSEvent) IsActive() bool {
	return s.Status == SOSStatusActive
}

func (s *SOSEvent) CanBeResolved() bool {
	return s.IsActive() && s.ResolvedAt == nil && s.CancelledAt == nil
}

func (s *SOSEvent) CanBeCancelled() bool {
	return s.Status == SOSStatusActive && s.RespondedAt == nil && s.CancelledAt == nil
}

func (s *SOSEvent) IsClosed() bool {
	return s.Status == SOSStatusResolved ||
		s.Status == SOSStatusFalseAlarm ||
		s.CancelledAt != nil
}

// Helper methods for EmergencyContact
func (e *EmergencyContact) IsDeleted() bool {
	return e.DeletedAt != nil
}

func (e *EmergencyContact) IsActive() bool {
	return !e.IsDeleted()
}
