package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// CraftsmanSpeciality represents a craftsman's specialty
type CraftsmanSpeciality struct {
	Name        string `json:"name" bson:"name"`
	Description string `json:"description" bson:"description"`
}

// Craftsman represents a craftsman profile
type Craftsman struct {
	ID          primitive.ObjectID    `json:"id,omitempty" bson:"_id,omitempty"`
	UserID      primitive.ObjectID    `json:"user_id" bson:"user_id"`
	Bio         string                `json:"bio" bson:"bio"`
	Specialties []CraftsmanSpeciality `json:"specialties" bson:"specialties"`
	Experience  int                   `json:"experience" bson:"experience"`
	Rating      float64               `json:"rating" bson:"rating"`
	Location    string                `json:"location" bson:"location"`
	ContactInfo ContactInformation    `json:"contact_info" bson:"contact_info"`
	IsVerified  bool                  `json:"is_verified" bson:"is_verified"`
	CreatedAt   time.Time             `json:"created_at" bson:"created_at"`
	UpdatedAt   time.Time             `json:"updated_at" bson:"updated_at"`
}

// ContactInformation represents contact details
type ContactInformation struct {
	Phone       string            `json:"phone" bson:"phone"`
	Website     string            `json:"website" bson:"website"`
	SocialMedia map[string]string `json:"social_media" bson:"social_media"`
}
