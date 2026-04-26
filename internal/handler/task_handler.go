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
	return ctx.MustGet("UserID").(uint)
}

// @Summary     タスク作成
// @Description 新しいタスクを作成する
// @Tags        tasks
// @Security    BearerAuth
// @Accept      json
// @Produce     json
// @Param       body body CreateTaskRequest true "タスク情報"
// @Success     201 {object} entity.Task
// @Failure     400 {object} map[string]string
// @Failure     500 {object} map[string]string
// @Router      /tasks [post]
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

// @Summary     タスク一覧取得
// @Description ログインユーザーのタスクを全件取得する
// @Tags        tasks
// @Security    BearerAuth
// @Produce     json
// @Success     200 {array} entity.Task
// @Failure     500 {object} map[string]string
// @Router      /tasks [get]
func (h *TaskHandler) GetTasks(ctx *gin.Context) {
	tasks, err := h.taskUsecase.GetAll(getUserID(ctx))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get tasks"})
		return
	}

	ctx.JSON(http.StatusOK, tasks)
}

// @Summary     タスク取得
// @Description 指定 ID のタスクを取得する
// @Tags        tasks
// @Security    BearerAuth
// @Produce     json
// @Param       id path int true "タスクID"
// @Success     200 {object} entity.Task
// @Failure     400 {object} map[string]string
// @Failure     404 {object} map[string]string
// @Router      /tasks/{id} [get]
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

// @Summary     タスク更新
// @Description 指定 ID のタスクを更新する
// @Tags        tasks
// @Security    BearerAuth
// @Accept      json
// @Produce     json
// @Param       id   path int               true "タスクID"
// @Param       body body UpdateTaskRequest true "更新情報"
// @Success     200 {object} entity.Task
// @Failure     400 {object} map[string]string
// @Failure     404 {object} map[string]string
// @Router      /tasks/{id} [put]
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

// @Summary     タスク削除
// @Description 指定 ID のタスクを削除する
// @Tags        tasks
// @Security    BearerAuth
// @Produce     json
// @Param       id path int true "タスクID"
// @Success     204
// @Failure     400 {object} map[string]string
// @Failure     404 {object} map[string]string
// @Router      /tasks/{id} [delete]
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
