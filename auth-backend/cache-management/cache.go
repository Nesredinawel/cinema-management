package cache

import (
	"context"
	"fmt"
	"log"
	"time"

	"auth-backend/config"

	"github.com/go-redis/redis/v8"
)

var ctx = context.Background()
var rdb *redis.Client

// Init initializes Redis client using values from Config
func Init(cfg *config.Config) error {
	addr := fmt.Sprintf("%s:%s", cfg.RedisHost, cfg.RedisPort)
	rdb = redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: cfg.RedisPassword, // optional
		DB:       0,
	})

	// Test connection
	if err := rdb.Ping(ctx).Err(); err != nil {
		return fmt.Errorf("failed to connect to Redis at %s: %v", addr, err)
	}

	log.Printf("✅ Connected to Redis at %s", addr)
	return nil
}

// ensureClient prevents nil pointer panic
func ensureClient() bool {
	if rdb == nil {
		log.Println("⚠️ Redis client is nil. Did you call cache.Init(cfg) in main.go?")
		return false
	}
	return true
}

// SaveOTP saves OTP with TTL
func SaveOTP(userID int, phone, otp string, ttlMinutes int) error {
	if !ensureClient() {
		return fmt.Errorf("redis client not initialized")
	}
	key := fmt.Sprintf("otp:%d:%s", userID, phone)
	return rdb.Set(ctx, key, otp, time.Duration(ttlMinutes)*time.Minute).Err()
}

// GetOTP retrieves OTP
func GetOTP(userID int, phone string) (string, error) {
	if !ensureClient() {
		return "", fmt.Errorf("redis client not initialized")
	}
	key := fmt.Sprintf("otp:%d:%s", userID, phone)
	return rdb.Get(ctx, key).Result()
}

// DeleteOTP removes OTP
func DeleteOTP(userID int, phone string) error {
	if !ensureClient() {
		return fmt.Errorf("redis client not initialized")
	}
	key := fmt.Sprintf("otp:%d:%s", userID, phone)
	return rdb.Del(ctx, key).Err()
}

// CanRequestOTP checks if OTP was requested recently
func CanRequestOTP(phone string, cooldown time.Duration) bool {
	if !ensureClient() {
		return true // fail-open (allow request) instead of crashing
	}

	key := fmt.Sprintf("otp_request:%s", phone)
	ttl, err := rdb.TTL(ctx, key).Result()
	if err == nil && ttl > 0 {
		return false
	}

	// Set cooldown marker
	if err := rdb.Set(ctx, key, "1", cooldown).Err(); err != nil {
		log.Printf("⚠️ Redis SET error in CanRequestOTP: %v", err)
	}
	return true
}

// IncrementFailedOTP increments failed attempts counter
func IncrementFailedOTP(userID int, phone string) (int64, error) {
	if !ensureClient() {
		return 0, fmt.Errorf("redis client not initialized")
	}

	key := fmt.Sprintf("otp_failed:%d:%s", userID, phone)
	count, err := rdb.Incr(ctx, key).Result()
	if err != nil {
		return 0, err
	}

	if count == 1 {
		// Expire counter after 5 min
		if err := rdb.Expire(ctx, key, 5*time.Minute).Err(); err != nil {
			log.Printf("⚠️ Redis EXPIRE error: %v", err)
		}
	}
	return count, nil
}
