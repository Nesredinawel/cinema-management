package controllers

import (
	"auth-backend/models"
	"auth-backend/utils"
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// EmailAuth → Login/Signup with email + password
func EmailAuth(c *gin.Context) {
	type Request struct {
		Name     string `json:"name"` // optional for signup
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required"`
	}

	var req Request
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Printf("❌ EmailAuth bind error: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	email := strings.TrimSpace(strings.ToLower(req.Email))

	// Hash password for new users
	hash, err := utils.HashPassword(req.Password)
	if err != nil {
		log.Printf("❌ EmailAuth hash password error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}

	newUser := &models.User{
		Name:         strings.TrimSpace(req.Name),
		Email:        &email,
		PasswordHash: hash,
		Role:         "customer",
		IsVerified:   false, // will verify later with phone OTP
	}

	// For customer, you may want to init loyalty points = 0
	extra := map[string]interface{}{
		"loyalty_points": 0,
	}

	user, created, err := models.CreateOrFetchUser(newUser, extra)
	if err != nil || user == nil {
		log.Printf("❌ EmailAuth CreateOrFetchUser error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create or fetch user"})
		return
	}

	if !created {
		// Existing user → check password
		if user.PasswordHash != "" && !utils.CheckPasswordHash(req.Password, user.PasswordHash) {
			log.Printf("❌ EmailAuth invalid password for user ID=%d", user.ID)
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
			return
		}
		log.Printf("ℹ️  EmailAuth existing user login: ID=%d", user.ID)
	} else {
		log.Printf("✅ EmailAuth new user created: ID=%d", user.ID)
	}

	// Generate JWT tokens
	accessToken, err := utils.GenerateAccessToken(user.ID, user.Role)
	if err != nil {
		log.Printf("❌ EmailAuth GenerateAccessToken error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate access token"})
		return
	}

	refreshToken, err := utils.GenerateRefreshToken(user.ID)
	if err != nil {
		log.Printf("❌ EmailAuth GenerateRefreshToken error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate refresh token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"access_token":             accessToken,
		"refresh_token":            refreshToken,
		"role":                     user.Role,
		"is_new_user":              created,
		"is_verified":              user.IsVerified,
		"needs_phone_verification": !user.IsVerified,
	})
}
