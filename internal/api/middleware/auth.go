package middleware

import (
	"context"
	"net/http"

	"github.com/bestruirui/bestsub/internal/api/models"
	"github.com/bestruirui/bestsub/internal/database"
	"github.com/bestruirui/bestsub/internal/utils/jwt"
	"github.com/bestruirui/bestsub/internal/utils/log"
	"github.com/gin-gonic/gin"
)

// AuthMiddleware JWT认证中间件
func Auth() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 从请求头中获取Authorization
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, models.ErrorResponse{
				Code:    http.StatusUnauthorized,
				Message: "Unauthorized",
				Error:   "Authorization header is required",
			})
			c.Abort()
			return
		}

		// 提取JWT令牌
		token, err := jwt.ExtractTokenFromHeader(authHeader)
		if err != nil {
			c.JSON(http.StatusUnauthorized, models.ErrorResponse{
				Code:    http.StatusUnauthorized,
				Message: "Unauthorized",
				Error:   "Invalid authorization header format",
			})
			c.Abort()
			return
		}

		// 验证JWT令牌
		claims, err := jwt.ValidateToken(token)
		if err != nil {
			log.Warnf("JWT validation failed: %v", err)
			c.JSON(http.StatusUnauthorized, models.ErrorResponse{
				Code:    http.StatusUnauthorized,
				Message: "Unauthorized",
				Error:   "Invalid or expired token",
			})
			c.Abort()
			return
		}

		// 检查令牌是否过期
		if jwt.IsTokenExpired(claims) {
			c.JSON(http.StatusUnauthorized, models.ErrorResponse{
				Code:    http.StatusUnauthorized,
				Message: "Unauthorized",
				Error:   "Token has expired",
			})
			c.Abort()
			return
		}

		// 验证会话是否存在且有效
		sessionRepo := database.Session()
		session, err := sessionRepo.GetByID(context.Background(), claims.SessionID)
		if err != nil {
			log.Errorf("Failed to get session: %v", err)
			c.JSON(http.StatusInternalServerError, models.ErrorResponse{
				Code:    http.StatusInternalServerError,
				Message: "Internal Server Error",
				Error:   "Failed to validate session",
			})
			c.Abort()
			return
		}

		if session == nil || !session.IsActive {
			c.JSON(http.StatusUnauthorized, models.ErrorResponse{
				Code:    http.StatusUnauthorized,
				Message: "Unauthorized",
				Error:   "Session is invalid or inactive",
			})
			c.Abort()
			return
		}

		// 验证令牌哈希是否匹配
		tokenHash := jwt.HashToken(token)
		if session.TokenHash != tokenHash {
			c.JSON(http.StatusUnauthorized, models.ErrorResponse{
				Code:    http.StatusUnauthorized,
				Message: "Unauthorized",
				Error:   "Token hash mismatch",
			})
			c.Abort()
			return
		}

		// 获取用户信息（单用户系统，直接从Auth表获取）
		authRepo := database.Auth()
		authInfo, err := authRepo.Get(context.Background())
		if err != nil {
			log.Errorf("Failed to get auth info: %v", err)
			c.JSON(http.StatusInternalServerError, models.ErrorResponse{
				Code:    http.StatusInternalServerError,
				Message: "Internal Server Error",
				Error:   "Failed to get user information",
			})
			c.Abort()
			return
		}

		if authInfo == nil {
			c.JSON(http.StatusUnauthorized, models.ErrorResponse{
				Code:    http.StatusUnauthorized,
				Message: "Unauthorized",
				Error:   "User not found",
			})
			c.Abort()
			return
		}

		// 将用户信息存储到上下文中
		c.Set("username", authInfo.UserName)
		c.Set("session_id", claims.SessionID)
		c.Set("claims", claims)

		c.Next()
	}
}
