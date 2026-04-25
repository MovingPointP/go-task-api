package main

import (
	"log"
	"os"

	"github.com/MovingPointP/go-task-api/internal/handler"
	"github.com/MovingPointP/go-task-api/internal/infrastructure/database"
	"github.com/MovingPointP/go-task-api/internal/infrastructure/persistence"
	"github.com/MovingPointP/go-task-api/internal/usecase"
	"github.com/joho/godotenv"
)

func main() {
	// .envの読み込み
	_ = godotenv.Load()

	db, err := database.NewDB()
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}

	// DI
	userRepo := persistence.NewUserRepository(db)
	taskRepo := persistence.NewTaskRepository(db)

	authUsecase := usecase.NewAuthUsecase(userRepo)
	taskUsecase := usecase.NewTaskUsecase(taskRepo)

	authHandler := handler.NewAuthHandler(authUsecase)
	taskHandler := handler.NewTaskHandler(taskUsecase)

	// ルーター
	r := handler.NewRouter(authHandler, taskHandler)

	// サーバー起動
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server starting on port %s", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}
