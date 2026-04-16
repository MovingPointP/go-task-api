package jwt

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// ペイロード
type Claims struct {
	UserID uint `json:"user_id"`
	jwt.RegisteredClaims
}

func GenerateToken(userID uint) (string, error) {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		return "", fmt.Errorf("JWT_SECRET is not set")
	}
	expirationHours, err := strconv.Atoi(os.Getenv("JWT_EXPIRATION_HOURS"))
	if err != nil || expirationHours <= 0 {
		expirationHours = 72
	}

	expiresAt := time.Now().Add(time.Duration(expirationHours) * time.Hour)

	claims := Claims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiresAt),  // 有効期限
			IssuedAt:  jwt.NewNumericDate(time.Now()), // 発行日時
		},
	}

	// 署名方式の決定
	// HS256: 256バイトの共通鍵
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	// エンコード
	return token.SignedString([]byte(secret))
}

func ParseToken(tokenString string) (*Claims, error) {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		return nil, fmt.Errorf("JWT_SECRET is not set")
	}

	// 署名アルゴリズムの検証とデコード
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(secret), nil
	})
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	return claims, nil
}
