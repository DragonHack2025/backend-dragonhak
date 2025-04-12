package main

import (
	"context"
	"log"
	"os"
	"strconv"
	"time"

	"backend-dragonhak/handlers"
	"backend-dragonhak/middleware"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	client        *mongo.Client
	rateLimiter   *middleware.RateLimiter
	emailVerifier *handlers.EmailVerifier
)

func init() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	// Set up MongoDB connection
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	clientOptions := options.Client().ApplyURI(os.Getenv("MONGODB_URI"))
	var err error
	client, err = mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	// Check the connection
	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Connected to MongoDB!")

	// Initialize collections
	db := client.Database("dragonhak")
	handlers.Collections.Users = db.Collection("users")
	handlers.Collections.CraftsmanProfiles = db.Collection("craftsman_profiles")
	handlers.Collections.Crafts = db.Collection("crafts")
	handlers.Collections.Workshops = db.Collection("workshops")

	// Initialize rate limiter
	rateLimiter = middleware.NewRateLimiter(os.Getenv("REDIS_ADDR"))
	emailVerifier = handlers.NewEmailVerifier(os.Getenv("REDIS_ADDR"))
}

func main() {
	// Set release mode in production
	if os.Getenv("GIN_MODE") == "release" {
		gin.SetMode(gin.ReleaseMode)
	}

	// Create router with security middleware
	router := gin.New()

	// Trust only localhost in development, configure for production
	if os.Getenv("GIN_MODE") == "release" {
		router.SetTrustedProxies([]string{"127.0.0.1"})
	} else {
		router.SetTrustedProxies(nil) // Trust all proxies in development
	}

	// Add security middleware
	router.Use(gin.Recovery())
	router.Use(gin.Logger())

	// Configure CORS
	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"http://localhost:3000"} // Add your frontend URL
	config.AllowCredentials = true
	config.AllowHeaders = []string{"Origin", "Content-Length", "Content-Type", "Authorization"}
	router.Use(cors.New(config))

	// Basic route
	router.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Welcome to DragonHak API",
		})
	})

	// Parse rate limit settings
	window, _ := strconv.Atoi(os.Getenv("RATE_LIMIT_WINDOW"))
	maxRequests, _ := strconv.Atoi(os.Getenv("RATE_LIMIT_MAX_REQUESTS"))

	// Auth routes with rate limiting
	authRoutes := router.Group("/api/auth")
	authRoutes.Use(rateLimiter.Limit(maxRequests, time.Duration(window)*time.Second))
	{
		authRoutes.POST("/login", handlers.Login)
		authRoutes.POST("/refresh", handlers.RefreshToken)
		authRoutes.POST("/logout", middleware.AuthMiddleware(os.Getenv("JWT_ACCESS_SECRET")), handlers.Logout)
	}

	// User routes
	userRoutes := router.Group("/api/users")
	userRoutes.Use(middleware.AuthMiddleware(os.Getenv("JWT_ACCESS_SECRET")))
	{
		userRoutes.GET("/", handlers.GetUsers)
		userRoutes.GET("/:id", handlers.GetUser)
		userRoutes.POST("/", handlers.CreateUser)
		userRoutes.PUT("/:id", handlers.UpdateUser)
		userRoutes.DELETE("/:id", handlers.DeleteUser)
		userRoutes.GET("/:id/badges", handlers.GetUserBadges)

		// Email verification routes
		userRoutes.POST("/verify/send", emailVerifier.SendVerificationEmail)
		userRoutes.GET("/verify", emailVerifier.VerifyEmail)
	}

	// Craftsman routes
	craftsmanRoutes := router.Group("/api/craftsmen")
	craftsmanRoutes.Use(middleware.AuthMiddleware(os.Getenv("JWT_ACCESS_SECRET")))
	{
		craftsmanRoutes.POST("/profile", handlers.CreateCraftsmanProfile)
		craftsmanRoutes.PUT("/profile/:id", handlers.UpdateCraftsmanProfile)
		craftsmanRoutes.POST("/crafts", handlers.CreateCraft)
		craftsmanRoutes.POST("/workshops", handlers.CreateWorkshop)
		craftsmanRoutes.GET("/:id/workshops", handlers.GetCraftsmanWorkshops)
	}

	// Customer routes
	customerRoutes := router.Group("/api/customers")
	customerRoutes.Use(middleware.AuthMiddleware(os.Getenv("JWT_ACCESS_SECRET")))
	{
		customerRoutes.GET("/search/craftsmen", handlers.SearchCraftsmen)
		customerRoutes.GET("/search/workshops", handlers.SearchWorkshops)
		customerRoutes.POST("/bookings", handlers.BookWorkshop)
		customerRoutes.GET("/:id/bookings", handlers.GetCustomerBookings)
	}

	// Badge routes
	badgeRoutes := router.Group("/api/badges")
	badgeRoutes.Use(middleware.AuthMiddleware(os.Getenv("JWT_ACCESS_SECRET")))
	{
		badgeRoutes.POST("/", handlers.CreateBadge)
		badgeRoutes.GET("/", handlers.GetBadges)
		badgeRoutes.POST("/:badgeId/award/:userId", handlers.AwardBadge)
	}

	// Marketplace routes
	marketplaceRoutes := router.Group("/api/marketplace")
	marketplaceRoutes.Use(middleware.AuthMiddleware(os.Getenv("JWT_ACCESS_SECRET")))
	{
		marketplaceRoutes.POST("/items", handlers.CreateMarketplaceItem)
		marketplaceRoutes.GET("/items", handlers.GetMarketplaceItems)
		marketplaceRoutes.GET("/items/:id", handlers.GetMarketplaceItem)
		marketplaceRoutes.POST("/auctions/:auctionId/bids", handlers.PlaceBid)
		marketplaceRoutes.GET("/auctions/:auctionId/bids", handlers.GetAuctionBids)
	}

	// Get port from environment variable or use default
	port := os.Getenv("PORT")
	if port == "" {
		port = "5000"
	}

	// Start server
	log.Printf("Server starting on port %s in %s mode", port, gin.Mode())
	if err := router.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
