package handlers

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/bestruirui/bestsub/internal/api/middleware"
	"github.com/bestruirui/bestsub/internal/api/router"
	"github.com/bestruirui/bestsub/internal/core/session"
	"github.com/bestruirui/bestsub/internal/database"
	"github.com/bestruirui/bestsub/internal/models/api"
	"github.com/bestruirui/bestsub/internal/models/auth"
	"github.com/bestruirui/bestsub/internal/utils"
	"github.com/bestruirui/bestsub/internal/utils/jwt"
	"github.com/bestruirui/bestsub/internal/utils/local"
	"github.com/bestruirui/bestsub/internal/utils/log"
	"github.com/bestruirui/bestsub/internal/utils/passwd"
	"github.com/gin-gonic/gin"
)

// authHandler 认证处理器
type authHandler struct{}

// init 函数用于自动注册路由
func init() {
	h := newAuthHandler()

	// 公开的认证路由（无需认证）
	router.NewGroupRouter("/api/v1/auth").
		AddRoute(
			router.NewRoute("/login", router.POST).
				Handle(h.login).
				WithDescription("User login"),
		).
		AddRoute(
			router.NewRoute("/refresh", router.POST).
				Handle(h.refreshToken).
				WithDescription("Refresh access token"),
		)

	// 需要认证的路由
	router.NewGroupRouter("/api/v1/auth").
		Use(middleware.Auth()).
		AddRoute(
			router.NewRoute("/logout", router.POST).
				Handle(h.logout).
				WithDescription("User logout"),
		).
		AddRoute(
			router.NewRoute("/user/password", router.POST).
				Handle(h.changePassword).
				WithDescription("Change user password"),
		).
		AddRoute(
			router.NewRoute("/user/name", router.POST).
				Handle(h.updateUsername).
				WithDescription("Update username"),
		).
		AddRoute(
			router.NewRoute("/user", router.GET).
				Handle(h.getUserInfo).
				WithDescription("Get user information"),
		).
		AddRoute(
			router.NewRoute("/sessions", router.GET).
				Handle(h.getSessions).
				WithDescription("Get user sessions"),
		).
		AddRoute(
			router.NewRoute("/sessions/:id", router.DELETE).
				Handle(h.deleteSession).
				WithDescription("Delete session"),
		)
}

// newAuthHandler 创建认证处理器
func newAuthHandler() *authHandler {
	return &authHandler{}
}

// login 用户登录
// @Summary 用户登录
// @Description 用户登录接口，验证用户名和密码，返回JWT令牌
// @Tags 认证
// @Accept json
// @Produce json
// @Param request body auth.LoginRequest true "登录请求"
// @Success 200 {object} api.ResponseSuccess{data=auth.LoginResponse} "登录成功"
// @Failure 400 {object} api.ResponseError "请求参数错误"
// @Failure 401 {object} api.ResponseError "用户名或密码错误"
// @Failure 500 {object} api.ResponseError "服务器内部错误"
// @Router /api/v1/auth/login [post]
func (h *authHandler) login(c *gin.Context) {
	var req auth.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, api.ResponseError{
			Code:    http.StatusBadRequest,
			Message: "Bad Request",
			Error:   "Invalid request format: " + err.Error(),
		})
		return
	}

	authRepo := database.Auth()
	authInfo, err := authRepo.Get(context.Background())

	if err != nil {
		log.Errorf("Failed to get auth info: %v from %s", err, c.ClientIP())
		c.JSON(http.StatusInternalServerError, api.ResponseError{
			Code:    http.StatusInternalServerError,
			Message: "Unauthorized",
			Error:   "Invalid username or password",
		})
		return
	}

	err = passwd.Verify(req.Password, authInfo.Password)
	if err != nil {
		log.Warnf("Login failed for user %s: %v from %s", req.Username, err, c.ClientIP())
		c.JSON(http.StatusUnauthorized, api.ResponseError{
			Code:    http.StatusUnauthorized,
			Message: "Unauthorized",
			Error:   "Invalid username or password",
		})
		return
	}

	sessionID, tempSess := session.GetUnActive()
	if tempSess == nil {
		log.Errorf("No unused session found from %s", c.ClientIP())
		c.JSON(http.StatusInternalServerError, api.ResponseError{
			Code:    http.StatusInternalServerError,
			Message: "Internal Server Error",
			Error:   "No unused session found",
		})
		return
	}

	tokenPair, err := jwt.GenerateTokenPair(sessionID, authInfo.UserName)
	if err != nil {
		log.Errorf("Failed to generate token pair: %v from %s", err, c.ClientIP())
		session.Disable(sessionID)
		c.JSON(http.StatusInternalServerError, api.ResponseError{
			Code:    http.StatusInternalServerError,
			Message: "Internal Server Error",
			Error:   "Failed to generate tokens",
		})
		return
	}

	now := uint32(local.Time().Unix())
	clientIP := utils.IPToUint32(c.ClientIP())

	tempSess.IsActive = true
	tempSess.ClientIP = clientIP
	tempSess.UserAgent = c.GetHeader("User-Agent")
	tempSess.CreatedAt = now
	tempSess.LastAccessAt = now
	tempSess.ExpiresAt = uint32(tokenPair.ExpiresAt.Unix())

	// 构建响应
	response := auth.LoginResponse{
		AccessToken:  tokenPair.AccessToken,
		RefreshToken: tokenPair.RefreshToken,
		ExpiresAt:    tokenPair.ExpiresAt,
		User: auth.Data{
			UserName:  authInfo.UserName,
			CreatedAt: authInfo.CreatedAt,
			UpdatedAt: authInfo.UpdatedAt,
		},
	}

	log.Infof("User %s logged in successfully from %s", req.Username, c.ClientIP())

	c.JSON(http.StatusOK, api.ResponseSuccess{
		Code:    http.StatusOK,
		Message: "Login successful",
		Data:    response,
	})
}

// refreshToken 刷新令牌
// @Summary 刷新访问令牌
// @Description 使用刷新令牌获取新的访问令牌
// @Tags 认证
// @Accept json
// @Produce json
// @Param request body auth.RefreshTokenRequest true "刷新令牌请求"
// @Success 200 {object} api.ResponseSuccess{data=auth.RefreshTokenResponse} "刷新成功"
// @Failure 400 {object} api.ResponseError "请求参数错误"
// @Failure 401 {object} api.ResponseError "刷新令牌无效"
// @Failure 500 {object} api.ResponseError "服务器内部错误"
// @Router /api/v1/auth/refresh [post]
func (h *authHandler) refreshToken(c *gin.Context) {
	var req auth.RefreshTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, api.ResponseError{
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
		c.JSON(http.StatusUnauthorized, api.ResponseError{
			Code:    http.StatusUnauthorized,
			Message: "Unauthorized",
			Error:   "Invalid refresh token",
		})
		return
	}

	sessionID := claims.SessionID

	sess, err := session.Get(sessionID)
	if err != nil {
		log.Errorf("Failed to get session by ID: %v", err)
		c.JSON(http.StatusUnauthorized, api.ResponseError{
			Code:    http.StatusUnauthorized,
			Message: "Unauthorized",
			Error:   "Session not found or expired",
		})
		return
	}

	// 验证客户端IP
	clientIP := utils.IPToUint32(c.ClientIP())

	// 验证IP和UserAgent
	if sess.ClientIP != clientIP {
		log.Warnf("Client IP mismatch during token refresh: session=%s, request=%s",
			utils.Uint32ToIP(sess.ClientIP), c.ClientIP())
		c.JSON(http.StatusUnauthorized, api.ResponseError{
			Code:    http.StatusUnauthorized,
			Message: "Unauthorized",
			Error:   "Client IP mismatch",
		})
		return
	}

	if sess.UserAgent != c.GetHeader("User-Agent") {
		log.Warnf("User agent mismatch during token refresh: session=%s, request=%s",
			sess.UserAgent, c.GetHeader("User-Agent"))
		c.JSON(http.StatusUnauthorized, api.ResponseError{
			Code:    http.StatusUnauthorized,
			Message: "Unauthorized",
			Error:   "User agent mismatch",
		})
		return
	}

	// 生成新的令牌对
	newTokenPair, err := jwt.GenerateTokenPair(sessionID, claims.Username)
	if err != nil {
		log.Errorf("Failed to generate new token pair: %v", err)
		c.JSON(http.StatusInternalServerError, api.ResponseError{
			Code:    http.StatusInternalServerError,
			Message: "Internal Server Error",
			Error:   "Failed to generate new tokens",
		})
		return
	}

	// 更新会话信息
	sess.ExpiresAt = uint32(newTokenPair.ExpiresAt.Unix())
	sess.LastAccessAt = uint32(local.Time().Unix())

	// 构建响应
	response := auth.RefreshTokenResponse{
		AccessToken:  newTokenPair.AccessToken,
		RefreshToken: newTokenPair.RefreshToken,
		ExpiresAt:    newTokenPair.ExpiresAt,
	}

	log.Infof("Token refreshed for session %d from %s", claims.SessionID, c.ClientIP())

	c.JSON(http.StatusOK, api.ResponseSuccess{
		Code:    http.StatusOK,
		Message: "Token refreshed successfully",
		Data:    response,
	})
}

// logout 用户登出
// @Summary 用户登出
// @Description 用户登出接口，使当前会话失效
// @Tags 认证
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} api.ResponseSuccess "登出成功"
// @Failure 401 {object} api.ResponseError "未授权"
// @Failure 500 {object} api.ResponseError "服务器内部错误"
// @Router /api/v1/auth/logout [post]
func (h *authHandler) logout(c *gin.Context) {
	sessionID, exists := c.Get("session_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, api.ResponseError{
			Code:    http.StatusUnauthorized,
			Message: "Unauthorized",
			Error:   "Session not found",
		})
		return
	}

	username, _ := c.Get("username")

	// 停用当前会话
	err := session.Disable(sessionID.(uint8))
	if err != nil {
		log.Errorf("Failed to disable session: %v", err)
		c.JSON(http.StatusInternalServerError, api.ResponseError{
			Code:    http.StatusInternalServerError,
			Message: "Internal Server Error",
			Error:   "Failed to logout",
		})
		return
	}

	log.Infof("User %s logged out successfully from %s", username, c.ClientIP())

	c.JSON(http.StatusOK, api.ResponseSuccess{
		Code:    http.StatusOK,
		Message: "Logout successful",
	})
}

// changePassword 修改密码
// @Summary 修改密码
// @Description 修改当前用户的密码
// @Tags 认证
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body auth.ChangePasswordRequest true "修改密码请求"
// @Success 200 {object} api.ResponseSuccess "密码修改成功"
// @Failure 400 {object} api.ResponseError "请求参数错误"
// @Failure 401 {object} api.ResponseError "未授权或旧密码错误"
// @Failure 500 {object} api.ResponseError "服务器内部错误"
// @Router /api/v1/auth/user/password [post]
func (h *authHandler) changePassword(c *gin.Context) {
	var req auth.ChangePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, api.ResponseError{
			Code:    http.StatusBadRequest,
			Message: "Bad Request",
			Error:   "Invalid request format: " + err.Error(),
		})
		return
	}

	username, exists := c.Get("username")
	if !exists {
		c.JSON(http.StatusUnauthorized, api.ResponseError{
			Code:    http.StatusUnauthorized,
			Message: "Unauthorized",
			Error:   "User not found in context",
		})
		return
	}

	// 获取当前用户信息并验证旧密码
	authRepo := database.Auth()
	authInfo, err := authRepo.Get(context.Background())
	if err != nil {
		log.Errorf("Failed to get auth info: %v", err)
		c.JSON(http.StatusInternalServerError, api.ResponseError{
			Code:    http.StatusInternalServerError,
			Message: "Internal Server Error",
			Error:   "Failed to get user information",
		})
		return
	}

	// 验证旧密码
	err = passwd.Verify(req.OldPassword, authInfo.Password)
	if err != nil {
		log.Warnf("Change password failed for user %s: old password verification failed from %s", username, c.ClientIP())
		c.JSON(http.StatusUnauthorized, api.ResponseError{
			Code:    http.StatusUnauthorized,
			Message: "Unauthorized",
			Error:   "Old password is incorrect",
		})
		return
	}

	// 加密新密码
	hashedPassword, err := passwd.Hash(req.NewPassword)
	if err != nil {
		log.Errorf("Failed to hash new password: %v", err)
		c.JSON(http.StatusInternalServerError, api.ResponseError{
			Code:    http.StatusInternalServerError,
			Message: "Internal Server Error",
			Error:   "Failed to process new password",
		})
		return
	}

	// 更新密码
	authInfo.Password = hashedPassword
	authInfo.UpdatedAt = time.Now()
	err = authRepo.Update(context.Background(), authInfo)
	if err != nil {
		log.Errorf("Failed to update password: %v", err)
		c.JSON(http.StatusInternalServerError, api.ResponseError{
			Code:    http.StatusInternalServerError,
			Message: "Internal Server Error",
			Error:   "Failed to update password",
		})
		return
	}

	// 删除所有会话（强制重新登录）
	session.DisableAll()

	log.Infof("Password changed successfully for user %s from %s", username.(string), c.ClientIP())

	c.JSON(http.StatusOK, api.ResponseSuccess{
		Code:    http.StatusOK,
		Message: "Password changed successfully. Please login again.",
	})
}

// getUserInfo 获取当前用户信息
// @Summary 获取用户信息
// @Description 获取当前登录用户的详细信息
// @Tags 认证
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} api.ResponseSuccess{data=auth.Data} "获取成功"
// @Failure 401 {object} api.ResponseError "未授权"
// @Failure 500 {object} api.ResponseError "服务器内部错误"
// @Router /api/v1/auth/user [get]
func (h *authHandler) getUserInfo(c *gin.Context) {
	_, exists := c.Get("username")
	if !exists {
		c.JSON(http.StatusUnauthorized, api.ResponseError{
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
		c.JSON(http.StatusInternalServerError, api.ResponseError{
			Code:    http.StatusInternalServerError,
			Message: "Internal Server Error",
			Error:   "Failed to get user information",
		})
		return
	}

	if authInfo == nil {
		c.JSON(http.StatusNotFound, api.ResponseError{
			Code:    http.StatusNotFound,
			Message: "Not Found",
			Error:   "User not found",
		})
		return
	}

	userInfo := auth.Data{
		UserName:  authInfo.UserName,
		CreatedAt: authInfo.CreatedAt,
		UpdatedAt: authInfo.UpdatedAt,
	}

	c.JSON(http.StatusOK, api.ResponseSuccess{
		Code:    http.StatusOK,
		Message: "User information retrieved successfully",
		Data:    userInfo,
	})
}

// getSessions 获取当前用户的所有会话
// @Summary 获取用户会话列表
// @Description 获取当前用户的所有活跃会话信息
// @Tags 认证
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} api.ResponseSuccess{data=auth.SessionListResponse} "获取成功"
// @Failure 401 {object} api.ResponseError "未授权"
// @Failure 500 {object} api.ResponseError "服务器内部错误"
// @Router /api/v1/auth/sessions [get]
func (h *authHandler) getSessions(c *gin.Context) {
	username, exists := c.Get("username")
	if !exists {
		c.JSON(http.StatusUnauthorized, api.ResponseError{
			Code:    http.StatusUnauthorized,
			Message: "Unauthorized",
			Error:   "User not found in context",
		})
		return
	}

	// 获取所有活跃会话
	sessions := session.GetActive()

	response := auth.SessionListResponse{
		Sessions: *sessions,
		Total:    len(*sessions),
	}

	log.Debugf("Retrieved %d sessions for user %s", len(*sessions), username)

	c.JSON(http.StatusOK, api.ResponseSuccess{
		Code:    http.StatusOK,
		Message: "Sessions retrieved successfully",
		Data:    response,
	})
}

// deleteSession 删除指定会话
// @Summary 删除会话
// @Description 删除指定ID的会话，使其失效
// @Tags 认证
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "会话ID"
// @Success 200 {object} api.ResponseSuccess "删除成功"
// @Failure 400 {object} api.ResponseError "请求参数错误"
// @Failure 401 {object} api.ResponseError "未授权"
// @Failure 404 {object} api.ResponseError "会话不存在"
// @Failure 500 {object} api.ResponseError "服务器内部错误"
// @Router /api/v1/auth/sessions/{id} [delete]
func (h *authHandler) deleteSession(c *gin.Context) {
	sessionIDStr := c.Param("id")
	if sessionIDStr == "" {
		c.JSON(http.StatusBadRequest, api.ResponseError{
			Code:    http.StatusBadRequest,
			Message: "Bad Request",
			Error:   "Session ID is required",
		})
		return
	}

	sessionID := uint8(0)
	if _, err := fmt.Sscanf(sessionIDStr, "%d", &sessionID); err != nil {
		c.JSON(http.StatusBadRequest, api.ResponseError{
			Code:    http.StatusBadRequest,
			Message: "Bad Request",
			Error:   "Invalid session ID format",
		})
		return
	}

	username, exists := c.Get("username")
	if !exists {
		c.JSON(http.StatusUnauthorized, api.ResponseError{
			Code:    http.StatusUnauthorized,
			Message: "Unauthorized",
			Error:   "User not found in context",
		})
		return
	}

	// 检查会话是否存在
	_, err := session.Get(sessionID)
	if err != nil {
		if err == session.ErrSessionNotFound {
			c.JSON(http.StatusNotFound, api.ResponseError{
				Code:    http.StatusNotFound,
				Message: "Not Found",
				Error:   "Session not found",
			})
		} else {
			log.Errorf("Failed to get session: %v", err)
			c.JSON(http.StatusInternalServerError, api.ResponseError{
				Code:    http.StatusInternalServerError,
				Message: "Internal Server Error",
				Error:   "Failed to get session",
			})
		}
		return
	}

	// 禁用会话
	err = session.Disable(sessionID)
	if err != nil {
		log.Errorf("Failed to disable session: %v", err)
		c.JSON(http.StatusInternalServerError, api.ResponseError{
			Code:    http.StatusInternalServerError,
			Message: "Internal Server Error",
			Error:   "Failed to delete session",
		})
		return
	}

	log.Infof("Session %d deleted by user %s from %s", sessionID, username.(string), c.ClientIP())

	c.JSON(http.StatusOK, api.ResponseSuccess{
		Code:    http.StatusOK,
		Message: "Session deleted successfully",
	})
}

// updateUsername 修改用户名
// @Summary 修改用户名
// @Description 修改当前用户的用户名
// @Tags 认证
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body auth.UpdateUserInfoRequest true "修改用户名请求"
// @Success 200 {object} api.ResponseSuccess "用户名修改成功"
// @Failure 400 {object} api.ResponseError "请求参数错误"
// @Failure 401 {object} api.ResponseError "未授权"
// @Failure 409 {object} api.ResponseError "用户名已存在"
// @Failure 500 {object} api.ResponseError "服务器内部错误"
// @Router /api/v1/auth/user/name [post]
func (h *authHandler) updateUsername(c *gin.Context) {
	var req auth.UpdateUserInfoRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, api.ResponseError{
			Code:    http.StatusBadRequest,
			Message: "Bad Request",
			Error:   "Invalid request format: " + err.Error(),
		})
		return
	}

	currentUsername, exists := c.Get("username")
	if !exists {
		c.JSON(http.StatusUnauthorized, api.ResponseError{
			Code:    http.StatusUnauthorized,
			Message: "Unauthorized",
			Error:   "User not found in context",
		})
		return
	}

	// 检查新用户名是否与当前用户名相同
	if req.Username == currentUsername.(string) {
		c.JSON(http.StatusBadRequest, api.ResponseError{
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
		c.JSON(http.StatusInternalServerError, api.ResponseError{
			Code:    http.StatusInternalServerError,
			Message: "Internal Server Error",
			Error:   "Failed to update username",
		})
		return
	}

	log.Infof("Username changed successfully from %s to %s from %s", currentUsername.(string), req.Username, c.ClientIP())

	c.JSON(http.StatusOK, api.ResponseSuccess{
		Code:    http.StatusOK,
		Message: "Username updated successfully.",
	})
}
