package types

import "to-do-list/internal/repository/db/model"

type CreateTaskRequest struct {
	Title    string `json:"title" binding:"required"`
	Content  string `json:"content"`
	Status   int    `json:"status"`
	Category string `json:"category"`
}

type ListTasksRequest struct {
	Page     int `form:"page" binding:"required,min=1"`
	PageSize int `form:"pageSize" binding:"required,min=1,max=1"`

	// 可选过滤条件
	Status *int `form:"status"` // 用指针是为了区分前端传了0，还是没传这个参数
}

type ListTasksResponse struct {
	Items []*model.Task `json:"items"`
	Total int64         `json:"total"`
}
