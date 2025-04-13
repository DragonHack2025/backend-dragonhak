package handlers

import (
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

	// Test cases
	tests := []struct {
		name           string
		payload        interface{}
		expectedStatus int
		expectedError  string
	}{
		{
			name: "Valid craftsman profile",
			payload: map[string]interface{}{
				"username":   "testcraftsman",
				"email":      "craftsman@example.com",
				"password":   "StrongP@ss123",
				"bio":        "Experienced woodworker",
				"experience": 5,
				"rating":     4.5,
				"location":   "New York",
				"contact_info": map[string]interface{}{
					"phone":   "+1234567890",
					"website": "https://example.com",
					"social_media": map[string]string{
						"instagram": "@woodworker",
					},
				},
				"is_verified": true,
			},
			expectedStatus: http.StatusCreated,
		},
		{
			name: "Invalid email format",
			payload: map[string]interface{}{
				"username":   "testcraftsman",
				"email":      "invalid-email",
				"password":   "StrongP@ss123",
				"bio":        "Experienced woodworker",
				"experience": 5,
				"rating":     4.5,
				"location":   "New York",
				"contact_info": map[string]interface{}{
					"phone":   "+1234567890",
					"website": "https://example.com",
				},
				"is_verified": true,
			},
			expectedStatus: http.StatusBadRequest,
			expectedError:  "Field validation for 'Email' failed",
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

				// Check tokens
				assert.NotEmpty(t, response["access_token"])
				assert.NotEmpty(t, response["refresh_token"])

				// Check user object
				user := response["user"].(map[string]interface{})
				assert.NotEmpty(t, user["id"])
				assert.Equal(t, tt.payload.(map[string]interface{})["email"], user["email"])
				assert.Equal(t, "craftsman", user["role"])

				// Check craftsman object
				craftsman := response["craftsman"].(map[string]interface{})
				assert.NotEmpty(t, craftsman["id"])
				assert.Equal(t, tt.payload.(map[string]interface{})["bio"], craftsman["bio"])
				assert.Equal(t, tt.payload.(map[string]interface{})["experience"], int(craftsman["experience"].(float64)))
				assert.Equal(t, tt.payload.(map[string]interface{})["rating"], craftsman["rating"])
				assert.Equal(t, tt.payload.(map[string]interface{})["location"], craftsman["location"])
				assert.Equal(t, tt.payload.(map[string]interface{})["is_verified"], craftsman["is_verified"])

				// Check contact info
				contactInfo := craftsman["contact_info"].(map[string]interface{})
				expectedContactInfo := tt.payload.(map[string]interface{})["contact_info"].(map[string]interface{})
				assert.Equal(t, expectedContactInfo["phone"], contactInfo["phone"])
				assert.Equal(t, expectedContactInfo["website"], contactInfo["website"])
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
