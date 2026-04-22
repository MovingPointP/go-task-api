package usecase_test

import (
	"os"
	"testing"

	"github.com/MovingPointP/go-task-api/internal/usecase"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAuthUsecase_Register(t *testing.T) {
	os.Setenv("JWT_SECRET", "test-secret")
	os.Setenv("JWT_EXPIRATION_HOURS", "1")

	uc := usecase.NewAuthUsecase(newMockUserRepository())

	// 正常系
	user, token, err := uc.Register("test@example.com", "test")
	require.NoError(t, err, "Register failed")
	assert.Equal(t, "test@example.com", user.Email)
	assert.NotEmpty(t, token)

	// メールアドレスの重複テスト
	_, _, err = uc.Register("test@example.com", "test")
	assert.Error(t, err, "expected error for duplicate email, got nil")
}

func TestAuthUsecase_Login(t *testing.T) {
	os.Setenv("JWT_SECRET", "test-secret")
	os.Setenv("JWT_EXPIRATION_HOURS", "1")

	uc := usecase.NewAuthUsecase(newMockUserRepository())

	// ユーザー登録
	_, _, err := uc.Register("login@example.com", "test")
	require.NoError(t, err, "Register failed")

	// 正常系
	user, token, err := uc.Login("login@example.com", "test")
	require.NoError(t, err, "Login failed")
	assert.Equal(t, "login@example.com", user.Email)
	assert.NotEmpty(t, token)

	// 誤ったパスワードのテスト
	_, _, err = uc.Login("login@example.com", "wrong")
	assert.Error(t, err, "expected error for wrong password, got nil")

	// 存在しない場合
	_, _, err = uc.Login("nobody@example.com", "test")
	assert.Error(t, err, "expected error for non-exsistent user, got nil")
}
