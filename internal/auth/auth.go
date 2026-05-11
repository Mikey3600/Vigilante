package auth

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// Claims represents our custom JWT claims containing tenant info.
type Claims struct {
	TenantID string `json:"tenant_id"`
	UserID   string `json:"user_id"`
	jwt.RegisteredClaims
}

// GenerateToken creates a new valid JWT token.
func GenerateToken(secretKey, tenantID, userID string, duration time.Duration) (string, error) {
	claims := Claims{
		TenantID: tenantID,
		UserID:   userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(duration)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secretKey))
}

// ParseToken validates a token and returns its claims.
func ParseToken(tokenString, secretKey string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(secretKey), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}
	return nil, errors.New("invalid token")
}
