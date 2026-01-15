package router

import (
	"github.com/LePhuocVuTien/SurvivalPro-Backend/infrastructure/ratelimit"
	"github.com/LePhuocVuTien/SurvivalPro-Backend/internal/handlers"
	"github.com/LePhuocVuTien/SurvivalPro-Backend/internal/middleware"
	"github.com/gofiber/fiber/v2"
)

// ================= AUTH =================
func authSetup(app *fiber.App, limiter *ratelimit.RedisLimiter) {

	api := app.Group("/api/v1")

	api.Post("/auth/login",
		middleware.RedisRateLimitMiddleware(limiter, "login"),
		handlers.Login,
	)

	api.Post("/auth/register",
		middleware.RedisRateLimitMiddleware(limiter, "register"),
		handlers.Register,
	)

	api.Post("/auth/forgot-password",
		middleware.RedisRateLimitMiddleware(limiter, "password_reset"),
		handlers.ForgotPassword,
	)

	api.Post("/auth/verify-email",
		middleware.RedisRateLimitMiddleware(limiter, "email_verify"),
		handlers.VerifyEmail,
	)
}
