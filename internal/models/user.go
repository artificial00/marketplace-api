package models

import "time"

type User struct {
	ID           int       `json:"id" db:"id"`
	Login        string    `json:"login" db:"login"`
	PasswordHash string    `json:"-" db:"password_hash"` // "-" скрывает поле в JSON
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time `json:"updated_at" db:"updated_at"`
}

// RegisterRequest структура для запроса регистрации
type RegisterRequest struct {
	Login    string `json:"login" binding:"required,min=3,max=50"`
	Password string `json:"password" binding:"required,min=6"`
}

// LoginRequest структура для запроса авторизации
type LoginRequest struct {
	Login    string `json:"login" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// AuthResponse структура ответа при успешной авторизации/регистрации
type AuthResponse struct {
	User  User   `json:"user"`
	Token string `json:"token"`
}
