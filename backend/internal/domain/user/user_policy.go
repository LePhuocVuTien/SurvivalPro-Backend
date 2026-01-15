package user

import "errors"

// ============================================================================
// DOMAIN ERRORS
// ============================================================================

var (
	// Account status errors
	ErrUserNotActive = errors.New("user account is not active")
	ErrUserPending   = errors.New("user account is pending activation")
	ErrUserSuspended = errors.New("user account is suspended")
	ErrUserBanned    = errors.New("user account is banned")
	ErrUserClosed    = errors.New("user account is closed")
	ErrUserDeleted   = errors.New("user account has been deleted")

	// Permission errors
	ErrPermissionDenied       = errors.New("permission denied")
	ErrCannotModifySelf       = errors.New("cannot modify own account in this way")
	ErrCannotModifyAdmin      = errors.New("cannot modify admin account")
	ErrCannotPromoteToAdmin   = errors.New("only admins can promote users to admin role")
	ErrCannotBanAdmin         = errors.New("cannot ban admin account")
	ErrCannotCloseAdmin       = errors.New("cannot close admin account")
	ErrInsufficientPermission = errors.New("insufficient permission for this action")

	// Status transition errors
	ErrInvalidStatusTransition = errors.New("invalid account status transition")
	ErrCannotReactivateClosed  = errors.New("cannot reactivate closed account")
	ErrCannotModifyBanned      = errors.New("cannot modify banned account")
	ErrAccountLocked           = errors.New("account is locked")
	ErrUserInactive            = errors.New("user account is inactive")
	ErrEmailNotVerified        = errors.New("email is not verified")
)

// ============================================================================
// LOGIN POLICY
// ============================================================================

// CanLogin checks if user can login to the system
// Business rule: User must be active and not deleted
func CanLogin(user *User) error {
	if user == nil {
		return errors.New("user is nil")
	}

	// Check soft delete
	if user.IsDeleted() {
		return ErrUserDeleted
	}

	// Check account status
	switch user.AccountStatus {
	case AccountPending:
		return ErrUserPending
	case AccountSuspended:
		return ErrUserSuspended
	case AccountBanned:
		return ErrUserBanned
	case AccountClosed:
		return ErrUserClosed
	case AccountActive:
		return nil
	default:
		return errors.New("unknown account status")
	}
}

// CanPerformAction checks if user can perform actions in the system
// Business rule: User must be active, not deleted, and email verified
func CanPerformAction(user *User) error {
	if err := CanLogin(user); err != nil {
		return err
	}

	// Additional check: email verification required for actions
	if !user.EmailVerified {
		return errors.New("email verification required")
	}

	return nil
}

// ============================================================================
// PROFILE MODIFICATION POLICY
// ============================================================================

// CanChangeProfile checks if actor can change target's profile
// Business rule:
// - Admin: can change anyone's profile
// - Leader: can change user's profile (not admin)
// - User: can only change own profile
func CanChangeProfile(actor *User, target *User) error {
	if actor == nil || target == nil {
		return errors.New("actor or target is nil")
	}

	// Check if actor can perform actions
	if err := CanPerformAction(actor); err != nil {
		return err
	}

	// Check if target is deleted
	if target.IsDeleted() {
		return ErrUserDeleted
	}

	// Admin can change anyone's profile
	if actor.IsAdmin() {
		return nil
	}

	// Leader can change user's profile (but not admin)
	if actor.IsLeader() {
		if target.IsAdmin() {
			return ErrCannotModifyAdmin
		}
		return nil
	}

	// Regular user can only change own profile
	if actor.ID == target.ID {
		return nil
	}

	return ErrPermissionDenied
}

// ============================================================================
// ROLE CHANGE POLICY
// ============================================================================

// CanChangeRole checks if actor can change target's role to newRole
// Business rule:
// - Only admin can promote to admin role
// - Admin can change anyone's role
// - Leader can change user's role to user/leader (not to admin)
// - User cannot change roles
// - Cannot change own role (safety measure)
func CanChangeRole(actor *User, target *User, newRole UserRole) error {
	if actor == nil || target == nil {
		return errors.New("actor or target is nil")
	}

	// Check if actor can perform actions
	if err := CanPerformAction(actor); err != nil {
		return err
	}

	// Check if target is deleted
	if target.IsDeleted() {
		return ErrUserDeleted
	}

	// Cannot change own role (safety measure)
	if actor.ID == target.ID {
		return ErrCannotModifySelf
	}

	// Validate new role
	if !isValidRole(newRole) {
		return errors.New("invalid role")
	}

	// Rule 1: Only admin can promote to admin
	if newRole == UserRoleAdmin {
		if !actor.IsAdmin() {
			return ErrCannotPromoteToAdmin
		}
	}

	// Rule 2: Admin has full control
	if actor.IsAdmin() {
		return nil
	}

	// Rule 3: Leader can change user roles (but not to admin, not admin users)
	if actor.IsLeader() {
		// Cannot modify admin
		if target.IsAdmin() {
			return ErrCannotModifyAdmin
		}

		// Cannot promote to admin
		if newRole == UserRoleAdmin {
			return ErrCannotPromoteToAdmin
		}

		// Can change user/leader roles
		return nil
	}

	// Rule 4: Regular user cannot change roles
	return ErrInsufficientPermission
}

// ============================================================================
// ACCOUNT STATUS CHANGE POLICY
// ============================================================================

// ValidStatusTransitions defines allowed status transitions
// Business rule: Only specific transitions are allowed
var ValidStatusTransitions = map[AccountStatus][]AccountStatus{
	AccountPending: {
		AccountActive, // Admin/Leader activates after review
		AccountClosed, // Admin/Leader closes before activation
	},
	AccountActive: {
		AccountSuspended, // Admin/Leader temporarily suspends
		AccountBanned,    // Admin permanently bans
		AccountClosed,    // Admin closes account
	},
	AccountSuspended: {
		AccountActive, // Admin/Leader reactivates
		AccountBanned, // Admin escalates to ban
		AccountClosed, // Admin closes
	},
	AccountBanned: {
		AccountClosed, // Admin archives banned account
		// Note: Cannot reactivate banned account directly
	},
	AccountClosed: {
		// Terminal state - no transitions allowed
	},
}

// CanChangeAccountStatus checks if actor can change target's account status to newStatus
// Business rule:
// - Admin: can change any status (except closed → *)
// - Leader: can activate, suspend, reactivate (but not ban/close admin)
// - User: cannot change status
// - Cannot change own status (safety measure, even for admin)
// - Must follow valid status transitions
func CanChangeAccountStatus(actor *User, target *User, newStatus AccountStatus) error {
	if actor == nil || target == nil {
		return errors.New("actor or target is nil")
	}

	// Check if actor can perform actions
	if err := CanPerformAction(actor); err != nil {
		return err
	}

	// Check if target is deleted
	if target.IsDeleted() {
		return ErrUserDeleted
	}

	// Rule 1: Cannot change own status (safety measure)
	if actor.ID == target.ID {
		return ErrCannotModifySelf
	}

	// Rule 2: Validate status transition
	if !isValidStatusTransition(target.AccountStatus, newStatus) {
		return ErrInvalidStatusTransition
	}

	// Rule 3: Admin has full control (except reactivating closed)
	if actor.IsAdmin() {
		// Even admin cannot reactivate closed accounts
		if target.AccountStatus == AccountClosed && newStatus != AccountClosed {
			return ErrCannotReactivateClosed
		}
		return nil
	}

	// Rule 4: Leader permissions
	if actor.IsLeader() {
		// Leader cannot ban or close admin
		if target.IsAdmin() {
			if newStatus == AccountBanned || newStatus == AccountClosed {
				return ErrCannotBanAdmin
			}
		}

		// Leader can: pending→active, active→suspended, suspended→active
		switch target.AccountStatus {
		case AccountPending:
			if newStatus == AccountActive || newStatus == AccountClosed {
				return nil
			}
		case AccountActive:
			if newStatus == AccountSuspended {
				return nil
			}
		case AccountSuspended:
			if newStatus == AccountActive {
				return nil
			}
		}

		// Leader cannot perform other status changes
		return ErrInsufficientPermission
	}

	// Rule 5: Regular user cannot change status
	return ErrInsufficientPermission
}

// ============================================================================
// DELETE POLICY
// ============================================================================

// CanDeleteUser checks if actor can delete (soft delete) target
// Business rule:
// - Admin: can delete anyone except self
// - Leader: can delete user (not admin)
// - User: cannot delete
func CanDeleteUser(actor *User, target *User) error {
	if actor == nil || target == nil {
		return errors.New("actor or target is nil")
	}

	// Check if actor can perform actions
	if err := CanPerformAction(actor); err != nil {
		return err
	}

	// Check if target is already deleted
	if target.IsDeleted() {
		return ErrUserDeleted
	}

	// Cannot delete self (safety measure)
	if actor.ID == target.ID {
		return ErrCannotModifySelf
	}

	// Admin can delete anyone
	if actor.IsAdmin() {
		return nil
	}

	// Leader can delete user (but not admin)
	if actor.IsLeader() {
		if target.IsAdmin() {
			return ErrCannotModifyAdmin
		}
		return nil
	}

	// Regular user cannot delete
	return ErrPermissionDenied
}

// ============================================================================
// VIEW POLICY
// ============================================================================

// CanViewProfile checks if actor can view target's profile
// Business rule:
// - Admin: can view all profiles
// - Leader: can view all profiles
// - User: can view own profile and other users' basic profiles
func CanViewProfile(actor *User, target *User) error {
	if actor == nil || target == nil {
		return errors.New("actor or target is nil")
	}

	// Check if actor can perform actions
	if err := CanPerformAction(actor); err != nil {
		return err
	}

	// Admin and Leader can view all profiles
	if actor.IsAdmin() || actor.IsLeader() {
		return nil
	}

	// User can view own profile
	if actor.ID == target.ID {
		return nil
	}

	// User can view other users' basic profiles
	// (Sensitive data filtering should be done at application layer)
	return nil
}

// CanViewSensitiveData checks if actor can view target's sensitive data
// Business rule: Only admin and leader can view sensitive data
func CanViewSensitiveData(actor *User, target *User) error {
	if actor == nil || target == nil {
		return errors.New("actor or target is nil")
	}

	// Check if actor can perform actions
	if err := CanPerformAction(actor); err != nil {
		return err
	}

	// Admin and Leader can view sensitive data
	if actor.IsAdmin() || actor.IsLeader() {
		return nil
	}

	// User can view own sensitive data
	if actor.ID == target.ID {
		return nil
	}

	return ErrPermissionDenied
}

// ============================================================================
// SURVIVAL KIT SPECIFIC POLICIES
// ============================================================================

// CanCreateKit checks if user can create a survival kit
// Business rule: User must be active and verified
func CanCreateKit(user *User) error {
	return CanPerformAction(user)
}

// CanShareKit checks if user can share a kit with others
// Business rule: User must be active and verified
func CanShareKit(user *User) error {
	return CanPerformAction(user)
}

// CanAccessSharedKit checks if user can access a kit shared by another user
// Business rule: User must be active (verification not required for viewing)
func CanAccessSharedKit(user *User) error {
	return CanLogin(user)
}

// CanManageTeam checks if user can manage team members
// Business rule: Only admin and leader can manage teams
func CanManageTeam(user *User) error {
	if user == nil {
		return errors.New("user is nil")
	}

	if err := CanPerformAction(user); err != nil {
		return err
	}

	if user.IsAdmin() || user.IsLeader() {
		return nil
	}

	return ErrPermissionDenied
}

// CanInviteMembers checks if user can invite members to the platform
// Business rule: Active users can invite (viral growth)
func CanInviteMembers(user *User) error {
	return CanPerformAction(user)
}

// ============================================================================
// HELPER FUNCTIONS
// ============================================================================

// isValidRole checks if role is valid
func isValidRole(role UserRole) bool {
	switch role {
	case UserRoleAdmin, UserRoleLeader, UserRoleUser:
		return true
	default:
		return false
	}
}

// GetAllowedStatusTransitions returns allowed status transitions for current status
func GetAllowedStatusTransitions(currentStatus AccountStatus) []AccountStatus {
	transitions, exists := ValidStatusTransitions[currentStatus]
	if !exists {
		return []AccountStatus{}
	}
	return transitions
}

// ============================================================================
// PERMISSION METHODS (Add to existing User model)
// ============================================================================

// Role hierarchy map (higher number = more permissions)
var roleHierarchy = map[UserRole]int{
	UserRoleAdmin:  3, // Highest
	UserRoleLeader: 2, // Middle
	UserRoleUser:   1, // Lowest
}

// HasPermission checks if user has permission based on role hierarchy
// Example:
//   - Admin has permission for Leader and User roles
//   - Leader has permission for User role only
//   - User has permission for User role only
//
// Usage:
//
//	user.HasPermission(UserRoleLeader) // true if admin, false if user
func (u *User) HasPermission(requiredRole UserRole) bool {
	userLevel, userExists := roleHierarchy[u.Role]
	requiredLevel, requiredExists := roleHierarchy[requiredRole]

	if !userExists || !requiredExists {
		return false
	}

	// User has permission if their level >= required level
	return userLevel >= requiredLevel
}

// CanLogin checks if user can login (combines multiple checks)
func (u *User) CanLogin() error {
	if !u.IsActive() {
		return ErrUserInactive
	}
	if !u.EmailVerified {
		return ErrEmailNotVerified
	}
	return nil
}

// CanAccessResource checks if user can access a resource owned by targetUserID
func (u *User) CanAccessResource(targetUserID int) bool {
	// Admin can access everything
	if u.IsAdmin() {
		return true
	}

	// Users can access their own resources
	return u.ID == targetUserID
}

// CanModifyResource checks if user can modify a resource owned by targetUserID
func (u *User) CanModifyResource(targetUserID int) bool {
	// Admin can modify everything
	if u.IsAdmin() {
		return true
	}

	// Users can only modify their own resources
	return u.ID == targetUserID
}

// CanManageUsers checks if user can manage other users
func (u *User) CanManageUsers() bool {
	// Only admin and leader can manage users
	return u.IsAdmin() || u.IsLeader()
}

// CanChangeRole checks if user can change roles
func (u *User) CanChangeRole() bool {
	// Only admin can change roles
	return u.IsAdmin()
}

// CanChangeStatus checks if user can change account status
func (u *User) CanChangeStatus() bool {
	// Only admin can change status
	return u.IsAdmin()
}

// CanDeleteUser checks if user can delete another user
func (u *User) CanDeleteUser(targetUserID int) bool {
	// Admin can delete anyone except themselves
	if u.IsAdmin() && u.ID != targetUserID {
		return true
	}

	// Users can delete themselves
	return u.ID == targetUserID
}

// CanViewUser checks if user can view another user's details
func (u *User) CanViewUser(targetUserID int) bool {
	// Admin and leader can view anyone
	if u.IsAdmin() || u.IsLeader() {
		return true
	}

	// Users can view themselves
	return u.ID == targetUserID
}

// ============================================================================
// ROLE VALIDATION
// ============================================================================

// IsValid checks if role is valid
func (r UserRole) IsValid() bool {
	switch r {
	case UserRoleAdmin, UserRoleLeader, UserRoleUser:
		return true
	default:
		return false
	}
}

// String returns string representation of role
func (r UserRole) String() string {
	return string(r)
}

// ============================================================================
// ACCOUNT STATUS VALIDATION
// ============================================================================

// IsValid checks if status is valid
func (s AccountStatus) IsValid() bool {
	switch s {
	case AccountPending, AccountActive, AccountSuspended, AccountBanned, AccountClosed:
		return true
	default:
		return false
	}
}

// String returns string representation of status
func (s AccountStatus) String() string {
	return string(s)
}

// ============================================================================
// STATUS TRANSITION POLICY
// ============================================================================

// isValidStatusTransition checks if status transition is allowed
func isValidStatusTransition(from, to AccountStatus) bool {
	// Define valid transitions
	validTransitions := map[AccountStatus][]AccountStatus{
		AccountPending: {
			AccountActive, // Verify email
			AccountClosed, // User closes before verifying
		},
		AccountActive: {
			AccountSuspended, // Admin suspends
			AccountBanned,    // Admin bans
			AccountClosed,    // User closes
		},
		AccountSuspended: {
			AccountActive, // Admin unsuspends
			AccountBanned, // Admin bans
			AccountClosed, // Admin closes
		},
		AccountBanned: {
			// Banned is terminal, only admin can unban (rare)
			AccountActive, // Admin unbans (rare case)
		},
		AccountClosed: {
			// Closed is terminal, cannot reopen
		},
	}

	allowedStatuses, exists := validTransitions[from]
	if !exists {
		return false
	}

	for _, allowed := range allowedStatuses {
		if allowed == to {
			return true
		}
	}

	return false
}

// ============================================================================
// PERMISSION EXAMPLES
// ============================================================================

// Example usage:
//
// 1. Check if user has permission for a role level:
//    if user.HasPermission(UserRoleLeader) {
//        // User is leader or admin
//    }
//
// 2. Check if user can login:
//    if user.CanLogin() {
//        // Generate JWT token
//    }
//
// 3. Check resource access:
//    if user.CanAccessResource(resourceOwnerID) {
//        // Allow access
//    }
//
// 4. Check management permissions:
//    if user.CanManageUsers() {
//        // Show user management UI
//    }
//
// 5. Check status transitions:
//    if user.CanTransitionTo(AccountSuspended) {
//        // Allow suspension
//    }
