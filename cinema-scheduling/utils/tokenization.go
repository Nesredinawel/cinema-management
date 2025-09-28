package utils

import (
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var jwtSecret []byte

// InitJWT loads the JWT secret from environment
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

// VerifyToken checks the JWT and returns the entity type and ID
func VerifyToken(tokenString string) (string, int, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})
	if err != nil || !token.Valid {
		return "", 0, fmt.Errorf("invalid token: %w", err)
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "", 0, fmt.Errorf("invalid claims")
	}

	entityType, ok := claims["entity_type"].(string)
	if !ok {
		return "", 0, fmt.Errorf("entity_type missing in token")
	}

	idFloat, ok := claims["id"].(float64)
	if !ok {
		return "", 0, fmt.Errorf("id missing in token")
	}

	return entityType, int(idFloat), nil
}
