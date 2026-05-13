package main

import (
	"api-gateway/internal/config"
	"api-gateway/internal/routes"
	"context"
	"fmt"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()

	cfg := config.Load()

	if err := config.InitRedis(cfg.REDIS_URL); err != nil {
		fmt.Printf("Failed to connect to Redis: %v\n", err)
		os.Exit(1)
	}

	if err := config.LoadInitialFlags(); err != nil {
		fmt.Printf("Failed to load flags: %v\n", err)
		os.Exit(1)
	}

	ctx := context.Background()
	config.StartSubscriber(ctx)

	r := gin.Default()
	routes.SetupRoutes(r, cfg)
	r.Run(":" + cfg.PORT)
}
