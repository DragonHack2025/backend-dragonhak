package handlers

import (
	"backend-dragonhak/models"
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestCreateUser(t *testing.T) {
	// Setup
	gin.SetMode(gin.TestMode)
	SetupTestDB(t)
	defer CleanupTestDB(t)

	router := gin.Default()
	router.POST("/users", CreateUser)

	// Test cases
	tests := []struct {
		name           string
		payload        interface{}
		expectedStatus int
		expectedError  string
	}{
		{
			name: "Valid user creation",
			payload: CreateUserRequest{
				Username: "testuser",
				Email:    "test@example.com",
				Password: "StrongP@ss123",
				Role:     "customer",
			},
			expectedStatus: http.StatusCreated,
		},
		{
			name: "Invalid email format",
			payload: CreateUserRequest{
				Username: "testuser",
				Email:    "invalid-email",
				Password: "StrongP@ss123",
				Role:     "customer",
			},
			expectedStatus: http.StatusBadRequest,
			expectedError:  "Invalid email format",
		},
		{
			name: "Weak password",
			payload: CreateUserRequest{
				Username: "testuser",
				Email:    "test@example.com",
				Password: "weak",
				Role:     "customer",
			},
			expectedStatus: http.StatusBadRequest,
			expectedError:  "Password must be at least 8 characters long",
		},
		{
			name: "Invalid role",
			payload: CreateUserRequest{
				Username: "testuser",
				Email:    "test@example.com",
				Password: "StrongP@ss123",
				Role:     "invalid_role",
			},
			expectedStatus: http.StatusBadRequest,
			expectedError:  "Invalid role",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create request
			jsonPayload, _ := json.Marshal(tt.payload)
			req, _ := http.NewRequest("POST", "/users", bytes.NewBuffer(jsonPayload))
			req.Header.Set("Content-Type", "application/json")

			// Create response recorder
			w := httptest.NewRecorder()

			// Perform request
			router.ServeHTTP(w, req)

			// Check status code
			assert.Equal(t, tt.expectedStatus, w.Code)

			if tt.expectedStatus == http.StatusCreated {
				// Parse response
				var response map[string]interface{}
				json.Unmarshal(w.Body.Bytes(), &response)

				// Check tokens
				assert.NotEmpty(t, response["access_token"])
				assert.NotEmpty(t, response["refresh_token"])

				// Check user object
				user, ok := response["user"].(map[string]interface{})
				assert.True(t, ok, "Response should contain a user object")
				assert.NotEmpty(t, user["id"])
				assert.Equal(t, tt.payload.(CreateUserRequest).Email, user["email"])
				assert.Equal(t, tt.payload.(CreateUserRequest).Role, user["role"])
			} else {
				// Check error message
				var response map[string]string
				json.Unmarshal(w.Body.Bytes(), &response)
				assert.Contains(t, response["error"], tt.expectedError)
			}
		})
	}
}

func TestGetUser(t *testing.T) {
	// Setup
	gin.SetMode(gin.TestMode)
	SetupTestDB(t)
	defer CleanupTestDB(t)

	router := gin.Default()
	router.GET("/users/:id", GetUser)

	// Create a test user
	userID := createTestUser(t)

	tests := []struct {
		name           string
		userID         string
		expectedStatus int
	}{
		{
			name:           "Valid user ID",
			userID:         userID.Hex(),
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Invalid user ID",
			userID:         "invalid-id",
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "Non-existent user ID",
			userID:         primitive.NewObjectID().Hex(),
			expectedStatus: http.StatusNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create request
			req, _ := http.NewRequest("GET", "/users/"+tt.userID, nil)

			// Create response recorder
			w := httptest.NewRecorder()

			// Perform request
			router.ServeHTTP(w, req)

			// Check status code
			assert.Equal(t, tt.expectedStatus, w.Code)

			if tt.expectedStatus == http.StatusOK {
				// Parse response
				var response map[string]interface{}
				json.Unmarshal(w.Body.Bytes(), &response)

				// Check response fields
				assert.Equal(t, tt.userID, response["id"])
				assert.Equal(t, "testuser", response["username"])
				assert.Equal(t, "test@example.com", response["email"])
				assert.Empty(t, response["password"]) // Password should not be returned
			}
		})
	}
}

func TestCreateCraftsmanProfile(t *testing.T) {
	// Setup
	gin.SetMode(gin.TestMode)
	SetupTestDB(t)
	defer CleanupTestDB(t)

	router := gin.Default()
	router.POST("/craftsmen", CreateCraftsmanProfile)

	// Create a test user
	userID := createTestUser(t)

	// Test cases
	tests := []struct {
		name           string
		payload        interface{}
		expectedStatus int
		expectedError  string
	}{
		{
			name: "Valid craftsman profile",
			payload: models.Craftsman{
				UserID:     userID,
				Bio:        "Experienced woodworker",
				Experience: 5,
				Rating:     4.5,
				Location:   "New York",
				ContactInfo: models.ContactInformation{
					Phone:   "+1234567890",
					Website: "https://example.com",
					SocialMedia: map[string]string{
						"instagram": "@woodworker",
					},
				},
				IsVerified: true,
			},
			expectedStatus: http.StatusCreated,
		},
		{
			name: "Invalid user ID",
			payload: models.Craftsman{
				UserID:     primitive.NewObjectID(),
				Bio:        "Experienced woodworker",
				Experience: 5,
				Rating:     4.5,
				Location:   "New York",
				ContactInfo: models.ContactInformation{
					Phone:   "+1234567890",
					Website: "https://example.com",
					SocialMedia: map[string]string{
						"instagram": "@woodworker",
					},
				},
				IsVerified: true,
			},
			expectedStatus: http.StatusBadRequest,
			expectedError:  "User not found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create request
			jsonPayload, _ := json.Marshal(tt.payload)
			req, _ := http.NewRequest("POST", "/craftsmen", bytes.NewBuffer(jsonPayload))
			req.Header.Set("Content-Type", "application/json")

			// Create response recorder
			w := httptest.NewRecorder()

			// Perform request
			router.ServeHTTP(w, req)

			// Check status code
			assert.Equal(t, tt.expectedStatus, w.Code)

			if tt.expectedStatus == http.StatusCreated {
				// Parse response
				var response map[string]interface{}
				json.Unmarshal(w.Body.Bytes(), &response)

				// Check response fields
				assert.NotEmpty(t, response["id"])
				assert.Equal(t, tt.payload.(models.Craftsman).Bio, response["bio"])
				assert.Equal(t, tt.payload.(models.Craftsman).Experience, int(response["experience"].(float64)))
				assert.Equal(t, tt.payload.(models.Craftsman).Rating, response["rating"])
				assert.Equal(t, tt.payload.(models.Craftsman).Location, response["location"])
				assert.Equal(t, tt.payload.(models.Craftsman).IsVerified, response["is_verified"])
			} else {
				// Check error message
				var response map[string]string
				json.Unmarshal(w.Body.Bytes(), &response)
				assert.Contains(t, response["error"], tt.expectedError)
			}
		})
	}
}

func TestCreateCraft(t *testing.T) {
	// Setup mock database
	SetupTestDB(t)
	defer CleanupTestDB(t)
}
