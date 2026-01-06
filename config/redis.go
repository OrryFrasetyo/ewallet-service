package config

import (
	"context"
	"fmt"
	"os"

	"github.com/redis/go-redis/v9"
)

func ConnectRedis() *redis.Client {
	dbHost := os.Getenv("REDIS_HOST")
	dbPort := os.Getenv("REDIS_PORT")
	dbPassword := os.Getenv("REDIS_PASSWORD")

	// fallback if run manual no docker (optional safety)
	if dbHost == "" {
		dbHost = "localhost"
		dbPort = "6379"
	}

	addr := fmt.Sprintf("%s:%s", dbHost, dbPort)

	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: dbPassword,
		DB:       0,
	})

	// test ping
	ctx := context.Background()
	_, err := client.Ping(ctx).Result()
	if err != nil {
		fmt.Printf("❌ Gagal connect ke Redis di %s: %v\n", addr, err)
		return nil
	}

	fmt.Println("✅ Connected to Redis successfully!")
	return client
}
