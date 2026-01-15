package router

import (
	"github.com/LePhuocVuTien/SurvivalPro-Backend/infrastructure/ratelimit"
	"github.com/LePhuocVuTien/SurvivalPro-Backend/internal/handlers"
	"github.com/gofiber/fiber/v2"
)

func Setup(app *fiber.App, limiter *ratelimit.RedisLimiter) {

	// Check Health
	newHandler := handlers.NewHandler(limiter)
	api := app.Group("/api/v1")
	api.Get("/health", newHandler.Check)

	authSetup(app, limiter)
	authenticatedSetup(app, limiter)
}
