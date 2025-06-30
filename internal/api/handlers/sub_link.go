package handlers

import (
	"context"
	"net/http"
	"strconv"
	"strings"

	"github.com/bestruirui/bestsub/internal/api/middleware"
	"github.com/bestruirui/bestsub/internal/api/models"
	"github.com/bestruirui/bestsub/internal/api/router"
	"github.com/bestruirui/bestsub/internal/core/subscription"
	"github.com/bestruirui/bestsub/internal/database"
	"github.com/bestruirui/bestsub/internal/models/sublink"
	"github.com/bestruirui/bestsub/internal/utils/log"
	"github.com/gin-gonic/gin"
)

// subLinkHandler 订阅链接处理器
type subLinkHandler struct{}

// init 函数用于自动注册路由
func init() {
	h := newSubLinkHandler()

	// 需要认证的订阅链接路由
	router.NewGroupRouter("/api/v1/sub-links").
		Use(middleware.Auth()).
		AddRoute(
			router.NewRoute("", router.POST).
				Handle(h.createSubLinks).
				WithDescription("Create subscription links"),
		).
		AddRoute(
			router.NewRoute("", router.GET).
				Handle(h.getSubLinks).
				WithDescription("Get subscription links or list all with pagination"),
		).
		AddRoute(
			router.NewRoute("", router.PATCH).
				Handle(h.updateSubLinks).
				WithDescription("Update subscription links"),
		).
		AddRoute(
			router.NewRoute("/:id", router.DELETE).
				Handle(h.deleteSubLink).
				WithDescription("Delete subscription link"),
		).
		AddRoute(
			router.NewRoute("/:id/fetch", router.POST).
				Handle(h.fetchSubLink).
				WithDescription("Fetch and parse subscription link content"),
		)
}

// newSubLinkHandler 创建订阅链接处理器
func newSubLinkHandler() *subLinkHandler {
	return &subLinkHandler{}
}

// createSubLinks 创建订阅链接
// @Summary 创建订阅链接
// @Description 创建单个订阅链接
// @Tags 订阅链接管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body models.SubLinkCreateRequest true "创建订阅链接请求"
// @Success 200 {object} models.SuccessResponse{data=sublink.Data} "创建成功"
// @Failure 400 {object} models.ErrorResponse "请求参数错误"
// @Failure 401 {object} models.ErrorResponse "未授权"
// @Failure 500 {object} models.ErrorResponse "服务器内部错误"
// @Router /api/v1/sub-links [post]
func (h *subLinkHandler) createSubLinks(c *gin.Context) {
	var req models.SubLinkCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Code:    http.StatusBadRequest,
			Message: "Bad Request",
			Error:   "Invalid request format: " + err.Error(),
		})
		return
	}

	// 验证必填字段
	if req.Name == "" || req.FetchConfig.URL == "" {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Code:    http.StatusBadRequest,
			Message: "Bad Request",
			Error:   "Name and URL are required",
		})
		return
	}

	// 验证URL格式
	if !strings.HasPrefix(req.FetchConfig.URL, "http://") && !strings.HasPrefix(req.FetchConfig.URL, "https://") {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Code:    http.StatusBadRequest,
			Message: "Bad Request",
			Error:   "Invalid URL scheme, must be http or https",
		})
		return
	}

	if req.FetchConfig.Timeout <= 0 {
		req.FetchConfig.Timeout = 5
	}
	if req.FetchConfig.UserAgent == "" {
		req.FetchConfig.UserAgent = "clash.meta"
	}

	subLinkRepo := database.SubLink()

	// 检查URL是否已存在
	existingLink, err := subLinkRepo.GetByURL(context.Background(), req.FetchConfig.URL)
	if err != nil {
		log.Errorf("Failed to check existing URL %s: %v", req.FetchConfig.URL, err)
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Code:    http.StatusInternalServerError,
			Message: "Internal Server Error",
			Error:   "Failed to check URL uniqueness",
		})
		return
	}
	if existingLink != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Code:    http.StatusBadRequest,
			Message: "Bad Request",
			Error:   "URL already exists",
		})
		return
	}

	// 创建数据库模型
	dbLink := &sublink.Data{
		BaseData: sublink.BaseData{
			Name: req.Name,
			FetchConfig: sublink.FetchConfig{
				URL:         req.FetchConfig.URL,
				Type:        req.FetchConfig.Type,
				UserAgent:   req.FetchConfig.UserAgent,
				ProxyEnable: req.FetchConfig.ProxyEnable,
				Timeout:     req.FetchConfig.Timeout,
				Retries:     req.FetchConfig.Retries,
			},
			IsEnabled: req.IsEnabled,
			Detector:  req.Detector,
			Notify:    req.Notify,
			CronExpr:  req.CronExpr,
		},
		LastStatus: "pending",
	}

	// 创建链接
	err = subLinkRepo.Create(context.Background(), dbLink)
	if err != nil {
		log.Errorf("Failed to create subscription link %s: %v", req.Name, err)
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Code:    http.StatusInternalServerError,
			Message: "Internal Server Error",
			Error:   "Failed to create subscription link",
		})
		return
	}

	log.Infof("Subscription link %s created successfully with ID %d", req.Name, dbLink.ID)

	c.JSON(http.StatusOK, models.SuccessResponse{
		Code:    http.StatusOK,
		Message: "Subscription link created successfully",
		Data:    dbLink,
	})
}

// getSubLinks 获取订阅链接
// @Summary 获取订阅链接
// @Tags 订阅链接管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param page query int false "页码" default(1)
// @Param page_size query int false "每页大小" default(10)
// @Param ids query string false "链接ID列表，逗号分隔"
// @Success 200 {object} models.SuccessResponse{data=[]sublink.Data} "获取成功"
// @Failure 400 {object} models.ErrorResponse "请求参数错误"
// @Failure 401 {object} models.ErrorResponse "未授权"
// @Failure 500 {object} models.ErrorResponse "服务器内部错误"
// @Router /api/v1/sub-links [get]
func (h *subLinkHandler) getSubLinks(c *gin.Context) {
	subLinkRepo := database.SubLink()

	// 检查是否有查询参数中的ids
	idsParam := c.Query("ids")

	// 如果有查询参数中的ids，执行批量获取
	if idsParam != "" {
		var targetIDs []int64

		// 处理查询参数中的IDs
		idStrs := strings.Split(idsParam, ",")
		for _, idStr := range idStrs {
			idStr = strings.TrimSpace(idStr)
			if id, err := strconv.ParseInt(idStr, 10, 64); err == nil {
				targetIDs = append(targetIDs, id)
			}
		}

		if len(targetIDs) == 0 {
			c.JSON(http.StatusBadRequest, models.ErrorResponse{
				Code:    http.StatusBadRequest,
				Message: "Bad Request",
				Error:   "No valid subscription link IDs provided",
			})
			return
		}

		// 批量获取链接
		var successLinks []sublink.Data
		for _, id := range targetIDs {
			dbLink, err := subLinkRepo.GetByID(context.Background(), id)
			if err != nil {
				log.Errorf("Failed to get subscription link %d: %v", id, err)
				continue
			}

			if dbLink == nil {
				log.Warnf("Subscription link %d not found", id)
				continue
			}

			// 转换为响应模型
			successLinks = append(successLinks, *dbLink)
		}

		response := models.SubLinkListResponse{
			Items: successLinks,
			Total: len(successLinks),
		}

		c.JSON(http.StatusOK, models.SuccessResponse{
			Code:    http.StatusOK,
			Message: "Subscription links retrieved successfully",
			Data:    response,
		})
		return
	}

	// 如果没有指定IDs，执行分页查询
	page := 1
	pageSize := 10

	if pageStr := c.Query("page"); pageStr != "" {
		if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
			page = p
		}
	}

	if pageSizeStr := c.Query("page_size"); pageSizeStr != "" {
		if ps, err := strconv.Atoi(pageSizeStr); err == nil && ps > 0 && ps <= 100 {
			pageSize = ps
		}
	}

	offset := (page - 1) * pageSize

	// 获取总数
	total, err := subLinkRepo.Count(context.Background())
	if err != nil {
		log.Errorf("Failed to get subscription links count: %v", err)
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Code:    http.StatusInternalServerError,
			Message: "Internal Server Error",
			Error:   "Failed to get subscription links count",
		})
		return
	}

	// 获取链接列表
	dbLinks, err := subLinkRepo.List(context.Background(), offset, pageSize)
	if err != nil {
		log.Errorf("Failed to get subscription links list: %v", err)
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Code:    http.StatusInternalServerError,
			Message: "Internal Server Error",
			Error:   "Failed to get subscription links list",
		})
		return
	}

	// 转换为响应模型
	items := make([]sublink.Data, 0, len(dbLinks))
	for _, dbLink := range dbLinks {
		items = append(items, *dbLink)
	}

	response := models.SubLinkListResponse{
		Items: items,
		Total: int(total),
	}

	log.Debugf("Retrieved %d subscription links (page %d, size %d)", len(items), page, pageSize)

	c.JSON(http.StatusOK, models.SuccessResponse{
		Code:    http.StatusOK,
		Message: "Subscription links retrieved successfully",
		Data:    response,
	})
}

// updateSubLinks 更新订阅链接
// @Summary 更新订阅链接
// @Description 根据请求体中的ID更新单个订阅链接信息
// @Tags 订阅链接管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body models.SubLinkUpdateRequest true "更新订阅链接请求"
// @Success 200 {object} models.SuccessResponse{data=sublink.Data} "更新成功"
// @Failure 400 {object} models.ErrorResponse "请求参数错误"
// @Failure 401 {object} models.ErrorResponse "未授权"
// @Failure 404 {object} models.ErrorResponse "订阅链接不存在"
// @Failure 500 {object} models.ErrorResponse "服务器内部错误"
// @Router /api/v1/sub-links [patch]
func (h *subLinkHandler) updateSubLinks(c *gin.Context) {
	var req models.SubLinkUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Code:    http.StatusBadRequest,
			Message: "Bad Request",
			Error:   "Invalid request format: " + err.Error(),
		})
		return
	}

	if req.ID == 0 {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Code:    http.StatusBadRequest,
			Message: "Bad Request",
			Error:   "ID is required for update operation",
		})
		return
	}

	subLinkRepo := database.SubLink()

	// 获取现有链接
	dbLink, err := subLinkRepo.GetByID(context.Background(), req.ID)
	if err != nil {
		log.Errorf("Failed to get subscription link %d: %v", req.ID, err)
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Code:    http.StatusInternalServerError,
			Message: "Internal Server Error",
			Error:   "Failed to get subscription link",
		})
		return
	}

	if dbLink == nil {
		c.JSON(http.StatusNotFound, models.ErrorResponse{
			Code:    http.StatusNotFound,
			Message: "Not Found",
			Error:   "Subscription link not found",
		})
		return
	}

	// 更新字段
	if req.Name != "" {
		dbLink.Name = req.Name
	}
	if req.FetchConfig.URL != "" {
		// 检查新URL是否与其他链接冲突
		existingLink, err := subLinkRepo.GetByURL(context.Background(), req.FetchConfig.URL)
		if err != nil {
			log.Errorf("Failed to check URL uniqueness for %s: %v", req.FetchConfig.URL, err)
			c.JSON(http.StatusInternalServerError, models.ErrorResponse{
				Code:    http.StatusInternalServerError,
				Message: "Internal Server Error",
				Error:   "Failed to check URL uniqueness",
			})
			return
		}
		if existingLink != nil && existingLink.ID != req.ID {
			c.JSON(http.StatusBadRequest, models.ErrorResponse{
				Code:    http.StatusBadRequest,
				Message: "Bad Request",
				Error:   "URL already exists in another link",
			})
			return
		}
		dbLink.FetchConfig.URL = req.FetchConfig.URL
	}
	if req.FetchConfig.Type != "" {
		dbLink.FetchConfig.Type = req.FetchConfig.Type
	}
	if req.FetchConfig.UserAgent != "" {
		dbLink.FetchConfig.UserAgent = req.FetchConfig.UserAgent
	}
	dbLink.IsEnabled = req.IsEnabled
	dbLink.FetchConfig.ProxyEnable = req.FetchConfig.ProxyEnable
	if req.Detector != nil {
		dbLink.Detector = req.Detector
	}
	if req.Notify != nil {
		dbLink.Notify = req.Notify
	}
	if req.CronExpr != "" {
		dbLink.CronExpr = req.CronExpr
	}

	// 更新链接
	err = subLinkRepo.Update(context.Background(), dbLink)
	if err != nil {
		log.Errorf("Failed to update subscription link %d: %v", req.ID, err)
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Code:    http.StatusInternalServerError,
			Message: "Internal Server Error",
			Error:   "Failed to update subscription link",
		})
		return
	}

	log.Infof("Subscription link %d updated successfully", req.ID)

	c.JSON(http.StatusOK, models.SuccessResponse{
		Code:    http.StatusOK,
		Message: "Subscription link updated successfully",
		Data:    dbLink,
	})
}

// deleteSubLink 删除订阅链接
// @Summary 删除订阅链接
// @Description 根据ID删除单个订阅链接
// @Tags 订阅链接管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "订阅链接ID"
// @Success 200 {object} models.SuccessResponse "删除成功"
// @Failure 400 {object} models.ErrorResponse "请求参数错误"
// @Failure 401 {object} models.ErrorResponse "未授权"
// @Failure 404 {object} models.ErrorResponse "订阅链接不存在"
// @Failure 500 {object} models.ErrorResponse "服务器内部错误"
// @Router /api/v1/sub-links/{id} [delete]
func (h *subLinkHandler) deleteSubLink(c *gin.Context) {
	idStr := c.Param("id")
	if idStr == "" {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Code:    http.StatusBadRequest,
			Message: "Bad Request",
			Error:   "Subscription link ID is required",
		})
		return
	}

	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Code:    http.StatusBadRequest,
			Message: "Bad Request",
			Error:   "Invalid subscription link ID format",
		})
		return
	}

	subLinkRepo := database.SubLink()

	// 检查链接是否存在
	dbLink, err := subLinkRepo.GetByID(context.Background(), id)
	if err != nil {
		log.Errorf("Failed to get subscription link %d: %v", id, err)
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Code:    http.StatusInternalServerError,
			Message: "Internal Server Error",
			Error:   "Failed to get subscription link",
		})
		return
	}

	if dbLink == nil {
		c.JSON(http.StatusNotFound, models.ErrorResponse{
			Code:    http.StatusNotFound,
			Message: "Not Found",
			Error:   "Subscription link not found",
		})
		return
	}

	// 删除链接
	err = subLinkRepo.Delete(context.Background(), id)
	if err != nil {
		log.Errorf("Failed to delete subscription link %d: %v", id, err)
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Code:    http.StatusInternalServerError,
			Message: "Internal Server Error",
			Error:   "Failed to delete subscription link",
		})
		return
	}

	log.Infof("Subscription link %d deleted successfully", id)

	c.JSON(http.StatusOK, models.SuccessResponse{
		Code:    http.StatusOK,
		Message: "Subscription link deleted successfully",
	})
}

// fetchSubLink 刷新订阅链接
// @Summary 刷新订阅链接
// @Description 从指定的订阅链接获取内容并解析，返回解析出的节点数量
// @Tags 订阅链接管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "订阅链接ID"
// @Success 200 {object} models.SuccessResponse{data=sublink.FetchResult} "获取成功"
// @Failure 400 {object} models.ErrorResponse "请求参数错误"
// @Failure 401 {object} models.ErrorResponse "未授权"
// @Failure 404 {object} models.ErrorResponse "订阅链接不存在"
// @Failure 500 {object} models.ErrorResponse "服务器内部错误"
// @Router /api/v1/sub-links/{id}/fetch [post]
func (h *subLinkHandler) fetchSubLink(c *gin.Context) {
	idStr := c.Param("id")
	if idStr == "" {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Code:    http.StatusBadRequest,
			Message: "Bad Request",
			Error:   "Subscription link ID is required",
		})
		return
	}

	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Code:    http.StatusBadRequest,
			Message: "Bad Request",
			Error:   "Invalid subscription link ID format",
		})
		return
	}

	subLinkRepo := database.SubLink()

	// 获取订阅链接信息
	dbLink, err := subLinkRepo.GetByID(context.Background(), id)
	if err != nil {
		log.Errorf("Failed to get subscription link %d: %v", id, err)
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Code:    http.StatusInternalServerError,
			Message: "Internal Server Error",
			Error:   "Failed to get subscription link",
		})
		return
	}

	if dbLink == nil {
		c.JSON(http.StatusNotFound, models.ErrorResponse{
			Code:    http.StatusNotFound,
			Message: "Not Found",
			Error:   "Subscription link not found",
		})
		return
	}

	// 检查链接是否启用
	if !dbLink.IsEnabled {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Code:    http.StatusBadRequest,
			Message: "Bad Request",
			Error:   "Subscription link is disabled",
		})
		return
	}

	// 使用订阅链接的配置获取内容
	result, err := subscription.Fetch(context.Background(), dbLink.FetchConfig)
	if err != nil {
		log.Errorf("Failed to fetch content from %s: %v", dbLink.FetchConfig.URL, err)
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Code:    http.StatusInternalServerError,
			Message: "Internal Server Error",
			Error:   "Failed to fetch subscription content: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, models.SuccessResponse{
		Code:    http.StatusOK,
		Message: "Subscription content fetched and parsed successfully",
		Data:    result,
	})
}
