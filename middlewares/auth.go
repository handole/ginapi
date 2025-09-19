package middlewares

import (
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

var jwtKey = []byte("secret_key_321") // ganti secret dari env

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header missing"})
			c.Abort()
			return
		}

		// format: Bearer token
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token format"})
			c.Abort()
			return
		}
		tokenString := parts[1]
		// parse token
		claims := &jwt.RegisteredClaims{}
		// parse token dengan key yg sama seperti saat generate token
		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			return jwtKey, nil
		})

		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		// simpan email di context
		c.Set("userID", claims.Subject)
		c.Next()
	}
}

// JWTMiddleware validasi token dan inject userID ke context
func JWTMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header missing"})
			c.Abort()
			return
		}

		// pastikan format: Bearer <token>
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid authorization format"})
			c.Abort()
			return
		}

		tokenString := parts[1]
		secret := []byte(os.Getenv("JWT_SECRET")) // ambil dari env

		// parse dan validasi token
		token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
			// validasi signing method
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwt.ErrSignatureInvalid
			}
			return secret, nil
		})

		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
			c.Abort()
			return
		}

		// ambil claims
		if claims, ok := token.Claims.(jwt.MapClaims); ok {
			if userID, ok := claims["userID"].(string); ok {
				// inject userID ke context
				c.Set("userID", userID)
			} else {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "userID not found in token"})
				c.Abort()
				return
			}
		}

		c.Next()
	}
}
