package main

import (
	"context"
	"crypto/tls"
	"log"
	"os"
	"strconv"
	"time"

	"backend-dragonhak/config"
	"backend-dragonhak/handlers"
	"backend-dragonhak/middleware"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

var (
	client        *mongo.Client
	rateLimiter   *middleware.RateLimiter
	emailVerifier *handlers.EmailVerifier
)

func init() {
	// Try to load .env file, but don't fail if it doesn't exist
	_ = godotenv.Load()

	// Get MongoDB URI from environment
	mongodbURI := os.Getenv("MONGODB_URI")
	if mongodbURI == "" {
		log.Fatal("MONGODB_URI environment variable is not set")
	}

	// Set up MongoDB connection with Stable API
	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	clientOptions := options.Client().
		ApplyURI(mongodbURI).
		SetServerAPIOptions(serverAPI).
		SetServerSelectionTimeout(30 * time.Second).
		SetSocketTimeout(30 * time.Second).
		SetConnectTimeout(30 * time.Second).
		SetTLSConfig(&tls.Config{
			InsecureSkipVerify: false,
			MinVersion:         tls.VersionTLS12,
		})

	// Add logging for connection attempt
	log.Println("Creating MongoDB client with options...")

	// Create context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	var err error
	client, err = mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatalf("Failed to create MongoDB client: %v", err)
	}

	// Add logging for connection check
	log.Println("Checking MongoDB connection...")

	config.InitCloudinary()

	// Send a ping to confirm a successful connection with retry
	var pingErr error
	for i := 0; i < 3; i++ {
		pingErr = client.Ping(ctx, readpref.Primary())
		if pingErr == nil {
			break
		}
		log.Printf("Connection attempt %d failed: %v", i+1, pingErr)
		time.Sleep(2 * time.Second)
	}

	if pingErr != nil {
		log.Fatalf("Failed to connect to MongoDB after retries: %v", pingErr)
	}

	log.Println("Successfully connected to MongoDB!")

	// Get database name from environment
	dbName := os.Getenv("MONGODB_DATABASE")
	if dbName == "" {
		log.Fatal("MONGODB_DATABASE environment variable is not set")
	}

	// Initialize collections with database name from environment
	db := client.Database(dbName)
	handlers.Collections.Users = db.Collection("users")
	handlers.Collections.Craftsmen = db.Collection("craftsmen")
	handlers.Collections.Crafts = db.Collection("crafts")
	handlers.Collections.Workshops = db.Collection("workshops")
	handlers.Collections.Auctions = db.Collection("auctions")
	handlers.Collections.Bids = db.Collection("bids")

	// Get Redis address from environment
	redisAddr := os.Getenv("REDIS_ADDR")
	if redisAddr == "" {
		log.Println("REDIS_ADDR not set, rate limiting and email verification will be disabled")
		rateLimiter = middleware.NewDummyRateLimiter()
		emailVerifier = handlers.NewDummyEmailVerifier()
		return
	}

	// Try to initialize rate limiter and email verifier
	rateLimiter = middleware.NewRateLimiter(redisAddr)
	emailVerifier = handlers.NewEmailVerifier(redisAddr)
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
	config.AllowOrigins = []string{"*"} // Allow all origins in production
	config.AllowMethods = []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}
	config.AllowHeaders = []string{"Origin", "Content-Type", "Accept", "Authorization", "X-Requested-With"}
	config.ExposeHeaders = []string{"Content-Length"}
	config.AllowCredentials = true
	router.Use(cors.New(config))

	// Health check endpoints
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "ok",
			"time":   time.Now().Unix(),
		})
	})

	router.GET("/health/db", func(c *gin.Context) {
		if err := client.Ping(context.TODO(), readpref.Primary()); err != nil {
			c.JSON(500, gin.H{
				"status": "error",
				"error":  err.Error(),
			})
			return
		}

		c.JSON(200, gin.H{
			"status":   "ok",
			"database": "connected",
		})
	})

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
		authRoutes.POST("/register", handlers.CreateUser)
		authRoutes.POST("/register/craftsman", handlers.CreateCraftsmanProfile)
		authRoutes.POST("/refresh", handlers.RefreshToken)
		authRoutes.POST("/logout", middleware.AuthMiddleware(os.Getenv("JWT_ACCESS_SECRET")), handlers.Logout)
	}

	// User routes
	userRoutes := router.Group("/api/users")
	{
		userRoutes.GET("/", handlers.GetUsers)
		userRoutes.GET("/:id", handlers.GetUser)
		userRoutes.Use(middleware.AuthMiddleware(os.Getenv("JWT_ACCESS_SECRET")))
		{
			userRoutes.PUT("/:id", handlers.UpdateUser)
			userRoutes.DELETE("/:id", handlers.DeleteUser)
			userRoutes.GET("/:id/badges", handlers.GetUserBadges)

			// Email verification routes
			userRoutes.POST("/verify/send", emailVerifier.SendVerificationEmail)
			userRoutes.GET("/verify", emailVerifier.VerifyEmail)
		}
	}

	// Craftsman routes
	craftsmanRoutes := router.Group("/api/craftsmen")
	{
		craftsmanRoutes.GET("", handlers.GetCraftsmen)
		craftsmanRoutes.Use(middleware.AuthMiddleware(os.Getenv("JWT_ACCESS_SECRET")))
		{
			craftsmanRoutes.GET("/:id", handlers.GetCraftsman)
			craftsmanRoutes.PUT("/:id", handlers.UpdateCraftsman)
			craftsmanRoutes.DELETE("/:id", handlers.DeleteCraftsman)
		}
	}

	// Customer routes
	customerRoutes := router.Group("/api/customers")
	{
		customerRoutes.GET("/search/craftsmen", handlers.SearchCraftsmen)
		customerRoutes.GET("/search/workshops", handlers.SearchWorkshops)
		customerRoutes.Use(middleware.AuthMiddleware(os.Getenv("JWT_ACCESS_SECRET")))
		{
			customerRoutes.POST("/bookings", handlers.BookWorkshop)
			customerRoutes.GET("/:id/bookings", handlers.GetCustomerBookings)
		}
	}

	// Badge routes
	badgeRoutes := router.Group("/api/badges")
	badgeRoutes.Use(middleware.AuthMiddleware(os.Getenv("JWT_ACCESS_SECRET")))
	{
		badgeRoutes.POST("/", handlers.CreateBadge)
		badgeRoutes.GET("/", handlers.GetBadges)
		badgeRoutes.POST("/:badgeId/award/:userId", handlers.AwardBadge)
	}

	// Image routes
	imageRoutes := router.Group("/api/images")
	{
		imageRoutes.GET("/:public_id", handlers.GetImage)
		imageRoutes.Use(middleware.AuthMiddleware(os.Getenv("JWT_ACCESS_SECRET")))
		{
			imageRoutes.POST("/upload", handlers.UploadImage)
			imageRoutes.DELETE("/:public_id", handlers.DeleteImage)
		}
	}

	// Auction routes
	auctionRoutes := router.Group("/api/auctions")
	{
		auctionRoutes.GET("/", handlers.GetAuctions)
		auctionRoutes.GET("/:id", handlers.GetAuction)
		auctionRoutes.Use(middleware.AuthMiddleware(os.Getenv("JWT_ACCESS_SECRET")))
		{
			auctionRoutes.POST("/", handlers.CreateAuction)
			auctionRoutes.POST("/:id/bids", handlers.PlaceBid)
			auctionRoutes.GET("/:id/bids", handlers.GetAuctionBids)
		}
	}

	// Get port from environment variable or use default
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Start server
	log.Printf("Server starting on port %s in %s mode", port, gin.Mode())
	if err := router.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
