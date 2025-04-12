package utils

import (
	"context"
	"testing"
	"time"
)

// TestDBConfig holds test database configuration
type TestDBConfig struct {
	URI      string
	Database string
}

// SetupTestDB now returns our mock client instead of a real MongoDB connection
func SetupTestDB(t *testing.T, config TestDBConfig) (*MockMongoClient, func()) {
	// Create a new mock client
	mockClient := NewMockMongoClient(t)

	// Return cleanup function
	cleanup := func() {
		// Clean up all mock collections
		CleanupMockDB(t, mockClient)
	}

	return mockClient, cleanup
}

// CreateTestUser modified to work with mock client
func CreateTestUser(t *testing.T, client *MockMongoClient, dbName string) (string, func()) {
	ctx := context.Background()

	collection := client.Collections["users"]

	// Create a test user
	user := map[string]interface{}{
		"username":  "testuser",
		"email":     "test@example.com",
		"password":  "password123",
		"role":      "customer",
		"createdAt": time.Now(),
		"updatedAt": time.Now(),
	}

	result, err := collection.InsertOne(ctx, user)
	if err != nil {
		t.Fatalf("Failed to create test user: %v", err)
	}

	// Return the user ID and cleanup function
	cleanup := func() {
		_, err := collection.DeleteOne(ctx, map[string]interface{}{"_id": result.InsertedID})
		if err != nil {
			t.Logf("Warning: Failed to delete test user: %v", err)
		}
	}

	return result.InsertedID.(string), cleanup
}
