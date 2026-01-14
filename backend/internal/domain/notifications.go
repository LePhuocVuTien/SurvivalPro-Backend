package models

import "time"

type Notification struct {
	ID int `json:"id" db:"id"`

	//Recipient(XOR: iether user_id or group_id, not both)
	UserID  *int `json:"user_id,omitempty" db:"user_id"`
	GroupID *int `json:"group_id,omitempty" db:"group_id"`

	SenderID *int `json:"sender_id,omitempty" db:"sender_id"`

	Title    string           `json:"title" db:"title"`
	Body     string           `json:"body" db:"body"`
	Type     NotificationType `json:"type" db:"type"`
	Priority PriorityLevel    `json:"priority" db:"priority"`

	Data       JSONB   `json:"data" db:"data"`
	ActionType *string `json:"action_type,omitempty" db:"action_type"`
	ActionURL  *string `json:"action_url,omitempty" vb:"action_url"`

	Sent   bool       `json:"sent" db:"sent"`
	SentAt *time.Time `json:"sent_at,omitempty" db:"sent_at"`

	IsRead bool       `json:"is_read" db:"is_read"`
	ReadAt *time.Time `json:"read_at,omitempty" db:"read_at"`

	ScheduledFor *time.Time `json:"scheduled_for,omitempty" db:"scheduled_for"`

	CreatedAt time.Time  `json:"created_at" db:"created_at"`
	ExpiresAt *time.Time `json:"expires_at" db:"expires_at"`
}

// ToResponse converts Notification to NotificationResponse
func (n *Notification) ToResponse() *NotificationResponse {
	return &NotificationResponse{
		ID:           n.ID,
		UserID:       n.UserID,
		GroupID:      n.GroupID,
		SenderID:     n.SenderID,
		Title:        n.Title,
		Body:         n.Body,
		Type:         n.Type,
		Priority:     n.Priority,
		Data:         n.Data,
		ActionType:   n.ActionType,
		ActionURL:    n.ActionURL,
		Sent:         n.Sent,
		SentAt:       n.SentAt,
		IsRead:       n.IsRead,
		ReadAt:       n.ReadAt,
		ScheduledFor: n.ScheduledFor,
		CreatedAt:    n.CreatedAt,
		ExpiresAt:    n.ExpiresAt,
	}
}

// MarkAsRead marks notification as read
func (n *Notification) MarkAsRead() {
	if !n.IsRead {
		now := time.Now()
		n.IsRead = true
		n.ReadAt = &now
	}
}

// MarkAsSent marks notification as sent
func (n *Notification) MarkAsSent() {
	if !n.Sent {
		now := time.Now()
		n.Sent = true
		n.SentAt = &now
	}
}

// Expired checks if notification is expired
func (n *Notification) IsExpired() bool {
	if n.ExpiresAt == nil {
		return false
	}
	return time.Now().After(*n.ExpiresAt)
}

// IsScheduled checks if notification is schedule for future
func (n *Notification) IsScheduled() bool {
	if n.ScheduledFor == nil {
		return false
	}

	return time.Now().Before(*n.ScheduledFor)
}

// ShouldSendNow checks if notification should be sent now
func (n *Notification) ShouldSendNow() bool {
	if n.Sent {
		return false
	}
	if n.IsExpired() {
		return false
	}
	if n.ScheduledFor == nil {
		return true
	}
	return time.Now().After(*n.ScheduledFor) || time.Now().Equal(*n.ScheduledFor)
}

// ValidateRecipient XOR constaint
func (n *Notification) ValidateRecipient() error {
	if n.UserID == nil && n.GroupID == nil {
		return ErrInvalidRecipient
	}
	if n.UserID != nil && n.GroupID != nil {
		return ErrInvalidRecipient
	}
	return nil
}

// ==================== Request/Response Models ====================

// NotificationCreate represents notifiocation creation request
type NotificationCreate struct {
	// Target (XOR: either UserID OR GroupID)
	UserID  *int `json:"user_id,omitempty"`
	GroupID *int `json:"group_id,omitempty"`

	Title        string           `json:"title" binding:"required,max=255"`
	Body         string           `json:"body" binding:"required,max=1000"`
	Type         NotificationType `json:"type" binding:"required"`
	Priority     *PriorityLevel   `json:"priority,omitempty"`
	Data         map[string]any   `json:"data,omitempty"`
	ActionType   *string          `json:"action_type,omitempty" binding:"omitempty,max=50"`
	ActionURL    *string          `json:"action_url,omitempty" binding:"omitempty,url"`
	ScheduledFor *time.Time       `json:"scheduled_for,omitempty"`
	ExpiresAt    *time.Time       `json:"expires_at,omitempty"`
}

// Validate validates notification create request
func (n *NotificationCreate) Validate() error {
	if n.UserID == nil && n.GroupID == nil {
		return ErrInvalidRecipient
	}

	if n.UserID != nil && n.GroupID != nil {
		return ErrInvalidRecipient
	}

	return nil
}

// GetPriority returns priority with default
func (n *NotificationCreate) GetPriority() PriorityLevel {
	if n.Priority == nil {
		return PriorityLevelMedium
	}
	return *n.Priority
}

// NotificaionReponse represents notification with sender info
type NotificationResponse struct {
	ID           int              `json:"id"`
	UserID       *int             `json:"user_id,omitempty"`
	GroupID      *int             `json:"group_id,omitempty"`
	SenderID     *int             `json:"sender_id,omitempty"`
	SenderName   *string          `json:"sender_name,omitempty"`
	Title        string           `json:"title"`
	Body         string           `json:"body"`
	Type         NotificationType `json:"type"`
	Priority     PriorityLevel    `json:"priority"`
	Data         JSONB            `json:"data,omitempty"`
	ActionType   *string          `json:"action_type,omitempty"`
	ActionURL    *string          `json:"action_url,omitempty"`
	Sent         bool             `json:"sent"`
	SentAt       *time.Time       `json:"sent_at,omitempty"`
	IsRead       bool             `json:"is_read"`
	ReadAt       *time.Time       `json:"read_at,omitempty"`
	ScheduledFor *time.Time       `json:"scheduled_for,omitempty"`
	CreatedAt    time.Time        `json:"created_at"`
	ExpiresAt    *time.Time       `json:"expires_at,omitempty"`
}

// NotificationListResponse represents paginated notification list
type NotificationListResponse struct {
	Notifications []*NotificationResponse `json:"notifications"`
	Total         int                     `json:"total"`
	Unread        int                     `json:"unread"`
	Limit         int                     `json:"limit"`
	Offset        int                     `json:"offset"`
	HasMore       bool                    `json:"has_more"`
}

// ==================== Delivery Tracking ====================

// NotificationDelivery represents notification delivery to a specific user
type NotificationDelivery struct {
	ID             int        `json:"id" db:"id"`
	NotificationId int        `json:"notification_id" data:"notification_id"`
	UserID         int        `json:"user_id" db:"user_id"`
	Delivered      bool       `json:"delivered" db:"delivered"`
	DeliveredAt    *time.Time `json:"delivered_at,omitempty" db:"delivered_at"`
	IsRead         bool       `json:"is_read" db:"is_read"`
	ReadAt         *time.Time `json:"read_at,omitempty" db:"read_at"`
	Failed         bool       `json:"failed" db:"failed"`
	ErrorMessage   *string    `json:"error_message,omitempty" db:"error_message"`
	CreatedAt      time.Time  `json:"created_at" db:"created_at"`
}

// MarkAsDelivered marks delivery as delivered
func (d *NotificationDelivery) MarkAsDelivered() {
	if !d.Delivered {
		now := time.Now()
		d.Delivered = true
		d.DeliveredAt = &now
		d.Failed = false
		d.ErrorMessage = nil
	}
}

// MarkAsFailed marks delivery as failed
func (d *NotificationDelivery) MarkAsFailed(errorMsg string) {
	d.Failed = true
	d.ErrorMessage = &errorMsg
	d.Delivered = false
}

// MarkAsRead marks delivery as read
func (d *NotificationDelivery) MaskAsRead() {
	if !d.IsRead {
		now := time.Now()
		d.IsRead = true
		d.ReadAt = &now
	}
}

// ==================== Bulk Operations ====================
// NotificationMarkRead represents mark as read request
type NotificationMarkRead struct {
	NotificationIDs []int `json:"notification_ids" binding:"required, min=1, max=100"`
}

// BulkNotificationCreate represents bulk notification creation
type BulkNotificationCreate struct {
	UserIDs      []int            `json:"user_ids" binding:"required,min=1,max=1000"`
	Title        string           `json:"title" binding:"required,max=255"`
	Body         string           `json:"body" binding:"required,max=1000"`
	Type         NotificationType `json:"type" binding:"required"`
	Priority     *PriorityLevel   `json:"priority,omitempty"`
	Data         map[string]any   `json:"data,omitempty"`
	ActionType   *string          `json:"action_type,omitempty" binding:"omitempty,max=50"`
	ActionURL    *string          `json:"action_url,omitempty" binding:"omitempty,url"`
	ScheduledFor *time.Time       `json:"scheduled_for,omitempty"`
	ExpiresAt    *time.Time       `json:"expires_at,omitempty"`
}

// GetPriority returns priority with default
func (b *BulkNotificationCreate) GetPriority() PriorityLevel {
	if b.Priority == nil {
		return PriorityLevelMedium
	}

	return *b.Priority
}

// ==================== Filter & Query ====================

// NotificationFilter represents filter options for notifications
type NotificationFilter struct {
	Type      *NotificationType `json:"type,omitempty" form:"type"`
	Priority  *PriorityLevel    `json:"priority,omitempty" form:"priority"`
	IsRead    *bool             `json:"is_read,omitempty" form:"is_read"`
	Limit     int               `json:"litmit" form:"limit" binding:"min=1,max=100"`
	Offset    int               `json:"offset" form:"offset" binding:"min=0"`
	SortBy    string            `json:"sort_by" form:"sort_by" binding:"oneof=created_at priority type ''"`
	SortOrder string            `json:"sort_order" form:"sort_order" binding:"oneof=asc desc ''"`
}

// SetDefault sets defualt values for filter
func (f *NotificationFilter) SetDefault() {
	if f.Limit == 0 {
		f.Limit = 20
	}
	if f.SortBy == "" {
		f.SortBy = "created_at"
	}
	if f.SortOrder == "" {
		f.SortOrder = "desc"
	}
}

// ==================== Statistics ====================

// NotificationStats represents notification statistics
type NotificationStats struct {
	TotalNotifications  int `json:"total_notifications"`
	UnreadNotifications int `json:"unread_notifications"`
	ReadNotifications   int `json:"read_notifications"`
	UrgentNotifications int `json:"ungent_notifications"` //high or critical priority
	ExpiringSoon        int `json:"expiring_soon"`        // expires in next 24h
}

// NotificationTypeCount represents count by notification type
type NotificationTypeCount struct {
	Type  NotificationType `json:"notification_type"`
	Count int              `json:"count"`
}

// NotificationPriorityCount represents count by priority
type NotificationPriorityCount struct {
	Priority PriorityLevel `json:"priority"`
	Count    int           `json:"count"`
}

// DetailedNotificationStats represents detailed statistics
type DetailedNotificationStats struct {
	NotificationStats
	ByType     []NotificationTypeCount     `json:"by_type"`
	ByPriority []NotificationPriorityCount `json:"by_priority"`
}
