package main

import (
	"fmt"
	"log"
	"os"

	"github.com/MovingPointP/go-task-api/docs"
	"github.com/MovingPointP/go-task-api/internal/handler"
	"github.com/MovingPointP/go-task-api/internal/infrastructure/database"
	"github.com/MovingPointP/go-task-api/internal/infrastructure/persistence"
	"github.com/MovingPointP/go-task-api/internal/usecase"
	"github.com/joho/godotenv"
)

// @title           go-task-api
// @version         1.0
// @description     タスク管理 REST API
// @BasePath        /api/v1
// @securityDefinitions.apikey BearerAuth
// @in              header
// @name            Authorization
func main() {
	// .envの読み込み
	_ = godotenv.Load()

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	host := os.Getenv("HOST")
	if host == "" {
		host = fmt.Sprintf("localhost:%s", port)
	}
	docs.SwaggerInfo.Host = host

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

	log.Printf("Server starting on port %s", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}
