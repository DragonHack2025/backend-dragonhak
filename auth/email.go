package auth

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"time"

	"github.com/redis/go-redis/v9"
)

const (
	verificationTokenLength = 32
	tokenExpiryDuration     = 24 * time.Hour
)

var (
	ErrTokenExpired = errors.New("verification token has expired")
	ErrTokenInvalid = errors.New("invalid verification token")
)

type EmailVerifier struct {
	client *redis.Client
}

func NewEmailVerifier(redisAddr string) *EmailVerifier {
	client := redis.NewClient(&redis.Options{
		Addr: redisAddr,
	})
	return &EmailVerifier{client: client}
}

// GenerateToken creates a new verification token for the given email
func (ev *EmailVerifier) GenerateToken(ctx context.Context, email string) (string, error) {
	// Generate random token
	tokenBytes := make([]byte, verificationTokenLength)
	if _, err := rand.Read(tokenBytes); err != nil {
		return "", err
	}
	token := hex.EncodeToString(tokenBytes)

	// Store token in Redis with expiry
	err := ev.client.Set(ctx, "email_verify:"+token, email, tokenExpiryDuration).Err()
	if err != nil {
		return "", err
	}

	return token, nil
}

// VerifyToken validates the token and returns the associated email
func (ev *EmailVerifier) VerifyToken(ctx context.Context, token string) (string, error) {
	// Get email from Redis
	email, err := ev.client.Get(ctx, "email_verify:"+token).Result()
	if err == redis.Nil {
		return "", ErrTokenInvalid
	} else if err != nil {
		return "", err
	}

	// Delete token after use
	ev.client.Del(ctx, "email_verify:"+token)

	return email, nil
}
