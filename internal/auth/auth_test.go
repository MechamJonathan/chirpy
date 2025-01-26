package auth

import (
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestCheckPas(t *testing.T) {
	password1 := "securePassword123"
	password2 := "evenMoreSecurePassword456"
	hash1, _ := HashPassword(password1)
	hash2, _ := HashPassword(password2)

	tests := []struct {
		name     string
		password string
		hash     string
		wantErr  bool
	}{
		{
			name:     "Correct Password",
			password: password1,
			hash:     hash1,
			wantErr:  false,
		},
		{
			name:     "Incorrect password",
			password: "wrongPassword",
			hash:     hash1,
			wantErr:  true,
		},
		{
			name:     "Password doesn't match different hash",
			password: password1,
			hash:     hash2,
			wantErr:  true,
		},
		{
			name:     "empty password",
			password: "",
			hash:     hash1,
			wantErr:  true,
		},
		{
			name:     "invalid hash",
			password: password1,
			hash:     "invalidhash",
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := CheckPasswordHash(tt.password, tt.hash)
			if (err != nil) != tt.wantErr {
				t.Errorf("CheckPasswordhash() error = %v, wantErr = %v", err, tt.wantErr)
			}
		})
	}

}

func TestJWTCreationAndValidation(t *testing.T) {
	// Create a test UUID and secret
	userId := uuid.New()
	secret := "your-test-secret"

	// Test token creation
	token, err := MakeJWT(userId, secret, time.Hour)
	if err != nil {
		t.Fatalf("Failed to create token: %v", err)
	}
	if token == "" {
		t.Fatal("Token is empty")
	}

	// Test token validation
	validatedId, err := ValidateJWT(token, secret)
	if err != nil {
		t.Fatalf("Failed to validate token: %v", err)
	}
	if validatedId != userId {
		t.Fatalf("User ID mismatch. Expected %v, got %v", userId, validatedId)
	}
}

func TestExpiredToken(t *testing.T) {
	// Create a test UUID and secret
	userId := uuid.New()
	secret := "your-test-secret"

	// Create token that's already expired (negative duration)
	token, err := MakeJWT(userId, secret, -time.Hour)
	if err != nil {
		t.Fatalf("Failed to create token: %v", err)
	}

	// Try to validate expired token
	_, err = ValidateJWT(token, secret)
	if err == nil {
		t.Fatal("Expected error for expired token, got nil")
	}
}

func TestInvalidSecretTroken(t *testing.T) {
	userId := uuid.New()
	secret := "secret"
	invalidSecret := "invalid-secret"

	token, err := MakeJWT(userId, secret, time.Hour)
	if err != nil {
		t.Fatalf("Failed to create token: %v", err)
	}

	_, err = ValidateJWT(token, invalidSecret)
	if err == nil {
		t.Fatal("Expected error for invalid secret, got nil")
	}
}
