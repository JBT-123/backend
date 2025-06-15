package main

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/How-to-get-ABG/backend/internal/config"
	_ "github.com/lib/pq"
)

func main() {
	cfg := config.New()
	connStr := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPassword, cfg.DBName,
	)
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	seeder := &Seed{DB: db}
	args := os.Args[1:]
	if len(args) > 0 {
		if err := seeder.Execute(args...); err != nil {
			panic(err)
		}
	} else {
		if err := seeder.Execute(); err != nil {
			panic(err)
		}
	}
}
