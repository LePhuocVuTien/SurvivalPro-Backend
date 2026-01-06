package models

import "time"

type UserLocation struct {
	ID         int       `json:"id"`
	UserID     int       `json:"user_id"`
	Latitude   float64   `json:"latitude"`
	Longtitude float64   `json"longtitude"`
	Timestamp  time.Time `json:"timestamp"`
}
