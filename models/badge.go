package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// BadgeCategory represents the category of a badge
type BadgeCategory string

const (
	BadgeCategoryLearning BadgeCategory = "learning"
	BadgeCategorySocial   BadgeCategory = "social"
	BadgeCategoryExpert   BadgeCategory = "expert"
)

// Badge represents a badge that can be earned by users
type Badge struct {
	ID          primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Name        string             `json:"name" bson:"name"`
	Description string             `json:"description" bson:"description"`
	Icon        string             `json:"icon" bson:"icon"`
	Category    BadgeCategory      `json:"category" bson:"category"`
	CreatedAt   time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt   time.Time          `json:"updated_at" bson:"updated_at"`
}
