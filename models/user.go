package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UserRole string

const (
	RoleAdmin     UserRole = "admin"
	RoleCraftsman UserRole = "craftsman"
	RoleCustomer  UserRole = "customer"
)

type BookingStatus string

const (
	BookingStatusPending   BookingStatus = "pending"
	BookingStatusConfirmed BookingStatus = "confirmed"
	BookingStatusCancelled BookingStatus = "cancelled"
	BookingStatusCompleted BookingStatus = "completed"
)

type User struct {
	ID        primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Name      string             `json:"name" bson:"name"`
	Surname   string             `json:"surname" bson:"surname"`
	Username  string             `json:"username" bson:"username"`
	Email     string             `json:"email" bson:"email"`
	Password  string             `json:"password" bson:"password"`
	Role      UserRole           `json:"role" bson:"role"`
	Badges    []Badge            `json:"badges,omitempty" bson:"badges,omitempty"` // Only for craftsmen and customers
	CreatedAt time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time          `json:"updated_at" bson:"updated_at"`
	// Email verification fields
	EmailVerified bool      `json:"email_verified" bson:"email_verified"`
	VerifiedAt    time.Time `json:"verified_at,omitempty" bson:"verified_at,omitempty"`
}

type Speciality struct {
	ID          primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Name        string             `json:"name" bson:"name"`
	Description string             `json:"description" bson:"description"`
	CreatedAt   time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt   time.Time          `json:"updated_at" bson:"updated_at"`
}

// Craft represents a specific craft that can be taught
type Craft struct {
	ID          primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Name        string             `json:"name" bson:"name"`
	Description string             `json:"description" bson:"description"`
	Category    string             `json:"category" bson:"category"`
	Difficulty  string             `json:"difficulty" bson:"difficulty"` // beginner, intermediate, advanced
	Duration    int                `json:"duration" bson:"duration"`     // in hours
	Price       float64            `json:"price" bson:"price"`
	CraftsmanID primitive.ObjectID `json:"craftsman_id" bson:"craftsman_id"`
	CreatedAt   time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt   time.Time          `json:"updated_at" bson:"updated_at"`
}

// Booking represents a workshop booking
type Booking struct {
	ID            primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	WorkshopID    primitive.ObjectID `json:"workshop_id" bson:"workshop_id"`
	CustomerID    primitive.ObjectID `json:"customer_id" bson:"customer_id"`
	Status        BookingStatus      `json:"status" bson:"status"` // confirmed, cancelled, completed
	PaymentStatus string             `json:"payment_status" bson:"payment_status"`
	CreatedAt     time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt     time.Time          `json:"updated_at" bson:"updated_at"`
}
