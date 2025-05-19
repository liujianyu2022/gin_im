package service

import (
	"fmt"
	"gin_im/api"
	"gin_im/config"
	"gin_im/model"
	"gin_im/repository"
	"gin_im/tools"
	"time"
)

type UserService struct {
	Repository *repository.UserRepository
}

func NewUserService(repository *repository.UserRepository) *UserService {
	return &UserService{
		Repository: repository,
	}
}

func (service *UserService) Register(request *api.RegisterRequest) (*model.User, error) {
	// 检查用户名是否已存在
	if user, _ := service.Repository.GetUserByName(request.Name); user != nil {
		return nil, api.ErrAlreadyExists
	}

	// 检查邮箱是否已存在
	if user, _ := service.Repository.GetUserByEmail(request.Email); user != nil {
		return nil, api.ErrAlreadyExists
	}

	// 检查电话是否已存在
	if user, _ := service.Repository.GetUserByPhone(request.Phone); user != nil {
		return nil, api.ErrAlreadyExists
	}

	hashedPassword, err := tools.HashPassword(request.Password)
	if err != nil {
		return nil, api.ErrBadRequest
	}

	user := &model.User{
		Name:     request.Name,
		Password: hashedPassword,
		Email:    request.Email,
		Phone:    request.Phone,
	}

	return service.Repository.CreateUser(user)
}

func (service *UserService) Login(request *api.LoginRequest, config *config.Config) (string, error) {
	user, err := service.Repository.GetUserByName(request.Name)
	if err != nil {
		return "", api.ErrInternalServer
	}
	if user == nil {
		return "", api.ErrNotFound
	}

	fmt.Println("user = ", user)

	// 验证密码
	if !tools.CheckPasswordHash(user.Password, request.Password) {
		return "", api.ErrWrongPassword
	}

	token, err := tools.GenerateToken(user.ID, config)
	if err != nil {
		return "", api.ErrInternalServer
	}

	fmt.Println("token = ", token)

	// 更新最后登录时间
	now := time.Now()
	user.LoginTime = uint64(now.Unix())

	err = service.Repository.UpdateUser(user)
	if err != nil {
		return "", api.ErrInternalServer
	}

	return token, nil
}

func (service *UserService) GetUserINformationByName(name string) (*model.User, error) {
	user, err := service.Repository.GetUserByName(name)
	return user, err
}
