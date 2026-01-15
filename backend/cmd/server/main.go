package main

import (
	"context"
	"log"

	"github.com/LePhuocVuTien/SurvivalPro-Backend/infrastructure/ratelimit"
	"github.com/LePhuocVuTien/SurvivalPro-Backend/infrastructure/shutdown"
	"github.com/LePhuocVuTien/SurvivalPro-Backend/internal/config"
	"github.com/LePhuocVuTien/SurvivalPro-Backend/internal/db"
	"github.com/LePhuocVuTien/SurvivalPro-Backend/internal/handlers"
	"github.com/LePhuocVuTien/SurvivalPro-Backend/internal/router"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"    // ✅
	"github.com/gofiber/fiber/v2/middleware/logger"  // ✅
	"github.com/gofiber/fiber/v2/middleware/recover" // ✅
)

var (
	// Global limiter instance
	limiter *ratelimit.RedisLimiter
)

func main() {
	// =========================================================================
	// LOAD CONFIGURATION
	// =========================================================================
	log.Println("🚀 Starting application...")
	config.MustLoad()

	// =========================================================================
	// INITIALIZE DATABASE
	// =========================================================================
	if err := db.InitDB(); err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer db.CloseDB()

	// =========================================================================
	// INITIALIZE REDIS & RATE LIMITER
	// =========================================================================
	limiter = ratelimit.InitializeRedisLimiter()
	defer limiter.Close()

	// Test Redis connection
	ctx := context.Background()
	if err := limiter.Ping(ctx); err != nil {
		log.Fatalf("❌ Redis connection failed: %v", err)
	}
	log.Println("✅ Redis rate limiter initialized successfully!")

	// =========================================================================
	// SETUP FIBER APP
	// =========================================================================
	app := fiber.New(fiber.Config{
		AppName:      config.Cfg.App.Name + " v" + config.Cfg.App.Version,
		ErrorHandler: handlers.ErrorHandler,
	})

	// Middleware
	app.Use(recover.New())
	app.Use(logger.New())
	app.Use(cors.New(cors.Config{
		AllowOrigins: config.Cfg.App.FrontendURL,
		AllowMethods: "GET,POST,PUT,DELETE,PATCH,OPTIONS",
		AllowHeaders: "Origin,Content-Type,Accept,Authorization",
	}))

	// Setup routes with rate limiting
	router.Setup(app, limiter)

	go shutdown.Graceful(shutdown.Resources{
		App: app,
		CloseDB: func() error {
			return db.CloseDB()
		},
		CloseLimiter: func() error {
			return limiter.Close()
		},
	})

	// Start server
	log.Printf("🚀 Server starting on port %s", config.Cfg.App.Port)
	log.Fatal(app.Listen(":" + config.Cfg.App.Port))
}
