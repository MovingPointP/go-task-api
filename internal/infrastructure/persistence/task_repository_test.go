package persistence_test

import (
	"testing"

	"github.com/MovingPointP/go-task-api/internal/domain/entity"
	"github.com/MovingPointP/go-task-api/internal/infrastructure/persistence"
)

func TestTaskRepository_Create(t *testing.T) {
	repo := persistence.NewTaskRepository(setupTestDB(t))

	task := &entity.Task{UserID: 1, Title: "テストタスク", Description: "説明"}

	if err := repo.Create(task); err != nil {
		t.Fatalf("Create failed: %v", err)
	}

	if task.ID == 0 {
		t.Error("expected task.ID to be set after Create")
	}
}

func TestTaskRepository_FindByID(t *testing.T) {
	repo := persistence.NewTaskRepository(setupTestDB(t))

	task := &entity.Task{UserID: 1, Title: "テストタスク"}
	repo.Create(task)

	// 存在する場合
	found, err := repo.FindByID(task.ID, 1)
	if err != nil {
		t.Fatalf("FindByID failed: %v", err)
	}
	if found == nil {
		t.Error("expected task, got nil")
	}
	if found.Title != "テストタスク" {
		t.Errorf("expected title テストタスク, got %s", found.Title)
	}

	// 存在しない場合
	notFound, err := repo.FindByID(9999, 1)
	if err != nil {
		t.Fatalf("FindByID failed: %v", err)
	}
	if notFound != nil {
		t.Error("expected nil for non-existent ID")
	}

	// 他ユーザーのタスク
	other, err := repo.FindByID(task.ID, 2)
	if err != nil {
		t.Fatalf("FindByID failed: %v", err)
	}
	if other != nil {
		t.Error("expected nil for wrong userID")
	}
}

func TestTaskRepository_FindAllByUserID(t *testing.T) {
	repo := persistence.NewTaskRepository(setupTestDB(t))

	repo.Create(&entity.Task{UserID: 1, Title: "タスク1"})
	repo.Create(&entity.Task{UserID: 1, Title: "タスク2"})
	repo.Create(&entity.Task{UserID: 2, Title: "他ユーザーのタスク"})

	tasks, err := repo.FindAllByUserID(1)
	if err != nil {
		t.Fatalf("FindAllByUserID failed: %v", err)
	}

	if len(tasks) != 2 {
		t.Errorf("expected 2 tasks for userID=1, got %d", len(tasks))
	}
}

func TestTaskRepository_Update(t *testing.T) {
	repo := persistence.NewTaskRepository(setupTestDB(t))

	task := &entity.Task{UserID: 1, Title: "元のタイトル", Completed: false}
	repo.Create(task)

	task.Title = "更新後のタイトル"
	task.Completed = true

	if err := repo.Update(task); err != nil {
		t.Fatalf("Update failed: %v", err)
	}

	updated, _ := repo.FindByID(task.ID, 1)
	if updated.Title != "更新後のタイトル" {
		t.Errorf("expected 更新後のタイトル, got %s", updated.Title)
	}
	if !updated.Completed {
		t.Error("expected Completed to be true after Update")
	}
}

func TestTaskRepository_Delete(t *testing.T) {
	repo := persistence.NewTaskRepository(setupTestDB(t))

	task := &entity.Task{UserID: 1, Title: "削除するタスク"}
	repo.Create(task)

	if err := repo.Delete(task.ID, 1); err != nil {
		t.Fatalf("Delete failed: %v", err)
	}

	deleted, _ := repo.FindByID(task.ID, 1)
	if deleted != nil {
		t.Error("expected nil after Delete")
	}
}
