package db

import (
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

func InitDB() error {

	dbURL := config.Cfg.DatabaseURL
	if dbURL == "" {
		return fmt.Errorf("❌ DATABASE_URL not set")
	}
	var err error

	if !strings.Contains(dbURL, "sslmode=") {
		dbURL += "?sslmode=disable"
	}

	DB, err = sql.Open("postgres", dbURL)
	if err != nil {
		return fmt.Errorf("❌ failed to open database: %w", err)
	}

	DB.SetMaxOpenConns(25)
	DB.SetMaxIdleConns(5)
	DB.SetConnMaxLifetime(30 * time.Minute)

	if err = DB.Ping(); err != nil {
		log.Println(string(dbURL))
		return fmt.Errorf("❌ Cannot ping database: %w", err)
	}

	createTable()

	log.Println("✅ Database connected successfully!")
	return nil
}

func createTable() {
	schema, _ := os.ReadFile("internal/db/schema.sql")
	_, err := DB.Exec(string(schema))
	if err != nil {
		log.Fatal("Error creating tables:", err)
	}
}

func CloseDB() {
	if DB != nil {
		DB.Close()
	}
}
