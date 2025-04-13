package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"backend-dragonhak/models"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestCreateCraftsmanProfileHandler(t *testing.T) {
	// Setup mock database
	SetupTestDB(t)
	defer CleanupTestDB(t)

	// Create test router
	router := gin.Default()
	router.POST("/api/craftsmen", CreateCraftsmanProfile)

	// Test cases
	tests := []struct {
		name           string
		payload        map[string]interface{}
		expectedStatus int
	}{
		{
			name: "Valid profile",
			payload: map[string]interface{}{
				"username":   "testcraftsman",
				"email":      "test@example.com",
				"password":   "password123",
				"bio":        "Test bio",
				"experience": 5,
				"rating":     4.5,
				"location":   "New York",
				"contact_info": map[string]interface{}{
					"phone":   "1234567890",
					"website": "https://example.com",
					"social_media": map[string]string{
						"instagram": "@craftsman",
					},
				},
				"is_verified": true,
			},
			expectedStatus: http.StatusCreated,
		},
		{
			name: "Missing bio",
			payload: map[string]interface{}{
				"username":   "testcraftsman2",
				"email":      "test2@example.com",
				"password":   "password123",
				"experience": 3,
				"rating":     4.0,
				"location":   "Los Angeles",
				"contact_info": map[string]interface{}{
					"phone":   "0987654321",
					"website": "https://example2.com",
				},
			},
			expectedStatus: http.StatusCreated,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Convert payload to JSON
			jsonData, err := json.Marshal(tt.payload)
			assert.NoError(t, err)

			// Create request
			req, err := http.NewRequest("POST", "/api/craftsmen", bytes.NewBuffer(jsonData))
			assert.NoError(t, err)
			req.Header.Set("Content-Type", "application/json")

			// Create response recorder
			w := httptest.NewRecorder()

			// Serve the request
			router.ServeHTTP(w, req)

			// Check status code
			assert.Equal(t, tt.expectedStatus, w.Code)

			if tt.expectedStatus == http.StatusCreated {
				// Parse response
				var response map[string]interface{}
				err = json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)

				// Check response fields
				assert.NotEmpty(t, response["access_token"])
				assert.NotEmpty(t, response["refresh_token"])

				// Check user fields
				user := response["user"].(map[string]interface{})
				assert.NotEmpty(t, user["id"])
				assert.Equal(t, tt.payload["email"], user["email"])
				assert.Equal(t, "craftsman", user["role"])

				// Check craftsman fields
				craftsman := response["craftsman"].(map[string]interface{})
				assert.NotEmpty(t, craftsman["id"])
				if bio, ok := tt.payload["bio"]; ok {
					assert.Equal(t, bio, craftsman["bio"])
				} else {
					assert.Equal(t, "", craftsman["bio"])
				}
				if exp, ok := tt.payload["experience"]; ok {
					assert.Equal(t, float64(exp.(int)), craftsman["experience"])
				}
				assert.Equal(t, tt.payload["rating"], craftsman["rating"])
				assert.Equal(t, tt.payload["location"], craftsman["location"])
				if tt.payload["contact_info"] != nil {
					assert.NotNil(t, craftsman["contact_info"])
				}
				if isVerified, ok := tt.payload["is_verified"]; ok {
					assert.Equal(t, isVerified, craftsman["is_verified"])
				} else {
					assert.Equal(t, false, craftsman["is_verified"])
				}

				// Check contact info
				contactInfo := craftsman["contact_info"].(map[string]interface{})
				expectedContactInfo := tt.payload["contact_info"].(map[string]interface{})
				assert.Equal(t, expectedContactInfo["phone"], contactInfo["phone"])
				assert.Equal(t, expectedContactInfo["website"], contactInfo["website"])
			}
		})
	}
}

// Helper function to create a test user
func createTestUser(t *testing.T) primitive.ObjectID {
	user := models.User{
		Username: "testuser",
		Email:    "test@example.com",
		Password: "password123",
		Role:     "craftsman",
	}

	result, err := Collections.Users.InsertOne(context.Background(), user)
	assert.NoError(t, err)

	return result.InsertedID.(primitive.ObjectID)
}

func TestCreateCraftHandler(t *testing.T) {
	// Setup mock database
	SetupTestDB(t)
	defer CleanupTestDB(t)

	// Create a test user first
	userID := createTestUser(t)

	// Create test router
	router := gin.Default()
	router.POST("/api/craftsmen/crafts", CreateCraft)

	tests := []struct {
		name       string
		craftData  models.Craft
		wantStatus int
	}{
		{
			name: "Valid craft creation",
			craftData: models.Craft{
				Name:        "Woodworking Basics",
				Description: "Learn basic woodworking techniques",
				Category:    "woodworking",
				Difficulty:  "beginner",
				Duration:    4,
				Price:       100.00,
				CraftsmanID: primitive.ObjectID{},
				CreatedAt:   time.Now(),
				UpdatedAt:   time.Now(),
			},
			wantStatus: http.StatusCreated,
		},
		{
			name: "Invalid craft data",
			craftData: models.Craft{
				Name:        "", // Empty name should be invalid
				Description: "Invalid craft",
				Duration:    -1, // Negative duration should be invalid
				CreatedAt:   time.Now(),
				UpdatedAt:   time.Now(),
			},
			wantStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Set the craftsman ID
			tt.craftData.CraftsmanID = userID

			// Convert craft to JSON
			jsonData, err := json.Marshal(tt.craftData)
			assert.NoError(t, err)

			// Create request
			req, err := http.NewRequest("POST", "/api/craftsmen/crafts", bytes.NewBuffer(jsonData))
			assert.NoError(t, err)
			req.Header.Set("Content-Type", "application/json")

			// Create response recorder
			w := httptest.NewRecorder()

			// Perform request
			router.ServeHTTP(w, req)

			// Check status code
			assert.Equal(t, tt.wantStatus, w.Code)

			if tt.wantStatus == http.StatusCreated {
				// Parse response
				var response map[string]interface{}
				err = json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)

				// Check response fields
				assert.NotEmpty(t, response["id"])
				assert.Equal(t, tt.craftData.Name, response["name"])
				assert.Equal(t, tt.craftData.Description, response["description"])
				assert.Equal(t, tt.craftData.Category, response["category"])
			}
		})
	}
}
