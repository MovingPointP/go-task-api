package handler

import (
	"net/http"
	"strconv"

	"github.com/MovingPointP/go-task-api/internal/usecase"
	"github.com/gin-gonic/gin"
)

type TaskHandler struct {
	taskUsecase usecase.TaskUsecase
}

func NewTaskHandler(taskusecase usecase.TaskUsecase) *TaskHandler {
	return &TaskHandler{taskUsecase: taskusecase}
}

// タスク作成のリクエストボディ
type CreateTaskRequest struct {
	Title       string `json:"title" binding:"required"`
	Description string `json:"description"`
}

// タスク更新のリクエストボディ
type UpdateTaskRequest struct {
	Title       string `json:"title" binding:"required"`
	Description string `json:"description"`
	Completed   bool   `json:"completed"`
}

func getUserID(ctx *gin.Context) uint {
	return ctx.MustGet("userID").(uint)
}

func (h *TaskHandler) CreateTask(ctx *gin.Context) {
	var req CreateTaskRequest
	if err := ctx.ShouldBind(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	task, err := h.taskUsecase.Create(getUserID(ctx), req.Title, req.Description)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create task"})
		return
	}

	ctx.JSON(http.StatusCreated, task)
}

func (h *TaskHandler) GetTasks(ctx *gin.Context) {
	tasks, err := h.taskUsecase.GetAll(getUserID(ctx))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get tasks"})
		return
	}

	ctx.JSON(http.StatusOK, tasks)
}

func (h *TaskHandler) GetTask(ctx *gin.Context) {
	taskID, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid task id"})
		return
	}

	task, err := h.taskUsecase.Get(uint(taskID), getUserID(ctx))
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, task)
}

func (h *TaskHandler) UpdateTask(ctx *gin.Context) {
	taskID, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid task id"})
		return
	}

	var req UpdateTaskRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	task, err := h.taskUsecase.Update(uint(taskID), getUserID(ctx), req.Title, req.Description, req.Completed)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, task)
}

func (h *TaskHandler) DeleteTask(ctx *gin.Context) {
	taskID, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid task id"})
		return
	}

	if err := h.taskUsecase.Delete(uint(taskID), getUserID(ctx)); err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusNoContent, nil)
}
