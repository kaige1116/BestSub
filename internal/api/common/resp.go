package common

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// APIResponse 统一API响应结构
type ResponseStruct struct {
	Code    int         `json:"code" example:"200"`                  // 状态码
	Message string      `json:"message" example:"success"`           // 响应消息
	Data    interface{} `json:"data,omitempty"`                      // 响应数据
	Error   string      `json:"error,omitempty" example:"error_msg"` // 错误信息
}

// SuccessResponse 成功响应
type ResponseSuccessStruct struct {
	Code    int         `json:"code" example:"200"`        // 状态码
	Message string      `json:"message" example:"success"` // 响应消息
	Data    interface{} `json:"data,omitempty"`            // 响应数据
}

// ErrorResponse 错误响应
type ResponseErrorStruct struct {
	Code    int    `json:"code" example:"400"`                     // 状态码
	Message string `json:"message" example:"Bad Request"`          // 响应消息
	Error   string `json:"error" example:"Invalid request format"` // 错误详情
}

// ValidationErrorResponse 验证错误响应
type ResponseValidationStruct struct {
	Code    int                    `json:"code" example:"422"`                  // 状态码
	Message string                 `json:"message" example:"Validation failed"` // 响应消息
	Error   string                 `json:"error" example:"Validation error"`    // 错误信息
	Details map[string]interface{} `json:"details,omitempty"`                   // 验证错误详情
}

// PaginationResponse 分页响应结构
type ResponsePaginationStruct struct {
	Page     int         `json:"page" example:"1"`       // 当前页码
	PageSize int         `json:"page_size" example:"10"` // 每页大小
	Total    int64       `json:"total" example:"100"`    // 总记录数
	Data     interface{} `json:"data"`                   // 数据列表
}

func ResponseSuccess(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, ResponseSuccessStruct{
		Code:    http.StatusOK,
		Message: "success",
		Data:    data,
	})
}

func ResponseError(c *gin.Context, code int, err error) {
	c.JSON(code, ResponseErrorStruct{
		Code:    code,
		Message: "error",
		Error:   err.Error(),
	})
	c.Abort()
}
