package handlers

import (
	"context"

	"github.com/LePhuocVuTien/SurvivalPro-Backend/internal/db"
	"github.com/gofiber/fiber/v2"
)

func (h *LimiterHandler) Check(c *fiber.Ctx) error {
	ctx := context.Background()

	// Check database
	dbHealth := "ok"
	if err := db.HealthCheck(); err != nil {
		dbHealth = "error: " + err.Error()
	}

	// Check Redis
	redisHealth := "ok"
	if err := h.Limiter.Ping(ctx); err != nil {
		redisHealth = "error: " + err.Error()
	}

	// Get rate limiter stats
	stats, _ := h.Limiter.GetStats(ctx)

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
