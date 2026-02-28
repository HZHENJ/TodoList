package routes

import (
	v1 "to-do-list/internal/api/v1"
	"to-do-list/internal/middleware"

	"github.com/gin-gonic/gin"
)

func NewRouter() *gin.Engine {
	r := gin.Default()

	// 挂载中间件
	r.Use(middleware.Cors())

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"msg":    "pong",
			"status": "ok",
		})
	})

	apiv1 := r.Group("/api/v1")
	{
		// 公开路由
		userPublic := apiv1.Group("/user")
		{
			userPublic.POST("/register", v1.UserRegister)
			userPublic.POST("/login", v1.UserLogin)

		}

		authed := apiv1.Group("")
		authed.Use(middleware.JWT())
		{
			authed.POST("/user/logout", v1.UserLogout)

			taskGroup := authed.Group("/task")
			{
				taskGroup.POST("/create", v1.CreateTask)
				taskGroup.GET("/list", v1.ListTasks)
			}
		}
	}
	return r
}
