package auth

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"net/http"
	"os"
	"strings"
)

// AuthMiddleware membuat gin middleware untuk autentikasi JWT.
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Ambil header Authorization
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is required"})
			return
		}

		// Cek format header "Bearer <token>"
		bearerToken := strings.Split(authHeader, " ")
		if len(bearerToken) != 2 || bearerToken[0] != "Bearer" {
			c.AbortWithStatusJSON(
				http.StatusUnauthorized,
				gin.H{"error": "Invalid token format, must be 'Bearer <token>'"},
			)
			return
		}

		// Validasi token
		tokenString := bearerToken[1]
		token, err := validateToken(tokenString)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok || !token.Valid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			return
		}

		// Ekstrak data pengguna dari claims token
		userID, okUserID := claims["user_id"].(float64) // JWT mengurai angka sebagai float64
		role, okRole := claims["role"].(string)

		if !okUserID || !okRole {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token claims"})
			return
		}

		// Set data pengguna ke dalam context Gin untuk digunakan di handler selanjutnya
		c.Set("userID", uint(userID))
		c.Set("userRole", role)

		c.Next() // Lanjutkan ke handler berikutnya jika token valid
	}
}

// validateToken mem-parsing dan memvalidasi string token JWT.
func validateToken(tokenString string) (*jwt.Token, error) {
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		return nil, fmt.Errorf("JWT_SECRET environment variable not set")
	}

	return jwt.Parse(
		tokenString, func(token *jwt.Token) (interface{}, error) {
			// Validasi algoritma signing
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(jwtSecret), nil
		},
	)
}

// AuthorizeRole membuat middleware untuk memeriksa peran (role) pengguna.
func AuthorizeRole(requiredRole string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userRole, exists := c.Get("userRole")
		if !exists {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "User role not found in context"})
			return
		}

		if userRole.(string) != requiredRole {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "You are not authorized to perform this action"})
			return
		}

		c.Next()
	}
}
