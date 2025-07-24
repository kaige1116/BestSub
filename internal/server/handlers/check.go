package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/bestruirui/bestsub/internal/core/task"
	"github.com/bestruirui/bestsub/internal/database/op"
	checkModel "github.com/bestruirui/bestsub/internal/models/check"
	taskModel "github.com/bestruirui/bestsub/internal/models/task"
	"github.com/bestruirui/bestsub/internal/modules/exec"
	"github.com/bestruirui/bestsub/internal/server/middleware"
	"github.com/bestruirui/bestsub/internal/server/resp"
	"github.com/bestruirui/bestsub/internal/server/router"
	"github.com/bestruirui/bestsub/internal/utils/log"
	"github.com/gin-gonic/gin"
)

func init() {

	router.NewGroupRouter("/api/v1/check").
		Use(middleware.Auth()).
		AddRoute(
			router.NewRoute("/type", router.GET).
				Handle(getCheckTypes),
		).
		AddRoute(
			router.NewRoute("/type/config", router.GET).
				Handle(getCheckTypeConfig),
		).
		AddRoute(
			router.NewRoute("", router.POST).
				Handle(createCheck),
		).
		AddRoute(
			router.NewRoute("", router.GET).
				Handle(getCheck),
		).
		AddRoute(
			router.NewRoute("", router.PUT).
				Handle(updateCheck),
		).
		AddRoute(
			router.NewRoute("/:id", router.DELETE).
				Handle(deleteCheck),
		).
		AddRoute(
			router.NewRoute("/:id/run", router.POST).
				Handle(runCheck),
		).
		AddRoute(
			router.NewRoute("/:id/stop", router.POST).
				Handle(stopCheck),
		)
}

// getCheckTypes 获取检测类型
// @Summary 获取检测类型
// @Description 获取检测类型
// @Tags 检测管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} resp.SuccessStruct{data=[]string} "获取成功"
// @Failure 400 {object} resp.ErrorStruct "请求参数错误"
// @Failure 401 {object} resp.ErrorStruct "未授权"
// @Failure 500 {object} resp.ErrorStruct "服务器内部错误"
// @Router /api/v1/check/type [get]
func getCheckTypes(c *gin.Context) {
	resp.Success(c, exec.GetTypes())
}

// getCheckTypeConfig 获取检测类型对应的配置项
// @Summary 获取检测类型对应的配置项
// @Description 获取检测类型对应的配置项
// @Tags 检测管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param type query string true "检测类型"
// @Success 200 {object} resp.SuccessStruct{data=map[string][]exec.Desc} "获取成功"
// @Failure 400 {object} resp.ErrorStruct "请求参数错误"
// @Failure 401 {object} resp.ErrorStruct "未授权"
// @Failure 500 {object} resp.ErrorStruct "服务器内部错误"
// @Router /api/v1/check/type/config [get]
func getCheckTypeConfig(c *gin.Context) {
	typ := c.Query("type")
	if typ == "" {
		resp.Success(c, exec.GetInfoMap())
	} else {
		resp.Success(c, exec.GetInfoMap()[typ])
	}
}

// createCheck 创建检测
// @Summary 创建检测
// @Description 创建单个检测
// @Tags 检测管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body checkModel.CreateRequest true "创建检测请求"
// @Success 200 {object} resp.SuccessStruct{data=checkModel.Response} "创建成功"
// @Failure 400 {object} resp.ErrorStruct "请求参数错误"
// @Failure 401 {object} resp.ErrorStruct "未授权"
// @Failure 500 {object} resp.ErrorStruct "服务器内部错误"
// @Router /api/v1/check [post]
func createCheck(c *gin.Context) {
	var req checkModel.CreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		resp.ErrorBadRequest(c)
		return
	}
	taskBytes, err := json.Marshal(req.Task)
	if err != nil {
		resp.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	taskStr := string(taskBytes)
	configBytes, err := json.Marshal(req.Config)
	if err != nil {
		resp.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	configStr := string(configBytes)
	checkData := checkModel.Data{
		Name:   req.Name,
		Enable: req.Enable,
		Task:   taskStr,
		Config: configStr,
		Result: "{}",
	}
	if err := op.CreateCheck(c.Request.Context(), &checkData); err != nil {
		log.Errorf("failed to create check: %v", err)
		resp.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	req.Task.ID = checkData.ID
	req.Task.Name = checkData.Name
	task.Check.Add(&req.Task, configStr)
	if req.Enable {
		task.Check.Enable(checkData.ID)
	}
	var taskConfig taskModel.Config
	if err := json.Unmarshal([]byte(checkData.Task), &taskConfig); err != nil {
		resp.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	var config any
	if err := json.Unmarshal([]byte(checkData.Config), &config); err != nil {
		resp.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	respCheckData := checkModel.Response{
		ID:     checkData.ID,
		Name:   checkData.Name,
		Enable: checkData.Enable,
		Task:   taskConfig,
		Config: config,
		Status: task.Check.Status(checkData.ID),
	}

	resp.Success(c, respCheckData)
}

// getCheck 获取检测列表
// @Summary 获取检测列表
// @Tags 检测管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id query int true "检测ID"
// @Success 200 {object} resp.SuccessStruct{data=[]checkModel.Response} "获取成功"
// @Failure 400 {object} resp.ErrorStruct "请求参数错误"
// @Failure 401 {object} resp.ErrorStruct "未授权"
// @Failure 500 {object} resp.ErrorStruct "服务器内部错误"
// @Router /api/v1/check [get]
func getCheck(c *gin.Context) {
	idStr := c.Query("id")
	if idStr == "" {
		checkList, err := op.GetCheckList(c.Request.Context())
		if err != nil {
			resp.Error(c, http.StatusInternalServerError, err.Error())
			return
		}
		var respCheckList = make([]checkModel.Response, len(checkList))
		for i := range checkList {
			var taskConfig taskModel.Config
			if err := json.Unmarshal([]byte(checkList[i].Task), &taskConfig); err != nil {
				resp.Error(c, http.StatusInternalServerError, err.Error())
				return
			}
			var config any
			if err := json.Unmarshal([]byte(checkList[i].Config), &config); err != nil {
				resp.Error(c, http.StatusInternalServerError, err.Error())
				return
			}
			var result taskModel.DBResult
			if err := json.Unmarshal([]byte(checkList[i].Result), &result); err != nil {
				resp.Error(c, http.StatusInternalServerError, err.Error())
				return
			}
			respCheckList[i] = checkModel.Response{
				ID:     checkList[i].ID,
				Name:   checkList[i].Name,
				Enable: checkList[i].Enable,
				Task:   taskConfig,
				Config: config,
				Result: result,
				Status: task.Check.Status(checkList[i].ID),
			}
		}
		resp.Success(c, respCheckList)
	} else {
		id, err := strconv.ParseUint(idStr, 10, 16)
		if err != nil {
			resp.ErrorBadRequest(c)
			return
		}
		check, err := op.GetCheckByID(uint16(id))
		if err != nil {
			resp.Error(c, http.StatusInternalServerError, err.Error())
			return
		}
		var respCheck = make([]checkModel.Response, 1)
		var taskConfig taskModel.Config
		if err := json.Unmarshal([]byte(check.Task), &taskConfig); err != nil {
			resp.Error(c, http.StatusInternalServerError, err.Error())
			return
		}
		var config any
		if err := json.Unmarshal([]byte(check.Config), &config); err != nil {
			resp.Error(c, http.StatusInternalServerError, err.Error())
			return
		}
		var result taskModel.DBResult
		if err := json.Unmarshal([]byte(check.Result), &result); err != nil {
			resp.Error(c, http.StatusInternalServerError, err.Error())
			return
		}
		respCheck[0] = checkModel.Response{
			ID:     check.ID,
			Name:   check.Name,
			Enable: check.Enable,
			Task:   taskConfig,
			Config: config,
			Result: result,
			Status: task.Check.Status(check.ID),
		}
		resp.Success(c, respCheck)
	}
}

// updateCheck 更新检测
// @Summary 更新检测
// @Description 根据请求体中的ID更新检测信息
// @Tags 检测管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body checkModel.UpdateRequest true "更新检测请求"
// @Success 200 {object} resp.SuccessStruct{data=checkModel.Response} "更新成功"
// @Failure 400 {object} resp.ErrorStruct "请求参数错误"
// @Failure 401 {object} resp.ErrorStruct "未授权"
// @Failure 404 {object} resp.ErrorStruct "检测不存在"
// @Failure 500 {object} resp.ErrorStruct "服务器内部错误"
// @Router /api/v1/check [put]
func updateCheck(c *gin.Context) {
	var req checkModel.UpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		resp.ErrorBadRequest(c)
		return
	}

	if req.ID == 0 {
		resp.Error(c, http.StatusBadRequest, "check id is required")
		return
	}

	taskBytes, err := json.Marshal(req.Task)
	if err != nil {
		resp.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	taskStr := string(taskBytes)
	configBytes, err := json.Marshal(req.Config)
	if err != nil {
		resp.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	configStr := string(configBytes)
	checkData := &checkModel.Data{
		ID:     req.ID,
		Name:   req.Name,
		Enable: req.Enable,
		Task:   taskStr,
		Config: configStr,
	}
	req.Task.ID = checkData.ID
	req.Task.Name = checkData.Name
	err = op.UpdateCheck(c.Request.Context(), checkData)
	if err != nil {
		resp.Error(c, http.StatusInternalServerError, "failed to update check")
		return
	}

	task.Check.Update(&req.Task, configStr)
	if req.Enable {
		task.Check.Enable(req.ID)
	} else {
		task.Check.Disable(req.ID)
	}

	log.Infof("Check %d updated from %s", req.ID, c.ClientIP())

	var taskConfig taskModel.Config
	if err := json.Unmarshal([]byte(checkData.Task), &taskConfig); err != nil {
		resp.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	var config any
	if err := json.Unmarshal([]byte(checkData.Config), &config); err != nil {
		resp.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	respCheckData := checkModel.Response{
		ID:     req.ID,
		Name:   req.Name,
		Enable: req.Enable,
		Task:   taskConfig,
		Config: config,
		Status: task.Check.Status(req.ID),
	}
	resp.Success(c, respCheckData)
}

// deleteCheck 删除检测
// @Summary 删除检测
// @Description 根据ID删除单个检测
// @Tags 检测管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "检测ID"
// @Success 200 {object} resp.SuccessStruct "删除成功"
// @Failure 400 {object} resp.ErrorStruct "请求参数错误"
// @Failure 401 {object} resp.ErrorStruct "未授权"
// @Failure 404 {object} resp.ErrorStruct "检测不存在"
// @Failure 500 {object} resp.ErrorStruct "服务器内部错误"
// @Router /api/v1/check/{id} [delete]
func deleteCheck(c *gin.Context) {
	idParam := c.Param("id")
	if idParam == "" {
		resp.Error(c, http.StatusBadRequest, "check id is required")
		return
	}

	id, err := strconv.ParseUint(idParam, 10, 16)
	if err != nil {
		resp.ErrorBadRequest(c)
		return
	}

	if err := op.DeleteCheck(c.Request.Context(), uint16(id)); err != nil {
		resp.Error(c, http.StatusInternalServerError, "failed to delete check")
		return
	}
	if err := task.Check.Remove(uint16(id)); err != nil {
		resp.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	resp.Success(c, nil)
}

// runCheck 手动运行检测
// @Summary 手动运行检测
// @Description 手动触发检测执行
// @Tags 检测管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "检测ID"
// @Success 200 {object} resp.SuccessStruct "运行成功"
// @Failure 400 {object} resp.ErrorStruct "请求参数错误"
// @Failure 401 {object} resp.ErrorStruct "未授权"
// @Failure 404 {object} resp.ErrorStruct "检测不存在"
// @Failure 500 {object} resp.ErrorStruct "服务器内部错误"
// @Router /api/v1/check/{id}/run [post]
func runCheck(c *gin.Context) {
	idParam := c.Param("id")
	if idParam == "" {
		resp.Error(c, http.StatusBadRequest, "check id is required")
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

	log.Infof("Check %d manually run from %s", id, c.ClientIP())

	resp.Success(c, nil)
}

// stopCheck 停止检测
// @Summary 停止检测
// @Description 停止正在运行的检测
// @Tags 检测管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "检测ID"
// @Success 200 {object} resp.SuccessStruct "停止成功"
// @Failure 400 {object} resp.ErrorStruct "请求参数错误"
// @Failure 401 {object} resp.ErrorStruct "未授权"
// @Failure 404 {object} resp.ErrorStruct "检测不存在"
// @Failure 500 {object} resp.ErrorStruct "服务器内部错误"
// @Router /api/v1/check/{id}/stop [post]
func stopCheck(c *gin.Context) {
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

	log.Infof("Check %d stopped from %s", id, c.ClientIP())

	resp.Success(c, nil)
}
