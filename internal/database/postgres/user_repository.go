package postgres

import (
	"database/sql"
	"fmt"

	"marketplace-api/internal/models"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

// CreateUser создает нового пользователя
func (r *UserRepository) CreateUser(login, passwordHash string) (*models.User, error) {
	query := `
		INSERT INTO users (login, password_hash) 
		VALUES ($1, $2) 
		RETURNING id, login, password_hash, created_at, updated_at
	`

	var user models.User
	err := r.db.QueryRow(query, login, passwordHash).Scan(
		&user.ID,
		&user.Login,
		&user.PasswordHash,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	return &user, nil
}

// GetUserByLogin получает пользователя по логину
func (r *UserRepository) GetUserByLogin(login string) (*models.User, error) {
	query := `
		SELECT id, login, password_hash, created_at, updated_at 
		FROM users 
		WHERE login = $1
	`

	var user models.User
	err := r.db.QueryRow(query, login).Scan(
		&user.ID,
		&user.Login,
		&user.PasswordHash,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	return &user, nil
}

// GetUserByID получает пользователя по ID
func (r *UserRepository) GetUserByID(id int) (*models.User, error) {
	query := `
		SELECT id, login, password_hash, created_at, updated_at 
		FROM users 
		WHERE id = $1
	`

	var user models.User
	err := r.db.QueryRow(query, id).Scan(
		&user.ID,
		&user.Login,
		&user.PasswordHash,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	return &user, nil
}

// UserExists проверяет, существует ли пользователь с таким логином
func (r *UserRepository) UserExists(login string) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM users WHERE login = $1)`

	var exists bool
	err := r.db.QueryRow(query, login).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("failed to check user existence: %w", err)
	}

	return exists, nil
}
