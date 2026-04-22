package persistence_test

import (
	"testing"

	"github.com/MovingPointP/go-task-api/internal/domain/entity"
	"github.com/MovingPointP/go-task-api/internal/infrastructure/persistence"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUserRepository_Create(t *testing.T) {
	repo := persistence.NewUserRepository(setupTestDB(t))

	t.Run("正常系", func(t *testing.T) {
		user := &entity.User{Email: "create@example.com", PasswordHash: "hash"}
		err := repo.Create(user)
		require.NoError(t, err, "Create failed")
		assert.NotZero(t, user.ID)
	})
}

func TestUserRepository_FindByEmail(t *testing.T) {
	repo := persistence.NewUserRepository(setupTestDB(t))
	repo.Create(&entity.User{Email: "found@example.com", PasswordHash: "hash"})

	t.Run("存在する場合", func(t *testing.T) {
		found, err := repo.FindByEmail("found@example.com")
		require.NoError(t, err, "FindByEmail failed")
		require.NotNil(t, found)
		assert.Equal(t, "found@example.com", found.Email)
	})

	t.Run("存在しない場合", func(t *testing.T) {
		notFound, err := repo.FindByEmail("nobody@example.com")
		require.NoError(t, err, "FindByEmail failed")
		assert.Nil(t, notFound)
	})
}

func TestUserRepository_FindByID(t *testing.T) {
	repo := persistence.NewUserRepository(setupTestDB(t))
	created := &entity.User{Email: "byid@example.com", PasswordHash: "hash"}
	repo.Create(created)

	t.Run("存在する場合", func(t *testing.T) {
		found, err := repo.FindByID(created.ID)
		require.NoError(t, err, "FindByID failed")
		require.NotNil(t, found)
		assert.Equal(t, created.ID, found.ID)
	})

	t.Run("存在しない場合", func(t *testing.T) {
		notFound, err := repo.FindByID(9999)
		require.NoError(t, err, "FindByID failed")
		assert.Nil(t, notFound)
	})
}
