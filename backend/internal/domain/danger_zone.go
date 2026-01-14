package models

import (
	"fmt"
	"time"
)

// DangerZone represents a danger zone area
type DangerZone struct {
	ID          int              `json:"id" db:"id"`
	Name        string           `json:"name" db:"name"`
	Description *string          `json:"description,omitempty" db:"description"`
	Area        GeographyPolygon `json:"-" db:"area"` // PostGIS raw data (EWKB/WKB)
	Severity    DangerSeverity   `json:"severity" db:"severity"`
	Type        DangerType       `json:"type" db:"type"`
	ActiveFrom  time.Time        `json:"active_from" db:"active_from"`
	ActiveUntil time.Time        `json:"active_until" db:"active_until"`
	CreatedBy   *int             `json:"created_by,omitempty" db:"created_by"`
	IsActive    bool             `json:"is_active" db:"is_active"`
	CreatedAt   time.Time        `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time        `json:"updated_at" db:"updated_at"`
	DeletedAt   *time.Time       `json:"-" db:"deleted_at"`
}

// DangerZoneCreate represents danger zone creation request
type DangerZoneCreate struct {
	Name        string         `json:"name" binding:"required,max=255"`
	Description *string        `json:"description,omitempty" binding:"omitempty,max=1000"`
	Area        Polygon        `json:"area" binding:"required"`
	Severity    DangerSeverity `json:"severity" binding:"required"`
	Type        DangerType     `json:"type" binding:"required"`
	ActiveFrom  time.Time      `json:"active_from" binding:"required"`
	ActiveUntil time.Time      `json:"active_until" binding:"required"`
}

// Validate validates the creation request
func (c *DangerZoneCreate) Validate() error {
	if !c.Area.IsValid() {
		return fmt.Errorf("invalid polygon: must be valid GeoJSON with closed rings")
	}

	if c.ActiveUntil.Before(c.ActiveFrom) || c.ActiveUntil.Equal(c.ActiveFrom) {
		return fmt.Errorf("active_until must be after active_from")
	}

	if c.ActiveFrom.Before(time.Now().Add(-24 * time.Hour)) {
		return fmt.Errorf("active_from cannot be more than 24 hours in the past")
	}

	if !c.Severity.IsValid() {
		return fmt.Errorf("invalid severity value")
	}

	if !c.Type.IsValid() {
		return fmt.Errorf("invalid danger type")
	}

	return nil
}

// DangerZoneUpdate represents danger zone update request
type DangerZoneUpdate struct {
	Name        *string         `json:"name,omitempty" binding:"omitempty,max=255"`
	Description *string         `json:"description,omitempty" binding:"omitempty,max=1000"`
	Area        *Polygon        `json:"area,omitempty"`
	Severity    *DangerSeverity `json:"severity,omitempty"`
	Type        *DangerType     `json:"type,omitempty"`
	ActiveFrom  *time.Time      `json:"active_from,omitempty"`
	ActiveUntil *time.Time      `json:"active_until,omitempty"`
	IsActive    *bool           `json:"is_active,omitempty"`
}

// Validate validates the update request
func (u *DangerZoneUpdate) Validate() error {
	if u.Area != nil && !u.Area.IsValid() {
		return fmt.Errorf("invalid polygon: must be valid GeoJSON with closed rings")
	}

	if u.ActiveFrom != nil && u.ActiveUntil != nil {
		if u.ActiveUntil.Before(*u.ActiveFrom) || u.ActiveUntil.Equal(*u.ActiveFrom) {
			return fmt.Errorf("active_until must be after active_from")
		}
	}

	if u.Severity != nil && !u.Severity.IsValid() {
		return fmt.Errorf("invalid severity value")
	}

	if u.Type != nil && !u.Type.IsValid() {
		return fmt.Errorf("invalid danger type")
	}

	return nil
}

// DangerZoneResponse represents danger zone with creator info
type DangerZoneResponse struct {
	ID              int            `json:"id"`
	Name            string         `json:"name"`
	Description     *string        `json:"description,omitempty"`
	Area            *Polygon       `json:"area,omitempty"` // Parsed from PostGIS, pointer để tránh panic
	Severity        DangerSeverity `json:"severity"`
	Type            DangerType     `json:"type"`
	ActiveFrom      time.Time      `json:"active_from"`
	ActiveUntil     time.Time      `json:"active_until"`
	CreatedBy       *int           `json:"created_by,omitempty"`
	CreatorName     *string        `json:"creator_name,omitempty"`
	IsActive        bool           `json:"is_active"`
	IsCurrentActive bool           `json:"is_current_active"`       // Computed: now between ActiveFrom and ActiveUntil
	UsersInZone     *int           `json:"users_in_zone,omitempty"` // Computed, pointer để omitempty hoạt động
	CreatedAt       time.Time      `json:"created_at"`
	UpdatedAt       time.Time      `json:"updated_at"`
}

// DangerZoneAlert represents an alert for a user in danger zone
type DangerZoneAlert struct {
	ID             int        `json:"id" db:"id"`
	DangerZoneID   int        `json:"danger_zone_id" db:"danger_zone_id"`
	UserID         int        `json:"user_id" db:"user_id"`
	EnteredAt      time.Time  `json:"entered_at" db:"entered_at"`
	ExitedAt       *time.Time `json:"exited_at,omitempty" db:"exited_at"`
	Notified       bool       `json:"notified" db:"notified"`
	NotifiedAt     *time.Time `json:"notified_at,omitempty" db:"notified_at"`
	Acknowledged   bool       `json:"acknowledged" db:"acknowledged"`
	AcknowledgedAt *time.Time `json:"acknowledged_at,omitempty" db:"acknowledged_at"`
	CreatedAt      time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at" db:"updated_at"`
}

// DangerZoneAlertResponse represents alert with zone and user details
type DangerZoneAlertResponse struct {
	ID              int            `json:"id"`
	DangerZoneID    int            `json:"danger_zone_id"`
	DangerZoneName  string         `json:"danger_zone_name"`
	Severity        DangerSeverity `json:"severity"`
	Type            DangerType     `json:"type"`
	UserID          int            `json:"user_id"`
	UserName        string         `json:"user_name"`
	UserLocation    *Point         `json:"user_location,omitempty"`
	EnteredAt       time.Time      `json:"entered_at"`
	ExitedAt        *time.Time     `json:"exited_at,omitempty"`
	Notified        bool           `json:"notified"`
	NotifiedAt      *time.Time     `json:"notified_at,omitempty"`
	Acknowledged    bool           `json:"acknowledged"`
	AcknowledgedAt  *time.Time     `json:"acknowledged_at,omitempty"`
	DurationMinutes *float64       `json:"duration_minutes,omitempty"` // Computed
	CreatedAt       time.Time      `json:"created_at"`
}

// CheckLocationInDangerZone represents request to check if location is in danger
type CheckLocationInDangerZone struct {
	Location Point `json:"location" binding:"required"`
	Radius   *int  `json:"radius,omitempty" binding:"omitempty,min=0,max=10000"` // Meters
}

// Validate validates the location check request
func (c *CheckLocationInDangerZone) Validate() error {
	if !c.Location.IsValid() {
		return fmt.Errorf("invalid location coordinates")
	}
	return nil
}

// CheckLocationResponse represents response for location check
type CheckLocationResponse struct {
	InDangerZone bool                  `json:"in_danger_zone"`
	DangerZones  []*DangerZoneResponse `json:"danger_zones,omitempty"`
	SafetyLevel  SafetyLevel           `json:"safety_level"` // Enum
	Warnings     []string              `json:"warnings,omitempty"`
}

// DangerZoneFilter represents filter options for danger zone queries
type DangerZoneFilter struct {
	Severity  *DangerSeverity `json:"severity,omitempty" form:"severity"`
	Type      *DangerType     `json:"type,omitempty" form:"type"`
	IsActive  *bool           `json:"is_active,omitempty" form:"is_active"`
	ActiveNow *bool           `json:"active_now,omitempty" form:"active_now"`
	Latitude  *float64        `json:"latitude,omitempty" form:"lat" binding:"omitempty,min=-90,max=90"`
	Longitude *float64        `json:"longitude,omitempty" form:"lng" binding:"omitempty,min=-180,max=180"`
	Radius    *int            `json:"radius,omitempty" form:"radius" binding:"omitempty,min=0,max=50000"` // Meters
	Limit     int             `json:"limit,omitempty" form:"limit" binding:"omitempty,min=1,max=100"`
	Offset    int             `json:"offset,omitempty" form:"offset" binding:"omitempty,min=0"`
	SortBy    *string         `json:"sort_by,omitempty" form:"sort_by"`       // severity, created_at, active_from
	SortOrder *string         `json:"sort_order,omitempty" form:"sort_order"` // asc, desc
}

// GetLocation returns Point if lat/lng are provided
func (f *DangerZoneFilter) GetLocation() *Point {
	if f.Latitude != nil && f.Longitude != nil {
		return &Point{
			Latitude:  *f.Latitude,
			Longitude: *f.Longitude,
		}
	}
	return nil
}

// DangerZoneStats represents danger zone statistics
type DangerZoneStats struct {
	TotalZones           int                    `json:"total_zones"`
	ActiveZones          int                    `json:"active_zones"`
	CurrentlyActiveZones int                    `json:"currently_active_zones"`
	UsersInDanger        int                    `json:"users_in_danger"`
	UnacknowledgedAlerts int                    `json:"unacknowledged_alerts"`
	BySeverity           map[DangerSeverity]int `json:"by_severity"`
	ByType               map[DangerType]int     `json:"by_type"`
	CriticalZones        []*DangerZoneResponse  `json:"critical_zones,omitempty"`
}

// AcknowledgeAlertRequest represents alert acknowledgement request
type AcknowledgeAlertRequest struct {
	Notes *string `json:"notes,omitempty" binding:"omitempty,max=500"`
}

// BulkAcknowledgeRequest represents bulk alert acknowledgement
type BulkAcknowledgeRequest struct {
	AlertIDs []int   `json:"alert_ids" binding:"required,min=1,max=100,dive,min=1"`
	Notes    *string `json:"notes,omitempty" binding:"omitempty,max=500"`
}

// Helper methods for DangerZone
func (d *DangerZone) IsCurrentlyActive() bool {
	now := time.Now()
	return d.IsActive &&
		now.After(d.ActiveFrom) &&
		now.Before(d.ActiveUntil)
}

func (d *DangerZone) IsExpired() bool {
	return time.Now().After(d.ActiveUntil)
}

func (d *DangerZone) IsDeleted() bool {
	return d.DeletedAt != nil
}

func (d *DangerZone) DaysRemaining() int {
	if d.IsExpired() {
		return 0
	}
	duration := time.Until(d.ActiveUntil)
	return int(duration.Hours() / 24)
}

func (d *DangerZone) GetAreaPolygon() (*Polygon, error) {
	return d.Area.ToPolygon()
}

// Helper methods for DangerZoneAlert
func (a *DangerZoneAlert) IsActive() bool {
	return a.ExitedAt == nil
}

func (a *DangerZoneAlert) NeedsAcknowledgement() bool {
	return !a.Acknowledged && a.IsActive()
}

func (a *DangerZoneAlert) CalculateDuration() *float64 {
	endTime := time.Now()
	if a.ExitedAt != nil {
		endTime = *a.ExitedAt
	}
	duration := endTime.Sub(a.EnteredAt).Minutes()
	return &duration
}
