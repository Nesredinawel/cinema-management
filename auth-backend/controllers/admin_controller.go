package controllers

import (
	"auth-backend/models"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// ---------------- CreateUserByAdmin ----------------
// Admin-only route to create staff/admin/customer
func CreateUserByAdmin(c *gin.Context) {
	var req struct {
		Name         string                 `json:"name" binding:"required"`
		Email        *string                `json:"email"`
		PhoneNumber  *string                `json:"phone_number"`
		Password     string                 `json:"password" binding:"required"`
		Role         string                 `json:"role" binding:"required"` // admin/staff/customer
		ExtraDetails map[string]interface{} `json:"extra,omitempty"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	req.Role = strings.ToLower(req.Role)
	if req.Role != "admin" && req.Role != "staff" && req.Role != "customer" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Role must be 'admin', 'staff', or 'customer'"})
		return
	}

	user := &models.User{
		Name:         req.Name,
		Email:        req.Email,
		PhoneNumber:  req.PhoneNumber,
		PasswordHash: req.Password, // hashed inside CreateOrFetchUser
		Role:         req.Role,
		IsVerified:   req.Role != "customer", // customers require verification
	}

	createdUser, created, err := models.CreateOrFetchUser(user, req.ExtraDetails)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}
	if !created {
		c.JSON(http.StatusConflict, gin.H{"error": "User already exists"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "User created successfully",
		"user":    createdUser,
	})
}

// ---------------- ChangeUserRole ----------------
// Admin can upgrade staff → admin
func ChangeUserRole(c *gin.Context) {
	var req struct {
		UserID int    `json:"user_id" binding:"required"`
		Role   string `json:"role" binding:"required"` // only 'admin'
		Level  string `json:"level"`                   // optional, admin level
	}

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

	if user.Role != "staff" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Only staff can be upgraded to admin"})
		return
	}

	user.Role = "admin"
	if err := models.UpdateUser(user, map[string]interface{}{"level": req.Level}); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user role"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "User role updated successfully",
		"user":    user,
	})
}

// ---------------- ListUsers ----------------
// Returns all users with role-specific extra info
func ListUsers(c *gin.Context) {
	users, err := models.GetAllUsers()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch users"})
		return
	}

	var response []gin.H
	for _, u := range users {
		roleExtra := map[string]interface{}{}

		switch u.Role {
		case "admin":
			level, _ := models.GetAdminLevel(u.ID)
			roleExtra["level"] = level
		case "staff":
			dept, _ := models.GetStaffDept(u.ID)
			roleExtra["dept"] = dept
		case "customer":
			points, _ := models.GetCustomerPoints(u.ID)
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
