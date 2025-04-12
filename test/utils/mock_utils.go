package utils

import (
	"context"
	"testing"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// MockCollection is a mock MongoDB collection
type MockCollection struct {
	Documents map[primitive.ObjectID]interface{}
}

// MockMongoClient is a mock MongoDB client for testing
type MockMongoClient struct {
	Collections map[string]Collection
}

// NewMockMongoClient creates a new mock MongoDB client
func NewMockMongoClient(t *testing.T) *MockMongoClient {
	// Create collections map
	collections := make(map[string]Collection)

	// Create mock collections
	collections["users"] = &MockCollection{
		Documents: make(map[primitive.ObjectID]interface{}),
	}
	collections["craftsman_profiles"] = &MockCollection{
		Documents: make(map[primitive.ObjectID]interface{}),
	}
	collections["crafts"] = &MockCollection{
		Documents: make(map[primitive.ObjectID]interface{}),
	}
	collections["workshops"] = &MockCollection{
		Documents: make(map[primitive.ObjectID]interface{}),
	}
	collections["bookings"] = &MockCollection{
		Documents: make(map[primitive.ObjectID]interface{}),
	}
	collections["badges"] = &MockCollection{
		Documents: make(map[primitive.ObjectID]interface{}),
	}
	collections["marketplace_items"] = &MockCollection{
		Documents: make(map[primitive.ObjectID]interface{}),
	}
	collections["auctions"] = &MockCollection{
		Documents: make(map[primitive.ObjectID]interface{}),
	}
	collections["bids"] = &MockCollection{
		Documents: make(map[primitive.ObjectID]interface{}),
	}

	return &MockMongoClient{
		Collections: collections,
	}
}

// SetupMockDB sets up the mock database with collections
func SetupMockDB(t *testing.T) *MockMongoClient {
	mockClient := NewMockMongoClient(t)
	return mockClient
}

// CleanupMockDB cleans up the mock database
func CleanupMockDB(t *testing.T, client *MockMongoClient) {
	// Clear all collections
	for _, collection := range client.Collections {
		if mc, ok := collection.(*MockCollection); ok {
			mc.Documents = make(map[primitive.ObjectID]interface{})
		}
	}
}

// InsertOne mocks the InsertOne operation
func (mc *MockCollection) InsertOne(ctx context.Context, document interface{}, opts ...*options.InsertOneOptions) (*mongo.InsertOneResult, error) {
	// Generate a new ObjectID for the document
	id := primitive.NewObjectID()

	// Store the document in our mock collection
	mc.Documents[id] = document

	// Return the result
	return &mongo.InsertOneResult{
		InsertedID: id,
	}, nil
}

// FindOne mocks the FindOne operation
func (mc *MockCollection) FindOne(ctx context.Context, filter interface{}, opts ...*options.FindOneOptions) *mongo.SingleResult {
	// Convert filter to bson.M
	filterMap, ok := filter.(bson.M)
	if !ok {
		return mongo.NewSingleResultFromDocument(nil, nil, nil)
	}

	// Find matching document
	for id, doc := range mc.Documents {
		if filterMap["_id"] == id {
			return mongo.NewSingleResultFromDocument(doc, nil, nil)
		}
	}

	return mongo.NewSingleResultFromDocument(nil, mongo.ErrNoDocuments, nil)
}

// Find mocks the Find operation
func (mc *MockCollection) Find(ctx context.Context, filter interface{}, opts ...*options.FindOptions) (*mongo.Cursor, error) {
	// For now, return all documents
	var docs []interface{}
	for _, doc := range mc.Documents {
		docs = append(docs, doc)
	}

	cursor, _ := mongo.NewCursorFromDocuments(docs, nil, nil)
	return cursor, nil
}

// DeleteOne mocks the DeleteOne operation
func (mc *MockCollection) DeleteOne(ctx context.Context, filter interface{}, opts ...*options.DeleteOptions) (*mongo.DeleteResult, error) {
	// Convert filter to bson.M
	filterMap, ok := filter.(bson.M)
	if !ok {
		return &mongo.DeleteResult{DeletedCount: 0}, nil
	}

	// Delete matching document
	for id := range mc.Documents {
		if filterMap["_id"] == id {
			delete(mc.Documents, id)
			return &mongo.DeleteResult{DeletedCount: 1}, nil
		}
	}

	return &mongo.DeleteResult{DeletedCount: 0}, nil
}

// Drop mocks the Drop operation
func (mc *MockCollection) Drop(ctx context.Context) error {
	mc.Documents = make(map[primitive.ObjectID]interface{})
	return nil
}
