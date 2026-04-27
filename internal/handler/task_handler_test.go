package handler_test

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/MovingPointP/go-task-api/internal/domain/entity"
	"github.com/MovingPointP/go-task-api/internal/handler"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const testUserID = uint(1)

func newTaskContext(method, url, body string) (*gin.Context, *httptest.ResponseRecorder) {
	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)
	var req *http.Request
	if body != "" {
		req, _ = http.NewRequest(method, url, strings.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
	} else {
		req, _ = http.NewRequest(method, url, nil)
	}
	ctx.Request = req
	ctx.Set("UserID", testUserID)
	return ctx, w
}

func TestTaskHandler_CreateTask(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("正常系", func(t *testing.T) {
		mock := NewMockTaskUsecase()
		h := handler.NewTaskHandler(mock)

		ctx, w := newTaskContext("POST", "/tasks", `{"title":"タイトル","description":"説明"}`)
		h.CreateTask(ctx)

		require.Equal(t, http.StatusCreated, w.Code)
		var task entity.Task
		require.NoError(t, json.Unmarshal(w.Body.Bytes(), &task))
		assert.Equal(t, "タイトル", task.Title)
		assert.Equal(t, "説明", task.Description)
		assert.Equal(t, testUserID, task.UserID)
	})

	t.Run("バリデーションエラー", func(t *testing.T) {
		mock := NewMockTaskUsecase()
		h := handler.NewTaskHandler(mock)

		ctx, w := newTaskContext("POST", "/tasks", `{"description":"説明"}`)
		h.CreateTask(ctx)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}

func TestTaskHandler_GetTasks(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("正常系", func(t *testing.T) {
		mock := NewMockTaskUsecase()
		mock.Create(testUserID, "タイトル1", "")
		mock.Create(testUserID, "タイトル2", "")
		mock.Create(2, "他ユーザーのタスク", "")
		h := handler.NewTaskHandler(mock)

		ctx, w := newTaskContext("GET", "/tasks", "")
		h.GetTasks(ctx)

		require.Equal(t, http.StatusOK, w.Code)
		var tasks []*entity.Task
		require.NoError(t, json.Unmarshal(w.Body.Bytes(), &tasks))
		assert.Len(t, tasks, 2)
	})
}

func TestTaskHandler_GetTask(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("正常系", func(t *testing.T) {
		mock := NewMockTaskUsecase()
		created, _ := mock.Create(testUserID, "タイトル", "説明")
		h := handler.NewTaskHandler(mock)

		ctx, w := newTaskContext("GET", "/tasks/1", "")
		ctx.Params = gin.Params{{Key: "id", Value: "1"}}
		h.GetTask(ctx)

		require.Equal(t, http.StatusOK, w.Code)
		var task entity.Task
		require.NoError(t, json.Unmarshal(w.Body.Bytes(), &task))
		assert.Equal(t, created.ID, task.ID)
		assert.Equal(t, "タイトル", task.Title)
	})

	t.Run("不正なタスクID", func(t *testing.T) {
		mock := NewMockTaskUsecase()
		h := handler.NewTaskHandler(mock)

		ctx, w := newTaskContext("GET", "/tasks/abc", "")
		ctx.Params = gin.Params{{Key: "id", Value: "abc"}}
		h.GetTask(ctx)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("存在しない", func(t *testing.T) {
		mock := NewMockTaskUsecase()
		h := handler.NewTaskHandler(mock)

		ctx, w := newTaskContext("GET", "/tasks/9999", "")
		ctx.Params = gin.Params{{Key: "id", Value: "9999"}}
		h.GetTask(ctx)

		assert.Equal(t, http.StatusNotFound, w.Code)
	})

	t.Run("他ユーザーのタスク", func(t *testing.T) {
		mock := NewMockTaskUsecase()
		mock.Create(2, "他ユーザーのタスク", "")
		h := handler.NewTaskHandler(mock)

		ctx, w := newTaskContext("GET", "/tasks/1", "")
		ctx.Params = gin.Params{{Key: "id", Value: "1"}}
		h.GetTask(ctx)

		assert.Equal(t, http.StatusNotFound, w.Code)
	})

	t.Run("DBエラー", func(t *testing.T) {
		h := handler.NewTaskHandler(&errTaskUsecase{err: errors.New("db error")})

		ctx, w := newTaskContext("GET", "/tasks/1", "")
		ctx.Params = gin.Params{{Key: "id", Value: "1"}}
		h.GetTask(ctx)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})
}

func TestTaskHandler_UpdateTask(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("正常系", func(t *testing.T) {
		mock := NewMockTaskUsecase()
		mock.Create(testUserID, "元のタイトル", "元の説明")
		h := handler.NewTaskHandler(mock)

		ctx, w := newTaskContext("PUT", "/tasks/1", `{"title":"更新後のタイトル","description":"更新後の説明","completed":true}`)
		ctx.Params = gin.Params{{Key: "id", Value: "1"}}
		h.UpdateTask(ctx)

		require.Equal(t, http.StatusOK, w.Code)
		var task entity.Task
		require.NoError(t, json.Unmarshal(w.Body.Bytes(), &task))
		assert.Equal(t, "更新後のタイトル", task.Title)
		assert.Equal(t, "更新後の説明", task.Description)
		assert.True(t, task.Completed)
	})

	t.Run("不正なタスクID", func(t *testing.T) {
		mock := NewMockTaskUsecase()
		h := handler.NewTaskHandler(mock)

		ctx, w := newTaskContext("PUT", "/tasks/abc", `{"title":"タイトル"}`)
		ctx.Params = gin.Params{{Key: "id", Value: "abc"}}
		h.UpdateTask(ctx)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("バリデーションエラー", func(t *testing.T) {
		mock := NewMockTaskUsecase()
		h := handler.NewTaskHandler(mock)

		ctx, w := newTaskContext("PUT", "/tasks/1", `{"description":"説明のみ"}`)
		ctx.Params = gin.Params{{Key: "id", Value: "1"}}
		h.UpdateTask(ctx)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("存在しない", func(t *testing.T) {
		mock := NewMockTaskUsecase()
		h := handler.NewTaskHandler(mock)

		ctx, w := newTaskContext("PUT", "/tasks/9999", `{"title":"タイトル"}`)
		ctx.Params = gin.Params{{Key: "id", Value: "9999"}}
		h.UpdateTask(ctx)

		assert.Equal(t, http.StatusNotFound, w.Code)
	})

	t.Run("DBエラー", func(t *testing.T) {
		h := handler.NewTaskHandler(&errTaskUsecase{err: errors.New("db error")})

		ctx, w := newTaskContext("PUT", "/tasks/1", `{"title":"タイトル"}`)
		ctx.Params = gin.Params{{Key: "id", Value: "1"}}
		h.UpdateTask(ctx)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})
}

func TestTaskHandler_DeleteTask(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("正常系", func(t *testing.T) {
		mock := NewMockTaskUsecase()
		mock.Create(testUserID, "タイトル", "")
		h := handler.NewTaskHandler(mock)

		ctx, w := newTaskContext("DELETE", "/tasks/1", "")
		ctx.Params = gin.Params{{Key: "id", Value: "1"}}
		h.DeleteTask(ctx)

		assert.Equal(t, http.StatusNoContent, w.Code)
	})

	t.Run("不正なタスクID", func(t *testing.T) {
		mock := NewMockTaskUsecase()
		h := handler.NewTaskHandler(mock)

		ctx, w := newTaskContext("DELETE", "/tasks/abc", "")
		ctx.Params = gin.Params{{Key: "id", Value: "abc"}}
		h.DeleteTask(ctx)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("存在しない", func(t *testing.T) {
		mock := NewMockTaskUsecase()
		h := handler.NewTaskHandler(mock)

		ctx, w := newTaskContext("DELETE", "/tasks/9999", "")
		ctx.Params = gin.Params{{Key: "id", Value: "9999"}}
		h.DeleteTask(ctx)

		assert.Equal(t, http.StatusNotFound, w.Code)
	})

	t.Run("DBエラー", func(t *testing.T) {
		h := handler.NewTaskHandler(&errTaskUsecase{err: errors.New("db error")})

		ctx, w := newTaskContext("DELETE", "/tasks/1", "")
		ctx.Params = gin.Params{{Key: "id", Value: "1"}}
		h.DeleteTask(ctx)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})
}
