package models

import "time"

// Group represents a survival group
type Group struct {
	ID          int     `json:"id" db:"id"`
	Name        string  `json:"name" db:"name"`
	Description *string `json:"desciption,omitempty" db:"desciption"`
	LeaderID    int     `json:"leader_id" db:"leader_id"`

	//Settings
	IsActive   bool `json:"is_active" db:"is_active"`
	MaxMembers int  `json:"max_members" db:"max_members"`

	AvatarURL *string `json:"avatar_url,omitempty" db:"avatar_url"`

	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

// GroupCreate represents group creation request
type GroupCreate struct {
	Name        string  `json:"name" binding:"required,max=255"`
	Description *string `json:"description,omitempty"`
	MaxMembers  *int    `json:"max_members,omitempty" binding:"omitempty,min=1,max=100"`
	AvatarURL   *string `json:"avatar_url,omitempty"`
}

// GroupUpdate represent group update request
type GroupUpdate struct {
	Name       *string `json:"name,omitempty" binding:"omitempty,max=255"`
	Desciption *string `json:"desrciption,omitempty"`
	MaxMembers *int    `json:"max_members,omitempty" binding:"omitempty,min=1,max=100"`
	AvatarURl  *string `json:"avatar_url,omitempty"`
	IsActive   *bool   `json:"is_active,omitempty"`
}

// GroupRepsonse represents group response with additional info
type GroupResponse struct {
	ID            int       `json:"id"`
	Name          string    `json:"name"`
	Description   *string   `json:"desciption,omitempty"`
	LeaderID      int       `json:"leader_id"`
	LeaderName    string    `json:"leader_name"`
	IsActive      bool      `json:"is_active"`
	MaxMembers    int       `json:"max_members"`
	ActiveMembers int       `json:"active_members"`
	AvatarURL     *string   `json:"avatar_url,omitempty"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

// GroupMember represents a member of a group
type GroupMember struct {
	ID       int          `json:"id" db:"id"`
	GroupID  string       `json:"group_id" db:"group_id"`
	UserID   string       `json:"user_id" db:"user_id"`
	Status   MemberStatus `json:"status" db:"status"`
	AddedBy  *int         `json:"added_by,omitempty" db:"added_by"`
	JoinedAt time.Time    `json:"joined_at" db:"joined_at"`
	LeftAt   *time.Time   `json:"left_at,omitemtpy" db:"left_at"`
}

// GroupMemberAt represents request to member to group
type GroupMemberAt struct {
	UserID int `json:"user_id" biding:"required"`
}

// GroupMemberResponse represents group member with user details
type GroupMemberResponse struct {
	ID        int          `json:"id"`
	GroupID   int          `json:"group_id"`
	UserID    int          `json:"user_id"`
	UserName  string       `json:"user_name"`
	UserEmail string       `json:"user_email"`
	AvatarURL *string      `json:"avatar_url,omitempty"`
	Status    MemberStatus `json:"status"`
	JoinedAt  time.Time    `json:"joined_at"`
	LeftAt    *time.Time   `json:"left_at,omitempty"`
}

// GroupWithMembers represents group with full member list
type GroupWithMembers struct {
	Group  *GroupResponse         `json:"group"`
	Member []*GroupMemberResponse `json:"members"`
}

// GroupSummary represents summary view from database view
type GroupSummary struct {
	ID            int    `json:"id" db:"id"`
	Name          string `json:"name" db:"name"`
	LeaderID      int    `json:"leader_id" db:"leader_id"`
	LeaderName    string `json:"leader_name" db:"leader_name"`
	ActiveMembers int    `json:"active_members" db:"active_members"`
	MaxMembers    int    `json:"max_members" db:"max_members"`
	IsActive      bool   `json:"is_active" db:"is_active"`
}
