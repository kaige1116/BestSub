package handlers

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/bestruirui/bestsub/internal/api/middleware"
	"github.com/bestruirui/bestsub/internal/api/models"
	"github.com/bestruirui/bestsub/internal/api/router"
	"github.com/bestruirui/bestsub/internal/database"
	dbModels "github.com/bestruirui/bestsub/internal/database/models"
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
// @Success 200 {object} models.SuccessResponse{data=models.ConfigItemsResponse} "获取成功"
// @Failure 401 {object} models.ErrorResponse "未授权"
// @Failure 500 {object} models.ErrorResponse "服务器内部错误"
// @Router /api/v1/config/items [get]
func (h *configHandler) getConfigItems(c *gin.Context) {
	group := c.Query("group")
	search := c.Query("search")

	configRepo := database.SystemConfig()
	var configs []dbModels.SystemConfig
	var err error

	if group != "" {
		// 按分组获取配置
		configs, err = configRepo.GetConfigsByGroup(context.Background(), group)
		if err != nil {
			log.Errorf("Failed to get configs by group %s: %v", group, err)
			c.JSON(http.StatusInternalServerError, models.ErrorResponse{
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
			c.JSON(http.StatusInternalServerError, models.ErrorResponse{
				Code:    http.StatusInternalServerError,
				Message: "Internal Server Error",
				Error:   "Failed to get configuration items",
			})
			return
		}

		configs = make([]dbModels.SystemConfig, 0, len(keys))
		for _, key := range keys {
			config, err := configRepo.GetByKey(context.Background(), key)
			if err != nil {
				log.Errorf("Failed to get config by key %s: %v", key, err)
				continue
			}
			if config != nil {
				configs = append(configs, *config)
			}
		}
	}

	// 如果有搜索关键字，进行过滤
	if search != "" {
		filteredConfigs := make([]dbModels.SystemConfig, 0)
		searchLower := strings.ToLower(search)
		for _, config := range configs {
			if strings.Contains(strings.ToLower(config.Key), searchLower) ||
				strings.Contains(strings.ToLower(config.Description), searchLower) {
				filteredConfigs = append(filteredConfigs, config)
			}
		}
		configs = filteredConfigs
	}

	// 转换为响应模型
	items := make([]models.ConfigItemResponse, 0, len(configs))
	for _, config := range configs {
		items = append(items, models.ConfigItemResponse{
			ConfigItemData: models.ConfigItemData{
				ID:          config.ID,
				Value:       config.Value,
				Description: config.Description,
			},
			GroupName: config.GroupName,
			Key:       config.Key,
			Type:      config.Type,
		})
	}

	response := models.ConfigItemsResponse{
		Items: items,
		Total: len(items),
	}

	log.Debugf("Retrieved %d configuration items", len(items))

	c.JSON(http.StatusOK, models.SuccessResponse{
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
// @Success 200 {object} models.SuccessResponse{data=models.ConfigItemResponse} "获取成功"
// @Failure 400 {object} models.ErrorResponse "请求参数错误"
// @Failure 401 {object} models.ErrorResponse "未授权"
// @Failure 404 {object} models.ErrorResponse "配置项不存在"
// @Failure 500 {object} models.ErrorResponse "服务器内部错误"
// @Router /api/v1/config/items/{id} [get]
func (h *configHandler) getConfigItem(c *gin.Context) {
	idStr := c.Param("id")
	if idStr == "" {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Code:    http.StatusBadRequest,
			Message: "Bad Request",
			Error:   "Configuration item ID is required",
		})
		return
	}

	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Code:    http.StatusBadRequest,
			Message: "Bad Request",
			Error:   "Invalid configuration item ID format",
		})
		return
	}

	// 由于接口中没有直接通过ID获取的方法，我们需要获取所有配置然后查找
	configRepo := database.SystemConfig()
	keys, err := configRepo.GetAllKeys(context.Background())
	if err != nil {
		log.Errorf("Failed to get all config keys: %v", err)
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Code:    http.StatusInternalServerError,
			Message: "Internal Server Error",
			Error:   "Failed to get configuration item",
		})
		return
	}

	var foundConfig *dbModels.SystemConfig
	for _, key := range keys {
		config, err := configRepo.GetByKey(context.Background(), key)
		if err != nil {
			log.Errorf("Failed to get config by key %s: %v", key, err)
			continue
		}
		if config != nil && config.ID == id {
			foundConfig = config
			break
		}
	}

	if foundConfig == nil {
		c.JSON(http.StatusNotFound, models.ErrorResponse{
			Code:    http.StatusNotFound,
			Message: "Not Found",
			Error:   "Configuration item not found",
		})
		return
	}

	response := models.ConfigItemResponse{
		ConfigItemData: models.ConfigItemData{
			ID:          foundConfig.ID,
			Value:       foundConfig.Value,
			Description: foundConfig.Description,
		},
		GroupName: foundConfig.GroupName,
		Key:       foundConfig.Key,
		Type:      foundConfig.Type,
	}

	log.Debugf("Retrieved configuration item %d: %s", id, foundConfig.Key)

	c.JSON(http.StatusOK, models.SuccessResponse{
		Code:    http.StatusOK,
		Message: "Configuration item retrieved successfully",
		Data:    response,
	})
}

// updateConfigItem 更新配置项
// @Summary 批量更新配置项
// @Description 根据请求数据中的ID批量更新配置项的值和描述
// @Tags 配置管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body models.UpdateConfigItemRequest true "更新配置项请求"
// @Success 200 {object} models.SuccessResponse{data=models.ConfigItemsResponse} "更新成功"
// @Failure 400 {object} models.ErrorResponse "请求参数错误"
// @Failure 401 {object} models.ErrorResponse "未授权"
// @Failure 500 {object} models.ErrorResponse "服务器内部错误"
// @Router /api/v1/config/items [patch]
func (h *configHandler) updateConfigItem(c *gin.Context) {
	var req models.UpdateConfigItemRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Code:    http.StatusBadRequest,
			Message: "Bad Request",
			Error:   "Invalid request format: " + err.Error(),
		})
		return
	}

	if len(req.Data) == 0 {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Code:    http.StatusBadRequest,
			Message: "Bad Request",
			Error:   "No configuration items to update",
		})
		return
	}

	configRepo := database.SystemConfig()

	// 获取所有配置项以便查找
	keys, err := configRepo.GetAllKeys(context.Background())
	if err != nil {
		log.Errorf("Failed to get all config keys: %v", err)
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Code:    http.StatusInternalServerError,
			Message: "Internal Server Error",
			Error:   "Failed to update configuration items",
		})
		return
	}

	// 创建ID到配置的映射
	configMap := make(map[int64]*dbModels.SystemConfig)
	for _, key := range keys {
		config, err := configRepo.GetByKey(context.Background(), key)
		if err != nil {
			log.Errorf("Failed to get config by key %s: %v", key, err)
			continue
		}
		if config != nil {
			configMap[config.ID] = config
		}
	}

	// 批量更新配置项
	updatedConfigs := make([]models.ConfigItemResponse, 0, len(req.Data))
	for _, item := range req.Data {
		foundConfig, exists := configMap[item.ID]
		if !exists {
			log.Warnf("Configuration item with ID %d not found, skipping", item.ID)
			continue
		}

		// 更新配置项的值
		foundConfig.Value = item.Value
		if item.Description != "" {
			foundConfig.Description = item.Description
		}

		err = configRepo.Update(context.Background(), foundConfig)
		if err != nil {
			log.Errorf("Failed to update config %s (ID: %d): %v", foundConfig.Key, item.ID, err)
			continue
		}

		// 重新获取更新后的配置以获取最新的更新时间
		updatedConfig, err := configRepo.GetByKey(context.Background(), foundConfig.Key)
		if err != nil {
			log.Errorf("Failed to get updated config %s: %v", foundConfig.Key, err)
			updatedConfig = foundConfig
		}

		updatedConfigs = append(updatedConfigs, models.ConfigItemResponse{
			ConfigItemData: models.ConfigItemData{
				ID:          updatedConfig.ID,
				Value:       updatedConfig.Value,
				Description: updatedConfig.Description,
			},
			GroupName: updatedConfig.GroupName,
			Key:       updatedConfig.Key,
			Type:      updatedConfig.Type,
		})

		log.Infof("Configuration item %d (%s) updated successfully", item.ID, foundConfig.Key)
	}

	if len(updatedConfigs) == 0 {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Code:    http.StatusBadRequest,
			Message: "Bad Request",
			Error:   "No valid configuration items were updated",
		})
		return
	}

	response := models.ConfigItemsResponse{
		Items: updatedConfigs,
		Total: len(updatedConfigs),
	}

	c.JSON(http.StatusOK, models.SuccessResponse{
		Code:    http.StatusOK,
		Message: fmt.Sprintf("Successfully updated %d configuration items", len(updatedConfigs)),
		Data:    response,
	})
}
