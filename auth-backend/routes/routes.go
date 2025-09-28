package routes

import (
	"auth-backend/controllers"
	"auth-backend/middleware"
	"net/http"

	"github.com/gin-gonic/gin"
)

// SetupRoutes sets up all HTTP routes
func SetupRoutes(router *gin.Engine) {
	// ---------------- Public routes ----------------
	public := router.Group("/api")
	{
		public.POST("/auth/email", controllers.EmailAuth)
		public.POST("/auth/google", controllers.GoogleLogin) // reads GOOGLE_CLIENT_ID internally
	}

	// ---------------- Protected routes (JWT) ----------------
	protected := router.Group("/api")
	protected.Use(middleware.AuthMiddleware()) // reads JWT_SECRET internally
	{
		// User profile
		protected.GET("/profile", func(c *gin.Context) {
			userID, _ := c.Get("user_id")
			role, _ := c.Get("role")
			isVerified, _ := c.Get("is_verified")
			c.JSON(http.StatusOK, gin.H{"user_id": userID, "role": role, "is_verified": isVerified})
		})

		// Phone OTP endpoints
		protected.POST("/auth/phone", controllers.PhoneAuth)
		protected.POST("/auth/verify-otp", controllers.VerifyOTP)
	}

	// ---------------- Admin routes ----------------
	admin := router.Group("/api/admin")
	admin.Use(middleware.AuthMiddleware(), middleware.RoleMiddleware("admin"))
	{
		admin.GET("/dashboard", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"message": "Admin dashboard"})
		})

		admin.POST("/create-user", controllers.CreateUserByAdmin)
		admin.POST("/change-role", controllers.ChangeUserRole)
		admin.GET("/users", controllers.ListUsers)
	}

	// ---------------- Staff routes ----------------
	staff := router.Group("/api/staff")
	staff.Use(middleware.AuthMiddleware(), middleware.RoleMiddleware("staff"))
	{
		staff.GET("/dashboard", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"message": "Staff dashboard"})
		})
	}

	// ---------------- Customer routes ----------------
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
