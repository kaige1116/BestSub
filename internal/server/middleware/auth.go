package middleware

import (
	"net/http"

	"github.com/bestruirui/bestsub/internal/config"
	"github.com/bestruirui/bestsub/internal/server/auth"
	"github.com/bestruirui/bestsub/internal/server/resp"
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
			resp.Error(c, http.StatusUnauthorized, "Authorization header is required")
			c.Abort()
			return
		}
		token := authHeader[7:]

		claims, err := auth.ValidateToken(token, config.Base().JWT.Secret)
		if err != nil {
			log.Warnf("JWT validation failed: %v", err)
			resp.Error(c, http.StatusUnauthorized, "Invalid or expired token")
			c.Abort()
			return
		}

		sess, err := auth.GetSession(claims.SessionID)
		if err != nil {
			log.Warnf("Session not found: %v", err)
			resp.Error(c, http.StatusUnauthorized, "Session not found or expired")
			c.Abort()
			return
		}

		if !sess.IsActive {
			log.Warnf("Session %d is not active", claims.SessionID)
			resp.Error(c, http.StatusUnauthorized, "Session is not active")
			c.Abort()
			return
		}

		if sess.HashAToken != xxhash.Sum64String(token) {
			log.Warnf("Token hash mismatch: session=%d, request=%d", sess.HashAToken, xxhash.Sum64String(token))
			resp.Error(c, http.StatusUnauthorized, "Token hash mismatch")
			c.Abort()
			return
		}

		clientIP := utils.IPToUint32(c.ClientIP())
		if sess.ClientIP != clientIP {
			log.Warnf("Client IP mismatch: session=%s, request=%s",
				utils.Uint32ToIP(sess.ClientIP), c.ClientIP())
			resp.Error(c, http.StatusUnauthorized, "Client IP mismatch")
			c.Abort()
			return
		}
		userAgent := c.GetHeader("User-Agent")
		if sess.UserAgent != userAgent {
			log.Warnf("User-Agent mismatch for session %d", claims.SessionID)
			resp.Error(c, http.StatusUnauthorized, "User-Agent mismatch")
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
			log.Warnf("WebSocket authentication failed: missing token, IP=%s", c.ClientIP())
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		claims, err := auth.ValidateToken(token, config.Base().JWT.Secret)
		if err != nil {
			log.Warnf("WebSocket JWT validation failed: %v, IP=%s", err, c.ClientIP())
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		sess, err := auth.GetSession(claims.SessionID)
		if err != nil {
			log.Warnf("WebSocket session not found: %v, IP=%s", err, c.ClientIP())
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
