package models

import "time"

type PushNotification struct {
	ID        int       `json:"id"`
	UserID    int       `json:"user_id"`
	Title     string    `json:"title"`
	Body      string    `json:"body"`
	Data      string    `json:"data"`
	Sent      string    `json:"sent"`
	CreatedAt time.Time `json:"created_at"`
}
