package usecase

import (
	"fmt"

	"github.com/MovingPointP/go-task-api/internal/domain/entity"
	"github.com/MovingPointP/go-task-api/internal/domain/repository"
	"github.com/MovingPointP/go-task-api/pkg/jwt"
	"golang.org/x/crypto/bcrypt"
)

type AuthUsecase interface {
	Register(email, password string) (*entity.User, string, error)
	Login(email, password string) (*entity.User, string, error)
}

type authUsecase struct {
	userRepo repository.UserRepository
}

// コンストラクタ
func NewAuthUsecase(userRepo repository.UserRepository) AuthUsecase {
	return &authUsecase{userRepo: userRepo}
}

func (u *authUsecase) Register(email, password string) (*entity.User, string, error) {
	// メールアドレスの重複チェック
	exsisting, err := u.userRepo.FindByEmail(email)
	if err != nil {
		return nil, "", fmt.Errorf("failed to check email: %w", err)
	}
	if exsisting != nil {
		return nil, "", entity.ErrEmailAlreadyInUse
	}

	// パスワードのハッシュ化
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, "", fmt.Errorf("failed to hash password: %w", err)
	}

	// ユーザーの作成
	user := &entity.User{
		Email:        email,
		PasswordHash: string(hash),
	}
	if err := u.userRepo.Create(user); err != nil {
		return nil, "", fmt.Errorf("failed to create user: %w", err)
	}

	// トークンの生成
	token, err := jwt.GenerateToken(user.ID)
	if err != nil {
		return nil, "", fmt.Errorf("failed to generate token: %w", err)
	}

	return user, token, nil

}

func (u *authUsecase) Login(email, password string) (*entity.User, string, error) {
	user, err := u.userRepo.FindByEmail(email)
	if err != nil {
		return nil, "", fmt.Errorf("failed to find user: %w", err)
	}
	if user == nil {
		return nil, "", entity.ErrInvalidCredentials
	}

	// パスワードの検証
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return nil, "", entity.ErrInvalidCredentials
	}

	// トークンの生成
	token, err := jwt.GenerateToken(user.ID)
	if err != nil {
		return nil, "", fmt.Errorf("failed to generate token: %w", err)
	}

	return user, token, nil
}
