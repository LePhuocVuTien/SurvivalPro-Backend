package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/LePhuocVuTien/SurvivalPro-Backend/handlers"
	"github.com/LePhuocVuTien/SurvivalPro-Backend/internal/config"
	"github.com/LePhuocVuTien/SurvivalPro-Backend/internal/db"
	"github.com/LePhuocVuTien/SurvivalPro-Backend/middleware"
	"github.com/LePhuocVuTien/SurvivalPro-Backend/redis"
	"github.com/gorilla/mux"
)

func main() {
	config.Load()

	// Initialize DB với error handling tốt hơn
	log.Println("Connecting to database...")
	if err := db.InitDB(); err != nil {
		log.Printf("⚠️  Database connection failed: %v", err)
		log.Println("Server will start without database connection")
	} else {
		defer db.CloseDB()
	}

	// Initialize Redis với error handling
	log.Println("Connecting to Redis...")
	if err := redis.InitRedis(); err != nil {
		log.Printf("⚠️  Redis connection failed: %v", err)
		log.Println("Server will start without Redis")
	} else {
		defer redis.CloseRedis()
	}

	// Start clearup goroutine for rate limiter
	go middleware.ClearupVisitors()

	router := mux.NewRouter()
	router.PathPrefix("/uploads/").Handler(http.StripPrefix("/uploads/", http.FileServer(http.Dir("internal/uploads"))))

	// Health check endpoint
	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "SurvivalPro API is running!")
	}).Methods("GET")

	// Public routes with rate limiting
	router.HandleFunc("/users/register", middleware.RateLimit(handlers.RegisterUser)).Methods("POST")
	router.HandleFunc("/users/login", middleware.RateLimit(handlers.LoginUser)).Methods("POST")
	router.HandleFunc("/guides", middleware.RateLimit(handlers.GetAllGuides)).Methods("GET")
	router.HandleFunc("/weather", middleware.RateLimit(handlers.GetCurrentWeather)).Methods("GET")

	// Protected routes
	router.HandleFunc("/users/profile", middleware.RateLimit(middleware.Auth(handlers.UpdateUserProfile))).Methods("PUT")
	router.HandleFunc("/users/avatar", middleware.RateLimit(middleware.Auth(handlers.UploadAvatar))).Methods("POST")
	router.HandleFunc("/users/push-token", middleware.RateLimit(middleware.Auth(handlers.RegisterPushToken))).Methods("POST")
	router.HandleFunc("/notifications", middleware.RateLimit(middleware.Auth(handlers.GetUserNotifications))).Methods("GET")

	router.HandleFunc("/checklist", middleware.RateLimit(middleware.Auth(handlers.GetUserChecklist))).Methods("GET")
	router.HandleFunc("/checklist", middleware.RateLimit(middleware.Auth(handlers.CreateChecklistItem))).Methods("POST")
	router.HandleFunc("/checklist/{id}", middleware.RateLimit(middleware.Auth(handlers.UpdateCheckListItem))).Methods("PUT")
	router.HandleFunc("/checklist/{id}", middleware.RateLimit(middleware.Auth(handlers.DeleteChecklistItem))).Methods("DELETE")

	router.HandleFunc("/location", middleware.RateLimit(middleware.Auth(handlers.SaveUserLocation))).Methods("POST")
	router.HandleFunc("/location", middleware.RateLimit(middleware.Auth(handlers.GetUserLocation))).Methods("GET")
	router.HandleFunc("/guides/upload-image", middleware.RateLimit(middleware.Auth(handlers.UploadGuideImage))).Methods("POST")

	// Enable CORS
	handler := middleware.EnableCORS(router)

	port := os.Getenv("APP_PORT")
	if port == "" {
		port = "8080" // Default cho local development
	}

	// Improtant: Listen 0.0.0.0, not localhost
	addr := "0.0.0.0:" + port
	log.Printf("🚀 Server running on %s", addr)

	if err := http.ListenAndServe(addr, handler); err != nil {
		log.Fatalf("❌ Server failed to start: %v", err)
	}

}
