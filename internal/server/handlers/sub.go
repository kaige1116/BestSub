package handlers

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/bestruirui/bestsub/internal/server/middleware"
	"github.com/bestruirui/bestsub/internal/server/resp"
	"github.com/bestruirui/bestsub/internal/server/router"
	taskcore "github.com/bestruirui/bestsub/internal/core/task"
	"github.com/bestruirui/bestsub/internal/database/op"
	dbc "github.com/bestruirui/bestsub/internal/models/common"
	"github.com/bestruirui/bestsub/internal/models/sub"
	"github.com/bestruirui/bestsub/internal/models/task"
	"github.com/bestruirui/bestsub/internal/utils/log"
	"github.com/gin-gonic/gin"
)

// init 函数用于自动注册路由
func init() {
	// 需要认证的订阅链接路由
	router.NewGroupRouter("/api/v1/sub").
		Use(middleware.Auth()).
		AddRoute(
			router.NewRoute("", router.POST).
				Handle(createSub),
		).
		AddRoute(
			router.NewRoute("", router.GET).
				Handle(getSubs),
		).
		AddRoute(
			router.NewRoute("", router.PATCH).
				Handle(updateSub),
		).
		AddRoute(
			router.NewRoute("/:id", router.DELETE).
				Handle(deleteSub),
		)
}

// createSub 创建订阅链接
// @Summary 创建订阅链接
// @Description 创建单个订阅链接
// @Tags 订阅管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body sub.CreateRequest true "创建订阅链接请求"
// @Success 200 {object} resp.SuccessStruct{data=sub.Response} "创建成功"
// @Failure 400 {object} resp.ErrorStruct "请求参数错误"
// @Failure 401 {object} resp.ErrorStruct "未授权"
// @Failure 500 {object} resp.ErrorStruct "服务器内部错误"
// @Router /api/v1/sub [post]
func createSub(c *gin.Context) {
	var req sub.CreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		resp.ErrorBadRequest(c)
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
		resp.Error(c, http.StatusInternalServerError, "failed to create subscription link")
		return
	}

	// 创建关联的任务
	var createdTasks []task.Data
	for _, taskReq := range req.Task {

		// 创建任务
		taskData, err := taskcore.AddTask(&taskReq)
		if err != nil {
			resp.Error(c, http.StatusInternalServerError, "failed to create task")
			return
		}

		// 建立订阅与任务的关联
		log.Debugf("建立订阅与任务的关联: %d, %d", subData.ID, taskData.ID)
		if err := op.SubRepo().AddTaskRelation(c.Request.Context(), subData.ID, taskData.ID); err != nil {
			resp.Error(c, http.StatusInternalServerError, "failed to add task relation")
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

	resp.Success(c, response)
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
// @Success 200 {object} resp.SuccessStruct{data=[]sub.Response} "获取成功"
// @Failure 400 {object} resp.ErrorStruct "请求参数错误"
// @Failure 401 {object} resp.ErrorStruct "未授权"
// @Failure 500 {object} resp.ErrorStruct "服务器内部错误"
// @Router /api/v1/sub [get]
func getSubs(c *gin.Context) {
	// 解析查询参数
	idsParam := c.Query("ids")

	// 如果指定了IDs，则获取指定的订阅链接
	if idsParam != "" {
		idStrs := strings.Split(idsParam, ",")
		var responses []sub.Response

		for _, idStr := range idStrs {
			id, err := strconv.ParseUint(strings.TrimSpace(idStr), 10, 16)
			if err != nil {
				resp.Error(c, http.StatusBadRequest, "invalid id")
				return
			}

			// 获取订阅链接
			subData, err := op.SubRepo().GetByID(c.Request.Context(), uint16(id))
			if err != nil {
				resp.Error(c, http.StatusInternalServerError, "failed to get subscription link")
				return
			}

			if subData == nil {
				resp.Error(c, http.StatusNotFound, fmt.Sprintf("subscription link %d not found", id))
				return
			}

			// 获取关联的任务
			tasks, err := op.TaskRepo().GetBySubID(c.Request.Context(), uint16(id))
			if err != nil {
				resp.Error(c, http.StatusInternalServerError, "failed to get tasks")
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

		resp.Success(c, responses)
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
		resp.Error(c, http.StatusInternalServerError, "failed to get subscription links")
		return
	}

	// 获取总数
	total, err := op.SubRepo().Count(c.Request.Context())
	if err != nil {
		resp.Error(c, http.StatusInternalServerError, "failed to get total count")
		return
	}

	// 构建响应列表
	var responses []sub.Response
	if subs != nil {
		for _, subData := range *subs {
			// 获取每个订阅的关联任务
			tasks, err := op.TaskRepo().GetBySubID(c.Request.Context(), subData.ID)
			if err != nil {
				resp.Error(c, http.StatusInternalServerError, "failed to get tasks")
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

	resp.Success(c, result)
}

// updateSub 更新订阅链接
// @Summary 更新订阅链接
// @Description 根据请求体中的ID更新订阅链接信息
// @Tags 订阅管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body sub.UpdateRequest true "更新订阅链接请求"
// @Success 200 {object} resp.SuccessStruct{data=sub.Response} "更新成功"
// @Failure 400 {object} resp.ErrorStruct "请求参数错误"
// @Failure 401 {object} resp.ErrorStruct "未授权"
// @Failure 404 {object} resp.ErrorStruct "订阅链接不存在"
// @Failure 500 {object} resp.ErrorStruct "服务器内部错误"
// @Router /api/v1/sub [patch]
func updateSub(c *gin.Context) {
	var req sub.UpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		resp.ErrorBadRequest(c)
		return
	}

	if req.ID == 0 {
		resp.Error(c, http.StatusBadRequest, "subscription link id is required")
		return
	}

	// 检查订阅链接是否存在
	existingSub, err := op.SubRepo().GetByID(c.Request.Context(), req.ID)
	if err != nil {
		resp.Error(c, http.StatusInternalServerError, "failed to get subscription link")
		return
	}

	if existingSub == nil {
		resp.Error(c, http.StatusNotFound, "subscription link not found")
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
		resp.Error(c, http.StatusInternalServerError, "failed to update subscription link")
		return
	}

	// 处理任务更新
	var updatedTasks []task.Data
	for _, taskReq := range req.Task {
		taskData, err := taskcore.UpdateTask(&taskReq)
		if err != nil {
			resp.Error(c, http.StatusInternalServerError, "failed to update task")
			return
		}

		updatedTasks = append(updatedTasks, *taskData)
	}

	// 获取更新后的订阅链接数据
	updatedSub, err := op.SubRepo().GetByID(c.Request.Context(), req.ID)
	if err != nil {
		resp.Error(c, http.StatusInternalServerError, "failed to get updated subscription link")
		return
	}

	// 构建响应
	response := sub.Response{
		Data: *updatedSub,
		Task: updatedTasks,
	}
	username, _ := c.Get("username")
	log.Infof("Subscription link %d updated by user %s from %s", req.ID, username, c.ClientIP())

	resp.Success(c, response)
}

// deleteSub 删除订阅链接
// @Summary 删除订阅链接
// @Description 根据ID删除单个订阅链接
// @Tags 订阅管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "订阅链接ID"
// @Success 200 {object} resp.SuccessStruct "删除成功"
// @Failure 400 {object} resp.ErrorStruct "请求参数错误"
// @Failure 401 {object} resp.ErrorStruct "未授权"
// @Failure 404 {object} resp.ErrorStruct "订阅链接不存在"
// @Failure 500 {object} resp.ErrorStruct "服务器内部错误"
// @Router /api/v1/sub/{id} [delete]
func deleteSub(c *gin.Context) {
	// 获取路径参数中的ID
	idParam := c.Param("id")
	if idParam == "" {
		resp.Error(c, http.StatusBadRequest, "subscription link id is required")
		return
	}

	id, err := strconv.ParseUint(idParam, 10, 16)
	if err != nil {
		resp.ErrorBadRequest(c)
		return
	}

	// 检查订阅链接是否存在
	existingSub, err := op.SubRepo().GetByID(c.Request.Context(), uint16(id))
	if err != nil {
		resp.Error(c, http.StatusInternalServerError, "failed to get subscription link")
		return
	}

	if existingSub == nil {
		resp.Error(c, http.StatusNotFound, "subscription link not found")
		return
	}

	taskIDs, err := op.TaskRepo().GetTaskIDsBySubID(c.Request.Context(), uint16(id))
	if err != nil {
		resp.Error(c, http.StatusInternalServerError, "failed to get task ids")
		return
	}

	for _, taskID := range taskIDs {
		if err := taskcore.RemoveTask(taskID); err != nil {
			resp.Error(c, http.StatusInternalServerError, "failed to remove task")
			return
		}
	}

	if err := op.SubRepo().Delete(c.Request.Context(), uint16(id)); err != nil {
		resp.Error(c, http.StatusInternalServerError, "failed to delete subscription link")
		return
	}
	username, _ := c.Get("username")
	log.Infof("Subscription link %d deleted by user %s from %s", id, username, c.ClientIP())

	resp.Success(c, nil)

}
