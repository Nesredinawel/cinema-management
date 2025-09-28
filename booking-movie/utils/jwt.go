package utils

import (
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var jwtSecret []byte

// InitJWT loads secret from env
func InitJWT() {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		panic("JWT_SECRET not set in environment")
	}
	jwtSecret = []byte(secret)
}

// GenerateToken creates a signed JWT with entity type and ID
func GenerateToken(entityType string, id int) (string, error) {
	claims := jwt.MapClaims{
		"entity_type": entityType,
		"id":          id,
		"issued_at":   time.Now().Unix(),
		"exp":         time.Now().Add(24 * time.Hour).Unix(),
		"issuer":      "cinema-scheduling",
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}

// VerifyEntityToken validates a token against entity type and ID
func VerifyEntityToken(entityType string, entityID int, tokenString string) error {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})
	if err != nil || !token.Valid {
		return fmt.Errorf("invalid token: %w", err)
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return fmt.Errorf("invalid claims in token")
	}

	// Check entity type
	if claims["entity_type"] != entityType {
		return fmt.Errorf("invalid entity type: expected %s", entityType)
	}

	// Check entity id
	id, ok := claims["id"].(float64)
	if !ok || int(id) != entityID {
		return fmt.Errorf("%s id mismatch", entityType)
	}

	// Optional: issuer validation
	if claims["issuer"] != "cinema-scheduling" {
		return fmt.Errorf("invalid issuer")
	}

	return nil
}
