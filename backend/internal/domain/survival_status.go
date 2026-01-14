package models

import "time"

// SurvivalStatusModel represent user survival status
type SurvivalStatusModel struct {
	UserID int            `json:"user_id" db:"user_id"`
	Status SurvivalStatus `json:"status" db:"status"`

	// Location (PostGIS geography point)
	Location         *Point   `json:"location,omitempty" db:"location"`
	LocationAccuracy *float64 `json:"location_accuracy" db:"location_accuracy"`

	BatteryLevel *int `json:"battery_level,omitempty" db:"battery_level"`
	IsOnline     bool `json:"is_online" db:"is_online"`

	//SOS
	SOSActive      bool       `json:"sos_active" db:"sos_active"`
	SOSActivatedAt *time.Time `json:"sos_activated,omitempty" db:"sos_activated_at"`

	LastUpdatedAt time.Time `json:"last_update_at" db:"last_updated_at"`
	LastSeenAt    time.Time `json:"last_seen_at" db:"last_seen_at"`
}

// SurvivalStatusUpdate represents status update request
type SurvivalStatusUpdate struct {
	Status       *SurvivalStatus `json:"status,omitempty"`
	Location     *Point          `json:"point,omitempty"`
	Accuracy     *float64        `json:"accuracy,omitempty"`
	BatteryLevel *int            `json:"battery_level,omitempty"`
	IsOnline     *bool           `json:"is_online,omitempty"`
}

// UserLocation represents a location history entry
type Userlocation struct {
	ID           int       `json:"id" db:"id"`
	UserID       int       `json:"user_id" db:"user_id"`
	Location     Point     `json:"location" db:"location"`
	Accuracy     *float64  `json:"accurac,omitempty" db:"accuracy"`  // meters
	Altitude     *float64  `json:"altitude,omitempty" db:"altitude"` // meters
	BatteryLevel *int      `json:"battery_level,omitempty" db:"battery_level"`
	IsMoving     *bool     `json:"is_moving,omitempty" db:"is_moving"`
	Speed        *float64  `json:"speed,omitempty" db:"speed"` // km/h
	Timestamp    time.Time `json:"timestamp" db:"timestamp"`
}

// LocationCreate represents location creation request
type LocationCreate struct {
	Location     Point    `json:"location" binding:"required"`
	Accuracy     *float64 `json:"accuracy,omitempty"`
	Altitude     *float64 `json:"altitude,omitempty"`
	BatteryLevel *int     `json:"battery_level,omitempty" binding:"omitempty,min=0,max=100"`
	IsMoving     *bool    `json:"is_moving,omitempty"`
	Speed        *float64 `json:"speed,omitempty"`
}

// LocationRequest represents location wiith user info
type LocationRequest struct {
	ID           int       `json:"id"`
	UserID       int       `json:"user_id"`
	UserName     string    `json:"user_name"`
	Location     Point     `json:"location"`
	Accuracy     *float64  `json:"accuracy,omitempty"`
	Altitude     *float64  `json:"altiitude,omitempty"`
	BatteryLevel *int      `json:"battery_level,omitempty"`
	IsMoving     *bool     `json:"is_moving,omitempty"`
	Speed        *float64  `json:"speed,omitempty"`
	Timestamp    time.Time `json:"timestamp"`
}

// NearByUser represents a user within range
type NearByUser struct {
	UserID       int            `json:"user_id"`
	Name         string         `json:"name"`
	AvatarURL    *string        `json:"avatar_url,omitempty"`
	Status       SurvivalStatus `json:"status"`
	Location     Point          `json:"location"`
	Distance     float64        `json:"distance"`
	BatteryLevel *int           `json:"battery_level,omitempty"`
	LastSeen     time.Time      `json:"last_seen"`
}

// DistanceQuery represent query for nearby users
type DistanceQuery struct {
	Location Point   `json:"location" binding:"required"`
	RadiusKM float64 `json:"radius_km" binding:"required,min=0.1,max=100"`
	Limit    int     `json:"limit,omitempty" binding:"omitempty,min=1,max=100"`
}
