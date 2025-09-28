package main

import (
	"cinema-scheduling/config"
	"cinema-scheduling/models"
	"cinema-scheduling/routes"
	"context"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	// ---------------- Load Config ----------------
	cfg := config.LoadConfig()

	// ---------------- Connect to Postgres ----------------
	var err error
	models.DB, err = pgxpool.New(context.Background(), cfg.PostgresURL)
	if err != nil {
		log.Fatalf("❌ Failed to connect to database: %v", err)
	}
	log.Println("✅ Connected to Postgres (Cinema Scheduling)")

	// ---------------- Setup HTTP routes ----------------
	router := gin.Default()
	routes.SetupRoutes(router, cfg)

	// ---------------- Run Server ----------------
	log.Printf("🎬 Cinema Scheduling service running on port %s", cfg.Port)
	if err := router.Run(":" + cfg.Port); err != nil {
		log.Fatalf("❌ Failed to start server: %v", err)
	}
}
