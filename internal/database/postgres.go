package database

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/How-to-get-ABG/backend/internal/config"
	_ "github.com/lib/pq"
)

type PostgresDB struct {
	db *sql.DB
}

func NewPostgresDB(cfg *config.Config) (*PostgresDB, error) {
	connStr := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPassword, cfg.DBName,
	)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("error opening database: %v", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("error connecting to the database: %v", err)
	}

	// Initialize schema
	if err := initSchema(db); err != nil {
		return nil, fmt.Errorf("error initializing schema: %v", err)
	}

	return &PostgresDB{db: db}, nil
}

func (p *PostgresDB) Close() error {
	return p.db.Close()
}

func (p *PostgresDB) CreateUser(email, hashedPassword string) error {
	query := `INSERT INTO users (email, password) VALUES ($1, $2)`
	_, err := p.db.Exec(query, email, hashedPassword)
	return err
}

func (p *PostgresDB) GetUserByEmail(email string) (*User, error) {
	query := `SELECT id, email, password FROM users WHERE email = $1`
	user := &User{}
	err := p.db.QueryRow(query, email).Scan(&user.ID, &user.Email, &user.Password)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (p *PostgresDB) GetUserByID(id int) (*User, error) {
	query := `SELECT id, email, password FROM users WHERE id = $1`
	user := &User{}
	err := p.db.QueryRow(query, id).Scan(&user.ID, &user.Email, &user.Password)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return user, nil
}

type User struct {
	ID       int
	Email    string
	Password string
}

func initSchema(db *sql.DB) error {
	// Read schema file
	schemaSQL, err := os.ReadFile("internal/database/schema.sql")
	if err != nil {
		return fmt.Errorf("error reading schema file: %v", err)
	}

	// Execute schema
	_, err = db.Exec(string(schemaSQL))
	if err != nil {
		return fmt.Errorf("error executing schema: %v", err)
	}

	return nil
}
