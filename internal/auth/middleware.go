package auth

import (
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
)

// JWTMiddleware validates the bearer token and injects tenant_id.
func JWTMiddleware() gin.HandlerFunc {
	secretKey := os.Getenv("JWT_SECRET")
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid authorization header format"})
			return
		}

		claims, err := ParseToken(parts[1], secretKey)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
			return
		}

		c.Set("tenant_id", claims.TenantID)
		c.Set("user_id", claims.UserID)
		c.Next()
	}
}

// GetTenantID extracts tenant context.
func GetTenantID(c *gin.Context) string {
	val, ok := c.Get("tenant_id")
	if !ok {
		return ""
	}
	str, _ := val.(string)
	return str
}
