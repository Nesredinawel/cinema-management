package utils

import (
	"fmt"
	"math/rand"
	"time"
)

// GenerateAndSendOTP generates and simulates sending OTP
func GenerateAndSendOTP(phone string) (string, error) {
	rand.Seed(time.Now().UnixNano())
	code := fmt.Sprintf("%06d", rand.Intn(1000000))
	// In production: integrate SMS provider (Twilio, etc.)
	fmt.Printf("Sending OTP %s to phone %s\n", code, phone)
	return code, nil
}
