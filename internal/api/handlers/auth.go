package handlers

import (
	"net/http"
	"strconv"

	"github.com/bestruirui/bestsub/internal/api/middleware"
	"github.com/bestruirui/bestsub/internal/api/router"
	"github.com/bestruirui/bestsub/internal/config"
	"github.com/bestruirui/bestsub/internal/core/session"
	"github.com/bestruirui/bestsub/internal/database"
	"github.com/bestruirui/bestsub/internal/models/api"
	"github.com/bestruirui/bestsub/internal/models/auth"
	"github.com/bestruirui/bestsub/internal/utils"
	"github.com/bestruirui/bestsub/internal/utils/local"
	"github.com/bestruirui/bestsub/internal/utils/log"
	"github.com/cespare/xxhash/v2"
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

	err := database.AuthVerify(req.Username, req.Password)
	if err != nil {
		log.Warnf("Login failed for user %s: %v from %s", req.Username, err, c.ClientIP())
		c.JSON(http.StatusUnauthorized, api.ResponseError{
			Code:    http.StatusUnauthorized,
			Message: "Unauthorized",
			Error:   "Invalid username or password",
		})
		return
	}

	sessionID, tempSess := session.GetOne()
	if tempSess == nil {
		log.Warnf("No unused session found from %s", c.ClientIP())
		c.JSON(http.StatusInternalServerError, api.ResponseError{
			Code:    http.StatusInternalServerError,
			Message: "Internal Server Error",
			Error:   "Too many devices, please log out of other devices",
		})
		return
	}

	tokenPair, err := middleware.GenerateTokenPair(sessionID, req.Username, config.Base().JWT.Secret)
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

	tempSess.IsActive = true
	tempSess.ClientIP = utils.IPToUint32(c.ClientIP())
	tempSess.UserAgent = c.GetHeader("User-Agent")
	tempSess.CreatedAt = now
	tempSess.LastAccessAt = now
	tempSess.ExpiresAt = uint32(tokenPair.AccessExpiresAt.Unix())
	tempSess.HashRToken = xxhash.Sum64String(tokenPair.RefreshToken)
	tempSess.HashAToken = xxhash.Sum64String(tokenPair.AccessToken)

	log.Infof("User %s logged in successfully from %s", req.Username, c.ClientIP())

	c.JSON(http.StatusOK, api.ResponseSuccess{
		Code:    http.StatusOK,
		Message: "Login successful",
		Data:    tokenPair,
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

	claims, err := middleware.ValidateToken(req.RefreshToken, config.Base().JWT.Secret)
	if err != nil {
		log.Warnf("Refresh token validation failed: %v", err)
		c.JSON(http.StatusUnauthorized, api.ResponseError{
			Code:    http.StatusUnauthorized,
			Message: "Unauthorized",
			Error:   "Invalid refresh token",
		})
		return
	}

	sess, err := session.Get(claims.SessionID)
	if err != nil {
		log.Warnf("Failed to get session by ID: %v", err)
		c.JSON(http.StatusUnauthorized, api.ResponseError{
			Code:    http.StatusUnauthorized,
			Message: "Unauthorized",
			Error:   "Session not found or expired",
		})
		return
	}

	if !sess.IsActive {
		log.Warnf("Session %d is not active", claims.SessionID)
		c.JSON(http.StatusUnauthorized, api.ResponseError{
			Code:    http.StatusUnauthorized,
			Message: "Unauthorized",
			Error:   "Session is not active",
		})
		return
	}

	if sess.HashRToken != xxhash.Sum64String(req.RefreshToken) {
		log.Warnf("Refresh token hash mismatch: session=%d, request=%d", sess.HashRToken, xxhash.Sum64String(req.RefreshToken))
		c.JSON(http.StatusUnauthorized, api.ResponseError{
			Code:    http.StatusUnauthorized,
			Message: "Unauthorized",
			Error:   "Refresh token hash mismatch",
		})
		return
	}

	clientIP := utils.IPToUint32(c.ClientIP())

	if sess.ClientIP != clientIP {
		log.Warnf("Client IP mismatch during token refresh: session=%s, request=%s", utils.Uint32ToIP(sess.ClientIP), c.ClientIP())
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

	newTokenPair, err := middleware.GenerateTokenPair(claims.SessionID, claims.Username, config.Base().JWT.Secret)
	if err != nil {
		log.Errorf("Failed to generate new token pair: %v", err)
		c.JSON(http.StatusInternalServerError, api.ResponseError{
			Code:    http.StatusInternalServerError,
			Message: "Internal Server Error",
			Error:   "Failed to generate new tokens",
		})
		return
	}

	sess.ExpiresAt = uint32(newTokenPair.AccessExpiresAt.Unix())
	sess.LastAccessAt = uint32(local.Time().Unix())
	sess.HashRToken = xxhash.Sum64String(newTokenPair.RefreshToken)
	sess.HashAToken = xxhash.Sum64String(newTokenPair.AccessToken)

	log.Infof("Token refreshed for session %d from %s", claims.SessionID, c.ClientIP())

	c.JSON(http.StatusOK, api.ResponseSuccess{
		Code:    http.StatusOK,
		Message: "Token refreshed successfully",
		Data:    newTokenPair,
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

	log.Infof("User logged out successfully from %s", c.ClientIP())

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

	err := database.AuthVerify(req.Username, req.OldPassword)
	if err != nil {
		log.Warnf("Change password failed for user %s: old password verification failed from %s", req.Username, c.ClientIP())
		c.JSON(http.StatusUnauthorized, api.ResponseError{
			Code:    http.StatusUnauthorized,
			Message: "Unauthorized",
			Error:   "Old password is incorrect",
		})
		return
	}

	err = database.AuthUpdatePassWord(req.NewPassword)
	if err != nil {
		log.Errorf("Failed to update password: %v", err)
		c.JSON(http.StatusInternalServerError, api.ResponseError{
			Code:    http.StatusInternalServerError,
			Message: "Internal Server Error",
			Error:   "Failed to update password",
		})
		return
	}

	session.DisableAll()

	log.Infof("Password changed successfully for user %s from %s", req.Username, c.ClientIP())

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
	authInfo, err := database.AuthGet()
	if err != nil {
		log.Errorf("Failed to get auth info from %s: %v", c.ClientIP(), err)
		c.JSON(http.StatusInternalServerError, api.ResponseError{
			Code:    http.StatusInternalServerError,
			Message: "Internal Server Error",
			Error:   "Failed to get auth info",
		})
		return
	}
	c.JSON(http.StatusOK, api.ResponseSuccess{
		Code:    http.StatusOK,
		Message: "User information retrieved successfully",
		Data:    authInfo,
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
	sessions := session.GetAll()
	response := auth.SessionListResponse{
		Sessions: *sessions,
		Total:    len(*sessions),
	}
	log.Debugf("Retrieved %d sessions", len(*sessions))
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

	sessionID, err := strconv.ParseUint(sessionIDStr, 10, 8)
	if err != nil {
		c.JSON(http.StatusBadRequest, api.ResponseError{
			Code:    http.StatusBadRequest,
			Message: "Bad Request",
			Error:   "Invalid session ID format",
		})
		return
	}

	_, err = session.Get(uint8(sessionID))
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

	err = session.Disable(uint8(sessionID))
	if err != nil {
		log.Errorf("Failed to disable session: %v", err)
		c.JSON(http.StatusInternalServerError, api.ResponseError{
			Code:    http.StatusInternalServerError,
			Message: "Internal Server Error",
			Error:   "Failed to delete session",
		})
		return
	}

	log.Infof("Session %d disabled by user from %s", sessionID, c.ClientIP())

	c.JSON(http.StatusOK, api.ResponseSuccess{
		Code:    http.StatusOK,
		Message: "Session disabled successfully",
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

	authInfo, err := database.AuthGet()
	if err != nil {
		log.Errorf("Failed to get auth info from %s: %v", c.ClientIP(), err)
		c.JSON(http.StatusInternalServerError, api.ResponseError{
			Code:    http.StatusInternalServerError,
			Message: "Internal Server Error",
			Error:   "Failed to get auth info",
		})
		return
	}

	if authInfo.UserName == req.Username {
		c.JSON(http.StatusBadRequest, api.ResponseError{
			Code:    http.StatusBadRequest,
			Message: "Bad Request",
			Error:   "New username cannot be the same as current username",
		})
		return
	}

	if err := database.AuthUpdateName(req.Username); err != nil {
		c.JSON(http.StatusInternalServerError, api.ResponseError{
			Code:    http.StatusInternalServerError,
			Message: "Internal Server Error",
			Error:   "Failed to update username",
		})
		return
	}

	log.Infof("Username changed successfully from %s to %s from %s", authInfo.UserName, req.Username, c.ClientIP())

	c.JSON(http.StatusOK, api.ResponseSuccess{
		Code:    http.StatusOK,
		Message: "Username updated successfully.",
	})
}
