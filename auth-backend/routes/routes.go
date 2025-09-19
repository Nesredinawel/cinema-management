package routes

import (
	"auth-backend/controllers"
	"auth-backend/middleware"
	"auth-backend/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(router *gin.Engine) {
	public := router.Group("/api")
	{
		// Auth endpoints (public)
		public.POST("/auth/email", controllers.EmailAuth)
		public.POST("/auth/google", controllers.GoogleLogin)
	}

	protected := router.Group("/api")
	protected.Use(middleware.AuthMiddleware())
	{
		// User profile
		protected.GET("/profile", func(c *gin.Context) {
			userID, _ := c.Get("user_id")
			role, _ := c.Get("role")
			isVerified, _ := c.Get("is_verified")
			c.JSON(http.StatusOK, gin.H{"user_id": userID, "role": role, "is_verified": isVerified})
		})

		// Phone OTP endpoints (JWT protected)
		protected.POST("/auth/phone", controllers.PhoneAuth)
		protected.POST("/auth/verify-otp", controllers.VerifyOTP)
	}

	// Admin routes
	admin := router.Group("/api/admin")
	admin.Use(middleware.AuthMiddleware(), middleware.RoleMiddleware("admin"))
	{
		admin.GET("/dashboard", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"message": "Admin dashboard"})
		})

		// Admin can create staff/admin/customer with extra details
		admin.POST("/create-user", func(c *gin.Context) {
			var body struct {
				Name         string                 `json:"name" binding:"required"`
				Email        *string                `json:"email"`
				PhoneNumber  *string                `json:"phone_number"`
				Password     string                 `json:"password" binding:"required"`
				Role         string                 `json:"role" binding:"required"` // admin/staff/customer
				ExtraDetails map[string]interface{} `json:"extra,omitempty"`         // level/dept/loyalty_points
			}

			if err := c.ShouldBindJSON(&body); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}

			user := &models.User{
				Name:         body.Name,
				Email:        body.Email,
				PhoneNumber:  body.PhoneNumber,
				PasswordHash: body.Password, // assume hashed in controller
				Role:         body.Role,
				IsVerified:   false,
			}

			createdUser, _, err := models.CreateOrFetchUser(user, body.ExtraDetails)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}

			c.JSON(http.StatusOK, gin.H{"message": "User created successfully", "user": createdUser})
		})

		// Admin can upgrade staff → admin
		admin.POST("/change-role", controllers.ChangeUserRole)

		// List all users
		admin.GET("/users", controllers.ListUsers)
	}

	// Staff routes
	staff := router.Group("/api/staff")
	staff.Use(middleware.AuthMiddleware(), middleware.RoleMiddleware("staff"))
	{
		staff.GET("/dashboard", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"message": "Staff dashboard"})
		})
	}

	// Customer routes
	customer := router.Group("/api/customer")
	customer.Use(middleware.AuthMiddleware(), middleware.RoleMiddleware("customer"))
	{
		customer.GET("/dashboard", func(c *gin.Context) {
			isVerified, _ := c.Get("is_verified")
			if !isVerified.(bool) {
				c.JSON(http.StatusOK, gin.H{"message": "Customer dashboard (limited, please verify phone)"})
				return
			}
			c.JSON(http.StatusOK, gin.H{"message": "Customer dashboard (full access)"})
		})
	}
}
