package handlers

import (
	"bytes"
	"fmt"
	"net/http"
	"slices"
	"strconv"

	"github.com/bestruirui/bestsub/internal/database/op"
	notifyModel "github.com/bestruirui/bestsub/internal/models/notify"
	"github.com/bestruirui/bestsub/internal/modules/notify"
	"github.com/bestruirui/bestsub/internal/server/middleware"
	"github.com/bestruirui/bestsub/internal/server/resp"
	"github.com/bestruirui/bestsub/internal/server/router"
	"github.com/bestruirui/bestsub/internal/utils/log"
	"github.com/gin-gonic/gin"
)

func init() {
	router.NewGroupRouter("/api/v1/notify").
		Use(middleware.Auth()).
		AddRoute(
			router.NewRoute("/channel", router.GET).
				Handle(getNotifyChannel),
		).
		AddRoute(
			router.NewRoute("/config", router.GET).
				Handle(getNotifyConfig),
		).
		AddRoute(
			router.NewRoute("/list", router.GET).
				Handle(getNotifyList),
		).
		AddRoute(
			router.NewRoute("/config", router.POST).
				Handle(createNotify),
		).
		AddRoute(
			router.NewRoute("/config", router.PUT).
				Handle(updateNotify),
		).
		AddRoute(
			router.NewRoute("/config", router.DELETE).
				Handle(deleteNotify),
		).
		AddRoute(
			router.NewRoute("/test", router.POST).
				Handle(testNotify),
		).
		AddRoute(
			router.NewRoute("/template", router.GET).
				Handle(getTemplates),
		).
		AddRoute(
			router.NewRoute("/template", router.PUT).
				Handle(updateTemplate),
		)
}

// getNotifyConfig 获取通知渠道
// @Summary 获取通知渠道
// @Tags 通知
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} resp.SuccessStruct{data=map[string][]notify.Desc} "获取成功"
// @Failure 400 {object} resp.ErrorStruct "请求参数错误"
// @Failure 401 {object} resp.ErrorStruct "未授权"
// @Failure 500 {object} resp.ErrorStruct "服务器内部错误"
// @Router /api/v1/notify/channel [get]
func getNotifyChannel(c *gin.Context) {
	channels := make([]string, 0, len(notify.GetNotifyInfoMap()))
	for channel := range notify.GetNotifyInfoMap() {
		channels = append(channels, channel)
	}
	resp.Success(c, channels)
}

// getNotifyConfig 获取通知配置
// @Summary 获取通知配置
// @Tags 通知
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param channel query string false "渠道"
// @Success 200 {object} resp.SuccessStruct{data=map[string][]notify.Desc} "获取成功"
// @Failure 400 {object} resp.ErrorStruct "请求参数错误"
// @Failure 401 {object} resp.ErrorStruct "未授权"
// @Failure 500 {object} resp.ErrorStruct "服务器内部错误"
// @Router /api/v1/notify/config [get]
func getNotifyConfig(c *gin.Context) {
	channel := c.Query("channel")
	if channel == "" {
		resp.Success(c, notify.GetNotifyInfoMap())
	} else {
		resp.Success(c, notify.GetNotifyInfoMap()[channel])
	}
}

// getNotifyList 获取通知列表
// @Summary 获取通知列表
// @Tags 通知
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} resp.SuccessStruct{data=[]notifyModel.Data} "获取成功"
// @Failure 400 {object} resp.ErrorStruct "请求参数错误"
// @Failure 401 {object} resp.ErrorStruct "未授权"
// @Failure 500 {object} resp.ErrorStruct "服务器内部错误"
// @Router /api/v1/notify/list [get]
func getNotifyList(c *gin.Context) {
	notifyList, err := op.GetNotifyList()
	if err != nil {
		resp.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	resp.Success(c, notifyList)
}

// createNotify 创建通知配置
// @Summary 创建通知配置
// @Description 创建单个通知配置
// @Tags 通知
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body notifyModel.CreateRequest true "创建通知配置请求"
// @Success 200 {object} resp.SuccessStruct{data=notifyModel.Data} "创建成功"
// @Failure 400 {object} resp.ErrorStruct "请求参数错误"
// @Failure 401 {object} resp.ErrorStruct "未授权"
// @Failure 500 {object} resp.ErrorStruct "服务器内部错误"
// @Router /api/v1/notify/config [post]
func createNotify(c *gin.Context) {
	var req notifyModel.CreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		resp.ErrorBadRequest(c)
		return
	}
	notifyData := &notifyModel.Data{
		Type:   req.Type,
		Config: req.Config,
	}
	types := notify.GetNotifyTypes()
	if !slices.Contains(types, req.Type) {
		resp.Error(c, http.StatusBadRequest, fmt.Sprintf("通知类型 %s 不存在", req.Type))
		return
	}
	if err := op.CreateNotify(c.Request.Context(), notifyData); err != nil {
		resp.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	log.Infof("Notify config %d created by from %s", notifyData.ID, c.ClientIP())
	resp.Success(c, notifyData)
}

// testNotify 测试通知配置
// @Summary 测试通知配置
// @Description 测试单个通知配置
// @Tags 通知
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body notifyModel.CreateRequest true "测试通知配置请求"
// @Success 200 {object} resp.SuccessStruct{data=notifyModel.Data} "测试成功"
// @Failure 400 {object} resp.ErrorStruct "请求参数错误"
// @Failure 401 {object} resp.ErrorStruct "未授权"
// @Failure 500 {object} resp.ErrorStruct "服务器内部错误"
// @Router /api/v1/notify/test [post]
func testNotify(c *gin.Context) {
	var req notifyModel.CreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		resp.ErrorBadRequest(c)
		return
	}
	types := notify.GetNotifyTypes()
	if !slices.Contains(types, req.Type) {
		resp.Error(c, http.StatusBadRequest, fmt.Sprintf("通知类型 %s 不存在", req.Type))
		return
	}
	notify, err := notify.GetNotify(req.Type, req.Config)
	if err != nil {
		resp.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	err = notify.Init()
	if err != nil {
		resp.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	var buf bytes.Buffer
	buf.WriteString("test")
	err = notify.Send("test", &buf)
	if err != nil {
		resp.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	log.Infof("Notify config %s tested by from %s", req.Type, c.ClientIP())
	resp.Success(c, nil)
}

// updateNotify 更新通知配置
// @Summary 更新通知配置
// @Description 根据请求体中的ID更新通知配置信息
// @Tags 通知
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body notifyModel.Data true "更新通知配置请求"
// @Success 200 {object} resp.SuccessStruct{data=notifyModel.Data} "更新成功"
// @Failure 400 {object} resp.ErrorStruct "请求参数错误"
// @Failure 401 {object} resp.ErrorStruct "未授权"
// @Failure 404 {object} resp.ErrorStruct "通知配置不存在"
// @Failure 500 {object} resp.ErrorStruct "服务器内部错误"
// @Router /api/v1/notify/config [put]
func updateNotify(c *gin.Context) {
	var req notifyModel.Data
	if err := c.ShouldBindJSON(&req); err != nil {
		resp.ErrorBadRequest(c)
		return
	}
	if err := op.UpdateNotify(c.Request.Context(), &req); err != nil {
		log.Errorf("Update notify config %d failed: %v", req.ID, err)
		resp.Error(c, http.StatusInternalServerError, "update notify config failed")
		return
	}
	log.Infof("Notify config %d updated by from %s", req.ID, c.ClientIP())
	resp.Success(c, req)
}

// deleteNotify 删除通知配置
// @Summary 删除通知配置
// @Description 根据ID删除单个通知配置
// @Tags 通知
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id query int true "通知配置ID"
// @Success 200 {object} resp.SuccessStruct "删除成功"
// @Failure 400 {object} resp.ErrorStruct "请求参数错误"
// @Failure 401 {object} resp.ErrorStruct "未授权"
// @Failure 404 {object} resp.ErrorStruct "通知配置不存在"
// @Failure 500 {object} resp.ErrorStruct "服务器内部错误"
// @Router /api/v1/notify/config [delete]
func deleteNotify(c *gin.Context) {
	// 获取查询参数中的ID
	idParam := c.Query("id")
	if idParam == "" {
		resp.ErrorBadRequest(c)
		return
	}

	id, err := strconv.ParseUint(idParam, 10, 16)
	if err != nil {
		resp.ErrorBadRequest(c)
		return
	}
	if err := op.DeleteNotify(c.Request.Context(), uint16(id)); err != nil {
		log.Errorf("Delete notify config %d failed: %v", id, err)
		resp.Error(c, http.StatusInternalServerError, "delete notify config failed")
		return
	}
	log.Infof("Notify config %d deleted by from %s", id, c.ClientIP())

	resp.Success(c, nil)
}

// getTemplates 获取通知模板
// @Summary 获取通知模板
// @Tags 通知
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} resp.SuccessStruct{data=[]notifyModel.Template} "获取成功"
// @Failure 400 {object} resp.ErrorStruct "请求参数错误"
// @Failure 401 {object} resp.ErrorStruct "未授权"
// @Failure 500 {object} resp.ErrorStruct "服务器内部错误"
// @Router /api/v1/notify/template [get]
func getTemplates(c *gin.Context) {
	notifyTemplateList, err := op.GetNotifyTemplateList()
	if err != nil {
		resp.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	resp.Success(c, notifyTemplateList)
}

// updateTemplate 更新通知模板
// @Summary 更新通知模板
// @Description 根据请求体中的ID更新通知模板信息
// @Tags 通知
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body notifyModel.Template true "更新通知模板请求"
// @Success 200 {object} resp.SuccessStruct{data=notifyModel.Template} "更新成功"
// @Failure 400 {object} resp.ErrorStruct "请求参数错误"
// @Failure 401 {object} resp.ErrorStruct "未授权"
// @Failure 404 {object} resp.ErrorStruct "通知模板不存在"
// @Failure 500 {object} resp.ErrorStruct "服务器内部错误"
// @Router /api/v1/notify/template [put]
func updateTemplate(c *gin.Context) {
	var req notifyModel.Template
	if err := c.ShouldBindJSON(&req); err != nil {
		resp.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	if err := op.UpdateNotifyTemplate(c.Request.Context(), &req); err != nil {
		resp.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	log.Infof("Notify template %s updated by from %s", req.Type, c.ClientIP())
	resp.Success(c, req)
}
