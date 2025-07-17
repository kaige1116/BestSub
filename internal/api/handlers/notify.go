package handlers

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/bestruirui/bestsub/internal/api/middleware"
	"github.com/bestruirui/bestsub/internal/api/router"
	"github.com/bestruirui/bestsub/internal/database"
	"github.com/bestruirui/bestsub/internal/models/api"
	"github.com/bestruirui/bestsub/internal/models/notify"
	"github.com/bestruirui/bestsub/internal/utils/log"
	"github.com/gin-gonic/gin"
)

type notifyHandler struct{}

func init() {
	h := newNotifyHandler()
	router.NewGroupRouter("/api/v1/notify").
		Use(middleware.Auth()).
		AddRoute(
			router.NewRoute("/config", router.GET).
				Handle(h.getNotifies).
				WithDescription("Get notify"),
		).
		AddRoute(
			router.NewRoute("/config", router.POST).
				Handle(h.createNotify).
				WithDescription("Create notify"),
		).
		AddRoute(
			router.NewRoute("/config", router.PUT).
				Handle(h.updateNotify).
				WithDescription("Update notify"),
		).
		AddRoute(
			router.NewRoute("/config", router.DELETE).
				Handle(h.deleteNotify).
				WithDescription("Delete notify"),
		).
		AddRoute(
			router.NewRoute("/template", router.GET).
				Handle(h.getTemplates).
				WithDescription("Get notify template"),
		).
		AddRoute(
			router.NewRoute("/template", router.POST).
				Handle(h.createTemplate).
				WithDescription("Create notify template"),
		).
		AddRoute(
			router.NewRoute("/template", router.PUT).
				Handle(h.updateTemplate).
				WithDescription("Update notify template"),
		).
		AddRoute(
			router.NewRoute("/template", router.DELETE).
				Handle(h.deleteTemplate).
				WithDescription("Delete notify template"),
		)
}

func newNotifyHandler() *notifyHandler {
	return &notifyHandler{}
}

// getNotifies 获取通知配置
// @Summary 获取通知配置
// @Tags 通知管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param page query int false "页码" default(1)
// @Param page_size query int false "每页大小" default(10)
// @Param ids query string false "配置ID列表，逗号分隔"
// @Success 200 {object} api.ResponseSuccess{data=[]notify.Data} "获取成功"
// @Failure 400 {object} api.ResponseError "请求参数错误"
// @Failure 401 {object} api.ResponseError "未授权"
// @Failure 500 {object} api.ResponseError "服务器内部错误"
// @Router /api/v1/notify/config [get]
func (h *notifyHandler) getNotifies(c *gin.Context) {
	// 解析查询参数
	idsParam := c.Query("ids")

	// 如果指定了IDs，则获取指定的通知配置
	if idsParam != "" {
		idStrs := strings.Split(idsParam, ",")
		var responses []notify.Data

		for _, idStr := range idStrs {
			id, err := strconv.ParseInt(strings.TrimSpace(idStr), 10, 64)
			if err != nil {
				c.JSON(http.StatusBadRequest, api.ResponseError{
					Code:    http.StatusBadRequest,
					Message: "无效的ID格式",
					Error:   err.Error(),
				})
				return
			}

			// 获取通知配置
			notifyData, err := database.NotifyRepo().GetByID(c.Request.Context(), id)
			if err != nil {
				c.JSON(http.StatusInternalServerError, api.ResponseError{
					Code:    http.StatusInternalServerError,
					Message: "获取通知配置失败",
					Error:   err.Error(),
				})
				return
			}

			if notifyData == nil {
				c.JSON(http.StatusNotFound, api.ResponseError{
					Code:    http.StatusNotFound,
					Message: fmt.Sprintf("通知配置 ID %d 不存在", id),
				})
				return
			}

			responses = append(responses, *notifyData)
		}

		username, _ := c.Get("username")
		log.Infof("Notify config list by user %s from %s", username, c.ClientIP())

		c.JSON(http.StatusOK, api.ResponseSuccess{
			Code:    http.StatusOK,
			Message: "获取成功",
			Data:    responses,
		})
		return
	}

	// 分页查询所有通知配置
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	offset := (page - 1) * pageSize

	// 获取通知配置列表
	notifies, err := database.NotifyRepo().List(c.Request.Context(), offset, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, api.ResponseError{
			Code:    http.StatusInternalServerError,
			Message: "获取通知配置列表失败",
			Error:   err.Error(),
		})
		return
	}

	// 获取总数
	total, err := database.NotifyRepo().Count(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, api.ResponseError{
			Code:    http.StatusInternalServerError,
			Message: "获取通知配置总数失败",
			Error:   err.Error(),
		})
		return
	}

	var responses []notify.Data
	if notifies != nil {
		responses = *notifies
	}

	// 构建分页响应
	result := map[string]any{
		"list":      responses,
		"total":     total,
		"page":      page,
		"page_size": pageSize,
	}

	username, _ := c.Get("username")
	log.Infof("Notify config list by user %s from %s", username, c.ClientIP())

	c.JSON(http.StatusOK, api.ResponseSuccess{
		Code:    http.StatusOK,
		Message: "获取成功",
		Data:    result,
	})
}

// createNotify 创建通知配置
// @Summary 创建通知配置
// @Description 创建单个通知配置
// @Tags 通知管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body notify.CreateRequest true "创建通知配置请求"
// @Success 200 {object} api.ResponseSuccess{data=notify.Data} "创建成功"
// @Failure 400 {object} api.ResponseError "请求参数错误"
// @Failure 401 {object} api.ResponseError "未授权"
// @Failure 500 {object} api.ResponseError "服务器内部错误"
// @Router /api/v1/notify/config [post]
func (h *notifyHandler) createNotify(c *gin.Context) {
	var req notify.CreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, api.ResponseError{
			Code:    http.StatusBadRequest,
			Message: "请求参数错误",
			Error:   err.Error(),
		})
		return
	}

	// 创建通知配置数据模型
	notifyData := &notify.Data{
		Name:        req.Name,
		Description: req.Description,
		Type:        req.Type,
		Config:      req.Config,
	}

	// 处理Enable字段
	if req.Enable != nil {
		notifyData.Enable = *req.Enable
	} else {
		notifyData.Enable = true // 默认启用
	}

	// 创建通知配置
	if err := database.NotifyRepo().Create(c.Request.Context(), notifyData); err != nil {
		c.JSON(http.StatusInternalServerError, api.ResponseError{
			Code:    http.StatusInternalServerError,
			Message: "创建通知配置失败",
			Error:   err.Error(),
		})
		return
	}

	username, _ := c.Get("username")
	log.Infof("Notify config %d created by user %s from %s", notifyData.ID, username, c.ClientIP())

	c.JSON(http.StatusOK, api.ResponseSuccess{
		Code:    http.StatusOK,
		Message: "创建成功",
		Data:    *notifyData,
	})
}

// updateNotify 更新通知配置
// @Summary 更新通知配置
// @Description 根据请求体中的ID更新通知配置信息
// @Tags 通知管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body notify.UpdateRequest true "更新通知配置请求"
// @Success 200 {object} api.ResponseSuccess{data=notify.Data} "更新成功"
// @Failure 400 {object} api.ResponseError "请求参数错误"
// @Failure 401 {object} api.ResponseError "未授权"
// @Failure 404 {object} api.ResponseError "通知配置不存在"
// @Failure 500 {object} api.ResponseError "服务器内部错误"
// @Router /api/v1/notify/config [put]
func (h *notifyHandler) updateNotify(c *gin.Context) {
	var req notify.UpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, api.ResponseError{
			Code:    http.StatusBadRequest,
			Message: "请求参数错误",
			Error:   err.Error(),
		})
		return
	}

	if req.ID == 0 {
		c.JSON(http.StatusBadRequest, api.ResponseError{
			Code:    http.StatusBadRequest,
			Message: "通知配置ID不能为空",
		})
		return
	}

	// 检查通知配置是否存在
	existingNotify, err := database.NotifyRepo().GetByID(c.Request.Context(), req.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, api.ResponseError{
			Code:    http.StatusInternalServerError,
			Message: "获取通知配置失败",
			Error:   err.Error(),
		})
		return
	}

	if existingNotify == nil {
		c.JSON(http.StatusNotFound, api.ResponseError{
			Code:    http.StatusNotFound,
			Message: "通知配置不存在",
		})
		return
	}

	// 更新通知配置数据
	notifyData := &notify.Data{
		ID:          req.ID,
		Name:        req.Name,
		Description: req.Description,
		Type:        req.Type,
		Config:      req.Config,
	}

	// 处理Enable字段
	if req.Enable != nil {
		notifyData.Enable = *req.Enable
	} else {
		notifyData.Enable = existingNotify.Enable
	}

	if err := database.NotifyRepo().Update(c.Request.Context(), notifyData); err != nil {
		c.JSON(http.StatusInternalServerError, api.ResponseError{
			Code:    http.StatusInternalServerError,
			Message: "更新通知配置失败",
			Error:   err.Error(),
		})
		return
	}

	// 获取更新后的通知配置数据
	updatedNotify, err := database.NotifyRepo().GetByID(c.Request.Context(), req.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, api.ResponseError{
			Code:    http.StatusInternalServerError,
			Message: "获取更新后的通知配置失败",
			Error:   err.Error(),
		})
		return
	}

	username, _ := c.Get("username")
	log.Infof("Notify config %d updated by user %s from %s", req.ID, username, c.ClientIP())

	c.JSON(http.StatusOK, api.ResponseSuccess{
		Code:    http.StatusOK,
		Message: "更新成功",
		Data:    *updatedNotify,
	})
}

// deleteNotify 删除通知配置
// @Summary 删除通知配置
// @Description 根据ID删除单个通知配置
// @Tags 通知管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id query int true "通知配置ID"
// @Success 200 {object} api.ResponseSuccess "删除成功"
// @Failure 400 {object} api.ResponseError "请求参数错误"
// @Failure 401 {object} api.ResponseError "未授权"
// @Failure 404 {object} api.ResponseError "通知配置不存在"
// @Failure 500 {object} api.ResponseError "服务器内部错误"
// @Router /api/v1/notify/config [delete]
func (h *notifyHandler) deleteNotify(c *gin.Context) {
	// 获取查询参数中的ID
	idParam := c.Query("id")
	if idParam == "" {
		c.JSON(http.StatusBadRequest, api.ResponseError{
			Code:    http.StatusBadRequest,
			Message: "通知配置ID不能为空",
		})
		return
	}

	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, api.ResponseError{
			Code:    http.StatusBadRequest,
			Message: "无效的ID格式",
			Error:   err.Error(),
		})
		return
	}

	// 检查通知配置是否存在
	existingNotify, err := database.NotifyRepo().GetByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, api.ResponseError{
			Code:    http.StatusInternalServerError,
			Message: "获取通知配置失败",
			Error:   err.Error(),
		})
		return
	}

	if existingNotify == nil {
		c.JSON(http.StatusNotFound, api.ResponseError{
			Code:    http.StatusNotFound,
			Message: "通知配置不存在",
		})
		return
	}

	// 删除通知配置
	if err := database.NotifyRepo().Delete(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, api.ResponseError{
			Code:    http.StatusInternalServerError,
			Message: "删除通知配置失败",
			Error:   err.Error(),
		})
		return
	}

	username, _ := c.Get("username")
	log.Infof("Notify config %d deleted by user %s from %s", id, username, c.ClientIP())

	c.JSON(http.StatusOK, api.ResponseSuccess{
		Code:    http.StatusOK,
		Message: "删除成功",
	})
}

// getTemplates 获取通知模板
// @Summary 获取通知模板
// @Tags 通知管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param page query int false "页码" default(1)
// @Param page_size query int false "每页大小" default(10)
// @Param ids query string false "模板ID列表，逗号分隔"
// @Success 200 {object} api.ResponseSuccess{data=[]notify.Template} "获取成功"
// @Failure 400 {object} api.ResponseError "请求参数错误"
// @Failure 401 {object} api.ResponseError "未授权"
// @Failure 500 {object} api.ResponseError "服务器内部错误"
// @Router /api/v1/notify/template [get]
func (h *notifyHandler) getTemplates(c *gin.Context) {
	// 解析查询参数
	idsParam := c.Query("ids")

	// 如果指定了IDs，则获取指定的通知模板
	if idsParam != "" {
		idStrs := strings.Split(idsParam, ",")
		var responses []notify.Template

		for _, idStr := range idStrs {
			id, err := strconv.ParseInt(strings.TrimSpace(idStr), 10, 64)
			if err != nil {
				c.JSON(http.StatusBadRequest, api.ResponseError{
					Code:    http.StatusBadRequest,
					Message: "无效的ID格式",
					Error:   err.Error(),
				})
				return
			}

			// 获取通知模板
			templateData, err := database.NotifyTemplateRepo().GetByID(c.Request.Context(), id)
			if err != nil {
				c.JSON(http.StatusInternalServerError, api.ResponseError{
					Code:    http.StatusInternalServerError,
					Message: "获取通知模板失败",
					Error:   err.Error(),
				})
				return
			}

			if templateData == nil {
				c.JSON(http.StatusNotFound, api.ResponseError{
					Code:    http.StatusNotFound,
					Message: fmt.Sprintf("通知模板 ID %d 不存在", id),
				})
				return
			}

			responses = append(responses, *templateData)
		}

		username, _ := c.Get("username")
		log.Infof("Notify template list by user %s from %s", username, c.ClientIP())

		c.JSON(http.StatusOK, api.ResponseSuccess{
			Code:    http.StatusOK,
			Message: "获取成功",
			Data:    responses,
		})
		return
	}

	// 分页查询所有通知模板
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	offset := (page - 1) * pageSize

	// 获取通知模板列表
	templates, err := database.NotifyTemplateRepo().List(c.Request.Context(), offset, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, api.ResponseError{
			Code:    http.StatusInternalServerError,
			Message: "获取通知模板列表失败",
			Error:   err.Error(),
		})
		return
	}

	// 获取总数
	total, err := database.NotifyTemplateRepo().Count(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, api.ResponseError{
			Code:    http.StatusInternalServerError,
			Message: "获取通知模板总数失败",
			Error:   err.Error(),
		})
		return
	}

	var responses []notify.Template
	if templates != nil {
		responses = *templates
	}

	// 构建分页响应
	result := map[string]any{
		"list":      responses,
		"total":     total,
		"page":      page,
		"page_size": pageSize,
	}

	username, _ := c.Get("username")
	log.Infof("Notify template list by user %s from %s", username, c.ClientIP())

	c.JSON(http.StatusOK, api.ResponseSuccess{
		Code:    http.StatusOK,
		Message: "获取成功",
		Data:    result,
	})
}

// createTemplate 创建通知模板
// @Summary 创建通知模板
// @Description 创建单个通知模板
// @Tags 通知管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body notify.TemplateCreateRequest true "创建通知模板请求"
// @Success 200 {object} api.ResponseSuccess{data=notify.Template} "创建成功"
// @Failure 400 {object} api.ResponseError "请求参数错误"
// @Failure 401 {object} api.ResponseError "未授权"
// @Failure 500 {object} api.ResponseError "服务器内部错误"
// @Router /api/v1/notify/template [post]
func (h *notifyHandler) createTemplate(c *gin.Context) {
	var req notify.TemplateCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, api.ResponseError{
			Code:    http.StatusBadRequest,
			Message: "请求参数错误",
			Error:   err.Error(),
		})
		return
	}

	// 创建通知模板数据模型
	templateData := &notify.Template{
		Name:        req.Name,
		Description: req.Description,
		Template:    req.Template,
	}

	// 创建通知模板
	if err := database.NotifyTemplateRepo().Create(c.Request.Context(), templateData); err != nil {
		c.JSON(http.StatusInternalServerError, api.ResponseError{
			Code:    http.StatusInternalServerError,
			Message: "创建通知模板失败",
			Error:   err.Error(),
		})
		return
	}

	username, _ := c.Get("username")
	log.Infof("Notify template %d created by user %s from %s", templateData.ID, username, c.ClientIP())

	c.JSON(http.StatusOK, api.ResponseSuccess{
		Code:    http.StatusOK,
		Message: "创建成功",
		Data:    *templateData,
	})
}

// updateTemplate 更新通知模板
// @Summary 更新通知模板
// @Description 根据请求体中的ID更新通知模板信息
// @Tags 通知管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body notify.TemplateUpdateRequest true "更新通知模板请求"
// @Success 200 {object} api.ResponseSuccess{data=notify.Template} "更新成功"
// @Failure 400 {object} api.ResponseError "请求参数错误"
// @Failure 401 {object} api.ResponseError "未授权"
// @Failure 404 {object} api.ResponseError "通知模板不存在"
// @Failure 500 {object} api.ResponseError "服务器内部错误"
// @Router /api/v1/notify/template [put]
func (h *notifyHandler) updateTemplate(c *gin.Context) {
	var req notify.TemplateUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, api.ResponseError{
			Code:    http.StatusBadRequest,
			Message: "请求参数错误",
			Error:   err.Error(),
		})
		return
	}

	if req.ID == 0 {
		c.JSON(http.StatusBadRequest, api.ResponseError{
			Code:    http.StatusBadRequest,
			Message: "通知模板ID不能为空",
		})
		return
	}

	// 检查通知模板是否存在
	existingTemplate, err := database.NotifyTemplateRepo().GetByID(c.Request.Context(), req.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, api.ResponseError{
			Code:    http.StatusInternalServerError,
			Message: "获取通知模板失败",
			Error:   err.Error(),
		})
		return
	}

	if existingTemplate == nil {
		c.JSON(http.StatusNotFound, api.ResponseError{
			Code:    http.StatusNotFound,
			Message: "通知模板不存在",
		})
		return
	}

	// 更新通知模板数据
	templateData := &notify.Template{
		ID:          req.ID,
		Name:        req.Name,
		Description: req.Description,
		Template:    req.Template,
	}

	if err := database.NotifyTemplateRepo().Update(c.Request.Context(), templateData); err != nil {
		c.JSON(http.StatusInternalServerError, api.ResponseError{
			Code:    http.StatusInternalServerError,
			Message: "更新通知模板失败",
			Error:   err.Error(),
		})
		return
	}

	// 获取更新后的通知模板数据
	updatedTemplate, err := database.NotifyTemplateRepo().GetByID(c.Request.Context(), req.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, api.ResponseError{
			Code:    http.StatusInternalServerError,
			Message: "获取更新后的通知模板失败",
			Error:   err.Error(),
		})
		return
	}

	username, _ := c.Get("username")
	log.Infof("Notify template %d updated by user %s from %s", req.ID, username, c.ClientIP())

	c.JSON(http.StatusOK, api.ResponseSuccess{
		Code:    http.StatusOK,
		Message: "更新成功",
		Data:    *updatedTemplate,
	})
}

// deleteTemplate 删除通知模板
// @Summary 删除通知模板
// @Description 根据ID删除单个通知模板
// @Tags 通知管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id query int true "通知模板ID"
// @Success 200 {object} api.ResponseSuccess "删除成功"
// @Failure 400 {object} api.ResponseError "请求参数错误"
// @Failure 401 {object} api.ResponseError "未授权"
// @Failure 404 {object} api.ResponseError "通知模板不存在"
// @Failure 500 {object} api.ResponseError "服务器内部错误"
// @Router /api/v1/notify/template [delete]
func (h *notifyHandler) deleteTemplate(c *gin.Context) {
	// 获取查询参数中的ID
	idParam := c.Query("id")
	if idParam == "" {
		c.JSON(http.StatusBadRequest, api.ResponseError{
			Code:    http.StatusBadRequest,
			Message: "通知模板ID不能为空",
		})
		return
	}

	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, api.ResponseError{
			Code:    http.StatusBadRequest,
			Message: "无效的ID格式",
			Error:   err.Error(),
		})
		return
	}

	// 检查通知模板是否存在
	existingTemplate, err := database.NotifyTemplateRepo().GetByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, api.ResponseError{
			Code:    http.StatusInternalServerError,
			Message: "获取通知模板失败",
			Error:   err.Error(),
		})
		return
	}

	if existingTemplate == nil {
		c.JSON(http.StatusNotFound, api.ResponseError{
			Code:    http.StatusNotFound,
			Message: "通知模板不存在",
		})
		return
	}

	// 删除通知模板
	if err := database.NotifyTemplateRepo().Delete(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, api.ResponseError{
			Code:    http.StatusInternalServerError,
			Message: "删除通知模板失败",
			Error:   err.Error(),
		})
		return
	}

	username, _ := c.Get("username")
	log.Infof("Notify template %d deleted by user %s from %s", id, username, c.ClientIP())

	c.JSON(http.StatusOK, api.ResponseSuccess{
		Code:    http.StatusOK,
		Message: "删除成功",
	})
}
