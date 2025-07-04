package handlers

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/bestruirui/bestsub/internal/api/middleware"
	"github.com/bestruirui/bestsub/internal/api/router"
	"github.com/bestruirui/bestsub/internal/database"
	"github.com/bestruirui/bestsub/internal/models/api"
	"github.com/bestruirui/bestsub/internal/models/system"
	"github.com/bestruirui/bestsub/internal/utils/log"
	"github.com/gin-gonic/gin"
)

// configHandler 配置处理器
type configHandler struct{}

// init 函数用于自动注册路由
func init() {
	h := newConfigHandler()

	// 需要认证的配置路由
	router.NewGroupRouter("/api/v1/config").
		Use(middleware.Auth()).
		AddRoute(
			router.NewRoute("/items", router.GET).
				Handle(h.getConfigItems).
				WithDescription("Get all configuration items"),
		).
		AddRoute(
			router.NewRoute("/items/:id", router.GET).
				Handle(h.getConfigItem).
				WithDescription("Get single configuration item"),
		).
		AddRoute(
			router.NewRoute("/items", router.PATCH).
				Handle(h.updateConfigItem).
				WithDescription("Batch update configuration items"),
		)
}

// newConfigHandler 创建配置处理器
func newConfigHandler() *configHandler {
	return &configHandler{}
}

// getConfigItems 获取所有配置项
// @Summary 获取所有配置项
// @Description 获取系统所有配置项，支持按分组过滤和关键字搜索
// @Tags 配置管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param group query string false "配置分组过滤"
// @Param search query string false "关键字搜索配置名或描述"
// @Success 200 {object} api.ResponseSuccess{data=system.ConfigItemsResponse} "获取成功"
// @Failure 401 {object} api.ResponseError "未授权"
// @Failure 500 {object} api.ResponseError "服务器内部错误"
// @Router /api/v1/config/items [get]
func (h *configHandler) getConfigItems(c *gin.Context) {
	group := c.Query("group")
	search := c.Query("search")

	configRepo := database.SystemConfig()
	var configs []system.Data
	var err error

	if group != "" {
		// 按分组获取配置
		configs, err = configRepo.GetConfigsByGroup(context.Background(), group)
		if err != nil {
			log.Errorf("Failed to get configs by group %s: %v", group, err)
			c.JSON(http.StatusInternalServerError, api.ResponseError{
				Code:    http.StatusInternalServerError,
				Message: "Internal Server Error",
				Error:   "Failed to get configuration items",
			})
			return
		}
	} else {
		// 获取所有配置键，然后逐个获取配置
		keys, err := configRepo.GetAllKeys(context.Background())
		if err != nil {
			log.Errorf("Failed to get all config keys: %v", err)
			c.JSON(http.StatusInternalServerError, api.ResponseError{
				Code:    http.StatusInternalServerError,
				Message: "Internal Server Error",
				Error:   "Failed to get configuration items",
			})
			return
		}

		configs = make([]system.Data, 0, len(keys))
		for _, key := range keys {
			config, err := configRepo.GetByKey(context.Background(), key)
			if err != nil {
				log.Warnf("Failed to get config for key %s: %v", key, err)
				continue
			}
			if config != nil {
				configs = append(configs, *config)
			}
		}
	}

	// 如果有搜索关键字，进行过滤
	if search != "" {
		filteredConfigs := make([]system.Data, 0)
		searchLower := strings.ToLower(search)
		for _, config := range configs {
			if strings.Contains(strings.ToLower(config.Key), searchLower) ||
				strings.Contains(strings.ToLower(config.Description), searchLower) {
				filteredConfigs = append(filteredConfigs, config)
			}
		}
		configs = filteredConfigs
	}

	response := system.ConfigItemsResponse{
		Data:  configs,
		Total: len(configs),
	}

	log.Debugf("Retrieved %d configuration items", len(configs))

	c.JSON(http.StatusOK, api.ResponseSuccess{
		Code:    http.StatusOK,
		Message: "Configuration items retrieved successfully",
		Data:    response,
	})
}

// getConfigItem 获取单个配置项
// @Summary 获取单个配置项
// @Description 根据ID获取指定的配置项详细信息
// @Tags 配置管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "配置项ID"
// @Success 200 {object} api.ResponseSuccess{data=system.Data} "获取成功"
// @Failure 400 {object} api.ResponseError "请求参数错误"
// @Failure 401 {object} api.ResponseError "未授权"
// @Failure 404 {object} api.ResponseError "配置项不存在"
// @Failure 500 {object} api.ResponseError "服务器内部错误"
// @Router /api/v1/config/items/{id} [get]
func (h *configHandler) getConfigItem(c *gin.Context) {
	idStr := c.Param("id")
	if idStr == "" {
		c.JSON(http.StatusBadRequest, api.ResponseError{
			Code:    http.StatusBadRequest,
			Message: "Bad Request",
			Error:   "Configuration item ID is required",
		})
		return
	}

	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, api.ResponseError{
			Code:    http.StatusBadRequest,
			Message: "Bad Request",
			Error:   "Invalid configuration item ID format",
		})
		return
	}

	// 由于SystemConfigRepository接口没有GetByID方法，我们需要获取所有配置然后查找
	configRepo := database.SystemConfig()
	keys, err := configRepo.GetAllKeys(context.Background())
	if err != nil {
		log.Errorf("Failed to get all config keys: %v", err)
		c.JSON(http.StatusInternalServerError, api.ResponseError{
			Code:    http.StatusInternalServerError,
			Message: "Internal Server Error",
			Error:   "Failed to get configuration item",
		})
		return
	}

	var targetConfig *system.Data
	for _, key := range keys {
		config, err := configRepo.GetByKey(context.Background(), key)
		if err != nil {
			continue
		}
		if config != nil && config.ID == id {
			targetConfig = config
			break
		}
	}

	if targetConfig == nil {
		c.JSON(http.StatusNotFound, api.ResponseError{
			Code:    http.StatusNotFound,
			Message: "Not Found",
			Error:   "Configuration item not found",
		})
		return
	}

	log.Debugf("Retrieved configuration item with ID %d", id)

	c.JSON(http.StatusOK, api.ResponseSuccess{
		Code:    http.StatusOK,
		Message: "Configuration item retrieved successfully",
		Data:    targetConfig,
	})
}

// updateConfigItem 更新配置项
// @Summary 批量更新配置项
// @Description 根据请求数据中的ID批量更新配置项的值和描述
// @Tags 配置管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body system.UpdateConfigItemRequest true "更新配置项请求"
// @Success 200 {object} api.ResponseSuccess{data=system.ConfigItemsResponse} "更新成功"
// @Failure 400 {object} api.ResponseError "请求参数错误"
// @Failure 401 {object} api.ResponseError "未授权"
// @Failure 500 {object} api.ResponseError "服务器内部错误"
// @Router /api/v1/config/items [patch]
func (h *configHandler) updateConfigItem(c *gin.Context) {
	var req system.UpdateConfigItemRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, api.ResponseError{
			Code:    http.StatusBadRequest,
			Message: "Bad Request",
			Error:   "Invalid request format: " + err.Error(),
		})
		return
	}

	if len(req.Data) == 0 {
		c.JSON(http.StatusBadRequest, api.ResponseError{
			Code:    http.StatusBadRequest,
			Message: "Bad Request",
			Error:   "No configuration items to update",
		})
		return
	}

	configRepo := database.SystemConfig()
	updatedConfigs := make([]system.Data, 0, len(req.Data))

	// 首先获取所有现有配置以便查找
	keys, err := configRepo.GetAllKeys(context.Background())
	if err != nil {
		log.Errorf("Failed to get all config keys: %v", err)
		c.JSON(http.StatusInternalServerError, api.ResponseError{
			Code:    http.StatusInternalServerError,
			Message: "Internal Server Error",
			Error:   "Failed to get configuration items",
		})
		return
	}

	// 创建ID到配置的映射
	configMap := make(map[int64]*system.Data)
	for _, key := range keys {
		config, err := configRepo.GetByKey(context.Background(), key)
		if err != nil {
			continue
		}
		if config != nil {
			configMap[config.ID] = config
		}
	}

	// 批量更新配置项
	for _, updateItem := range req.Data {
		existingConfig, exists := configMap[updateItem.ID]
		if !exists {
			log.Warnf("Configuration item with ID %d not found, skipping", updateItem.ID)
			continue
		}

		// 更新配置值和描述
		existingConfig.Value = updateItem.Value
		existingConfig.Description = updateItem.Description

		err := configRepo.Update(context.Background(), existingConfig)
		if err != nil {
			log.Errorf("Failed to update config item %d: %v", updateItem.ID, err)
			c.JSON(http.StatusInternalServerError, api.ResponseError{
				Code:    http.StatusInternalServerError,
				Message: "Internal Server Error",
				Error:   fmt.Sprintf("Failed to update configuration item %d", updateItem.ID),
			})
			return
		}

		updatedConfigs = append(updatedConfigs, *existingConfig)
		log.Debugf("Updated configuration item %d (%s)", updateItem.ID, existingConfig.Key)
	}

	response := system.ConfigItemsResponse{
		Data:  updatedConfigs,
		Total: len(updatedConfigs),
	}

	log.Infof("Successfully updated %d configuration items", len(updatedConfigs))

	c.JSON(http.StatusOK, api.ResponseSuccess{
		Code:    http.StatusOK,
		Message: "Configuration items updated successfully",
		Data:    response,
	})
}
