package jwt_test

import (
	"os"
	"testing"

	"github.com/MovingPointP/go-task-api/pkg/jwt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGenerateAndParseToken(t *testing.T) {
	os.Setenv("JWT_SECRET", "test-secret-key-for-unit-test")
	os.Setenv("JWT_EXPIRATION_HOURS", "1")

	// トークン生成
	token, err := jwt.GenerateToken(1)
	require.NoError(t, err, "GenerateToken failed")
	require.NotEmpty(t, token)

	// トークン復元
	claims, err := jwt.ParseToken(token)
	require.NoError(t, err, "ParseToken failed")
	assert.Equal(t, uint(1), claims.UserID)
}

func TestParseToken_InvalidToken(t *testing.T) {
	os.Setenv("JWT_SECRET", "test-secret-key-for-unit-test")

	_, err := jwt.ParseToken("invalid token")
	assert.Error(t, err, "expected error for invalid token, got nil")
}

func TestParseToken_WrongSecret(t *testing.T) {
	os.Setenv("JWT_SECRET", "secret-a")
	token, _ := jwt.GenerateToken(1)

	os.Setenv("JWT_SECRET", "secret-b")
	_, err := jwt.ParseToken(token)
	assert.Error(t, err, "expected error for token signed with different secret, got nil")
}
