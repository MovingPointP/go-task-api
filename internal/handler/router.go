package handler

import (
	"net/http"

	"github.com/MovingPointP/go-task-api/internal/handler/middleware"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func NewRouter(authHandler *AuthHandler, taskHandler *TaskHandler) *gin.Engine {
	r := gin.Default()

	// CORSヘッダーの設定
	r.Use(func(ctx *gin.Context) {
		ctx.Header("Access-Control-Allow-Origin", "*")
		ctx.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		ctx.Header("Access-Control-Allow-Headers", "Authorization, Content-Type")
		if ctx.Request.Method == "OPTIONS" {
			ctx.AbortWithStatus(http.StatusNoContent)
			return
		}
		ctx.Next()
	})

	// Swagger
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	v1 := r.Group("api/v1")
	{
		// 認証不要
		auth := v1.Group("/auth")
		{
			auth.POST("/register", authHandler.Register)
			auth.POST("/login", authHandler.Login)
		}

		//　認証必要
		tasks := v1.Group("/tasks")
		tasks.Use(middleware.AuthMiddleware())
		{
			tasks.POST("", taskHandler.CreateTask)
			tasks.GET("", taskHandler.GetTasks)
			tasks.GET("/:id", taskHandler.GetTask)
			tasks.PUT("/:id", taskHandler.UpdateTask)
			tasks.DELETE("/:id", taskHandler.DeleteTask)
		}
	}
	return r
}
