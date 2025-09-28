package main

import (
	cache "auth-backend/cache-management"
	"auth-backend/config"
	"auth-backend/jobs"
	"auth-backend/models"
	"auth-backend/routes"
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
	log.Println("✅ Connected to Postgres")

	// ---------------- Initialize Redis ----------------
	if err := cache.Init(cfg); err != nil {
		log.Fatalf("❌ Failed to connect to Redis: %v", err)
	}
	log.Println("✅ Connected to Redis")

	// ---------------- Start background jobs ----------------
	go jobs.RunOTPCleanup()

	// ---------------- Setup HTTP routes ----------------
	router := gin.Default()
	routes.SetupRoutes(router) // no config needed here

	// ---------------- Run Server ----------------
	log.Printf("🚀 Server running on port %s", cfg.Port)
	if err := router.Run(":" + cfg.Port); err != nil {
		log.Fatalf("❌ Failed to start server: %v", err)
	}
}
