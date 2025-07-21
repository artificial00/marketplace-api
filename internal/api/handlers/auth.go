package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"marketplace-api/internal/models"
	"marketplace-api/internal/service"
	"marketplace-api/pkg/utils"
)

type AuthHandler struct {
	authService service.AuthServiceInterface
}

func NewAuthHandler(authService service.AuthServiceInterface) *AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}

// Register регистрирует нового пользователя
// @Summary Регистрация пользователя
// @Description Регистрирует нового пользователя в системе
// @Tags auth
// @Accept json
// @Produce json
// @Param user body models.RegisterRequest true "Данные пользователя"
// @Success 201 {object} models.AuthResponse
// @Failure 400 {object} utils.ErrorResponse
// @Failure 409 {object} utils.ErrorResponse
// @Router /auth/register [post]
func (h *AuthHandler) Register(c *gin.Context) {
	var req models.RegisterRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "Invalid request format")
		return
	}

	response, err := h.authService.Register(req)
	if err != nil {
		if err.Error() == "user already exists" {
			utils.Conflict(c, "User with this login already exists")
			return
		}
		if err.Error() == "invalid login format" {
			utils.BadRequest(c, "Login must be 3-50 characters long and contain only letters, numbers, and underscores")
			return
		}
		if err.Error() == "password must be at least 6 characters long and contain letters and digits" {
			utils.BadRequest(c, err.Error())
			return
		}

		utils.InternalError(c, "Registration failed")
		return
	}

	utils.SendSuccess(c, http.StatusCreated, response, "User registered successfully")
}

// Login авторизует пользователя
// @Summary Авторизация пользователя
// @Description Авторизует пользователя в системе
// @Tags auth
// @Accept json
// @Produce json
// @Param credentials body models.LoginRequest true "Учетные данные"
// @Success 200 {object} models.AuthResponse
// @Failure 400 {object} utils.ErrorResponse
// @Failure 401 {object} utils.ErrorResponse
// @Router /auth/login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	var req models.LoginRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "Invalid request format")
		return
	}

	response, err := h.authService.Login(req)
	if err != nil {
		utils.Unauthorized(c, "Invalid login or password")
		return
	}

	utils.SendSuccess(c, http.StatusOK, response, "Login successful")
}

// Me возвращает информацию о текущем пользователе
// @Summary Получить информацию о текущем пользователе
// @Description Возвращает информацию о авторизованном пользователе
// @Tags auth
// @Security Bearer
// @Produce json
// @Success 200 {object} models.User
// @Failure 401 {object} utils.ErrorResponse
// @Router /auth/me [get]
func (h *AuthHandler) Me(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		utils.Unauthorized(c, "User not found in context")
		return
	}

	id, ok := userID.(int)
	if !ok {
		utils.InternalError(c, "Invalid user ID format")
		return
	}

	user, err := h.authService.GetUserByID(id)
	if err != nil {
		utils.InternalError(c, "Failed to get user info")
		return
	}

	utils.SendSuccess(c, http.StatusOK, user, "")
}
