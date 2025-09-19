// jobs/otp_cleanup.go
package jobs

import (
	"auth-backend/models"
	"log"
	"os"
	"strconv"
	"time"
)

func getExpireMinutes() int {
	expireMinutes := 5
	if val := os.Getenv("OTP_EXPIRE_MINUTES"); val != "" {
		if v, err := strconv.Atoi(val); err == nil {
			expireMinutes = v
		}
	}
	return expireMinutes
}

func getRetentionDays() int {
	retention := 5
	if val := os.Getenv("OTP_RETENTION_DAYS"); val != "" {
		if v, err := strconv.Atoi(val); err == nil {
			retention = v
		}
	}
	return retention
}

// RunOTPCleanup runs periodically to mark expired OTPs and delete old history
func RunOTPCleanup() {
	ticker := time.NewTicker(1 * time.Minute)

	go func() {
		for range ticker.C {
			expireMinutes := getExpireMinutes()
			retentionDays := getRetentionDays()

			// 1. Mark expired OTPs
			if count, err := models.MarkExpiredOTPs(expireMinutes); err != nil {
				log.Printf("❌ OTP expiration update failed: %v", err)
			} else {
				log.Printf("✅ Marked %d OTPs expired (older than %d minutes)", count, expireMinutes)
			}

			// 2. Delete old OTP history
			if count, err := models.DeleteOldOTPs(retentionDays); err != nil {
				log.Printf("❌ Old OTP deletion failed: %v", err)
			} else {
				log.Printf("🗑️ Deleted %d OTP log entries older than %d days", count, retentionDays)
			}
		}
	}()
}
