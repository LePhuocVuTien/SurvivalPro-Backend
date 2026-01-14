package ratelimit

import (
	"errors"
	"fmt"
	"time"
)

// ============================================================================
// STORE TYPE (Type-safe enum)
// ============================================================================

// StoreType represents the storage backend for rate limiting
type StoreType string

const (
	// StoreMemory uses in-memory storage (fast, good for single server)
	StoreMemory StoreType = "memory"

	// StoreRedis uses Redis storage (distributed, good for multiple servers)
	StoreRedis StoreType = "redis"

	// StoreDatabase uses database storage (persistent, audit trail)
	StoreDatabase StoreType = "database"
)

// IsValid checks if store type is valid
func (s StoreType) IsValid() bool {
	switch s {
	case StoreMemory, StoreRedis, StoreDatabase:
		return true
	default:
		return false
	}
}

// String returns string representation
func (s StoreType) String() string {
	return string(s)
}

// ============================================================================
// CONFIGURATION
// ============================================================================

// Config represents rate limiter configuration
type Config struct {
	// Enabled controls whether rate limiting is active
	// Set to false in development to disable rate limiting
	Enabled bool

	// DefaultRules contains the default rate limit rules to apply
	DefaultRules []*RateLimitRule

	// StoreType determines which storage backend to use
	// Options: StoreMemory, StoreRedis, StoreDatabase
	StoreType StoreType

	// Redis configuration (only used if StoreType is StoreRedis)
	RedisAddr     string // Redis server address (e.g., "localhost:6379")
	RedisPassword string // Redis password (empty if no auth)
	RedisDB       int    // Redis database number (0-15)

	// CleanupInterval for memory store (how often to clean expired entries)
	// Only used with StoreMemory. Default is 5 minutes.
	CleanupInterval time.Duration
}

// Validate validates the configuration
func (c *Config) Validate() error {
	// Check store type
	if !c.StoreType.IsValid() {
		return fmt.Errorf("invalid store type: %s", c.StoreType)
	}

	// Validate Redis config if using Redis
	if c.StoreType == StoreRedis {
		if c.RedisAddr == "" {
			return errors.New("redis address is required when using Redis store")
		}
	}

	// Validate rules
	for _, rule := range c.DefaultRules {
		if rule.Action == "" {
			return errors.New("rate limit rule must have an action")
		}
		if rule.MaxAttempts <= 0 {
			return fmt.Errorf("rule %s: max_attempts must be positive", rule.Action)
		}
		if rule.WindowSizeSeconds <= 0 {
			return fmt.Errorf("rule %s: window_size must be positive", rule.Action)
		}
		if rule.BlockDurationSeconds < 0 {
			return fmt.Errorf("rule %s: block_duration cannot be negative", rule.Action)
		}
	}

	return nil
}

// ============================================================================
// DEFAULT CONFIGURATIONS
// ============================================================================

// DefaultConfig returns default rate limiter configuration
// This is suitable for most applications with balanced rate limiting
func DefaultConfig() *Config {
	return &Config{
		Enabled:         true,
		DefaultRules:    DefaultRules(),
		StoreType:       StoreMemory,
		CleanupInterval: 5 * time.Minute,
	}
}

// DefaultRules returns default rate limit rules
// These rules provide reasonable protection against common attacks
func DefaultRules() []*RateLimitRule {
	rules := []*RateLimitRule{
		{
			Action:               ActionLogin,
			MaxAttempts:          5,
			WindowSizeSeconds:    300,  // 5 minutes
			BlockDurationSeconds: 1800, // 30 minutes
			IsActive:             true,
		},
		{
			Action:               ActionPasswordReset,
			MaxAttempts:          3,
			WindowSizeSeconds:    3600, // 1 hour
			BlockDurationSeconds: 3600, // 1 hour
			IsActive:             true,
		},
		{
			Action:               ActionEmailVerify,
			MaxAttempts:          5,
			WindowSizeSeconds:    3600, // 1 hour
			BlockDurationSeconds: 3600, // 1 hour
			IsActive:             true,
		},
		{
			Action:               ActionResendEmail,
			MaxAttempts:          3,
			WindowSizeSeconds:    600,  // 10 minutes
			BlockDurationSeconds: 1800, // 30 minutes
			IsActive:             true,
		},
		{
			Action:               ActionRegistration,
			MaxAttempts:          3,
			WindowSizeSeconds:    3600, // 1 hour
			BlockDurationSeconds: 7200, // 2 hours
			IsActive:             true,
		},
		{
			Action:               ActionOTPRequest,
			MaxAttempts:          5,
			WindowSizeSeconds:    300, // 5 minutes
			BlockDurationSeconds: 600, // 10 minutes
			IsActive:             true,
		},
	}

	// LoadDurations converts second-based config to time.Duration fields
	// This is done for runtime convenience - seconds are the source of truth
	for _, rule := range rules {
		rule.LoadDurations()
	}

	return rules
}

// ============================================================================
// PRESET CONFIGURATIONS
// ============================================================================

// StrictConfig returns strict rate limiting configuration
// Use this for high-security environments or when under attack
//
// Characteristics:
// - Lower attempt limits (3 login attempts vs 5)
// - Longer block durations (1 hour vs 30 minutes)
// - Recommended to use Redis for distributed systems
func StrictConfig() *Config {
	return &Config{
		Enabled: true,
		DefaultRules: []*RateLimitRule{
			NewRule(ActionLogin).
				MaxAttempts(3).
				Window(5 * time.Minute).
				BlockFor(1 * time.Hour).
				Build(),
			NewRule(ActionPasswordReset).
				MaxAttempts(2).
				Window(1 * time.Hour).
				BlockFor(2 * time.Hour).
				Build(),
			NewRule(ActionRegistration).
				MaxAttempts(2).
				Window(1 * time.Hour).
				BlockFor(4 * time.Hour).
				Build(),
			NewRule(ActionEmailVerify).
				MaxAttempts(3).
				Window(1 * time.Hour).
				BlockFor(2 * time.Hour).
				Build(),
			NewRule(ActionOTPRequest).
				MaxAttempts(3).
				Window(5 * time.Minute).
				BlockFor(15 * time.Minute).
				Build(),
		},
		StoreType:       StoreRedis,
		CleanupInterval: 5 * time.Minute,
	}
}

// LenientConfig returns lenient rate limiting configuration
// Use this for development or trusted environments
//
// Characteristics:
// - Higher attempt limits (10 login attempts vs 5)
// - Shorter block durations (15 minutes vs 30 minutes)
// - Uses memory store for simplicity
func LenientConfig() *Config {
	return &Config{
		Enabled: true,
		DefaultRules: []*RateLimitRule{
			NewRule(ActionLogin).
				MaxAttempts(10).
				Window(5 * time.Minute).
				BlockFor(15 * time.Minute).
				Build(),
			NewRule(ActionPasswordReset).
				MaxAttempts(5).
				Window(1 * time.Hour).
				BlockFor(30 * time.Minute).
				Build(),
			NewRule(ActionRegistration).
				MaxAttempts(5).
				Window(1 * time.Hour).
				BlockFor(1 * time.Hour).
				Build(),
			NewRule(ActionEmailVerify).
				MaxAttempts(10).
				Window(1 * time.Hour).
				BlockFor(30 * time.Minute).
				Build(),
			NewRule(ActionOTPRequest).
				MaxAttempts(10).
				Window(5 * time.Minute).
				BlockFor(5 * time.Minute).
				Build(),
		},
		StoreType:       StoreMemory,
		CleanupInterval: 5 * time.Minute,
	}
}

// DevelopmentConfig returns development-friendly configuration
// Rate limiting is DISABLED to make development easier
//
// Use this only in local development - never in production!
func DevelopmentConfig() *Config {
	return &Config{
		Enabled:         false, // Disabled for development
		DefaultRules:    DefaultRules(),
		StoreType:       StoreMemory,
		CleanupInterval: 5 * time.Minute,
	}
}

// ProductionConfig returns production-ready configuration
// Uses Redis for distributed rate limiting across multiple servers
//
// You must provide Redis connection details:
//
//	cfg := ProductionConfig()
//	cfg.RedisAddr = "redis:6379"
//	cfg.RedisPassword = "secret"
func ProductionConfig() *Config {
	return &Config{
		Enabled:         true,
		DefaultRules:    DefaultRules(),
		StoreType:       StoreRedis,
		RedisAddr:       "", // Must be set by caller
		RedisPassword:   "", // Must be set by caller
		RedisDB:         0,
		CleanupInterval: 5 * time.Minute,
	}
}

// ============================================================================
// RULE BUILDER (Fluent API)
// ============================================================================

// RuleBuilder helps build custom rate limit rules with a fluent API
//
// Example usage:
//
//	rule := NewRule(ActionLogin).
//	    MaxAttempts(5).
//	    Window(5 * time.Minute).
//	    BlockFor(30 * time.Minute).
//	    Build()
type RuleBuilder struct {
	rule *RateLimitRule
}

// NewRule creates a new rule builder for the specified action
func NewRule(action string) *RuleBuilder {
	return &RuleBuilder{
		rule: &RateLimitRule{
			Action:   action,
			IsActive: true,
		},
	}
}

// MaxAttempts sets the maximum number of attempts allowed in the window
func (b *RuleBuilder) MaxAttempts(n int) *RuleBuilder {
	b.rule.MaxAttempts = n
	return b
}

// Window sets the time window for counting attempts
// Example: Window(5 * time.Minute) allows MaxAttempts within 5 minutes
func (b *RuleBuilder) Window(d time.Duration) *RuleBuilder {
	b.rule.WindowSize = d
	b.rule.WindowSizeSeconds = int(d.Seconds())
	return b
}

// BlockFor sets how long to block after exceeding the limit
// Example: BlockFor(30 * time.Minute) blocks for 30 minutes
func (b *RuleBuilder) BlockFor(d time.Duration) *RuleBuilder {
	b.rule.BlockDuration = d
	b.rule.BlockDurationSeconds = int(d.Seconds())
	return b
}

// Disabled marks the rule as inactive
// Inactive rules are ignored by the rate limiter
func (b *RuleBuilder) Disabled() *RuleBuilder {
	b.rule.IsActive = false
	return b
}

// Enabled marks the rule as active (default state)
func (b *RuleBuilder) Enabled() *RuleBuilder {
	b.rule.IsActive = true
	return b
}

// Build returns the constructed rule
// It automatically calls LoadDurations() to ensure Duration fields are set
func (b *RuleBuilder) Build() *RateLimitRule {
	// Ensure Duration fields are populated from Seconds fields
	// This is important because:
	// 1. Seconds are the source of truth (stored in DB)
	// 2. Durations are runtime helpers (not stored)
	// 3. If user only set Seconds, Duration fields would be zero without this
	b.rule.LoadDurations()
	return b.rule
}

// ============================================================================
// CONFIGURATION HELPERS
// ============================================================================

// MergeRules merges multiple rule sets, with later rules overriding earlier ones
// This is useful for combining default rules with custom overrides
func MergeRules(ruleSets ...[]*RateLimitRule) []*RateLimitRule {
	ruleMap := make(map[string]*RateLimitRule)

	// Process rules in order - later rules override earlier ones
	for _, rules := range ruleSets {
		for _, rule := range rules {
			ruleMap[rule.Action] = rule
		}
	}

	// Convert map back to slice
	merged := make([]*RateLimitRule, 0, len(ruleMap))
	for _, rule := range ruleMap {
		merged = append(merged, rule)
	}

	return merged
}

// DisableRule creates a disabled version of a rule
// Useful for selectively disabling specific actions
func DisableRule(action string) *RateLimitRule {
	return &RateLimitRule{
		Action:   action,
		IsActive: false,
	}
}

// ============================================================================
// ENVIRONMENT-BASED CONFIGURATION
// ============================================================================

// ConfigFromEnv returns configuration based on environment
// Supports: "development", "staging", "production"
func ConfigFromEnv(env string) *Config {
	switch env {
	case "development", "dev":
		return DevelopmentConfig()
	case "staging":
		return LenientConfig()
	case "production", "prod":
		return ProductionConfig()
	default:
		return DefaultConfig()
	}
}

// ============================================================================
// EXAMPLE CUSTOM CONFIGURATIONS
// ============================================================================

// APIConfig returns configuration suitable for public APIs
// More lenient for general API calls, stricter for auth endpoints
func APIConfig() *Config {
	return &Config{
		Enabled: true,
		DefaultRules: []*RateLimitRule{
			// Authentication endpoints - strict
			NewRule(ActionLogin).
				MaxAttempts(5).
				Window(5 * time.Minute).
				BlockFor(30 * time.Minute).
				Build(),

			// General API calls - lenient
			NewRule(ActionAPICall).
				MaxAttempts(100).
				Window(1 * time.Minute).
				BlockFor(1 * time.Minute).
				Build(),

			// Registration - moderate
			NewRule(ActionRegistration).
				MaxAttempts(3).
				Window(1 * time.Hour).
				BlockFor(2 * time.Hour).
				Build(),
		},
		StoreType:       StoreRedis,
		CleanupInterval: 5 * time.Minute,
	}
}

// WebAppConfig returns configuration suitable for web applications
// Balanced between security and user experience
func WebAppConfig() *Config {
	return &Config{
		Enabled: true,
		DefaultRules: []*RateLimitRule{
			NewRule(ActionLogin).
				MaxAttempts(5).
				Window(10 * time.Minute).
				BlockFor(30 * time.Minute).
				Build(),
			NewRule(ActionPasswordReset).
				MaxAttempts(3).
				Window(1 * time.Hour).
				BlockFor(1 * time.Hour).
				Build(),
			NewRule(ActionRegistration).
				MaxAttempts(5).
				Window(24 * time.Hour).
				BlockFor(24 * time.Hour).
				Build(),
		},
		StoreType:       StoreMemory,
		CleanupInterval: 5 * time.Minute,
	}
}
