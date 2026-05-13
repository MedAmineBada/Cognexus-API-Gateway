package middleware

import (
	"api-gateway/internal/config"
	"github.com/gin-gonic/gin"
)

func ValidateRBAC() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		role, exists := ctx.Get("user_role")
		if !exists {
			ctx.JSON(401, gin.H{"error": "unauthorized"})
			ctx.Abort()
			return
		}

		normalizedPath := config.NormalizePath(ctx.Request.URL.Path)
		key := ctx.Request.Method + ":" + normalizedPath

		allowedRoles, exists := config.RouteRoles[key]
		if !exists {
			ctx.Next()
			return
		}

		for _, allowedRole := range allowedRoles {
			if allowedRole == role.(string) {
				ctx.Next()
				return
			}
		}

		ctx.JSON(403, gin.H{"error": "forbidden"})
		ctx.Abort()
	}
}
