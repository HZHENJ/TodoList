package service

import (
	"context"
	"time"
	"to-do-list/internal/repository/db/dao"
	"to-do-list/internal/repository/db/model"
	"to-do-list/pkg/e"
	"to-do-list/pkg/utils"
	"to-do-list/types"
)

type UserService struct {
	userDao dao.UserDao
}

func NewUserService(userDao dao.UserDao) *UserService {
	return &UserService{
		userDao: userDao,
	}
}

// Register 用户注册
func (s *UserService) Register(ctx context.Context, req *types.UserRegisterRequest) (interface{}, int) {
	// 1.检查邮箱是否被注册过
	exist, err := s.userDao.CheckEmailExist(ctx, req.Email)
	if err != nil {
		return nil, e.ERROR
	}
	if exist {
		return nil, e.ERROR_USER_EXIST
	}

	// 2.自动生成随机的Username和Nickname
	randomUsername := utils.GenerateDefaultUsername()

	// 3.创建新的用户
	newUser := &model.User{
		Email:    req.Email,
		Username: randomUsername,
		// TODO Nickname暂时和Username相同
		Nickname: randomUsername,
	}

	// 4.密码加密
	if err := newUser.SetPassword(req.Password); err != nil {
		return nil, e.ERROR
	}

	// 存入数据库
	if err := s.userDao.CreateUser(ctx, newUser); err != nil {
		return nil, e.ERROR
	}

	return nil, e.SUCCESS
}

// Login 用户登陆
func (s *UserService) Login(ctx context.Context, req *types.UserLoginRequest) (interface{}, int) {
	// 1. 检查用户是否存在
	user, err := s.userDao.FindByEmail(ctx, req.Email)
	if err != nil {
		return nil, e.ERROR_USER_NOT_EXIST
	}

	// 2. 校验密码是否正确
	if !user.CheckPassword(req.Password) {
		return nil, e.ERROR_USER_WRONG_PWD
	}

	// 3. 生成JWT token
	// 确定JWT中存放的内容
	// TODO authority需要设置一下，现在1代表普通用户
	token, err := utils.GenerateToken(user.ID, user.Email, user.Username, 1)
	if err != nil {
		return nil, e.ERROR
	}

	resp := &types.UserLoginResponse{
		Token: token,
		User: types.UserInfo{
			Id:       user.ID,
			Nickname: user.Nickname,
			Email:    user.Email,
			Username: user.Username,
		},
	}

	// 4. 返回结果
	return resp, e.SUCCESS
}

// Logout 用户退出登陆
func (s *UserService) Logout(ctx context.Context, token string) (interface{}, int) {
	// 1. 解析token
	claims, err := utils.ParseToken(token)
	if err != nil {
		// 如果解析失败了，可能Token本身就过期或者被篡改了，说明Token失效了
		// 既然本身就是无效的那么就就认定用户已经是属于“未登录/登出”状态了
		return nil, e.SUCCESS
	}

	// 判断token是否过期
	now := time.Now().Unix()
	exp := claims.ExpiresAt.Unix()
	remain := exp - now

	if remain <= 0 {
		return nil, e.SUCCESS
	}

	// TODO 将Token放进redis中

	return nil, e.SUCCESS
}
