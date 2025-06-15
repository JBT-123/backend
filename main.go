package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"reflect"

	"github.com/How-to-get-ABG/backend/internal/config"
	"github.com/How-to-get-ABG/backend/internal/database"
	"github.com/How-to-get-ABG/backend/internal/handlers"
	"github.com/How-to-get-ABG/backend/internal/middleware"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"golang.org/x/crypto/bcrypt"
)

type Seed struct {
	DB *sql.DB
}

func (s *Seed) UserSeed() {
	email := "test@example.com"
	password := "password"
	hash, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	_, err := s.DB.Exec(`INSERT INTO users (email, password) VALUES ($1, $2) ON CONFLICT (email) DO NOTHING`, email, string(hash))
	if err != nil {
		log.Printf("Error seeding user: %v\n", err)
	} else {
		log.Println("Seeded user: test@example.com / password")
	}
}

func (s *Seed) Execute(methods ...string) error {
	seedType := reflect.TypeOf(s)
	seedValue := reflect.ValueOf(s)

	if len(methods) == 0 {
		for i := 0; i < seedType.NumMethod(); i++ {
			method := seedType.Method(i)
			if method.Name == "Execute" {
				continue
			}
			log.Printf("Running seeder: %s\n", method.Name)
			seedValue.MethodByName(method.Name).Call(nil)
		}
		return nil
	}

	for _, name := range methods {
		method := seedValue.MethodByName(name)
		if !method.IsValid() {
			return fmt.Errorf("Seeder method %s not found", name)
		}
		log.Printf("Running seeder: %s\n", name)
		method.Call(nil)
	}
	return nil
}

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Printf("Warning: .env file not found")
	}

	cfg := config.New()

	// Check for seeder argument
	if len(os.Args) > 1 && os.Args[1] == "seed" {
		dbObj, err := database.NewPostgresDB(cfg)
		if err != nil {
			log.Fatalf("Failed to connect to database: %v", err)
		}
		defer dbObj.Close()

		seeder := &Seed{DB: dbObj.SQLDB()}
		if len(os.Args) > 2 {
			if err := seeder.Execute(os.Args[2:]...); err != nil {
				log.Fatalf("Seeding error: %v", err)
			}
		} else {
			if err := seeder.Execute(); err != nil {
				log.Fatalf("Seeding error: %v", err)
			}
		}
		return
	}

	db, err := database.NewPostgresDB(cfg)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	handlers.InitHandlers(db, cfg)
	middleware.InitMiddleware(cfg)

	router := mux.NewRouter()

	router.HandleFunc("/api/register", handlers.Register).Methods("POST")
	router.HandleFunc("/api/login", handlers.Login).Methods("POST")

	protected := router.PathPrefix("/api").Subrouter()
	protected.Use(middleware.AuthMiddleware)
	protected.HandleFunc("/me", handlers.GetUserProfile).Methods("GET")

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server starting on port %s", port)
	log.Fatal(http.ListenAndServe(":"+port, router))
}
