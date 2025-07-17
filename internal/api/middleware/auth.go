package middleware

import (
	"net/http"

	"github.com/bestruirui/bestsub/internal/config"
	"github.com/bestruirui/bestsub/internal/core/session"
	"github.com/bestruirui/bestsub/internal/models/api"
	"github.com/bestruirui/bestsub/internal/utils"
	"github.com/bestruirui/bestsub/internal/utils/log"
	"github.com/cespare/xxhash/v2"
	"github.com/gin-gonic/gin"
)

// Auth JWT认证中间件
func Auth() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, api.ResponseError{
				Code:    http.StatusUnauthorized,
				Message: "Unauthorized",
				Error:   "Authorization header is required",
			})
			c.Abort()
			return
		}
		token := authHeader[7:]

		claims, err := ValidateToken(token, config.Base().JWT.Secret)
		if err != nil {
			log.Warnf("JWT validation failed: %v", err)
			c.JSON(http.StatusUnauthorized, api.ResponseError{
				Code:    http.StatusUnauthorized,
				Message: "Unauthorized",
				Error:   "Invalid or expired token",
			})
			c.Abort()
			return
		}

		sess, err := session.Get(claims.SessionID)
		if err != nil {
			log.Warnf("Session not found: %v", err)
			c.JSON(http.StatusUnauthorized, api.ResponseError{
				Code:    http.StatusUnauthorized,
				Message: "Unauthorized",
				Error:   "Session not found or expired",
			})
			c.Abort()
			return
		}

		if !sess.IsActive {
			log.Warnf("Session %d is not active", claims.SessionID)
			c.JSON(http.StatusUnauthorized, api.ResponseError{
				Code:    http.StatusUnauthorized,
				Message: "Unauthorized",
				Error:   "Session is not active",
			})
			c.Abort()
			return
		}

		if sess.HashAToken != xxhash.Sum64String(token) {
			log.Warnf("Token hash mismatch: session=%d, request=%d", sess.HashAToken, xxhash.Sum64String(token))
			c.JSON(http.StatusUnauthorized, api.ResponseError{
				Code:    http.StatusUnauthorized,
				Message: "Unauthorized",
				Error:   "Token hash mismatch",
			})
			c.Abort()
			return
		}

		clientIP := utils.IPToUint32(c.ClientIP())
		if sess.ClientIP != clientIP {
			log.Warnf("Client IP mismatch: session=%s, request=%s",
				utils.Uint32ToIP(sess.ClientIP), c.ClientIP())
			c.JSON(http.StatusUnauthorized, api.ResponseError{
				Code:    http.StatusUnauthorized,
				Message: "Unauthorized",
				Error:   "Client IP mismatch",
			})
			c.Abort()
			return
		}
		userAgent := c.GetHeader("User-Agent")
		if sess.UserAgent != userAgent {
			log.Warnf("User-Agent mismatch for session %d", claims.SessionID)
			c.JSON(http.StatusUnauthorized, api.ResponseError{
				Code:    http.StatusUnauthorized,
				Message: "Unauthorized",
				Error:   "User-Agent mismatch",
			})
			c.Abort()
			return
		}

		c.Set("session_id", claims.SessionID)
		c.Next()
	}
}

// WSAuth WebSocket专用认证中间件
// WebSocket连接的认证处理与普通HTTP请求不同，需要特殊处理
func WSAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.Query("token")

		if token == "" {
			log.Warnf("WebSocket认证失败: 缺少token, IP=%s", c.ClientIP())
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		claims, err := ValidateToken(token, config.Base().JWT.Secret)
		if err != nil {
			log.Warnf("WebSocket JWT验证失败: %v, IP=%s", err, c.ClientIP())
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		sess, err := session.Get(claims.SessionID)
		if err != nil {
			log.Warnf("WebSocket会话未找到: %v, IP=%s", err, c.ClientIP())
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		if !sess.IsActive {
			log.Warnf("WebSocket session is not active, SessionID=%d, IP=%s", claims.SessionID, c.ClientIP())
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		clientIP := utils.IPToUint32(c.ClientIP())

		if sess.ClientIP != clientIP {
			log.Warnf("WebSocket client IP mismatch: session=%s, request=%s",
				utils.Uint32ToIP(sess.ClientIP), c.ClientIP())
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		c.Set("session_id", claims.SessionID)
		c.Next()
	}
}
