package handler_test

import (
	"github.com/MovingPointP/go-task-api/internal/domain/entity"
)

type mockAuthUsecase struct {
	user  *entity.User
	token string
	err   error
}

func NewMockAuthUsecase(user *entity.User, token string, err error) *mockAuthUsecase {
	return &mockAuthUsecase{
		user:  user,
		token: token,
		err:   err,
	}
}

func (m *mockAuthUsecase) Register(email, password string) (*entity.User, string, error) {
	return m.user, m.token, m.err
}
func (m *mockAuthUsecase) Login(email, password string) (*entity.User, string, error) {
	return m.user, m.token, m.err
}
