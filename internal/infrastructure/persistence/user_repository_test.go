package persistence_test

import (
	"testing"

	"github.com/MovingPointP/go-task-api/internal/domain/entity"
	"github.com/MovingPointP/go-task-api/internal/infrastructure/persistence"
)

func TestUserRepository_Create(t *testing.T) {
	repo := persistence.NewUserRepository(setupTestDB(t))

	user := &entity.User{Email: "create@example.com", PasswordHash: "hash"}
	if err := repo.Create(user); err != nil {
		t.Fatalf("Create failed: %v", err)
	}

	if user.ID == 0 {
		t.Error("expected user.ID to be set after Create")
	}
}

func TestUserRepository_FindByEmail(t *testing.T) {
	repo := persistence.NewUserRepository(setupTestDB(t))

	repo.Create(&entity.User{Email: "found@example.com", PasswordHash: "hash"})

	// 存在する場合
	found, err := repo.FindByEmail("found@example.com")
	if err != nil {
		t.Fatalf("FindByEmail failed: %v", err)
	}
	if found == nil {
		t.Fatalf("expected user, got nil")
	}
	if found.Email != "found@example.com" {
		t.Errorf("expected found@example.com, got %s", found.Email)
	}

	// 存在しない場合
	notFound, err := repo.FindByEmail("nobody@example.com")
	if err != nil {
		t.Fatalf("FindByEmail failed: %v", err)
	}

	if notFound != nil {
		t.Error("expected nil for non-existent email")
	}
}

func TestUserRepository_FindByID(t *testing.T) {
	repo := persistence.NewUserRepository(setupTestDB(t))

	created := &entity.User{Email: "byid@example.com", PasswordHash: "hash"}

	repo.Create(created)

	// 存在する場合
	found, err := repo.FindByID(created.ID)
	if err != nil {
		t.Fatalf("FindByID failed: %v", err)
	}
	if found == nil || found.ID != created.ID {
		t.Errorf("expected user with ID %d", created.ID)
	}

	// 存在しない場合
	notFound, err := repo.FindByID(9999)
	if err != nil {
		t.Fatalf("FindByID failed: %v", err)
	}
	if notFound != nil {
		t.Error("expected nil for non-existent ID")
	}
}
