package models

import "time"

type CheckListItem struct {
	ID          int       `json:"id"`
	UserID      int       `json:"user_id"`
	Title       string    `json:"title"`
	Category    string    `json:"category"`
	Description string    `json:"description,omitempty"`
	IsChecked   bool      `json:"is_checked"`
	CreatedAt   time.Time `json:"created_at"`
}
