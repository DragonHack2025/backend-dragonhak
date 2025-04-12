package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// WorkshopBooking represents a booking for a workshop
type WorkshopBooking struct {
	ID            primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	WorkshopID    primitive.ObjectID `json:"workshop_id" bson:"workshop_id"`
	UserID        primitive.ObjectID `json:"user_id" bson:"user_id"`
	Status        string             `json:"status" bson:"status"`
	PaymentStatus string             `json:"payment_status" bson:"payment_status"`
	CreatedAt     time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt     time.Time          `json:"updated_at" bson:"updated_at"`
}

// Workshop represents a workshop event
type Workshop struct {
	ID              primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Title           string             `json:"title" bson:"title"`
	Description     string             `json:"description" bson:"description"`
	Date            time.Time          `json:"date" bson:"date"`
	Duration        int                `json:"duration" bson:"duration"`
	MaxParticipants int                `json:"max_participants" bson:"max_participants"`
	CurrentStudents int                `json:"current_students" bson:"current_students"`
	Price           float64            `json:"price" bson:"price"`
	Location        string             `json:"location" bson:"location"`
	CraftsmanID     primitive.ObjectID `json:"craftsman_id" bson:"craftsman_id"`
	CreatedAt       time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt       time.Time          `json:"updated_at" bson:"updated_at"`
}
