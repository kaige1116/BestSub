package handlers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/bestruirui/bestsub/internal/config"
	"github.com/bestruirui/bestsub/internal/database/op"
	authModel "github.com/bestruirui/bestsub/internal/models/auth"
	notifyModel "github.com/bestruirui/bestsub/internal/models/notify"
	"github.com/bestruirui/bestsub/internal/modules/notify"
	"github.com/bestruirui/bestsub/internal/server/auth"
	"github.com/bestruirui/bestsub/internal/server/middleware"
	"github.com/bestruirui/bestsub/internal/server/resp"
	"github.com/bestruirui/bestsub/internal/server/router"
	"github.com/bestruirui/bestsub/internal/utils"
	"github.com/bestruirui/bestsub/internal/utils/log"
	"github.com/cespare/xxhash/v2"
	"github.com/gin-gonic/gin"
)

func init() {

	router.NewGroupRouter("/api/v1/auth").
		AddRoute(
			router.NewRoute("/login", router.POST).
				Handle(login),
		).
		AddRoute(
			router.NewRoute("/refresh", router.POST).
				Handle(refreshToken),
		)

	router.NewGroupRouter("/api/v1/auth").
		Use(middleware.Auth()).
		AddRoute(
			router.NewRoute("/logout", router.POST).
				Handle(logout),
		).
		AddRoute(
			router.NewRoute("/user/password", router.POST).
				Handle(changePassword),
		).
		AddRoute(
			router.NewRoute("/user/name", router.POST).
				Handle(updateUsername),
		).
		AddRoute(
			router.NewRoute("/user", router.GET).
				Handle(getUserInfo),
		).
		AddRoute(
			router.NewRoute("/sessions", router.GET).
				Handle(getSessions),
		).
		AddRoute(
			router.NewRoute("/sessions/:id", router.DELETE).
				Handle(deleteSession),
		)
}

// login 用户登录
// @Summary 用户登录
// @Description 用户登录接口，验证用户名和密码，返回JWT令牌
// @Tags 认证
// @Accept json
// @Produce json
// @Param request body authModel.LoginRequest true "登录请求"
// @Success 200 {object} resp.SuccessStruct{data=authModel.LoginResponse} "登录成功"
// @Failure 400 {object} resp.ErrorStruct "请求参数错误"
// @Failure 401 {object} resp.ErrorStruct "用户名或密码错误"
// @Failure 500 {object} resp.ErrorStruct "服务器内部错误"
// @Router /api/v1/auth/login [post]
func login(c *gin.Context) {
	var req authModel.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		resp.ErrorBadRequest(c)
		return
	}

	err := op.AuthVerify(req.Username, req.Password)
	if err != nil {
		log.Warnf("Login failed for user %s: %v from %s", req.Username, err, c.ClientIP())
		go notify.SendSystemNotify(notifyModel.TypeLoginFailed, "登录失败", authModel.LoginNotify{
			Username:  req.Username,
			IP:        c.ClientIP(),
			Time:      time.Now().Format("2006-01-02 15:04:05"),
			Msg:       "登录失败，用户名或密码错误",
			UserAgent: c.GetHeader("User-Agent"),
		})
		resp.Error(c, http.StatusUnauthorized, "username or password error")
		return
	}

	sessionID, tempSess := auth.GetOneSession()
	if tempSess == nil {
		log.Warnf("No unused session found from %s", c.ClientIP())
		go notify.SendSystemNotify(notifyModel.TypeLoginFailed, "登录失败", authModel.LoginNotify{
			Username:  req.Username,
			IP:        c.ClientIP(),
			Time:      time.Now().Format("2006-01-02 15:04:05"),
			Msg:       "登录失败，没有找到空闲的会话",
			UserAgent: c.GetHeader("User-Agent"),
		})
		resp.Error(c, http.StatusUnauthorized, "please logout other devices first")
		return
	}

	tokenPair, err := auth.GenerateTokenPair(sessionID, req.Username, config.Base().JWT.Secret)
	if err != nil {
		log.Errorf("Failed to generate token pair: %v from %s", err, c.ClientIP())
		auth.DisableSession(sessionID)
		go notify.SendSystemNotify(notifyModel.TypeLoginFailed, "登录失败", authModel.LoginNotify{
			Username:  req.Username,
			IP:        c.ClientIP(),
			Time:      time.Now().Format("2006-01-02 15:04:05"),
			Msg:       "登录失败，生成令牌失败",
			UserAgent: c.GetHeader("User-Agent"),
		})
		resp.Error(c, http.StatusInternalServerError, "failed to generate token pair")
		return
	}

	now := uint32(time.Now().Unix())

	tempSess.IsActive = true
	tempSess.ClientIP = utils.IPToUint32(c.ClientIP())
	tempSess.UserAgent = c.GetHeader("User-Agent")
	tempSess.CreatedAt = now
	tempSess.LastAccessAt = now
	tempSess.ExpiresAt = uint32(tokenPair.RefreshExpiresAt.Unix())
	tempSess.HashRToken = xxhash.Sum64String(tokenPair.RefreshToken)
	tempSess.HashAToken = xxhash.Sum64String(tokenPair.AccessToken)
	log.Infof("User %s logged in successfully from %s", req.Username, c.ClientIP())
	go notify.SendSystemNotify(notifyModel.TypeLoginSuccess, "登录成功", authModel.LoginNotify{
		Username:  req.Username,
		IP:        c.ClientIP(),
		Time:      time.Now().Format("2006-01-02 15:04:05"),
		Msg:       "登录成功",
		UserAgent: c.GetHeader("User-Agent"),
	})

	resp.Success(c, tokenPair)
}

// refreshToken 刷新令牌
// @Summary 刷新访问令牌
// @Description 使用刷新令牌获取新的访问令牌
// @Tags 认证
// @Accept json
// @Produce json
// @Param request body authModel.RefreshTokenRequest true "刷新令牌请求"
// @Success 200 {object} resp.SuccessStruct{data=authModel.LoginResponse} "刷新成功"
// @Failure 400 {object} resp.ErrorStruct "请求参数错误"
// @Failure 401 {object} resp.ErrorStruct "刷新令牌无效"
// @Failure 500 {object} resp.ErrorStruct "服务器内部错误"
// @Router /api/v1/auth/refresh [post]
func refreshToken(c *gin.Context) {
	var req authModel.RefreshTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		resp.ErrorBadRequest(c)
		return
	}

	claims, err := auth.ValidateToken(req.RefreshToken, config.Base().JWT.Secret)
	if err != nil {
		log.Warnf("Refresh token validation failed: %v", err)
		resp.Error(c, http.StatusUnauthorized, "refresh token validation failed")
		return
	}

	sess, err := auth.GetSession(claims.SessionID)
	if err != nil {
		log.Warnf("Failed to get session by ID: %v", err)
		resp.Error(c, http.StatusUnauthorized, "session not found")
		return
	}

	if !sess.IsActive {
		log.Warnf("Session %d is not active", claims.SessionID)
		resp.Error(c, http.StatusUnauthorized, "session is not active")
		return
	}

	if sess.HashRToken != xxhash.Sum64String(req.RefreshToken) {
		log.Warnf("Refresh token hash mismatch: session=%d, request=%d", sess.HashRToken, xxhash.Sum64String(req.RefreshToken))
		resp.Error(c, http.StatusUnauthorized, "refresh token hash mismatch")
		return
	}

	clientIP := utils.IPToUint32(c.ClientIP())

	if sess.ClientIP != clientIP {
		log.Warnf("Client IP mismatch during token refresh: session=%s, request=%s", utils.Uint32ToIP(sess.ClientIP), c.ClientIP())
		resp.Error(c, http.StatusUnauthorized, "client IP mismatch")
		return
	}

	if sess.UserAgent != c.GetHeader("User-Agent") {
		log.Warnf("User agent mismatch during token refresh: session=%s, request=%s",
			sess.UserAgent, c.GetHeader("User-Agent"))
		resp.Error(c, http.StatusUnauthorized, "user agent mismatch")
		return
	}

	newTokenPair, err := auth.GenerateTokenPair(claims.SessionID, claims.Username, config.Base().JWT.Secret)
	if err != nil {
		log.Errorf("Failed to generate new token pair: %v", err)
		resp.Error(c, http.StatusInternalServerError, "failed to generate new token pair")
		return
	}

	sess.ExpiresAt = uint32(newTokenPair.RefreshExpiresAt.Unix())
	sess.LastAccessAt = uint32(time.Now().Unix())
	sess.HashRToken = xxhash.Sum64String(newTokenPair.RefreshToken)
	sess.HashAToken = xxhash.Sum64String(newTokenPair.AccessToken)

	log.Infof("Token refreshed for session %d from %s", claims.SessionID, c.ClientIP())

	resp.Success(c, newTokenPair)
}

// logout 用户登出
// @Summary 用户登出
// @Description 用户登出接口，使当前会话失效
// @Tags 认证
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} resp.SuccessStruct "登出成功"
// @Failure 401 {object} resp.ErrorStruct "未授权"
// @Failure 500 {object} resp.ErrorStruct "服务器内部错误"
// @Router /api/v1/auth/logout [post]
func logout(c *gin.Context) {
	sessionID, exists := c.Get("session_id")
	if !exists {
		resp.Error(c, http.StatusUnauthorized, "session not found")
		return
	}

	err := auth.DisableSession(sessionID.(uint8))
	if err != nil {
		log.Errorf("Failed to disable session: %v", err)
		resp.Error(c, http.StatusInternalServerError, "failed to disable session")
		return
	}

	log.Infof("User logged out successfully from %s", c.ClientIP())

	resp.Success(c, nil)
}

// changePassword 修改密码
// @Summary 修改密码
// @Description 修改当前用户的密码
// @Tags 认证
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body authModel.ChangePasswordRequest true "修改密码请求"
// @Success 200 {object} resp.SuccessStruct "密码修改成功"
// @Failure 400 {object} resp.ErrorStruct "请求参数错误"
// @Failure 401 {object} resp.ErrorStruct "未授权或旧密码错误"
// @Failure 500 {object} resp.ErrorStruct "服务器内部错误"
// @Router /api/v1/auth/user/password [post]
func changePassword(c *gin.Context) {
	var req authModel.ChangePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		resp.ErrorBadRequest(c)
		return
	}

	err := op.AuthVerify(req.Username, req.OldPassword)
	if err != nil {
		log.Warnf("Change password failed for user %s: old password verification failed from %s", req.Username, c.ClientIP())
		resp.Error(c, http.StatusUnauthorized, "old password verification failed")
		return
	}

	err = op.AuthUpdatePassWord(req.NewPassword)
	if err != nil {
		log.Errorf("Failed to update password: %v", err)
		resp.Error(c, http.StatusInternalServerError, "failed to update password")
		return
	}

	auth.DisableAllSession()

	log.Infof("Password changed successfully for user %s from %s", req.Username, c.ClientIP())

	resp.Success(c, nil)
}

// getUserInfo 获取当前用户信息
// @Summary 获取用户信息
// @Description 获取当前登录用户的详细信息
// @Tags 认证
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} resp.SuccessStruct{data=authModel.Data} "获取成功"
// @Failure 401 {object} resp.ErrorStruct "未授权"
// @Failure 500 {object} resp.ErrorStruct "服务器内部错误"
// @Router /api/v1/auth/user [get]
func getUserInfo(c *gin.Context) {
	authInfo, err := op.AuthGet()
	if err != nil {
		log.Errorf("Failed to get auth info from %s: %v", c.ClientIP(), err)
		resp.Error(c, http.StatusInternalServerError, "failed to get auth info")
		return
	}
	resp.Success(c, authInfo)
}

// getSessions 获取当前用户的所有会话
// @Summary 获取用户会话列表
// @Description 获取当前用户的所有活跃会话信息
// @Tags 认证
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} resp.SuccessStruct{data=authModel.SessionListResponse} "获取成功"
// @Failure 401 {object} resp.ErrorStruct "未授权"
// @Failure 500 {object} resp.ErrorStruct "服务器内部错误"
// @Router /api/v1/auth/sessions [get]
func getSessions(c *gin.Context) {
	sessions := auth.GetAllSession()
	response := authModel.SessionListResponse{
		Sessions: *sessions,
		Total:    uint8(len(*sessions)),
	}
	resp.Success(c, response)
}

// deleteSession 删除指定会话
// @Summary 删除会话
// @Description 删除指定ID的会话，使其失效
// @Tags 认证
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "会话ID"
// @Success 200 {object} resp.SuccessStruct "删除成功"
// @Failure 400 {object} resp.ErrorStruct "请求参数错误"
// @Failure 401 {object} resp.ErrorStruct "未授权"
// @Failure 404 {object} resp.ErrorStruct "会话不存在"
// @Failure 500 {object} resp.ErrorStruct "服务器内部错误"
// @Router /api/v1/auth/sessions/{id} [delete]
func deleteSession(c *gin.Context) {
	sessionIDStr := c.Param("id")
	if sessionIDStr == "" {
		resp.ErrorBadRequest(c)
		return
	}

	sessionID, err := strconv.ParseUint(sessionIDStr, 10, 8)
	if err != nil {
		resp.ErrorBadRequest(c)
		return
	}

	_, err = auth.GetSession(uint8(sessionID))
	if err != nil {
		if err == auth.ErrSessionNotFound {
			resp.Error(c, http.StatusNotFound, "session not found")
		} else {
			log.Errorf("Failed to get session: %v", err)
			resp.Error(c, http.StatusInternalServerError, "failed to get session")
		}
		return
	}

	err = auth.DisableSession(uint8(sessionID))
	if err != nil {
		log.Errorf("Failed to disable session: %v", err)
		resp.Error(c, http.StatusInternalServerError, "failed to disable session")
		return
	}

	log.Infof("Session %d disabled by user from %s", sessionID, c.ClientIP())

	resp.Success(c, nil)
}

// updateUsername 修改用户名
// @Summary 修改用户名
// @Description 修改当前用户的用户名
// @Tags 认证
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body authModel.UpdateUserInfoRequest true "修改用户名请求"
// @Success 200 {object} resp.SuccessStruct "用户名修改成功"
// @Failure 400 {object} resp.ErrorStruct "请求参数错误"
// @Failure 401 {object} resp.ErrorStruct "未授权"
// @Failure 409 {object} resp.ErrorStruct "用户名已存在"
// @Failure 500 {object} resp.ErrorStruct "服务器内部错误"
// @Router /api/v1/auth/user/name [post]
func updateUsername(c *gin.Context) {
	var req authModel.UpdateUserInfoRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		resp.ErrorBadRequest(c)
		return
	}

	authInfo, err := op.AuthGet()
	if err != nil {
		log.Errorf("Failed to get auth info from %s: %v", c.ClientIP(), err)
		resp.Error(c, http.StatusInternalServerError, "failed to get auth info")
		return
	}

	if authInfo.UserName == req.Username {
		resp.Error(c, http.StatusBadRequest, "new username cannot be the same as current username")
		return
	}

	if err := op.AuthUpdateName(req.Username); err != nil {
		resp.Error(c, http.StatusInternalServerError, "failed to update username")
		return
	}

	log.Infof("Username changed successfully from %s to %s from %s", authInfo.UserName, req.Username, c.ClientIP())

	resp.Success(c, nil)
}
