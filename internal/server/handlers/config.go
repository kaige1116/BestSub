package handlers

import (
	"context"
	"net/http"

	"github.com/bestruirui/bestsub/internal/database/op"
	"github.com/bestruirui/bestsub/internal/models/config"
	"github.com/bestruirui/bestsub/internal/server/middleware"
	"github.com/bestruirui/bestsub/internal/server/resp"
	"github.com/bestruirui/bestsub/internal/server/router"
	"github.com/bestruirui/bestsub/internal/utils/log"
	"github.com/gin-gonic/gin"
)

func init() {

	router.NewGroupRouter("/api/v1/config").
		Use(middleware.Auth()).
		AddRoute(
			router.NewRoute("/item", router.GET).
				Handle(getConfigItems),
		).
		AddRoute(
			router.NewRoute("/item", router.PUT).
				Handle(updateConfigItem),
		)
}

// getConfigItems 获取配置项
// @Summary 获取配置项
// @Description 获取系统所有配置项，支持按分组过滤和关键字搜索
// @Tags 配置
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param group query string false "分组名称"
// @Success 200 {object} resp.SuccessStruct{data=[]config.GroupAdvance} "获取成功"
// @Failure 401 {object} resp.ErrorStruct "未授权"
// @Failure 500 {object} resp.ErrorStruct "服务器内部错误"
// @Router /api/v1/config/item [get]
func getConfigItems(c *gin.Context) {
	group := c.Query("group")
	var dbConfig []config.GroupAdvance
	var err error
	if group == "" {
		dbConfig, err = op.GetAllConfig(context.Background())
		if err != nil {
			log.Errorf("Failed to get all config: %v", err)
			resp.Error(c, http.StatusInternalServerError, "failed to get all config")
			return
		}
	} else {
		dbConfig, err = op.GetConfigByGroup(group)
		if err != nil {
			log.Errorf("Failed to get config by group: %v", err)
			resp.Error(c, http.StatusInternalServerError, "failed to get config by group")
			return
		}
	}

	resp.Success(c, dbConfig)
}

// updateConfigItem 更新配置项
// @Summary 更新配置项
// @Description 根据请求数据中的ID批量更新配置项的值和描述
// @Tags 配置
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body []config.UpdateAdvance  true "更新配置项请求"
// @Success 200 {object} resp.SuccessStruct{data=[]config.UpdateAdvance} "更新成功"
// @Failure 400 {object} resp.ErrorStruct "请求参数错误"
// @Failure 401 {object} resp.ErrorStruct "未授权"
// @Failure 500 {object} resp.ErrorStruct "服务器内部错误"
// @Router /api/v1/config/item [put]
func updateConfigItem(c *gin.Context) {
	var req []config.UpdateAdvance
	if err := c.ShouldBindJSON(&req); err != nil {
		resp.ErrorBadRequest(c)
		return
	}

	err := op.UpdateConfig(context.Background(), &req)
	if err != nil {
		log.Errorf("Failed to update config: %v", err)
		resp.Error(c, http.StatusInternalServerError, "failed to update config")
		return
	}

	resp.Success(c, req)
}
