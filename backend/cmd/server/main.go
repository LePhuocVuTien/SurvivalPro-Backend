package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	// "your-project/infrastructure/ratelimit"
	// "your-project/internal/database"
	// "your-project/middleware"

	"github.com/LePhuocVuTien/SurvivalPro-Backend/infrastructure/ratelimit"
	"github.com/LePhuocVuTien/SurvivalPro-Backend/internal/config"
	"github.com/LePhuocVuTien/SurvivalPro-Backend/internal/db"
	"github.com/LePhuocVuTien/SurvivalPro-Backend/internal/middleware"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"    // ✅
	"github.com/gofiber/fiber/v2/middleware/logger"  // ✅
	"github.com/gofiber/fiber/v2/middleware/recover" // ✅
	"github.com/redis/go-redis/v9"
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
	limiter = initializeRedisLimiter()
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
		ErrorHandler: customErrorHandler,
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
	//setupRoutes(app, limiter)

	// Start server
	log.Printf("🚀 Server starting on port %s", config.Cfg.App.Port)
	log.Fatal(app.Listen(":" + config.Cfg.App.Port))
}

// ============================================================================
// REDIS LIMITER INITIALIZATION
// ============================================================================

func initializeRedisLimiter() *ratelimit.RedisLimiter {
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
	rules := []*ratelimit.RateLimitRule{
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
	return ratelimit.NewRedisLimiter(client, rules)
}

// ============================================================================
// ROUTE SETUP WITH RATE LIMITING
// ============================================================================

func setupRoutes(app *fiber.App) {
	api := app.Group("/api/v1")

	// =========================================================================
	// PUBLIC ROUTES WITH RATE LIMITING
	// =========================================================================

	// Login with rate limiting
	api.Post("/auth/login",
		middleware.RedisRateLimitMiddleware(limiter, "login"),
		handleLogin,
	)

	// Register with rate limiting
	api.Post("/auth/register",
		middleware.RedisRateLimitMiddleware(limiter, "register"),
		handleRegister,
	)

	// Password reset with rate limiting
	api.Post("/auth/forgot-password",
		middleware.RedisRateLimitMiddleware(limiter, "password_reset"),
		handleForgotPassword,
	)

	// Email verification with rate limiting
	api.Post("/auth/verify-email",
		middleware.RedisRateLimitMiddleware(limiter, "email_verify"),
		handleVerifyEmail,
	)

	// =========================================================================
	// AUTHENTICATED ROUTES WITH API RATE LIMITING
	// =========================================================================

	auth := api.Group("/",
		// Your JWT middleware here
		middleware.RedisRateLimitMiddleware(limiter, "api"),
	)

	// User routes
	users := auth.Group("/users")
	users.Get("/", handleListUsers)
	users.Get("/:id", handleGetUser)
	users.Put("/:id", handleUpdateUser)
	users.Delete("/:id", handleDeleteUser)

	// File upload with specific rate limiting
	auth.Post("/upload",
		middleware.RedisRateLimitMiddleware(limiter, "upload"),
		handleUpload,
	)

	// =========================================================================
	// ADMIN ROUTES (rate limit management)
	// =========================================================================

	admin := auth.Group("/admin", middleware.RequireAdmin())

	// Get rate limit status
	admin.Get("/rate-limits/:identifier/:action", handleGetRateLimitStatus)

	// Reset rate limit
	admin.Delete("/rate-limits/:identifier/:action", handleResetRateLimit)

	// Block user
	admin.Post("/rate-limits/:identifier/:action/block", handleBlockUser)

	// Unblock user
	admin.Post("/rate-limits/:identifier/:action/unblock", handleUnblockUser)

	// Get rate limit statistics
	admin.Get("/rate-limits/stats", handleGetRateLimitStats)

	// List all rules
	admin.Get("/rate-limits/rules", handleListRules)

	// =========================================================================
	// HEALTH CHECK
	// =========================================================================

	api.Get("/health", handleHealthCheck)
}

// ============================================================================
// HANDLERS
// ============================================================================

// users.Get("/:id", handlers.HandleGetUser)
func handleLogin(c *fiber.Ctx) error {
	// Your login logic here
	return c.JSON(fiber.Map{
		"message": "Login successful",
		"token":   "jwt_token_here",
	})
}

func handleRegister(c *fiber.Ctx) error {
	// Your registration logic here
	return c.JSON(fiber.Map{
		"message": "Registration successful",
	})
}

func handleForgotPassword(c *fiber.Ctx) error {
	// Your password reset logic here
	return c.JSON(fiber.Map{
		"message": "Password reset email sent",
	})
}

func handleVerifyEmail(c *fiber.Ctx) error {
	// Your email verification logic here
	return c.JSON(fiber.Map{
		"message": "Email verified",
	})
}

func handleListUsers(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{"users": []interface{}{}})
}

func handleGetUser(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{"user": fiber.Map{}})
}

func handleUpdateUser(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{"message": "User updated"})
}

func handleDeleteUser(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{"message": "User deleted"})
}

func handleUpload(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{"message": "File uploaded"})
}

func handleHealthCheck(c *fiber.Ctx) error {
	ctx := context.Background()

	// Check database
	dbHealth := "ok"
	if err := db.HealthCheck(); err != nil {
		dbHealth = "error: " + err.Error()
	}

	// Check Redis
	redisHealth := "ok"
	if err := limiter.Ping(ctx); err != nil {
		redisHealth = "error: " + err.Error()
	}

	// Get rate limiter stats
	stats, _ := limiter.GetStats(ctx)

	return c.JSON(fiber.Map{
		"status": "ok",
		"database": fiber.Map{
			"status": dbHealth,
		},
		"redis": fiber.Map{
			"status": redisHealth,
		},
		"rate_limiter": stats,
	})
}

// ============================================================================
// ADMIN HANDLERS FOR RATE LIMIT MANAGEMENT
// ============================================================================

func handleGetRateLimitStatus(c *fiber.Ctx) error {
	identifier := c.Params("identifier")
	action := c.Params("action")

	ctx := context.Background()
	status, err := limiter.GetStatus(ctx, identifier, action)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "Failed to get rate limit status",
		})
	}

	return c.JSON(status)
}

func handleResetRateLimit(c *fiber.Ctx) error {
	identifier := c.Params("identifier")
	action := c.Params("action")

	ctx := context.Background()
	if err := limiter.Reset(ctx, identifier, action); err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "Failed to reset rate limit",
		})
	}

	return c.JSON(fiber.Map{
		"message": "Rate limit reset successfully",
	})
}

func handleBlockUser(c *fiber.Ctx) error {
	identifier := c.Params("identifier")
	action := c.Params("action")

	// Parse duration from request body
	var req struct {
		Duration string `json:"duration"` // e.g., "1h", "30m"
	}
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "Invalid request",
		})
	}

	duration, err := time.ParseDuration(req.Duration)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "Invalid duration format",
		})
	}

	ctx := context.Background()
	if err := limiter.Block(ctx, identifier, action, duration); err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "Failed to block user",
		})
	}

	return c.JSON(fiber.Map{
		"message": "User blocked successfully",
	})
}

func handleUnblockUser(c *fiber.Ctx) error {
	identifier := c.Params("identifier")
	action := c.Params("action")

	ctx := context.Background()
	if err := limiter.Unblock(ctx, identifier, action); err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "Failed to unblock user",
		})
	}

	return c.JSON(fiber.Map{
		"message": "User unblocked successfully",
	})
}

func handleGetRateLimitStats(c *fiber.Ctx) error {
	ctx := context.Background()
	stats, err := limiter.GetStats(ctx)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "Failed to get statistics",
		})
	}

	return c.JSON(stats)
}

func handleListRules(c *fiber.Ctx) error {
	rules := limiter.ListRules()
	return c.JSON(fiber.Map{
		"rules": rules,
	})
}

// ============================================================================
// ERROR HANDLER
// ============================================================================

func customErrorHandler(c *fiber.Ctx, err error) error {
	code := fiber.StatusInternalServerError
	message := "Internal Server Error"

	if e, ok := err.(*fiber.Error); ok {
		code = e.Code
		message = e.Message
	}

	if config.Cfg.App.Debug {
		log.Printf("Error: %v", err)
	}

	return c.Status(code).JSON(fiber.Map{
		"error":   message,
		"details": err.Error(),
	})
}

// ============================================================================
// GRACEFUL SHUTDOWN
// ============================================================================

func gracefulShutdown(app *fiber.App) {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	<-quit

	log.Println("🛑 Shutting down server...")

	// Shutdown Fiber
	if err := app.Shutdown(); err != nil {
		log.Printf("Error during server shutdown: %v", err)
	}

	// Close database
	if err := db.CloseDB(); err != nil {
		log.Printf("Error closing database: %v", err)
	}

	// Close Redis rate limiter
	if err := limiter.Close(); err != nil {
		log.Printf("Error closing Redis: %v", err)
	}

	log.Println("✅ Server stopped gracefully")
	os.Exit(0)
}
