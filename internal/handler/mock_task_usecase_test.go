package handler_test

import (
	"github.com/MovingPointP/go-task-api/internal/domain/entity"
)

type mockTaskUsecase struct {
	tasks  map[uint]*entity.Task
	nextID uint
}

func NewMockTaskUsecase() *mockTaskUsecase {
	return &mockTaskUsecase{
		tasks:  make(map[uint]*entity.Task),
		nextID: 1,
	}
}

func (m *mockTaskUsecase) Create(userID uint, title, description string) (*entity.Task, error) {
	task := &entity.Task{
		ID:          m.nextID,
		UserID:      userID,
		Title:       title,
		Description: description,
	}
	m.tasks[m.nextID] = task
	m.nextID++
	return task, nil
}

func (m *mockTaskUsecase) Get(taskID, userID uint) (*entity.Task, error) {
	task, ok := m.tasks[taskID]
	if !ok || task.UserID != userID {
		return nil, entity.ErrTaskNotFound
	}
	return task, nil
}

func (m *mockTaskUsecase) GetAll(userID uint) ([]*entity.Task, error) {
	var tasks []*entity.Task
	for _, task := range m.tasks {
		if task.UserID == userID {
			tasks = append(tasks, task)
		}
	}
	return tasks, nil
}

func (m *mockTaskUsecase) Update(taskID, userID uint, title, description string, completed bool) (*entity.Task, error) {
	task, ok := m.tasks[taskID]
	if !ok || task.UserID != userID {
		return nil, entity.ErrTaskNotFound
	}
	task.Title = title
	task.Description = description
	task.Completed = completed
	return task, nil
}

func (m *mockTaskUsecase) Delete(taskID, userID uint) error {
	task, ok := m.tasks[taskID]
	if !ok || task.UserID != userID {
		return entity.ErrTaskNotFound
	}
	delete(m.tasks, taskID)
	return nil
}
