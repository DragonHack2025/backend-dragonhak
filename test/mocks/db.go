package mocks

import (
	"context"

	"backend-dragonhak/models"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// MockCollections provides mock database collections for testing
type MockCollections struct {
	Users        []models.User
	Crafts       []models.Craft
	Workshops    []models.Workshop
	Bookings     []models.WorkshopBooking
	Reviews      []models.Review
	Badges       []models.Badge
	UserBadges   []models.UserBadge
	Transactions []models.Transaction
	Craftsmen    []models.Craftsman
}

// NewMockCollections creates a new instance of mock collections
func NewMockCollections() *MockCollections {
	return &MockCollections{
		Users:        make([]models.User, 0),
		Crafts:       make([]models.Craft, 0),
		Workshops:    make([]models.Workshop, 0),
		Bookings:     make([]models.WorkshopBooking, 0),
		Reviews:      make([]models.Review, 0),
		Badges:       make([]models.Badge, 0),
		UserBadges:   make([]models.UserBadge, 0),
		Transactions: make([]models.Transaction, 0),
		Craftsmen:    make([]models.Craftsman, 0),
	}
}

// MockDB provides a mock database implementation for testing
type MockDB struct {
	Collections *MockCollections
}

// NewMockDB creates a new mock database instance
func NewMockDB() *MockDB {
	return &MockDB{
		Collections: NewMockCollections(),
	}
}

// InsertOne mocks the MongoDB InsertOne operation
func (m *MockDB) InsertOne(ctx context.Context, collection string, document interface{}) (primitive.ObjectID, error) {
	// Generate a new ObjectID
	id := primitive.NewObjectID()

	// Set the ID field based on the document type
	switch doc := document.(type) {
	case *models.User:
		doc.ID = id
		m.Collections.Users = append(m.Collections.Users, *doc)
	case *models.Craftsman:
		doc.ID = id
		m.Collections.Craftsmen = append(m.Collections.Craftsmen, *doc)
	case *models.Craft:
		doc.ID = id
		m.Collections.Crafts = append(m.Collections.Crafts, *doc)
	case *models.Workshop:
		doc.ID = id
		m.Collections.Workshops = append(m.Collections.Workshops, *doc)
	case *models.WorkshopBooking:
		doc.ID = id
		m.Collections.Bookings = append(m.Collections.Bookings, *doc)
	case *models.Review:
		doc.ID = id
		m.Collections.Reviews = append(m.Collections.Reviews, *doc)
	case *models.Badge:
		doc.ID = id
		m.Collections.Badges = append(m.Collections.Badges, *doc)
	case *models.UserBadge:
		doc.ID = id
		m.Collections.UserBadges = append(m.Collections.UserBadges, *doc)
	case *models.Transaction:
		doc.ID = id
		m.Collections.Transactions = append(m.Collections.Transactions, *doc)
	}

	return id, nil
}

// FindOne mocks the MongoDB FindOne operation
func (m *MockDB) FindOne(ctx context.Context, collection string, filter interface{}, result interface{}) error {
	// Implementation depends on the collection and filter type
	// This is a simplified version
	switch collection {
	case "users":
		if users, ok := result.(*[]models.User); ok {
			*users = m.Collections.Users
		}
	case "craftsmen":
		if craftsmen, ok := result.(*[]models.Craftsman); ok {
			*craftsmen = m.Collections.Craftsmen
		}
	case "crafts":
		if crafts, ok := result.(*[]models.Craft); ok {
			*crafts = m.Collections.Crafts
		}
	case "workshops":
		if workshops, ok := result.(*[]models.Workshop); ok {
			*workshops = m.Collections.Workshops
		}
	case "bookings":
		if bookings, ok := result.(*[]models.WorkshopBooking); ok {
			*bookings = m.Collections.Bookings
		}
	case "reviews":
		if reviews, ok := result.(*[]models.Review); ok {
			*reviews = m.Collections.Reviews
		}
	case "badges":
		if badges, ok := result.(*[]models.Badge); ok {
			*badges = m.Collections.Badges
		}
	case "user_badges":
		if userBadges, ok := result.(*[]models.UserBadge); ok {
			*userBadges = m.Collections.UserBadges
		}
	case "transactions":
		if transactions, ok := result.(*[]models.Transaction); ok {
			*transactions = m.Collections.Transactions
		}
	}

	return nil
}

// UpdateOne mocks the MongoDB UpdateOne operation
func (m *MockDB) UpdateOne(ctx context.Context, collection string, filter interface{}, update interface{}) error {
	// Implementation depends on the collection and filter type
	// This is a simplified version
	return nil
}

// DeleteOne mocks the MongoDB DeleteOne operation
func (m *MockDB) DeleteOne(ctx context.Context, collection string, filter interface{}) error {
	// Implementation depends on the collection and filter type
	// This is a simplified version
	return nil
}
