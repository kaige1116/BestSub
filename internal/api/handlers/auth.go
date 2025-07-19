package handlers

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/bestruirui/bestsub/internal/api/common"
	"github.com/bestruirui/bestsub/internal/api/middleware"
	"github.com/bestruirui/bestsub/internal/api/router"
	"github.com/bestruirui/bestsub/internal/config"
	"github.com/bestruirui/bestsub/internal/database/op"
	"github.com/bestruirui/bestsub/internal/models/auth"
	"github.com/bestruirui/bestsub/internal/utils"
	"github.com/bestruirui/bestsub/internal/utils/local"
	"github.com/bestruirui/bestsub/internal/utils/log"
	"github.com/cespare/xxhash/v2"
	"github.com/gin-gonic/gin"
)

// init 函数用于自动注册路由
func init() {

	// 公开的认证路由（无需认证）
	router.NewGroupRouter("/api/v1/auth").
		AddRoute(
			router.NewRoute("/login", router.POST).
				Handle(login).
				WithDescription("User login"),
		).
		AddRoute(
			router.NewRoute("/refresh", router.POST).
				Handle(refreshToken).
				WithDescription("Refresh access token"),
		)

	// 需要认证的路由
	router.NewGroupRouter("/api/v1/auth").
		Use(middleware.Auth()).
		AddRoute(
			router.NewRoute("/logout", router.POST).
				Handle(logout).
				WithDescription("User logout"),
		).
		AddRoute(
			router.NewRoute("/user/password", router.POST).
				Handle(changePassword).
				WithDescription("Change user password"),
		).
		AddRoute(
			router.NewRoute("/user/name", router.POST).
				Handle(updateUsername).
				WithDescription("Update username"),
		).
		AddRoute(
			router.NewRoute("/user", router.GET).
				Handle(getUserInfo).
				WithDescription("Get user information"),
		).
		AddRoute(
			router.NewRoute("/sessions", router.GET).
				Handle(getSessions).
				WithDescription("Get user sessions"),
		).
		AddRoute(
			router.NewRoute("/sessions/:id", router.DELETE).
				Handle(deleteSession).
				WithDescription("Delete session"),
		)
}

// login 用户登录
// @Summary 用户登录
// @Description 用户登录接口，验证用户名和密码，返回JWT令牌
// @Tags 认证
// @Accept json
// @Produce json
// @Param request body auth.LoginRequest true "登录请求"
// @Success 200 {object} common.ResponseSuccessStruct{data=auth.LoginResponse} "登录成功"
// @Failure 400 {object} common.ResponseErrorStruct "请求参数错误"
// @Failure 401 {object} common.ResponseErrorStruct "用户名或密码错误"
// @Failure 500 {object} common.ResponseErrorStruct "服务器内部错误"
// @Router /api/v1/auth/login [post]
func login(c *gin.Context) {
	var req auth.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		common.ResponseError(c, http.StatusBadRequest, err)
		return
	}

	err := op.AuthVerify(req.Username, req.Password)
	if err != nil {
		log.Warnf("Login failed for user %s: %v from %s", req.Username, err, c.ClientIP())
		common.ResponseError(c, http.StatusUnauthorized, err)
		return
	}

	sessionID, tempSess := common.GetOneSession()
	if tempSess == nil {
		log.Warnf("No unused session found from %s", c.ClientIP())
		common.ResponseError(c, http.StatusUnauthorized, err)
		return
	}

	tokenPair, err := common.GenerateTokenPair(sessionID, req.Username, config.Base().JWT.Secret)
	if err != nil {
		log.Errorf("Failed to generate token pair: %v from %s", err, c.ClientIP())
		common.DisableSession(sessionID)
		common.ResponseError(c, http.StatusInternalServerError, err)
		return
	}

	now := uint32(local.Time().Unix())

	tempSess.IsActive = true
	tempSess.ClientIP = utils.IPToUint32(c.ClientIP())
	tempSess.UserAgent = c.GetHeader("User-Agent")
	tempSess.CreatedAt = now
	tempSess.LastAccessAt = now
	tempSess.ExpiresAt = uint32(tokenPair.RefreshExpiresAt.Unix())
	tempSess.HashRToken = xxhash.Sum64String(tokenPair.RefreshToken)
	tempSess.HashAToken = xxhash.Sum64String(tokenPair.AccessToken)

	log.Infof("User %s logged in successfully from %s", req.Username, c.ClientIP())

	common.ResponseSuccess(c, tokenPair)
}

// refreshToken 刷新令牌
// @Summary 刷新访问令牌
// @Description 使用刷新令牌获取新的访问令牌
// @Tags 认证
// @Accept json
// @Produce json
// @Param request body auth.RefreshTokenRequest true "刷新令牌请求"
// @Success 200 {object} common.ResponseSuccessStruct{data=auth.LoginResponse} "刷新成功"
// @Failure 400 {object} common.ResponseErrorStruct "请求参数错误"
// @Failure 401 {object} common.ResponseErrorStruct "刷新令牌无效"
// @Failure 500 {object} common.ResponseErrorStruct "服务器内部错误"
// @Router /api/v1/auth/refresh [post]
func refreshToken(c *gin.Context) {
	var req auth.RefreshTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		common.ResponseError(c, http.StatusBadRequest, err)
		return
	}

	claims, err := common.ValidateToken(req.RefreshToken, config.Base().JWT.Secret)
	if err != nil {
		log.Warnf("Refresh token validation failed: %v", err)
		common.ResponseError(c, http.StatusUnauthorized, fmt.Errorf("refresh token validation failed"))
		return
	}

	sess, err := common.GetSession(claims.SessionID)
	if err != nil {
		log.Warnf("Failed to get session by ID: %v", err)
		common.ResponseError(c, http.StatusUnauthorized, err)
		return
	}

	if !sess.IsActive {
		log.Warnf("Session %d is not active", claims.SessionID)
		common.ResponseError(c, http.StatusUnauthorized, fmt.Errorf("session is not active"))
		return
	}

	if sess.HashRToken != xxhash.Sum64String(req.RefreshToken) {
		log.Warnf("Refresh token hash mismatch: session=%d, request=%d", sess.HashRToken, xxhash.Sum64String(req.RefreshToken))
		common.ResponseError(c, http.StatusUnauthorized, err)
		return
	}

	clientIP := utils.IPToUint32(c.ClientIP())

	if sess.ClientIP != clientIP {
		log.Warnf("Client IP mismatch during token refresh: session=%s, request=%s", utils.Uint32ToIP(sess.ClientIP), c.ClientIP())
		common.ResponseError(c, http.StatusUnauthorized, err)
		return
	}

	if sess.UserAgent != c.GetHeader("User-Agent") {
		log.Warnf("User agent mismatch during token refresh: session=%s, request=%s",
			sess.UserAgent, c.GetHeader("User-Agent"))
		common.ResponseError(c, http.StatusUnauthorized, err)
		return
	}

	newTokenPair, err := common.GenerateTokenPair(claims.SessionID, claims.Username, config.Base().JWT.Secret)
	if err != nil {
		log.Errorf("Failed to generate new token pair: %v", err)
		common.ResponseError(c, http.StatusInternalServerError, err)
		return
	}

	sess.ExpiresAt = uint32(newTokenPair.AccessExpiresAt.Unix())
	sess.LastAccessAt = uint32(local.Time().Unix())
	sess.HashRToken = xxhash.Sum64String(newTokenPair.RefreshToken)
	sess.HashAToken = xxhash.Sum64String(newTokenPair.AccessToken)

	log.Infof("Token refreshed for session %d from %s", claims.SessionID, c.ClientIP())

	common.ResponseSuccess(c, newTokenPair)
}

// logout 用户登出
// @Summary 用户登出
// @Description 用户登出接口，使当前会话失效
// @Tags 认证
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} common.ResponseSuccessStruct "登出成功"
// @Failure 401 {object} common.ResponseErrorStruct "未授权"
// @Failure 500 {object} common.ResponseErrorStruct "服务器内部错误"
// @Router /api/v1/auth/logout [post]
func logout(c *gin.Context) {
	sessionID, exists := c.Get("session_id")
	if !exists {
		common.ResponseError(c, http.StatusUnauthorized, errors.New("session not found"))
		return
	}

	err := common.DisableSession(sessionID.(uint8))
	if err != nil {
		log.Errorf("Failed to disable session: %v", err)
		common.ResponseError(c, http.StatusInternalServerError, err)
		return
	}

	log.Infof("User logged out successfully from %s", c.ClientIP())

	common.ResponseSuccess(c, nil)
}

// changePassword 修改密码
// @Summary 修改密码
// @Description 修改当前用户的密码
// @Tags 认证
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body auth.ChangePasswordRequest true "修改密码请求"
// @Success 200 {object} common.ResponseSuccessStruct "密码修改成功"
// @Failure 400 {object} common.ResponseErrorStruct "请求参数错误"
// @Failure 401 {object} common.ResponseErrorStruct "未授权或旧密码错误"
// @Failure 500 {object} common.ResponseErrorStruct "服务器内部错误"
// @Router /api/v1/auth/user/password [post]
func changePassword(c *gin.Context) {
	var req auth.ChangePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		common.ResponseError(c, http.StatusBadRequest, err)
		return
	}

	err := op.AuthVerify(req.Username, req.OldPassword)
	if err != nil {
		log.Warnf("Change password failed for user %s: old password verification failed from %s", req.Username, c.ClientIP())
		common.ResponseError(c, http.StatusUnauthorized, err)
		return
	}

	err = op.AuthUpdatePassWord(req.NewPassword)
	if err != nil {
		log.Errorf("Failed to update password: %v", err)
		common.ResponseError(c, http.StatusInternalServerError, err)
		return
	}

	common.DisableAllSession()

	log.Infof("Password changed successfully for user %s from %s", req.Username, c.ClientIP())

	common.ResponseSuccess(c, nil)
}

// getUserInfo 获取当前用户信息
// @Summary 获取用户信息
// @Description 获取当前登录用户的详细信息
// @Tags 认证
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} common.ResponseSuccessStruct{data=auth.Data} "获取成功"
// @Failure 401 {object} common.ResponseErrorStruct "未授权"
// @Failure 500 {object} common.ResponseErrorStruct "服务器内部错误"
// @Router /api/v1/auth/user [get]
func getUserInfo(c *gin.Context) {
	authInfo, err := op.AuthGet()
	if err != nil {
		log.Errorf("Failed to get auth info from %s: %v", c.ClientIP(), err)
		common.ResponseError(c, http.StatusInternalServerError, err)
		return
	}
	common.ResponseSuccess(c, authInfo)
}

// getSessions 获取当前用户的所有会话
// @Summary 获取用户会话列表
// @Description 获取当前用户的所有活跃会话信息
// @Tags 认证
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} common.ResponseSuccessStruct{data=auth.SessionListResponse} "获取成功"
// @Failure 401 {object} common.ResponseErrorStruct "未授权"
// @Failure 500 {object} common.ResponseErrorStruct "服务器内部错误"
// @Router /api/v1/auth/sessions [get]
func getSessions(c *gin.Context) {
	sessions := common.GetAllSession()
	response := auth.SessionListResponse{
		Sessions: *sessions,
		Total:    uint8(len(*sessions)),
	}
	common.ResponseSuccess(c, response)
}

// deleteSession 删除指定会话
// @Summary 删除会话
// @Description 删除指定ID的会话，使其失效
// @Tags 认证
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "会话ID"
// @Success 200 {object} common.ResponseSuccessStruct "删除成功"
// @Failure 400 {object} common.ResponseErrorStruct "请求参数错误"
// @Failure 401 {object} common.ResponseErrorStruct "未授权"
// @Failure 404 {object} common.ResponseErrorStruct "会话不存在"
// @Failure 500 {object} common.ResponseErrorStruct "服务器内部错误"
// @Router /api/v1/auth/sessions/{id} [delete]
func deleteSession(c *gin.Context) {
	sessionIDStr := c.Param("id")
	if sessionIDStr == "" {
		common.ResponseError(c, http.StatusBadRequest, errors.New("session ID is required"))
		return
	}

	sessionID, err := strconv.ParseUint(sessionIDStr, 10, 8)
	if err != nil {
		common.ResponseError(c, http.StatusBadRequest, err)
		return
	}

	_, err = common.GetSession(uint8(sessionID))
	if err != nil {
		if err == common.ErrSessionNotFound {
			common.ResponseError(c, http.StatusNotFound, err)
		} else {
			log.Errorf("Failed to get session: %v", err)
			common.ResponseError(c, http.StatusInternalServerError, err)
		}
		return
	}

	err = common.DisableSession(uint8(sessionID))
	if err != nil {
		log.Errorf("Failed to disable session: %v", err)
		common.ResponseError(c, http.StatusInternalServerError, err)
		return
	}

	log.Infof("Session %d disabled by user from %s", sessionID, c.ClientIP())

	common.ResponseSuccess(c, nil)
}

// updateUsername 修改用户名
// @Summary 修改用户名
// @Description 修改当前用户的用户名
// @Tags 认证
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body auth.UpdateUserInfoRequest true "修改用户名请求"
// @Success 200 {object} common.ResponseSuccessStruct "用户名修改成功"
// @Failure 400 {object} common.ResponseErrorStruct "请求参数错误"
// @Failure 401 {object} common.ResponseErrorStruct "未授权"
// @Failure 409 {object} common.ResponseErrorStruct "用户名已存在"
// @Failure 500 {object} common.ResponseErrorStruct "服务器内部错误"
// @Router /api/v1/auth/user/name [post]
func updateUsername(c *gin.Context) {
	var req auth.UpdateUserInfoRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		common.ResponseError(c, http.StatusBadRequest, err)
		return
	}

	authInfo, err := op.AuthGet()
	if err != nil {
		log.Errorf("Failed to get auth info from %s: %v", c.ClientIP(), err)
		common.ResponseError(c, http.StatusInternalServerError, err)
		return
	}

	if authInfo.UserName == req.Username {
		common.ResponseError(c, http.StatusBadRequest, errors.New("new username cannot be the same as current username"))
		return
	}

	if err := op.AuthUpdateName(req.Username); err != nil {
		common.ResponseError(c, http.StatusInternalServerError, err)
		return
	}

	log.Infof("Username changed successfully from %s to %s from %s", authInfo.UserName, req.Username, c.ClientIP())

	common.ResponseSuccess(c, nil)
}
