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

	// Create a test user first
	userID := createTestUser(t)

	// Create test router
	router := gin.Default()
	router.POST("/api/craftsmen/profile", CreateCraftsmanProfile)

	// Test cases
	tests := []struct {
		name           string
		payload        models.Craftsman
		expectedStatus int
	}{
		{
			name: "Valid profile",
			payload: models.Craftsman{
				UserID:     userID,
				Bio:        "Test bio",
				Experience: 5,
				Rating:     4.5,
				Location:   "New York",
				ContactInfo: models.ContactInformation{
					Phone:   "1234567890",
					Website: "https://example.com",
					SocialMedia: map[string]string{
						"instagram": "@craftsman",
					},
				},
				IsVerified: true,
			},
			expectedStatus: http.StatusCreated,
		},
		{
			name: "Missing bio",
			payload: models.Craftsman{
				UserID:     userID,
				Experience: 3,
				Rating:     4.0,
				Location:   "Los Angeles",
				ContactInfo: models.ContactInformation{
					Phone:   "0987654321",
					Website: "https://example2.com",
				},
			},
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Convert payload to JSON
			jsonData, err := json.Marshal(tt.payload)
			assert.NoError(t, err)

			// Create request
			req, err := http.NewRequest("POST", "/api/craftsmen/profile", bytes.NewBuffer(jsonData))
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
				assert.NotEmpty(t, response["id"])
				assert.Equal(t, tt.payload.Bio, response["bio"])

				assert.Equal(t, tt.payload.Experience, int(response["experience"].(float64)))
				assert.Equal(t, tt.payload.Rating, response["rating"])
				assert.Equal(t, tt.payload.Location, response["location"])

				// Convert contact info from interface{} to models.ContactInformation
				contactInfoJSON, err := json.Marshal(response["contact_info"])
				assert.NoError(t, err)
				var contactInfo models.ContactInformation
				err = json.Unmarshal(contactInfoJSON, &contactInfo)
				assert.NoError(t, err)
				assert.Equal(t, tt.payload.ContactInfo, contactInfo)

				assert.Equal(t, tt.payload.IsVerified, response["is_verified"])
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
