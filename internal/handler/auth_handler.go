package handler

import (
	"errors"
	"net/http"

	"github.com/MovingPointP/go-task-api/internal/domain/entity"
	"github.com/MovingPointP/go-task-api/internal/usecase"
	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	authUsecase usecase.AuthUsecase
}

func NewAuthHandler(authUsecase usecase.AuthUsecase) *AuthHandler {
	return &AuthHandler{authUsecase: authUsecase}
}

// ユーザー登録のリクエストボディ
type RegisterRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
}

// ログインのリクエストボディ
type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// 認証レスポンス
type AuthResponse struct {
	Token string `json:"token"`
	Email string `json:"email"`
}

// @Summary     ユーザー登録
// @Description メールアドレスとパスワードでユーザーを新規登録し、JWTを返す
// @Tags        auth
// @Accept      json
// @Produce     json
// @Param       body body RegisterRequest true "登録情報"
// @Success     201 {object} AuthResponse
// @Failure     400 {object} map[string]string
// @Failure     409 {object} map[string]string
// @Failure     500 {object} map[string]string
// @Router      /auth/register [post]
func (h *AuthHandler) Register(ctx *gin.Context) {
	var req RegisterRequest

	// リクエストボディをバインドしてバリデーション
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, token, err := h.authUsecase.Register(req.Email, req.Password)
	if err != nil {
		if errors.Is(err, entity.ErrEmailAlreadyInUse) {
			ctx.JSON(http.StatusConflict, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to register user"})
		return
	}

	ctx.JSON(http.StatusCreated, AuthResponse{
		Token: token,
		Email: user.Email,
	})
}

// @Summary     ログイン
// @Description メールアドレスとパスワードでログインし、JWTを返す
// @Tags        auth
// @Accept      json
// @Produce     json
// @Param       body body LoginRequest true "ログイン情報"
// @Success     200 {object} AuthResponse
// @Failure     400 {object} map[string]string
// @Failure     401 {object} map[string]string
// @Failure     500 {object} map[string]string
// @Router      /auth/login [post]
func (h *AuthHandler) Login(ctx *gin.Context) {
	var req LoginRequest

	// リクエストボディをバインドしてバリデーション
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, token, err := h.authUsecase.Login(req.Email, req.Password)
	if err != nil {
		if errors.Is(err, entity.ErrInvalidCredentials) {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to login"})
		return
	}

	ctx.JSON(http.StatusOK, AuthResponse{
		Token: token,
		Email: user.Email,
	})
}
