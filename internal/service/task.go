package service

import (
	"context"
	"to-do-list/internal/repository/db/dao"
	"to-do-list/internal/repository/db/model"
	"to-do-list/pkg/e"
	"to-do-list/types"
)

type TaskService struct {
	taskDao dao.TaskDao
}

func NewTaskService(taskDao dao.TaskDao) *TaskService {
	return &TaskService{
		taskDao: taskDao,
	}
}

// CreateTask 创建任务
func (s *TaskService) CreateTask(cxt context.Context, userId uint, req *types.CreateTaskRequest) (interface{}, int) {
	newTask := &model.Task{
		UserId:   userId,
		Title:    req.Title,
		Content:  req.Content,
		Status:   req.Status,
		Category: req.Category,
	}

	err := s.taskDao.CreateTask(cxt, newTask)
	if err != nil {
		return nil, e.ERROR
	}
	return newTask, e.SUCCESS
}

// ListTasks 列出所有任务列表
func (s *TaskService) ListTasks(ctx context.Context, userId uint, req *types.ListTasksRequest) (interface{}, int) {
	tasks, total, err := s.taskDao.ListTasks(ctx, userId, req.Page, req.PageSize)
	if err != nil {
		return nil, e.ERROR
	}

	resp := &types.ListTasksResponse{
		Items: tasks,
		Total: total,
	}

	return resp, e.SUCCESS
}

// UpdateTask 更新任务
func (s *TaskService) UpdateTask(ctx context.Context, taskId, userId uint, req *types.UpdateTaskRequest) (interface{}, int) {
	updateData := make(map[string]interface{})
	if req.Title != "" {
		updateData["title"] = req.Title
	}
	if req.Content != "" {
		updateData["content"] = req.Content
	}
	if req.Category != "" {
		updateData["category"] = req.Category
	}
	if req.Status != nil {
		updateData["status"] = *req.Status
	}

	if len(updateData) == 0 {
		return nil, e.SUCCESS
	}

	err := s.taskDao.UpdateTask(ctx, taskId, userId, updateData)
	if err != nil {
		return nil, e.ERROR
	}
	return nil, e.SUCCESS
}

// DeleteTask
func (s *TaskService) DeleteTask(ctx context.Context, taskId, userId uint) (interface{}, int) {
	if err := s.taskDao.DeleteTask(ctx, taskId, userId); err != nil {
		return nil, e.ERROR
	}
	return nil, e.SUCCESS
}
