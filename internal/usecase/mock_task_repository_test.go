package usecase_test

import "github.com/MovingPointP/go-task-api/internal/domain/entity"

type mockTaskRepository struct {
	tasks  map[uint]*entity.Task
	nextID uint
}

func NewMockTaskRepository() *mockTaskRepository {
	return &mockTaskRepository{
		tasks:  make(map[uint]*entity.Task),
		nextID: 1,
	}
}

func (m *mockTaskRepository) Create(task *entity.Task) error {
	task.ID = m.nextID
	m.tasks[m.nextID] = task
	m.nextID++

	return nil
}

func (m *mockTaskRepository) FindByID(id, userID uint) (*entity.Task, error) {
	task, ok := m.tasks[id]
	if !ok || task.UserID != userID {
		return nil, nil
	}

	return task, nil
}

func (m *mockTaskRepository) FindAllByUserID(userID uint) ([]*entity.Task, error) {
	var tasks []*entity.Task

	for _, task := range m.tasks {
		if task.UserID == userID {
			tasks = append(tasks, task)
		}
	}

	return tasks, nil
}

func (m *mockTaskRepository) Update(task *entity.Task) error {
	_, ok := m.tasks[task.ID]
	if !ok {
		return nil
	}

	m.tasks[task.ID] = task

	return nil
}

func (m *mockTaskRepository) Delete(id, userID uint) error {
	task, ok := m.tasks[id]
	if !ok || task.UserID != userID {
		return nil
	}

	delete(m.tasks, id)

	return nil
}
