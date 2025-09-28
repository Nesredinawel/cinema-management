package routes

import (
	"booking-movie/config"
	"booking-movie/controllers"
	"booking-movie/middleware"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine, cfg *config.Config) {
	api := r.Group("/api/v1")
	{
		// Apply BookingAuthMiddleware to all booking endpoints
		api.Use(middleware.BookingAuthMiddleware())

		api.POST("/bookings", controllers.CreateBookingHandler)

	}
}
