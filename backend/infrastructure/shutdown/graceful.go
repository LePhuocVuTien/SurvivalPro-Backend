package shutdown

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/gofiber/fiber/v2"
)

type Resources struct {
	App          *fiber.App
	CloseDB      func() error
	CloseLimiter func() error
}

func Graceful(r Resources) {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	<-quit
	log.Println("🛑 Shutting down server...")

	if r.App != nil {
		if err := r.App.Shutdown(); err != nil {
			log.Printf("Shutdown app error: %v", err)
		}
	}

	if r.CloseDB != nil {
		if err := r.CloseDB(); err != nil {
			log.Printf("Close DB error: %v", err)
		}
	}

	if r.CloseLimiter != nil {
		if err := r.CloseLimiter(); err != nil {
			log.Printf("Close limiter error: %v", err)
		}
	}

	log.Println("✅ Server stopped gracefully")
}
