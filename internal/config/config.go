package config

import (
	"fmt"
	"log"
	"net/url"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	AppEnv         string
	AppPort        string
	DatabaseURL    string
	RedisURL       string
	JWTSecret      []byte
	OpenWeatherKey string
}

type DBConfig struct {
	User     string
	Password string
	Host     string
	Port     string
	Name     string
	SSLMode  string
}

var Cfg *Config

func Load() {

	if os.Getenv("APP_ENV") != "production" {
		if err := godotenv.Load(); err != nil {
			log.Println("⚠️ No .env file found, using system env")
		}
	}

	var databaseURL string

	// Option 1: Dùng DATABASE_URL trực tiếp (Railway production)
	if os.Getenv("APP_ENV") == "production" {
		dbURL := os.Getenv("DATABASE_URL")
		log.Println("✅ Using DATABASE_URL from environment", dbURL)
		databaseURL = dbURL
	} else {
		// Option 2: Build từ các biến riêng lẻ (local dev)
		log.Println("⚠️ Building DATABASE_URL from individual variables")

		requiredVars := []string{
			"DB_HOST", "DB_USER", "DB_PASSWORD", "DB_NAME",
		}

		for _, v := range requiredVars {
			if os.Getenv(v) == "" {
				log.Fatalf("❌ Missing required environment variable: %s", v)
			}
		}

		db := DBConfig{
			User:     os.Getenv("DB_USER"),
			Password: os.Getenv("DB_PASSWORD"),
			Host:     os.Getenv("DB_HOST"),
			Port:     getEnvOrDefault("DB_PORT", "5432"),
			Name:     os.Getenv("DB_NAME"),
			SSLMode:  getEnvOrDefault("DB_SSLMODE", "disable"),
		}

		databaseURL = fmt.Sprintf(
			"postgres://%s:%s@%s:%s/%s?sslmode=%s",
			db.User,
			url.QueryEscape(db.Password),
			db.Host,
			db.Port,
			db.Name,
			db.SSLMode,
		)
	}

	// Các biến bắt buộc khác
	requiredVars := []string{
		"REDIS_URL", "JWT_SECRET", "OPENWEATHER_API_KEY",
	}

	for _, v := range requiredVars {
		if os.Getenv(v) == "" {
			log.Fatalf("❌ Missing required environment variable: %s", v)
		} else {
			log.Printf("%q : %q", v, os.Getenv(v))
		}
	}

	Cfg = &Config{
		AppEnv:         getEnvOrDefault("APP_ENV", "development"),
		AppPort:        getEnvOrDefault("PORT", "8080"),
		DatabaseURL:    databaseURL,
		RedisURL:       os.Getenv("REDIS_URL"),
		JWTSecret:      []byte(os.Getenv("JWT_SECRET")),
		OpenWeatherKey: os.Getenv("OPENWEATHER_API_KEY"),
	}

	log.Println("✅ Config loaded successfully")
	// Log để debug (không log password)
	log.Printf("   App Env: %s", Cfg.AppEnv)
	log.Printf("   Port: %s", Cfg.AppPort)
	log.Printf("   DB Host: %s", extractHostFromURL(databaseURL))
}

// Helper function
func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// Extract host from DATABASE_URL for logging (without password)
func extractHostFromURL(dbURL string) string {
	if parsedURL, err := url.Parse(dbURL); err == nil {
		return parsedURL.Hostname()
	}
	return "unknown"
}
