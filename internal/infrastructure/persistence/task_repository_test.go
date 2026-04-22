package persistence_test

import (
	"testing"

	"github.com/MovingPointP/go-task-api/internal/domain/entity"
	"github.com/MovingPointP/go-task-api/internal/infrastructure/persistence"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTaskRepository_Create(t *testing.T) {
	repo := persistence.NewTaskRepository(setupTestDB(t))

	t.Run("正常系", func(t *testing.T) {
		task := &entity.Task{UserID: 1, Title: "テストタスク", Description: "説明"}
		err := repo.Create(task)
		require.NoError(t, err, "Create failed")
		assert.NotZero(t, task.ID)
	})
}

func TestTaskRepository_FindByID(t *testing.T) {
	repo := persistence.NewTaskRepository(setupTestDB(t))
	task := &entity.Task{UserID: 1, Title: "テストタスク"}
	repo.Create(task)

	t.Run("存在する場合", func(t *testing.T) {
		found, err := repo.FindByID(task.ID, 1)
		require.NoError(t, err, "FindByID failed")
		require.NotNil(t, found)
		assert.Equal(t, "テストタスク", found.Title)
	})

	t.Run("存在しない場合", func(t *testing.T) {
		notFound, err := repo.FindByID(9999, 1)
		require.NoError(t, err, "FindByID failed")
		assert.Nil(t, notFound)
	})

	t.Run("他ユーザーのタスク", func(t *testing.T) {
		other, err := repo.FindByID(task.ID, 2)
		require.NoError(t, err, "FindByID failed")
		assert.Nil(t, other)
	})
}

func TestTaskRepository_FindAllByUserID(t *testing.T) {
	repo := persistence.NewTaskRepository(setupTestDB(t))

	repo.Create(&entity.Task{UserID: 1, Title: "タスク1"})
	repo.Create(&entity.Task{UserID: 1, Title: "タスク2"})
	repo.Create(&entity.Task{UserID: 2, Title: "他ユーザーのタスク"})

	t.Run("正常系", func(t *testing.T) {
		tasks, err := repo.FindAllByUserID(1)
		require.NoError(t, err, "FindAllByUserID failed")
		assert.Len(t, tasks, 2)
	})
}

func TestTaskRepository_Update(t *testing.T) {
	repo := persistence.NewTaskRepository(setupTestDB(t))

	task := &entity.Task{UserID: 1, Title: "元のタイトル", Completed: false}
	repo.Create(task)

	task.Title = "更新後のタイトル"
	task.Completed = true

	t.Run("正常系", func(t *testing.T) {
		err := repo.Update(task)
		require.NoError(t, err, "Update failed")

		updated, _ := repo.FindByID(task.ID, 1)
		assert.Equal(t, "更新後のタイトル", updated.Title)
		assert.True(t, updated.Completed)
	})
}

func TestTaskRepository_Delete(t *testing.T) {
	repo := persistence.NewTaskRepository(setupTestDB(t))

	task := &entity.Task{UserID: 1, Title: "削除するタスク"}
	repo.Create(task)

	t.Run("正常系", func(t *testing.T) {
		err := repo.Delete(task.ID, 1)
		require.NoError(t, err, "Delete failed")

		deleted, _ := repo.FindByID(task.ID, 1)
		assert.Nil(t, deleted)
	})
}
