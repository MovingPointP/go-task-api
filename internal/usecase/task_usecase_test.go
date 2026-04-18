package usecase_test

import (
	"testing"

	"github.com/MovingPointP/go-task-api/internal/usecase"
)

func TestTaskUsecase_CreateTask(t *testing.T) {
	tc := usecase.NewTaskUsecase(NewMockTaskRepository())

	// 正常系
	task, err := tc.Create(1, "タイトル", "説明")
	if err != nil {
		t.Fatalf("Create failed: %v", err)
	}
	if task.UserID != 1 {
		t.Errorf("expected userID 1, got %d", task.UserID)
	}
	if task.Title != "タイトル" {
		t.Errorf("expected title タイトル, got %s", task.Title)
	}
	if task.Description != "説明" {
		t.Errorf("expected description 説明, got %s", task.Description)
	}
	if task.Completed {
		t.Error("expected Completed false, got true")
	}
}

func TestTaskUsecase_GetTask(t *testing.T) {
	tc := usecase.NewTaskUsecase(NewMockTaskRepository())

	task, _ := tc.Create(1, "タイトル", "説明")

	// 正常系
	found, err := tc.Get(task.ID, 1)
	if err != nil {
		t.Fatalf("Get failed: %v", err)
	}
	if found.UserID != 1 {
		t.Errorf("expected userID 1, got %d", found.UserID)
	}
	if found.Title != "タイトル" {
		t.Errorf("expected title タイトル, got %s", found.Title)
	}
	if found.Description != "説明" {
		t.Errorf("expected description 説明, got %s", found.Description)
	}
	if task.Completed {
		t.Error("expected Completed false, got true")
	}

	// 存在しない場合
	_, err = tc.Get(9999, 1)
	if err == nil {
		t.Error("expected error for non-exsistent task, got nil")
	}

	// 他ユーザーのタスク
	_, err = tc.Get(task.ID, 2)
	if err == nil {
		t.Error("expected error for other user task, got nil")
	}
}

func TestTaskUsecase_GetAll(t *testing.T) {
	tc := usecase.NewTaskUsecase(NewMockTaskRepository())

	tc.Create(1, "タイトル1", "")
	tc.Create(1, "タイトル2", "")
	tc.Create(2, "タイトル3", "")

	// 正常系
	tasks, err := tc.GetAll(1)
	if err != nil {
		t.Fatalf("GetAll failed: %v", err)
	}

	if len(tasks) != 2 {
		t.Errorf("expected 2 tasks for userID=1, got %d", len(tasks))
	}
}

func TestTaskUsecase_Update(t *testing.T) {
	tc := usecase.NewTaskUsecase(NewMockTaskRepository())

	task, _ := tc.Create(1, "元のタイトル", "元の説明")

	// 正常系
	updated, err := tc.Update(task.ID, 1, "更新後のタイトル", "更新後の説明", true)
	if err != nil {
		t.Fatalf("Update failed: %v", err)
	}
	if updated.Title != "更新後のタイトル" {
		t.Errorf("expected title 更新後のタイトル, got %s", updated.Title)
	}
	if updated.Description != "更新後の説明" {
		t.Errorf("expected description 更新後の説明, got %s", updated.Description)
	}
	if !updated.Completed {
		t.Error("expected completed true, got false")
	}

	// 他ユーザーのタスク
	_, err = tc.Update(task.ID, 2, "他ユーザーのタイトル", "他ユーザーの説明", true)
	if err == nil {
		t.Error("expected error for other user task, got nil")
	}
}

func TestTaskUsecase_Delete(t *testing.T) {
	tc := usecase.NewTaskUsecase(NewMockTaskRepository())

	task, _ := tc.Create(1, "タイトル", "説明")

	// 他ユーザーのタスク
	if err := tc.Delete(task.ID, 2); err == nil {
		t.Error("expected error for other user task, got nil")
	}

	// 正常系
	if err := tc.Delete(task.ID, 1); err != nil {
		t.Fatalf("Delete failed: %v", err)
	}
	deleted, _ := tc.Get(task.ID, 1)
	if deleted != nil {
		t.Error("expected nil after Delete")
	}
}
