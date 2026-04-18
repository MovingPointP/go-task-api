package usecase_test

import (
	"os"
	"testing"

	"github.com/MovingPointP/go-task-api/internal/usecase"
)

func TestAuthUsecase_Register(t *testing.T) {
	os.Setenv("JWT_SECRET", "test-secret")
	os.Setenv("JWT_EXPIRATION_HOURS", "1")

	uc := usecase.NewAuthUsecase(newMockUserRepository())

	// 正常系
	user, token, err := uc.Register("test@example.com", "test")
	if err != nil {
		t.Fatalf("Register failed")
	}
	if user.Email != "test@example.com" {
		t.Errorf("expected email test@example.com, got %s", user.Email)
	}
	if token == "" {
		t.Error("expected non-empty token")
	}

	// メールアドレスの重複テスト
	_, _, err = uc.Register("test@example.com", "test")
	if err == nil {
		t.Error("expected error for duplicate email, got nil")
	}
}

func TestAuthUsecase_Login(t *testing.T) {
	os.Setenv("JWT_SECRET", "test-secret")
	os.Setenv("JWT_EXPIRATION_HOURS", "1")

	uc := usecase.NewAuthUsecase(newMockUserRepository())

	// ユーザー登録
	_, _, err := uc.Register("login@example.com", "test")
	if err != nil {
		t.Fatalf("Register failed: %v", err)
	}

	// 正常系
	user, token, err := uc.Login("login@example.com", "test")
	if err != nil {
		t.Fatalf("Login failed: %v", err)
	}
	if user.Email != "login@example.com" {
		t.Errorf("expected email login@example.com, got %s", user.Email)
	}
	if token == "" {
		t.Error("expected non-empty token")
	}

	// 誤ったパスワードのテスト
	_, _, err = uc.Login("login@example.com", "wrong")
	if err == nil {
		t.Error("expected error for wrong password, got nil")
	}

	// 存在しない場合
	_, _, err = uc.Login("nobody@example.com", "test")
	if err == nil {
		t.Error("expected error for non-exsistent user, got nil")
	}
}
