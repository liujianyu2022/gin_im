// @title User API
// @version 1.0
// @description This is user API.
// @termsOfService http://swagger.io/terms/
// @host localhost:8080
// @BasePath /api

package handler

import (
	"fmt"
	"gin_im/api"
	"gin_im/config"
	"gin_im/service"

	"github.com/asaskevich/govalidator"
	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	Service *service.UserService
	Config  *config.Config
}

func NewUserHandler(service *service.UserService, config *config.Config) *UserHandler {
	return &UserHandler{
		Service: service,
		Config:  config,
	}
}

// @Summary Register
// @Description Register
// @Tags user
// @Accept  json
// @Produce  json
// @Param request body api.RegisterRequest true "User registration information"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Router /user/register [post]
func (handler *UserHandler) Register(ctx *gin.Context) {
	var request api.RegisterRequest

	err := ctx.ShouldBind(&request)
	if err != nil {
		api.HandleError(ctx, api.ErrBadRequest, nil)
		return
	}

	_, err = handler.Service.Register(&request)
	if err != nil {
		api.HandleError(ctx, err, nil)
		return
	}

	api.HandleSuccess(ctx, "register successfully!")
}

// @Summary Login
// @Description Login
// @Tags user
// @Accept  json
// @Produce  json
// @Param request body api.LoginRequest true "User Login information"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Router /user/login [post]
func (handler *UserHandler) Login(ctx *gin.Context) {
	var request api.LoginRequest

	err := ctx.ShouldBind(&request)
	if err != nil {
		api.HandleError(ctx, api.ErrBadRequest, nil)
		return
	}

	token, err := handler.Service.Login(&request, handler.Config)
	if err != nil {
		api.HandleError(ctx, err, nil)
		return
	}

	data := api.LoginResponse{
		Token: token,
	}

	api.HandleSuccess(ctx, data)
}

// @Summary Get User Information
// @Description Get user information by token
// @Tags user
// @Accept  json
// @Produce  json
// @Security ApiKeyAuth
// @Param Authorization header string true "Bearer {token}"
// @Success 200 {object} map[string]any
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /user/information [get]
func (handler *UserHandler) GetUserInformation(ctx *gin.Context) {
	userId, existed := ctx.Get("userId")
	if !existed {
		api.HandleError(ctx, api.ErrBadRequest, nil)
		return
	}

	// 类型断言
	userIdUint, ok := userId.(uint)
	if !ok {
		api.HandleError(ctx, api.ErrBadRequest, nil)
		return
	}

	fmt.Println("userId = ", userIdUint)

	user, err := handler.Service.GetUserInformationById(userIdUint)
	if err != nil {
		api.HandleError(ctx, err, nil)
		return
	}

	api.HandleSuccess(ctx, user)
}

// Produce 声明 API 返回的响应格式（如 JSON、XML、HTML 等）

// @Summary Update
// @Description Update
// @Tags user
// @Accept  multipart/form-data
// @Produce  json
// @Param Authorization header string true "Bearer {token}"
// @Param name formData string false "用户名"
// @Param password formData string false "密码"
// @Param email formData string false "邮箱"
// @Param phone formData string false "电话"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Router /user/update [put]
func (handler *UserHandler) UpdateUser(ctx *gin.Context) {
	userId, existed := ctx.Get("userId")
	if !existed {
		api.HandleError(ctx, api.ErrBadRequest, "")
		return
	}

	userIdUint, ok := userId.(uint)
	if !ok {
		api.HandleError(ctx, api.ErrBadRequest, "")
		return
	}

	updateUser, err := handler.Service.GetUserInformationById(userIdUint)

	if name := ctx.PostForm("name"); name != "" {
		updateUser.Name = name
	}
	if password := ctx.PostForm("password"); password != "" {
		updateUser.Password = password
	}
	if email := ctx.PostForm("email"); email != "" {
		updateUser.Email = email
	}
	if phone := ctx.PostForm("phone"); phone != "" {
		updateUser.Phone = phone
	}

	_, err = govalidator.ValidateStruct(updateUser)
	if err != nil {
		api.HandleError(ctx, api.ErrBadRequest, "")
		return
	}

	_, err = handler.Service.UpdateUser(updateUser)
	if err != nil {
		api.HandleError(ctx, err, "")
		return
	}

	api.HandleSuccess(ctx, updateUser)
}
