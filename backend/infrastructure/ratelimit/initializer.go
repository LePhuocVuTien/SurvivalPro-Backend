package ratelimit

import (
	"context"
	"log"
	"time"

	"github.com/LePhuocVuTien/SurvivalPro-Backend/internal/config"
	"github.com/redis/go-redis/v9"
)

// ============================================================================
// REDIS LIMITER INITIALIZATION
// ============================================================================

func InitializeRedisLimiter() *RedisLimiter {
	// Create Redis client
	opt, err := redis.ParseURL(config.Cfg.Redis.URL)
	if err != nil {
		log.Fatalf("❌ Invalid REDIS_URL: %v", err)
	}

	// Create Redis client with connection pool
	client := redis.NewClient(&redis.Options{
		Addr:         opt.Addr,
		Password:     opt.Password,
		DB:           opt.DB,
		PoolSize:     config.Cfg.Redis.PoolSize,
		MinIdleConns: 5,
		MaxRetries:   3,
	})

	// Test connection
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := client.Ping(ctx).Err(); err != nil {
		log.Fatalf("❌ Cannot connect to Redis: %v", err)
	}

	// Define rate limit rules
	rules := []*RateLimitRule{
		// Login: 5 attempts per 5 minutes, block for 30 minutes
		{
			Action:        "login",
			MaxAttempts:   5,
			WindowSize:    5 * time.Minute,
			BlockDuration: 30 * time.Minute,
			IsActive:      true,
		},
		// Register: 3 attempts per hour, block for 2 hours
		{
			Action:        "register",
			MaxAttempts:   3,
			WindowSize:    1 * time.Hour,
			BlockDuration: 2 * time.Hour,
			IsActive:      true,
		},
		// Password reset: 3 attempts per hour, block for 1 hour
		{
			Action:        "password_reset",
			MaxAttempts:   3,
			WindowSize:    1 * time.Hour,
			BlockDuration: 1 * time.Hour,
			IsActive:      true,
		},
		// Email verify: 5 attempts per hour, block for 1 hour
		{
			Action:        "email_verify",
			MaxAttempts:   5,
			WindowSize:    1 * time.Hour,
			BlockDuration: 1 * time.Hour,
			IsActive:      true,
		},
		// API: 100 requests per minute, block for 5 minutes
		{
			Action:        "api",
			MaxAttempts:   100,
			WindowSize:    1 * time.Minute,
			BlockDuration: 5 * time.Minute,
			IsActive:      true,
		},
		// Upload: 10 per hour, block for 1 hour
		{
			Action:        "upload",
			MaxAttempts:   10,
			WindowSize:    1 * time.Hour,
			BlockDuration: 1 * time.Hour,
			IsActive:      true,
		},
	}

	// Create rate limiter
	return NewRedisLimiter(client, rules)
}
