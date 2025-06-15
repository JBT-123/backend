package main

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

func (s *Seed) UserSeed() {
	email := "test@example.com"
	password := "password"
	hash, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	_, err := s.DB.Exec(`INSERT INTO users (email, password) VALUES ($1, $2) ON CONFLICT (email) DO NOTHING`, email, string(hash))
	if err != nil {
		fmt.Printf("Error seeding user: %v\n", err)
	} else {
		fmt.Println("Seeded user: test@example.com / password")
	}
}
