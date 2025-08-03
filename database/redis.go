package database

import (
	"context"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
)

var RD *redis.Client

func RedisConnect() {
	db, err := strconv.Atoi(os.Getenv("REDIS_DB"))
	if err != nil {
		log.Fatalf("Invalid REDIS_DB: %v", err)
	}

	protocol, err := strconv.Atoi(os.Getenv("REDIS_PROTOCOL"))
	if err != nil {
		log.Fatalf("Invalid REDIS_PROTOCOL: %v", err)
	}

	RD = redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_ADDR"),
		Password: os.Getenv("REDIS_PASSWORD"),
		DB:       db,
		Protocol: protocol,
	})

	ctx := context.Background()
	err = RD.Set(ctx, "check", "success", 5*time.Second).Err()
	if err != nil {
		panic(err)
	}
}
