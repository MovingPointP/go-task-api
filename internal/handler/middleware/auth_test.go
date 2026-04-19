package middleware_test

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/MovingPointP/go-task-api/internal/handler/middleware"
	"github.com/MovingPointP/go-task-api/pkg/jwt"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAuthMiddleware(t *testing.T) {
	os.Setenv("JWT_SECRET", "test-secret")
	gin.SetMode(gin.TestMode)

	t.Run("正常系: 200が返る", func(t *testing.T) {
		g := gin.New()
		g.Use(middleware.AuthMiddleware())
		g.GET("/test", func(ctx *gin.Context) {
			uid, _ := ctx.Get("UserID")
			ctx.JSON(http.StatusOK, gin.H{"user_id": uid})
		})

		validToken, err := jwt.GenerateToken(123)
		require.NoError(t, err, "GenerateToken failed")

		req := httptest.NewRequest(http.MethodGet, "/test", nil)
		req.Header.Set("Authorization", "Bearer "+validToken)
		w := httptest.NewRecorder()

		g.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.JSONEq(t, `{"user_id": 123}`, w.Body.String())
	})

	t.Run("異常系: ヘッダーがない場合", func(t *testing.T) {
		g := gin.New()
		g.Use(middleware.AuthMiddleware())

		req := httptest.NewRequest(http.MethodGet, "/test", nil)
		w := httptest.NewRecorder()

		g.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
		assert.Contains(t, w.Body.String(), "Authorization header is required")
	})

	t.Run("異常系: トークンの形式が違う場合", func(t *testing.T) {
		g := gin.New()
		g.Use(middleware.AuthMiddleware())

		req := httptest.NewRequest(http.MethodGet, "/test", nil)
		req.Header.Set("Authorization", "Wrong format")
		w := httptest.NewRecorder()

		g.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
		assert.Contains(t, w.Body.String(), "Authorization header format must be Bearer {token}")
	})

	t.Run("異常系: 誤ったトークンの場合", func(t *testing.T) {
		g := gin.New()
		g.Use(middleware.AuthMiddleware())

		wrongToken := "wrong-token"

		req := httptest.NewRequest(http.MethodGet, "/test", nil)
		req.Header.Set("Authorization", "Bearer "+wrongToken)
		w := httptest.NewRecorder()

		g.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
		assert.Contains(t, w.Body.String(), "Invalid or expired token")
	})
}
