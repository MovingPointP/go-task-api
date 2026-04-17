package persistence

import (
	"errors"

	"github.com/MovingPointP/go-task-api/internal/domain/entity"
	"github.com/MovingPointP/go-task-api/internal/domain/repository"
	"gorm.io/gorm"
)

type gormTaskRepository struct {
	db *gorm.DB
}

// コンストラクタ
func NewTaskRepository(db *gorm.DB) repository.TaskRepository {
	return &gormTaskRepository{db: db}
}

func (r *gormTaskRepository) Create(task *entity.Task) error {
	return r.db.Create(task).Error
}

func (r *gormTaskRepository) FindByID(id, userID uint) (*entity.Task, error) {
	var task entity.Task

	if err := r.db.Where("id = ? AND user_id = ?", id, userID).First(&task).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &task, nil
}

func (r *gormTaskRepository) FindAllByUserID(userID uint) ([]*entity.Task, error) {
	var tasks []*entity.Task

	if err := r.db.Where("user_id = ?", userID).Order("created_at DESC").Find(&tasks).Error; err != nil {
		return nil, err
	}
	return tasks, nil
}

func (r *gormTaskRepository) Update(task *entity.Task) error {
	return r.db.Save(task).Error
}

func (r *gormTaskRepository) Delete(id, userID uint) error {
	return r.db.Where("id = ? AND user_id = ?", id, userID).Delete(&entity.Task{}).Error
}
