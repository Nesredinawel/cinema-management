package models

import (
	"context"
	"fmt"
	"time"
)

// OTP History Model
type OTP struct {
	ID             int
	UserID         int
	Phone          string
	Code           string
	Status         string // "SENT", "VERIFIED", "FAILED", "EXPIRED"
	FailedAttempts int
	CreatedAt      time.Time
	VerifiedAt     *time.Time
}

// ---------------- Save OTP Request to History ----------------
func SaveOTPRequest(userID int, phone, code string) error {
	query := `
        INSERT INTO otp_history (user_id, phone, code, status, failed_attempts, created_at)
        VALUES ($1, $2, $3, 'SENT', 0, NOW())`
	_, err := DB.Exec(context.Background(), query, userID, phone, code)
	return err
}

// ---------------- Mark OTP Verified ----------------
func MarkOTPVerified(userID int, phone, code string) error {
	query := `
        UPDATE otp_history
        SET status='VERIFIED', verified_at=NOW()
        WHERE user_id=$1 AND phone=$2 AND code=$3
        ORDER BY created_at DESC
        LIMIT 1`
	_, err := DB.Exec(context.Background(), query, userID, phone, code)
	return err
}

// ---------------- Mark OTP Failed Attempt ----------------
func MarkOTPFailed(userID int, phone, code string) error {
	query := `
        UPDATE otp_history
        SET failed_attempts = failed_attempts + 1, status='FAILED'
        WHERE user_id=$1 AND phone=$2 AND code=$3
        ORDER BY created_at DESC
        LIMIT 1`
	_, err := DB.Exec(context.Background(), query, userID, phone, code)
	return err
}

// ---------------- Mark OTP Expired ----------------
// Update OTPs older than N minutes in otp_history
func MarkExpiredOTPs(expireMinutes int) (int64, error) {
	interval := fmt.Sprintf("%d minutes", expireMinutes)
	cmdTag, err := DB.Exec(context.Background(), `
        UPDATE otp_history
        SET status = 'EXPIRED'
        WHERE status = 'SENT' AND created_at < NOW() - $1::INTERVAL
    `, interval)

	if err != nil {
		return 0, err
	}
	return cmdTag.RowsAffected(), nil
}

// ---------------- Delete Old OTPs ----------------
// Delete OTP history older than N days
func DeleteOldOTPs(retentionDays int) (int64, error) {
	interval := fmt.Sprintf("%d days", retentionDays)
	cmdTag, err := DB.Exec(context.Background(), `
        DELETE FROM otp_history WHERE created_at < NOW() - $1::INTERVAL
    `, interval)

	if err != nil {
		return 0, err
	}
	return cmdTag.RowsAffected(), nil
}
