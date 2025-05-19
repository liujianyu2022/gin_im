// @title User API
// @version 1.0
// @description This is user API.
// @termsOfService http://swagger.io/terms/
// @host localhost:8080
// @BasePath /api

package handler

import (
	"gin_im/api"
	"gin_im/config"
	"gin_im/service"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	Service *service.UserService
	Config  *config.Config
}

func NewUserHandler(service *service.UserService, config *config.Config) *UserHandler {
	return &UserHandler{
		Service: service,
		Config: config,
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
// @Description Get user information by name
// @Tags user
// @Accept  json
// @Produce  json
// @Param name query string true "User name"
// @Success 200 {object} map[string]any
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /user/information [get]
func (handler *UserHandler) GetUserInformation(ctx *gin.Context) {
	name := ctx.Query("name")
	if name == "" {
		api.HandleError(ctx, api.ErrBadRequest, nil)
		return
	}

	user, err := handler.Service.GetUserINformationByName(name)
	if err != nil {
		api.HandleError(ctx, err, nil)
		return
	}

	api.HandleSuccess(ctx, user)
}

func (handler *UserHandler) UpdateUser(ctx *gin.Context) {

}
