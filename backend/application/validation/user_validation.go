package validation

import (
	"errors"
	"regexp"
	"strings"

	"github.com/LePhuocVuTien/SurvivalPro-Backend/internal/domain/user"
)

// ============================================================================
// INPUT VALIDATION
// ============================================================================
// This layer handles INPUT FORMAT validation ONLY:
// - Email format
// - Password format
// - Required fields
// - Length constraints
// - Regex patterns
//
// BUSINESS RULES (can login? can delete?) belong in domain/user/user_policy.go

var (
	// Format errors
	ErrInvalidEmail          = errors.New("invalid email format")
	ErrPasswordTooShort      = errors.New("password must be at least 8 characters")
	ErrPasswordNoUppercase   = errors.New("password must contain at least one uppercase letter")
	ErrPasswordNoLowercase   = errors.New("password must contain at least one lowercase letter")
	ErrPasswordNoNumber      = errors.New("password must contain at least one number")
	ErrPasswordNoSpecialChar = errors.New("password must contain at least one special character")
	ErrInvalidPhone          = errors.New("invalid phone format")
	ErrInvalidRole           = errors.New("invalid role")
	ErrInvalidAccountStatus  = errors.New("invalid account status")
	ErrInvalidProvider       = errors.New("invalid auth provider")
	ErrNameTooShort          = errors.New("name must be at least 2 characters")
	ErrNameTooLong           = errors.New("name must not exceed 100 characters")
	ErrReasonRequired        = errors.New("reason is required")
	ErrReasonTooShort        = errors.New("reason must be at least 10 characters")
	ErrTokenRequired         = errors.New("token is required")
	ErrInvalid2FACode        = errors.New("2FA code must be 6 digits")

	// Required field errors
	ErrEmailRequired    = errors.New("email is required")
	ErrPasswordRequired = errors.New("password is required")
	ErrNameRequired     = errors.New("name is required")
	ErrTokenEmpty       = errors.New("token cannot be empty")
	ErrUserIDRequired   = errors.New("user_id is required")
)

// ============================================================================
// FORMAT VALIDATION FUNCTIONS
// ============================================================================

// ValidateEmail validates email format
func ValidateEmail(email string) error {
	if email == "" {
		return ErrEmailRequired
	}

	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)
	if !emailRegex.MatchString(email) {
		return ErrInvalidEmail
	}

	return nil
}

// ValidatePassword validates password format (NOT business rules)
func ValidatePassword(password string) []error {
	var errs []error

	if password == "" {
		return []error{ErrPasswordRequired}
	}

	if len(password) < 8 {
		errs = append(errs, ErrPasswordTooShort)
	}

	if !regexp.MustCompile(`[A-Z]`).MatchString(password) {
		errs = append(errs, ErrPasswordNoUppercase)
	}

	if !regexp.MustCompile(`[a-z]`).MatchString(password) {
		errs = append(errs, ErrPasswordNoLowercase)
	}

	if !regexp.MustCompile(`[0-9]`).MatchString(password) {
		errs = append(errs, ErrPasswordNoNumber)
	}

	if !regexp.MustCompile(`[!@#$%^&*()_+\-=\[\]{}|;:,.<>?]`).MatchString(password) {
		errs = append(errs, ErrPasswordNoSpecialChar)
	}

	return errs
}

// ValidatePhone validates phone number format
func ValidatePhone(phone string) error {
	if phone == "" {
		return nil // Phone is optional
	}

	normalized := NormalizePhone(phone)

	// International format
	internationalPattern := `^\+[1-9]\d{7,14}$`
	if matched, _ := regexp.MatchString(internationalPattern, normalized); matched {
		return nil
	}

	// Local format
	localPattern := `^0\d{9,10}$`
	if matched, _ := regexp.MatchString(localPattern, normalized); matched {
		return nil
	}

	return ErrInvalidPhone
}

// ValidateRole validates role enum
func ValidateRole(role user.UserRole) error {
	switch role {
	case user.UserRoleAdmin, user.UserRoleLeader, user.UserRoleUser:
		return nil
	default:
		return ErrInvalidRole
	}
}

// ValidateAccountStatus validates account status enum
func ValidateAccountStatus(status user.AccountStatus) error {
	switch status {
	case user.AccountPending, user.AccountActive,
		user.AccountSuspended, user.AccountBanned, user.AccountClosed:
		return nil
	default:
		return ErrInvalidAccountStatus
	}
}

// ValidateAuthProvider validates auth provider enum
func ValidateAuthProvider(provider user.AuthProvider) error {
	switch provider {
	case user.ProviderGoogle, user.ProviderFacebook, user.ProviderApple:
		return nil
	default:
		return ErrInvalidProvider
	}
}

// ValidateName validates name format
func ValidateName(name string) error {
	if name == "" {
		return ErrNameRequired
	}
	if len(name) < 2 {
		return ErrNameTooShort
	}
	if len(name) > 100 {
		return ErrNameTooLong
	}
	return nil
}

// Validate2FACode validates 2FA code format
func Validate2FACode(code string) error {
	if len(code) != 6 {
		return ErrInvalid2FACode
	}
	if !regexp.MustCompile(`^\d{6}$`).MatchString(code) {
		return ErrInvalid2FACode
	}
	return nil
}

// ============================================================================
// DTO INPUT VALIDATION
// ============================================================================

// ValidateUserCreateRequest validates user creation INPUT
func ValidateUserCreateRequest(req *user.UserCreateRequest) []error {
	var errs []error

	if err := ValidateEmail(req.Email); err != nil {
		errs = append(errs, err)
	}

	errs = append(errs, ValidatePassword(req.Password)...)

	if err := ValidateName(req.Name); err != nil {
		errs = append(errs, err)
	}

	if req.Phone != nil && *req.Phone != "" {
		if err := ValidatePhone(*req.Phone); err != nil {
			errs = append(errs, err)
		}
	}

	if req.Role != "" {
		if err := ValidateRole(req.Role); err != nil {
			errs = append(errs, err)
		}
	}

	return errs
}

// ValidateUserUpdateRequest validates user update INPUT
func ValidateUserUpdateRequest(req *user.UserUpdateRequest) []error {
	var errs []error

	if req.Name != nil {
		if err := ValidateName(*req.Name); err != nil {
			errs = append(errs, err)
		}
	}

	if req.Phone != nil && *req.Phone != "" {
		if err := ValidatePhone(*req.Phone); err != nil {
			errs = append(errs, err)
		}
	}

	return errs
}

// ValidateLoginRequest validates login INPUT
func ValidateLoginRequest(req *user.LoginRequest) []error {
	var errs []error

	if err := ValidateEmail(req.Email); err != nil {
		errs = append(errs, err)
	}

	if req.Password == "" {
		errs = append(errs, ErrPasswordRequired)
	}

	if req.TwoFactorCode != nil && *req.TwoFactorCode != "" {
		if err := Validate2FACode(*req.TwoFactorCode); err != nil {
			errs = append(errs, err)
		}
	}

	return errs
}

// ValidateChangePasswordRequest validates change password INPUT
func ValidateChangePasswordRequest(req *user.ChangePasswordRequest) []error {
	var errs []error

	if req.OldPassword == "" {
		errs = append(errs, errors.New("old password is required"))
	}

	errs = append(errs, ValidatePassword(req.NewPassword)...)

	return errs
}

// ValidateForgotPasswordRequest validates forgot password INPUT
func ValidateForgotPasswordRequest(req *user.ForgotPasswordRequest) []error {
	var errs []error

	if err := ValidateEmail(req.Email); err != nil {
		errs = append(errs, err)
	}

	return errs
}

// ValidateResetPasswordRequest validates reset password INPUT
func ValidateResetPasswordRequest(req *user.ResetPasswordRequest) []error {
	var errs []error

	if req.Token == "" {
		errs = append(errs, ErrTokenRequired)
	}

	errs = append(errs, ValidatePassword(req.NewPassword)...)

	return errs
}

// ValidateUpdateAccountStatusRequest validates update account status INPUT
func ValidateUpdateAccountStatusRequest(req *user.UpdateAccountStatusRequest) []error {
	var errs []error

	if req.UserID <= 0 {
		errs = append(errs, ErrUserIDRequired)
	}

	if err := ValidateAccountStatus(req.Status); err != nil {
		errs = append(errs, err)
	}

	if req.Reason == "" {
		errs = append(errs, ErrReasonRequired)
	} else if len(req.Reason) < 10 {
		errs = append(errs, ErrReasonTooShort)
	}

	return errs
}

// ValidateSocialLoginRequest validates social login INPUT
func ValidateSocialLoginRequest(req *user.SocialLoginRequest) []error {
	var errs []error

	if err := ValidateAuthProvider(req.Provider); err != nil {
		errs = append(errs, err)
	}

	if req.AccessToken == "" {
		errs = append(errs, errors.New("access_token is required"))
	}

	if req.ProviderUserID == "" {
		errs = append(errs, errors.New("provider_user_id is required"))
	}

	return errs
}

// ============================================================================
// HELPER FUNCTIONS
// ============================================================================

// NormalizePhone normalizes phone number format
func NormalizePhone(phone string) string {
	phone = strings.ReplaceAll(phone, " ", "")
	phone = strings.ReplaceAll(phone, "-", "")
	phone = strings.ReplaceAll(phone, "(", "")
	phone = strings.ReplaceAll(phone, ")", "")
	return phone
}

// NormalizeEmail normalizes email format
func NormalizeEmail(email string) string {
	return strings.ToLower(strings.TrimSpace(email))
}

// FormatValidationErrors formats validation errors for API response
func FormatValidationErrors(errs []error) map[string]string {
	errors := make(map[string]string)
	for _, err := range errs {
		switch err {
		case ErrInvalidEmail, ErrEmailRequired:
			errors["email"] = err.Error()
		case ErrPasswordTooShort, ErrPasswordNoUppercase, ErrPasswordNoLowercase,
			ErrPasswordNoNumber, ErrPasswordNoSpecialChar, ErrPasswordRequired:
			if _, exists := errors["password"]; !exists {
				errors["password"] = err.Error()
			} else {
				errors["password"] += "; " + err.Error()
			}
		case ErrInvalidPhone:
			errors["phone"] = err.Error()
		case ErrInvalidRole:
			errors["role"] = err.Error()
		case ErrNameTooShort, ErrNameTooLong, ErrNameRequired:
			errors["name"] = err.Error()
		default:
			errors["general"] = err.Error()
		}
	}
	return errors
}
