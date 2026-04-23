package handler_test

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/MovingPointP/go-task-api/internal/domain/entity"
	"github.com/MovingPointP/go-task-api/internal/handler"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func newAuthContext(method, url, body string) (*gin.Context, *httptest.ResponseRecorder) {
	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)
	ctx.Request, _ = http.NewRequest(method, url, strings.NewReader(body))
	ctx.Request.Header.Set("Content-Type", "application/json")
	return ctx, w
}

func TestAuthHandler_Register(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("正常系", func(t *testing.T) {
		email := "test@example.com"
		password := "password123"

		h := handler.NewAuthHandler(
			NewMockAuthUsecase(
				&entity.User{Email: email},
				"mock-token",
				nil,
			))

		body := `{"email":"` + email + `", "password":"` + password + `"}`
		ctx, w := newAuthContext("POST", "/register", body)

		h.Register(ctx)

		assert.Equal(t, http.StatusCreated, w.Code)
		assert.JSONEq(t, `{"token": "mock-token", "email": "test@example.com" }`, w.Body.String())
	})

	t.Run("メールアドレス重複", func(t *testing.T) {
		h := handler.NewAuthHandler(
			NewMockAuthUsecase(nil, "", entity.ErrEmailAlreadyInUse),
		)

		ctx, w := newAuthContext("POST", "/register", `{"email":"duplicate@example.com", "password":"password123"}`)

		h.Register(ctx)

		assert.Equal(t, http.StatusConflict, w.Code)
		assert.JSONEq(t, `{"error": "email already in use"}`, w.Body.String())
	})

	t.Run("その他サーバーエラー", func(t *testing.T) {
		h := handler.NewAuthHandler(
			NewMockAuthUsecase(nil, "", errors.New("unexpected error")),
		)

		ctx, w := newAuthContext("POST", "/register", `{"email":"test@example.com", "password":"password123"}`)

		h.Register(ctx)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.JSONEq(t, `{"error": "failed to register user"}`, w.Body.String())
	})
}

func TestAuthHandler_Login(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("正常系", func(t *testing.T) {
		email := "test@example.com"
		password := "password123"

		h := handler.NewAuthHandler(
			NewMockAuthUsecase(
				&entity.User{Email: email},
				"mock-token",
				nil,
			))

		body := `{"email":"` + email + `", "password":"` + password + `"}`
		ctx, w := newAuthContext("POST", "/login", body)

		h.Login(ctx)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.JSONEq(t, `{"token": "mock-token", "email": "test@example.com"}`, w.Body.String())
	})

	t.Run("認証情報不正", func(t *testing.T) {
		h := handler.NewAuthHandler(
			NewMockAuthUsecase(nil, "", entity.ErrInvalidCredentials),
		)

		ctx, w := newAuthContext("POST", "/login", `{"email":"test@example.com", "password":"wrongpassword"}`)

		h.Login(ctx)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
		assert.JSONEq(t, `{"error": "invalid email or password"}`, w.Body.String())
	})

	t.Run("その他サーバーエラー", func(t *testing.T) {
		h := handler.NewAuthHandler(
			NewMockAuthUsecase(nil, "", errors.New("unexpected error")),
		)

		ctx, w := newAuthContext("POST", "/login", `{"email":"test@example.com", "password":"password123"}`)

		h.Login(ctx)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.JSONEq(t, `{"error": "failed to login"}`, w.Body.String())
	})
}
