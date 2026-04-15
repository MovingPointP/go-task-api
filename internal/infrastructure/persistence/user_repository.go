package persistence

import (
	"errors"

	"github.com/MovingPointP/go-task-api/internal/domain/entity"
	"github.com/MovingPointP/go-task-api/internal/domain/repository"
	"gorm.io/gorm"
)

type gormUserRepository struct {
	db *gorm.DB
}

// コンストラクタ
func NewUserRepository(db *gorm.DB) repository.UserRepository {
	return &gormUserRepository{db: db}
}

func (r *gormUserRepository) Create(user *entity.User) error {
	return r.db.Create(user).Error
}

func (r *gormUserRepository) FindByEmail(email string) (*entity.User, error) {
	var user entity.User
	if err := r.db.Where("email = ?", email).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func (r *gormUserRepository) FindByID(id uint) (*entity.User, error) {
	var user entity.User
	if err := r.db.First(&user, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}
