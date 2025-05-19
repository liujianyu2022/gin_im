package api

import (
	"errors"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Response struct {
	Code    int    `json: "code"`
	Message string `json: "message"`
	Data    any    `json: "data"`
}

// 业务状态码范围（兼容HTTP语义）
const (
	CodeSuccess       = 200 // 成功
	CodeBadRequest    = 400 // 参数错误
	CodeUnauthorized  = 401 // 未授权
	CodeForbidden     = 403 // 无权限
	CodeNotFound      = 404 // 资源不存在
	CodeConflict      = 409 // 新增：用于资源冲突
	CodeWrongPassword = 410
	CodeServerError   = 500 // 服务端错误
	CodeDatabase      = 501 // 新增：数据库错误
	// ...可扩展其他业务码
)

// 预定义错误
var (
	ErrBadRequest     = newError(CodeBadRequest, "Invalid Request Parameters")
	ErrUnauthorized   = newError(CodeUnauthorized, "未授权访问")
	ErrForbidden      = newError(CodeForbidden, "禁止访问")
	ErrNotFound       = newError(CodeNotFound, "资源不存在")
	ErrAlreadyExists  = newError(CodeConflict, "resource already exists")
	ErrWrongPassword  = newError(CodeWrongPassword, "Invalid Password")
	ErrInternalServer = newError(CodeServerError, "服务器内部错误")
	// ...其他业务错误
)

// errorCodeMap 存储错误到状态码的映射
var errorCodeMap = map[error]int{}

func newError(code int, msg string) error {
	err := errors.New(msg) // 新建一个 error对象
	errorCodeMap[err] = code
	return err
}

// ------------------------------------------------------------------
// 统一响应处理（始终返回HTTP 200）
// ------------------------------------------------------------------

func HandleResponse(ctx *gin.Context, code int, message string, data any) {
	if data == nil {
		data = struct{}{}
	}

	// 记录服务端错误（实际项目应接入日志系统）
	if code >= CodeServerError {
		log.Printf("[ERROR] code = %d, msg = %s", code, message)
	}

	response := Response{
		Code:    code,
		Message: message,
		Data:    data,
	}

	ctx.JSON(http.StatusOK, response)
}

// 快捷方法
func HandleSuccess(ctx *gin.Context, data any) {
	HandleResponse(ctx, CodeSuccess, "success", data)
}

func HandleError(ctx *gin.Context, err error, data any) {
	code := CodeServerError // 默认服务端错误

	customCode, exists := errorCodeMap[err]
	if exists {
		code = customCode
	}

	// 生产环境隐藏敏感错误详情
	var message string = err.Error()
	if code >= CodeServerError && gin.Mode() != gin.DebugMode {
		message = "internal server error"
	}

	HandleResponse(ctx, code, message, data)
}