package redis

import (
	"context"
	"fmt"
	"log"
	"r-builder/config"
	"time"

	"github.com/go-redis/redis/v8"
)

var ctx = context.Background()
var rdb *redis.Client

func InitRedis() {
	cfg := config.LoadConfig()
	rdb = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", cfg.Redis_Host, cfg.Redis_Port),
		Password: "",
		DB:       0,
	})

	_, err := rdb.Ping(ctx).Result()
	if err != nil {
		log.Fatalf("Could not connect to Redis: %v", err)
	}
}

func GetVersion(key string) (string, error) {
	val, err := rdb.Get(ctx, key).Result()
	if err == redis.Nil {
		return "", nil
	}
	if err != nil {
		log.Printf("Error getting value from Redis: %v", err)
		return "", err
	}
	return val, nil
}

func SetVersion(key string, value interface{}, expiration time.Duration) error {
	err := rdb.Set(ctx, key, value, expiration).Err()
	if err != nil {
		log.Printf("Error setting value in Redis: %v", err)
		return err
	}
	return nil
}

func DeleteVersion(key string) error {
	err := rdb.Del(ctx, key).Err()
	if err != nil {
		log.Printf("Error deleting key from Redis: %v", err)
		return err
	}
	return nil
}
