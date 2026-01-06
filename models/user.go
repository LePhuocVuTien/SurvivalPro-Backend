package models

import "time"

type User struct {
	ID           int       `json:"id"`
	Email        string    `json:"email"`
	Password     string    `json:"password,omitempty"`
	Name         string    `json:"name"`
	BloodType    string    `json:"blood_type"`
	Allergies    string    `json:"allergies"`
	EmergencyNum string    `json:"emergency_num"`
	AvatarURL    string    `json:"avatar_url"`
	PushToken    string    `json:"push_token"`
	CreatedAt    time.Time `json:"created_at"`
}
