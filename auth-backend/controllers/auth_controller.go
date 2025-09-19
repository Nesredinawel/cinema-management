package controllers

import (
	cache "auth-backend/cache-management"
	"auth-backend/models"
	"auth-backend/utils"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

// ---------------- PhoneAuth → request OTP ----------------
func PhoneAuth(c *gin.Context) {
	log.Println("📲 [PhoneAuth] Request received")

	type Request struct {
		Phone string `json:"phone" binding:"required"`
	}

	var req Request
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Printf("❌ [PhoneAuth] Bind error: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	req.Phone = strings.TrimSpace(req.Phone)
	log.Printf("📲 [PhoneAuth] Phone after trim: %s", req.Phone)

	// 🔑 Get user_id from JWT
	userIDVal, exists := c.Get("user_id")
	if !exists {
		log.Println("⚠️ [PhoneAuth] user_id not found in context")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	userID, ok := userIDVal.(int)
	if !ok {
		log.Printf("❌ [PhoneAuth] user_id type assertion failed: %#v", userIDVal)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user ID type"})
		return
	}
	log.Printf("✅ [PhoneAuth] user_id from JWT: %d", userID)

	user, err := models.GetUserByID(userID)
	if err != nil {
		log.Printf("❌ [PhoneAuth] DB error fetching user %d: %v", userID, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}
	if user == nil {
		log.Printf("❌ [PhoneAuth] User not found: %d", userID)
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}
	log.Printf("👤 [PhoneAuth] User found: ID=%d, Verified=%v, Phone=%v", user.ID, user.IsVerified, user.PhoneNumber)

	// Skip OTP if already verified
	if user.IsVerified && user.PhoneNumber != nil && *user.PhoneNumber == req.Phone {
		log.Printf("ℹ️ [PhoneAuth] User already verified with phone %s", req.Phone)
		c.JSON(http.StatusOK, gin.H{
			"message":     "User already verified, OTP not required",
			"is_verified": true,
		})
		return
	}

	// Rate limiting check (1 minute cooldown per phone)
	if !cache.CanRequestOTP(req.Phone, 1*time.Minute) {
		log.Printf("⚠️ [PhoneAuth] OTP request too soon for phone %s", req.Phone)
		c.JSON(http.StatusTooManyRequests, gin.H{"error": "OTP recently sent, please wait"})
		return
	}

	// Generate OTP
	log.Printf("🔑 [PhoneAuth] Generating OTP for %s", req.Phone)
	otp, err := utils.GenerateAndSendOTP(req.Phone)
	if err != nil {
		log.Printf("❌ [PhoneAuth] Failed to send OTP to %s: %v", req.Phone, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send OTP"})
		return
	}
	log.Printf("✅ [PhoneAuth] OTP generated: %s (not shown to user)", otp)

	// Save OTP in cache (expires in 5 minutes)
	if err := cache.SaveOTP(userID, req.Phone, otp, 5); err != nil {
		log.Printf("❌ [PhoneAuth] Failed to save OTP in cache: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save OTP"})
		return
	}
	log.Println("💾 [PhoneAuth] OTP saved in cache")

	// Save OTP history in DB
	if err := models.SaveOTPRequest(userID, req.Phone, otp); err != nil {
		log.Printf("⚠️ [PhoneAuth] Failed to save OTP request in DB: %v", err)
	} else {
		log.Println("📜 [PhoneAuth] OTP history saved in DB")
	}

	log.Printf("✅ [PhoneAuth] OTP sent to %s for user ID=%d", req.Phone, user.ID)
	c.JSON(http.StatusOK, gin.H{"message": "OTP sent for phone verification"})
}

// ---------------- VerifyOTP → verify phone ----------------
func VerifyOTP(c *gin.Context) {
	log.Println("📲 [VerifyOTP] Request received")

	type Request struct {
		Phone string `json:"phone" binding:"required"`
		OTP   string `json:"otp" binding:"required"`
	}

	var req Request
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Printf("❌ [VerifyOTP] Bind error: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	req.Phone = strings.TrimSpace(req.Phone)
	log.Printf("📲 [VerifyOTP] Phone=%s, OTP=%s", req.Phone, req.OTP)

	// 🔑 Get user_id from JWT
	userIDVal, exists := c.Get("user_id")
	if !exists {
		log.Println("⚠️ [VerifyOTP] user_id not found in context")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	userID, ok := userIDVal.(int)
	if !ok {
		log.Printf("❌ [VerifyOTP] user_id type assertion failed: %#v", userIDVal)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user ID type"})
		return
	}
	log.Printf("✅ [VerifyOTP] user_id from JWT: %d", userID)

	user, err := models.GetUserByID(userID)
	if err != nil {
		log.Printf("❌ [VerifyOTP] DB error fetching user %d: %v", userID, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}
	if user == nil {
		log.Printf("❌ [VerifyOTP] User not found: %d", userID)
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}
	log.Printf("👤 [VerifyOTP] User found: ID=%d, Verified=%v, Phone=%v", user.ID, user.IsVerified, user.PhoneNumber)

	// Already verified
	if user.IsVerified && user.PhoneNumber != nil && *user.PhoneNumber == req.Phone {
		log.Printf("ℹ️ [VerifyOTP] User already verified with phone %s", req.Phone)
		c.JSON(http.StatusOK, gin.H{
			"message":     "Account already verified",
			"is_verified": true,
			"role":        user.Role,
		})
		return
	}

	// Get OTP from cache
	cachedOTP, err := cache.GetOTP(userID, req.Phone)
	if err != nil {
		log.Printf("❌ [VerifyOTP] Failed to get OTP from cache for phone %s: %v", req.Phone, err)
	}
	log.Printf("🔑 [VerifyOTP] Cached OTP=%s, Provided OTP=%s", cachedOTP, req.OTP)

	if err != nil || cachedOTP != req.OTP {
		log.Println("❌ [VerifyOTP] Invalid or expired OTP")
		_, _ = cache.IncrementFailedOTP(userID, req.Phone)
		_ = models.MarkOTPFailed(userID, req.Phone, req.OTP)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired OTP"})
		return
	}

	// Ensure phone number uniqueness
	existingUser, _ := models.GetUserByPhone(req.Phone)
	if existingUser != nil && existingUser.ID != user.ID {
		log.Printf("⚠️ [VerifyOTP] Phone %s already used by user ID=%d", req.Phone, existingUser.ID)
		c.JSON(http.StatusConflict, gin.H{"error": "Phone number already in use"})
		return
	}

	// Update user as verified
	log.Printf("🔄 [VerifyOTP] Updating user %d as verified", user.ID)
	user.IsVerified = true
	user.PhoneNumber = &req.Phone
	if err := models.UpdateUser(user, nil); err != nil { // ✅ pass nil for extra
		log.Printf("❌ [VerifyOTP] Failed to update user %d: %v", user.ID, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user"})
		return
	}
	log.Printf("✅ [VerifyOTP] User %d updated as verified", user.ID)

	// Mark OTP as verified in DB
	if err := models.MarkOTPVerified(userID, req.Phone, req.OTP); err != nil {
		log.Printf("⚠️ [VerifyOTP] Failed to mark OTP as verified in DB: %v", err)
	}
	log.Println("📜 [VerifyOTP] OTP marked as verified in DB")

	// Delete OTP from cache
	if err := cache.DeleteOTP(userID, req.Phone); err != nil {
		log.Printf("⚠️ [VerifyOTP] Failed to delete OTP from cache: %v", err)
	} else {
		log.Println("🗑️ [VerifyOTP] OTP deleted from cache")
	}

	// Issue tokens
	accessToken, err := utils.GenerateAccessToken(user.ID, user.Role)
	if err != nil {
		log.Printf("⚠️ [VerifyOTP] Failed to generate access token: %v", err)
	}
	refreshToken, err := utils.GenerateRefreshToken(user.ID)
	if err != nil {
		log.Printf("⚠️ [VerifyOTP] Failed to generate refresh token: %v", err)
	}

	log.Printf("✅ [VerifyOTP] Phone verification successful for user %d", user.ID)
	c.JSON(http.StatusOK, gin.H{
		"access_token":  accessToken,
		"refresh_token": refreshToken,
		"role":          user.Role,
		"is_verified":   true,
		"message":       "Phone verification successful",
	})
}
