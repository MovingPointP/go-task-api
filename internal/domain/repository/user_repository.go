package repository

import "github.com/MovingPointP/go-task-api/internal/domain/entity"

type UserRepository interface {
	Create(user *entity.User) error
	FindByEmail(email string) (*entity.User, error)
	FindByID(id uint) (*entity.User, error)
}
