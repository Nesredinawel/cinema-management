package controllers

import (
	"auth-backend/models"
	"auth-backend/utils"
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// GoogleLogin → Login/Signup with Google
func GoogleLogin(c *gin.Context) {
	type Request struct {
		IDToken string `json:"id_token" binding:"required"`
	}

	var req Request
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Printf("❌ GoogleLogin bind error: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	payload, err := utils.VerifyGoogleToken(req.IDToken)
	if err != nil {
		log.Printf("❌ GoogleLogin verify token error: %v", err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid Google token"})
		return
	}

	emailVal, ok := payload.Claims["email"].(string)
	if !ok || emailVal == "" {
		log.Printf("❌ GoogleLogin missing email in token payload")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Google token missing email"})
		return
	}

	nameVal, ok := payload.Claims["name"].(string)
	if !ok || nameVal == "" {
		nameVal = "Google User"
	}

	email := strings.TrimSpace(strings.ToLower(emailVal))
	name := strings.TrimSpace(nameVal)
	googleID := payload.Subject

	log.Printf("👉 GoogleLogin called with Email=%s, Name=%s, GoogleID=%s", email, name, googleID)

	newUser := &models.User{
		Name:       name,
		Email:      &email,
		GoogleID:   &googleID,
		Role:       "customer",
		IsVerified: false, // will verify later with phone
	}

	// Extra role details map (for future extension, e.g., loyalty points)
	extra := map[string]interface{}{
		"loyalty_points": 0,
	}

	user, created, err := models.CreateOrFetchUser(newUser, extra)
	if err != nil || user == nil {
		log.Printf("❌ GoogleLogin CreateOrFetchUser error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create or fetch user"})
		return
	}

	// Generate JWT tokens
	accessToken, err := utils.GenerateAccessToken(user.ID, user.Role)
	if err != nil {
		log.Printf("❌ GoogleLogin GenerateAccessToken error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate access token"})
		return
	}

	refreshToken, err := utils.GenerateRefreshToken(user.ID)
	if err != nil {
		log.Printf("❌ GoogleLogin GenerateRefreshToken error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate refresh token"})
		return
	}

	log.Printf("✅ GoogleLogin success for user ID=%d (new=%v)", user.ID, created)
	c.JSON(http.StatusOK, gin.H{
		"access_token":             accessToken,
		"refresh_token":            refreshToken,
		"role":                     user.Role,
		"is_new_user":              created,
		"is_verified":              user.IsVerified,
		"needs_phone_verification": !user.IsVerified,
	})
}
