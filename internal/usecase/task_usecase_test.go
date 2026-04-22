package usecase_test

import (
	"testing"

	"github.com/MovingPointP/go-task-api/internal/usecase"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTaskUsecase_CreateTask(t *testing.T) {
	tc := usecase.NewTaskUsecase(NewMockTaskRepository())

	t.Run("正常系", func(t *testing.T) {
		task, err := tc.Create(1, "タイトル", "説明")
		require.NoError(t, err, "Create failed")
		assert.Equal(t, uint(1), task.UserID)
		assert.Equal(t, "タイトル", task.Title)
		assert.Equal(t, "説明", task.Description)
		assert.False(t, task.Completed)
	})
}

func TestTaskUsecase_GetTask(t *testing.T) {
	tc := usecase.NewTaskUsecase(NewMockTaskRepository())
	task, _ := tc.Create(1, "タイトル", "説明")

	t.Run("正常系", func(t *testing.T) {
		found, err := tc.Get(task.ID, 1)
		require.NoError(t, err, "Get failed")
		assert.Equal(t, uint(1), found.UserID)
		assert.Equal(t, "タイトル", found.Title)
		assert.Equal(t, "説明", found.Description)
		assert.False(t, found.Completed)
	})

	t.Run("存在しない場合", func(t *testing.T) {
		_, err := tc.Get(9999, 1)
		assert.Error(t, err, "expected error for non-exsistent task, got nil")
	})

	t.Run("他ユーザーのタスク", func(t *testing.T) {
		_, err := tc.Get(task.ID, 2)
		assert.Error(t, err, "expected error for other user task, got nil")
	})
}

func TestTaskUsecase_GetAll(t *testing.T) {
	tc := usecase.NewTaskUsecase(NewMockTaskRepository())

	tc.Create(1, "タイトル1", "")
	tc.Create(1, "タイトル2", "")
	tc.Create(2, "タイトル3", "")

	t.Run("正常系", func(t *testing.T) {
		tasks, err := tc.GetAll(1)
		require.NoError(t, err, "GetAll failed")
		assert.Len(t, tasks, 2)
	})
}

func TestTaskUsecase_Update(t *testing.T) {
	tc := usecase.NewTaskUsecase(NewMockTaskRepository())
	task, _ := tc.Create(1, "元のタイトル", "元の説明")

	t.Run("正常系", func(t *testing.T) {
		updated, err := tc.Update(task.ID, 1, "更新後のタイトル", "更新後の説明", true)
		require.NoError(t, err, "Update failed")
		assert.Equal(t, "更新後のタイトル", updated.Title)
		assert.Equal(t, "更新後の説明", updated.Description)
		assert.True(t, updated.Completed)
	})

	t.Run("他ユーザーのタスク", func(t *testing.T) {
		_, err := tc.Update(task.ID, 2, "他ユーザーのタイトル", "他ユーザーの説明", true)
		assert.Error(t, err, "expected error for other user task, got nil")
	})
}

func TestTaskUsecase_Delete(t *testing.T) {
	tc := usecase.NewTaskUsecase(NewMockTaskRepository())
	task, _ := tc.Create(1, "タイトル", "説明")

	t.Run("他ユーザーのタスク", func(t *testing.T) {
		err := tc.Delete(task.ID, 2)
		assert.Error(t, err, "expected error for other user task, got nil")
	})

	t.Run("正常系", func(t *testing.T) {
		err := tc.Delete(task.ID, 1)
		require.NoError(t, err, "Delete failed")
		deleted, _ := tc.Get(task.ID, 1)
		assert.Nil(t, deleted)
	})
}
