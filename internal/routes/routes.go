package routes

import (
	"api-gateway/internal/config"
	"api-gateway/internal/handlers"
	"api-gateway/internal/middleware"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine, cfg *config.Config) {
	api := r.Group("/api/v1/")

	api.Any("/exam/*path", middleware.ValidateJWT(cfg), handlers.ForwardToExamService(cfg))
}
