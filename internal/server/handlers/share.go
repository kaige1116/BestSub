package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/bestruirui/bestsub/internal/database/op"
	"github.com/bestruirui/bestsub/internal/models/share"
	"github.com/bestruirui/bestsub/internal/server/middleware"
	"github.com/bestruirui/bestsub/internal/server/resp"
	"github.com/bestruirui/bestsub/internal/server/router"
	"github.com/bestruirui/bestsub/internal/utils/log"
	"github.com/gin-gonic/gin"
)

func init() {
	router.NewGroupRouter("/api/v1/share").
		Use(middleware.Auth()).
		AddRoute(
			router.NewRoute("", router.POST).
				Handle(createShare),
		).
		AddRoute(
			router.NewRoute("", router.GET).
				Handle(getShare),
		).
		AddRoute(
			router.NewRoute("", router.PUT).
				Handle(updateShare),
		).
		AddRoute(
			router.NewRoute("/:id", router.DELETE).
				Handle(deleteShare),
		)
	router.NewGroupRouter("/api/v1/share").
		AddRoute(
			router.NewRoute("/:token", router.GET).
				Handle(getShareContent),
		)
}

// @Summary 创建分享链接
// @Description 创建分享链接
// @Tags 分享
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param data body share.CreateRequest true "分享数据"
// @Success 200 {object} resp.SuccessStruct{data=share.Response} "创建成功"
// @Failure 401 {object} resp.ErrorStruct "未授权"
// @Failure 500 {object} resp.ErrorStruct "服务器内部错误"
// @Router /api/v1/share [post]
func createShare(c *gin.Context) {
	var req share.CreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Errorf("createShare: %v", err)
		resp.ErrorBadRequest(c)
		return
	}
	configBytes, err := json.Marshal(req.Config)
	if err != nil {
		resp.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	configStr := string(configBytes)
	data := share.Data{
		Enable: req.Enable,
		Name:   req.Name,
		Token:  req.Token,
		Config: configStr,
	}
	if err := op.CreateShare(c.Request.Context(), &data); err != nil {
		resp.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	resp.Success(c, share.Response{
		ID:          data.ID,
		Name:        data.Name,
		Token:       data.Token,
		Enable:      data.Enable,
		AccessCount: data.AccessCount,
		Config:      req.Config,
	})
}

// @Summary 获取分享链接
// @Description 获取分享链接
// @Tags 分享
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} resp.SuccessStruct{data=[]share.Response} "获取成功"
// @Failure 401 {object} resp.ErrorStruct "未授权"
// @Failure 500 {object} resp.ErrorStruct "服务器内部错误"
// @Router /api/v1/share [get]
func getShare(c *gin.Context) {
	shares, err := op.GetShareList(c.Request.Context())
	if err != nil {
		resp.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	var result = make([]share.Response, 0, len(shares))
	for _, v := range shares {
		var config share.Config
		if err := json.Unmarshal([]byte(v.Config), &config); err != nil {
			resp.Error(c, http.StatusInternalServerError, err.Error())
			return
		}
		result = append(result, share.Response{
			ID:          v.ID,
			Name:        v.Name,
			Token:       v.Token,
			Enable:      v.Enable,
			AccessCount: v.AccessCount,
			Config:      config,
		})
	}
	resp.Success(c, result)
}

// @Summary 更新分享链接
// @Description 更新分享链接
// @Tags 分享
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param data body share.UpdateRequest true "分享数据"
// @Success 200 {object} resp.SuccessStruct{data=share.Response} "更新成功"
// @Failure 401 {object} resp.ErrorStruct "未授权"
// @Failure 500 {object} resp.ErrorStruct "服务器内部错误"
// @Router /api/v1/share [put]
func updateShare(c *gin.Context) {
	var req share.UpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		resp.ErrorBadRequest(c)
		return
	}
	configBytes, err := json.Marshal(req.Config)
	if err != nil {
		resp.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	configStr := string(configBytes)
	data := share.Data{
		ID:     req.ID,
		Enable: req.Enable,
		Name:   req.Name,
		Token:  req.Token,
		Config: configStr,
	}
	if err := op.UpdateShare(c.Request.Context(), &data); err != nil {
		resp.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	resp.Success(c, share.Response{
		ID:          data.ID,
		Name:        data.Name,
		Token:       data.Token,
		Enable:      data.Enable,
		AccessCount: data.AccessCount,
		Config:      req.Config,
	})
}

// @Summary 删除分享链接
// @Description 删除分享链接
// @Tags 分享
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "分享ID"
// @Success 200 {object} resp.SuccessStruct "删除成功"
// @Failure 401 {object} resp.ErrorStruct "未授权"
// @Failure 500 {object} resp.ErrorStruct "服务器内部错误"
// @Router /api/v1/share/{id} [delete]
func deleteShare(c *gin.Context) {
	id := c.Param("id")
	idUint, err := strconv.ParseUint(id, 10, 16)
	if err != nil {
		resp.ErrorBadRequest(c)
		return
	}
	if err := op.DeleteShare(c.Request.Context(), uint16(idUint)); err != nil {
		resp.ErrorBadRequest(c)
		return
	}
	resp.Success(c, nil)
}

// @Summary 获取订阅内容
// @Description 获取订阅内容
// @Tags 分享
// @Accept json
// @Produce json
// @Param token path string true "分享token"
// @Success 200 {object} resp.SuccessStruct "获取成功"
// @Failure 500 {object} resp.ErrorStruct "服务器内部错误"
// @Router /api/v1/share/{token} [get]
func getShareContent(c *gin.Context) {
	token := c.Param("token")
	shareData, err := op.GetShareByToken(c.Request.Context(), token)
	if err != nil {
		resp.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	if !shareData.Enable {
		resp.Error(c, http.StatusInternalServerError, "share not enable")
		return
	}
	var config share.Config
	if err := json.Unmarshal([]byte(shareData.Config), &config); err != nil {
		resp.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	if config.Expires.Before(time.Now()) {
		resp.Error(c, http.StatusInternalServerError, "share expired")
		return
	}
	if config.MaxAccessCount > 0 && config.MaxAccessCount <= shareData.AccessCount {
		resp.Error(c, http.StatusInternalServerError, "share access count exceeded")
		return
	}
	op.UpdateShareAccessCount(c.Request.Context(), shareData.ID)
	// TODO: 获取订阅内容
	resp.Success(c, config)
}
