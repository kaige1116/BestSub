package resp

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type ResponseStruct struct {
	Code    int         `json:"code" example:"200"`
	Message string      `json:"message" example:"success"`
	Data    interface{} `json:"data,omitempty"`
}

type ResponsePaginationStruct struct {
	Page     int         `json:"page" example:"1"`
	PageSize int         `json:"page_size" example:"10"`
	Total    uint16      `json:"total" example:"100"`
	Data     interface{} `json:"data"`
}

func Success(c *gin.Context, data any) {
	c.JSON(http.StatusOK, ResponseStruct{
		Code:    http.StatusOK,
		Message: "success",
		Data:    data,
	})
}

func Error(c *gin.Context, code int, err string) {
	c.JSON(code, ResponseStruct{
		Code:    code,
		Message: err,
	})
	c.Abort()
}
func ErrorBadRequest(c *gin.Context) {
	c.JSON(http.StatusBadRequest, ResponseStruct{
		Code:    http.StatusBadRequest,
		Message: "bad request",
	})
	c.Abort()
}
