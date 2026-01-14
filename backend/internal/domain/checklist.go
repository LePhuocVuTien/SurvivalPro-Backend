package models

import "time"

// ChecklistItem represents a checklist item
type ChecklistItem struct {
	ID              int               `json:"id" db:"id"`
	UserID          int               `json:"user_id" db:"user_id"`
	Title           string            `json:"title" db:"title"`
	Description     *string           `json:"description,omitempty" db:"description"`
	Category        ChecklistCategory `json:"category" db:"category"`
	IsChecked       bool              `json:"is_checked" db:"is_checked"`
	Priority        PriorityLevel     `json:"priority" db:"priority"`
	QuantityNeeded  *int              `json:"quantity_needed,omitempty" db:"quantity_needed"`
	QuantityCurrent int               `json:"quantity_current" db:"quantity_current"`
	DueDate         *time.Time        `json:"due_date,omitempty" db:"due_date"`
	CheckedAt       *time.Time        `json:"checked_at,omitempty" db:"checked_at"`
	CreatedAt       time.Time         `json:"created_at" db:"created_at"`
	UpdatedAt       time.Time         `json:"updated_at" db:"updated_at"`
}

// ChecklistItemCreate represents checklist item creation request
type ChecklistItemCreate struct {
	Title           string            `json:"title" binding:"required,max=255"`
	Description     *string           `json:"description,omitempty"`
	Category        ChecklistCategory `json:"category" binding:"required"`
	Priority        *PriorityLevel    `json:"priority,omitempty"`
	QuantityNeeded  *int              `json:"quantity_needed,omitempty" binding:"omitempty,min=1"`
	QuantityCurrent *int              `json:"quantity_current" binding:"omitempty,min=0"`
	DueDate         *time.Time        `json:"due_date,omitempty"`
}

// ChecklistItemUpdate represents checklist item update request
type ChecklisitItemUpdate struct {
	Title           *string            `json:"title,omitempty" binding:"omitempty,max=255"`
	Descripition    *string            `json:"description,omitempty"`
	Category        *ChecklistCategory `json:"category,omitempty"`
	IsChecked       *bool              `json:"is_check,omitempty"`
	Priority        *PriorityLevel     `json:"priority,omitempty"`
	QuantityNeeded  *int               `json:"quantity_needed,omitempty" binding:"omitempty,min=1"`
	QuantityCurrent *int               `json:"quantity_current,omitempty" binding:"omitempty,min=0"`
	DueDate         *time.Time         `json:"due_date,omitempty"`
}

// ChecklistSummary represents checklist summary by category
type ChecklistSummary struct {
	Category      ChecklistCategory `json:"category"`
	TotalItems    int               `json:"totalItems"`
	CheckedItems  int               `json:"checked_items"`
	UrgentItems   int               `json:"urgent_items"` // high or critical priority
	OverdueItems  int               `json:"overdue_items"`
	CompletionPct float64           `json:"completion_pct"`
}

// ChecklistStats represents overall checklist statistics
type ChecklistStats struct {
	TotalItems     int                 `json:"total_items"`
	CheckedItems   int                 `json:"checked_items"`
	UncheckedItems int                 `json:"unchecked_items"`
	UrgentItems    int                 `json:"urgent_items"`
	OverdueItems   int                 `json:"overdue_items"`
	CompletionPct  float64             `json:"completion_pct"`
	ByCategory     []*ChecklistSummary `json:"by_category"`
	LastUpdated    *time.Time          `json:"last_updated,omitempty"`
}

// ChecklistFilter represents fileter options for checklist update
type ChecklistFilter struct {
	Category  *ChecklistCategory `json:"category,omitempty" form:"category"`
	IsChecked *bool              `json:"is_checked,omitempty" form:"is_checked"`
	Priority  *PriorityLevel     `json:"priority,omitempty" form:"priority"`
	Overdue   *bool              `json:"overdue" form:"overdue"`
	Limit     int                `json:"limit,omitempty" form:"limit" binding:"omitempty.min=1,max=100"`
	Offset    int                `json:"offset,omitempty" form:"offset" binding:"omitempty,min=0"`
}
