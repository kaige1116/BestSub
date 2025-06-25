package handlers

import (
	"context"
	"fmt"
	"net/http"

	"github.com/bestruirui/bestsub/internal/api/models"
	"github.com/bestruirui/bestsub/internal/database"
	dbModels "github.com/bestruirui/bestsub/internal/database/models"
	"github.com/bestruirui/bestsub/internal/utils/jwt"
	"github.com/bestruirui/bestsub/internal/utils/log"
	"github.com/bestruirui/bestsub/internal/utils/time"
	"github.com/gin-gonic/gin"
)

// AuthHandler 认证处理器
type AuthHandler struct{}

// NewAuthHandler 创建认证处理器
func NewAuthHandler() *AuthHandler {
	return &AuthHandler{}
}

// Login 用户登录
// @Summary 用户登录
// @Description 用户登录接口，验证用户名和密码，返回JWT令牌
// @Tags 认证
// @Accept json
// @Produce json
// @Param request body models.LoginRequest true "登录请求"
// @Success 200 {object} models.SuccessResponse{data=models.LoginResponse} "登录成功"
// @Failure 400 {object} models.ErrorResponse "请求参数错误"
// @Failure 401 {object} models.ErrorResponse "用户名或密码错误"
// @Failure 500 {object} models.ErrorResponse "服务器内部错误"
// @Router /api/v1/auth/login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	var req models.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Code:    http.StatusBadRequest,
			Message: "Bad Request",
			Error:   "Invalid request format: " + err.Error(),
		})
		return
	}

	// 验证用户名和密码
	authRepo := database.Auth()
	err := authRepo.VerifyPassword(context.Background(), req.Username, req.Password)
	if err != nil {
		log.Warnf("Login failed for user %s: %v", req.Username, err)
		c.JSON(http.StatusUnauthorized, models.ErrorResponse{
			Code:    http.StatusUnauthorized,
			Message: "Unauthorized",
			Error:   "Invalid username or password",
		})
		return
	}

	// 获取用户信息
	authInfo, err := authRepo.Get(context.Background())
	if err != nil {
		log.Errorf("Failed to get auth info: %v", err)
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Code:    http.StatusInternalServerError,
			Message: "Internal Server Error",
			Error:   "Failed to get user information",
		})
		return
	}

	// 创建会话记录
	sessionRepo := database.Session()
	session := &dbModels.Session{
		IPAddress: c.ClientIP(),
		UserAgent: c.GetHeader("User-Agent"),
		IsActive:  true,
	}

	// 先创建会话以获取ID
	err = sessionRepo.Create(context.Background(), session)
	if err != nil {
		log.Errorf("Failed to create session: %v", err)
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Code:    http.StatusInternalServerError,
			Message: "Internal Server Error",
			Error:   "Failed to create session",
		})
		return
	}

	// 生成JWT令牌对
	tokenPair, err := jwt.GenerateTokenPair(session.ID)
	if err != nil {
		log.Errorf("Failed to generate token pair: %v", err)
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Code:    http.StatusInternalServerError,
			Message: "Internal Server Error",
			Error:   "Failed to generate tokens",
		})
		return
	}

	// 更新会话信息
	session.TokenHash = tokenPair.TokenHash
	session.RefreshToken = tokenPair.RefreshToken
	session.ExpiresAt = tokenPair.ExpiresAt
	err = sessionRepo.Update(context.Background(), session)
	if err != nil {
		log.Errorf("Failed to update session: %v", err)
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Code:    http.StatusInternalServerError,
			Message: "Internal Server Error",
			Error:   "Failed to update session",
		})
		return
	}

	// 构建响应
	response := models.LoginResponse{
		AccessToken:  tokenPair.AccessToken,
		RefreshToken: tokenPair.RefreshToken,
		ExpiresAt:    tokenPair.ExpiresAt,
		User: models.UserInfo{
			Username:  authInfo.UserName,
			CreatedAt: authInfo.CreatedAt,
			UpdatedAt: authInfo.UpdatedAt,
		},
	}

	log.Infof("User %s logged in successfully from %s", req.Username, c.ClientIP())

	c.JSON(http.StatusOK, models.SuccessResponse{
		Code:    http.StatusOK,
		Message: "Login successful",
		Data:    response,
	})
}

// RefreshToken 刷新令牌
// @Summary 刷新访问令牌
// @Description 使用刷新令牌获取新的访问令牌
// @Tags 认证
// @Accept json
// @Produce json
// @Param request body models.RefreshTokenRequest true "刷新令牌请求"
// @Success 200 {object} models.SuccessResponse{data=models.RefreshTokenResponse} "刷新成功"
// @Failure 400 {object} models.ErrorResponse "请求参数错误"
// @Failure 401 {object} models.ErrorResponse "刷新令牌无效"
// @Failure 500 {object} models.ErrorResponse "服务器内部错误"
// @Router /api/v1/auth/refresh [post]
func (h *AuthHandler) RefreshToken(c *gin.Context) {
	var req models.RefreshTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Code:    http.StatusBadRequest,
			Message: "Bad Request",
			Error:   "Invalid request format: " + err.Error(),
		})
		return
	}

	// 验证刷新令牌
	claims, err := jwt.ValidateRefreshToken(req.RefreshToken)
	if err != nil {
		log.Warnf("Refresh token validation failed: %v", err)
		c.JSON(http.StatusUnauthorized, models.ErrorResponse{
			Code:    http.StatusUnauthorized,
			Message: "Unauthorized",
			Error:   "Invalid refresh token",
		})
		return
	}

	// 验证会话是否存在且有效
	sessionRepo := database.Session()
	session, err := sessionRepo.GetByRefreshToken(context.Background(), req.RefreshToken)
	if err != nil {
		log.Errorf("Failed to get session by refresh token: %v", err)
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Code:    http.StatusInternalServerError,
			Message: "Internal Server Error",
			Error:   "Failed to validate session",
		})
		return
	}

	if session == nil || !session.IsActive {
		c.JSON(http.StatusUnauthorized, models.ErrorResponse{
			Code:    http.StatusUnauthorized,
			Message: "Unauthorized",
			Error:   "Session is invalid or inactive",
		})
		return
	}

	// 生成新的令牌对
	newTokenPair, err := jwt.GenerateTokenPair(session.ID)
	if err != nil {
		log.Errorf("Failed to generate new token pair: %v", err)
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Code:    http.StatusInternalServerError,
			Message: "Internal Server Error",
			Error:   "Failed to generate new tokens",
		})
		return
	}

	// 更新会话信息
	session.TokenHash = newTokenPair.TokenHash
	session.RefreshToken = newTokenPair.RefreshToken
	session.ExpiresAt = newTokenPair.ExpiresAt
	session.IPAddress = c.ClientIP()
	session.UserAgent = c.GetHeader("User-Agent")
	err = sessionRepo.Update(context.Background(), session)
	if err != nil {
		log.Errorf("Failed to update session: %v", err)
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Code:    http.StatusInternalServerError,
			Message: "Internal Server Error",
			Error:   "Failed to update session",
		})
		return
	}

	// 构建响应
	response := models.RefreshTokenResponse{
		AccessToken:  newTokenPair.AccessToken,
		RefreshToken: newTokenPair.RefreshToken,
		ExpiresAt:    newTokenPair.ExpiresAt,
	}

	log.Infof("Token refreshed for session %d", claims.SessionID)

	c.JSON(http.StatusOK, models.SuccessResponse{
		Code:    http.StatusOK,
		Message: "Token refreshed successfully",
		Data:    response,
	})
}

// Logout 用户登出
// @Summary 用户登出
// @Description 用户登出接口，使当前会话失效
// @Tags 认证
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} models.SuccessResponse "登出成功"
// @Failure 401 {object} models.ErrorResponse "未授权"
// @Failure 500 {object} models.ErrorResponse "服务器内部错误"
// @Router /api/v1/auth/logout [post]
func (h *AuthHandler) Logout(c *gin.Context) {
	sessionID, exists := c.Get("session_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, models.ErrorResponse{
			Code:    http.StatusUnauthorized,
			Message: "Unauthorized",
			Error:   "Session not found",
		})
		return
	}

	username, _ := c.Get("username")

	// 停用当前会话
	sessionRepo := database.Session()
	session, err := sessionRepo.GetByID(context.Background(), sessionID.(int64))
	if err != nil {
		log.Errorf("Failed to get session: %v", err)
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Code:    http.StatusInternalServerError,
			Message: "Internal Server Error",
			Error:   "Failed to get session",
		})
		return
	}

	if session != nil {
		session.IsActive = false
		err = sessionRepo.Update(context.Background(), session)
		if err != nil {
			log.Errorf("Failed to deactivate session: %v", err)
			c.JSON(http.StatusInternalServerError, models.ErrorResponse{
				Code:    http.StatusInternalServerError,
				Message: "Internal Server Error",
				Error:   "Failed to logout",
			})
			return
		}
	}

	log.Infof("User %s logged out successfully", username)

	c.JSON(http.StatusOK, models.SuccessResponse{
		Code:    http.StatusOK,
		Message: "Logout successful",
	})
}

// ChangePassword 修改密码
// @Summary 修改密码
// @Description 修改当前用户的密码
// @Tags 认证
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body models.ChangePasswordRequest true "修改密码请求"
// @Success 200 {object} models.SuccessResponse "密码修改成功"
// @Failure 400 {object} models.ErrorResponse "请求参数错误"
// @Failure 401 {object} models.ErrorResponse "未授权或旧密码错误"
// @Failure 500 {object} models.ErrorResponse "服务器内部错误"
// @Router /api/v1/auth/change-password [post]
func (h *AuthHandler) ChangePassword(c *gin.Context) {
	var req models.ChangePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Code:    http.StatusBadRequest,
			Message: "Bad Request",
			Error:   "Invalid request format: " + err.Error(),
		})
		return
	}

	username, exists := c.Get("username")
	if !exists {
		c.JSON(http.StatusUnauthorized, models.ErrorResponse{
			Code:    http.StatusUnauthorized,
			Message: "Unauthorized",
			Error:   "User not found in context",
		})
		return
	}

	// 验证旧密码
	authRepo := database.Auth()
	err := authRepo.VerifyPassword(context.Background(), username.(string), req.OldPassword)
	if err != nil {
		log.Warnf("Change password failed for user %s: old password verification failed", username)
		c.JSON(http.StatusUnauthorized, models.ErrorResponse{
			Code:    http.StatusUnauthorized,
			Message: "Unauthorized",
			Error:   "Old password is incorrect",
		})
		return
	}

	// 获取当前用户信息
	authInfo, err := authRepo.Get(context.Background())
	if err != nil {
		log.Errorf("Failed to get auth info: %v", err)
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Code:    http.StatusInternalServerError,
			Message: "Internal Server Error",
			Error:   "Failed to get user information",
		})
		return
	}

	// 更新密码
	authInfo.Password = req.NewPassword
	authInfo.UpdatedAt = time.Now()
	err = authRepo.Update(context.Background(), authInfo)
	if err != nil {
		log.Errorf("Failed to update password: %v", err)
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Code:    http.StatusInternalServerError,
			Message: "Internal Server Error",
			Error:   "Failed to update password",
		})
		return
	}

	// 删除所有会话（强制重新登录）
	sessionRepo := database.Session()
	err = sessionRepo.DeleteAll(context.Background())
	if err != nil {
		log.Errorf("Failed to delete all sessions: %v", err)
	}

	log.Infof("Password changed successfully for user %s", username)

	c.JSON(http.StatusOK, models.SuccessResponse{
		Code:    http.StatusOK,
		Message: "Password changed successfully. Please login again.",
	})
}

// GetUserInfo 获取当前用户信息
// @Summary 获取用户信息
// @Description 获取当前登录用户的详细信息
// @Tags 认证
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} models.SuccessResponse{data=models.UserInfo} "获取成功"
// @Failure 401 {object} models.ErrorResponse "未授权"
// @Failure 500 {object} models.ErrorResponse "服务器内部错误"
// @Router /api/v1/auth/user [get]
func (h *AuthHandler) GetUserInfo(c *gin.Context) {
	_, exists := c.Get("username")
	if !exists {
		c.JSON(http.StatusUnauthorized, models.ErrorResponse{
			Code:    http.StatusUnauthorized,
			Message: "Unauthorized",
			Error:   "User not found in context",
		})
		return
	}

	// 获取用户信息
	authRepo := database.Auth()
	authInfo, err := authRepo.Get(context.Background())
	if err != nil {
		log.Errorf("Failed to get auth info: %v", err)
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Code:    http.StatusInternalServerError,
			Message: "Internal Server Error",
			Error:   "Failed to get user information",
		})
		return
	}

	if authInfo == nil {
		c.JSON(http.StatusNotFound, models.ErrorResponse{
			Code:    http.StatusNotFound,
			Message: "Not Found",
			Error:   "User not found",
		})
		return
	}

	userInfo := models.UserInfo{
		Username:  authInfo.UserName,
		CreatedAt: authInfo.CreatedAt,
		UpdatedAt: authInfo.UpdatedAt,
	}

	c.JSON(http.StatusOK, models.SuccessResponse{
		Code:    http.StatusOK,
		Message: "User information retrieved successfully",
		Data:    userInfo,
	})
}

// GetSessions 获取当前用户的所有会话
// @Summary 获取用户会话列表
// @Description 获取当前用户的所有活跃会话信息
// @Tags 认证
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} models.SuccessResponse{data=models.SessionListResponse} "获取成功"
// @Failure 401 {object} models.ErrorResponse "未授权"
// @Failure 500 {object} models.ErrorResponse "服务器内部错误"
// @Router /api/v1/auth/sessions [get]
func (h *AuthHandler) GetSessions(c *gin.Context) {
	username, exists := c.Get("username")
	if !exists {
		c.JSON(http.StatusUnauthorized, models.ErrorResponse{
			Code:    http.StatusUnauthorized,
			Message: "Unauthorized",
			Error:   "User not found in context",
		})
		return
	}

	// 获取所有活跃会话
	sessionRepo := database.Session()
	sessions, err := sessionRepo.GetAllActive(context.Background())
	if err != nil {
		log.Errorf("Failed to get active sessions: %v", err)
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Code:    http.StatusInternalServerError,
			Message: "Internal Server Error",
			Error:   "Failed to get sessions",
		})
		return
	}

	// 转换为响应模型
	sessionInfos := make([]models.SessionInfo, 0, len(sessions))
	for _, session := range sessions {
		sessionInfos = append(sessionInfos, models.SessionInfo{
			ID:        session.ID,
			IPAddress: session.IPAddress,
			UserAgent: session.UserAgent,
			IsActive:  session.IsActive,
			ExpiresAt: session.ExpiresAt,
			CreatedAt: session.CreatedAt,
			UpdatedAt: session.UpdatedAt,
		})
	}

	response := models.SessionListResponse{
		Sessions: sessionInfos,
		Total:    len(sessionInfos),
	}

	log.Debugf("Retrieved %d sessions for user %s", len(sessionInfos), username)

	c.JSON(http.StatusOK, models.SuccessResponse{
		Code:    http.StatusOK,
		Message: "Sessions retrieved successfully",
		Data:    response,
	})
}

// DeleteSession 删除指定会话
// @Summary 删除会话
// @Description 删除指定ID的会话，使其失效
// @Tags 认证
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "会话ID"
// @Success 200 {object} models.SuccessResponse "删除成功"
// @Failure 400 {object} models.ErrorResponse "请求参数错误"
// @Failure 401 {object} models.ErrorResponse "未授权"
// @Failure 404 {object} models.ErrorResponse "会话不存在"
// @Failure 500 {object} models.ErrorResponse "服务器内部错误"
// @Router /api/v1/auth/sessions/{id} [delete]
func (h *AuthHandler) DeleteSession(c *gin.Context) {
	sessionIDStr := c.Param("id")
	if sessionIDStr == "" {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Code:    http.StatusBadRequest,
			Message: "Bad Request",
			Error:   "Session ID is required",
		})
		return
	}

	sessionID := int64(0)
	if _, err := fmt.Sscanf(sessionIDStr, "%d", &sessionID); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Code:    http.StatusBadRequest,
			Message: "Bad Request",
			Error:   "Invalid session ID format",
		})
		return
	}

	username, exists := c.Get("username")
	if !exists {
		c.JSON(http.StatusUnauthorized, models.ErrorResponse{
			Code:    http.StatusUnauthorized,
			Message: "Unauthorized",
			Error:   "User not found in context",
		})
		return
	}

	// 获取会话信息
	sessionRepo := database.Session()
	session, err := sessionRepo.GetByID(context.Background(), sessionID)
	if err != nil {
		log.Errorf("Failed to get session: %v", err)
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Code:    http.StatusInternalServerError,
			Message: "Internal Server Error",
			Error:   "Failed to get session",
		})
		return
	}

	if session == nil {
		c.JSON(http.StatusNotFound, models.ErrorResponse{
			Code:    http.StatusNotFound,
			Message: "Not Found",
			Error:   "Session not found",
		})
		return
	}

	// 删除会话
	err = sessionRepo.Delete(context.Background(), sessionID)
	if err != nil {
		log.Errorf("Failed to delete session: %v", err)
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Code:    http.StatusInternalServerError,
			Message: "Internal Server Error",
			Error:   "Failed to delete session",
		})
		return
	}

	log.Infof("Session %d deleted by user %s", sessionID, username)

	c.JSON(http.StatusOK, models.SuccessResponse{
		Code:    http.StatusOK,
		Message: "Session deleted successfully",
	})
}

// UpdateUsername 修改用户名
// @Summary 修改用户名
// @Description 修改当前用户的用户名
// @Tags 认证
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body models.UpdateUserInfoRequest true "修改用户名请求"
// @Success 200 {object} models.SuccessResponse "用户名修改成功"
// @Failure 400 {object} models.ErrorResponse "请求参数错误"
// @Failure 401 {object} models.ErrorResponse "未授权"
// @Failure 409 {object} models.ErrorResponse "用户名已存在"
// @Failure 500 {object} models.ErrorResponse "服务器内部错误"
// @Router /api/v1/auth/update-username [post]
func (h *AuthHandler) UpdateUsername(c *gin.Context) {
	var req models.UpdateUserInfoRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Code:    http.StatusBadRequest,
			Message: "Bad Request",
			Error:   "Invalid request format: " + err.Error(),
		})
		return
	}

	currentUsername, exists := c.Get("username")
	if !exists {
		c.JSON(http.StatusUnauthorized, models.ErrorResponse{
			Code:    http.StatusUnauthorized,
			Message: "Unauthorized",
			Error:   "User not found in context",
		})
		return
	}

	// 检查新用户名是否与当前用户名相同
	if req.Username == currentUsername.(string) {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Code:    http.StatusBadRequest,
			Message: "Bad Request",
			Error:   "New username cannot be the same as current username",
		})
		return
	}

	// 更新用户名
	authRepo := database.Auth()
	err := authRepo.UpdateUsername(context.Background(), req.Username)
	if err != nil {
		log.Errorf("Failed to update username from %s to %s: %v", currentUsername, req.Username, err)
		// 根据错误类型返回不同的状态码
		// 这里假设数据库层会返回适当的错误，如果用户名已存在等
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Code:    http.StatusInternalServerError,
			Message: "Internal Server Error",
			Error:   "Failed to update username",
		})
		return
	}

	log.Infof("Username changed successfully from %s to %s", currentUsername, req.Username)

	c.JSON(http.StatusOK, models.SuccessResponse{
		Code:    http.StatusOK,
		Message: "Username updated successfully.",
	})
}
