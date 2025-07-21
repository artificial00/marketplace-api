package postgres

import (
	"fmt"

	"marketplace-api/internal/models"
)

type AuthRepository struct {
	userRepo *UserRepository
}

func NewAuthRepository(userRepo *UserRepository) *AuthRepository {
	return &AuthRepository{
		userRepo: userRepo,
	}
}

// RegisterUser регистрирует нового пользователя (обертка с обработкой ошибок)
func (r *AuthRepository) RegisterUser(login, passwordHash string) (*models.User, error) {
	exists, err := r.userRepo.UserExists(login)
	if err != nil {
		return nil, fmt.Errorf("failed to check user existence: %w", err)
	}

	if exists {
		return nil, fmt.Errorf(models.ErrUserExists)
	}

	user, err := r.userRepo.CreateUser(login, passwordHash)
	if err != nil {
		return nil, fmt.Errorf("failed to register user: %w", err)
	}

	return user, nil
}

// AuthenticateUser получает пользователя для аутентификации
func (r *AuthRepository) AuthenticateUser(login string) (*models.User, error) {
	user, err := r.userRepo.GetUserByLogin(login)
	if err != nil {
		if err.Error() == "user not found" {
			return nil, fmt.Errorf(models.ErrUserNotFound)
		}
		return nil, fmt.Errorf("failed to authenticate user: %w", err)
	}

	return user, nil
}

// ValidateUserToken проверяет существование пользователя по ID (для валидации токенов)
func (r *AuthRepository) ValidateUserToken(userID int) (*models.User, error) {
	user, err := r.userRepo.GetUserByID(userID)
	if err != nil {
		if err.Error() == "user not found" {
			return nil, fmt.Errorf(models.ErrUserNotFound)
		}
		return nil, fmt.Errorf("failed to validate user token: %w", err)
	}

	return user, nil
}
