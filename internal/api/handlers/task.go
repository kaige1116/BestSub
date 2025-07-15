package handlers

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/bestruirui/bestsub/internal/api/middleware"
	"github.com/bestruirui/bestsub/internal/api/router"
	"github.com/bestruirui/bestsub/internal/core/task"
	"github.com/bestruirui/bestsub/internal/models/api"
	taskModel "github.com/bestruirui/bestsub/internal/models/task"
	"github.com/bestruirui/bestsub/internal/utils/log"
	"github.com/gin-gonic/gin"
)

// taskHandler 任务处理器
type taskHandler struct{}

// init 函数用于自动注册路由
func init() {
	h := newTaskHandler()

	// 需要认证的任务路由
	router.NewGroupRouter("/api/v1/tasks").
		Use(middleware.Auth()).
		AddRoute(
			router.NewRoute("", router.POST).
				Handle(h.createTask).
				WithDescription("Create task"),
		).
		AddRoute(
			router.NewRoute("", router.GET).
				Handle(h.getTasks).
				WithDescription("Get tasks or list all with pagination"),
		).
		AddRoute(
			router.NewRoute("/:id", router.GET).
				Handle(h.getTask).
				WithDescription("Get task by ID"),
		).
		AddRoute(
			router.NewRoute("", router.PATCH).
				Handle(h.updateTask).
				WithDescription("Update task"),
		).
		AddRoute(
			router.NewRoute("/:id", router.DELETE).
				Handle(h.deleteTask).
				WithDescription("Delete task"),
		).
		AddRoute(
			router.NewRoute("/:id/run", router.POST).
				Handle(h.runTask).
				WithDescription("Run task manually"),
		).
		AddRoute(
			router.NewRoute("/:id/stop", router.POST).
				Handle(h.stopTask).
				WithDescription("Stop running task"),
		)
}

// newTaskHandler 创建任务处理器
func newTaskHandler() *taskHandler {
	return &taskHandler{}
}

// createTask 创建任务
// @Summary 创建任务
// @Description 创建单个任务
// @Tags 任务管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body taskModel.CreateRequest true "创建任务请求"
// @Success 200 {object} api.ResponseSuccess{data=taskModel.Data} "创建成功"
// @Failure 400 {object} api.ResponseError "请求参数错误"
// @Failure 401 {object} api.ResponseError "未授权"
// @Failure 500 {object} api.ResponseError "服务器内部错误"
// @Router /api/v1/tasks [post]
func (h *taskHandler) createTask(c *gin.Context) {
	var req taskModel.CreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, api.ResponseError{
			Code:    http.StatusBadRequest,
			Message: "请求参数错误",
			Error:   err.Error(),
		})
		return
	}

	// 创建任务
	taskData, err := task.AddTask(&req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, api.ResponseError{
			Code:    http.StatusInternalServerError,
			Message: "创建任务失败",
			Error:   err.Error(),
		})
		return
	}

	username, _ := c.Get("username")
	log.Infof("Task %d created by user %s from %s", taskData.ID, username, c.ClientIP())

	c.JSON(http.StatusOK, api.ResponseSuccess{
		Code:    http.StatusOK,
		Message: "创建成功",
		Data:    taskData,
	})
}

// getTasks 获取任务列表
// @Summary 获取任务列表
// @Tags 任务管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param page query int false "页码" default(1)
// @Param page_size query int false "每页大小" default(10)
// @Param ids query string false "任务ID列表，逗号分隔"
// @Success 200 {object} api.ResponseSuccess{data=[]taskModel.Data} "获取成功"
// @Failure 400 {object} api.ResponseError "请求参数错误"
// @Failure 401 {object} api.ResponseError "未授权"
// @Failure 500 {object} api.ResponseError "服务器内部错误"
// @Router /api/v1/tasks [get]
func (h *taskHandler) getTasks(c *gin.Context) {
	// 解析查询参数
	idsParam := c.Query("ids")

	// 如果指定了IDs，则获取指定的任务
	if idsParam != "" {
		idStrs := strings.Split(idsParam, ",")
		var tasks []taskModel.Data

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

			// 获取任务
			taskData, err := task.GetTask(id)
			if err != nil {
				c.JSON(http.StatusInternalServerError, api.ResponseError{
					Code:    http.StatusInternalServerError,
					Message: "获取任务失败",
					Error:   err.Error(),
				})
				return
			}

			if taskData == nil {
				c.JSON(http.StatusNotFound, api.ResponseError{
					Code:    http.StatusNotFound,
					Message: fmt.Sprintf("任务 ID %d 不存在", id),
				})
				return
			}

			tasks = append(tasks, *taskData)
		}

		c.JSON(http.StatusOK, api.ResponseSuccess{
			Code:    http.StatusOK,
			Message: "获取成功",
			Data:    tasks,
		})
		return
	}

	// 分页查询所有任务
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	offset := (page - 1) * pageSize

	// 获取任务列表
	tasks, total, err := task.ListTasks(offset, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, api.ResponseError{
			Code:    http.StatusInternalServerError,
			Message: "获取任务列表失败",
			Error:   err.Error(),
		})
		return
	}

	// 构建分页响应
	result := map[string]interface{}{
		"list":      tasks,
		"total":     total,
		"page":      page,
		"page_size": pageSize,
	}

	username, _ := c.Get("username")
	log.Infof("Task list requested by user %s from %s", username, c.ClientIP())

	c.JSON(http.StatusOK, api.ResponseSuccess{
		Code:    http.StatusOK,
		Message: "获取成功",
		Data:    result,
	})
}

// getTask 获取单个任务
// @Summary 获取任务详情
// @Tags 任务管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "任务ID"
// @Success 200 {object} api.ResponseSuccess{data=taskModel.Data} "获取成功"
// @Failure 400 {object} api.ResponseError "请求参数错误"
// @Failure 401 {object} api.ResponseError "未授权"
// @Failure 404 {object} api.ResponseError "任务不存在"
// @Failure 500 {object} api.ResponseError "服务器内部错误"
// @Router /api/v1/tasks/{id} [get]
func (h *taskHandler) getTask(c *gin.Context) {
	// 获取路径参数中的ID
	idParam := c.Param("id")
	if idParam == "" {
		c.JSON(http.StatusBadRequest, api.ResponseError{
			Code:    http.StatusBadRequest,
			Message: "任务ID不能为空",
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

	// 获取任务
	taskData, err := task.GetTask(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, api.ResponseError{
			Code:    http.StatusInternalServerError,
			Message: "获取任务失败",
			Error:   err.Error(),
		})
		return
	}

	if taskData == nil {
		c.JSON(http.StatusNotFound, api.ResponseError{
			Code:    http.StatusNotFound,
			Message: "任务不存在",
		})
		return
	}

	username, _ := c.Get("username")
	log.Infof("Task %d requested by user %s from %s", id, username, c.ClientIP())

	c.JSON(http.StatusOK, api.ResponseSuccess{
		Code:    http.StatusOK,
		Message: "获取成功",
		Data:    taskData,
	})
}

// updateTask 更新任务
// @Summary 更新任务
// @Description 根据请求体中的ID更新任务信息
// @Tags 任务管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body taskModel.UpdateRequest true "更新任务请求"
// @Success 200 {object} api.ResponseSuccess{data=taskModel.Data} "更新成功"
// @Failure 400 {object} api.ResponseError "请求参数错误"
// @Failure 401 {object} api.ResponseError "未授权"
// @Failure 404 {object} api.ResponseError "任务不存在"
// @Failure 500 {object} api.ResponseError "服务器内部错误"
// @Router /api/v1/tasks [patch]
func (h *taskHandler) updateTask(c *gin.Context) {
	var req taskModel.UpdateRequest
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
			Message: "任务ID不能为空",
		})
		return
	}

	// 更新任务
	taskData, err := task.UpdateTask(&req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, api.ResponseError{
			Code:    http.StatusInternalServerError,
			Message: "更新任务失败",
			Error:   err.Error(),
		})
		return
	}

	username, _ := c.Get("username")
	log.Infof("Task %d updated by user %s from %s", req.ID, username, c.ClientIP())

	c.JSON(http.StatusOK, api.ResponseSuccess{
		Code:    http.StatusOK,
		Message: "更新成功",
		Data:    taskData,
	})
}

// deleteTask 删除任务
// @Summary 删除任务
// @Description 根据ID删除单个任务
// @Tags 任务管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "任务ID"
// @Success 200 {object} api.ResponseSuccess "删除成功"
// @Failure 400 {object} api.ResponseError "请求参数错误"
// @Failure 401 {object} api.ResponseError "未授权"
// @Failure 404 {object} api.ResponseError "任务不存在"
// @Failure 500 {object} api.ResponseError "服务器内部错误"
// @Router /api/v1/tasks/{id} [delete]
func (h *taskHandler) deleteTask(c *gin.Context) {
	// 获取路径参数中的ID
	idParam := c.Param("id")
	if idParam == "" {
		c.JSON(http.StatusBadRequest, api.ResponseError{
			Code:    http.StatusBadRequest,
			Message: "任务ID不能为空",
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

	// 检查任务是否存在
	existingTask, err := task.GetTask(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, api.ResponseError{
			Code:    http.StatusInternalServerError,
			Message: "获取任务失败",
			Error:   err.Error(),
		})
		return
	}

	if existingTask == nil {
		c.JSON(http.StatusNotFound, api.ResponseError{
			Code:    http.StatusNotFound,
			Message: "任务不存在",
		})
		return
	}

	// 删除任务
	if err := task.RemoveTaskWithDb(id); err != nil {
		c.JSON(http.StatusInternalServerError, api.ResponseError{
			Code:    http.StatusInternalServerError,
			Message: "删除任务失败",
			Error:   err.Error(),
		})
		return
	}

	username, _ := c.Get("username")
	log.Infof("Task %d deleted by user %s from %s", id, username, c.ClientIP())

	c.JSON(http.StatusOK, api.ResponseSuccess{
		Code:    http.StatusOK,
		Message: "删除成功",
	})
}

// runTask 手动运行任务
// @Summary 手动运行任务
// @Description 手动触发任务执行
// @Tags 任务管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "任务ID"
// @Success 200 {object} api.ResponseSuccess "运行成功"
// @Failure 400 {object} api.ResponseError "请求参数错误"
// @Failure 401 {object} api.ResponseError "未授权"
// @Failure 404 {object} api.ResponseError "任务不存在"
// @Failure 500 {object} api.ResponseError "服务器内部错误"
// @Router /api/v1/tasks/{id}/run [post]
func (h *taskHandler) runTask(c *gin.Context) {
	// 获取路径参数中的ID
	idParam := c.Param("id")
	if idParam == "" {
		c.JSON(http.StatusBadRequest, api.ResponseError{
			Code:    http.StatusBadRequest,
			Message: "任务ID不能为空",
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

	// 检查任务是否存在
	existingTask, err := task.GetTask(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, api.ResponseError{
			Code:    http.StatusInternalServerError,
			Message: "获取任务失败",
			Error:   err.Error(),
		})
		return
	}

	if existingTask == nil {
		c.JSON(http.StatusNotFound, api.ResponseError{
			Code:    http.StatusNotFound,
			Message: "任务不存在",
		})
		return
	}

	// 运行任务
	if err := task.RunTask(id); err != nil {
		c.JSON(http.StatusInternalServerError, api.ResponseError{
			Code:    http.StatusInternalServerError,
			Message: "运行任务失败",
			Error:   err.Error(),
		})
		return
	}

	username, _ := c.Get("username")
	log.Infof("Task %d manually run by user %s from %s", id, username, c.ClientIP())

	c.JSON(http.StatusOK, api.ResponseSuccess{
		Code:    http.StatusOK,
		Message: "任务已开始运行",
	})
}

// stopTask 停止任务
// @Summary 停止任务
// @Description 停止正在运行的任务
// @Tags 任务管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "任务ID"
// @Success 200 {object} api.ResponseSuccess "停止成功"
// @Failure 400 {object} api.ResponseError "请求参数错误"
// @Failure 401 {object} api.ResponseError "未授权"
// @Failure 404 {object} api.ResponseError "任务不存在"
// @Failure 500 {object} api.ResponseError "服务器内部错误"
// @Router /api/v1/tasks/{id}/stop [post]
func (h *taskHandler) stopTask(c *gin.Context) {
	// 获取路径参数中的ID
	idParam := c.Param("id")
	if idParam == "" {
		c.JSON(http.StatusBadRequest, api.ResponseError{
			Code:    http.StatusBadRequest,
			Message: "任务ID不能为空",
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

	// 检查任务是否存在
	existingTask, err := task.GetTask(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, api.ResponseError{
			Code:    http.StatusInternalServerError,
			Message: "获取任务失败",
			Error:   err.Error(),
		})
		return
	}

	if existingTask == nil {
		c.JSON(http.StatusNotFound, api.ResponseError{
			Code:    http.StatusNotFound,
			Message: "任务不存在",
		})
		return
	}

	// 停止任务
	if err := task.StopTask(id); err != nil {
		c.JSON(http.StatusInternalServerError, api.ResponseError{
			Code:    http.StatusInternalServerError,
			Message: "停止任务失败",
			Error:   err.Error(),
		})
		return
	}

	username, _ := c.Get("username")
	log.Infof("Task %d stopped by user %s from %s", id, username, c.ClientIP())

	c.JSON(http.StatusOK, api.ResponseSuccess{
		Code:    http.StatusOK,
		Message: "任务已停止",
	})
}
