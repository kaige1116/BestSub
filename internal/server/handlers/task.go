package handlers

import (
	"encoding/json"
	"net/http"
	"slices"
	"strconv"

	"github.com/bestruirui/bestsub/internal/core/task"
	"github.com/bestruirui/bestsub/internal/database/op"
	taskModel "github.com/bestruirui/bestsub/internal/models/task"
	"github.com/bestruirui/bestsub/internal/modules/exec"
	"github.com/bestruirui/bestsub/internal/server/middleware"
	"github.com/bestruirui/bestsub/internal/server/resp"
	"github.com/bestruirui/bestsub/internal/server/router"
	"github.com/bestruirui/bestsub/internal/utils/desc"
	"github.com/bestruirui/bestsub/internal/utils/log"
	"github.com/gin-gonic/gin"
)

func init() {

	router.NewGroupRouter("/api/v1/task").
		Use(middleware.Auth()).
		AddRoute(
			router.NewRoute("/type", router.GET).
				Handle(getTaskTypes),
		).
		AddRoute(
			router.NewRoute("/type/config", router.GET).
				Handle(getTaskTypeConfig),
		).
		AddRoute(
			router.NewRoute("/config", router.GET).
				Handle(getTaskConfig),
		).
		AddRoute(
			router.NewRoute("", router.POST).
				Handle(createTask),
		).
		AddRoute(
			router.NewRoute("", router.GET).
				Handle(getTaskList),
		).
		AddRoute(
			router.NewRoute("", router.PUT).
				Handle(updateTask),
		).
		AddRoute(
			router.NewRoute("/:id", router.DELETE).
				Handle(deleteTask),
		).
		AddRoute(
			router.NewRoute("/:id/run", router.POST).
				Handle(runTask),
		).
		AddRoute(
			router.NewRoute("/:id/stop", router.POST).
				Handle(stopTask),
		)
}

// createTask 获取任务类型
// @Summary 获取任务类型
// @Description 获取任务类型
// @Tags 任务管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} resp.SuccessStruct{data=[]string} "获取成功"
// @Failure 400 {object} resp.ErrorStruct "请求参数错误"
// @Failure 401 {object} resp.ErrorStruct "未授权"
// @Failure 500 {object} resp.ErrorStruct "服务器内部错误"
// @Router /api/v1/task/type [get]
func getTaskTypes(c *gin.Context) {
	resp.Success(c, exec.GetTypes())
}

// createTask 获取任务类型对应的配置项
// @Summary 获取任务类型对应的配置项
// @Description 获取任务类型对应的配置项
// @Tags 任务管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param type query string true "任务类型"
// @Success 200 {object} resp.SuccessStruct{data=map[string][]exec.Desc} "获取成功"
// @Failure 400 {object} resp.ErrorStruct "请求参数错误"
// @Failure 401 {object} resp.ErrorStruct "未授权"
// @Failure 500 {object} resp.ErrorStruct "服务器内部错误"
// @Router /api/v1/task/type/config [get]
func getTaskTypeConfig(c *gin.Context) {
	typ := c.Query("type")
	if typ == "" {
		resp.Success(c, exec.GetInfoMap())
	} else {
		resp.Success(c, exec.GetInfoMap()[typ])
	}
}

// createTask 获取任务配置项
// @Summary 获取任务配置项
// @Description 获取任务配置项
// @Tags 任务管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} resp.SuccessStruct{data=[]desc.Data} "获取成功"
// @Failure 400 {object} resp.ErrorStruct "请求参数错误"
// @Failure 401 {object} resp.ErrorStruct "未授权"
// @Failure 500 {object} resp.ErrorStruct "服务器内部错误"
// @Router /api/v1/task/config [get]
func getTaskConfig(c *gin.Context) {
	resp.Success(c, desc.Gen(taskModel.Data{}))
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
// @Router /api/v1/task [post]
func createTask(c *gin.Context) {
	var req taskModel.CreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		resp.ErrorBadRequest(c)
		return
	}
	taskData := taskModel.Data{
		Name:   req.Name,
		Enable: req.Enable,
		Config: req.Config,
		Extra:  req.Extra,
	}
	if err := op.CreateTask(c.Request.Context(), &taskData); err != nil {
		log.Errorf("failed to create task: %v", err)
		resp.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	var cronTaskConfig taskModel.Config
	err := json.Unmarshal([]byte(req.Config), &cronTaskConfig)
	if err != nil {
		log.Errorf("failed to unmarshal config: %v", err)
		resp.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	cronTaskConfig.ID = taskData.ID
	cronTaskConfig.Extra = req.Extra
	task.Check.Add(&cronTaskConfig)

	resp.Success(c, taskData)
}

// getTasks 获取任务列表
// @Summary 获取任务列表
// @Tags 任务管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} resp.SuccessStruct{data=[]taskModel.Response} "获取成功"
// @Failure 400 {object} resp.ErrorStruct "请求参数错误"
// @Failure 401 {object} resp.ErrorStruct "未授权"
// @Failure 500 {object} resp.ErrorStruct "服务器内部错误"
// @Router /api/v1/task [get]
func getTaskList(c *gin.Context) {
	taskList, err := op.GetTaskList(c.Request.Context())
	if err != nil {
		resp.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	runningTaskIds := task.Check.GetRunningTaskID()
	var respTaskList = make([]taskModel.Response, len(taskList))
	for i := range taskList {
		respTaskList[i] = taskModel.Response{
			Data: taskList[i],
		}
		if slices.Contains(runningTaskIds, taskList[i].ID) {
			respTaskList[i].Status = "running"
		} else {
			respTaskList[i].Status = "stopped"
		}
	}
	resp.Success(c, respTaskList)
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
// @Router /api/v1/task [put]
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

	taskData := &taskModel.Data{
		ID:     req.ID,
		Name:   req.Name,
		Enable: req.Enable,
		Extra:  req.Extra,
	}

	err := op.UpdateTask(c.Request.Context(), taskData)
	if err != nil {
		resp.Error(c, http.StatusInternalServerError, "failed to update task")
		return
	}

	var cronTaskConfig taskModel.Config
	err = json.Unmarshal([]byte(req.Config), &cronTaskConfig)
	if err != nil {
		resp.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	cronTaskConfig.ID = req.ID
	cronTaskConfig.Extra = req.Extra
	task.Check.Update(&cronTaskConfig)

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
// @Router /api/v1/task/{id} [delete]
func deleteTask(c *gin.Context) {
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

	if err := op.DeleteTask(c.Request.Context(), uint16(id)); err != nil {
		resp.Error(c, http.StatusInternalServerError, "failed to delete task")
		return
	}
	if err := task.Check.Remove(uint16(id)); err != nil {
		resp.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

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
// @Router /api/v1/task/{id}/run [post]
func runTask(c *gin.Context) {
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

	if err := task.Check.Run(uint16(id)); err != nil {
		resp.Error(c, http.StatusInternalServerError, err.Error())
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
// @Router /api/v1/task/{id}/stop [post]
func stopTask(c *gin.Context) {
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

	if err := task.Check.StopTask(uint16(id)); err != nil {
		resp.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	log.Infof("Task %d stopped from %s", id, c.ClientIP())

	resp.Success(c, nil)
}
