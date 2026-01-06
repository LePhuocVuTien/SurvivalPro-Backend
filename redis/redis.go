package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/LePhuocVuTien/SurvivalPro-Backend/internal/config"
	"github.com/redis/go-redis/v9"
)

var (
	Client *redis.Client
	Ctx    = context.Background()
)

func InitRedis() error {

	// Don't have Redis, -> skip
	if config.Cfg.RedisURL == "" {
		return fmt.Errorf("⚠️ REDIS_URL not set, skip Redis")
	}
	opt, err := redis.ParseURL(config.Cfg.RedisURL)
	if err != nil {
		return fmt.Errorf("❌ Invalid REDIS_URL: %w", err)
	}

	Client = redis.NewClient(opt)

	if err := Client.Ping(Ctx).Err(); err != nil {
		return fmt.Errorf("❌ Connot connect to Redis: %w", err)
	}

	log.Println("✅ Redis connected successfully!")
	return nil
}

func CloseRedis() {
	if Client != nil {
		Client.Close()
	}
}

func CacheGet(key string) (string, error) {
	return Client.Get(Ctx, key).Result()
}

func CacheSet(key string, value interface{}, expiration time.Duration) error {
	jsonData, err := json.Marshal(value)
	if err != nil {
		return err
	}
	return Client.Set(Ctx, key, jsonData, expiration).Err()
}

func CacheDelete(key string) error {
	return Client.Del(Ctx, key).Err()
}
