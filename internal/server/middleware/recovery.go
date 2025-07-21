package middleware

import (
	"net/http"

	"github.com/bestruirui/bestsub/internal/server/resp"
	"github.com/bestruirui/bestsub/internal/utils/log"
	"github.com/gin-gonic/gin"
)

// Middleware 恢复中间件
func Recovery() gin.HandlerFunc {
	return gin.CustomRecovery(func(c *gin.Context, recovered interface{}) {
		log.Errorf("Panic recovered: %v", recovered)
		resp.Error(c, http.StatusInternalServerError, "An unexpected error occurred")
		c.Abort()
	})
}
