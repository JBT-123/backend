package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/How-to-get-ABG/backend/internal/config"
	"github.com/How-to-get-ABG/backend/internal/database"
	"github.com/How-to-get-ABG/backend/internal/models"
	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
)

var (
	db  *database.PostgresDB
	cfg *config.Config
)

func InitHandlers(database *database.PostgresDB, config *config.Config) {
	db = database
	cfg = config
}

func Register(w http.ResponseWriter, r *http.Request) {
	var req models.RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Check if user already exists
	existingUser, err := db.GetUserByEmail(req.Email)
	if err != nil {
		http.Error(w, "Error checking user existence", http.StatusInternalServerError)
		return
	}
	if existingUser != nil {
		http.Error(w, "User already exists", http.StatusConflict)
		return
	}
	
	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "Error processing password", http.StatusInternalServerError)
		return
	}

	username:= fmt.Sprintf("guest_%d", time.Now().Unix())

	// Create user in database
	if err := db.CreateUser(req.Email, string(hashedPassword), username); err != nil {
		http.Error(w, "Error creating user", http.StatusInternalServerError)
		return
	}

	// Get the created user
	user, err := db.GetUserByEmail(req.Email)
	if err != nil {
		http.Error(w, "Error retrieving created user", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(models.User{
		ID:    user.ID,
		Email: user.Email,
	})
}

func Login(w http.ResponseWriter, r *http.Request) {
	var req models.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Get user from database
	user, err := db.GetUserByEmail(req.Email)
	if err != nil {
		http.Error(w, "Error retrieving user", http.StatusInternalServerError)
		return
	}
	if user == nil {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	// Verify password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	// Generate JWT token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID,
		"email":   user.Email,
		"exp":     time.Now().Add(time.Hour * 72).Unix(),
	})

	tokenString, err := token.SignedString([]byte(cfg.JWTSecret))
	if err != nil {
		http.Error(w, "Error generating token", http.StatusInternalServerError)
		return
	}

	response := models.AuthResponse{
		Token: tokenString,
		User: models.User{
			ID:    user.ID,
			Email: user.Email,
		},
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func GetUserProfile(w http.ResponseWriter, r *http.Request) {
	// Get user ID from context (set by middleware)
	userID := r.Context().Value("user_id").(int)

	// Get user from database
	user, err := db.GetUserByID(userID)
	if err != nil {
		http.Error(w, "Error retrieving user", http.StatusInternalServerError)
		return
	}
	if user == nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(models.User{
		ID:    user.ID,
		Email: user.Email,
	})
}
