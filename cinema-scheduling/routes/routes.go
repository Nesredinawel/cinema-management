package routes

import (
	"cinema-scheduling/config"
	"cinema-scheduling/controllers"
	"cinema-scheduling/middleware"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(router *gin.Engine, cfg *config.Config) {
	// ---------------- Admin & Manager Routes ----------------
	adminGroup := router.Group("/api/admin")
	adminGroup.Use(middleware.JWTAuthMiddleware(cfg.JWTSecret, "admin", "manager"))
	{
		// ---------------- Movies ----------------
		adminGroup.POST("/movies", controllers.AddMovie)
		adminGroup.GET("/movies", controllers.ListMovies)
		adminGroup.GET("/movies/:movie_id", controllers.GetMovie)
		adminGroup.PUT("/movies/:movie_id", controllers.UpdateMovie)
		adminGroup.DELETE("/movies/:movie_id", controllers.DeleteMovie)

		// ---------------- Genres ----------------
		adminGroup.POST("/genres", controllers.AddGenre)
		adminGroup.GET("/genres", controllers.ListGenres)
		adminGroup.GET("/genres/:genre_id", controllers.GetGenre)
		adminGroup.PUT("/genres/:genre_id", controllers.UpdateGenre)
		adminGroup.DELETE("/genres/:genre_id", controllers.DeleteGenre)

		// ---------------- Schedules ----------------
		adminGroup.POST("/schedules", controllers.AddSchedule)
		adminGroup.GET("/schedules", controllers.ListSchedules)
		adminGroup.GET("/schedules/:schedule_id", controllers.GetSchedule)
		adminGroup.PUT("/schedules/:schedule_id", controllers.UpdateSchedule)
		adminGroup.DELETE("/schedules/:schedule_id", controllers.DeleteSchedule)

		// ---------------- Snacks ----------------
		adminGroup.POST("/snacks", controllers.AddSnack)
		adminGroup.GET("/snacks", controllers.ListSnacks)
		adminGroup.GET("/snacks/:snack_id", controllers.GetSnack)
		adminGroup.PUT("/snacks/:snack_id", controllers.UpdateSnack)
		adminGroup.DELETE("/snacks/:snack_id", controllers.DeleteSnack)

		// ---------------- Halls ----------------
		adminGroup.POST("/halls", controllers.AddHall)
		adminGroup.GET("/halls", controllers.ListHalls)
		adminGroup.GET("/halls/:hall_id", controllers.GetHall)
		adminGroup.PUT("/halls/:hall_id", controllers.UpdateHall)
		adminGroup.DELETE("/halls/:hall_id", controllers.DeleteHall)

		// ---------------- Schedule-specific Snacks ----------------
		adminGroup.POST("/schedules/:schedule_id/snacks", controllers.AddScheduleSnack)
		adminGroup.GET("/schedules/:schedule_id/snacks", controllers.ListScheduleSnacks)
		// Updated to fetch specific snack for a schedule
		adminGroup.GET("/schedules/:schedule_id/snacks/:snack_id", controllers.GetScheduleSnack)
		adminGroup.PUT("/schedules/:schedule_id/snacks/:snack_id", controllers.UpdateScheduleSnack)
		adminGroup.DELETE("/schedules/:schedule_id/snacks/:snack_id", controllers.DeleteScheduleSnack)
	}

	// ---------------- Public Routes ----------------
	publicGroup := router.Group("/api")
	{
		// Movies
		publicGroup.GET("/movies", controllers.ListMovies)
		publicGroup.GET("/movies/:movie_id", controllers.GetMovie)

		// Schedules
		publicGroup.GET("/schedules", controllers.ListSchedules)
		publicGroup.GET("/schedules/:schedule_id", controllers.GetSchedule)

		// Snacks
		publicGroup.GET("/snacks", controllers.ListSnacks)
		publicGroup.GET("/snacks/:snack_id", controllers.GetSnack)

		// Halls
		publicGroup.GET("/halls", controllers.ListHalls)
		publicGroup.GET("/halls/:hall_id", controllers.GetHall)

		// Schedule-specific Snacks
		publicGroup.GET("/schedules/:schedule_id/snacks", controllers.ListScheduleSnacks)
		// Updated public route to fetch specific snack for a schedule
		publicGroup.GET("/schedules/:schedule_id/snacks/:snack_id", controllers.GetScheduleSnack)
	}
}
