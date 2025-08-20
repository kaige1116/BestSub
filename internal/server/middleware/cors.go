package middleware

import (
	"net/http"
	"strings"

	"github.com/bestruirui/bestsub/internal/server/resp"
	"github.com/bestruirui/bestsub/internal/utils"
	"github.com/bestruirui/bestsub/internal/utils/log"
	"github.com/gin-gonic/gin"
)

func Cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		if utils.IsDebug() {
			c.Header("Access-Control-Allow-Origin", "*")
			c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
			c.Header("Access-Control-Allow-Headers", "Authorization, Content-Type, X-Requested-With")
			c.Header("Access-Control-Expose-Headers", "Content-Length")
			if c.Request.Method == "OPTIONS" {
				c.AbortWithStatus(204)
				return
			}
		} else {
			origin := c.Request.Header.Get("Origin")
			if origin != "" {
				originWithoutScheme := origin
				if strings.HasPrefix(origin, "http://") {
					originWithoutScheme = origin[7:]
				} else if strings.HasPrefix(origin, "https://") {
					originWithoutScheme = origin[8:]
				}

				if originWithoutScheme != c.Request.Host {
					log.Warnf("CORS %s: Origin=%s, Host=%s, Method=%s, Path=%s, IP=%s, UserAgent=%s",
						"BLOCKED", origin, c.Request.Host, c.Request.Method, c.Request.URL.Path,
						c.ClientIP(), c.GetHeader("User-Agent"))
					resp.Error(c, http.StatusForbidden, "CORS policy violation")
					return
				}
			}
		}
		c.Next()
	}
}
