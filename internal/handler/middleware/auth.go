package middleware

import (
	"net/http"
	"strings"

	"github.com/MovingPointP/go-task-api/pkg/jwt"
	"github.com/gin-gonic/gin"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// ヘッダー取得
		authHeader := ctx.GetHeader("Authorization")
		if authHeader == "" {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized,
				gin.H{"error": "Authorization header is required"},
			)
			return
		}

		// "Bearer <token>"の形式かチェック
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized,
				gin.H{"error": "Authorization header format must be Bearer {token}"},
			)
			return
		}

		// トークンの検証
		claims, err := jwt.ParseToken(parts[1])
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized,
				gin.H{"error": "Invalid or expired token"},
			)
			return
		}

		// コンテキストにセット
		ctx.Set("UserID", claims.UserID)
		ctx.Next()
	}
}
