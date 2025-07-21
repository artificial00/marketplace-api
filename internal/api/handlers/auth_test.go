package handlers

import (
	"bytes"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"marketplace-api/internal/models"
	mockservice "marketplace-api/internal/service/mocks"
)

func TestAuthHandler_Register(t *testing.T) {
	type mockBehavior func(s *mockservice.MockAuthService, req models.RegisterRequest)

	testTable := []struct {
		name                 string
		requestBody          string
		request              models.RegisterRequest
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:        "OK",
			requestBody: `{"login":"artificial00","password":"password123"}`,
			request: models.RegisterRequest{
				Login:    "artificial00",
				Password: "password123",
			},
			mockBehavior: func(s *mockservice.MockAuthService, req models.RegisterRequest) {
				response := &models.AuthResponse{
					User: models.User{
						ID:        1,
						Login:     "artificial00",
						CreatedAt: time.Date(2025, 7, 21, 19, 56, 37, 0, time.UTC),
						UpdatedAt: time.Date(2025, 7, 21, 19, 56, 37, 0, time.UTC),
					},
					Token: "jwt.token.here",
				}
				s.EXPECT().Register(req).Return(response, nil)
			},
			expectedStatusCode:   http.StatusCreated,
			expectedResponseBody: `{"message":"User registered successfully","data":{"user":{"id":1,"login":"artificial00","created_at":"2025-07-21T19:56:37Z","updated_at":"2025-07-21T19:56:37Z"},"token":"jwt.token.here"}}`,
		},
		{
			name:                 "Invalid request format",
			requestBody:          `{"login":"artificial00","password":}`,
			expectedStatusCode:   http.StatusBadRequest,
			expectedResponseBody: `{"error":"bad_request", "message":"Invalid request format"}`,
		},
		{
			name:                 "Missing login field",
			requestBody:          `{"password":"password123"}`,
			expectedStatusCode:   http.StatusBadRequest,
			expectedResponseBody: `{"error":"bad_request", "message":"Invalid request format"}`,
		},
		{
			name:                 "Missing password field",
			requestBody:          `{"login":"artificial00"}`,
			expectedStatusCode:   http.StatusBadRequest,
			expectedResponseBody: `{"error":"bad_request", "message":"Invalid request format"}`,
		},
		{
			name:        "User already exists",
			requestBody: `{"login":"existing_user","password":"password123"}`,
			request: models.RegisterRequest{
				Login:    "existing_user",
				Password: "password123",
			},
			mockBehavior: func(s *mockservice.MockAuthService, req models.RegisterRequest) {
				s.EXPECT().Register(req).Return(nil, errors.New("user already exists"))
			},
			expectedStatusCode:   http.StatusConflict,
			expectedResponseBody: `{"error":"conflict", "message":"User with this login already exists"}`,
		},
		{
			name:        "Internal server error",
			requestBody: `{"login":"testuser","password":"password123"}`,
			request: models.RegisterRequest{
				Login:    "testuser",
				Password: "password123",
			},
			mockBehavior: func(s *mockservice.MockAuthService, req models.RegisterRequest) {
				s.EXPECT().Register(req).Return(nil, errors.New("database connection failed"))
			},
			expectedStatusCode:   http.StatusInternalServerError,
			expectedResponseBody: `{"error":"internal_error", "message":"Registration failed"}`,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			authService := mockservice.NewMockAuthService(c)

			if testCase.mockBehavior != nil {
				testCase.mockBehavior(authService, testCase.request)
			}

			handler := NewAuthHandler(authService)

			w := httptest.NewRecorder()

			gin.SetMode(gin.TestMode)
			ctx, r := gin.CreateTestContext(w)

			r.POST("/auth/register", handler.Register)

			ctx.Request, _ = http.NewRequest("POST", "/auth/register", bytes.NewBufferString(testCase.requestBody))
			ctx.Request.Header.Set("Content-Type", "application/json")

			r.ServeHTTP(w, ctx.Request)

			assert.Equal(t, testCase.expectedStatusCode, w.Code)
			assert.JSONEq(t, testCase.expectedResponseBody, w.Body.String())
		})
	}
}

func TestAuthHandler_Login(t *testing.T) {
	type mockBehavior func(s *mockservice.MockAuthService, req models.LoginRequest)

	testTable := []struct {
		name                 string
		requestBody          string
		request              models.LoginRequest
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:        "OK",
			requestBody: `{"login":"artificial00","password":"password123"}`,
			request: models.LoginRequest{
				Login:    "artificial00",
				Password: "password123",
			},
			mockBehavior: func(s *mockservice.MockAuthService, req models.LoginRequest) {
				response := &models.AuthResponse{
					User: models.User{
						ID:        1,
						Login:     "artificial00",
						CreatedAt: time.Date(2025, 7, 21, 19, 56, 37, 0, time.UTC),
						UpdatedAt: time.Date(2025, 7, 21, 19, 56, 37, 0, time.UTC),
					},
					Token: "jwt.token.here",
				}
				s.EXPECT().Login(req).Return(response, nil)
			},
			expectedStatusCode:   http.StatusOK,
			expectedResponseBody: `{"message":"Login successful","data":{"user":{"id":1,"login":"artificial00","created_at":"2025-07-21T19:56:37Z","updated_at":"2025-07-21T19:56:37Z"},"token":"jwt.token.here"}}`,
		},
		{
			name:                 "Invalid request format",
			requestBody:          `{"login":"artificial00","password":}`,
			expectedStatusCode:   http.StatusBadRequest,
			expectedResponseBody: `{"error":"bad_request", "message":"Invalid request format"}`,
		},
		{
			name:                 "Missing login field",
			requestBody:          `{"password":"password123"}`,
			expectedStatusCode:   http.StatusBadRequest,
			expectedResponseBody: `{"error":"bad_request", "message":"Invalid request format"}`,
		},
		{
			name:                 "Missing password field",
			requestBody:          `{"login":"artificial00"}`,
			expectedStatusCode:   http.StatusBadRequest,
			expectedResponseBody: `{"error":"bad_request", "message":"Invalid request format"}`,
		},
		{
			name:        "Invalid credentials",
			requestBody: `{"login":"artificial00","password":"wrongpassword"}`,
			request: models.LoginRequest{
				Login:    "artificial00",
				Password: "wrongpassword",
			},
			mockBehavior: func(s *mockservice.MockAuthService, req models.LoginRequest) {
				s.EXPECT().Login(req).Return(nil, errors.New("invalid credentials"))
			},
			expectedStatusCode:   http.StatusUnauthorized,
			expectedResponseBody: `{"error":"unauthorized", "message":"Invalid login or password"}`,
		},
		{
			name:        "User not found",
			requestBody: `{"login":"nonexistent","password":"password123"}`,
			request: models.LoginRequest{
				Login:    "nonexistent",
				Password: "password123",
			},
			mockBehavior: func(s *mockservice.MockAuthService, req models.LoginRequest) {
				s.EXPECT().Login(req).Return(nil, errors.New("user not found"))
			},
			expectedStatusCode:   http.StatusUnauthorized,
			expectedResponseBody: `{"error":"unauthorized", "message":"Invalid login or password"}`,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			authService := mockservice.NewMockAuthService(c)

			if testCase.mockBehavior != nil {
				testCase.mockBehavior(authService, testCase.request)
			}

			handler := NewAuthHandler(authService)

			w := httptest.NewRecorder()

			gin.SetMode(gin.TestMode)
			ctx, r := gin.CreateTestContext(w)

			r.POST("/auth/login", handler.Login)

			ctx.Request, _ = http.NewRequest("POST", "/auth/login", bytes.NewBufferString(testCase.requestBody))
			ctx.Request.Header.Set("Content-Type", "application/json")

			r.ServeHTTP(w, ctx.Request)

			assert.Equal(t, testCase.expectedStatusCode, w.Code)
			assert.JSONEq(t, testCase.expectedResponseBody, w.Body.String())
		})
	}
}

func TestAuthHandler_Me(t *testing.T) {
	type mockBehavior func(s *mockservice.MockAuthService, userID int)

	testTable := []struct {
		name                 string
		userID               interface{}
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:   "OK",
			userID: 1,
			mockBehavior: func(s *mockservice.MockAuthService, userID int) {
				user := &models.User{
					ID:        1,
					Login:     "artificial00",
					CreatedAt: time.Date(2025, 7, 21, 19, 56, 37, 0, time.UTC),
					UpdatedAt: time.Date(2025, 7, 21, 19, 56, 37, 0, time.UTC),
				}
				s.EXPECT().GetUserByID(userID).Return(user, nil)
			},
			expectedStatusCode:   http.StatusOK,
			expectedResponseBody: `{"data":{"id":1,"login":"artificial00","created_at":"2025-07-21T19:56:37Z","updated_at":"2025-07-21T19:56:37Z"}}`,
		},
		{
			name:                 "User not found in context",
			expectedStatusCode:   http.StatusUnauthorized,
			expectedResponseBody: `{"error":"unauthorized", "message":"User not found in context"}`,
		},
		{
			name:                 "Invalid user ID format (string)",
			userID:               "invalid",
			expectedStatusCode:   http.StatusInternalServerError,
			expectedResponseBody: `{"error":"internal_error", "message":"Invalid user ID format"}`,
		},
		{
			name:                 "Invalid user ID format (float)",
			userID:               1.5,
			expectedStatusCode:   http.StatusInternalServerError,
			expectedResponseBody: `{"error":"internal_error", "message":"Invalid user ID format"}`,
		},
		{
			name:   "Service error - user not found",
			userID: 999,
			mockBehavior: func(s *mockservice.MockAuthService, userID int) {
				s.EXPECT().GetUserByID(userID).Return(nil, errors.New("user not found"))
			},
			expectedStatusCode:   http.StatusInternalServerError,
			expectedResponseBody: `{"error":"internal_error", "message":"Failed to get user info"}`,
		},
		{
			name:   "Service error - database error",
			userID: 1,
			mockBehavior: func(s *mockservice.MockAuthService, userID int) {
				s.EXPECT().GetUserByID(userID).Return(nil, errors.New("database connection failed"))
			},
			expectedStatusCode:   http.StatusInternalServerError,
			expectedResponseBody: `{"error":"internal_error", "message":"Failed to get user info"}`,
		},
		{
			name:   "Zero user ID",
			userID: 0,
			mockBehavior: func(s *mockservice.MockAuthService, userID int) {
				s.EXPECT().GetUserByID(userID).Return(nil, errors.New("invalid user id"))
			},
			expectedStatusCode:   http.StatusInternalServerError,
			expectedResponseBody: `{"error":"internal_error", "message":"Failed to get user info"}`,
		},
		{
			name:   "Negative user ID",
			userID: -1,
			mockBehavior: func(s *mockservice.MockAuthService, userID int) {
				s.EXPECT().GetUserByID(userID).Return(nil, errors.New("invalid user id"))
			},
			expectedStatusCode:   http.StatusInternalServerError,
			expectedResponseBody: `{"error":"internal_error", "message":"Failed to get user info"}`,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			authService := mockservice.NewMockAuthService(c)

			if testCase.mockBehavior != nil {
				if userID, ok := testCase.userID.(int); ok {
					testCase.mockBehavior(authService, userID)
				}
			}

			handler := NewAuthHandler(authService)

			w := httptest.NewRecorder()

			gin.SetMode(gin.TestMode)
			ctx, r := gin.CreateTestContext(w)

			r.Use(func(ctx *gin.Context) {
				if testCase.userID != nil {
					ctx.Set("user_id", testCase.userID)
				}
			})

			r.GET("/auth/me", handler.Me)

			ctx.Request, _ = http.NewRequest("GET", "/auth/me", nil)

			r.ServeHTTP(w, ctx.Request)

			assert.Equal(t, testCase.expectedStatusCode, w.Code)
			assert.JSONEq(t, testCase.expectedResponseBody, w.Body.String())
		})
	}
}
