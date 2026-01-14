package handlers

import (
	"database/sql"
	"strconv"
	"strings"

	"github.com/LePhuocVuTien/SurvivalPro-Backend/internal/db"
	"github.com/LePhuocVuTien/SurvivalPro-Backend/internal/domain/user"
	"github.com/LePhuocVuTien/SurvivalPro-Backend/internal/middleware"
	"github.com/gofiber/fiber/v2"
)

// ============================================================================
// USER HANDLERS
// ============================================================================

// HandleListUsers returns list of users (Admin/Leader only)
func HandleListUsers(c *fiber.Ctx) error {
	currentUser := middleware.GetUserFromContext(c)

	// Parse pagination
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "20"))
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}

	// Parse filters
	role := c.Query("role")     // Filter by role
	status := c.Query("status") // Filter by status
	search := c.Query("search") // Search by name/email
	includeDeleted := c.Query("include_deleted") == "true"

	// Only admin can see deleted users
	if includeDeleted && !currentUser.IsAdmin() {
		includeDeleted = false
	}

	// Build query
	query := `
		SELECT 
			id, email, name, role, account_status,
			phone, avatar_url, email_verified, phone_verified,
			created_at, updated_at, deleted_at
		FROM users
		WHERE 1=1
	`
	args := []interface{}{}
	argCount := 1

	// Filter by role
	if role != "" {
		query += ` AND role = $` + strconv.Itoa(argCount)
		args = append(args, role)
		argCount++
	}

	// Filter by status
	if status != "" {
		query += ` AND account_status = $` + strconv.Itoa(argCount)
		args = append(args, status)
		argCount++
	}

	// Search by name or email
	if search != "" {
		query += ` AND (name ILIKE $` + strconv.Itoa(argCount) +
			` OR email ILIKE $` + strconv.Itoa(argCount) + `)`
		args = append(args, "%"+search+"%")
		argCount++
	}

	// Include/exclude deleted
	if !includeDeleted {
		query += ` AND deleted_at IS NULL`
	}

	// Count total
	countQuery := `SELECT COUNT(*) FROM (` + query + `) AS total`
	var total int
	err := db.DB.QueryRow(countQuery, args...).Scan(&total)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "Failed to count users",
		})
	}

	// Add pagination
	offset := (page - 1) * limit
	query += ` ORDER BY created_at DESC LIMIT $` + strconv.Itoa(argCount) +
		` OFFSET $` + strconv.Itoa(argCount+1)
	args = append(args, limit, offset)

	// Execute query
	rows, err := db.DB.Query(query, args...)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "Failed to fetch users",
		})
	}
	defer rows.Close()

	// Parse results
	users := []user.User{}
	for rows.Next() {
		var u user.User
		err := rows.Scan(
			&u.ID, &u.Email, &u.Name, &u.Role, &u.AccountStatus,
			&u.Phone, &u.AvatarURL, &u.EmailVerified, &u.PhoneVerified,
			&u.CreatedAt, &u.UpdatedAt, &u.DeletedAt,
		)
		if err != nil {
			continue
		}
		users = append(users, u)
	}

	// Calculate pagination
	totalPages := (total + limit - 1) / limit

	return c.JSON(fiber.Map{
		"users": users,
		"pagination": fiber.Map{
			"page":        page,
			"limit":       limit,
			"total":       total,
			"total_pages": totalPages,
		},
	})
}

// HandleGetUser returns a single user by ID
func HandleGetUser(c *fiber.Ctx) error {
	currentUser := middleware.GetUserFromContext(c)
	userID, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "Invalid user ID",
		})
	}

	// Check permission
	if !currentUser.CanViewUser(userID) {
		return c.Status(403).JSON(fiber.Map{
			"error": "Forbidden - you can only view your own profile or you need higher permissions",
		})
	}

	// Query user
	var u user.User
	query := `
		SELECT 
			id, email, name, role, account_status, status_changed_at,
			phone, avatar_url, blood_type, allergies,
			emergency_contact_name, emergency_contact_phone,
			email_verified, phone_verified,
			created_at, updated_at, deleted_at
		FROM users
		WHERE id = $1
	`

	// Admin can see deleted users
	if !currentUser.IsAdmin() {
		query += ` AND deleted_at IS NULL`
	}

	err = db.DB.QueryRow(query, userID).Scan(
		&u.ID, &u.Email, &u.Name, &u.Role, &u.AccountStatus, &u.StatusChangedAt,
		&u.Phone, &u.AvatarURL, &u.BloodType, &u.Allergies,
		&u.EmergencyContactName, &u.EmergencyContactPhone,
		&u.EmailVerified, &u.PhoneVerified,
		&u.CreatedAt, &u.UpdatedAt, &u.DeletedAt,
	)

	if err == sql.ErrNoRows {
		return c.Status(404).JSON(fiber.Map{
			"error": "User not found",
		})
	}
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "Failed to fetch user",
		})
	}

	return c.JSON(fiber.Map{
		"user": u,
	})
}

// HandleUpdateUser updates user information
func HandleUpdateUser(c *fiber.Ctx) error {
	currentUser := middleware.GetUserFromContext(c)
	userID, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "Invalid user ID",
		})
	}

	// Check permission
	if !currentUser.CanModifyResource(userID) {
		return c.Status(403).JSON(fiber.Map{
			"error": "Forbidden - you can only update your own profile",
		})
	}

	// Parse request
	type UpdateUserRequest struct {
		Name                  *string `json:"name"`
		Phone                 *string `json:"phone"`
		AvatarURL             *string `json:"avatar_url"`
		BloodType             *string `json:"blood_type"`
		Allergies             *string `json:"allergies"`
		EmergencyContactName  *string `json:"emergency_contact_name"`
		EmergencyContactPhone *string `json:"emergency_contact_phone"`
	}

	var req UpdateUserRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Validate
	if req.Name != nil && strings.TrimSpace(*req.Name) == "" {
		return c.Status(400).JSON(fiber.Map{
			"error": "Name cannot be empty",
		})
	}

	// Build update query
	updates := []string{}
	args := []interface{}{}
	argCount := 1

	if req.Name != nil {
		updates = append(updates, "name = $"+strconv.Itoa(argCount))
		args = append(args, *req.Name)
		argCount++
	}
	if req.Phone != nil {
		updates = append(updates, "phone = $"+strconv.Itoa(argCount))
		args = append(args, *req.Phone)
		argCount++
	}
	if req.AvatarURL != nil {
		updates = append(updates, "avatar_url = $"+strconv.Itoa(argCount))
		args = append(args, *req.AvatarURL)
		argCount++
	}
	if req.BloodType != nil {
		updates = append(updates, "blood_type = $"+strconv.Itoa(argCount))
		args = append(args, *req.BloodType)
		argCount++
	}
	if req.Allergies != nil {
		updates = append(updates, "allergies = $"+strconv.Itoa(argCount))
		args = append(args, *req.Allergies)
		argCount++
	}
	if req.EmergencyContactName != nil {
		updates = append(updates, "emergency_contact_name = $"+strconv.Itoa(argCount))
		args = append(args, *req.EmergencyContactName)
		argCount++
	}
	if req.EmergencyContactPhone != nil {
		updates = append(updates, "emergency_contact_phone = $"+strconv.Itoa(argCount))
		args = append(args, *req.EmergencyContactPhone)
		argCount++
	}

	if len(updates) == 0 {
		return c.Status(400).JSON(fiber.Map{
			"error": "No fields to update",
		})
	}

	// Add updated_at and updated_by
	updates = append(updates, "updated_at = NOW()")
	updates = append(updates, "updated_by = $"+strconv.Itoa(argCount))
	args = append(args, currentUser.ID)
	argCount++

	// Add WHERE clause
	args = append(args, userID)

	// Execute update
	query := `UPDATE users SET ` + strings.Join(updates, ", ") +
		` WHERE id = $` + strconv.Itoa(argCount) + ` AND deleted_at IS NULL`

	result, err := db.DB.Exec(query, args...)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "Failed to update user",
		})
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return c.Status(404).JSON(fiber.Map{
			"error": "User not found",
		})
	}

	return c.JSON(fiber.Map{
		"message": "User updated successfully",
	})
}

// HandleDeleteUser soft deletes a user
func HandleDeleteUser(c *fiber.Ctx) error {
	currentUser := middleware.GetUserFromContext(c)
	userID, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "Invalid user ID",
		})
	}

	// Check permission
	if !currentUser.CanDeleteUser(userID) {
		return c.Status(403).JSON(fiber.Map{
			"error": "Forbidden - you can only delete your own account or you need admin permissions",
		})
	}

	// Soft delete
	query := `
		UPDATE users 
		SET deleted_at = NOW(), deleted_by = $1
		WHERE id = $2 AND deleted_at IS NULL
	`

	result, err := db.DB.Exec(query, currentUser.ID, userID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "Failed to delete user",
		})
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return c.Status(404).JSON(fiber.Map{
			"error": "User not found or already deleted",
		})
	}

	return c.JSON(fiber.Map{
		"message": "User deleted successfully",
	})
}

// ============================================================================
// ADMIN HANDLERS
// ============================================================================

// HandleChangeUserRole changes user role (Admin only)
func HandleChangeUserRole(c *fiber.Ctx) error {
	currentUser := middleware.GetUserFromContext(c)
	userID, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "Invalid user ID",
		})
	}

	// Only admin can change roles
	if !currentUser.IsAdmin() {
		return c.Status(403).JSON(fiber.Map{
			"error": "Forbidden - only admins can change roles",
		})
	}

	// Parse request
	type ChangeRoleRequest struct {
		Role user.UserRole `json:"role"`
	}

	var req ChangeRoleRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Validate role
	if req.Role != user.UserRoleAdmin &&
		req.Role != user.UserRoleLeader &&
		req.Role != user.UserRoleUser {
		return c.Status(400).JSON(fiber.Map{
			"error": "Invalid role. Must be: admin, leader, or user",
		})
	}

	// Cannot change own role
	if userID == currentUser.ID {
		return c.Status(400).JSON(fiber.Map{
			"error": "Cannot change your own role",
		})
	}

	// Update role
	query := `
		UPDATE users 
		SET role = $1, updated_at = NOW(), updated_by = $2
		WHERE id = $3 AND deleted_at IS NULL
	`

	result, err := db.DB.Exec(query, req.Role, currentUser.ID, userID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "Failed to change role",
		})
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return c.Status(404).JSON(fiber.Map{
			"error": "User not found",
		})
	}

	return c.JSON(fiber.Map{
		"message":  "Role changed successfully",
		"new_role": req.Role,
	})
}

// HandleChangeAccountStatus changes account status (Admin only)
func HandleChangeAccountStatus(c *fiber.Ctx) error {
	currentUser := middleware.GetUserFromContext(c)
	userID, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "Invalid user ID",
		})
	}

	// Only admin can change status
	if !currentUser.IsAdmin() {
		return c.Status(403).JSON(fiber.Map{
			"error": "Forbidden - only admins can change account status",
		})
	}

	// Parse request
	type ChangeStatusRequest struct {
		Status user.AccountStatus `json:"status"`
		Reason *string            `json:"reason"`
	}

	var req ChangeStatusRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Validate status
	validStatuses := []user.AccountStatus{
		user.AccountActive,
		user.AccountSuspended,
		user.AccountBanned,
		user.AccountClosed,
	}
	isValid := false
	for _, s := range validStatuses {
		if req.Status == s {
			isValid = true
			break
		}
	}
	if !isValid {
		return c.Status(400).JSON(fiber.Map{
			"error": "Invalid status",
		})
	}

	// Cannot change own status
	if userID == currentUser.ID {
		return c.Status(400).JSON(fiber.Map{
			"error": "Cannot change your own account status",
		})
	}

	// Get current status
	var currentStatus user.AccountStatus
	err = db.DB.QueryRow(`SELECT account_status FROM users WHERE id = $1`, userID).
		Scan(&currentStatus)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{
			"error": "User not found",
		})
	}

	// Update status
	query := `
		UPDATE users 
		SET account_status = $1, 
			status_changed_at = NOW(),
			status_changed_by = $2,
			status_reason = $3,
			updated_at = NOW(),
			updated_by = $2
		WHERE id = $4 AND deleted_at IS NULL
	`

	result, err := db.DB.Exec(query, req.Status, currentUser.ID, req.Reason, userID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "Failed to change status",
		})
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return c.Status(404).JSON(fiber.Map{
			"error": "User not found",
		})
	}

	// Log status change
	logQuery := `
		INSERT INTO account_status_changes 
		(user_id, old_status, new_status, reason, changed_by, created_at)
		VALUES ($1, $2, $3, $4, $5, NOW())
	`
	db.DB.Exec(logQuery, userID, currentStatus, req.Status, req.Reason, currentUser.ID)

	return c.JSON(fiber.Map{
		"message":    "Account status changed successfully",
		"new_status": req.Status,
	})
}

// HandleRestoreUser restores a soft-deleted user (Admin only)
func HandleRestoreUser(c *fiber.Ctx) error {
	currentUser := middleware.GetUserFromContext(c)
	userID, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "Invalid user ID",
		})
	}

	// Only admin can restore
	if !currentUser.IsAdmin() {
		return c.Status(403).JSON(fiber.Map{
			"error": "Forbidden - only admins can restore users",
		})
	}

	// Restore user
	query := `
		UPDATE users 
		SET deleted_at = NULL, deleted_by = NULL, updated_at = NOW(), updated_by = $1
		WHERE id = $2 AND deleted_at IS NOT NULL
	`

	result, err := db.DB.Exec(query, currentUser.ID, userID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "Failed to restore user",
		})
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return c.Status(404).JSON(fiber.Map{
			"error": "User not found or not deleted",
		})
	}

	return c.JSON(fiber.Map{
		"message": "User restored successfully",
	})
}
