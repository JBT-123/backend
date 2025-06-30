package database

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/How-to-get-ABG/backend/internal/config"
	"github.com/How-to-get-ABG/backend/internal/models"
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

func (p *PostgresDB) CreateUser(email, hashedPassword, username string) error {
	var id int;
	query := `INSERT INTO users (email, password, username) VALUES ($1, $2, $3) RETURNING id`
	err := p.db.QueryRow(query, email, hashedPassword, username).Scan(&id)
	if err != nil {
		fmt.Println("Problem retrieving id", err)
	}
	query = `INSERT INTO preferences (user_id) VALUES ($1)`
	_, err = p.db.Exec(query, id)
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

func (p *PostgresDB) Seed() error {
	// Check if test user exists
	testEmail := "test@example.com"
	existing, err := p.GetUserByEmail(testEmail)
	if err != nil {
		return err
	}
	if existing != nil {
		return nil // Already seeded
	}
	// Hash password
	hashed := "$2a$10$7a8Qw1Qw1Qw1Qw1Qw1Qw1u1Qw1Qw1Qw1Qw1Qw1Qw1Qw1Qw1Qw1Qw1" // bcrypt hash for 'password'
	return p.CreateUser(testEmail, hashed, "nil")
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

	statements := strings.Split(string(schemaSQL), ";")

    for _, stmt := range statements {
        stmt = strings.TrimSpace(stmt)
        if stmt == "" {
            continue
        }
        _, err := db.Exec(stmt)
        if err != nil {
            return fmt.Errorf("error executing statement %q: %v", stmt, err)
        }
    }

    return nil
}

func (p *PostgresDB) SQLDB() *sql.DB {
	return p.db
}


func (p *PostgresDB) UpdateUserDB(i *models.UserInfo, id int) error {
	b, _ := json.Marshal(i)
	var m map[string]interface{}
	json.Unmarshal(b, &m)

	for key, val := range m {
		if val == "" || val == nil {
			continue
		}

		query := fmt.Sprintf(`
			UPDATE users
			SET %s = $1
			WHERE id = $2
		`, key)

		_, err := p.db.Exec(query, val, id)
		if err != nil {
			fmt.Println("Error updating db:", err)
			return err
		}
	}

	return nil
}

func (p *PostgresDB) UpdateUserPref(i *models.UserPref, id int) error{
	b, _ := json.Marshal(i)
	var m map[string]any
	json.Unmarshal(b, &m)

	for key, val := range m {
		if val == "" || val == nil {
			continue
		}

		query := fmt.Sprintf(`
			UPDATE preferences
			SET %s = $1
			WHERE user_id = $2
		`, key)

		_, err := p.db.Exec(query, val, id)
		if err != nil {
			fmt.Println("Error updating db:", err)
			return err
		}
	}

	return nil
}
//models.UserConnectionRes
func (p *PostgresDB) UserLinkedPref(id int) ([]*models.UserConnectionRes, error){
	query:= `SELECT preferred_genders, interests FROM preferences WHERE user_id = $1`
	pref:=&models.UserPref{};
	err:= p.db.QueryRow(query, id).Scan(&pref.PrefGender, &pref.Interest)
	if err!=nil{
		fmt.Println("An error has occured: ", err)
		return nil, err
	}
	query= `
        SELECT  u.username, u.first_name, u.last_name, u.gender, u.birthdate, u.bio, u.profile_pic_url
        FROM users u
        JOIN preferences pref ON u.id = pref.user_id
        WHERE pref.preferred_genders = $1 AND pref.interests = $2
    `
	userInfo:=[]*models.UserConnectionRes{};
	rows, err:= p.db.Query(query, pref.PrefGender, pref.Interest)
	 if err != nil {
            fmt.Println("Err in query:", err)
            return nil, err
        }
	defer rows.Close()
	for rows.Next(){
		var user models.UserConnectionRes

        err := rows.Scan(
            &user.UserInfo.Username,
            &user.UserInfo.FirstName,
            &user.UserInfo.LastName,
            &user.UserInfo.Gender,
            &user.UserInfo.BirthDate,
            &user.UserInfo.Bio,
            &user.UserInfo.ProfilePicUrl,
        )
		if err = rows.Err(); err != nil {
		fmt.Println("Row iteration error:", err)
		return nil, err
	}
        // You can fill UserPref here if needed (you already know prefGender and interests)
        user.UserPref.PrefGender = pref.PrefGender
        user.UserPref.Interest = pref.Interest
        userInfo = append(userInfo, &user)
    }
	
	if err!=nil{
		fmt.Println("An error has occured: ", err)
		return nil, err
	}
	return userInfo, nil
}
