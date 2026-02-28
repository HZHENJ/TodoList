package v1

import (
	"net/http"
	"to-do-list/internal/repository/db"
	"to-do-list/internal/repository/db/dao"
	"to-do-list/internal/service"
	"to-do-list/pkg/ctl"
	"to-do-list/pkg/e"
	"to-do-list/types"

	"github.com/gin-gonic/gin"
)

// CreateTask 创建用户接口
func CreateTask(c *gin.Context) {
	app := ctl.NewWrapper(c)
	var req types.CreateTaskRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		app.Error(e.INVALID_PARAMS, err)
		return
	}

	uid, exists := c.Get("UserId")
	if !exists {
		app.Response(http.StatusUnauthorized, e.ERROR_AUTH_CHECK_TOKEN_FAIL, nil)
		return
	}
	userId := uid.(uint)

	taskDao := dao.NewTaskDao(db.DB)
	taskService := service.NewTaskService(taskDao)

	data, code := taskService.CreateTask(c.Request.Context(), userId, &req)
	if code != e.SUCCESS {
		app.Response(http.StatusOK, code, nil)
		return
	}

	app.Success(data)
}

// ListTasks
func ListTasks(c *gin.Context) {
	app := ctl.NewWrapper(c)
	var req types.ListTasksRequest

	if err := c.ShouldBindQuery(&req); err != nil {
		app.Error(e.INVALID_PARAMS, err)
		return
	}

	// 从 JWT 中间件的 Context 中安全获取当前操作人的 ID
	uId, exists := c.Get("userID")
	if !exists {
		app.Response(http.StatusUnauthorized, e.ERROR_AUTH_CHECK_TOKEN_FAIL, nil)
		return
	}
	userId := uId.(uint)

	taskDao := dao.NewTaskDao(db.DB)
	taskService := service.NewTaskService(taskDao)

	// 调用业务逻辑
	data, code := taskService.ListTasks(c.Request.Context(), userId, &req)
	if code != e.SUCCESS {
		app.Response(http.StatusOK, code, nil)
		return
	}

	app.Success(data)
}
