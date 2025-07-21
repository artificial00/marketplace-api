package handlers

import (
	"bytes"
	"errors"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"marketplace-api/internal/models"
	mockservice "marketplace-api/internal/service/mocks"
)

func stringPtr(s string) *string {
	return &s
}

func intPtr(i int) *int {
	return &i
}

func float64Ptr(f float64) *float64 {
	return &f
}

func TestListingHandler_CreateListing(t *testing.T) {
	type mockBehavior func(s *mockservice.MockListingService, userID int, req models.CreateListingRequest)

	testTable := []struct {
		name                 string
		requestBody          string
		userID               interface{}
		request              models.CreateListingRequest
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:        "OK with image URL",
			requestBody: `{"title":"iPhone 15","description":"Brand new iPhone 15 Pro Max","price":120000.50,"image_url":"https://example.com/iphone15.jpg"}`,
			userID:      1,
			request: models.CreateListingRequest{
				Title:       "iPhone 15",
				Description: "Brand new iPhone 15 Pro Max",
				Price:       120000.50,
				ImageURL:    stringPtr("https://example.com/iphone15.jpg"),
			},
			mockBehavior: func(s *mockservice.MockListingService, userID int, req models.CreateListingRequest) {
				listing := &models.Listing{
					ID:          1,
					Title:       "iPhone 15",
					Description: "Brand new iPhone 15 Pro Max",
					Price:       120000.50,
					ImageURL:    stringPtr("https://example.com/iphone15.jpg"),
					UserID:      1,
					CreatedAt:   time.Date(2025, 7, 21, 20, 28, 29, 0, time.UTC),
					UpdatedAt:   time.Date(2025, 7, 21, 20, 28, 29, 0, time.UTC),
				}
				s.EXPECT().CreateListing(userID, req).Return(listing, nil)
			},
			expectedStatusCode:   http.StatusCreated,
			expectedResponseBody: `{"message":"Listing created successfully","data":{"id":1,"title":"iPhone 15","description":"Brand new iPhone 15 Pro Max","price":120000.5,"image_url":"https://example.com/iphone15.jpg","user_id":1,"created_at":"2025-07-21T20:28:29Z","updated_at":"2025-07-21T20:28:29Z"}}`,
		},
		{
			name:        "OK without image URL",
			requestBody: `{"title":"MacBook Pro","description":"Latest MacBook Pro 16 inch","price":250000.00}`,
			userID:      1,
			request: models.CreateListingRequest{
				Title:       "MacBook Pro",
				Description: "Latest MacBook Pro 16 inch",
				Price:       250000.00,
				ImageURL:    nil,
			},
			mockBehavior: func(s *mockservice.MockListingService, userID int, req models.CreateListingRequest) {
				listing := &models.Listing{
					ID:          2,
					Title:       "MacBook Pro",
					Description: "Latest MacBook Pro 16 inch",
					Price:       250000.00,
					ImageURL:    nil,
					UserID:      1,
					CreatedAt:   time.Date(2025, 7, 21, 20, 28, 29, 0, time.UTC),
					UpdatedAt:   time.Date(2025, 7, 21, 20, 28, 29, 0, time.UTC),
				}
				s.EXPECT().CreateListing(userID, req).Return(listing, nil)
			},
			expectedStatusCode:   http.StatusCreated,
			expectedResponseBody: `{"message":"Listing created successfully","data":{"id":2,"title":"MacBook Pro","description":"Latest MacBook Pro 16 inch","price":250000,"image_url":null,"user_id":1,"created_at":"2025-07-21T20:28:29Z","updated_at":"2025-07-21T20:28:29Z"}}`,
		},
		{
			name:                 "User not found in context",
			requestBody:          `{"title":"iPhone 15","description":"Brand new iPhone 15","price":120000.00}`,
			expectedStatusCode:   http.StatusUnauthorized,
			expectedResponseBody: `{"error":"unauthorized", "message":"User not found in context"}`,
		},
		{
			name:                 "Invalid request format - malformed JSON",
			requestBody:          `{"title":"iPhone 15","description":"Brand new iPhone 15","price":}`,
			userID:               1,
			expectedStatusCode:   http.StatusBadRequest,
			expectedResponseBody: `{"error":"bad_request", "message":"Invalid request format: invalid character '}' looking for beginning of value"}`,
		},
		{
			name:                 "Missing title field",
			requestBody:          `{"description":"Brand new iPhone 15","price":120000.00}`,
			userID:               1,
			expectedStatusCode:   http.StatusBadRequest,
			expectedResponseBody: `{"error":"bad_request", "message":"Invalid request format: Key: 'CreateListingRequest.Title' Error:Field validation for 'Title' failed on the 'required' tag"}`,
		},
		{
			name:        "Internal server error",
			requestBody: `{"title":"iPhone 15","description":"Brand new iPhone 15","price":120000.00}`,
			userID:      1,
			request: models.CreateListingRequest{
				Title:       "iPhone 15",
				Description: "Brand new iPhone 15",
				Price:       120000.00,
			},
			mockBehavior: func(s *mockservice.MockListingService, userID int, req models.CreateListingRequest) {
				s.EXPECT().CreateListing(userID, req).Return(nil, errors.New("database connection failed"))
			},
			expectedStatusCode:   http.StatusInternalServerError,
			expectedResponseBody: `{"error":"internal_error", "message":"Failed to create listing"}`,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			listingService := mockservice.NewMockListingService(c)

			if testCase.mockBehavior != nil {
				if userID, ok := testCase.userID.(int); ok {
					testCase.mockBehavior(listingService, userID, testCase.request)
				}
			}

			handler := NewListingHandler(listingService)

			w := httptest.NewRecorder()
			gin.SetMode(gin.TestMode)
			ctx, r := gin.CreateTestContext(w)

			r.Use(func(ctx *gin.Context) {
				if testCase.userID != nil {
					ctx.Set("user_id", testCase.userID)
				}
			})

			r.POST("/listings", handler.CreateListing)

			ctx.Request, _ = http.NewRequest("POST", "/listings", bytes.NewBufferString(testCase.requestBody))
			ctx.Request.Header.Set("Content-Type", "application/json")

			r.ServeHTTP(w, ctx.Request)

			assert.Equal(t, testCase.expectedStatusCode, w.Code)
			assert.JSONEq(t, testCase.expectedResponseBody, w.Body.String())
		})
	}
}

func TestListingHandler_GetListing(t *testing.T) {
	type mockBehavior func(s *mockservice.MockListingService, id int, currentUserID *int)

	testTable := []struct {
		name                 string
		listingID            string
		userID               interface{}
		currentUserID        *int
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:          "OK without user",
			listingID:     "1",
			currentUserID: nil,
			mockBehavior: func(s *mockservice.MockListingService, id int, currentUserID *int) {
				listing := &models.Listing{
					ID:          1,
					Title:       "iPhone 15",
					Description: "Brand new iPhone 15 Pro Max",
					Price:       120000.00,
					ImageURL:    stringPtr("https://example.com/iphone15.jpg"),
					UserID:      1,
					CreatedAt:   time.Date(2025, 7, 21, 20, 28, 29, 0, time.UTC),
					UpdatedAt:   time.Date(2025, 7, 21, 20, 28, 29, 0, time.UTC),
				}
				s.EXPECT().GetListingByID(id, currentUserID).Return(listing, nil)
			},
			expectedStatusCode:   http.StatusOK,
			expectedResponseBody: `{"data":{"id":1,"title":"iPhone 15","description":"Brand new iPhone 15 Pro Max","price":120000,"image_url":"https://example.com/iphone15.jpg","user_id":1,"created_at":"2025-07-21T20:28:29Z","updated_at":"2025-07-21T20:28:29Z"}}`,
		},
		{
			name:          "OK with user",
			listingID:     "2",
			userID:        1,
			currentUserID: intPtr(1),
			mockBehavior: func(s *mockservice.MockListingService, id int, currentUserID *int) {
				listing := &models.Listing{
					ID:          2,
					Title:       "MacBook Pro",
					Description: "16-inch MacBook Pro with M2 chip",
					Price:       250000.00,
					ImageURL:    nil,
					UserID:      2,
					CreatedAt:   time.Date(2025, 7, 21, 20, 28, 29, 0, time.UTC),
					UpdatedAt:   time.Date(2025, 7, 21, 20, 28, 29, 0, time.UTC),
				}
				s.EXPECT().GetListingByID(id, currentUserID).Return(listing, nil)
			},
			expectedStatusCode:   http.StatusOK,
			expectedResponseBody: `{"data":{"id":2,"title":"MacBook Pro","description":"16-inch MacBook Pro with M2 chip","price":250000,"image_url":null,"user_id":2,"created_at":"2025-07-21T20:28:29Z","updated_at":"2025-07-21T20:28:29Z"}}`,
		},
		{
			name:                 "Invalid listing ID - non-numeric",
			listingID:            "abc",
			expectedStatusCode:   http.StatusBadRequest,
			expectedResponseBody: `{"error":"bad_request", "message":"Invalid listing ID"}`,
		},
		{
			name:      "Listing not found",
			listingID: "999",
			mockBehavior: func(s *mockservice.MockListingService, id int, currentUserID *int) {
				s.EXPECT().GetListingByID(id, currentUserID).Return(nil, errors.New("listing not found"))
			},
			expectedStatusCode:   http.StatusNotFound,
			expectedResponseBody: `{"error":"not_found", "message":"Listing not found"}`,
		},
		{
			name:      "Invalid listing ID from service",
			listingID: "-1",
			mockBehavior: func(s *mockservice.MockListingService, id int, currentUserID *int) {
				s.EXPECT().GetListingByID(id, currentUserID).Return(nil, errors.New("invalid listing ID"))
			},
			expectedStatusCode:   http.StatusBadRequest,
			expectedResponseBody: `{"error":"bad_request", "message":"invalid listing ID"}`,
		},
		{
			name:      "Zero listing ID",
			listingID: "0",
			mockBehavior: func(s *mockservice.MockListingService, id int, currentUserID *int) {
				s.EXPECT().GetListingByID(id, currentUserID).Return(nil, errors.New("invalid listing ID"))
			},
			expectedStatusCode:   http.StatusBadRequest,
			expectedResponseBody: `{"error":"bad_request", "message":"invalid listing ID"}`,
		},
		{
			name:      "Internal server error",
			listingID: "1",
			mockBehavior: func(s *mockservice.MockListingService, id int, currentUserID *int) {
				s.EXPECT().GetListingByID(id, currentUserID).Return(nil, errors.New("database connection failed"))
			},
			expectedStatusCode:   http.StatusInternalServerError,
			expectedResponseBody: `{"error":"internal_error", "message":"Failed to get listing"}`,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			listingService := mockservice.NewMockListingService(c)

			if testCase.mockBehavior != nil {
				if id, err := strconv.Atoi(testCase.listingID); err == nil {
					testCase.mockBehavior(listingService, id, testCase.currentUserID)
				}
			}

			handler := NewListingHandler(listingService)

			w := httptest.NewRecorder()
			gin.SetMode(gin.TestMode)
			ctx, r := gin.CreateTestContext(w)

			r.Use(func(ctx *gin.Context) {
				if testCase.userID != nil {
					ctx.Set("user_id", testCase.userID)
				}
			})

			r.GET("/listings/:id", handler.GetListing)

			url := "/listings/" + testCase.listingID
			ctx.Request, _ = http.NewRequest("GET", url, nil)

			r.ServeHTTP(w, ctx.Request)

			assert.Equal(t, testCase.expectedStatusCode, w.Code)
			assert.JSONEq(t, testCase.expectedResponseBody, w.Body.String())
		})
	}
}

func TestListingHandler_UpdateListing(t *testing.T) {
	type mockBehavior func(s *mockservice.MockListingService, id int, userID int, req models.UpdateListingRequest)

	testTable := []struct {
		name                 string
		listingID            string
		requestBody          string
		userID               interface{}
		request              models.UpdateListingRequest
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:        "OK - full update",
			listingID:   "1",
			requestBody: `{"title":"iPhone 15 Pro Updated","description":"Updated description","price":130000.00,"image_url":"https://example.com/updated.jpg"}`,
			userID:      1,
			request: models.UpdateListingRequest{
				Title:       stringPtr("iPhone 15 Pro Updated"),
				Description: stringPtr("Updated description"),
				Price:       float64Ptr(130000.00),
				ImageURL:    stringPtr("https://example.com/updated.jpg"),
			},
			mockBehavior: func(s *mockservice.MockListingService, id int, userID int, req models.UpdateListingRequest) {
				listing := &models.Listing{
					ID:          1,
					Title:       "iPhone 15 Pro Updated",
					Description: "Updated description",
					Price:       130000.00,
					ImageURL:    stringPtr("https://example.com/updated.jpg"),
					UserID:      1,
					CreatedAt:   time.Date(2025, 7, 21, 20, 28, 29, 0, time.UTC),
					UpdatedAt:   time.Date(2025, 7, 21, 20, 30, 0, 0, time.UTC),
				}
				s.EXPECT().UpdateListing(id, userID, req).Return(listing, nil)
			},
			expectedStatusCode:   http.StatusOK,
			expectedResponseBody: `{"message":"Listing updated successfully","data":{"id":1,"title":"iPhone 15 Pro Updated","description":"Updated description","price":130000,"image_url":"https://example.com/updated.jpg","user_id":1,"created_at":"2025-07-21T20:28:29Z","updated_at":"2025-07-21T20:30:00Z"}}`,
		},
		{
			name:        "OK - partial update",
			listingID:   "2",
			requestBody: `{"price":140000.00}`,
			userID:      1,
			request: models.UpdateListingRequest{
				Price: float64Ptr(140000.00),
			},
			mockBehavior: func(s *mockservice.MockListingService, id int, userID int, req models.UpdateListingRequest) {
				listing := &models.Listing{
					ID:          2,
					Title:       "MacBook Pro",
					Description: "16-inch MacBook Pro",
					Price:       140000.00,
					ImageURL:    nil,
					UserID:      1,
					CreatedAt:   time.Date(2025, 7, 21, 20, 28, 29, 0, time.UTC),
					UpdatedAt:   time.Date(2025, 7, 21, 20, 30, 0, 0, time.UTC),
				}
				s.EXPECT().UpdateListing(id, userID, req).Return(listing, nil)
			},
			expectedStatusCode:   http.StatusOK,
			expectedResponseBody: `{"message":"Listing updated successfully","data":{"id":2,"title":"MacBook Pro","description":"16-inch MacBook Pro","price":140000,"image_url":null,"user_id":1,"created_at":"2025-07-21T20:28:29Z","updated_at":"2025-07-21T20:30:00Z"}}`,
		},
		{
			name:                 "User not found in context",
			listingID:            "1",
			requestBody:          `{"title":"Updated title"}`,
			expectedStatusCode:   http.StatusUnauthorized,
			expectedResponseBody: `{"error":"unauthorized", "message":"User not found in context"}`,
		},
		{
			name:                 "Invalid listing ID - non-numeric",
			listingID:            "abc",
			requestBody:          `{"title":"Updated title"}`,
			userID:               1,
			expectedStatusCode:   http.StatusBadRequest,
			expectedResponseBody: `{"error":"bad_request", "message":"Invalid listing ID"}`,
		},
		{
			name:                 "Invalid request format - malformed JSON",
			listingID:            "1",
			requestBody:          `{"title":"Updated title","price":}`,
			userID:               1,
			expectedStatusCode:   http.StatusBadRequest,
			expectedResponseBody: `{"error":"bad_request", "message":"Invalid request format: invalid character '}' looking for beginning of value"}`,
		},
		{
			name:        "Listing not found",
			listingID:   "999",
			requestBody: `{"title":"Updated title"}`,
			userID:      1,
			request: models.UpdateListingRequest{
				Title: stringPtr("Updated title"),
			},
			mockBehavior: func(s *mockservice.MockListingService, id int, userID int, req models.UpdateListingRequest) {
				s.EXPECT().UpdateListing(id, userID, req).Return(nil, errors.New("listing not found"))
			},
			expectedStatusCode:   http.StatusNotFound,
			expectedResponseBody: `{"error":"not_found", "message":"Listing not found"}`,
		},
		{
			name:        "Access denied - not owner",
			listingID:   "1",
			requestBody: `{"title":"Updated title"}`,
			userID:      2,
			request: models.UpdateListingRequest{
				Title: stringPtr("Updated title"),
			},
			mockBehavior: func(s *mockservice.MockListingService, id int, userID int, req models.UpdateListingRequest) {
				s.EXPECT().UpdateListing(id, userID, req).Return(nil, errors.New("access denied: you can only edit your own listings"))
			},
			expectedStatusCode:   http.StatusForbidden,
			expectedResponseBody: `{"error":"forbidden", "message":"You can only edit your own listings"}`,
		},
		{
			name:        "Invalid listing ID from service",
			listingID:   "-1",
			requestBody: `{"title":"Updated title"}`,
			userID:      1,
			request: models.UpdateListingRequest{
				Title: stringPtr("Updated title"),
			},
			mockBehavior: func(s *mockservice.MockListingService, id int, userID int, req models.UpdateListingRequest) {
				s.EXPECT().UpdateListing(id, userID, req).Return(nil, errors.New("invalid listing ID"))
			},
			expectedStatusCode:   http.StatusBadRequest,
			expectedResponseBody: `{"error":"bad_request", "message":"invalid listing ID"}`,
		},
		{
			name:        "No fields to update",
			listingID:   "1",
			requestBody: `{}`,
			userID:      1,
			request:     models.UpdateListingRequest{},
			mockBehavior: func(s *mockservice.MockListingService, id int, userID int, req models.UpdateListingRequest) {
				s.EXPECT().UpdateListing(id, userID, req).Return(nil, errors.New("no fields to update"))
			},
			expectedStatusCode:   http.StatusBadRequest,
			expectedResponseBody: `{"error":"bad_request", "message":"no fields to update"}`,
		},
		{
			name:        "Internal server error",
			listingID:   "1",
			requestBody: `{"title":"Updated title"}`,
			userID:      1,
			request: models.UpdateListingRequest{
				Title: stringPtr("Updated title"),
			},
			mockBehavior: func(s *mockservice.MockListingService, id int, userID int, req models.UpdateListingRequest) {
				s.EXPECT().UpdateListing(id, userID, req).Return(nil, errors.New("database connection failed"))
			},
			expectedStatusCode:   http.StatusInternalServerError,
			expectedResponseBody: `{"error":"internal_error", "message":"Failed to update listing"}`,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			listingService := mockservice.NewMockListingService(c)

			if testCase.mockBehavior != nil {
				if id, err := strconv.Atoi(testCase.listingID); err == nil {
					if userID, ok := testCase.userID.(int); ok {
						testCase.mockBehavior(listingService, id, userID, testCase.request)
					}
				}
			}

			handler := NewListingHandler(listingService)

			w := httptest.NewRecorder()
			gin.SetMode(gin.TestMode)
			ctx, r := gin.CreateTestContext(w)

			r.Use(func(ctx *gin.Context) {
				if testCase.userID != nil {
					ctx.Set("user_id", testCase.userID)
				}
			})

			r.PUT("/listings/:id", handler.UpdateListing)

			url := "/listings/" + testCase.listingID
			ctx.Request, _ = http.NewRequest("PUT", url, bytes.NewBufferString(testCase.requestBody))
			ctx.Request.Header.Set("Content-Type", "application/json")

			r.ServeHTTP(w, ctx.Request)

			assert.Equal(t, testCase.expectedStatusCode, w.Code)
			assert.JSONEq(t, testCase.expectedResponseBody, w.Body.String())
		})
	}
}

func TestListingHandler_DeleteListing(t *testing.T) {
	type mockBehavior func(s *mockservice.MockListingService, id int, userID int)

	testTable := []struct {
		name                 string
		listingID            string
		userID               interface{}
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:      "OK",
			listingID: "1",
			userID:    1,
			mockBehavior: func(s *mockservice.MockListingService, id int, userID int) {
				s.EXPECT().DeleteListing(id, userID).Return(nil)
			},
			expectedStatusCode:   http.StatusOK,
			expectedResponseBody: `{"message":"Listing deleted successfully"}`,
		},
		{
			name:                 "User not found in context",
			listingID:            "1",
			expectedStatusCode:   http.StatusUnauthorized,
			expectedResponseBody: `{"error":"unauthorized", "message":"User not found in context"}`,
		},
		{
			name:                 "Invalid listing ID - non-numeric",
			listingID:            "abc",
			userID:               1,
			expectedStatusCode:   http.StatusBadRequest,
			expectedResponseBody: `{"error":"bad_request", "message":"Invalid listing ID"}`,
		},
		{
			name:      "Listing not found",
			listingID: "999",
			userID:    1,
			mockBehavior: func(s *mockservice.MockListingService, id int, userID int) {
				s.EXPECT().DeleteListing(id, userID).Return(errors.New("listing not found"))
			},
			expectedStatusCode:   http.StatusNotFound,
			expectedResponseBody: `{"error":"not_found", "message":"Listing not found"}`,
		},
		{
			name:      "Access denied - not owner",
			listingID: "1",
			userID:    2,
			mockBehavior: func(s *mockservice.MockListingService, id int, userID int) {
				s.EXPECT().DeleteListing(id, userID).Return(errors.New("access denied: you can only delete your own listings"))
			},
			expectedStatusCode:   http.StatusForbidden,
			expectedResponseBody: `{"error":"forbidden", "message":"You can only delete your own listings"}`,
		},
		{
			name:      "Invalid listing ID from service",
			listingID: "-1",
			userID:    1,
			mockBehavior: func(s *mockservice.MockListingService, id int, userID int) {
				s.EXPECT().DeleteListing(id, userID).Return(errors.New("invalid listing ID"))
			},
			expectedStatusCode:   http.StatusBadRequest,
			expectedResponseBody: `{"error":"bad_request", "message":"invalid listing ID"}`,
		},
		{
			name:      "Zero listing ID",
			listingID: "0",
			userID:    1,
			mockBehavior: func(s *mockservice.MockListingService, id int, userID int) {
				s.EXPECT().DeleteListing(id, userID).Return(errors.New("invalid listing ID"))
			},
			expectedStatusCode:   http.StatusBadRequest,
			expectedResponseBody: `{"error":"bad_request", "message":"invalid listing ID"}`,
		},
		{
			name:      "Internal server error",
			listingID: "1",
			userID:    1,
			mockBehavior: func(s *mockservice.MockListingService, id int, userID int) {
				s.EXPECT().DeleteListing(id, userID).Return(errors.New("database connection failed"))
			},
			expectedStatusCode:   http.StatusInternalServerError,
			expectedResponseBody: `{"error":"internal_error", "message":"Failed to delete listing"}`,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			listingService := mockservice.NewMockListingService(c)

			if testCase.mockBehavior != nil {
				if id, err := strconv.Atoi(testCase.listingID); err == nil {
					if userID, ok := testCase.userID.(int); ok {
						testCase.mockBehavior(listingService, id, userID)
					}
				}
			}

			handler := NewListingHandler(listingService)

			w := httptest.NewRecorder()
			gin.SetMode(gin.TestMode)
			ctx, r := gin.CreateTestContext(w)

			r.Use(func(ctx *gin.Context) {
				if testCase.userID != nil {
					ctx.Set("user_id", testCase.userID)
				}
			})

			r.DELETE("/listings/:id", handler.DeleteListing)

			url := "/listings/" + testCase.listingID
			ctx.Request, _ = http.NewRequest("DELETE", url, nil)

			r.ServeHTTP(w, ctx.Request)

			assert.Equal(t, testCase.expectedStatusCode, w.Code)
			assert.JSONEq(t, testCase.expectedResponseBody, w.Body.String())
		})
	}
}
