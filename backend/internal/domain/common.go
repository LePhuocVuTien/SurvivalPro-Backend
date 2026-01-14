package models

import (
	"fmt"
	"time"
)

// ==================== Response Structures ====================

// APIResponse is a generic API response wrapper (Go 1.18+ generic version)
type APIResponse[T any] struct {
	Success bool        `json:"success"`
	Message string      `json:"message,omitempty"`
	Data    T           `json:"data,omitempty"`
	Error   *APIError   `json:"error,omitempty"`
	Meta    *Pagination `json:"meta,omitempty"` // Chỉ dùng Pagination, bỏ Meta trùng lặp
}

// APIResponseAny is a non-generic version for cases where type is unknown
type APIResponseAny struct {
	Success bool        `json:"success"`
	Message string      `json:"message,omitempty"`
	Data    any         `json:"data,omitempty"`
	Error   *APIError   `json:"error,omitempty"`
	Meta    *Pagination `json:"meta,omitempty"`
}

// APIError represents an API error
type APIError struct {
	Code    string         `json:"code"`
	Message string         `json:"message"`
	Details map[string]any `json:"details,omitempty"` // Go 1.18+
}

// ==================== Pagination ====================

// Pagination represents pagination information
type Pagination struct {
	Page       int  `json:"page"`
	PageSize   int  `json:"page_size"`
	TotalPages int  `json:"total_pages"`
	TotalCount int  `json:"total_count"`
	HasNext    bool `json:"has_next"`
	HasPrev    bool `json:"has_prev"`
}

// PaginationParams represents pagination query parameters
type PaginationParams struct {
	Page     int `form:"page" json:"page" binding:"omitempty,min=1"`
	PageSize int `form:"page_size" json:"page_size" binding:"omitempty,min=1,max=100"`
}

// GetOffset returns the offset for SQL queries (pure function)
func (p PaginationParams) GetOffset() int {
	page := p.Page
	if page <= 0 {
		page = 1
	}
	return (page - 1) * p.GetPageSize()
}

// GetPage returns the page with default value (pure function)
func (p PaginationParams) GetPage() int {
	if p.Page <= 0 {
		return 1
	}
	return p.Page
}

// GetPageSize returns the page size with default and max limit (pure function)
func (p PaginationParams) GetPageSize() int {
	if p.PageSize <= 0 {
		return 20 // Default page size
	}
	if p.PageSize > 100 {
		return 100 // Max page size
	}
	return p.PageSize
}

// GetLimit returns the limit for SQL queries (same as PageSize)
func (p PaginationParams) GetLimit() int {
	return p.GetPageSize()
}

// CalculatePagination calculates pagination metadata
func CalculatePagination(page, pageSize, totalCount int) *Pagination {
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 20
	}

	totalPages := 0
	if totalCount > 0 {
		totalPages = (totalCount + pageSize - 1) / pageSize
	}
	// Không set totalPages = 1 khi totalCount = 0, để semantics đúng

	return &Pagination{
		Page:       page,
		PageSize:   pageSize,
		TotalPages: totalPages,
		TotalCount: totalCount,
		HasNext:    page < totalPages,
		HasPrev:    page > 1,
	}
}

// HealthCheck represents a health check response
type HealthCheck struct {
	Status        HealthStatus `json:"status"`
	Timestamp     time.Time    `json:"timestamp"`
	UptimeSeconds int64        `json:"uptime_seconds"`
	Version       string       `json:"version"`
	Environment   string       `json:"environment,omitempty"`
}

// ==================== Simple Responses ====================

// IDResponse represents a response containing just an ID
type IDResponse struct {
	ID int `json:"id"`
}

// IDsResponse represents a response containing multiple IDs
type IDsResponse struct {
	IDs []int `json:"ids"`
}

// MessageResponse represents a simple message response
type MessageResponse struct {
	Message string `json:"message"`
}

// SuccessResponse represents a simple success response
type SuccessResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message,omitempty"`
}

// ==================== Validation Errors ====================

// ValidationError represents a field validation error
type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
	Value   any    `json:"value,omitempty"`
}

// ValidationErrors represents multiple validation errors
type ValidationErrors struct {
	Errors []*ValidationError `json:"errors"`
}

// AddError adds a validation error
func (v *ValidationErrors) AddError(field, message string, value any) {
	v.Errors = append(v.Errors, &ValidationError{
		Field:   field,
		Message: message,
		Value:   value,
	})
}

// HasErrors checks if there are any errors
func (v *ValidationErrors) HasErrors() bool {
	return len(v.Errors) > 0
}

// Error implements error interface
func (v *ValidationErrors) Error() string {
	if len(v.Errors) == 0 {
		return "no validation errors"
	}
	if len(v.Errors) == 1 {
		return fmt.Sprintf("%s: %s", v.Errors[0].Field, v.Errors[0].Message)
	}
	return fmt.Sprintf("%d validation errors", len(v.Errors))
}

// ToAPIError converts ValidationErrors to APIError
func (v *ValidationErrors) ToAPIError() *APIError {
	return &APIError{
		Code:    "VALIDATION_ERROR",
		Message: "Validation failed",
		Details: map[string]any{
			"errors": v.Errors,
		},
	}
}

// ==================== Error Response ====================

// ErrorResponse represents an error response
type ErrorResponse struct {
	Error   string         `json:"error"`
	Message string         `json:"message"`
	Details map[string]any `json:"details,omitempty"`
	Code    string         `json:"code,omitempty"`
}

// ==================== Sort Parameters ====================

// SortOrder represents sort order enum
type SortOrder string

const (
	SortOrderAsc  SortOrder = "asc"
	SortOrderDesc SortOrder = "desc"
)

// IsValid checks if sort order is valid
func (s SortOrder) IsValid() bool {
	return s == SortOrderAsc || s == SortOrderDesc
}

// SortParams represents sorting parameters
type SortParams struct {
	SortBy    string    `form:"sort_by" json:"sort_by" binding:"omitempty,max=50"`
	SortOrder SortOrder `form:"sort_order" json:"sort_order" binding:"omitempty,oneof=asc desc"`
}

// GetSortOrder returns the sort order with default (pure function)
func (s SortParams) GetSortOrder() SortOrder {
	if s.SortOrder == "" || !s.SortOrder.IsValid() {
		return SortOrderDesc
	}
	return s.SortOrder
}

// GetSortBy returns the sort field with validation against whitelist (pure function)
func (s SortParams) GetSortBy(defaultField string, allowedFields map[string]string) string {
	if s.SortBy == "" {
		return defaultField
	}

	// Validate and map to actual DB column
	if dbColumn, ok := allowedFields[s.SortBy]; ok {
		return dbColumn
	}

	return defaultField
}

// ValidateSortBy validates sort field against whitelist
func (s SortParams) ValidateSortBy(allowedFields map[string]string) error {
	if s.SortBy == "" {
		return nil
	}

	if _, ok := allowedFields[s.SortBy]; !ok {
		return fmt.Errorf("invalid sort_by field: %s", s.SortBy)
	}

	return nil
}

// ==================== Date Range Parameters ====================

// DateRangeParams represents date range query parameters
type DateRangeParams struct {
	StartDate *time.Time `form:"start_date" json:"start_date,omitempty" time_format:"2006-01-02"`
	EndDate   *time.Time `form:"end_date" json:"end_date,omitempty" time_format:"2006-01-02"`
}

// Validate validates date range
func (d DateRangeParams) Validate() error {
	if d.StartDate != nil && d.EndDate != nil {
		if d.EndDate.Before(*d.StartDate) {
			return fmt.Errorf("end_date must be after start_date")
		}
	}
	return nil
}

// GetStartDate returns start date or default (pure function)
func (d DateRangeParams) GetStartDate(defaultDaysAgo int) time.Time {
	if d.StartDate == nil {
		return time.Now().AddDate(0, 0, -defaultDaysAgo)
	}
	return *d.StartDate
}

// GetEndDate returns end date or default (now) (pure function)
func (d DateRangeParams) GetEndDate() time.Time {
	if d.EndDate == nil {
		return time.Now()
	}
	return *d.EndDate
}

// ==================== Search Parameters ====================

// SearchParams represents search query parameters
type SearchParams struct {
	Query string `form:"q" json:"q" binding:"omitempty,max=200"`
}

// HasQuery checks if search query is provided
func (s SearchParams) HasQuery() bool {
	return s.Query != ""
}

// GetQuery returns trimmed query
func (s SearchParams) GetQuery() string {
	return s.Query
}

// ==================== Filter Parameters (Combined) ====================

// FilterParams combines common filter parameters
type FilterParams struct {
	PaginationParams
	SortParams
	DateRangeParams
	SearchParams
}

// Validate validates all filter params
func (f FilterParams) Validate(allowedSortFields map[string]string) error {
	if err := f.DateRangeParams.Validate(); err != nil {
		return err
	}

	if err := f.SortParams.ValidateSortBy(allowedSortFields); err != nil {
		return err
	}

	return nil
}

// ==================== Location Parameters ====================

// LocationParams represents location query parameters
type LocationParams struct {
	Latitude  *float64 `form:"lat" json:"latitude,omitempty" binding:"omitempty,min=-90,max=90"`
	Longitude *float64 `form:"lng" json:"longitude,omitempty" binding:"omitempty,min=-180,max=180"`
	Radius    *int     `form:"radius" json:"radius,omitempty" binding:"omitempty,min=0,max=50000"` // meters
}

// HasLocation checks if location params are provided
func (l LocationParams) HasLocation() bool {
	return l.Latitude != nil && l.Longitude != nil
}

// GetLocation returns Point if valid
func (l LocationParams) GetLocation() *Point {
	if !l.HasLocation() {
		return nil
	}
	return &Point{
		Latitude:  *l.Latitude,
		Longitude: *l.Longitude,
	}
}

// Validate validates location params
func (l LocationParams) Validate() error {
	if (l.Latitude == nil) != (l.Longitude == nil) {
		return fmt.Errorf("both latitude and longitude must be provided")
	}

	if l.HasLocation() {
		point := l.GetLocation()
		if !point.IsValid() {
			return fmt.Errorf("invalid location coordinates")
		}
	}

	return nil
}

// ==================== Stats Parameters ====================

// StatsPeriod represents statistics time period enum
type StatsPeriod string

const (
	StatsPeriodToday     StatsPeriod = "today"
	StatsPeriodYesterday StatsPeriod = "yesterday"
	StatsPeriodWeek      StatsPeriod = "week"
	StatsPeriodMonth     StatsPeriod = "month"
	StatsPeriodYear      StatsPeriod = "year"
	StatsPeriodCustom    StatsPeriod = "custom"
)

// IsValid checks if stats period is valid
func (s StatsPeriod) IsValid() bool {
	switch s {
	case StatsPeriodToday, StatsPeriodYesterday, StatsPeriodWeek,
		StatsPeriodMonth, StatsPeriodYear, StatsPeriodCustom:
		return true
	default:
		return false
	}
}

// GetDateRange returns start and end date for the period
func (s StatsPeriod) GetDateRange() (time.Time, time.Time) {
	now := time.Now()
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())

	switch s {
	case StatsPeriodToday:
		return today, now
	case StatsPeriodYesterday:
		yesterday := today.AddDate(0, 0, -1)
		return yesterday, today
	case StatsPeriodWeek:
		weekAgo := today.AddDate(0, 0, -7)
		return weekAgo, now
	case StatsPeriodMonth:
		monthAgo := today.AddDate(0, -1, 0)
		return monthAgo, now
	case StatsPeriodYear:
		yearAgo := today.AddDate(-1, 0, 0)
		return yearAgo, now
	default:
		return today, now
	}
}

// StatsParams represents statistics query parameters
type StatsParams struct {
	Period    StatsPeriod `form:"period" json:"period,omitempty" binding:"omitempty,oneof=today yesterday week month year custom"`
	StartDate *time.Time  `form:"start_date" json:"start_date,omitempty" time_format:"2006-01-02"`
	EndDate   *time.Time  `form:"end_date" json:"end_date,omitempty" time_format:"2006-01-02"`
}

// GetDateRange returns date range based on period or custom dates
func (s StatsParams) GetDateRange() (time.Time, time.Time) {
	if s.Period == StatsPeriodCustom && s.StartDate != nil && s.EndDate != nil {
		return *s.StartDate, *s.EndDate
	}

	if s.Period != "" && s.Period.IsValid() {
		return s.Period.GetDateRange()
	}

	// Default to today
	return StatsPeriodToday.GetDateRange()
}

// Validate validates stats params
func (s StatsParams) Validate() error {
	if s.Period != "" && !s.Period.IsValid() {
		return fmt.Errorf("invalid period")
	}

	if s.Period == StatsPeriodCustom {
		if s.StartDate == nil || s.EndDate == nil {
			return fmt.Errorf("start_date and end_date are required for custom period")
		}
		if s.EndDate.Before(*s.StartDate) {
			return fmt.Errorf("end_date must be after start_date")
		}
	}

	return nil
}

// ==================== Bulk Response ====================

// BulkItemResult represents result of a single item in bulk operation
type BulkItemResult struct {
	ID      int     `json:"id"`
	Success bool    `json:"success"`
	Error   *string `json:"error,omitempty"`
}

// BulkResponse represents a bulk operation response
type BulkResponse struct {
	TotalRequested int               `json:"total_requested"`
	Successful     int               `json:"successful"`
	Failed         int               `json:"failed"`
	Results        []*BulkItemResult `json:"results,omitempty"`
}

// IsFullSuccess checks if all items succeeded
func (b *BulkResponse) IsFullSuccess() bool {
	return b.Failed == 0 && b.Successful == b.TotalRequested
}

// IsPartialSuccess checks if some items succeeded
func (b *BulkResponse) IsPartialSuccess() bool {
	return b.Successful > 0 && b.Failed > 0
}

// SuccessRate returns success rate percentage
func (b *BulkResponse) SuccessRate() float64 {
	if b.TotalRequested == 0 {
		return 0
	}
	return float64(b.Successful) / float64(b.TotalRequested) * 100
}

// ==================== Helper Functions ====================

// NewSuccessResponse creates a generic success response
func NewSuccessResponse[T any](message string, data T) *APIResponse[T] {
	return &APIResponse[T]{
		Success: true,
		Message: message,
		Data:    data,
	}
}

// NewSuccessResponseAny creates a non-generic success response
func NewSuccessResponseAny(message string, data any) *APIResponseAny {
	return &APIResponseAny{
		Success: true,
		Message: message,
		Data:    data,
	}
}

// NewPaginatedResponse creates a paginated success response
func NewPaginatedResponse[T any](message string, data T, pagination *Pagination) *APIResponse[T] {
	return &APIResponse[T]{
		Success: true,
		Message: message,
		Data:    data,
		Meta:    pagination,
	}
}

// NewErrorResponse creates a generic error response
func NewErrorResponse[T any](code, message string, details map[string]any) *APIResponse[T] {
	var zero T
	return &APIResponse[T]{
		Success: false,
		Data:    zero,
		Error: &APIError{
			Code:    code,
			Message: message,
			Details: details,
		},
	}
}

// NewErrorResponseAny creates a non-generic error response
func NewErrorResponseAny(code, message string, details map[string]any) *APIResponseAny {
	return &APIResponseAny{
		Success: false,
		Error: &APIError{
			Code:    code,
			Message: message,
			Details: details,
		},
	}
}

// NewValidationErrorResponse creates a validation error response
func NewValidationErrorResponse[T any](errors *ValidationErrors) *APIResponse[T] {
	var zero T
	return &APIResponse[T]{
		Success: false,
		Data:    zero,
		Error:   errors.ToAPIError(),
	}
}

// NewValidationErrorResponseAny creates a non-generic validation error response
func NewValidationErrorResponseAny(errors *ValidationErrors) *APIResponseAny {
	return &APIResponseAny{
		Success: false,
		Error:   errors.ToAPIError(),
	}
}
