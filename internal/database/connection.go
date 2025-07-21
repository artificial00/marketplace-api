package database

import (
	"database/sql"
	"fmt"
	"time"
)

// NewConnection создает новое подключение к PostgreSQL
func NewConnection(databaseURL string) (*sql.DB, error) {
	db, err := sql.Open("postgres", databaseURL)
	if err != nil {
		return nil, fmt.Errorf("failed to open database connection: %w", err)
	}

	db.SetMaxOpenConns(25)                 // Максимальное количество открытых соединений
	db.SetMaxIdleConns(25)                 // Максимальное количество неактивных соединений
	db.SetConnMaxLifetime(5 * time.Minute) // Максимальное время жизни соединения

	if err := db.Ping(); err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return db, nil
}

// CreateTables создает необходимые таблицы в базе данных
func CreateTables(db *sql.DB) error {
	createUsersTable := `
	CREATE TABLE IF NOT EXISTS users (
		id SERIAL PRIMARY KEY,
		login VARCHAR(50) UNIQUE NOT NULL,
		password_hash VARCHAR(255) NOT NULL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	)`

	createListingsTable := `
	CREATE TABLE IF NOT EXISTS listings (
		id SERIAL PRIMARY KEY,
		title VARCHAR(255) NOT NULL,
		description TEXT NOT NULL,
		image_url VARCHAR(500),
		price DECIMAL(10, 2) NOT NULL CHECK (price >= 0),
		user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	)`

	queries := []string{
		createUsersTable,
		createListingsTable,
	}

	for _, query := range queries {
		if _, err := db.Exec(query); err != nil {
			return fmt.Errorf("failed to execute query: %w", err)
		}
	}

	return nil
}

// CheckConnection проверяет, что соединение с базой данных активно
func CheckConnection(db *sql.DB) error {
	if err := db.Ping(); err != nil {
		return fmt.Errorf("database connection is not available: %w", err)
	}
	return nil
}
