package redis

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/LePhuocVuTien/SurvivalPro-Backend/internal/config"
	"github.com/redis/go-redis/v9"
)

var (
	Client *redis.Client
	ctx    = context.Background()

	// Common errors
	ErrNotFound      = errors.New("key not found in cache")
	ErrMarshalJSON   = errors.New("failed to marshal data to JSON")
	ErrUnmarshalJSON = errors.New("failed to unmarshal JSON data")
)

// CacheConfig holds cache-specific configuration
type CacheConfig struct {
	DefaultExpiration time.Duration
	CleanupInterval   time.Duration
}

var cacheConfig = CacheConfig{
	DefaultExpiration: 5 * time.Minute,
	CleanupInterval:   10 * time.Minute,
}

// InitRedis initializes Redis connection with configuration
func InitRedis() error {
	// Check if Redis URL is configured
	if config.Cfg.Redis.URL == "" {
		log.Println("⚠️  REDIS_URL not set, Redis features will be disabled")
		return fmt.Errorf("Redis URL not configured")
	}

	// Parse Redis URL
	opt, err := redis.ParseURL(config.Cfg.Redis.URL)
	if err != nil {
		return fmt.Errorf("❌ Invalid REDIS_URL: %w", err)
	}

	// Override with config if provided
	if config.Cfg.Redis.PoolSize > 0 {
		opt.PoolSize = config.Cfg.Redis.PoolSize
	}

	// Create Redis client
	Client = redis.NewClient(opt)

	// Test connection with timeout
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	if err := Client.Ping(ctx).Err(); err != nil {
		return fmt.Errorf("❌ Cannot connect to Redis: %w", err)
	}

	// Log connection info
	log.Println("✅ Redis connected successfully!")
	log.Printf("   Host: %s:%s", config.Cfg.Redis.Host, config.Cfg.Redis.Port)
	log.Printf("   DB: %d", config.Cfg.Redis.DB)
	log.Printf("   Pool Size: %d", opt.PoolSize)

	return nil
}

// CloseRedis gracefully closes Redis connection
func CloseRedis() error {
	if Client != nil {
		if err := Client.Close(); err != nil {
			return fmt.Errorf("failed to close Redis connection: %w", err)
		}
		log.Println("✅ Redis connection closed")
	}
	return nil
}

// IsConnected checks if Redis is connected
func IsConnected() bool {
	if Client == nil {
		return false
	}

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	return Client.Ping(ctx).Err() == nil
}

// ============================================================================
// Basic Cache Operations
// ============================================================================

// Get retrieves a value from cache as string
func Get(key string) (string, error) {
	if Client == nil {
		return "", fmt.Errorf("Redis client not initialized")
	}

	val, err := Client.Get(ctx, key).Result()
	if err == redis.Nil {
		return "", ErrNotFound
	}
	if err != nil {
		return "", fmt.Errorf("failed to get key %s: %w", key, err)
	}

	return val, nil
}

// Set stores a string value in cache
func Set(key string, value interface{}, expiration time.Duration) error {
	if Client == nil {
		return fmt.Errorf("Redis client not initialized")
	}

	return Client.Set(ctx, key, value, expiration).Err()
}

// Delete removes a key from cache
func Delete(key string) error {
	if Client == nil {
		return fmt.Errorf("Redis client not initialized")
	}

	return Client.Del(ctx, key).Err()
}

// DeleteMultiple removes multiple keys from cache
func DeleteMultiple(keys ...string) error {
	if Client == nil {
		return fmt.Errorf("Redis client not initialized")
	}

	if len(keys) == 0 {
		return nil
	}

	return Client.Del(ctx, keys...).Err()
}

// Exists checks if a key exists in cache
func Exists(key string) (bool, error) {
	if Client == nil {
		return false, fmt.Errorf("Redis client not initialized")
	}

	result, err := Client.Exists(ctx, key).Result()
	if err != nil {
		return false, err
	}

	return result > 0, nil
}

// Expire sets expiration time for a key
func Expire(key string, expiration time.Duration) error {
	if Client == nil {
		return fmt.Errorf("Redis client not initialized")
	}

	return Client.Expire(ctx, key, expiration).Err()
}

// TTL returns the remaining time to live of a key
func TTL(key string) (time.Duration, error) {
	if Client == nil {
		return 0, fmt.Errorf("Redis client not initialized")
	}

	return Client.TTL(ctx, key).Result()
}

// ============================================================================
// JSON Cache Operations
// ============================================================================

// GetJSON retrieves and unmarshals JSON data from cache
func GetJSON(key string, dest interface{}) error {
	if Client == nil {
		return fmt.Errorf("Redis client not initialized")
	}

	val, err := Client.Get(ctx, key).Result()
	if err == redis.Nil {
		return ErrNotFound
	}
	if err != nil {
		return fmt.Errorf("failed to get key %s: %w", key, err)
	}

	if err := json.Unmarshal([]byte(val), dest); err != nil {
		return fmt.Errorf("%w: %v", ErrUnmarshalJSON, err)
	}

	return nil
}

// SetJSON marshals and stores JSON data in cache
func SetJSON(key string, value interface{}, expiration time.Duration) error {
	if Client == nil {
		return fmt.Errorf("Redis client not initialized")
	}

	jsonData, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("%w: %v", ErrMarshalJSON, err)
	}

	return Client.Set(ctx, key, jsonData, expiration).Err()
}

// GetOrSetJSON gets value from cache, or sets it if not found
func GetOrSetJSON(key string, dest interface{}, expiration time.Duration, fetchFn func() (interface{}, error)) error {
	// Try to get from cache
	err := GetJSON(key, dest)
	if err == nil {
		return nil // Found in cache
	}
	if err != ErrNotFound {
		return err // Real error
	}

	// Not in cache, fetch data
	data, err := fetchFn()
	if err != nil {
		return fmt.Errorf("failed to fetch data: %w", err)
	}

	// Store in cache
	if err := SetJSON(key, data, expiration); err != nil {
		log.Printf("⚠️  Failed to cache data for key %s: %v", key, err)
		// Don't return error, just log it
	}

	// Marshal to dest
	jsonData, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("%w: %v", ErrMarshalJSON, err)
	}

	return json.Unmarshal(jsonData, dest)
}

// ============================================================================
// Pattern-based Operations
// ============================================================================

// DeletePattern deletes all keys matching a pattern
func DeletePattern(pattern string) error {
	if Client == nil {
		return fmt.Errorf("Redis client not initialized")
	}

	var cursor uint64
	var keys []string

	for {
		var batch []string
		var err error

		batch, cursor, err = Client.Scan(ctx, cursor, pattern, 100).Result()
		if err != nil {
			return fmt.Errorf("failed to scan keys: %w", err)
		}

		keys = append(keys, batch...)

		if cursor == 0 {
			break
		}
	}

	if len(keys) > 0 {
		return Client.Del(ctx, keys...).Err()
	}

	return nil
}

// GetKeysByPattern retrieves all keys matching a pattern
func GetKeysByPattern(pattern string) ([]string, error) {
	if Client == nil {
		return nil, fmt.Errorf("Redis client not initialized")
	}

	var cursor uint64
	var keys []string

	for {
		var batch []string
		var err error

		batch, cursor, err = Client.Scan(ctx, cursor, pattern, 100).Result()
		if err != nil {
			return nil, fmt.Errorf("failed to scan keys: %w", err)
		}

		keys = append(keys, batch...)

		if cursor == 0 {
			break
		}
	}

	return keys, nil
}

// ============================================================================
// Hash Operations
// ============================================================================

// HSet sets field in hash
func HSet(key, field string, value interface{}) error {
	if Client == nil {
		return fmt.Errorf("Redis client not initialized")
	}

	return Client.HSet(ctx, key, field, value).Err()
}

// HGet gets field from hash
func HGet(key, field string) (string, error) {
	if Client == nil {
		return "", fmt.Errorf("Redis client not initialized")
	}

	val, err := Client.HGet(ctx, key, field).Result()
	if err == redis.Nil {
		return "", ErrNotFound
	}
	return val, err
}

// HGetAll gets all fields from hash
func HGetAll(key string) (map[string]string, error) {
	if Client == nil {
		return nil, fmt.Errorf("Redis client not initialized")
	}

	return Client.HGetAll(ctx, key).Result()
}

// HDel deletes field from hash
func HDel(key string, fields ...string) error {
	if Client == nil {
		return fmt.Errorf("Redis client not initialized")
	}

	return Client.HDel(ctx, key, fields...).Err()
}

// ============================================================================
// List Operations
// ============================================================================

// LPush pushes value to head of list
func LPush(key string, values ...interface{}) error {
	if Client == nil {
		return fmt.Errorf("Redis client not initialized")
	}

	return Client.LPush(ctx, key, values...).Err()
}

// RPush pushes value to tail of list
func RPush(key string, values ...interface{}) error {
	if Client == nil {
		return fmt.Errorf("Redis client not initialized")
	}

	return Client.RPush(ctx, key, values...).Err()
}

// LPop pops value from head of list
func LPop(key string) (string, error) {
	if Client == nil {
		return "", fmt.Errorf("Redis client not initialized")
	}

	val, err := Client.LPop(ctx, key).Result()
	if err == redis.Nil {
		return "", ErrNotFound
	}
	return val, err
}

// RPop pops value from tail of list
func RPop(key string) (string, error) {
	if Client == nil {
		return "", fmt.Errorf("Redis client not initialized")
	}

	val, err := Client.RPop(ctx, key).Result()
	if err == redis.Nil {
		return "", ErrNotFound
	}
	return val, err
}

// LRange gets range of values from list
func LRange(key string, start, stop int64) ([]string, error) {
	if Client == nil {
		return nil, fmt.Errorf("Redis client not initialized")
	}

	return Client.LRange(ctx, key, start, stop).Result()
}

// LLen gets length of list
func LLen(key string) (int64, error) {
	if Client == nil {
		return 0, fmt.Errorf("Redis client not initialized")
	}

	return Client.LLen(ctx, key).Result()
}

// ============================================================================
// Set Operations
// ============================================================================

// SAdd adds members to set
func SAdd(key string, members ...interface{}) error {
	if Client == nil {
		return fmt.Errorf("Redis client not initialized")
	}

	return Client.SAdd(ctx, key, members...).Err()
}

// SMembers gets all members of set
func SMembers(key string) ([]string, error) {
	if Client == nil {
		return nil, fmt.Errorf("Redis client not initialized")
	}

	return Client.SMembers(ctx, key).Result()
}

// SIsMember checks if member is in set
func SIsMember(key string, member interface{}) (bool, error) {
	if Client == nil {
		return false, fmt.Errorf("Redis client not initialized")
	}

	return Client.SIsMember(ctx, key, member).Result()
}

// SRem removes members from set
func SRem(key string, members ...interface{}) error {
	if Client == nil {
		return fmt.Errorf("Redis client not initialized")
	}

	return Client.SRem(ctx, key, members...).Err()
}

// ============================================================================
// Counter Operations
// ============================================================================

// Incr increments counter
func Incr(key string) (int64, error) {
	if Client == nil {
		return 0, fmt.Errorf("Redis client not initialized")
	}

	return Client.Incr(ctx, key).Result()
}

// IncrBy increments counter by value
func IncrBy(key string, value int64) (int64, error) {
	if Client == nil {
		return 0, fmt.Errorf("Redis client not initialized")
	}

	return Client.IncrBy(ctx, key, value).Result()
}

// Decr decrements counter
func Decr(key string) (int64, error) {
	if Client == nil {
		return 0, fmt.Errorf("Redis client not initialized")
	}

	return Client.Decr(ctx, key).Result()
}

// DecrBy decrements counter by value
func DecrBy(key string, value int64) (int64, error) {
	if Client == nil {
		return 0, fmt.Errorf("Redis client not initialized")
	}

	return Client.DecrBy(ctx, key, value).Result()
}

// ============================================================================
// Helper Functions
// ============================================================================

// FlushDB clears all keys in current database (use with caution!)
func FlushDB() error {
	if Client == nil {
		return fmt.Errorf("Redis client not initialized")
	}

	return Client.FlushDB(ctx).Err()
}

// FlushAll clears all keys in all databases (use with extreme caution!)
func FlushAll() error {
	if Client == nil {
		return fmt.Errorf("Redis client not initialized")
	}

	return Client.FlushAll(ctx).Err()
}

// Info returns Redis server information
func Info(section ...string) (string, error) {
	if Client == nil {
		return "", fmt.Errorf("Redis client not initialized")
	}

	return Client.Info(ctx, section...).Result()
}

// DBSize returns number of keys in current database
func DBSize() (int64, error) {
	if Client == nil {
		return 0, fmt.Errorf("Redis client not initialized")
	}

	return Client.DBSize(ctx).Result()
}

// ============================================================================
// Backward Compatibility (deprecated but kept for compatibility)
// ============================================================================

// CacheGet is deprecated, use Get instead
func CacheGet(key string) (string, error) {
	log.Println("⚠️  CacheGet is deprecated, use Get instead")
	return Get(key)
}

// CacheSet is deprecated, use SetJSON instead
func CacheSet(key string, value interface{}, expiration time.Duration) error {
	log.Println("⚠️  CacheSet is deprecated, use SetJSON instead")
	return SetJSON(key, value, expiration)
}

// CacheDelete is deprecated, use Delete instead
func CacheDelete(key string) error {
	log.Println("⚠️  CacheDelete is deprecated, use Delete instead")
	return Delete(key)
}
