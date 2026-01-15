package router

import (
	"github.com/LePhuocVuTien/SurvivalPro-Backend/infrastructure/ratelimit"
	"github.com/LePhuocVuTien/SurvivalPro-Backend/internal/handlers"
	"github.com/LePhuocVuTien/SurvivalPro-Backend/internal/middleware"
	"github.com/gofiber/fiber/v2"
)

// ================= AUTHENTICATED =================
func authenticatedSetup(app *fiber.App, limiter *ratelimit.RedisLimiter) {
	api := app.Group("/api/v1")

	auth := api.Group("/",
		middleware.RedisRateLimitMiddleware(limiter, "api"),
	)

	users := auth.Group("/users")
	users.Get("/", handlers.ListUsers)
	users.Get("/:id", handlers.GetUser)
	users.Put("/:id", handlers.UpdateUser)
	users.Delete("/:id", handlers.DeleteUser)

	auth.Post("/upload",
		middleware.RedisRateLimitMiddleware(limiter, "upload"),
		handlers.Upload,
	)

	admin := auth.Group("/admin", middleware.RequireAdmin())

	newHandler := handlers.NewHandler(limiter)
	admin.Get("/rate-limits/:identifier/:action", newHandler.GetRateLimitStatus)
	admin.Delete("/rate-limits/:identifier/:action", newHandler.ResetRateLimit)
	admin.Post("/rate-limits/:identifier/:action/block", newHandler.BlockUser)
	admin.Post("/rate-limits/:identifier/:action/unblock", newHandler.UnblockUser)
	admin.Get("/rate-limits/stats", newHandler.GetRateLimitStats)
	admin.Get("/rate-limits/rules", newHandler.ListRules)
}
