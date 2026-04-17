package jwt_test

import (
	"os"
	"testing"

	"github.com/MovingPointP/go-task-api/pkg/jwt"
)

func TestGenerateAndParseToken(t *testing.T) {
	os.Setenv("JWT_SECRET", "test-secret-key-for-unit-test")
	os.Setenv("JWT_EXPIRATION_HOURS", "1")

	// トークン生成
	token, err := jwt.GenerateToken(1)
	if err != nil {
		t.Fatalf("GenerateToken failed: %v", err)
	}
	if token == "" {
		t.Fatal("expected non-empty token")
	}

	// トークン復元
	claims, err := jwt.ParseToken(token)
	if err != nil {
		t.Fatalf("ParseToken failed: %v", err)
	}
	if claims.UserID != 1 {
		t.Errorf("expected UserID 1, got %d", claims.UserID)
	}
}

func TestParseToken_InvalidToken(t *testing.T) {
	os.Setenv("JWT_SECRET", "test-secret-key-for-unit-test")

	_, err := jwt.ParseToken("invalid token")
	if err == nil {
		t.Error("expected error for invalid token, got nil")
	}
}

func TestParseToken_WrongSecret(t *testing.T) {
	os.Setenv("JWT_SECRET", "secret-a")
	token, _ := jwt.GenerateToken(1)

	os.Setenv("JWT_SECRET", "secret-b")
	_, err := jwt.ParseToken(token)
	if err == nil {
		t.Error("expected error for token signed with different secret, got nil")
	}
}
