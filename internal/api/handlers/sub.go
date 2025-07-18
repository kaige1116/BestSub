package handlers

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/bestruirui/bestsub/internal/api/common"
	"github.com/bestruirui/bestsub/internal/api/middleware"
	"github.com/bestruirui/bestsub/internal/api/router"
	taskcore "github.com/bestruirui/bestsub/internal/core/task"
	"github.com/bestruirui/bestsub/internal/database/op"
	dbc "github.com/bestruirui/bestsub/internal/models/common"
	"github.com/bestruirui/bestsub/internal/models/sub"
	"github.com/bestruirui/bestsub/internal/models/task"
	"github.com/bestruirui/bestsub/internal/utils/log"
	"github.com/gin-gonic/gin"
)

// subLinkHandler 订阅链接处理器
type subLinkHandler struct{}

// init 函数用于自动注册路由
func init() {
	h := newSubLinkHandler()

	// 需要认证的订阅链接路由
	router.NewGroupRouter("/api/v1/sub").
		Use(middleware.Auth()).
		AddRoute(
			router.NewRoute("", router.POST).
				Handle(h.createSub).
				WithDescription("Create subscription links"),
		).
		AddRoute(
			router.NewRoute("", router.GET).
				Handle(h.getSubs).
				WithDescription("Get subscription links or list all with pagination"),
		).
		AddRoute(
			router.NewRoute("", router.PATCH).
				Handle(h.updateSub).
				WithDescription("Update subscription links"),
		).
		AddRoute(
			router.NewRoute("/:id", router.DELETE).
				Handle(h.deleteSub).
				WithDescription("Delete subscription link"),
		)
}

// newSubLinkHandler 创建订阅链接处理器
func newSubLinkHandler() *subLinkHandler {
	return &subLinkHandler{}
}

// createSub 创建订阅链接
// @Summary 创建订阅链接
// @Description 创建单个订阅链接
// @Tags 订阅管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body sub.CreateRequest true "创建订阅链接请求"
// @Success 200 {object} common.ResponseSuccessStruct{data=sub.Response} "创建成功"
// @Failure 400 {object} common.ResponseErrorStruct "请求参数错误"
// @Failure 401 {object} common.ResponseErrorStruct "未授权"
// @Failure 500 {object} common.ResponseErrorStruct "服务器内部错误"
// @Router /api/v1/sub [post]
func (h *subLinkHandler) createSub(c *gin.Context) {
	var req sub.CreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		common.ResponseError(c, http.StatusBadRequest, err)
		return
	}

	// 创建订阅链接数据模型
	subData := &sub.Data{
		BaseDbModel: dbc.BaseDbModel{
			Name:        req.Name,
			Description: req.Description,
		},
		URL: req.URL,
	}

	// 创建订阅链接
	if err := op.SubRepo().Create(c.Request.Context(), subData); err != nil {
		common.ResponseError(c, http.StatusInternalServerError, err)
		return
	}

	// 创建关联的任务
	var createdTasks []task.Data
	for _, taskReq := range req.Task {

		// 创建任务
		taskData, err := taskcore.AddTask(&taskReq)
		if err != nil {
			common.ResponseError(c, http.StatusInternalServerError, err)
			return
		}

		// 建立订阅与任务的关联
		log.Debugf("建立订阅与任务的关联: %d, %d", subData.ID, taskData.ID)
		if err := op.SubRepo().AddTaskRelation(c.Request.Context(), subData.ID, taskData.ID); err != nil {
			common.ResponseError(c, http.StatusInternalServerError, err)
			return
		}

		createdTasks = append(createdTasks, *taskData)
	}

	// 构建响应
	response := sub.Response{
		Data: *subData,
		Task: createdTasks,
	}

	username, _ := c.Get("username")
	log.Infof("Subscription link %d created by user %s from %s", subData.ID, username, c.ClientIP())

	common.ResponseSuccess(c, response)
}

// getSubs 获取订阅链接
// @Summary 获取订阅链接
// @Tags 订阅管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param page query int false "页码" default(1)
// @Param page_size query int false "每页大小" default(10)
// @Param ids query string false "链接ID列表，逗号分隔"
// @Success 200 {object} common.ResponseSuccessStruct{data=[]sub.Response} "获取成功"
// @Failure 400 {object} common.ResponseErrorStruct "请求参数错误"
// @Failure 401 {object} common.ResponseErrorStruct "未授权"
// @Failure 500 {object} common.ResponseErrorStruct "服务器内部错误"
// @Router /api/v1/sub [get]
func (h *subLinkHandler) getSubs(c *gin.Context) {
	// 解析查询参数
	idsParam := c.Query("ids")

	// 如果指定了IDs，则获取指定的订阅链接
	if idsParam != "" {
		idStrs := strings.Split(idsParam, ",")
		var responses []sub.Response

		for _, idStr := range idStrs {
			id, err := strconv.ParseInt(strings.TrimSpace(idStr), 10, 64)
			if err != nil {
				common.ResponseError(c, http.StatusBadRequest, err)
				return
			}

			// 获取订阅链接
			subData, err := op.SubRepo().GetByID(c.Request.Context(), id)
			if err != nil {
				common.ResponseError(c, http.StatusInternalServerError, err)
				return
			}

			if subData == nil {
				common.ResponseError(c, http.StatusNotFound, fmt.Errorf("订阅链接 ID %d 不存在", id))
				return
			}

			// 获取关联的任务
			tasks, err := op.TaskRepo().GetBySubID(c.Request.Context(), id)
			if err != nil {
				common.ResponseError(c, http.StatusInternalServerError, err)
				return
			}

			var taskList []task.Data
			if tasks != nil {
				taskList = *tasks
			}

			responses = append(responses, sub.Response{
				Data: *subData,
				Task: taskList,
			})
		}

		common.ResponseSuccess(c, responses)
		return
	}

	// 分页查询所有订阅链接
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	offset := (page - 1) * pageSize

	// 获取订阅链接列表
	subs, err := op.SubRepo().List(c.Request.Context(), offset, pageSize)
	if err != nil {
		common.ResponseError(c, http.StatusInternalServerError, err)
		return
	}

	// 获取总数
	total, err := op.SubRepo().Count(c.Request.Context())
	if err != nil {
		common.ResponseError(c, http.StatusInternalServerError, err)
		return
	}

	// 构建响应列表
	var responses []sub.Response
	if subs != nil {
		for _, subData := range *subs {
			// 获取每个订阅的关联任务
			tasks, err := op.TaskRepo().GetBySubID(c.Request.Context(), subData.ID)
			if err != nil {
				common.ResponseError(c, http.StatusInternalServerError, err)
				return
			}

			var taskList []task.Data
			if tasks != nil {
				taskList = *tasks
			}

			responses = append(responses, sub.Response{
				Data: subData,
				Task: taskList,
			})
		}
	}

	// 构建分页响应
	result := map[string]interface{}{
		"list":      responses,
		"total":     total,
		"page":      page,
		"page_size": pageSize,
	}

	username, _ := c.Get("username")
	log.Infof("Subscription link list by user %s from %s", username, c.ClientIP())

	common.ResponseSuccess(c, result)
}

// updateSub 更新订阅链接
// @Summary 更新订阅链接
// @Description 根据请求体中的ID更新订阅链接信息
// @Tags 订阅管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body sub.UpdateRequest true "更新订阅链接请求"
// @Success 200 {object} common.ResponseSuccessStruct{data=sub.Response} "更新成功"
// @Failure 400 {object} common.ResponseErrorStruct "请求参数错误"
// @Failure 401 {object} common.ResponseErrorStruct "未授权"
// @Failure 404 {object} common.ResponseErrorStruct "订阅链接不存在"
// @Failure 500 {object} common.ResponseErrorStruct "服务器内部错误"
// @Router /api/v1/sub [patch]
func (h *subLinkHandler) updateSub(c *gin.Context) {
	var req sub.UpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		common.ResponseError(c, http.StatusBadRequest, err)
		return
	}

	if req.ID == 0 {
		common.ResponseError(c, http.StatusBadRequest, fmt.Errorf("订阅链接ID不能为空"))
		return
	}

	// 检查订阅链接是否存在
	existingSub, err := op.SubRepo().GetByID(c.Request.Context(), req.ID)
	if err != nil {
		common.ResponseError(c, http.StatusInternalServerError, err)
		return
	}

	if existingSub == nil {
		common.ResponseError(c, http.StatusNotFound, fmt.Errorf("订阅链接不存在"))
		return
	}

	// 更新订阅链接数据
	subData := &sub.Data{
		BaseDbModel: dbc.BaseDbModel{
			ID:          req.ID,
			Name:        req.Name,
			Description: req.Description,
		},
		URL: req.URL,
	}

	if err := op.SubRepo().Update(c.Request.Context(), subData); err != nil {
		common.ResponseError(c, http.StatusInternalServerError, err)
		return
	}

	// 处理任务更新
	var updatedTasks []task.Data
	for _, taskReq := range req.Task {
		taskData, err := taskcore.UpdateTask(&taskReq)
		if err != nil {
			common.ResponseError(c, http.StatusInternalServerError, err)
			return
		}

		updatedTasks = append(updatedTasks, *taskData)
	}

	// 获取更新后的订阅链接数据
	updatedSub, err := op.SubRepo().GetByID(c.Request.Context(), req.ID)
	if err != nil {
		common.ResponseError(c, http.StatusInternalServerError, err)
		return
	}

	// 构建响应
	response := sub.Response{
		Data: *updatedSub,
		Task: updatedTasks,
	}
	username, _ := c.Get("username")
	log.Infof("Subscription link %d updated by user %s from %s", req.ID, username, c.ClientIP())

	common.ResponseSuccess(c, response)
}

// deleteSub 删除订阅链接
// @Summary 删除订阅链接
// @Description 根据ID删除单个订阅链接
// @Tags 订阅管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "订阅链接ID"
// @Success 200 {object} common.ResponseSuccessStruct "删除成功"
// @Failure 400 {object} common.ResponseErrorStruct "请求参数错误"
// @Failure 401 {object} common.ResponseErrorStruct "未授权"
// @Failure 404 {object} common.ResponseErrorStruct "订阅链接不存在"
// @Failure 500 {object} common.ResponseErrorStruct "服务器内部错误"
// @Router /api/v1/sub/{id} [delete]
func (h *subLinkHandler) deleteSub(c *gin.Context) {
	// 获取路径参数中的ID
	idParam := c.Param("id")
	if idParam == "" {
		common.ResponseError(c, http.StatusBadRequest, fmt.Errorf("订阅链接ID不能为空"))
		return
	}

	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		common.ResponseError(c, http.StatusBadRequest, err)
		return
	}

	// 检查订阅链接是否存在
	existingSub, err := op.SubRepo().GetByID(c.Request.Context(), id)
	if err != nil {
		common.ResponseError(c, http.StatusInternalServerError, err)
		return
	}

	if existingSub == nil {
		common.ResponseError(c, http.StatusNotFound, fmt.Errorf("订阅链接不存在"))
		return
	}

	taskIDs, err := op.TaskRepo().GetTaskIDsBySubID(c.Request.Context(), id)
	if err != nil {
		common.ResponseError(c, http.StatusInternalServerError, err)
		return
	}

	for _, taskID := range taskIDs {
		if err := taskcore.RemoveTask(taskID); err != nil {
			common.ResponseError(c, http.StatusInternalServerError, err)
			return
		}
	}

	// 删除订阅链接（数据库触发器会自动删除关联的任务）
	if err := op.SubRepo().Delete(c.Request.Context(), id); err != nil {
		common.ResponseError(c, http.StatusInternalServerError, err)
		return
	}
	username, _ := c.Get("username")
	log.Infof("Subscription link %d deleted by user %s from %s", id, username, c.ClientIP())

	common.ResponseSuccess(c, nil)

}
