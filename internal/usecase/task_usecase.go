package usecase

import (
	"errors"
	"fmt"

	"github.com/MovingPointP/go-task-api/internal/domain/entity"
	"github.com/MovingPointP/go-task-api/internal/domain/repository"
)

type TaskUsecase interface {
	Create(userID uint, title, description string) (*entity.Task, error)
	Get(taskID, userID uint) (*entity.Task, error)
	GetAll(userID uint) ([]*entity.Task, error)
	Update(taskID, userID uint, title, description string, completed bool) (*entity.Task, error)
	Delete(taskID, userID uint) error
}

type taskUsecase struct {
	taskRepo repository.TaskRepository
}

// コンストラクタ
func NewTaskUsecase(taskRepo repository.TaskRepository) TaskUsecase {
	return &taskUsecase{taskRepo: taskRepo}
}

func (t *taskUsecase) Create(userID uint, title, description string) (*entity.Task, error) {
	task := &entity.Task{
		UserID:      userID,
		Title:       title,
		Description: description,
	}
	if err := t.taskRepo.Create(task); err != nil {
		return nil, fmt.Errorf("failed to create task: %w", err)
	}
	return task, nil
}

func (t *taskUsecase) Get(taskID, userID uint) (*entity.Task, error) {
	task, err := t.taskRepo.FindByID(taskID, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get task: %w", err)
	}
	if task == nil {
		return nil, errors.New("task not found")
	}
	return task, nil
}

func (t *taskUsecase) GetAll(userID uint) ([]*entity.Task, error) {
	tasks, err := t.taskRepo.FindAllByUserID(userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get tasks: %w", err)
	}
	return tasks, nil
}

func (t *taskUsecase) Update(taskID, userID uint, title, description string, completed bool) (*entity.Task, error) {
	// タスクの取得
	task, err := t.taskRepo.FindByID(taskID, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get task: %w", err)
	}
	if task == nil {
		return nil, errors.New("task not found")
	}

	task.Title = title
	task.Description = description
	task.Completed = completed

	// タスクの更新
	if err := t.taskRepo.Update(task); err != nil {
		return nil, fmt.Errorf("failed to update task: %w", err)
	}
	return task, nil
}

func (t *taskUsecase) Delete(taskID, userID uint) error {
	// タスクの取得
	task, err := t.taskRepo.FindByID(taskID, userID)
	if err != nil {
		return fmt.Errorf("failed to get task: %w", err)
	}
	if task == nil {
		return errors.New("task not found")
	}

	// タスクの削除
	if err := t.taskRepo.Delete(taskID, userID); err != nil {
		return fmt.Errorf("failed to delete task: %w", err)
	}
	return nil
}
