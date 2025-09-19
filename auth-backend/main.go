package main

import (
	"auth-backend/cache-management"
	"auth-backend/jobs"
	"auth-backend/models"
	"auth-backend/routes"
	"context"
	"fmt"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	// ---------------- Build Postgres URL ----------------
	pgURL := fmt.Sprintf("postgres://%s:%s@%s:%s/%s",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_NAME"),
	)

	// ---------------- Connect to Postgres ----------------
	var err error
	models.DB, err = pgxpool.New(context.Background(), pgURL)
	if err != nil {
		log.Fatalf("❌ Failed to connect to database: %v", err)
	}
	log.Println("✅ Connected to Postgres")

	// ---------------- Initialize Redis ----------------
	if err := cache.Init(); err != nil {
		log.Fatalf("❌ Failed to connect to Redis: %v", err)
	}
	log.Println("✅ Connected to Redis")

	// ---------------- Start background jobs ----------------
	go jobs.RunOTPCleanup()

	// ---------------- Setup HTTP routes ----------------
	router := gin.Default()
	routes.SetupRoutes(router)

	// ---------------- Run Server ----------------
	port := os.Getenv("PORT")
	if port == "" {
		port = "8081"
	}
	log.Printf("🚀 Server running on port %s", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatalf("❌ Failed to start server: %v", err)
	}
}
