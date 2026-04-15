package persistence_test

import (
	"testing"

	"github.com/MovingPointP/go-task-api/internal/domain/entity"
	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func setupTestDB(t *testing.T) *gorm.DB {
	t.Helper()

	config := &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	}
	db, err := gorm.Open(sqlite.Open(":memory:"), config)

	if err != nil {
		t.Fatalf("failed to open in-memory db: %v", err)
	}

	if err := db.AutoMigrate(&entity.User{}, &entity.Task{}); err != nil {
		t.Fatalf("failed to migrate: %v", err)
	}
	return db
}
