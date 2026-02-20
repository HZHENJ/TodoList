package middleware

import (
	"net/http"
	"strings"
	"to-do-list/pkg/e"
	"to-do-list/pkg/utils"

	"github.com/gin-gonic/gin"
)

// JWT 鉴权中间件
func JWT() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code": e.ERROR_AUTH_CHECK_TOKEN_FAIL,
				"msg":  e.GetMsg(e.ERROR_AUTH_CHECK_TOKEN_FAIL),
				"data": nil,
			})
			c.Abort()
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code": e.ERROR_AUTH_CHECK_TOKEN_FAIL,
				"msg":  e.GetMsg(e.ERROR_AUTH_CHECK_TOKEN_FAIL),
				"data": nil,
			})
			c.Abort()
			return
		}
		tokenString := parts[1]

		// TODO Redis 黑名单

		// 解析Token
		claims, err := utils.ParseToken(tokenString)
		if err != nil {
			// 如果走到这里，通常是因为篡改，或者是 jwt 包底层校验出它自然过期了
			c.JSON(http.StatusUnauthorized, gin.H{
				"code": e.ERROR_AUTH_CHECK_TOKEN_TIMEOUT,
				"msg":  e.GetMsg(e.ERROR_AUTH_CHECK_TOKEN_FAIL),
				"data": nil,
			})
			c.Abort()
			return
		}

		// 这里可以其他的内容
		c.Set("Username", claims.Username)

		c.Next()
	}
}
