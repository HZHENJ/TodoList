package v1

import (
	"net/http"
	"strings"
	"to-do-list/internal/repository/db"
	"to-do-list/internal/repository/db/dao"
	"to-do-list/internal/service"
	"to-do-list/pkg/ctl"
	"to-do-list/pkg/e"
	"to-do-list/types"

	"github.com/gin-gonic/gin"
)

// UserRegister 用户注册接口
func UserRegister(c *gin.Context) {
	app := ctl.NewWrapper(c)
	var req types.UserRegisterRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		app.Error(e.INVALID_PARAMS, err)
		return
	}

	userDao := dao.NewUserDao(db.DB)
	userService := service.NewUserService(userDao)

	data, code := userService.Register(c.Request.Context(), &req)
	if code != e.SUCCESS {
		app.Response(http.StatusOK, code, nil)
		return
	}

	app.Success(data)
}

// UserLogin 用户登陆接口
func UserLogin(c *gin.Context) {
	app := ctl.NewWrapper(c)
	var req types.UserLoginRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		app.Error(e.INVALID_PARAMS, err)
		return
	}

	userDao := dao.NewUserDao(db.DB)
	userService := service.NewUserService(userDao)

	data, code := userService.Login(c.Request.Context(), &req)
	if code != e.SUCCESS {
		app.Response(http.StatusOK, code, nil)
		return
	}
	app.Success(data)
}

// UserLogout 用户登出接口
func UserLogout(c *gin.Context) {
	app := ctl.NewWrapper(c)
	authHeader := c.GetHeader("Authorization")

	if authHeader == "" {
		app.Response(http.StatusOK, e.ERROR_AUTH_CHECK_TOKEN_FAIL, nil)
		return
	}

	parts := strings.SplitN(authHeader, " ", 2)
	if len(parts) != 2 || parts[0] != "Bearer" {
		app.Response(http.StatusOK, e.ERROR_AUTH_CHECK_TOKEN_FAIL, nil)
		return
	}
	token := parts[1]

	userDao := dao.NewUserDao(db.DB)
	userService := service.NewUserService(userDao)

	data, code := userService.Logout(c.Request.Context(), token)
	if code != e.SUCCESS {
		app.Response(http.StatusOK, code, nil)
		return
	}
	app.Success(data)
}
