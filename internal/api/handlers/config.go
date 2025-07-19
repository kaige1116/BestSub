package handlers

import (
	"context"
	"net/http"

	"github.com/bestruirui/bestsub/internal/api/middleware"
	"github.com/bestruirui/bestsub/internal/api/resp"
	"github.com/bestruirui/bestsub/internal/api/router"
	"github.com/bestruirui/bestsub/internal/database/op"
	"github.com/bestruirui/bestsub/internal/models/system"
	"github.com/bestruirui/bestsub/internal/utils/log"
	"github.com/gin-gonic/gin"
)

// init 函数用于自动注册路由
func init() {

	// 需要认证的配置路由
	router.NewGroupRouter("/api/v1/config").
		Use(middleware.Auth()).
		AddRoute(
			router.NewRoute("/items", router.GET).
				Handle(getConfigItems).
				WithDescription("Get all configuration items"),
		).
		AddRoute(
			router.NewRoute("/items", router.PATCH).
				Handle(updateConfigItem).
				WithDescription("Batch update configuration items"),
		)
}

// getConfigItems 获取所有配置项
// @Summary 获取所有配置项
// @Description 获取系统所有配置项，支持按分组过滤和关键字搜索
// @Tags 配置管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} resp.SuccessStruct{data=[]system.GroupData} "获取成功"
// @Failure 401 {object} resp.ErrorStruct "未授权"
// @Failure 500 {object} resp.ErrorStruct "服务器内部错误"
// @Router /api/v1/config/items [get]
func getConfigItems(c *gin.Context) {

	dbConfig, err := op.GetAllConfig(context.Background())
	if err != nil {
		log.Errorf("Failed to get all config: %v", err)
		resp.Error(c, http.StatusInternalServerError, "failed to get all config")
		return
	}

	resp.Success(c, dbConfig)
}

// updateConfigItem 更新配置项
// @Summary 批量更新配置项
// @Description 根据请求数据中的ID批量更新配置项的值和描述
// @Tags 配置管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body []system.UpdateData  true "更新配置项请求"
// @Success 200 {object} resp.SuccessStruct{data=[]system.UpdateData} "更新成功"
// @Failure 400 {object} resp.ErrorStruct "请求参数错误"
// @Failure 401 {object} resp.ErrorStruct "未授权"
// @Failure 500 {object} resp.ErrorStruct "服务器内部错误"
// @Router /api/v1/config/items [patch]
func updateConfigItem(c *gin.Context) {
	var req []system.UpdateData
	if err := c.ShouldBindJSON(&req); err != nil {
		resp.ErrorBadRequest(c)
		return
	}

	err := op.UpdateConfig(context.Background(), &req)
	if err != nil {
		resp.Error(c, http.StatusInternalServerError, "failed to update config")
		return
	}

	resp.Success(c, req)
}
