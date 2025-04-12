package handlers

import (
	"backend-dragonhak/models"
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// MockCollection implements a mock MongoDB collection
type MockCollection struct {
	Data []interface{}
}

func (mc *MockCollection) InsertOne(ctx context.Context, document interface{}, opts ...*options.InsertOneOptions) (*mongo.InsertOneResult, error) {
	// Generate a new ObjectID for the document
	id := primitive.NewObjectID()

	// If the document is a struct with an ID field, set it
	switch v := document.(type) {
	case *models.User:
		v.ID = id
		document = v
	case *models.Workshop:
		v.ID = id
		document = v
	case *models.Craft:
		v.ID = id
		document = v
	case *models.Craftsman:
		v.ID = id
		document = v
	case *models.Booking:
		v.ID = id
		document = v
	case models.User:
		v.ID = id
		document = v
	case models.Workshop:
		v.ID = id
		document = v
	case models.Craft:
		v.ID = id
		document = v
	case models.Craftsman:
		v.ID = id
		document = v
	case models.Booking:
		v.ID = id
		document = v
	}

	mc.Data = append(mc.Data, document)
	return &mongo.InsertOneResult{InsertedID: id}, nil
}

// FindOne mocks the FindOne operation
func (mc *MockCollection) FindOne(ctx context.Context, filter interface{}, opts ...*options.FindOneOptions) *mongo.SingleResult {
	// Convert filter to bson.M
	filterMap, ok := filter.(bson.M)
	if !ok {
		return mongo.NewSingleResultFromDocument(nil, mongo.ErrNoDocuments, nil)
	}

	// Find matching document
	for _, doc := range mc.Data {
		// Handle different document types
		switch v := doc.(type) {
		case *models.Workshop:
			if v.ID == filterMap["_id"] {
				return mongo.NewSingleResultFromDocument(v, nil, nil)
			}
		case *models.User:
			if v.ID == filterMap["_id"] || v.Email == filterMap["email"] {
				return mongo.NewSingleResultFromDocument(v, nil, nil)
			}
		case *models.Craft:
			if v.ID == filterMap["_id"] {
				return mongo.NewSingleResultFromDocument(v, nil, nil)
			}
		case *models.Craftsman:
			if v.ID == filterMap["_id"] {
				return mongo.NewSingleResultFromDocument(v, nil, nil)
			}
		case *models.Booking:
			if v.ID == filterMap["_id"] {
				return mongo.NewSingleResultFromDocument(v, nil, nil)
			}
		case models.Workshop:
			if v.ID == filterMap["_id"] {
				return mongo.NewSingleResultFromDocument(&v, nil, nil)
			}
		case models.User:
			if v.ID == filterMap["_id"] || v.Email == filterMap["email"] {
				return mongo.NewSingleResultFromDocument(&v, nil, nil)
			}
		case models.Craft:
			if v.ID == filterMap["_id"] {
				return mongo.NewSingleResultFromDocument(&v, nil, nil)
			}
		case models.Craftsman:
			if v.ID == filterMap["_id"] {
				return mongo.NewSingleResultFromDocument(&v, nil, nil)
			}
		case models.Booking:
			if v.ID == filterMap["_id"] {
				return mongo.NewSingleResultFromDocument(&v, nil, nil)
			}
		case bson.M:
			if id, ok := v["_id"].(primitive.ObjectID); ok {
				if id == filterMap["_id"] {
					return mongo.NewSingleResultFromDocument(v, nil, nil)
				}
			}
		}
	}

	// No matching document found
	return mongo.NewSingleResultFromDocument(nil, mongo.ErrNoDocuments, nil)
}

// UpdateOne mocks the UpdateOne operation
func (mc *MockCollection) UpdateOne(ctx context.Context, filter interface{}, update interface{}, opts ...*options.UpdateOptions) (*mongo.UpdateResult, error) {
	filterMap, ok := filter.(bson.M)
	if !ok {
		return nil, mongo.ErrNoDocuments
	}

	updateMap, ok := update.(bson.M)
	if !ok {
		return nil, mongo.ErrNoDocuments
	}

	// Find and update the document
	for i, doc := range mc.Data {
		switch v := doc.(type) {
		case *models.Workshop:
			if v.ID == filterMap["_id"] {
				if inc, ok := updateMap["$inc"].(bson.M); ok {
					if val, ok := inc["current_students"].(int); ok {
						v.CurrentStudents += val
						mc.Data[i] = v
						return &mongo.UpdateResult{MatchedCount: 1, ModifiedCount: 1}, nil
					}
				}
			}
		case models.Workshop:
			if v.ID == filterMap["_id"] {
				if inc, ok := updateMap["$inc"].(bson.M); ok {
					if val, ok := inc["current_students"].(int); ok {
						v.CurrentStudents += val
						mc.Data[i] = v
						return &mongo.UpdateResult{MatchedCount: 1, ModifiedCount: 1}, nil
					}
				}
			}
		case bson.M:
			if id, ok := v["_id"].(primitive.ObjectID); ok {
				if id == filterMap["_id"] {
					if inc, ok := updateMap["$inc"].(bson.M); ok {
						if val, ok := inc["current_students"].(int); ok {
							if current, ok := v["current_students"].(int); ok {
								v["current_students"] = current + val
								mc.Data[i] = v
								return &mongo.UpdateResult{MatchedCount: 1, ModifiedCount: 1}, nil
							}
						}
					}
				}
			}
		}
	}

	return &mongo.UpdateResult{MatchedCount: 0, ModifiedCount: 0}, nil
}

// DeleteOne mocks the DeleteOne operation
func (mc *MockCollection) DeleteOne(ctx context.Context, filter interface{}, opts ...*options.DeleteOptions) (*mongo.DeleteResult, error) {
	filterMap, ok := filter.(bson.M)
	if !ok {
		return nil, mongo.ErrNoDocuments
	}

	// Find and delete the document
	for i, doc := range mc.Data {
		switch v := doc.(type) {
		case *models.Workshop:
			if v.ID == filterMap["_id"] {
				mc.Data = append(mc.Data[:i], mc.Data[i+1:]...)
				return &mongo.DeleteResult{DeletedCount: 1}, nil
			}
		case *models.User:
			if v.ID == filterMap["_id"] {
				mc.Data = append(mc.Data[:i], mc.Data[i+1:]...)
				return &mongo.DeleteResult{DeletedCount: 1}, nil
			}
		case *models.Craft:
			if v.ID == filterMap["_id"] {
				mc.Data = append(mc.Data[:i], mc.Data[i+1:]...)
				return &mongo.DeleteResult{DeletedCount: 1}, nil
			}
		case *models.Craftsman:
			if v.ID == filterMap["_id"] {
				mc.Data = append(mc.Data[:i], mc.Data[i+1:]...)
				return &mongo.DeleteResult{DeletedCount: 1}, nil
			}
		case *models.Booking:
			if v.ID == filterMap["_id"] {
				mc.Data = append(mc.Data[:i], mc.Data[i+1:]...)
				return &mongo.DeleteResult{DeletedCount: 1}, nil
			}
		case models.Workshop:
			if v.ID == filterMap["_id"] {
				mc.Data = append(mc.Data[:i], mc.Data[i+1:]...)
				return &mongo.DeleteResult{DeletedCount: 1}, nil
			}
		case models.User:
			if v.ID == filterMap["_id"] {
				mc.Data = append(mc.Data[:i], mc.Data[i+1:]...)
				return &mongo.DeleteResult{DeletedCount: 1}, nil
			}
		case models.Craft:
			if v.ID == filterMap["_id"] {
				mc.Data = append(mc.Data[:i], mc.Data[i+1:]...)
				return &mongo.DeleteResult{DeletedCount: 1}, nil
			}
		case models.Craftsman:
			if v.ID == filterMap["_id"] {
				mc.Data = append(mc.Data[:i], mc.Data[i+1:]...)
				return &mongo.DeleteResult{DeletedCount: 1}, nil
			}
		case models.Booking:
			if v.ID == filterMap["_id"] {
				mc.Data = append(mc.Data[:i], mc.Data[i+1:]...)
				return &mongo.DeleteResult{DeletedCount: 1}, nil
			}
		case bson.M:
			if id, ok := v["_id"].(primitive.ObjectID); ok {
				if id == filterMap["_id"] {
					mc.Data = append(mc.Data[:i], mc.Data[i+1:]...)
					return &mongo.DeleteResult{DeletedCount: 1}, nil
				}
			}
		}
	}

	return &mongo.DeleteResult{DeletedCount: 0}, nil
}

// Find mocks the Find operation
func (mc *MockCollection) Find(ctx context.Context, filter interface{}, opts ...*options.FindOptions) (*mongo.Cursor, error) {
	cursor, _ := mongo.NewCursorFromDocuments(mc.Data, nil, nil)
	return cursor, nil
}

// SetupTestDB initializes test collections with mock data
func SetupTestDB(t *testing.T) {
	Collections = struct {
		Users             Collection
		CraftsmanProfiles Collection
		Crafts            Collection
		Workshops         Collection
		Badges            Collection
		Marketplace       Collection
		Auctions          Collection
		Bids              Collection
		Bookings          Collection
	}{
		Users:             &MockCollection{Data: make([]interface{}, 0)},
		CraftsmanProfiles: &MockCollection{Data: make([]interface{}, 0)},
		Crafts:            &MockCollection{Data: make([]interface{}, 0)},
		Workshops:         &MockCollection{Data: make([]interface{}, 0)},
		Badges:            &MockCollection{Data: make([]interface{}, 0)},
		Marketplace:       &MockCollection{Data: make([]interface{}, 0)},
		Auctions:          &MockCollection{Data: make([]interface{}, 0)},
		Bids:              &MockCollection{Data: make([]interface{}, 0)},
		Bookings:          &MockCollection{Data: make([]interface{}, 0)},
	}
}

// CreateTestUser creates a test user and returns its ID
func CreateTestUser(t *testing.T) primitive.ObjectID {
	user := models.User{
		Username:  "testuser",
		Email:     "test@example.com",
		Password:  "hashedpassword",
		Role:      "user",
		CreatedAt: time.Now(),
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	result, err := Collections.Users.InsertOne(ctx, user)
	assert.NoError(t, err)

	return result.InsertedID.(primitive.ObjectID)
}

// CleanupTestDB cleans up the test database
func CleanupTestDB(t *testing.T) {
	Collections = struct {
		Users             Collection
		CraftsmanProfiles Collection
		Crafts            Collection
		Workshops         Collection
		Badges            Collection
		Marketplace       Collection
		Auctions          Collection
		Bids              Collection
		Bookings          Collection
	}{
		Users:             &MockCollection{Data: make([]interface{}, 0)},
		CraftsmanProfiles: &MockCollection{Data: make([]interface{}, 0)},
		Crafts:            &MockCollection{Data: make([]interface{}, 0)},
		Workshops:         &MockCollection{Data: make([]interface{}, 0)},
		Badges:            &MockCollection{Data: make([]interface{}, 0)},
		Marketplace:       &MockCollection{Data: make([]interface{}, 0)},
		Auctions:          &MockCollection{Data: make([]interface{}, 0)},
		Bids:              &MockCollection{Data: make([]interface{}, 0)},
		Bookings:          &MockCollection{Data: make([]interface{}, 0)},
	}
}
