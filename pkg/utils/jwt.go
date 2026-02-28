package utils

import (
	"errors"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// TODO 这里先用全局
var jwtSecret = []byte(os.Getenv("JWT_SECRET"))

// Claims 自定义载荷
type Claims struct {
	UserId    uint   `json:"user_id"`
	Email     string `json:"email"`
	Username  string `json:"username"`
	Authority int    `json:"authority"`

	jwt.RegisteredClaims
}

// GenerateToken 签发用户Token
func GenerateToken(userId uint, email, username string, authority int) (string, error) {
	now := time.Now()
	exp := now.Add(time.Hour * 1)
	claims := Claims{
		UserId:    userId,
		Email:     email,
		Username:  username,
		Authority: authority,

		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(exp),
			IssuedAt:  jwt.NewNumericDate(now),
			Issuer:    "to-do-list",
		},
	}
	// 指定签名方法HS256
	tokenClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token, err := tokenClaims.SignedString(jwtSecret)
	return token, err
}

// ParseToken 验证用户token
func ParseToken(tokenString string) (*Claims, error) {
	// 提供解析Token所需的密钥
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		// 校验签名算法是否预期（防范被篡改为 NONE 算法的漏洞）
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return jwtSecret, nil
	})
	if err != nil {
		return nil, err
	}
	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}
	return nil, errors.New("invalid token")
}
