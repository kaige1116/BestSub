package handlers

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/bestruirui/bestsub/internal/api/middleware"
	"github.com/bestruirui/bestsub/internal/api/resp"
	"github.com/bestruirui/bestsub/internal/api/router"
	"github.com/bestruirui/bestsub/internal/core/task"
	taskModel "github.com/bestruirui/bestsub/internal/models/task"
	"github.com/bestruirui/bestsub/internal/utils/log"
	"github.com/gin-gonic/gin"
)

// init 函数用于自动注册路由
func init() {

	// 需要认证的任务路由
	router.NewGroupRouter("/api/v1/tasks").
		Use(middleware.Auth()).
		AddRoute(
			router.NewRoute("", router.POST).
				Handle(createTask).
				WithDescription("Create task"),
		).
		AddRoute(
			router.NewRoute("", router.GET).
				Handle(getTasks).
				WithDescription("Get tasks or list all with pagination"),
		).
		AddRoute(
			router.NewRoute("/:id", router.GET).
				Handle(getTask).
				WithDescription("Get task by ID"),
		).
		AddRoute(
			router.NewRoute("", router.PATCH).
				Handle(updateTask).
				WithDescription("Update task"),
		).
		AddRoute(
			router.NewRoute("/:id", router.DELETE).
				Handle(deleteTask).
				WithDescription("Delete task"),
		).
		AddRoute(
			router.NewRoute("/:id/run", router.POST).
				Handle(runTask).
				WithDescription("Run task manually"),
		).
		AddRoute(
			router.NewRoute("/:id/stop", router.POST).
				Handle(stopTask).
				WithDescription("Stop running task"),
		)
}

// createTask 创建任务
// @Summary 创建任务
// @Description 创建单个任务
// @Tags 任务管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body taskModel.CreateRequest true "创建任务请求"
// @Success 200 {object} resp.SuccessStruct{data=taskModel.Data} "创建成功"
// @Failure 400 {object} resp.ErrorStruct "请求参数错误"
// @Failure 401 {object} resp.ErrorStruct "未授权"
// @Failure 500 {object} resp.ErrorStruct "服务器内部错误"
// @Router /api/v1/tasks [post]
func createTask(c *gin.Context) {
	var req taskModel.CreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		resp.ErrorBadRequest(c)
		return
	}

	taskData, err := task.AddTask(&req)
	if err != nil {
		resp.Error(c, http.StatusInternalServerError, "failed to create task")
		return
	}

	log.Infof("Task %d created from %s", taskData.ID, c.ClientIP())

	resp.Success(c, taskData)
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
// @Success 200 {object} resp.SuccessStruct{data=[]taskModel.Data} "获取成功"
// @Failure 400 {object} resp.ErrorStruct "请求参数错误"
// @Failure 401 {object} resp.ErrorStruct "未授权"
// @Failure 500 {object} resp.ErrorStruct "服务器内部错误"
// @Router /api/v1/tasks [get]
func getTasks(c *gin.Context) {
	// 解析查询参数
	idsParam := c.Query("ids")

	// 如果指定了IDs，则获取指定的任务
	if idsParam != "" {
		idStrs := strings.Split(idsParam, ",")
		var tasks []taskModel.Data

		for _, idStr := range idStrs {
			id, err := strconv.ParseUint(strings.TrimSpace(idStr), 10, 16)
			if err != nil {
				resp.Error(c, http.StatusBadRequest, "invalid id")
				return
			}

			// 获取任务
			taskData, err := task.GetTask(uint16(id))
			if err != nil {
				resp.Error(c, http.StatusInternalServerError, "failed to get task")
				return
			}

			if taskData == nil {
				resp.Error(c, http.StatusNotFound, fmt.Sprintf("task %d not found", id))
				return
			}

			tasks = append(tasks, *taskData)
		}

		resp.Success(c, tasks)
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
		resp.Error(c, http.StatusInternalServerError, "failed to get tasks")
		return
	}

	// 构建分页响应
	result := map[string]interface{}{
		"list":      tasks,
		"total":     total,
		"page":      page,
		"page_size": pageSize,
	}

	log.Infof("Task list requested from %s", c.ClientIP())

	resp.Success(c, result)
}

// getTask 获取单个任务
// @Summary 获取任务详情
// @Tags 任务管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "任务ID"
// @Success 200 {object} resp.SuccessStruct{data=taskModel.Data} "获取成功"
// @Failure 400 {object} resp.ErrorStruct "请求参数错误"
// @Failure 401 {object} resp.ErrorStruct "未授权"
// @Failure 404 {object} resp.ErrorStruct "任务不存在"
// @Failure 500 {object} resp.ErrorStruct "服务器内部错误"
// @Router /api/v1/tasks/{id} [get]
func getTask(c *gin.Context) {
	// 获取路径参数中的ID
	idParam := c.Param("id")
	if idParam == "" {
		resp.Error(c, http.StatusBadRequest, "task id is required")
		return
	}

	id, err := strconv.ParseUint(idParam, 10, 16)
	if err != nil {
		resp.ErrorBadRequest(c)
		return
	}

	// 获取任务
	taskData, err := task.GetTask(uint16(id))
	if err != nil {
		resp.Error(c, http.StatusInternalServerError, "failed to get task")
		return
	}

	if taskData == nil {
		resp.Error(c, http.StatusNotFound, fmt.Sprintf("task %d not found", id))
		return
	}

	log.Infof("Task %d requested from %s", id, c.ClientIP())

	resp.Success(c, taskData)
}

// updateTask 更新任务
// @Summary 更新任务
// @Description 根据请求体中的ID更新任务信息
// @Tags 任务管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body taskModel.UpdateRequest true "更新任务请求"
// @Success 200 {object} resp.SuccessStruct{data=taskModel.Data} "更新成功"
// @Failure 400 {object} resp.ErrorStruct "请求参数错误"
// @Failure 401 {object} resp.ErrorStruct "未授权"
// @Failure 404 {object} resp.ErrorStruct "任务不存在"
// @Failure 500 {object} resp.ErrorStruct "服务器内部错误"
// @Router /api/v1/tasks [patch]
func updateTask(c *gin.Context) {
	var req taskModel.UpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		resp.ErrorBadRequest(c)
		return
	}

	if req.ID == 0 {
		resp.Error(c, http.StatusBadRequest, "task id is required")
		return
	}

	// 更新任务
	taskData, err := task.UpdateTask(&req)
	if err != nil {
		resp.Error(c, http.StatusInternalServerError, "failed to update task")
		return
	}

	log.Infof("Task %d updated from %s", req.ID, c.ClientIP())

	resp.Success(c, taskData)
}

// deleteTask 删除任务
// @Summary 删除任务
// @Description 根据ID删除单个任务
// @Tags 任务管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "任务ID"
// @Success 200 {object} resp.SuccessStruct "删除成功"
// @Failure 400 {object} resp.ErrorStruct "请求参数错误"
// @Failure 401 {object} resp.ErrorStruct "未授权"
// @Failure 404 {object} resp.ErrorStruct "任务不存在"
// @Failure 500 {object} resp.ErrorStruct "服务器内部错误"
// @Router /api/v1/tasks/{id} [delete]
func deleteTask(c *gin.Context) {
	// 获取路径参数中的ID
	idParam := c.Param("id")
	if idParam == "" {
		resp.Error(c, http.StatusBadRequest, "task id is required")
		return
	}

	id, err := strconv.ParseUint(idParam, 10, 16)
	if err != nil {
		resp.ErrorBadRequest(c)
		return
	}

	// 检查任务是否存在
	existingTask, err := task.GetTask(uint16(id))
	if err != nil {
		resp.Error(c, http.StatusInternalServerError, "failed to get task")
		return
	}

	if existingTask == nil {
		resp.Error(c, http.StatusNotFound, fmt.Sprintf("task %d not found", id))
		return
	}

	// 删除任务
	if err := task.RemoveTaskWithDb(uint16(id)); err != nil {
		resp.Error(c, http.StatusInternalServerError, "failed to delete task")
		return
	}

	log.Infof("Task %d deleted from %s", id, c.ClientIP())

	resp.Success(c, nil)
}

// runTask 手动运行任务
// @Summary 手动运行任务
// @Description 手动触发任务执行
// @Tags 任务管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "任务ID"
// @Success 200 {object} resp.SuccessStruct "运行成功"
// @Failure 400 {object} resp.ErrorStruct "请求参数错误"
// @Failure 401 {object} resp.ErrorStruct "未授权"
// @Failure 404 {object} resp.ErrorStruct "任务不存在"
// @Failure 500 {object} resp.ErrorStruct "服务器内部错误"
// @Router /api/v1/tasks/{id}/run [post]
func runTask(c *gin.Context) {
	// 获取路径参数中的ID
	idParam := c.Param("id")
	if idParam == "" {
		resp.Error(c, http.StatusBadRequest, "task id is required")
		return
	}

	id, err := strconv.ParseUint(idParam, 10, 16)
	if err != nil {
		resp.ErrorBadRequest(c)
		return
	}

	// 检查任务是否存在
	existingTask, err := task.GetTask(uint16(id))
	if err != nil {
		resp.Error(c, http.StatusInternalServerError, "failed to get task")
		return
	}

	if existingTask == nil {
		resp.Error(c, http.StatusNotFound, fmt.Sprintf("task %d not found", id))
		return
	}

	// 运行任务
	if err := task.RunTask(uint16(id)); err != nil {
		resp.Error(c, http.StatusInternalServerError, "failed to run task")
		return
	}

	log.Infof("Task %d manually run from %s", id, c.ClientIP())

	resp.Success(c, nil)
}

// stopTask 停止任务
// @Summary 停止任务
// @Description 停止正在运行的任务
// @Tags 任务管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "任务ID"
// @Success 200 {object} resp.SuccessStruct "停止成功"
// @Failure 400 {object} resp.ErrorStruct "请求参数错误"
// @Failure 401 {object} resp.ErrorStruct "未授权"
// @Failure 404 {object} resp.ErrorStruct "任务不存在"
// @Failure 500 {object} resp.ErrorStruct "服务器内部错误"
// @Router /api/v1/tasks/{id}/stop [post]
func stopTask(c *gin.Context) {
	// 获取路径参数中的ID
	idParam := c.Param("id")
	if idParam == "" {
		resp.Error(c, http.StatusBadRequest, "task id is required")
		return
	}

	id, err := strconv.ParseUint(idParam, 10, 16)
	if err != nil {
		resp.ErrorBadRequest(c)
		return
	}

	// 检查任务是否存在
	existingTask, err := task.GetTask(uint16(id))
	if err != nil {
		resp.Error(c, http.StatusInternalServerError, "failed to get task")
		return
	}

	if existingTask == nil {
		resp.Error(c, http.StatusNotFound, fmt.Sprintf("task %d not found", id))
		return
	}

	// 停止任务
	if err := task.StopTask(uint16(id)); err != nil {
		resp.Error(c, http.StatusInternalServerError, "failed to stop task")
		return
	}

	log.Infof("Task %d stopped from %s", id, c.ClientIP())

	resp.Success(c, nil)
}
