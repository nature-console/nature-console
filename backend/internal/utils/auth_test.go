package utils

import (
	"strings"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func TestGenerateToken(t *testing.T) {
	tests := []struct {
		name     string
		userID   uint
		email    string
		wantErr  bool
	}{
		{
			name:    "Valid user data",
			userID:  1,
			email:   "test@example.com",
			wantErr: false,
		},
		{
			name:    "Zero user ID",
			userID:  0,
			email:   "test@example.com",
			wantErr: false, // Zero ID should be allowed
		},
		{
			name:    "Empty email",
			userID:  1,
			email:   "",
			wantErr: false, // Empty email should be allowed
		},
		{
			name:    "Large user ID",
			userID:  999999,
			email:   "test@example.com",
			wantErr: false,
		},
		{
			name:    "Unicode email",
			userID:  1,
			email:   "tëst@example.com",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			token, err := GenerateToken(tt.userID, tt.email)
			
			if (err != nil) != tt.wantErr {
				t.Errorf("GenerateToken() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			
			if !tt.wantErr {
				if token == "" {
					t.Error("GenerateToken() returned empty token")
				}
				
				// Verify token structure (should have 3 parts separated by dots)
				parts := strings.Split(token, ".")
				if len(parts) != 3 {
					t.Errorf("Invalid token format, expected 3 parts, got %d", len(parts))
				}
				
				// Verify token can be parsed
				claims, err := ValidateToken(token)
				if err != nil {
					t.Errorf("Generated token is not valid: %v", err)
				}
				
				if claims.UserID != tt.userID {
					t.Errorf("Token UserID = %v, want %v", claims.UserID, tt.userID)
				}
				
				if claims.Email != tt.email {
					t.Errorf("Token Email = %v, want %v", claims.Email, tt.email)
				}
				
				// Check expiration is in the future
				if claims.ExpiresAt.Time.Before(time.Now()) {
					t.Error("Token should not be expired immediately after creation")
				}
				
				// Check expiration is approximately 24 hours from now
				expectedExpiry := time.Now().Add(24 * time.Hour)
				timeDiff := claims.ExpiresAt.Time.Sub(expectedExpiry)
				if timeDiff > time.Minute || timeDiff < -time.Minute {
					t.Errorf("Token expiry time difference too large: %v", timeDiff)
				}
			}
		})
	}
}

func TestValidateToken(t *testing.T) {
	// Generate a valid token for testing
	validUserID := uint(123)
	validEmail := "test@example.com"
	validToken, err := GenerateToken(validUserID, validEmail)
	if err != nil {
		t.Fatalf("Failed to generate valid token for testing: %v", err)
	}

	tests := []struct {
		name      string
		token     string
		wantErr   bool
		wantUserID uint
		wantEmail string
	}{
		{
			name:       "Valid token",
			token:      validToken,
			wantErr:    false,
			wantUserID: validUserID,
			wantEmail:  validEmail,
		},
		{
			name:    "Empty token",
			token:   "",
			wantErr: true,
		},
		{
			name:    "Invalid token format",
			token:   "invalid.token",
			wantErr: true,
		},
		{
			name:    "Malformed token",
			token:   "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.malformed.signature",
			wantErr: true,
		},
		{
			name:    "Token with wrong signature",
			token:   "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			claims, err := ValidateToken(tt.token)
			
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateToken() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			
			if !tt.wantErr {
				if claims == nil {
					t.Error("ValidateToken() returned nil claims for valid token")
					return
				}
				
				if claims.UserID != tt.wantUserID {
					t.Errorf("Claims UserID = %v, want %v", claims.UserID, tt.wantUserID)
				}
				
				if claims.Email != tt.wantEmail {
					t.Errorf("Claims Email = %v, want %v", claims.Email, tt.wantEmail)
				}
			}
		})
	}
}

func TestTokenExpiration(t *testing.T) {
	// Test that expired tokens are rejected
	t.Run("Expired token", func(t *testing.T) {
		// Create a token with past expiration
		claims := &Claims{
			UserID: 1,
			Email:  "test@example.com",
			RegisteredClaims: jwt.RegisteredClaims{
				ExpiresAt: jwt.NewNumericDate(time.Now().Add(-time.Hour)), // Expired 1 hour ago
				IssuedAt:  jwt.NewNumericDate(time.Now().Add(-2 * time.Hour)),
			},
		}
		
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		tokenString, err := token.SignedString(jwtSecret)
		if err != nil {
			t.Fatalf("Failed to create expired token: %v", err)
		}
		
		_, err = ValidateToken(tokenString)
		if err == nil {
			t.Error("Expected error for expired token")
		}
	})
}

func TestHashPassword(t *testing.T) {
	tests := []struct {
		name     string
		password string
		wantErr  bool
	}{
		{
			name:     "Simple password",
			password: "password123",
			wantErr:  false,
		},
		{
			name:     "Empty password",
			password: "",
			wantErr:  false, // bcrypt allows empty passwords
		},
		{
			name:     "Long password",
			password: strings.Repeat("a", 70), // bcrypt max is 72 bytes
			wantErr:  false,
		},
		{
			name:     "Special characters",
			password: "p@ssw0rd!@#$%^&*()",
			wantErr:  false,
		},
		{
			name:     "Unicode password",
			password: "pásswörd日本語",
			wantErr:  false,
		},
		{
			name:     "Very long password",
			password: strings.Repeat("password", 9), // 72 characters (8*9)
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hash, err := HashPassword(tt.password)
			
			if (err != nil) != tt.wantErr {
				t.Errorf("HashPassword() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			
			if !tt.wantErr {
				if hash == "" {
					t.Error("HashPassword() returned empty hash")
				}
				
				if hash == tt.password {
					t.Error("Hash should not be equal to original password")
				}
				
				// Verify hash starts with bcrypt prefix
				if !strings.HasPrefix(hash, "$2a$14$") && !strings.HasPrefix(hash, "$2b$14$") {
					t.Errorf("Hash format unexpected: %s", hash)
				}
				
				// Verify the hash can be used to verify the password
				if !CheckPasswordHash(tt.password, hash) {
					t.Error("Generated hash does not verify against original password")
				}
			}
		})
	}
}

func TestCheckPasswordHash(t *testing.T) {
	// Generate test hash
	password := "testpassword123"
	hash, err := HashPassword(password)
	if err != nil {
		t.Fatalf("Failed to generate hash for testing: %v", err)
	}
	
	tests := []struct {
		name     string
		password string
		hash     string
		want     bool
	}{
		{
			name:     "Correct password",
			password: password,
			hash:     hash,
			want:     true,
		},
		{
			name:     "Wrong password",
			password: "wrongpassword",
			hash:     hash,
			want:     false,
		},
		{
			name:     "Empty password with non-empty hash",
			password: "",
			hash:     hash,
			want:     false,
		},
		{
			name:     "Non-empty password with empty hash",
			password: password,
			hash:     "",
			want:     false,
		},
		{
			name:     "Empty password and empty hash",
			password: "",
			hash:     "",
			want:     false,
		},
		{
			name:     "Invalid hash format",
			password: password,
			hash:     "invalid_hash",
			want:     false,
		},
		{
			name:     "Case sensitive password",
			password: "TESTPASSWORD123",
			hash:     hash,
			want:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := CheckPasswordHash(tt.password, tt.hash)
			if got != tt.want {
				t.Errorf("CheckPasswordHash() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPasswordHashUniqueness(t *testing.T) {
	password := "samepassword"
	
	// Generate multiple hashes for the same password
	hash1, err1 := HashPassword(password)
	hash2, err2 := HashPassword(password)
	
	if err1 != nil || err2 != nil {
		t.Fatalf("Failed to generate hashes: %v, %v", err1, err2)
	}
	
	// Hashes should be different due to random salt
	if hash1 == hash2 {
		t.Error("Different hash calls should produce different hashes due to salt")
	}
	
	// But both should verify against the same password
	if !CheckPasswordHash(password, hash1) {
		t.Error("First hash should verify password")
	}
	if !CheckPasswordHash(password, hash2) {
		t.Error("Second hash should verify password")
	}
}

func TestTokenGeneration_Uniqueness(t *testing.T) {
	userID := uint(1)
	email := "test@example.com"
	
	// Generate multiple tokens for the same user
	token1, err1 := GenerateToken(userID, email)
	time.Sleep(time.Second) // Ensure different timestamps
	token2, err2 := GenerateToken(userID, email)
	
	if err1 != nil || err2 != nil {
		t.Fatalf("Failed to generate tokens: %v, %v", err1, err2)
	}
	
	// Tokens should be different due to different IssuedAt times
	if token1 == token2 {
		t.Error("Different token generations should produce different tokens")
	}
	
	// But both should be valid for the same user
	claims1, err1 := ValidateToken(token1)
	claims2, err2 := ValidateToken(token2)
	
	if err1 != nil || err2 != nil {
		t.Fatalf("Failed to validate tokens: %v, %v", err1, err2)
	}
	
	if claims1.UserID != userID || claims2.UserID != userID {
		t.Error("Both tokens should have the same UserID")
	}
	
	if claims1.Email != email || claims2.Email != email {
		t.Error("Both tokens should have the same Email")
	}
}

func TestClaims_Structure(t *testing.T) {
	claims := Claims{
		UserID: 123,
		Email:  "test@example.com",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	
	// Test that Claims implements jwt.Claims interface
	var _ jwt.Claims = &claims
	
	// Test that Claims can be used to create tokens
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &claims)
	tokenString, err := token.SignedString(jwtSecret)
	if err != nil {
		t.Fatalf("Failed to create token with Claims: %v", err)
	}
	
	// Test that the token can be parsed back
	parsedClaims, err := ValidateToken(tokenString)
	if err != nil {
		t.Fatalf("Failed to validate token with Claims: %v", err)
	}
	
	if parsedClaims.UserID != claims.UserID {
		t.Errorf("UserID mismatch: got %v, want %v", parsedClaims.UserID, claims.UserID)
	}
	
	if parsedClaims.Email != claims.Email {
		t.Errorf("Email mismatch: got %v, want %v", parsedClaims.Email, claims.Email)
	}
}

func TestPasswordHashingCost(t *testing.T) {
	password := "testpassword"
	
	// Test that the bcrypt cost is set to 14 as expected
	hash, err := HashPassword(password)
	if err != nil {
		t.Fatalf("Failed to hash password: %v", err)
	}
	
	// The cost should be visible in the hash
	if !strings.HasPrefix(hash, "$2a$14$") && !strings.HasPrefix(hash, "$2b$14$") {
		t.Errorf("Expected cost 14 in hash, got hash: %s", hash)
	}
}

func BenchmarkHashPassword(b *testing.B) {
	password := "benchmarkpassword123"
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := HashPassword(password)
		if err != nil {
			b.Fatalf("HashPassword failed: %v", err)
		}
	}
}

func BenchmarkCheckPasswordHash(b *testing.B) {
	password := "benchmarkpassword123"
	hash, err := HashPassword(password)
	if err != nil {
		b.Fatalf("Failed to generate hash: %v", err)
	}
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		CheckPasswordHash(password, hash)
	}
}

func BenchmarkGenerateToken(b *testing.B) {
	userID := uint(123)
	email := "benchmark@example.com"
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := GenerateToken(userID, email)
		if err != nil {
			b.Fatalf("GenerateToken failed: %v", err)
		}
	}
}

func BenchmarkValidateToken(b *testing.B) {
	userID := uint(123)
	email := "benchmark@example.com"
	token, err := GenerateToken(userID, email)
	if err != nil {
		b.Fatalf("Failed to generate token: %v", err)
	}
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := ValidateToken(token)
		if err != nil {
			b.Fatalf("ValidateToken failed: %v", err)
		}
	}
}