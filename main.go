package main

import (
	"api-gateway/internal/config"
	"api-gateway/internal/routes"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()

	r := gin.Default()
	cfg := config.Load()

	routes.SetupRoutes(r, cfg)

	r.Run(":" + cfg.PORT)
}
