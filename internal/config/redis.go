package config

import (
	"api-gateway/internal/store"
	"context"
	"encoding/json"
	"fmt"
	"github.com/redis/go-redis/v9"
	"time"
)

type FlagMessage struct {
	Flag      string `json:"flag"`
	Enabled   bool   `json:"enabled"`
	Reason    string `json:"reason"`
	Timestamp int64  `json:"timestamp"`
}

var client *redis.Client

func InitRedis(redisURL string) error {
	opts, err := redis.ParseURL(redisURL)
	if err != nil {
		return fmt.Errorf("failed to parse redis url: %w", err)
	}

	client = redis.NewClient(opts)

	maxTime := 30
	interval := 2
	elapsed := 0

	for elapsed < maxTime {
		ctx := context.Background()
		_, err := client.Ping(ctx).Result()
		if err == nil {
			fmt.Println("Redis connected")
			return nil
		}
		fmt.Printf("Redis connection failed: %v. Retrying in %ds...\n", err, interval)
		time.Sleep(time.Duration(interval) * time.Second)
		elapsed += interval
	}

	return fmt.Errorf("could not connect to Redis after 30 seconds")
}

func LoadInitialFlags() error {
	ctx := context.Background()

	result, err := client.HGetAll(ctx, "feature_flags").Result()
	if err != nil {
		return fmt.Errorf("failed to load flags: %w", err)
	}

	flagMap := map[string]bool{}
	for k, v := range result {
		flagMap[k] = v == "1"
	}

	store.LoadAll(flagMap)
	fmt.Printf("Loaded %d flags from Redis\n", len(flagMap))
	return nil
}

func StartSubscriber(ctx context.Context) {
	go func() {
		pubsub := client.Subscribe(ctx, "system:flags:changed")
		defer pubsub.Close()

		fmt.Println("Subscribed to system:flags:changed")

		for {
			select {
			case <-ctx.Done():
				fmt.Println("Subscriber stopped")
				return
			case msg := <-pubsub.Channel():
				var flagMsg FlagMessage
				if err := json.Unmarshal([]byte(msg.Payload), &flagMsg); err != nil {
					fmt.Printf("Failed to parse flag message: %v\n", err)
					continue
				}

				store.Set(flagMsg.Flag, flagMsg.Enabled)
				fmt.Printf("Flag updated: %s = %v reason: %s\n", flagMsg.Flag, flagMsg.Enabled, flagMsg.Reason)
			}
		}
	}()
}
