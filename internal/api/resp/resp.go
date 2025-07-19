package resp

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// SuccessResponse 成功响应
type SuccessStruct struct {
	Code    int         `json:"code" example:"200"`        // 状态码
	Message string      `json:"message" example:"success"` // 响应消息
	Data    interface{} `json:"data,omitempty"`            // 响应数据
}

// ErrorResponse 错误响应
type ErrorStruct struct {
	Code    int    `json:"code" example:"400"`                     // 状态码
	Message string `json:"message" example:"error"`                // 响应消息
	Error   string `json:"error" example:"Invalid request format"` // 错误详情
}

// PaginationResponse 分页响应结构
type ResponsePaginationStruct struct {
	Page     int         `json:"page" example:"1"`       // 当前页码
	PageSize int         `json:"page_size" example:"10"` // 每页大小
	Total    uint16      `json:"total" example:"100"`    // 总记录数
	Data     interface{} `json:"data"`                   // 数据列表
}

func Success(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, SuccessStruct{
		Code:    http.StatusOK,
		Message: "success",
		Data:    data,
	})
}

func Error(c *gin.Context, code int, err string) {
	c.JSON(code, ErrorStruct{
		Code:    code,
		Message: "error",
		Error:   err,
	})
	c.Abort()
}
func ErrorBadRequest(c *gin.Context) {
	c.JSON(http.StatusBadRequest, ErrorStruct{
		Code:    http.StatusBadRequest,
		Message: "error",
		Error:   "bad request",
	})
	c.Abort()
}
