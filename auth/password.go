package auth

import (
	"errors"
	"unicode"

	"golang.org/x/crypto/bcrypt"
)

var (
	ErrPasswordTooShort  = errors.New("password must be at least 8 characters long")
	ErrPasswordNoUpper   = errors.New("password must contain at least one uppercase letter")
	ErrPasswordNoLower   = errors.New("password must contain at least one lowercase letter")
	ErrPasswordNoNumber  = errors.New("password must contain at least one number")
	ErrPasswordNoSpecial = errors.New("password must contain at least one special character")
	ErrPasswordCommon    = errors.New("password is too common or easily guessable")
)

// commonPasswords is a list of commonly used passwords that should be rejected
var commonPasswords = map[string]bool{
	"password": true,
	"123456":   true,
	"qwerty":   true,
	"admin":    true,
	"letmein":  true,
	"welcome":  true,
}

// ValidatePassword checks if a password meets security requirements
func ValidatePassword(password string) error {
	// Check length
	if len(password) < 8 {
		return ErrPasswordTooShort
	}

	// Check for common passwords
	if commonPasswords[password] {
		return ErrPasswordCommon
	}

	var hasUpper, hasLower, hasNumber, hasSpecial bool
	for _, char := range password {
		switch {
		case unicode.IsUpper(char):
			hasUpper = true
		case unicode.IsLower(char):
			hasLower = true
		case unicode.IsNumber(char):
			hasNumber = true
		case unicode.IsPunct(char) || unicode.IsSymbol(char):
			hasSpecial = true
		}
	}

	if !hasUpper {
		return ErrPasswordNoUpper
	}
	if !hasLower {
		return ErrPasswordNoLower
	}
	if !hasNumber {
		return ErrPasswordNoNumber
	}
	if !hasSpecial {
		return ErrPasswordNoSpecial
	}

	return nil
}

// HashPassword creates a bcrypt hash of the password
func HashPassword(password string) (string, error) {
	// Validate password before hashing
	if err := ValidatePassword(password); err != nil {
		return "", err
	}

	// Generate hash with cost of 12
	hash, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return "", err
	}

	return string(hash), nil
}

// ComparePassword compares a password with its hash
func ComparePassword(password, hash string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
}

// GenerateTemporaryPassword generates a secure temporary password
func GenerateTemporaryPassword() (string, error) {
	// This is a placeholder. In a real implementation, you would:
	// 1. Generate a cryptographically secure random string
	// 2. Ensure it meets password requirements
	// 3. Return it along with its hash
	return "", errors.New("not implemented")
}
