package dao

import (
	"context"
	"to-do-list/internal/repository/db/model"

	"gorm.io/gorm"
)

type TaskDao interface {
	CreateTask(ctx context.Context, task *model.Task) error
	ListTasks(ctx context.Context, userId uint, page, pageSize int) ([]*model.Task, int64, error)
	UpdateTask(ctx context.Context, taskId, userId uint, updateData map[string]interface{}) error
	DeleteTask(ctx context.Context, taskId, userId uint) error
}

type taskDao struct {
	db *gorm.DB
}

func NewTaskDao(db *gorm.DB) TaskDao {
	return &taskDao{db: db}
}

// CreateTask 创建用户
func (dao *taskDao) CreateTask(ctx context.Context, task *model.Task) error {
	return dao.db.WithContext(ctx).Create(task).Error
}

// ListTasks 列出所有任务 分页
func (dao *taskDao) ListTasks(ctx context.Context, userId uint, page, pageSize int) ([]*model.Task, int64, error) {
	var tasks []*model.Task
	var total int64

	query := dao.db.WithContext(ctx).Model(&model.Task{}).Where("user_id = ?", userId)

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	err := query.
		Preload("User").
		Offset(offset).
		Limit(pageSize).
		Order("created_at desc").
		Find(&tasks).Error
	return tasks, total, err
}

// UpdateTask 更新任务
func (dao *taskDao) UpdateTask(ctx context.Context, taskId, userId uint, updateData map[string]interface{}) error {
	if len(updateData) == 0 {
		return nil
	} // 没有做任何修改

	result := dao.db.WithContext(ctx).
		Model(&model.Task{}).
		Where("id = ? AND user_id = ?", taskId, userId).
		Updates(updateData)

	return result.Error
}

// DeleteTask 删除任务
func (dao *taskDao) DeleteTask(ctx context.Context, taskId, userId uint) error {
	result := dao.db.WithContext(ctx).
		Where("id = ? AND user_id = ?", taskId, userId).
		Delete(&model.Task{})
	return result.Error
}
