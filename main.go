package main

import (
	"log"
	"net/http"
	"os"

	"github.com/How-to-get-ABG/backend/internal/config"
	"github.com/How-to-get-ABG/backend/internal/database"
	"github.com/How-to-get-ABG/backend/internal/handlers"
	"github.com/How-to-get-ABG/backend/internal/middleware"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Printf("Warning: .env file not found")
	}

	// Initialize configuration
	cfg := config.New()

	// Initialize database connection
	db, err := database.NewPostgresDB(cfg)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Initialize handlers and middleware
	handlers.InitHandlers(db, cfg)
	middleware.InitMiddleware(cfg)

	// Initialize router
	router := mux.NewRouter()

	// Public routes
	router.HandleFunc("/api/register", handlers.Register).Methods("POST")
	router.HandleFunc("/api/login", handlers.Login).Methods("POST")

	// Protected routes
	protected := router.PathPrefix("/api").Subrouter()
	protected.Use(middleware.AuthMiddleware)
	protected.HandleFunc("/me", handlers.GetUserProfile).Methods("GET")

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server starting on port %s", port)
	log.Fatal(http.ListenAndServe(":"+port, router))
}
