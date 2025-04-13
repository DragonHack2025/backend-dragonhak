package fixtures

import (
	"time"

	"backend-dragonhak/models"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// SampleUsers provides test user data
var SampleUsers = []models.User{
	{
		ID:        primitive.NewObjectID(),
		Username:  "testuser1",
		Email:     "test1@example.com",
		Password:  "hashedpassword1",
		Role:      models.RoleCustomer,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	},
	{
		ID:        primitive.NewObjectID(),
		Username:  "testuser2",
		Email:     "test2@example.com",
		Password:  "hashedpassword2",
		Role:      models.RoleAdmin,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	},
}

// SampleCraftsmen provides test craftsman data
var SampleCraftsmen = []models.Craftsman{
	{
		ID:         primitive.NewObjectID(),
		UserID:     SampleUsers[0].ID,
		Bio:        "Experienced woodworker",
		Experience: 5,
		Rating:     4.5,
		Location:   "New York",
		ContactInfo: models.ContactInformation{
			Phone:   "+1234567890",
			Website: "https://example.com",
			SocialMedia: map[string]string{
				"instagram": "@woodworker",
			},
		},
		IsVerified: true,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	},
}

// SampleCrafts provides test craft data
var SampleCrafts = []models.Craft{
	{
		ID:          primitive.NewObjectID(),
		Name:        "Wooden Chair",
		Description: "Handcrafted wooden chair",
		Category:    "furniture",
		Difficulty:  "intermediate",
		Duration:    4,
		Price:       150.00,
		CraftsmanID: SampleUsers[0].ID,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	},
}

// SampleWorkshops provides test workshop data
var SampleWorkshops = []models.Workshop{
	{
		ID:              primitive.NewObjectID(),
		Title:           "Woodworking Basics",
		Description:     "Learn basic woodworking techniques",
		Date:            time.Now().Add(24 * time.Hour),
		Duration:        4,
		MaxParticipants: 10,
		Price:           50.00,
		Location:        "Workshop Studio",
		CraftsmanID:     SampleUsers[0].ID,
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	},
}

// SampleBookings provides test booking data
var SampleBookings = []models.WorkshopBooking{
	{
		ID:            primitive.NewObjectID(),
		WorkshopID:    SampleWorkshops[0].ID,
		UserID:        SampleUsers[0].ID,
		Status:        "confirmed",
		PaymentStatus: "paid",
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	},
}

// SampleReviews provides test review data
var SampleReviews = []models.Review{
	{
		ID:         primitive.NewObjectID(),
		WorkshopID: SampleWorkshops[0].ID,
		UserID:     SampleUsers[0].ID,
		Rating:     5,
		Comment:    "Great workshop!",
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	},
}

// SampleBadges provides test badge data
var SampleBadges = []models.Badge{
	{
		ID:          primitive.NewObjectID(),
		Name:        "First Workshop",
		Description: "Attended first workshop",
		Icon:        "https://example.com/badge1.png",
		Category:    models.BadgeCategoryLearning,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	},
}

// SampleTransactions provides test transaction data
var SampleTransactions = []models.Transaction{
	{
		ID:            primitive.NewObjectID(),
		UserID:        SampleUsers[0].ID,
		Amount:        50.00,
		Currency:      "USD",
		Status:        models.TransactionStatusCompleted,
		Description:   "Workshop payment",
		ReferenceID:   "ref123",
		PaymentMethod: "credit_card",
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	},
}
