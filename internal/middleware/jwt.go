package middleware

import (
	"api-gateway/internal/config"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"strings"
)

func ValidateAccessJWT(cfg *config.Config) gin.HandlerFunc {
	return func(context *gin.Context) {
		authHeader := context.GetHeader("Authorization")
		if strings.TrimSpace(authHeader) == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			context.JSON(401, gin.H{"error": "invalid or missing access token"})
			context.Abort()
			return
		}

		authHeader = strings.Replace(authHeader, "Bearer ", "", 1)

		token, err := jwt.Parse(authHeader, func(token *jwt.Token) (interface{}, error) {
			return []byte(cfg.JWT_SECRET), nil
		})

		if err != nil {
			context.JSON(401, gin.H{"error": "expired or invalid access token"})
			context.Abort()
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			context.JSON(401, gin.H{"error": "invalid token claims"})
			context.Abort()
			return
		}

		userID, ok := claims["sub"].(string)
		if !ok {
			context.JSON(401, gin.H{"error": "invalid user id"})
			context.Abort()
			return
		}

		userRole, ok := claims["role"].(string)
		if !ok {
			context.JSON(401, gin.H{"error": "invalid user role"})
			context.Abort()
			return
		}

		tokenType, ok := claims["type"].(string)
		if !ok || tokenType != "access" {
			context.JSON(401, gin.H{"error": "invalid token type"})
			context.Abort()
			return
		}

		context.Set("user_id", userID)
		context.Set("user_role", userRole)
		context.Next()
	}
}
