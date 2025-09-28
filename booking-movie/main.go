package main

import (
	"booking-movie/config"
	"booking-movie/models"
	"booking-movie/routes"
	"context"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	cfg := config.LoadConfig()

	var err error
	models.DB, err = pgxpool.New(context.Background(), cfg.PostgresURL)
	if err != nil {
		log.Fatalf("❌ Failed to connect DB: %v", err)
	}
	log.Println("✅ Connected to Postgres (Cinema Booking)")

	router := gin.Default()
	routes.SetupRoutes(router, cfg)

	log.Printf("🎟️ Cinema Booking service running on port %s", cfg.Port)
	if err := router.Run(":" + cfg.Port); err != nil {
		log.Fatalf("❌ Failed to start server: %v", err)
	}
}
