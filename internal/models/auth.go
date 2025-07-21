package models

// Константы для ошибок аутентификации
const (
	ErrUserExists   = "user_exists"
	ErrUserNotFound = "user_not_found"
)

// AuthError структура для ошибок аутентификации
type AuthError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

// ValidationError структура для ошибок валидации
type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}
