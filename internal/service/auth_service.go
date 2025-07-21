package service

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
	"marketplace-api/internal/database/postgres"
	"marketplace-api/internal/models"
	"marketplace-api/pkg/utils"
)

type AuthService struct {
	userRepo  *postgres.UserRepository
	jwtSecret string
}

func NewAuthService(userRepo *postgres.UserRepository, jwtSecret string) *AuthService {
	return &AuthService{
		userRepo:  userRepo,
		jwtSecret: jwtSecret,
	}
}

type AuthServiceInterface interface {
	Register(req models.RegisterRequest) (*models.AuthResponse, error)
	Login(req models.LoginRequest) (*models.AuthResponse, error)
	GetUserByID(id int) (*models.User, error)
}

// Register регистрирует нового пользователя
func (s *AuthService) Register(req models.RegisterRequest) (*models.AuthResponse, error) {
	if !utils.ValidateLogin(req.Login) {
		return nil, fmt.Errorf("invalid login format")
	}

	if !utils.ValidatePassword(req.Password) {
		return nil, fmt.Errorf("password must be at least 6 characters long and contain letters and digits")
	}

	exists, err := s.userRepo.UserExists(req.Login)
	if err != nil {
		return nil, fmt.Errorf("failed to check user existence: %w", err)
	}
	if exists {
		return nil, fmt.Errorf("user already exists")
	}

	// Хешируем пароль
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	user, err := s.userRepo.CreateUser(req.Login, string(passwordHash))
	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	token, err := utils.GenerateToken(user.ID, user.Login, s.jwtSecret)
	if err != nil {
		return nil, fmt.Errorf("failed to generate token: %w", err)
	}

	return &models.AuthResponse{
		User:  *user,
		Token: token,
	}, nil
}

// Login авторизует пользователя
func (s *AuthService) Login(req models.LoginRequest) (*models.AuthResponse, error) {
	user, err := s.userRepo.GetUserByLogin(req.Login)
	if err != nil {
		return nil, fmt.Errorf("invalid credentials")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password))
	if err != nil {
		return nil, fmt.Errorf("invalid credentials")
	}

	token, err := utils.GenerateToken(user.ID, user.Login, s.jwtSecret)
	if err != nil {
		return nil, fmt.Errorf("failed to generate token: %w", err)
	}

	return &models.AuthResponse{
		User:  *user,
		Token: token,
	}, nil
}

// GetUserByID получает пользователя по ID
func (s *AuthService) GetUserByID(id int) (*models.User, error) {
	return s.userRepo.GetUserByID(id)
}
