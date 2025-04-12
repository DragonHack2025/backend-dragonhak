package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Review represents a review for a workshop or craftsman
type Review struct {
	ID          primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	UserID      primitive.ObjectID `json:"user_id" bson:"user_id"`
	WorkshopID  primitive.ObjectID `json:"workshop_id,omitempty" bson:"workshop_id,omitempty"`
	CraftsmanID primitive.ObjectID `json:"craftsman_id,omitempty" bson:"craftsman_id,omitempty"`
	Rating      int                `json:"rating" bson:"rating"`
	Comment     string             `json:"comment" bson:"comment"`
	CreatedAt   time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt   time.Time          `json:"updated_at" bson:"updated_at"`
}

// UserBadge represents a badge earned by a user
type UserBadge struct {
	ID       primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	UserID   primitive.ObjectID `json:"user_id" bson:"user_id"`
	BadgeID  primitive.ObjectID `json:"badge_id" bson:"badge_id"`
	EarnedAt time.Time          `json:"earned_at" bson:"earned_at"`
	Progress int                `json:"progress" bson:"progress"` // Progress towards next level (0-100)
}

// TransactionStatus represents the status of a transaction
type TransactionStatus string

const (
	TransactionStatusPending   TransactionStatus = "pending"
	TransactionStatusCompleted TransactionStatus = "completed"
	TransactionStatusFailed    TransactionStatus = "failed"
	TransactionStatusRefunded  TransactionStatus = "refunded"
)

// Transaction represents a financial transaction
type Transaction struct {
	ID            primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	UserID        primitive.ObjectID `json:"user_id" bson:"user_id"`
	Amount        float64            `json:"amount" bson:"amount"`
	Currency      string             `json:"currency" bson:"currency"`
	Status        TransactionStatus  `json:"status" bson:"status"`
	Description   string             `json:"description" bson:"description"`
	ReferenceID   string             `json:"reference_id" bson:"reference_id"`
	PaymentMethod string             `json:"payment_method" bson:"payment_method"`
	CreatedAt     time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt     time.Time          `json:"updated_at" bson:"updated_at"`
}
