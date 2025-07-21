package utils

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message,omitempty"`
}

type SuccessResponse struct {
	Data    interface{} `json:"data,omitempty"`
	Message string      `json:"message,omitempty"`
}

// SendError отправляет ошибку в JSON формате
func SendError(c *gin.Context, statusCode int, err string, message ...string) {
	response := ErrorResponse{
		Error: err,
	}

	if len(message) > 0 {
		response.Message = message[0]
	}

	c.JSON(statusCode, response)
}

// SendSuccess отправляет успешный ответ в JSON формате
func SendSuccess(c *gin.Context, statusCode int, data interface{}, message ...string) {
	response := SuccessResponse{
		Data: data,
	}

	if len(message) > 0 {
		response.Message = message[0]
	}

	c.JSON(statusCode, response)
}

// BadRequest отправляет ошибку 400
func BadRequest(c *gin.Context, message string) {
	SendError(c, http.StatusBadRequest, "bad_request", message)
}

// Unauthorized отправляет ошибку 401
func Unauthorized(c *gin.Context, message string) {
	SendError(c, http.StatusUnauthorized, "unauthorized", message)
}

// InternalError отправляет ошибку 500
func InternalError(c *gin.Context, message string) {
	SendError(c, http.StatusInternalServerError, "internal_error", message)
}

// Conflict отправляет ошибку 409
func Conflict(c *gin.Context, message string) {
	SendError(c, http.StatusConflict, "conflict", message)
}

func NotFound(c *gin.Context, message string) {
	SendError(c, http.StatusNotFound, "not_found", message)
}

func Forbidden(c *gin.Context, message string) {
	SendError(c, http.StatusForbidden, "forbidden", message)
}
