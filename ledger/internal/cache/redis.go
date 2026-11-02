package cache

import (
	"context"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
)

var Client *redis.Client

func Connect() error {
	addr := getEnv("REDIS_ADDR", "localhost:6379")
	dbStr := getEnv("REDIS_DB", "0")
	password := os.Getenv("REDIS_PASSWORD")

	db, err := strconv.Atoi(dbStr)
	if err != nil {
		return err
	}

	Client = redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := Client.Ping(ctx).Err(); err != nil {
		return err
	}

	log.Println("Connected to Redis successfully")
	return nil
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}