package middleware

import (
	"net/http"
	"os"
	"path"
	"strings"

	"github.com/bestruirui/bestsub/internal/config"
	"github.com/gin-gonic/gin"
)

var staticExtensions = map[string]bool{
	".js":  true,
	".css": true,
	".mjs": true,

	".png":  true,
	".jpg":  true,
	".jpeg": true,
	".gif":  true,
	".svg":  true,
	".ico":  true,
	".webp": true,
	".avif": true,
	".bmp":  true,

	".woff":  true,
	".woff2": true,
	".ttf":   true,
	".eot":   true,
	".otf":   true,

	".xml":  true,
	".json": true,
	".txt":  true,
	".pdf":  true,
}

var (
	cacheOneHourHeader  = "public, max-age=3600"
	cacheOneWeekHeader  = "public, max-age=604800"
	cacheOneMonthHeader = "public, max-age=2592000"
	cacheOneYearHeader  = "public, max-age=31536000"

	fileServer http.Handler
)

func init() {

	fileServer = http.FileServer(http.Dir(config.Base().Server.UIPath))
}

func Static() gin.HandlerFunc {
	return func(c *gin.Context) {
		reqPath := c.Request.URL.Path

		if strings.HasPrefix(reqPath, "/api") {
			c.Next()
			return
		}

		if reqPath == "/" || reqPath == "" {
			reqPath = "/index.html"
		}

		ext := path.Ext(reqPath)

		switch ext {
		case ".js", ".css", ".mjs":
			c.Header("Cache-Control", cacheOneYearHeader)
		case ".png", ".jpg", ".jpeg", ".gif", ".svg", ".ico", ".webp", ".avif", ".bmp":
			c.Header("Cache-Control", cacheOneMonthHeader)
		case ".woff", ".woff2", ".ttf", ".eot", ".otf":
			c.Header("Cache-Control", cacheOneYearHeader)
		case ".xml", ".json", ".txt":
			c.Header("Cache-Control", cacheOneHourHeader)
		case ".html":
			c.Header("Cache-Control", cacheOneHourHeader)
		default:
			c.Header("Cache-Control", cacheOneWeekHeader)
		}

		if _, err := os.Stat(config.Base().Server.UIPath + reqPath); err != nil {
			if staticExtensions[ext] {
				c.Status(http.StatusNotFound)
				return
			}
			c.Request.URL.Path = "/index.html"
		}
		fileServer.ServeHTTP(c.Writer, c.Request)
	}
}
