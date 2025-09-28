package config

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port           string
	DBHost         string
	DBPort         string
	DBUser         string
	DBPassword     string
	DBName         string
	PostgresURL    string
	JWTSecret      string
	JWTExpiryHours int
	GoogleClientID string
	RedisHost      string
	RedisPort      string
	RedisPassword  string
}

// LoadConfig reads environment variables and returns a Config struct
func LoadConfig() *Config {
	// Load .env file if exists
	if err := godotenv.Load(); err != nil {
		log.Println(".env file not found, using system environment variables")
	}

	port := getEnv("PORT", "8081")
	dbHost := getEnv("DB_HOST", "localhost")
	dbPort := getEnv("DB_PORT", "5432")
	dbUser := getEnv("DB_USER", "postgres")
	dbPassword := getEnv("DB_PASSWORD", "")
	dbName := getEnv("DB_NAME", "cinema_auth")
	jwtSecret := getEnv("JWT_SECRET", "secret")
	googleClientID := getEnv("GOOGLE_CLIENT_ID", "")
	redisHost := getEnv("REDIS_HOST", "localhost")
	redisPort := getEnv("REDIS_PORT", "6379")
	redisPassword := getEnv("REDIS_PASSWORD", "")

	// Build Postgres URL
	postgresURL := fmt.Sprintf("postgres://%s:%s@%s:%s/%s",
		dbUser, dbPassword, dbHost, dbPort, dbName,
	)

	return &Config{
		Port:           port,
		DBHost:         dbHost,
		DBPort:         dbPort,
		DBUser:         dbUser,
		DBPassword:     dbPassword,
		DBName:         dbName,
		PostgresURL:    postgresURL,
		JWTSecret:      jwtSecret,
		JWTExpiryHours: 72,
		GoogleClientID: googleClientID,
		RedisHost:      redisHost,
		RedisPort:      redisPort,
		RedisPassword:  redisPassword,
	}
}

// getEnv returns env variable or fallback value
func getEnv(key, fallback string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return fallback
}
