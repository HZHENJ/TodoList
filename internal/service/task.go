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

// ListTasks
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

// DeleteTask

// UpdateTask
