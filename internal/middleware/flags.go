package middleware

import (
	"api-gateway/internal/config"
	"api-gateway/internal/store"

	"github.com/gin-gonic/gin"
)

func CheckFlag() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		flagName := config.GetFlagForPath(ctx.Request.URL.Path, ctx.Request.Method)

		if flagName == "" {
			ctx.Next()
			return
		}

		if !store.IsEnabled(flagName) {
			ctx.JSON(503, gin.H{"error": "this feature is currently unavailable"})
			ctx.Abort()
			return
		}

		ctx.Next()
	}
}
