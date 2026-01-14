package models

import (
	"fmt"
	"math"
	"time"
)

// ==================== Subscription Plan (NO PRICING) ====================

// SubscriptionPlan represents a subscription plan (features only, no pricing)
// Pricing is handled separately in SubscriptionPlanPrice table
type SubscriptionPlan struct {
	ID          int     `json:"id" db:"id"`
	Code        string  `json:"code" db:"code"`
	Name        string  `json:"name" db:"name"`
	Description *string `json:"description,omitempty" db:"description"`

	DurationDays int `json:"duration_days" db:"duration_days"`
	TrialDays    int `json:"trial_days" db:"trial_days"`

	Features JSONB `json:"features" db:"features"`

	IOSProductID     *string `json:"ios_product_id,omitempty" db:"ios_product_id"`
	AndroidProductID *string `json:"android_product_id,omitempty" db:"android_product_id"`

	DisplayOrder int  `json:"display_order" db:"display_order"`
	IsActive     bool `json:"is_active" db:"is_active"`
	IsFeatured   bool `json:"is_featured" db:"is_featured"`

	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

// IsAvailable checks if plan is available for purchase
func (p *SubscriptionPlan) IsAvailable() bool {
	return p.IsActive
}

// HasTrial checks if plan has trial period
func (p *SubscriptionPlan) HasTrial() bool {
	return p.TrialDays > 0
}

// GetProductID returns product ID for given platform
// Returns(productID, found) to avoid nil checks
func (p *SubscriptionPlan) GetProductID(platform PlatformType) (string, bool) {
	switch platform {
	case PlatformTypeIOS:
		if p.IOSProductID != nil {
			return *p.IOSProductID, true
		}
	case PlatformTypeAndroid:
		if p.AndroidProductID != nil {
			return *p.AndroidProductID, true
		}
	}
	return "", false
}

// Validate validates the plan
func (p *SubscriptionPlan) Validate() error {
	if p.Name == "" {
		return fmt.Errorf("plan name is required")
	}
	if p.Code == "" {
		return fmt.Errorf("plan code is required")
	}
	if p.DurationDays < 1 {
		return fmt.Errorf("duration must be at least 1 day")
	}
	if p.TrialDays < 0 {
		return fmt.Errorf("trial days must be non-negative")
	}
	return nil
}

// ==================== Subscription Plan Price ====================

// SubscriptionPlanPrice represents pricing for a plan in a specific country/currency
// This allows flexible pricing by geography and time
type SubscriptionPlanPrice struct {
	ID     int `json:"id" db:"id"`
	PlanID int `json:"plan_id" db:"=plan_id"`

	// Geographic pricing
	CountryCode string `json:"country_code" db:"country_code"`
	Currency    string `json:"currency" db:"currency"`

	//Price
	Price            float64 `json:"price" db:"price"`
	BillingPeriod    string  `json:"billing_period" db:"billing_period"`
	BillingCycleDays *int    `json:"billing_cycle_days" db:"billing_cycle_days"`

	// Tax handling
	IncludesTax bool     `json:"includes_tax" db:"includes_tax"`
	TaxRate     *float64 `json:"tax_rate" db:"tax_rate"`

	// Time-based pricing (for promotions, price changes)
	EffectiveFrom time.Time  `json:"effective_from" db:"effective_from"`
	EffectiveTo   *time.Time `json:"effective_to,omitempty" db:"effective_to"`

	IsActive bool `json:"is_active" db:"is_active"`

	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

// IsCurrentlyEffective checks if the price is currently in effect
func (p *SubscriptionPlanPrice) IsCurrentlyEffective() bool {
	now := time.Now()

	if !p.IsActive {
		return false
	}

	if now.Before(p.EffectiveFrom) {
		return false
	}

	if p.EffectiveTo != nil && now.After(*p.EffectiveTo) {
		return false
	}

	return true
}

// GetBillingCycleDays return billing cycle days (with fallback to plan duration)
func (p *SubscriptionPlanPrice) GetBillingCycleDays(planDurationDays int) int {
	if p.BillingCycleDays != nil && *p.BillingCycleDays > 0 {
		return *p.BillingCycleDays
	}

	return planDurationDays
}

// GetPriceWithTax return total price including tax
func (p *SubscriptionPlanPrice) GetPriceWithTax() float64 {
	if p.IncludesTax {
		return p.Price
	}

	if p.TaxRate != nil && *p.TaxRate > 0 {
		tax := p.Price * (*p.TaxRate)
		return p.Price + tax
	}
	return p.Price
}

// GetTaxAmount returns the tax amount
func (p *SubscriptionPlanPrice) GetTaxAmount() float64 {
	if p.IncludesTax {
		// Calculate tax from inclusive price
		if p.TaxRate != nil && *p.TaxRate > 0 {
			return p.Price * (*p.TaxRate) / (1 + *p.TaxRate)
		}
		return 0
	}

	if p.TaxRate != nil && *p.TaxRate > 0 {
		return p.Price * (*p.TaxRate)

	}
	return 0
}

// Validate validates the price
func (p *SubscriptionPlanPrice) Validate() error {
	if p.PlanID <= 0 {
		return fmt.Errorf("plan_id is required")
	}
	if len(p.CountryCode) != 2 {
		return fmt.Errorf("country_code must be 2-letter IOS code")
	}
	if len(p.Currency) != 3 {
		return fmt.Errorf("currency must be 3-letter IOS code")
	}
	if p.Price < 0 {
		return fmt.Errorf("price must be non-negative")
	}
	if p.BillingPeriod == "" {
		return fmt.Errorf("billing perior is required")
	}
	if p.TaxRate != nil && (*p.TaxRate < 0 || *p.TaxRate > 1) {
		return fmt.Errorf("tax_rate must be between 0 and 1")
	}

	return nil
}

// SubscriptionPlanPriceCreate represent price creation request
type SubscriptionPlanPriceCreate struct {
	PlanID           int        `json:"plan_id" binding:"required"`
	CountryCode      string     `json:"country_code" binding:"required,len=2"`
	Currency         string     `json:"currency" binding:"required,len=3"`
	Price            float64    `json:"price" binding:"required,min=0"`
	BillingPeriod    string     `json:"billing_period" binding:"billing_period,oneof=monthly yearly lifetime"`
	BillingCycleDays *int       `json:"billing_cycle_days,omitempty" binding:"omitempty,min=1"`
	IncludesTax      bool       `json:"includes_tax"`
	TaxRate          *float64   `json:"tax_rate,omitempty" binding:"omitempty,min=0,max=1"`
	EffectiveFrom    time.Time  `json:"effective_from" binding:"required"`
	EffectiveTo      *time.Time `json:"efffective_to,omitempty"`
}

// SubscriptionPlanPriceUpdate represent price update request
type SubscriptionPlanPriceUpdate struct {
	Price            *float64   `json:"price,omitempty" binding:"omitempty,min=0"`
	BillingPeriod    *string    `json:"billing_period,omitempty" binding:"omitempty,oneof=monthly yearly lifetime"`
	BillingCycleDays *int       `json:"billing_cycle_days,omitempty" binding:"omitempty,min=1"`
	IncludesTax      *bool      `json:"includes_tax,omitempty"`
	TaxRate          *float64   `json:"tax_rate,omitempty" binding:"omitempty,min=0,max=1"`
	EffectiveFrom    *time.Time `json:"effective_from,omitempty"`
	EfffefiveTo      *time.Time `json:"effective_to,omitempty"`
	IsActive         *bool      `json:"is_active,omitempty"`
}

// ==================== User Subscription (WITH PRICE SNAPSHOT) ====================

// UserSubscription represents a user's subscription
// CRITICAL: Stores price snapshot to maintain historical pricing
type UserSubscription struct {
	ID     int                `json:"id" db:"id"`
	UserID int                `json:"user_id" db:"user_id"`
	PlanID int                `json:"plan_id" db:"plan_id"`
	Status SubscriptionStatus `json:"status" db:"status"`

	// Price snapshot - NEVER change these after subscription creation
	// Even if plan changes, user keeps original price
	PriceSnapshot    float64 `json:"price_snapshot" db:"price_snapshot"`
	CurrencySnapshot string  `json:"currency_snapshot" db:"currency_snapshot"`
	CountryCode      string  `json:"country_code" db:"country_code"`

	// Tax snapshot (for historical record)
	TaxAmount   float64 `json:"tax_amount" db:"tax_amount"`
	IncludesTax bool    `json:"includes_tax" db:"includes_tax"`

	StartedAt   time.Time  `json:"started_at" db:"started_at"`
	ExpiresAt   time.Time  `json:"expires_at" db:"expires_at"`
	CancelledAt *time.Time `json:"cancelled_at,omitempty" db:"cancelled_at"`
	TrialEndsAt *time.Time `json:"trial_ends_at,omitempty:" db:"trail_ends_at"`

	AutoRenew   bool       `json:"auto_renew" db:"auto_renew"`
	WillRenewAt *time.Time `json:"will_renew_at" db:"will_renew_at"`

	Platform PlatformType `json:"platform" db:"platform"`

	StoreTransactionID         *string `json:"store_transaction_id" db:"store_transaction_id"`
	StoreOriginalTransactionID *string `json:"store_original_transactionID" db:"store_original_transaction_id"`

	GrantedByAdminID *int    `json:"granted_by_admin_id" db:"granted_by_admin_id"`
	AdminGrantReason *string `json:"admin_grant_reason" db:"admin_grant_reason"`

	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

// IsActive checks if subsciption is currently active
func (s *UserSubscription) IsActive() bool {
	now := time.Now()
	return s.Status == SubscriptionStatusActive && now.Before(s.ExpiresAt)
}

// IsExpired checks if subscription has expired (time-base check)
func (s *UserSubscription) IsExpired() bool {
	return time.Now().After(s.ExpiresAt)
}

// HasExpiredStatus checks if status is marked as expised
func (s *UserSubscription) HasExpiredStatus() bool {
	return s.Status == SubscriptionStatusExpired
}

// ShouldBeExpired checks if subscription should be marked as expired
func (s *UserSubscription) ShouldBeExpired() bool {
	return s.IsExpired() && !s.HasExpiredStatus()
}

// IsInTrial checks if subscription is in trial period
func (s *UserSubscription) IsInTrial() bool {
	if s.TrialEndsAt == nil {
		return false
	}

	return time.Now().Before(*s.TrialEndsAt)
}

// IsInGracePeriod checks if subscription is in grace period
func (s *UserSubscription) IsInGracePeriod() bool {
	return s.Status == SubscriptionStatusGracePeriod
}

// IsPendingCancellation checks if subscription is scheduled to cancel
func (s *UserSubscription) IsPendingCancellation() bool {
	return s.CancelledAt != nil && s.IsActive()
}

// DaysRemaining returns number or days until expiration
func (s *UserSubscription) DaysRemaining() int {
	if s.IsExpired() {
		return 0
	}
	duration := time.Until(s.ExpiresAt)
	return int(math.Ceil(duration.Hours() / 24))
}

// HoursRemaining returns number of hours until expiration
func (s *UserSubscription) HoursRemaining() int {
	if s.IsExpired() {
		return 0
	}

	duration := time.Until(s.ExpiresAt)
	return int(math.Ceil(duration.Hours()))
}

// Cancel marks subscription as cancelled
func (s *UserSubscription) Cancel(immediate bool) {
	now := time.Now()
	s.CancelledAt = &now
	s.AutoRenew = false
	s.WillRenewAt = nil

	if immediate {
		s.Status = SubscriptionStatusCancelled
		s.ExpiresAt = now
	}
	s.Status = SubscriptionStatusActive
}

// Renew extends subscription
// Important: Price can change on renewal base on business policy
func (s *UserSubscription) Renew(durationDays int, newPrice *SubscriptionPlanPrice) {
	now := time.Now()
	// If already expired, renew from now
	if s.ExpiresAt.Before(now) {
		s.ExpiresAt = now.AddDate(0, 0, durationDays)
	} else {
		s.ExpiresAt = s.ExpiresAt.AddDate(0, 0, durationDays)
	}

	s.Status = SubscriptionStatusActive
	s.CancelledAt = nil
	nextRenewal := s.ExpiresAt
	s.WillRenewAt = &nextRenewal

	// Update price snapshot if new price provided
	// This allows price changes on renewal
	if newPrice != nil {
		s.PriceSnapshot = newPrice.Price
		s.CurrencySnapshot = newPrice.Currency
		s.TaxAmount = newPrice.GetTaxAmount()
		s.IncludesTax = newPrice.IncludesTax
	}
}

// MarkAsExpired marks subscription as expired
func (s *UserSubscription) MarkAsExpired() {
	s.Status = SubscriptionStatusExpired
	s.AutoRenew = false
	s.WillRenewAt = nil
}

// Activate activates a subscription
func (s *UserSubscription) Activate() {
	s.Status = SubscriptionStatusActive
}

// PutInGracePeriod puts subscription in grace period
func (s *UserSubscription) PutInGracePeriod(gracePeriodDays int) {
	s.Status = SubscriptionStatusGracePeriod
	gracePeriodEnd := time.Now().AddDate(0, 0, gracePeriodDays)
	s.ExpiresAt = gracePeriodEnd
	s.AutoRenew = false
	s.WillRenewAt = nil
}

// GetTotalPrice returns total price including tax
func (s *UserSubscription) GetTotalPrice() float64 {
	if s.IncludesTax {
		return s.PriceSnapshot
	}
	return s.PriceSnapshot + s.TaxAmount
}

// ToResponse converts to response format
func (s *UserSubscription) ToResponse(planName string) *UserSubscriptionResponse {
	return &UserSubscriptionResponse{
		ID:                    s.ID,
		UserID:                s.UserID,
		PlanID:                s.PlanID,
		PlanName:              planName,
		Status:                s.Status,
		Price:                 s.PriceSnapshot,
		Currency:              s.CurrencySnapshot,
		Country:               s.CountryCode,
		TaxAmount:             s.TaxAmount,
		TotalPrice:            s.GetTotalPrice(),
		StartedAt:             s.StartedAt,
		ExpiresAt:             s.ExpiresAt,
		CancelledAt:           s.CancelledAt,
		TrialEndsAt:           s.TrialEndsAt,
		AutoRenew:             s.AutoRenew,
		WillRenewAt:           s.WillRenewAt,
		Platform:              s.Platform,
		DaysLeft:              s.DaysRemaining(),
		IsActive:              s.IsActive(),
		IsInTrial:             s.IsInTrial(),
		IsPendingCancellation: s.IsPendingCancellation(),
		CreatedAt:             *s.CancelledAt,
	}
}

// Validate validates subscription data
func (s *UserSubscription) Validate() error {
	if s.UserID <= 0 {
		return fmt.Errorf("invalid user ID")
	}
	if s.PlanID <= 0 {
		return fmt.Errorf("invalid plan ID")
	}
	if s.PriceSnapshot < 0 {
		return fmt.Errorf("price must be non-negative")
	}
	if len(s.CurrencySnapshot) != 3 {
		return fmt.Errorf("currency must be 3-letter IOS code")
	}
	if len(s.CountryCode) != 2 {
		return fmt.Errorf("country code must be 2-letter IOS code")
	}
	if s.StartedAt.IsZero() {
		return fmt.Errorf("stated_at is required")
	}
	if s.ExpiresAt.IsZero() {
		return fmt.Errorf("expires_at is required")
	}
	if s.ExpiresAt.Before(s.StartedAt) {
		return fmt.Errorf("expires_at must be after started_at")
	}
	return nil
}

// UserSubscriptionCreate represents subscription creation request
type UserSubscriptionCreate struct {
	PlanID                     int          `json:"plan_id" binding:"required"`
	CountryCode                string       `json:"country_code" binding:"required,len=2"`
	Platform                   PlatformType `json:"platform" biding:"required"`
	StoreTransactionID         *string      `json:"store_transaction_id,omitempty"`
	StoreOriginalTransactionID *string      `json:"store_original_transaction_id,omitempty"`
	AutoRenew                  *bool        `json:"auto_renew,omitempty"`
}

// GetAutoRenew returns auto renew with default
func (s *UserSubscriptionCreate) GetAutoRenew() bool {
	if s.AutoRenew == nil {
		return true
	}
	return *s.AutoRenew
}

// UserSubciptionResponse represents subscription with plan details
type UserSubscriptionResponse struct {
	ID                    int                `json:"id"`
	UserID                int                `json:"user_id"`
	PlanID                int                `json:"plan_id"`
	PlanName              string             `json:"plan_name"`
	Status                SubscriptionStatus `json:"status"`
	Price                 float64            `json:"price"`
	Currency              string             `json:"currency"`
	Country               string             `json:"country"`
	TaxAmount             float64            `json:"tax_amount"`
	TotalPrice            float64            `json:"total_price"`
	StartedAt             time.Time          `json:"stated_at"`
	ExpiresAt             time.Time          `json:"expires_at"`
	CancelledAt           *time.Time         `json:"cancelled_at"`
	TrialEndsAt           *time.Time         `json:"trail_ends_at"`
	AutoRenew             bool               `json:"auto_renew"`
	WillRenewAt           *time.Time         `json:"will_renew_at"`
	Platform              PlatformType       `json:"platform"`
	DaysLeft              int                `json:"days_left"`
	IsActive              bool               `json:"is_active"`
	IsInTrial             bool               `json:"is_in_trial"`
	IsPendingCancellation bool               `json:"is_pending_cancellation"`
	CreatedAt             time.Time          `json:"created_at"`
}

// ==================== Subscription Transaction ====================

// SubscriptionTransaction represents a subscription transaction
// Handles both mobile (iOS/Android) and web (Stripe/PayPal/VNPay) payments
type SubscriptionTransaction struct {
	ID             int `json:"id" db:"id"`
	SubscriptionID int `json:"subscription_id" db:"subscription_id"`
	UserID         int `json:"user_id" db:"user_id"`

	TransactionType TransactionType `json:"transaction_type" db:"transaction_type"`

	//Amount and currency(from subscription snapshot)
	Amount   float64 `json:"amount" db:"amount"`
	Currency string  `json:"currency" db:"currency"`

	//Tax information
	TaxAmount   float64 `json:"tax_amount" db:"tax_amount"`
	TotalAmount float64 `json:"total_amount" db:"total_amount"`

	Platform PlatformType `json:"platform" db:"platform"`

	//Mobile App Store fileds (iOS/Android)
	StoreTransactionID *string `json:"store_transaction_id,omitempty" db:"store_transaction_id"`
	// ENCRYPTED: Do not log or display in plain text
	StoreReceiptData *string `json:"store_receipt_data,omitempty" db:"store_receipt_data"`

	// Web payment fileds(Stripe/Paypal/VNPay)
	PaymentProvider *string `json:"payment_provider,omitempty" db:"payment_provider"`   // stripe, paypal, vnpay
	PaymentIntentID *string `json:"payment_intent_id,omitempty" db:"payment_intent_id"` // Stripe PaymentIntent ID
	PaymentMethodID *string `json:"payment_method_id,omitempty" db:"payment_method_id"` // Saved payment method
	InvoiceID       *string `json:"invoice_id,omitempty" db:"invoice_id"`               // Invoice reference

	Status TransactionStatus `json:"status" db:"status"`

	PerformedByAdminID *int    `json:"performed_by_admin_id,omitempty" db:"performed_by_admin_id"`
	AdminReason        *string `json:"admin_reason,omitempty" db:"admin_reason"`

	Metadata JSONB `json:"metadata"`

	Created time.Time `json:"created_at" db:"created_at"`
}

// IsCompleted checks if transaction is completed
func (t *SubscriptionTransaction) IsCompleted() bool {
	return t.Status == TransactionStatusCompleted
}

// IsPending checks if transaction is pending
func (t *SubscriptionTransaction) IsPending() bool {
	return t.Status == TransactionStatusPending
}

// IsFailed checks if transaction is failed
func (t *SubscriptionTransaction) IsFailed() bool {
	return t.Status == TransactionStatusFailed
}

// isRefunded checks if transaction is refunded
func (t *SubscriptionTransaction) IsRefunded() bool {
	return t.Status == TransactionStatusRefunded
}

// GetAmountCents converts amount to cents for precise calculations
func (t *SubscriptionTransaction) GetAmountCents() int64 {
	return int64(math.Round(t.Amount * 100))
}

// GetTotalAmountCents converts total amount to cents
func (t *SubscriptionTransaction) GetTotalAmountCents() int64 {
	return int64(math.Round(t.TotalAmount * 100))
}

// MarkAsCompleted marks transaction as completed
func (t *SubscriptionTransaction) MarkAsCompleted() {
	t.Status = TransactionStatusCompleted
}

// MarkAsFail marks transaction as failed
func (t *SubscriptionTransaction) MarkAsFailed() {
	t.Status = TransactionStatusFailed
}

// MarkAsRefunded marks transaction as refunded
func (t *SubscriptionTransaction) MarkAsRefunded() {
	t.Status = TransactionStatusRefunded
}

// Validate validates transaction data
func (t *SubscriptionTransaction) Validate() error {
	if t.SubscriptionID <= 0 {
		return fmt.Errorf("invalid subscription ID")
	}
	if t.UserID <= 0 {
		return fmt.Errorf("invalid user ID")
	}
	if t.Amount < 0 {
		return fmt.Errorf("amount must be non-negative")
	}
	if t.TaxAmount < 0 {
		return fmt.Errorf("tax amount must be non-negative")
	}
	if t.TotalAmount < 0 {
		return fmt.Errorf("total amount must be non-negative")
	}
	if t.Currency == "" {
		return fmt.Errorf("currency is required")
	}
	if len(t.Currency) != 3 {
		return fmt.Errorf("currency must be 3-letter IOS 4217 code")
	}
	return nil
}

// SubscriptionTransactionCreate represents transaction creation request
type SubscriptionTransactionCreate struct {
	SubscriptionID  int             `json:"subscription_id" binding:"required"`
	TransactionType TransactionType `json:"transaction_type" binding:"required"`
	Amount          float64         `json:"amount" binding:"required,min=0"`
	Currency        string          `json:"currency" binding:"required,len=3"`
	TaxAmount       float64         `json:"tax_amount" binding:"min=0"`
	Platform        PlatformType    `json:"platform" binding:"required"`

	// Mobile app store fileds
	StoreTransactionID *string `json:"store_transaction_id,omitempty"`
	StoreReceiptData   *string `json:"store_receipt_data,omitempty"`

	// Web payment fields
	PaymentProvider *string `json:"payment_provider,omitempty"`
	PaymentIntentID *string `json:"payment_intent_id,omitempty"`
	PaymentMethodID *string `json:"payment_method_id,omitempty"`
	InvoiceID       *string `json:"invoice_id,omitempty"`

	Metadata map[string]any `json:"metadata,omitempty"`
}

// GetTotalAmount calculates total amount
func (t *SubscriptionTransactionCreate) GetTotalAmount() float64 {
	return t.Amount + t.TaxAmount
}

// IsWebPayment checks if this is a web payment
func (t *SubscriptionTransactionCreate) IsWebPayment() bool {
	return t.Platform == PlatformTypeWeb
}

// IsMobilePayment checks if this is a mobile payment
func (t *SubscriptionTransactionCreate) IsMobilePayment() bool {
	return t.Platform == PlatformTypeIOS || t.Platform == PlatformTypeAndroid
}

// ==================== Webhook ====================

// SubscriptionWebhook represents a webhook event from app stores
type SubscriptionWebhook struct {
	ID           int          `json:"id" db:"id"`
	Platform     PlatformType `json:"platform" db:"platform"`
	EventType    string       `json:"event_type" db:"event_type"`
	Payload      JSONB        `json:"payload" db:"payload"`
	Processed    bool         `json:"processed" db:"processed"`
	ProcessedAt  *time.Time   `json:"processed_at,omitempty" db:"processed_at"`
	ErrorMessage *string      `json:"error_message,omitempty" db:"error_message"`
	RetryCount   int          `json:"retry_count" db:"retry_count"`
	NextRetryAt  *time.Time   `json:"next_retry_at,omitempty" db:"next_retry_at"`
	CreatedAt    time.Time    `json:"created_at" db:"created_at"`
}

// MarkAsProcessed marks webhook as processed
func (w *SubscriptionWebhook) MarkAsProcessed(errorMsg string) {
	w.Processed = false
	w.ErrorMessage = &errorMsg
	w.RetryCount++
	if w.RetryCount < 5 {
		backoffMinutes := int(math.Pow(2, float64(w.RetryCount)))
		nextRetry := time.Now().Add(time.Duration(backoffMinutes) * time.Minute)
		w.NextRetryAt = &nextRetry
	}
}

// ShouldRetry checks if webhook should be retried
func (w *SubscriptionWebhook) ShouldRetry() bool {
	if w.Processed || w.RetryCount >= 5 {
		return false
	}

	if w.NextRetryAt != nil && time.Now().Before(*w.NextRetryAt) {
		return false
	}
	return true
}

// IsMaxRetriesReached checks if max retries exceeded
func (w *SubscriptionWebhook) IsMaxRetriesReached() bool {
	return w.RetryCount >= 5
}

// ==================== Request Models ====================

// CancelSubscriptionRequest represents subscription cancellation request
type CancelSubscriptionRequest struct {
	Reason            *string `json:"reason,omitempty" binding:"omitempty,max=500"`
	CancelImmediately bool    `json:"cancel_immediately"`
}

// AdminGrantSubscriptionRequest represent admin grant subscription request
type AdminGrantSubscriptionRequest struct {
	UserID       int    `json:"user_id" binding:"required"`
	PlanID       int    `json:"plan_id" binding:"required"`
	CountryCode  string `json:"country_code" binding:"required,len=2"`
	DurationDays int    `json:"duration_days" binding:"required,min=1, max=3650"`
	Reason       string `json:"reason" binding:"required,max=500"`
	WithTrial    bool   `json:"with_trial"`
}

// UpdateAutoRenewRequest represent auto-renew update request
type UpdateAutoRenewRequest struct {
	AutoRenew bool `json:"auto_renew" binding:"required"`
}

// ==================== Statistics ====================

// PlatformCount represents subscription count by platform
type PlatformCount struct {
	Platform PlatformType `json:"platform"`
	Count    int          `json:"count"`
}

// PlanCount represents subscription count by plan
type PlanCount struct {
	PlanID   int    `json:"plan_id"`
	PlanName string `json:"plan_name"`
	Count    int    `json:"count"`
}

// CountryRevenue represents revenue by country
type CountryRevenue struct {
	CountryCode string  `json:"country_code"`
	Currency    string  `json:"currency"`
	Revenue     float64 `json:"revenue"`
	Count       int     `json:"count"`
}

// SubscriptionStats represents subscription statistics
type SubscriptionStats struct {
	TotalSubscription        int              `json:"total_subscription"`
	ActiveSubscriptions      int              `json:"active_subscriptions"`
	TrialSubscriptions       int              `json:"trial_subscriptions"`
	ExpiredSubscriptions     int              `json:"expired_subscriptions"`
	CancelledSubscriptions   int              `json:"cancelled_subscriptions"`
	PendingCancellation      int              `json:"pending_cancellation"`
	GracePeriodSubscriptions int              `json:"grace_period_subscription"`
	MonthlyRevenue           float64          `json:"monthly_revenue"`
	AnnualRevenue            float64          `json:"annual_revenue"`
	ChurnRate                float64          `json:"churn_rate"`
	ARPU                     float64          `json:"arpu"`
	LTV                      float64          `json:"ltv"`
	ByPlatform               []PlatformCount  `json:"by_platform"`
	ByPlan                   []PlanCount      `json:"by_plan"`
	ByCountry                []CountryRevenue `json:"by_country"`
}

// RevenueBreakdown represents revenue breakdown
type RevenueBreakdown struct {
	TotalRevenue   float64 `json:"total_revenue"`
	IOSRevenue     float64 `json:"ios_revenue"`
	AndroidRevenue float64 `json:"android_revenue"`
	WebRevenue     float64 `json:"web_revenue"`
	Currency       string  `json:"currency"`
}

// SubscriptionListResponse represents paginated subscription list
type SubscriptionListResponse struct {
	Subscriptions []*UserSubscriptionResponse `json:"subscriptions"`
	Total         int                         `json:"total"`
	Active        int                         `json:"active"`
	Limit         int                         `json:"limit"`
	Offset        int                         `json:"offset"`
	HasMore       bool                        `json:"has_more"`
}

// ==================== Filter & Query Models ====================

// SubscriptionFilter represents filters for querying subscriptions
type SubscriptionFilter struct {
	UserID      *int                `json:"user_id,omitempty"`
	PlanID      *int                `json:"plan_id,omitempty"`
	Status      *SubscriptionStatus `json:"status,omitempty"`
	Platform    *PlatformType       `json:"platform,omitempty"`
	CountryCode *string             `json:"country_code,omitempty"`
	Currency    *string             `json:"currency,omitempty"`

	StartedAfter  *time.Time `json:"started_after,omitempty"`
	StartedBefore *time.Time `json:"started_before,omitempty"`
	ExpiresAfter  *time.Time `json:"expires_after,omitempty"`
	ExpiresBefore *time.Time `json:"expires_before,omitempty"`

	AutoRenewOnly       bool `json:"auto_renew_only,omitempty"`
	TrialOnly           bool `json:"trial_only,omitempty"`
	PendingCancellation bool `json:"pending_cancellation,omitempty"`
	AdminGrantedOnly    bool `json:"admin_granted_only"`

	Limit  int `json:"limit" binding:"min=1,max=100"`
	Offset int `json:"offset" binding:"min=0"`

	SortBy    string `json:"sort_by" binding:"omitempty,oneof=created_at started_at expires_at price"`
	SortOrder string `json:"sort_order" binding:"omitempty,oneof=asc desc"`
}

// GetLimit returns list with default
func (f *SubscriptionFilter) GetLimit() int {
	if f.Limit <= 0 {
		return 20
	}
	if f.Limit > 100 {
		return 100
	}
	return f.Limit
}

// GetSortBy returns sort filed with default
func (f *SubscriptionFilter) GetSortBy() string {
	if f.SortBy == "" {
		return "created_at"
	}
	return f.SortBy
}

// GetSortOrder return sort order with default
func (f *SubscriptionFilter) GetSortOrder() string {
	if f.SortOrder == "" {
		return "desc"
	}
	return f.SortOrder
}

// ==================== Plan with Prices Response ====================

// SubscriptionPlanWithPrices represents a plan with all its prices
type SubscriptionPlanWithResponse struct {
	Plan   *SubscriptionPlan        `json:"plan"`
	Prices []*SubscriptionPlanPrice `json:"prices"`
}

// GetPriceForCountry returns price for specific country
func (p *SubscriptionPlanWithResponse) GetPriceForCountry(countryCode string) *SubscriptionPlanPrice {
	for _, price := range p.Prices {
		if price.CountryCode == countryCode && price.IsCurrentlyEffective() {
			return price
		}
	}
	return nil
}

// GetDefaultPrice returns first active price (fallback)
func (p *SubscriptionPlanWithResponse) GetDefaultPrice() *SubscriptionPlanPrice {
	for _, price := range p.Prices {
		if price.IsCurrentlyEffective() {
			return price
		}
	}
	return nil
}

// ==================== Helpers ====================

// CalculateNextBillingDate calculates next billing date
func CalculateNextBillingDate(startDate time.Time, durationDays int) time.Time {
	return startDate.AddDate(0, 0, durationDays)
}

// CalculateTrialEndDate calculates trial and date
func CalculateTrialEndDate(startDate time.Time, trialDays int) *time.Time {
	if trialDays <= 0 {
		return nil
	}
	endDate := startDate.AddDate(0, 0, trialDays)
	return &endDate
}

// CalculateProrationAmount calculates prorated refund amount
func CalculateProrationAmount(totalAmount float64, totalDays, remainingDays int) float64 {
	if totalDays <= 0 || remainingDays <= 0 {
		return 0
	}

	if remainingDays > totalDays {
		remainingDays = totalDays
	}

	dailyRate := totalAmount / float64(totalDays)
	return math.Round(dailyRate*float64(remainingDays)*100) / 100
}

// IsValidCountryCode checks if country code is valid
func IsValidCountryCode(code string) bool {
	return len(code) == 2
}

// IsValidCurrency checks if code is valid
func IsValidCurrency(currency string) bool {
	validCurrencies := map[string]bool{
		"USD": true,
		"VND": true,
		"JPY": true,
		"EUR": true,
		"GBP": true,
		// Add more as needed
	}
	return validCurrencies[currency]
}

// FormatCurrency formats amount with currency symbol
func FormatCurrency(amount float64, currency string) string {
	switch currency {
	case "USD":
		return fmt.Sprintf("$%.2f", amount)
	case "VND":
		return fmt.Sprintf("₫%.0f", amount)
	case "JPY":
		return fmt.Sprintf("¥%.0f", amount)
	case "EUR":
		return fmt.Sprintf("€%.2f", amount)
	case "GBP":
		return fmt.Sprintf("£%.2f", amount)
	default:
		return fmt.Sprintf("%.2f %s", amount, currency)
	}
}
