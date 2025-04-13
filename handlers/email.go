package handlers

import (
	"context"
	"net/http"
	"time"

	"backend-dragonhak/auth"
	"backend-dragonhak/models"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type EmailVerifier struct {
	verifier *auth.EmailVerifier
}

func NewEmailVerifier(redisAddr string) *EmailVerifier {
	return &EmailVerifier{
		verifier: auth.NewEmailVerifier(redisAddr),
	}
}

func NewDummyEmailVerifier() *EmailVerifier {
	return &EmailVerifier{
		verifier: &auth.EmailVerifier{},
	}
}

// SendVerificationEmail generates a verification token and sends it to the user's email
func (ev *EmailVerifier) SendVerificationEmail(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	userID := c.GetString("userID")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	objID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	var user models.User
	err = Collections.Users.FindOne(ctx, bson.M{"_id": objID}).Decode(&user)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	if user.EmailVerified {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Email already verified"})
		return
	}

	token, err := ev.verifier.GenerateToken(ctx, user.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate verification token"})
		return
	}

	// TODO: Send email with verification link
	// For now, just return the token in the response
	c.JSON(http.StatusOK, gin.H{
		"message": "Verification email sent",
		"token":   token, // Remove this in production
	})
}

// VerifyEmail handles the email verification process
func (ev *EmailVerifier) VerifyEmail(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	token := c.Query("token")
	if token == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Verification token is required"})
		return
	}

	email, err := ev.verifier.VerifyToken(ctx, token)
	if err != nil {
		status := http.StatusInternalServerError
		if err == auth.ErrTokenInvalid || err == auth.ErrTokenExpired {
			status = http.StatusBadRequest
		}
		c.JSON(status, gin.H{"error": err.Error()})
		return
	}

	// Update user's verification status
	update := bson.M{
		"$set": bson.M{
			"email_verified": true,
			"verified_at":    time.Now(),
		},
	}

	result, err := Collections.Users.UpdateOne(ctx, bson.M{"email": email}, update)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update verification status"})
		return
	}

	if result.MatchedCount == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Email verified successfully"})
}
