package usecase_test

import (
	"github.com/MovingPointP/go-task-api/internal/domain/entity"
)

type mockUserRepository struct {
	users  map[uint]*entity.User
	nextID uint
}

func newMockUserRepository() *mockUserRepository {
	return &mockUserRepository{
		users:  make(map[uint]*entity.User),
		nextID: 1,
	}
}

func (m *mockUserRepository) Create(user *entity.User) error {
	user.ID = m.nextID
	m.users[user.ID] = user
	m.nextID++
	return nil
}

func (m *mockUserRepository) FindByEmail(email string) (*entity.User, error) {
	for _, u := range m.users {
		if u.Email == email {
			return u, nil
		}
	}
	return nil, nil
}

func (m *mockUserRepository) FindByID(id uint) (*entity.User, error) {
	user, ok := m.users[id]
	if !ok {
		return nil, nil
	}
	return user, nil
}
