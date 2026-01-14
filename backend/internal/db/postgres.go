package db

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/LePhuocVuTien/SurvivalPro-Backend/internal/config"
	_ "github.com/lib/pq"
)

var DB *sql.DB

// InitDB initializes database connection using config
func InitDB() error {
	// Get database URL from config
	dbURL := config.Cfg.Database.URL
	if dbURL == "" {
		return fmt.Errorf("❌ Database URL not configured")
	}

	var err error

	// Open database connection
	DB, err = sql.Open("postgres", dbURL)
	if err != nil {
		return fmt.Errorf("❌ Failed to open database: %w", err)
	}

	// Set connection pool settings from config
	DB.SetMaxOpenConns(config.Cfg.Database.MaxOpenConns)
	DB.SetMaxIdleConns(config.Cfg.Database.MaxIdleConns)
	DB.SetConnMaxLifetime(config.Cfg.Database.ConnMaxLifetime)
	DB.SetConnMaxIdleTime(config.Cfg.Database.ConnMaxIdleTime)

	// Test connection
	if err = DB.Ping(); err != nil {
		return fmt.Errorf("❌ Cannot ping database: %w", err)
	}

	// Create tables from schema
	if err := createTables(); err != nil {
		return fmt.Errorf("❌ Failed to create tables: %w", err)
	}

	log.Println("✅ Database connected successfully!")
	log.Printf("   Host: %s:%s", config.Cfg.Database.Host, config.Cfg.Database.Port)
	log.Printf("   Database: %s", config.Cfg.Database.Name)
	log.Printf("   Max Open Connections: %d", config.Cfg.Database.MaxOpenConns)
	log.Printf("   Max Idle Connections: %d", config.Cfg.Database.MaxIdleConns)

	return nil
}

// createTables creates database tables from schema file
func createTables() error {
	// Try multiple possible schema file locations
	schemaPaths := []string{
		"internal/db/schema.sql",
		"internal/database/schema.sql",
		"db/schema.sql",
		"schema.sql",
	}

	var schema []byte
	var err error
	var foundPath string

	for _, path := range schemaPaths {
		schema, err = os.ReadFile(path)
		if err == nil {
			foundPath = path
			break
		}
	}

	if foundPath == "" {
		log.Println("⚠️  No schema.sql file found, skipping table creation")
		return nil
	}

	// Execute schema
	_, err = DB.Exec(string(schema))
	if err != nil {
		return fmt.Errorf("error executing schema from %s: %w", foundPath, err)
	}

	log.Printf("✅ Tables created/updated from %s", foundPath)
	return nil
}

// CloseDB gracefully closes database connection
func CloseDB() error {
	if DB != nil {
		if err := DB.Close(); err != nil {
			return fmt.Errorf("failed to close database: %w", err)
		}
		log.Println("✅ Database connection closed")
	}
	return nil
}

// HealthCheck checks if database is healthy
func HealthCheck() error {
	if DB == nil {
		return fmt.Errorf("database not initialized")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := DB.PingContext(ctx); err != nil {
		return fmt.Errorf("database health check failed: %w", err)
	}

	return nil
}

// GetStats returns database connection pool statistics
func GetStats() sql.DBStats {
	if DB == nil {
		return sql.DBStats{}
	}
	return DB.Stats()
}

// ExecuteSchemaFile executes a specific SQL file
func ExecuteSchemaFile(filepath string) error {
	schema, err := os.ReadFile(filepath)
	if err != nil {
		return fmt.Errorf("failed to read schema file: %w", err)
	}

	_, err = DB.Exec(string(schema))
	if err != nil {
		return fmt.Errorf("failed to execute schema: %w", err)
	}

	log.Printf("✅ Executed schema from %s", filepath)
	return nil
}

// RunMigrations runs database migrations
func RunMigrations() error {
	// Check if migrations directory exists
	migrationPaths := []string{
		"internal/db/migrations",
		"internal/database/migrations",
		"db/migrations",
		"migrations",
	}

	var migrationsDir string
	for _, path := range migrationPaths {
		if _, err := os.Stat(path); err == nil {
			migrationsDir = path
			break
		}
	}

	if migrationsDir == "" {
		log.Println("⚠️  No migrations directory found, skipping migrations")
		return nil
	}

	// Read migration files
	files, err := os.ReadDir(migrationsDir)
	if err != nil {
		return fmt.Errorf("failed to read migrations directory: %w", err)
	}

	// Execute each migration file in order
	for _, file := range files {
		if file.IsDir() {
			continue
		}

		if !strings.HasSuffix(file.Name(), ".sql") {
			continue
		}

		filepath := fmt.Sprintf("%s/%s", migrationsDir, file.Name())
		log.Printf("Running migration: %s", file.Name())

		if err := ExecuteSchemaFile(filepath); err != nil {
			return fmt.Errorf("failed to run migration %s: %w", file.Name(), err)
		}
	}

	log.Println("✅ All migrations completed successfully")
	return nil
}
