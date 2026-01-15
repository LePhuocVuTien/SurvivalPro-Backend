package handlers

import (
	"log"

	"github.com/LePhuocVuTien/SurvivalPro-Backend/internal/config"
	"github.com/gofiber/fiber/v2"
)

// ============================================================================
// ERROR HANDLER
// ============================================================================

func ErrorHandler(c *fiber.Ctx, err error) error {
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
