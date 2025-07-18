package middleware

import (
	"fmt"
	"net/http"

	"github.com/bestruirui/bestsub/internal/api/common"
	"github.com/bestruirui/bestsub/internal/utils/log"
	"github.com/gin-gonic/gin"
)

// Middleware 恢复中间件
func Recovery() gin.HandlerFunc {
	return gin.CustomRecovery(func(c *gin.Context, recovered interface{}) {
		log.Errorf("Panic recovered: %v", recovered)

		common.ResponseError(c, http.StatusInternalServerError, fmt.Errorf("An unexpected error occurred"))
		c.Abort()
	})
}
