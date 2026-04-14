package repository

import "github.com/MovingPointP/go-task-api/internal/domain/entity"

type TaskRepository interface {
	Create(task *entity.Task) error
	FindByID(id, userID uint) (*entity.Task, error)
	FindAllByUserID(userID uint) ([]*entity.Task, error)
	Update(task *entity.Task) error
	Delete(id, userID uint) error
}
