package config

import (
	"fmt"
	"log"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/joho/godotenv"
)

// Config holds all application configuration
type Config struct {
	App      AppConfig
	Database DatabaseConfig
	Redis    RedisConfig
	JWT      JWTConfig
	Security SecurityConfig
	External ExternalConfig
	Email    EmailConfig
	Storage  StorageConfig
}

// AppConfig contains application-level configuration
type AppConfig struct {
	Env         string
	Port        string
	Name        string
	Version     string
	Debug       bool
	BaseURL     string
	FrontendURL string
}

// DatabaseConfig contains database connection configuration
type DatabaseConfig struct {
	URL             string
	Host            string
	Port            string
	User            string
	Password        string
	Name            string
	SSLMode         string
	MaxOpenConns    int
	MaxIdleConns    int
	ConnMaxLifetime time.Duration
	ConnMaxIdleTime time.Duration
}

// RedisConfig contains Redis connection configuration
type RedisConfig struct {
	URL      string
	Host     string
	Port     string
	Password string
	DB       int
	PoolSize int
}

// JWTConfig contains JWT token configuration
type JWTConfig struct {
	Secret               []byte
	AccessTokenDuration  time.Duration
	RefreshTokenDuration time.Duration
	ResetTokenDuration   time.Duration
	VerifyTokenDuration  time.Duration
	Issuer               string
	Audience             string
}

// SecurityConfig contains security-related configuration
type SecurityConfig struct {
	BcryptCost           int
	MaxLoginAttempts     int
	LockoutDuration      time.Duration
	SessionTimeout       time.Duration
	PasswordMinLength    int
	RequireSpecialChar   bool
	RequireNumber        bool
	RequireUppercase     bool
	PasswordHistoryCount int
	TwoFactorEnabled     bool
}

// ExternalConfig contains external API configuration
type ExternalConfig struct {
	OpenWeatherAPIKey  string
	OpenWeatherBaseURL string
	// Add other external APIs here
}

// EmailConfig contains email service configuration
type EmailConfig struct {
	SMTPHost     string
	SMTPPort     int
	SMTPUser     string
	SMTPPassword string
	FromEmail    string
	FromName     string
	UseTLS       bool
}

// StorageConfig contains file storage configuration
type StorageConfig struct {
	Type             string // "local", "s3", "gcs"
	LocalPath        string
	S3Bucket         string
	S3Region         string
	S3AccessKey      string
	S3SecretKey      string
	MaxUploadSize    int64
	AllowedFileTypes []string
}

// Global config instance
var Cfg *Config

// Load loads configuration from environment variables
func Load() error {
	// Load .env file in non-production environments
	if env := os.Getenv("APP_ENV"); env != "production" && env != "staging" {
		if err := godotenv.Load(); err != nil {
			log.Println("⚠️  No .env file found, using system environment variables")
		}
	}

	cfg := &Config{}

	// Load all configurations
	if err := loadAppConfig(&cfg.App); err != nil {
		return fmt.Errorf("failed to load app config: %w", err)
	}

	if err := loadDatabaseConfig(&cfg.Database); err != nil {
		return fmt.Errorf("failed to load database config: %w", err)
	}

	if err := loadRedisConfig(&cfg.Redis); err != nil {
		return fmt.Errorf("failed to load redis config: %w", err)
	}

	if err := loadJWTConfig(&cfg.JWT); err != nil {
		return fmt.Errorf("failed to load JWT config: %w", err)
	}

	if err := loadSecurityConfig(&cfg.Security); err != nil {
		return fmt.Errorf("failed to load security config: %w", err)
	}

	if err := loadExternalConfig(&cfg.External); err != nil {
		return fmt.Errorf("failed to load external config: %w", err)
	}

	if err := loadEmailConfig(&cfg.Email); err != nil {
		return fmt.Errorf("failed to load email config: %w", err)
	}

	if err := loadStorageConfig(&cfg.Storage); err != nil {
		return fmt.Errorf("failed to load storage config: %w", err)
	}

	// Validate configuration
	if err := cfg.Validate(); err != nil {
		return fmt.Errorf("configuration validation failed: %w", err)
	}

	Cfg = cfg

	// Log loaded configuration (without sensitive data)
	logConfiguration()

	return nil
}

// MustLoad loads configuration and panics on error
func MustLoad() {
	if err := Load(); err != nil {
		log.Fatalf("❌ Failed to load configuration: %v", err)
	}
}

func loadAppConfig(cfg *AppConfig) error {
	cfg.Env = getEnvOrDefault("APP_ENV", "development")
	cfg.Port = getEnvOrDefault("PORT", "8080")
	cfg.Name = getEnvOrDefault("APP_NAME", "SurvivalKitAPI")
	cfg.Version = getEnvOrDefault("APP_VERSION", "1.0.0")
	cfg.Debug = getBoolEnv("APP_DEBUG", cfg.Env == "development")
	cfg.BaseURL = getEnvOrDefault("APP_BASE_URL", "http://localhost:"+cfg.Port)
	cfg.FrontendURL = getEnvOrDefault("FRONTEND_URL", "http://localhost:3000")

	return nil
}

func loadDatabaseConfig(cfg *DatabaseConfig) error {
	// Option 1: Use DATABASE_URL directly (Railway, Heroku, etc.)
	if dbURL := os.Getenv("DATABASE_URL"); dbURL != "" {
		log.Println("✅ Using DATABASE_URL from environment")

		// Parse DATABASE_URL to extract components
		parsed, err := url.Parse(dbURL)
		if err != nil {
			return fmt.Errorf("invalid DATABASE_URL: %w", err)
		}

		cfg.URL = dbURL
		cfg.Host = parsed.Hostname()
		cfg.Port = parsed.Port()
		if cfg.Port == "" {
			cfg.Port = "5432"
		}
		cfg.User = parsed.User.Username()
		cfg.Password, _ = parsed.User.Password()
		cfg.Name = strings.TrimPrefix(parsed.Path, "/")

		// Extract sslmode from query params
		query := parsed.Query()
		cfg.SSLMode = query.Get("sslmode")
		if cfg.SSLMode == "" {
			cfg.SSLMode = "require"
		}
	} else {
		// Option 2: Build from individual variables
		log.Println("⚠️  Building DATABASE_URL from individual variables")

		requiredVars := map[string]string{
			"DB_HOST":     "Database host",
			"DB_USER":     "Database user",
			"DB_PASSWORD": "Database password",
			"DB_NAME":     "Database name",
		}

		for key, desc := range requiredVars {
			if os.Getenv(key) == "" {
				return fmt.Errorf("missing required environment variable: %s (%s)", key, desc)
			}
		}

		cfg.Host = os.Getenv("DB_HOST")
		cfg.Port = getEnvOrDefault("DB_PORT", "5432")
		cfg.User = os.Getenv("DB_USER")
		cfg.Password = os.Getenv("DB_PASSWORD")
		cfg.Name = os.Getenv("DB_NAME")
		cfg.SSLMode = getEnvOrDefault("DB_SSLMODE", "disable")

		cfg.URL = fmt.Sprintf(
			"postgres://%s:%s@%s:%s/%s?sslmode=%s",
			cfg.User,
			url.QueryEscape(cfg.Password),
			cfg.Host,
			cfg.Port,
			cfg.Name,
			cfg.SSLMode,
		)
	}

	// Connection pool settings
	cfg.MaxOpenConns = getIntEnv("DB_MAX_OPEN_CONNS", 25)
	cfg.MaxIdleConns = getIntEnv("DB_MAX_IDLE_CONNS", 5)
	cfg.ConnMaxLifetime = getDurationEnv("DB_CONN_MAX_LIFETIME", 5*time.Minute)
	cfg.ConnMaxIdleTime = getDurationEnv("DB_CONN_MAX_IDLE_TIME", 5*time.Minute)

	return nil
}

func loadRedisConfig(cfg *RedisConfig) error {
	// Option 1: Use REDIS_URL
	if redisURL := os.Getenv("REDIS_URL"); redisURL != "" {
		log.Println("✅ Using REDIS_URL from environment")

		parsed, err := url.Parse(redisURL)
		if err != nil {
			return fmt.Errorf("invalid REDIS_URL: %w", err)
		}

		cfg.URL = redisURL
		cfg.Host = parsed.Hostname()
		cfg.Port = parsed.Port()
		if cfg.Port == "" {
			cfg.Port = "6379"
		}

		if parsed.User != nil {
			cfg.Password, _ = parsed.User.Password()
		}

		// Extract DB from path (e.g., /0, /1)
		if len(parsed.Path) > 1 {
			if db, err := strconv.Atoi(strings.TrimPrefix(parsed.Path, "/")); err == nil {
				cfg.DB = db
			}
		}
	} else {
		// Option 2: Build from individual variables
		cfg.Host = getEnvOrDefault("REDIS_HOST", "localhost")
		cfg.Port = getEnvOrDefault("REDIS_PORT", "6379")
		cfg.Password = os.Getenv("REDIS_PASSWORD")
		cfg.DB = getIntEnv("REDIS_DB", 0)

		cfg.URL = fmt.Sprintf("redis://:%s@%s:%s/%d",
			cfg.Password,
			cfg.Host,
			cfg.Port,
			cfg.DB,
		)
	}

	cfg.PoolSize = getIntEnv("REDIS_POOL_SIZE", 10)

	return nil
}

func loadJWTConfig(cfg *JWTConfig) error {
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		return fmt.Errorf("JWT_SECRET is required")
	}

	cfg.Secret = []byte(jwtSecret)
	cfg.AccessTokenDuration = getDurationEnv("JWT_ACCESS_DURATION", 15*time.Minute)
	cfg.RefreshTokenDuration = getDurationEnv("JWT_REFRESH_DURATION", 7*24*time.Hour)
	cfg.ResetTokenDuration = getDurationEnv("JWT_RESET_DURATION", 1*time.Hour)
	cfg.VerifyTokenDuration = getDurationEnv("JWT_VERIFY_DURATION", 24*time.Hour)
	cfg.Issuer = getEnvOrDefault("JWT_ISSUER", "healthcare-api")
	cfg.Audience = getEnvOrDefault("JWT_AUDIENCE", "healthcare-app")

	return nil
}

func loadSecurityConfig(cfg *SecurityConfig) error {
	cfg.BcryptCost = getIntEnv("BCRYPT_COST", 12)
	cfg.MaxLoginAttempts = getIntEnv("MAX_LOGIN_ATTEMPTS", 5)
	cfg.LockoutDuration = getDurationEnv("LOCKOUT_DURATION", 30*time.Minute)
	cfg.SessionTimeout = getDurationEnv("SESSION_TIMEOUT", 24*time.Hour)
	cfg.PasswordMinLength = getIntEnv("PASSWORD_MIN_LENGTH", 8)
	cfg.RequireSpecialChar = getBoolEnv("PASSWORD_REQUIRE_SPECIAL", true)
	cfg.RequireNumber = getBoolEnv("PASSWORD_REQUIRE_NUMBER", true)
	cfg.RequireUppercase = getBoolEnv("PASSWORD_REQUIRE_UPPERCASE", true)
	cfg.PasswordHistoryCount = getIntEnv("PASSWORD_HISTORY_COUNT", 5)
	cfg.TwoFactorEnabled = getBoolEnv("TWO_FACTOR_ENABLED", false)

	return nil
}

func loadExternalConfig(cfg *ExternalConfig) error {
	cfg.OpenWeatherAPIKey = os.Getenv("OPENWEATHER_API_KEY")
	cfg.OpenWeatherBaseURL = getEnvOrDefault("OPENWEATHER_BASE_URL", "https://api.openweathermap.org/data/2.5")

	// Warn if API key is missing but don't fail
	if cfg.OpenWeatherAPIKey == "" {
		log.Println("⚠️  OPENWEATHER_API_KEY not set, weather features will be disabled")
	}

	return nil
}

func loadEmailConfig(cfg *EmailConfig) error {
	cfg.SMTPHost = os.Getenv("SMTP_HOST")
	cfg.SMTPPort = getIntEnv("SMTP_PORT", 587)
	cfg.SMTPUser = os.Getenv("SMTP_USER")
	cfg.SMTPPassword = os.Getenv("SMTP_PASSWORD")
	cfg.FromEmail = os.Getenv("SMTP_FROM_EMAIL")
	cfg.FromName = getEnvOrDefault("SMTP_FROM_NAME", "Healthcare App")
	cfg.UseTLS = getBoolEnv("SMTP_USE_TLS", true)

	// Email is optional, just warn if not configured
	if cfg.SMTPHost == "" {
		log.Println("⚠️  SMTP not configured, email features will be disabled")
	}

	return nil
}

func loadStorageConfig(cfg *StorageConfig) error {
	cfg.Type = getEnvOrDefault("STORAGE_TYPE", "local")
	cfg.LocalPath = getEnvOrDefault("STORAGE_LOCAL_PATH", "./uploads")
	cfg.S3Bucket = os.Getenv("S3_BUCKET")
	cfg.S3Region = os.Getenv("S3_REGION")
	cfg.S3AccessKey = os.Getenv("S3_ACCESS_KEY")
	cfg.S3SecretKey = os.Getenv("S3_SECRET_KEY")
	cfg.MaxUploadSize = int64(getIntEnv("MAX_UPLOAD_SIZE_MB", 10)) * 1024 * 1024 // Convert to bytes

	allowedTypes := os.Getenv("ALLOWED_FILE_TYPES")
	if allowedTypes == "" {
		cfg.AllowedFileTypes = []string{".jpg", ".jpeg", ".png", ".gif", ".pdf", ".doc", ".docx"}
	} else {
		cfg.AllowedFileTypes = strings.Split(allowedTypes, ",")
	}

	return nil
}

// Validate validates the entire configuration
func (c *Config) Validate() error {
	// Validate App
	if c.App.Port == "" {
		return fmt.Errorf("app port is required")
	}

	// Validate Database
	if c.Database.URL == "" {
		return fmt.Errorf("database URL is required")
	}
	if c.Database.MaxOpenConns < c.Database.MaxIdleConns {
		return fmt.Errorf("max open connections must be >= max idle connections")
	}

	// Validate Redis
	if c.Redis.URL == "" {
		return fmt.Errorf("redis URL is required")
	}

	// Validate JWT
	if len(c.JWT.Secret) < 32 {
		return fmt.Errorf("JWT secret must be at least 32 characters")
	}
	if c.JWT.AccessTokenDuration <= 0 {
		return fmt.Errorf("JWT access token duration must be positive")
	}
	if c.JWT.RefreshTokenDuration <= c.JWT.AccessTokenDuration {
		return fmt.Errorf("JWT refresh token duration must be greater than access token duration")
	}

	// Validate Security
	if c.Security.BcryptCost < 10 || c.Security.BcryptCost > 14 {
		log.Printf("⚠️  Warning: Bcrypt cost %d is outside recommended range (10-14)", c.Security.BcryptCost)
	}
	if c.Security.PasswordMinLength < 8 {
		return fmt.Errorf("password minimum length must be at least 8")
	}

	// Validate Storage
	if c.Storage.Type == "s3" {
		if c.Storage.S3Bucket == "" || c.Storage.S3Region == "" {
			return fmt.Errorf("S3 bucket and region are required when using S3 storage")
		}
	}

	return nil
}

// IsProduction returns true if running in production environment
func (c *Config) IsProduction() bool {
	return c.App.Env == "production"
}

// IsDevelopment returns true if running in development environment
func (c *Config) IsDevelopment() bool {
	return c.App.Env == "development"
}

// IsStaging returns true if running in staging environment
func (c *Config) IsStaging() bool {
	return c.App.Env == "staging"
}

// Helper functions

func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getIntEnv(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
		log.Printf("⚠️  Invalid integer value for %s: %s, using default: %d", key, value, defaultValue)
	}
	return defaultValue
}

func getBoolEnv(key string, defaultValue bool) bool {
	if value := os.Getenv(key); value != "" {
		if boolValue, err := strconv.ParseBool(value); err == nil {
			return boolValue
		}
		log.Printf("⚠️  Invalid boolean value for %s: %s, using default: %t", key, value, defaultValue)
	}
	return defaultValue
}

func getDurationEnv(key string, defaultValue time.Duration) time.Duration {
	if value := os.Getenv(key); value != "" {
		if duration, err := time.ParseDuration(value); err == nil {
			return duration
		}
		log.Printf("⚠️  Invalid duration value for %s: %s, using default: %s", key, value, defaultValue)
	}
	return defaultValue
}

func logConfiguration() {
	log.Println("=" + strings.Repeat("=", 70))
	log.Println("✅ Configuration loaded successfully")
	log.Println("=" + strings.Repeat("=", 70))

	log.Printf("📱 Application:")
	log.Printf("   Environment: %s", Cfg.App.Env)
	log.Printf("   Name: %s v%s", Cfg.App.Name, Cfg.App.Version)
	log.Printf("   Port: %s", Cfg.App.Port)
	log.Printf("   Debug: %t", Cfg.App.Debug)
	log.Printf("   Base URL: %s", Cfg.App.BaseURL)
	log.Printf("   Frontend URL: %s", Cfg.App.FrontendURL)

	log.Printf("🗄️  Database:")
	log.Printf("   Host: %s:%s", Cfg.Database.Host, Cfg.Database.Port)
	log.Printf("   Name: %s", Cfg.Database.Name)
	log.Printf("   SSL Mode: %s", Cfg.Database.SSLMode)
	log.Printf("   Max Open Connections: %d", Cfg.Database.MaxOpenConns)
	log.Printf("   Max Idle Connections: %d", Cfg.Database.MaxIdleConns)

	log.Printf("🔴 Redis:")
	log.Printf("   Host: %s:%s", Cfg.Redis.Host, Cfg.Redis.Port)
	log.Printf("   DB: %d", Cfg.Redis.DB)
	log.Printf("   Pool Size: %d", Cfg.Redis.PoolSize)

	log.Printf("🔐 JWT:")
	log.Printf("   Access Token Duration: %s", Cfg.JWT.AccessTokenDuration)
	log.Printf("   Refresh Token Duration: %s", Cfg.JWT.RefreshTokenDuration)
	log.Printf("   Issuer: %s", Cfg.JWT.Issuer)

	log.Printf("🛡️  Security:")
	log.Printf("   Bcrypt Cost: %d", Cfg.Security.BcryptCost)
	log.Printf("   Max Login Attempts: %d", Cfg.Security.MaxLoginAttempts)
	log.Printf("   Lockout Duration: %s", Cfg.Security.LockoutDuration)
	log.Printf("   Two-Factor: %t", Cfg.Security.TwoFactorEnabled)

	log.Printf("📧 Email:")
	if Cfg.Email.SMTPHost != "" {
		log.Printf("   SMTP Host: %s:%d", Cfg.Email.SMTPHost, Cfg.Email.SMTPPort)
		log.Printf("   From: %s <%s>", Cfg.Email.FromName, Cfg.Email.FromEmail)
	} else {
		log.Printf("   Status: Not configured")
	}

	log.Printf("💾 Storage:")
	log.Printf("   Type: %s", Cfg.Storage.Type)
	if Cfg.Storage.Type == "local" {
		log.Printf("   Path: %s", Cfg.Storage.LocalPath)
	} else if Cfg.Storage.Type == "s3" {
		log.Printf("   Bucket: %s", Cfg.Storage.S3Bucket)
		log.Printf("   Region: %s", Cfg.Storage.S3Region)
	}
	log.Printf("   Max Upload Size: %d MB", Cfg.Storage.MaxUploadSize/(1024*1024))

	log.Printf("🌐 External APIs:")
	if Cfg.External.OpenWeatherAPIKey != "" {
		log.Printf("   OpenWeather: Configured")
	} else {
		log.Printf("   OpenWeather: Not configured")
	}

	log.Println("=" + strings.Repeat("=", 70))
}
