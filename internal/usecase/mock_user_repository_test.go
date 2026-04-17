package usecase_test

import (
	"github.com/MovingPointP/go-task-api/internal/domain/entity"
)

type mockUserRepository struct {
	users  map[string]*entity.User
	nextID uint
}

func newMockUserRepository() *mockUserRepository {
	return &mockUserRepository{
		users:  make(map[string]*entity.User),
		nextID: 1,
	}
}

func (m *mockUserRepository) Create(user *entity.User) error {
	user.ID = m.nextID
	m.nextID++
	m.users[user.Email] = user
	return nil
}

func (m *mockUserRepository) FindByEmail(email string) (*entity.User, error) {
	user, ok := m.users[email]
	if !ok {
		return nil, nil
	}
	return user, nil
}

func (m *mockUserRepository) FindByID(id uint) (*entity.User, error) {
	for _, u := range m.users {
		if u.ID == id {
			return u, nil
		}
	}
	return nil, nil
}
