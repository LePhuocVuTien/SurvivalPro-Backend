package handlers

import (
	"context"
	"time"

	"github.com/gofiber/fiber/v2"
)

// ============================================================================
// ADMIN HANDLERS FOR RATE LIMIT MANAGEMENT
// ============================================================================

func (h *LimiterHandler) GetRateLimitStatus(c *fiber.Ctx) error {
	identifier := c.Params("identifier")
	action := c.Params("action")

	ctx := context.Background()
	status, err := h.Limiter.GetStatus(ctx, identifier, action)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "Failed to get rate limit status",
		})
	}

	return c.JSON(status)
}

func (h *LimiterHandler) ResetRateLimit(c *fiber.Ctx) error {
	identifier := c.Params("identifier")
	action := c.Params("action")

	ctx := context.Background()
	if err := h.Limiter.Reset(ctx, identifier, action); err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "Failed to reset rate limit",
		})
	}

	return c.JSON(fiber.Map{
		"message": "Rate limit reset successfully",
	})
}

func (h *LimiterHandler) BlockUser(c *fiber.Ctx) error {
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
	if err := h.Limiter.Block(ctx, identifier, action, duration); err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "Failed to block user",
		})
	}

	return c.JSON(fiber.Map{
		"message": "User blocked successfully",
	})
}

func (h *LimiterHandler) UnblockUser(c *fiber.Ctx) error {
	identifier := c.Params("identifier")
	action := c.Params("action")

	ctx := context.Background()
	if err := h.Limiter.Unblock(ctx, identifier, action); err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "Failed to unblock user",
		})
	}

	return c.JSON(fiber.Map{
		"message": "User unblocked successfully",
	})
}

func (h *LimiterHandler) GetRateLimitStats(c *fiber.Ctx) error {
	ctx := context.Background()
	stats, err := h.Limiter.GetStats(ctx)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "Failed to get statistics",
		})
	}

	return c.JSON(stats)
}

func (h *LimiterHandler) ListRules(c *fiber.Ctx) error {
	rules := h.Limiter.ListRules()
	return c.JSON(fiber.Map{
		"rules": rules,
	})
}
