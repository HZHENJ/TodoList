package dao

import (
	"context"
	"to-do-list/internal/repository/db/model"

	"gorm.io/gorm"
)

// UserDao represent interface
type UserDao interface {
	CheckEmailExist(ctx context.Context, email string) (bool, error)
	CreateUser(ctx context.Context, user *model.User) error
	FindByEmail(ctx context.Context, email string) (*model.User, error)
}

// userDao means external packet cannot use it directly
// must via following function
type userDao struct {
	db *gorm.DB
}

// NewUserDao
func NewUserDao(db *gorm.DB) UserDao {
	return &userDao{
		db: db,
	}
}

// CheckEmailExist 检查邮箱是否已经存在
func (dao *userDao) CheckEmailExist(ctx context.Context, email string) (bool, error) {
	var count int64

	// 使用Count相对First性能更好，只需要知道是否存在
	err := dao.db.WithContext(ctx).
		Model(&model.User{}).
		Where("email = ?", email).
		Count(&count).Error

	if err != nil {
		return false, err
	}

	return count > 0, nil
}

// CreateUser 插入用户
func (dao *userDao) CreateUser(ctx context.Context, user *model.User) error {
	return dao.db.WithContext(ctx).Create(user).Error
}

// FindByEmail 通过Email查找用户
func (dao *userDao) FindByEmail(ctx context.Context, email string) (*model.User, error) {
	var user model.User
	err := dao.db.WithContext(ctx).
		Where("email = ?", email).
		First(&user).Error
	return &user, err
}
