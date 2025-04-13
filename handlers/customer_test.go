package handlers

import (
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

func TestSearchCraftsmen(t *testing.T) {
	// Setup mock database
	SetupTestDB(t)
	defer CleanupTestDB(t)

	// Create test router
	router := gin.Default()
	router.GET("/api/customers/craftsmen/search", SearchCraftsmen)

	// Create test craftsman profile
	craftsmanProfile := models.Craftsman{
		UserID: primitive.NewObjectID(),
		Bio:    "Experienced woodworker",
		Specialties: []models.CraftsmanSpeciality{
			{
				Name:        "woodworking",
				Description: "Experienced woodworker",	
			},
		},
		Experience: 10,
		Rating:     4.5,
		Location:   "New York",
		ContactInfo: models.ContactInformation{
			Phone:   "123-456-7890",
			Website: "www.example.com",
			SocialMedia: map[string]string{
				"instagram": "@craftsman",
			},
		},
		IsVerified: true,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	// Insert the profile into the database
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	_, err := Collections.Craftsmen.InsertOne(ctx, craftsmanProfile)
	assert.NoError(t, err)

	tests := []struct {
		name       string
		query      string
		wantStatus int
	}{
		{
			name:       "Search by specialty",
			query:      "?specialty=woodworking",
			wantStatus: http.StatusOK,
		},
		{
			name:       "Search by location",
			query:      "?location=New York",
			wantStatus: http.StatusOK,
		},
		{
			name:       "Search by rating",
			query:      "?minRating=4.0",
			wantStatus: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create request
			req, err := http.NewRequest("GET", "/api/customers/craftsmen/search"+tt.query, nil)
			assert.NoError(t, err)

			// Create response recorder
			w := httptest.NewRecorder()

			// Perform request
			router.ServeHTTP(w, req)

			// Check status code
			assert.Equal(t, tt.wantStatus, w.Code)

			if tt.wantStatus == http.StatusOK {
				// Parse response
				var response []models.Craftsman
				err = json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)

				// Check that we got at least one result
				assert.Greater(t, len(response), 0)
			}
		})
	}
}

func TestBookWorkshop(t *testing.T) {
	// Setup mock database
	SetupTestDB(t)
	defer CleanupTestDB(t)

	// Create test router
	router := gin.Default()
	router.POST("/api/customers/workshops/:id/book", BookWorkshop)

	// Create test user (customer)
	customerID := CreateTestUser(t)

	// Create test craftsman user and profile
	craftsmanID := CreateTestUser(t)

	// Create and insert craftsman profile
	craftsmanProfile := models.Craftsman{
		UserID: craftsmanID,
		Bio:    "Experienced woodworker",
		Specialties: []models.CraftsmanSpeciality{
			{
				Name:        "woodworking",
				Description: "Experienced woodworker",
			},
		},
		Experience: 10,
		Rating:     4.5,
		Location:   "New York",
		ContactInfo: models.ContactInformation{
			Phone:   "123-456-7890",
			Website: "www.example.com",
			SocialMedia: map[string]string{
				"instagram": "@craftsman",
			},
		},
		IsVerified: true,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	// Insert the profile into the database
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	_, err := Collections.Craftsmen.InsertOne(ctx, craftsmanProfile)
	assert.NoError(t, err)

	// Create test workshop
	workshop := models.Workshop{
		ID:              primitive.NewObjectID(),
		CraftsmanID:     craftsmanID,
		Title:           "Woodworking Workshop",
		Description:     "Learn woodworking basics",
		Date:            time.Now().Add(24 * time.Hour),
		Duration:        4,
		MaxParticipants: 10,
		Price:           100.00,
		Location:        "New York",
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}

	// Insert the workshop into the database
	workshopResult, err := Collections.Workshops.InsertOne(ctx, workshop)
	assert.NoError(t, err)
	workshopID := workshopResult.InsertedID.(primitive.ObjectID)

	tests := []struct {
		name       string
		workshopID string
		wantStatus int
	}{
		{
			name:       "Valid workshop booking",
			workshopID: workshopID.Hex(),
			wantStatus: http.StatusCreated,
		},
		{
			name:       "Invalid workshop ID",
			workshopID: "invalid-id",
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "Non-existent workshop ID",
			workshopID: primitive.NewObjectID().Hex(),
			wantStatus: http.StatusNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create request
			req, err := http.NewRequest("POST", "/api/customers/workshops/"+tt.workshopID+"/book", nil)
			assert.NoError(t, err)

			// Add customer ID to header
			req.Header.Set("X-Customer-ID", customerID.Hex())

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
				assert.Equal(t, tt.workshopID, response["workshop_id"])
				assert.Equal(t, customerID.Hex(), response["customer_id"])
			}
		})
	}
}
