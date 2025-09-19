package routes

import (
	"cinema-scheduling/controllers"
	"cinema-scheduling/middleware"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(router *gin.Engine) {
	// ---------------- Admin & Manager Routes ----------------
	adminGroup := router.Group("/api/admin")
	adminGroup.Use(middleware.JWTAuthMiddleware("admin", "manager"))
	{
		// ---------------- Movies ----------------
		adminGroup.POST("/movies", controllers.AddMovie)
		adminGroup.GET("/movies", controllers.ListMovies)
		adminGroup.GET("/movies/:movie_id", controllers.GetMovie)
		adminGroup.PUT("/movies/:movie_id", controllers.UpdateMovie)    // update movie
		adminGroup.DELETE("/movies/:movie_id", controllers.DeleteMovie) // delete movie

		// ---------------- Schedules ----------------
		adminGroup.POST("/schedules", controllers.AddSchedule)
		adminGroup.GET("/schedules", controllers.ListSchedules)
		adminGroup.GET("/schedules/:schedule_id", controllers.GetSchedule)
		adminGroup.PUT("/schedules/:schedule_id", controllers.UpdateSchedule)    // update schedule
		adminGroup.DELETE("/schedules/:schedule_id", controllers.DeleteSchedule) // delete schedule

		// ---------------- Snacks ----------------
		adminGroup.POST("/snacks", controllers.AddSnack)
		adminGroup.GET("/snacks", controllers.ListSnacks)
		adminGroup.GET("/snacks/:snack_id", controllers.GetSnack)
		adminGroup.PUT("/snacks/:snack_id", controllers.UpdateSnack)    // update snack
		adminGroup.DELETE("/snacks/:snack_id", controllers.DeleteSnack) // delete snack

		// ---------------- Halls ----------------
		adminGroup.POST("/halls", controllers.AddHall)
		adminGroup.GET("/halls", controllers.ListHalls)
		adminGroup.GET("/halls/:hall_id", controllers.GetHall)
		adminGroup.PUT("/halls/:hall_id", controllers.UpdateHall)    // update hall
		adminGroup.DELETE("/halls/:hall_id", controllers.DeleteHall) // delete hall

		// ---------------- Schedule-specific Snacks ----------------
		adminGroup.POST("/schedules/:schedule_id/snacks", controllers.AddScheduleSnack)
		adminGroup.GET("/schedules/:schedule_id/snacks", controllers.ListScheduleSnacks)
		adminGroup.PUT("/schedules/snacks/:schedule_snack_id", controllers.UpdateScheduleSnack)
		adminGroup.DELETE("/schedules/snacks/:schedule_snack_id", controllers.DeleteScheduleSnack)
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

		// Snacks (general list)
		publicGroup.GET("/snacks", controllers.ListSnacks)
		publicGroup.GET("/snacks/:snack_id", controllers.GetSnack)

		// Halls
		publicGroup.GET("/halls", controllers.ListHalls)
		publicGroup.GET("/halls/:hall_id", controllers.GetHall)
	}
}
