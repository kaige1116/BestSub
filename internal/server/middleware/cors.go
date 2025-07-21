package middleware

import (
	"github.com/bestruirui/bestsub/internal/utils"
	"github.com/bestruirui/bestsub/internal/utils/log"
	"github.com/gin-gonic/gin"
)

// CORS中间件
func Cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		if utils.IsDebug() {
			c.Header("Access-Control-Allow-Origin", "*")
		} else {
			origin := c.Request.Header.Get("Origin")
			if origin != "" {
				log.Warnf("CORS %s: Origin=%s, Host=%s, Method=%s, Path=%s, IP=%s, UserAgent=%s",
					"BLOCKED", origin, c.Request.Host, c.Request.Method, c.Request.URL.Path,
					c.ClientIP(), c.GetHeader("User-Agent"))
			}
		}
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Authorization")
		c.Header("Access-Control-Expose-Headers", "Content-Length")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	}
}
