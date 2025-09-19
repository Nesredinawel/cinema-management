package controllers

import (
	"auth-backend/models"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// ---------------- CreateUserByAdmin ----------------
// Admin-only route to create staff or admin users
func CreateUserByAdmin(c *gin.Context) {
	type Request struct {
		Name     string `json:"name" binding:"required"`
		Email    string `json:"email" binding:"required,email"`
		Phone    string `json:"phone" binding:"required"`
		Password string `json:"password" binding:"required,min=6"`
		Role     string `json:"role" binding:"required"` // "staff" or "admin"
		Dept     string `json:"dept"`                    // optional, for staff
		Level    string `json:"level"`                   // optional, for admin
	}

	var req Request
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validate role
	if req.Role != "staff" && req.Role != "admin" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Role must be 'staff' or 'admin'"})
		return
	}

	user := &models.User{
		Name:         req.Name,
		Email:        &req.Email,
		PhoneNumber:  &req.Phone,
		PasswordHash: req.Password, // plain password, will be hashed inside CreateOrFetchUser
		Role:         req.Role,
		IsVerified:   true, // admin-created users are verified by default
	}

	// Prepare extra map for role extension
	extra := make(map[string]interface{})
	if req.Role == "staff" {
		extra["dept"] = req.Dept
	} else if req.Role == "admin" {
		extra["level"] = req.Level
	}

	// Create user (or fetch if exists)
	createdUser, created, err := models.CreateOrFetchUser(user, extra)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}
	if !created {
		c.JSON(http.StatusConflict, gin.H{"error": "User with this email or phone already exists"})
		return
	}

	// Build role_extra for response
	roleExtra := gin.H{}
	switch createdUser.Role {
	case "admin":
		level, err := models.GetAdminLevel(createdUser.ID)
		if err != nil {
			level = "unknown"
		}
		roleExtra["level"] = level
	case "staff":
		dept, err := models.GetStaffDept(createdUser.ID)
		if err != nil {
			dept = "unknown"
		}
		roleExtra["dept"] = dept
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "User created successfully",
		"user": gin.H{
			"id":          createdUser.ID,
			"name":        createdUser.Name,
			"email":       createdUser.Email,
			"phone":       createdUser.PhoneNumber,
			"role":        createdUser.Role,
			"is_verified": createdUser.IsVerified,
			"role_extra":  roleExtra,
		},
	})
}

// ---------------- ChangeUserRole ----------------
// Allows admin to upgrade staff to admin
func ChangeUserRole(c *gin.Context) {
	type Request struct {
		UserID int    `json:"user_id" binding:"required"`
		Role   string `json:"role" binding:"required"` // only "admin" allowed
		Level  string `json:"level"`                   // optional, admin level
	}

	var req Request
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	req.Role = strings.ToLower(req.Role)
	if req.Role != "admin" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Role can only be upgraded to 'admin'"})
		return
	}

	user, err := models.GetUserByID(req.UserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch user"})
		return
	}
	if user == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	if user.Role == "admin" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User is already an admin"})
		return
	}

	if user.Role != "staff" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Only staff can be upgraded to admin"})
		return
	}

	// Update role and role_id
	user.Role = "admin"
	user.RoleID = 1 // admin role_id

	// Update role extension details
	extra := map[string]interface{}{"level": req.Level}
	if err := models.UpdateUser(user, extra); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user role"})
		return
	}

	// Build role_extra for response
	roleExtra := gin.H{}
	level, err := models.GetAdminLevel(user.ID)
	if err != nil {
		level = "unknown"
	}
	roleExtra["level"] = level

	c.JSON(http.StatusOK, gin.H{
		"message": "User role updated successfully",
		"user": gin.H{
			"id":          user.ID,
			"name":        user.Name,
			"role":        user.Role,
			"email":       user.Email,
			"phone":       user.PhoneNumber,
			"is_verified": user.IsVerified,
			"role_extra":  roleExtra,
		},
	})
}

// ---------------- ListUsers ----------------
// Returns all users with role extension details
func ListUsers(c *gin.Context) {
	users, err := models.GetAllUsers()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch users"})
		return
	}

	var response []gin.H
	for _, u := range users {
		roleExtra := gin.H{}
		switch u.Role {
		case "admin":
			level, err := models.GetAdminLevel(u.ID)
			if err != nil {
				level = "unknown"
			}
			roleExtra["level"] = level
		case "staff":
			dept, err := models.GetStaffDept(u.ID)
			if err != nil {
				dept = "unknown"
			}
			roleExtra["dept"] = dept
		case "customer":
			points, err := models.GetCustomerPoints(u.ID)
			if err != nil {
				points = 0
			}
			roleExtra["loyalty_points"] = points
		}

		response = append(response, gin.H{
			"id":          u.ID,
			"name":        u.Name,
			"email":       u.Email,
			"phone":       u.PhoneNumber,
			"role":        u.Role,
			"is_verified": u.IsVerified,
			"created_at":  u.CreatedAt,
			"updated_at":  u.UpdatedAt,
			"role_extra":  roleExtra,
		})
	}

	c.JSON(http.StatusOK, gin.H{"users": response})
}
