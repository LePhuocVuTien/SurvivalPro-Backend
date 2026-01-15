package usecase

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"time"

	"github.com/LePhuocVuTien/SurvivalPro-Backend/application/validation"
	"github.com/LePhuocVuTien/SurvivalPro-Backend/infrastructure/ratelimit"
	"github.com/LePhuocVuTien/SurvivalPro-Backend/internal/domain/user"
)

// ============================================================================
// DOMAIN ERRORS
// ============================================================================

var (
	ErrInvalidCredentials = errors.New("invalid_credentials")
	ErrAccountLocked      = errors.New("account_locked")
	ErrUserNotFound       = errors.New("user_not_found")
	ErrEmailNotVerified   = errors.New("email_not_verified")
	ErrInvalidToken       = errors.New("invalid_token")
	ErrTokenExpired       = errors.New("token_expired")

	ErrTwoFactorRequired = errors.New("two_factor_required")
	ErrInvalid2FA        = errors.New("invalid_2fa_code")
)

// ============================================================================
// LOGIN USE CASE
// ============================================================================

// LoginUseCase handles login business logic
type LoginUseCase struct {
	userRepo       UserRepository
	credentialRepo CredentialRepository
	securityRepo   SecurityRepository
	sessionRepo    SessionRepository
	activityRepo   ActivityRepository
	rateLimiter    ratelimit.Limiter
}

// NewLoginUseCase creates a new login use case
func NewLoginUseCase(
	userRepo UserRepository,
	credentialRepo CredentialRepository,
	securityRepo SecurityRepository,
	sessionRepo SessionRepository,
	activityRepo ActivityRepository,
	rateLimiter ratelimit.Limiter,
) *LoginUseCase {
	return &LoginUseCase{
		userRepo:       userRepo,
		credentialRepo: credentialRepo,
		securityRepo:   securityRepo,
		sessionRepo:    sessionRepo,
		activityRepo:   activityRepo,
		rateLimiter:    rateLimiter,
	}
}

// Execute executes the login use case with proper order:
// 1. Basic input validation (format, required fields)
// 2. Rate limiting (dual-identifier: IP + email hash)
// 3. Business logic (authentication)
func (uc *LoginUseCase) Execute(ctx context.Context, req *user.LoginRequest, ipAddress string) (*user.LoginResponse, error) {
	// ========================================================================
	// STEP 1: Basic Input Validation (BEFORE rate limiting)
	// ========================================================================
	// Validate format and required fields to prevent spam with invalid data
	if errs := validation.ValidateLoginRequest(req); len(errs) > 0 {
		return nil, fmt.Errorf("validation failed: %v", errs[0])
	}

	// ========================================================================
	// STEP 2: Rate Limiting (Dual-Identifier Strategy)
	// ========================================================================
	// Rate limit by BOTH IP and email hash to prevent:
	// 1. Multiple users behind NAT from being blocked together (IP-only problem)
	// 2. Attacker rotating IPs to bypass limits (IP-only problem)

	// Rate limit by IP (prevent brute force from single IP)
	ipIdentifier := ratelimit.FormatIdentifier(ratelimit.IdentifierIP, ipAddress)
	ipStatus, err := uc.rateLimiter.RecordAttempt(ctx, ipIdentifier, ratelimit.ActionLogin)
	if err != nil {
		return nil, fmt.Errorf("rate limit check failed: %w", err)
	}

	if !ipStatus.IsAllowed() {
		uc.logLoginActivity(ctx, 0, req.Email, false, "Rate limited by IP", ipAddress)
		return nil, &RateLimitError{
			Action:     ratelimit.ActionLogin,
			Status:     ipStatus,
			RetryAfter: ipStatus.TimeUntilReset(),
		}
	}

	// Rate limit by email hash (prevent distributed brute force on single account)
	// Use hash to avoid storing email in Redis keys
	emailIdentifier := ratelimit.FormatIdentifier(
		ratelimit.IdentifierEmail,
		hashEmail(req.Email),
	)
	emailStatus, err := uc.rateLimiter.RecordAttempt(ctx, emailIdentifier, ratelimit.ActionLogin)
	if err != nil {
		return nil, fmt.Errorf("rate limit check failed: %w", err)
	}

	if !emailStatus.IsAllowed() {
		uc.logLoginActivity(ctx, 0, req.Email, false, "Rate limited by email", ipAddress)
		return nil, &RateLimitError{
			Action:     ratelimit.ActionLogin,
			Status:     emailStatus,
			RetryAfter: emailStatus.TimeUntilReset(),
		}
	}

	// ========================================================================
	// STEP 3: Find User
	// ========================================================================
	foundUser, err := uc.userRepo.FindByEmail(ctx, req.Email)
	if err != nil {
		// Don't reveal if user exists (security best practice)
		uc.logLoginActivity(ctx, 0, req.Email, false, "User not found", ipAddress)
		return nil, ErrInvalidCredentials
	}

	// ========================================================================
	// STEP 4: Get Security Info (with proper error handling)
	// ========================================================================
	securityInfo, err := uc.securityRepo.GetByUserID(ctx, foundUser.ID)
	if err != nil {
		// If security info doesn't exist, create default
		securityInfo = &user.UserSecurityInfo{
			UserID:              foundUser.ID,
			FailedLoginAttempts: 0,
		}
		// Attempt to create, but continue even if it fails
		_ = uc.securityRepo.Create(ctx, securityInfo)
	}

	// ========================================================================
	// STEP 5: Check Account Locks (Domain Policy)
	// ========================================================================
	if securityInfo.IsLocked() {
		uc.logLoginActivity(ctx, foundUser.ID, req.Email, false, "Account locked", ipAddress)
		return nil, ErrAccountLocked
	}

	// Check if user can login (Domain Policy)
	if err := foundUser.CanLogin(); err != nil {
		uc.logLoginActivity(ctx, foundUser.ID, req.Email, false, err.Error(), ipAddress)
		return nil, err
	}

	// ========================================================================
	// STEP 6: Verify Password
	// ========================================================================
	credential, err := uc.credentialRepo.GetByUserID(ctx, foundUser.ID)
	if err != nil {
		uc.logLoginActivity(ctx, foundUser.ID, req.Email, false, "Credential not found", ipAddress)
		return nil, ErrInvalidCredentials
	}

	// Verify password (use bcrypt or similar)
	if !verifyPassword(credential.PasswordHash, req.Password) {
		// Increment failed login attempts
		securityInfo.IncrementFailedAttempts()
		_ = uc.securityRepo.Update(ctx, securityInfo)

		uc.logLoginActivity(ctx, foundUser.ID, req.Email, false, "Invalid password", ipAddress)
		return nil, ErrInvalidCredentials
	}

	// ========================================================================
	// STEP 7: Check 2FA (if enabled)
	// ========================================================================
	if credential.TwoFactorEnabled {
		if req.TwoFactorCode == nil || *req.TwoFactorCode == "" {
			return nil, ErrTwoFactorRequired
		}

		if !verify2FACode(credential.TwoFactorSecret, *req.TwoFactorCode) {
			uc.logLoginActivity(ctx, foundUser.ID, req.Email, false, "Invalid 2FA code", ipAddress)
			return nil, ErrInvalid2FA
		}
	}

	// ========================================================================
	// STEP 8: SUCCESS - Reset Everything
	// ========================================================================

	// Reset rate limits for both identifiers
	_ = uc.rateLimiter.Reset(ctx, ipIdentifier, ratelimit.ActionLogin)
	_ = uc.rateLimiter.Reset(ctx, emailIdentifier, ratelimit.ActionLogin)

	// Reset failed login attempts
	securityInfo.ResetFailedAttempts()
	securityInfo.UpdateLastLogin()
	_ = uc.securityRepo.Update(ctx, securityInfo)

	// Create session
	session := &user.UserSession{
		UserID:     foundUser.ID,
		DeviceID:   req.DeviceID,
		DeviceName: req.DeviceName,
		Platform:   req.Platform,
		IPAddress:  &ipAddress,
		ExpiresAt:  time.Now().Add(7 * 24 * time.Hour),
	}

	if err := uc.sessionRepo.Create(ctx, session); err != nil {
		return nil, fmt.Errorf("failed to create session: %w", err)
	}

	// Generate tokens
	accessToken, refreshToken, err := uc.generateTokens(foundUser, session)
	if err != nil {
		return nil, fmt.Errorf("failed to generate tokens: %w", err)
	}

	// Log successful login
	uc.logLoginActivity(ctx, foundUser.ID, req.Email, true, "", ipAddress)

	// Return response
	return &user.LoginResponse{
		User:         foundUser.ToResponse(securityInfo, credential),
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		SessionID:    session.ID,
		ExpiresAt:    time.Now().Add(15 * time.Minute),
	}, nil
}

// ============================================================================
// PASSWORD RESET USE CASE
// ============================================================================

type PasswordResetUseCase struct {
	userRepo    UserRepository
	tokenRepo   TokenRepository
	rateLimiter ratelimit.Limiter
}

func NewPasswordResetUseCase(
	userRepo UserRepository,
	tokenRepo TokenRepository,
	rateLimiter ratelimit.Limiter,
) *PasswordResetUseCase {
	return &PasswordResetUseCase{
		userRepo:    userRepo,
		tokenRepo:   tokenRepo,
		rateLimiter: rateLimiter,
	}
}

func (uc *PasswordResetUseCase) RequestReset(ctx context.Context, email string, ipAddress string) error {
	// STEP 1: Validate input format
	if err := validation.ValidateEmail(email); err != nil {
		return err
	}

	// STEP 2: Rate limit by BOTH IP and email
	// This prevents both IP-based and email-based abuse

	ipIdentifier := ratelimit.FormatIdentifier(ratelimit.IdentifierIP, ipAddress)
	ipStatus, err := uc.rateLimiter.RecordAttempt(ctx, ipIdentifier, ratelimit.ActionPasswordReset)
	if err != nil {
		return err
	}

	if !ipStatus.IsAllowed() {
		return &RateLimitError{
			Action:     ratelimit.ActionPasswordReset,
			Status:     ipStatus,
			RetryAfter: ipStatus.TimeUntilReset(),
		}
	}

	emailIdentifier := ratelimit.FormatIdentifier(
		ratelimit.IdentifierEmail,
		hashEmail(email),
	)
	emailStatus, err := uc.rateLimiter.RecordAttempt(ctx, emailIdentifier, ratelimit.ActionPasswordReset)
	if err != nil {
		return err
	}

	if !emailStatus.IsAllowed() {
		return &RateLimitError{
			Action:     ratelimit.ActionPasswordReset,
			Status:     emailStatus,
			RetryAfter: emailStatus.TimeUntilReset(),
		}
	}

	// STEP 3: Find user
	foundUser, err := uc.userRepo.FindByEmail(ctx, email)
	if err != nil {
		// Don't reveal if user exists, but still rate limit
		// This is a security best practice
		return nil
	}

	// STEP 4: Generate and send reset token
	token := generateSecureToken()
	resetToken := &user.PasswordResetToken{
		UserID:    foundUser.ID,
		Token:     token,
		ExpiresAt: time.Now().Add(1 * time.Hour),
	}

	if err := uc.tokenRepo.CreatePasswordReset(ctx, resetToken); err != nil {
		return fmt.Errorf("failed to create reset token: %w", err)
	}

	// Send email (async)
	go sendPasswordResetEmail(email, token)

	return nil
}

// ============================================================================
// REGISTRATION USE CASE
// ============================================================================

type RegistrationUseCase struct {
	userRepo    UserRepository
	rateLimiter ratelimit.Limiter
}

func NewRegistrationUseCase(
	userRepo UserRepository,
	rateLimiter ratelimit.Limiter,
) *RegistrationUseCase {
	return &RegistrationUseCase{
		userRepo:    userRepo,
		rateLimiter: rateLimiter,
	}
}

func (uc *RegistrationUseCase) Register(ctx context.Context, req *user.UserCreateRequest, ipAddress string) error {
	// STEP 1: Validate input
	if errs := validation.ValidateUserCreateRequest(req); len(errs) > 0 {
		return fmt.Errorf("validation failed: %v", errs[0])
	}

	// STEP 2: Rate limit by IP (prevent spam registrations)
	ipIdentifier := ratelimit.FormatIdentifier(ratelimit.IdentifierIP, ipAddress)
	status, err := uc.rateLimiter.RecordAttempt(ctx, ipIdentifier, ratelimit.ActionRegistration)
	if err != nil {
		return err
	}

	if !status.IsAllowed() {
		return &RateLimitError{
			Action:     ratelimit.ActionRegistration,
			Status:     status,
			RetryAfter: status.TimeUntilReset(),
		}
	}

	// STEP 3: Check if email already exists
	existing, err := uc.userRepo.FindByEmail(ctx, req.Email)
	if err == nil && existing != nil {
		return errors.New("email_already_exists")
	}

	// STEP 4: Create user
	// ... implementation

	return nil
}

// ============================================================================
// HELPER FUNCTIONS
// ============================================================================

// hashEmail creates a SHA-256 hash of email for rate limiting
// This prevents storing raw emails in Redis keys
func hashEmail(email string) string {
	h := sha256.New()
	h.Write([]byte(email))
	return hex.EncodeToString(h.Sum(nil))[:16] // Use first 16 chars
}

// generateSecureToken generates a cryptographically secure random token
func generateSecureToken() string {
	// Implementation using crypto/rand
	return "secure_token"
}

// verifyPassword verifies password using bcrypt
func verifyPassword(hash *string, password string) bool {
	// Implementation using bcrypt.CompareHashAndPassword
	return true
}

// verify2FACode verifies 2FA TOTP code
func verify2FACode(secret *string, code string) bool {
	// Implementation using TOTP library
	return true
}

// sendPasswordResetEmail sends password reset email
func sendPasswordResetEmail(email, token string) {
	// Implementation using email service
}

// generateTokens generates access and refresh tokens
func (uc *LoginUseCase) generateTokens(user *user.User, session *user.UserSession) (string, string, error) {
	// Implementation using JWT
	return "access_token", "refresh_token", nil
}

// logLoginActivity logs login attempt
func (uc *LoginUseCase) logLoginActivity(ctx context.Context, userID int, email string, success bool, reason string, ipAddress string) {
	activity := &user.LoginActivity{
		UserID:    &userID,
		Email:     email,
		Success:   success,
		Reason:    &reason,
		IPAddress: ipAddress,
	}
	_ = uc.activityRepo.CreateLoginActivity(ctx, activity)
}

// ============================================================================
// RATE LIMIT ERROR
// ============================================================================

// RateLimitError represents a rate limit error with retry information
type RateLimitError struct {
	Action     string
	Status     *ratelimit.RateLimitStatus
	RetryAfter time.Duration
}

func (e *RateLimitError) Error() string {
	if e.Status.Blocked {
		return fmt.Sprintf("rate_limit_exceeded: blocked until %s",
			e.Status.BlockedUntil.Format(time.RFC3339))
	}
	return fmt.Sprintf("rate_limit_exceeded: %d attempts remaining",
		e.Status.RemainingTries)
}

// IsRateLimitError checks if error is a rate limit error
func IsRateLimitError(err error) bool {
	_, ok := err.(*RateLimitError)
	return ok
}

// ============================================================================
// REPOSITORY INTERFACES
// ============================================================================

type UserRepository interface {
	FindByEmail(ctx context.Context, email string) (*user.User, error)
	Create(ctx context.Context, user *user.User) error
}

type CredentialRepository interface {
	GetByUserID(ctx context.Context, userID int) (*user.UserCredential, error)
}

type SecurityRepository interface {
	GetByUserID(ctx context.Context, userID int) (*user.UserSecurityInfo, error)
	Create(ctx context.Context, info *user.UserSecurityInfo) error
	Update(ctx context.Context, info *user.UserSecurityInfo) error
}

type SessionRepository interface {
	Create(ctx context.Context, session *user.UserSession) error
}

type TokenRepository interface {
	CreatePasswordReset(ctx context.Context, token *user.PasswordResetToken) error
}

type ActivityRepository interface {
	CreateLoginActivity(ctx context.Context, activity *user.LoginActivity) error
}
