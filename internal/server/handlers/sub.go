package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/bestruirui/bestsub/internal/core/task"
	"github.com/bestruirui/bestsub/internal/database/op"
	"github.com/bestruirui/bestsub/internal/models/sub"
	taskModel "github.com/bestruirui/bestsub/internal/models/task"
	execer "github.com/bestruirui/bestsub/internal/modules/exec/executor"
	"github.com/bestruirui/bestsub/internal/server/middleware"
	"github.com/bestruirui/bestsub/internal/server/resp"
	"github.com/bestruirui/bestsub/internal/server/router"
	"github.com/bestruirui/bestsub/internal/utils/log"
	"github.com/gin-gonic/gin"
)

func init() {
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
			router.NewRoute("", router.PUT).
				Handle(updateSub),
		).
		AddRoute(
			router.NewRoute("/:id", router.DELETE).
				Handle(deleteSub),
		).
		AddRoute(
			router.NewRoute("/refresh/:id", router.POST).
				Handle(refreshSub),
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
	subData := sub.Data{
		Name:     req.Name,
		Enable:   req.Enable,
		CronExpr: req.CronExpr,
		Config:   "",
		Result:   "{}",
	}
	if err := op.CreateSub(c.Request.Context(), &subData); err != nil {
		log.Errorf("failed to create sub: %v", err)
		resp.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	req.Config.SubID = subData.ID
	configBytes, err := json.Marshal(req.Config)
	if err != nil {
		resp.ErrorBadRequest(c)
		return
	}
	configStr := string(configBytes)
	subData.Config = configStr
	if err := op.UpdateSub(c.Request.Context(), &subData); err != nil {
		log.Errorf("failed to update sub: %v", err)
		resp.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	var taskData = &taskModel.Config{
		ID:       subData.ID,
		Name:     subData.Name,
		CronExpr: subData.CronExpr,
		LogLevel: "info",
		Timeout:  60,
		Type:     "fetch",
	}
	task.Fetch.Add(taskData, configStr)
	if subData.Enable {
		task.Fetch.Enable(subData.ID)
	} else {
		task.Fetch.Disable(subData.ID)
	}
	resp.Success(c, sub.Response{
		ID:       subData.ID,
		Name:     subData.Name,
		Enable:   subData.Enable,
		CronExpr: subData.CronExpr,
		Config: execer.Fetch{
			SubUrl:    req.Config.SubUrl,
			Proxy:     req.Config.Proxy,
			FailTimes: req.Config.FailTimes,
		},
		Status: task.Fetch.Status(subData.ID),
	})
}

// getSubs 获取订阅链接
// @Summary 获取订阅链接
// @Tags 订阅管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id query int true "链接ID"
// @Success 200 {object} resp.SuccessStruct{data=[]sub.Response} "获取成功"
// @Failure 400 {object} resp.ErrorStruct "请求参数错误"
// @Failure 401 {object} resp.ErrorStruct "未授权"
// @Failure 500 {object} resp.ErrorStruct "服务器内部错误"
// @Router /api/v1/sub [get]
func getSubs(c *gin.Context) {
	idStr := c.Query("id")
	if idStr == "" {
		subList, err := op.GetSubList(c.Request.Context())
		if err != nil {
			resp.Error(c, http.StatusInternalServerError, err.Error())
			return
		}
		var respSubList = make([]sub.Response, len(subList))
		for i := range subList {
			var config execer.Fetch
			if err := json.Unmarshal([]byte(subList[i].Config), &config); err != nil {
				resp.Error(c, http.StatusInternalServerError, err.Error())
				return
			}
			var result taskModel.DBResult
			if err := json.Unmarshal([]byte(subList[i].Result), &result); err != nil {
				resp.Error(c, http.StatusInternalServerError, err.Error())
				return
			}
			respSubList[i] = sub.Response{
				ID:       subList[i].ID,
				Name:     subList[i].Name,
				Enable:   subList[i].Enable,
				CronExpr: subList[i].CronExpr,
				Config:   config,
				Status:   task.Fetch.Status(subList[i].ID),
				Result:   result,
			}
		}
		resp.Success(c, respSubList)
	} else {
		id, err := strconv.ParseUint(idStr, 10, 16)
		if err != nil {
			resp.ErrorBadRequest(c)
			return
		}
		subData, err := op.GetSubByID(c.Request.Context(), uint16(id))
		if err != nil {
			resp.Error(c, http.StatusInternalServerError, err.Error())
			return
		}
		var respSub = [1]sub.Response{}
		var config execer.Fetch
		if err := json.Unmarshal([]byte(subData.Config), &config); err != nil {
			resp.Error(c, http.StatusInternalServerError, err.Error())
			return
		}
		var result taskModel.DBResult
		if err := json.Unmarshal([]byte(subData.Result), &result); err != nil {
			resp.Error(c, http.StatusInternalServerError, err.Error())
			return
		}
		respSub[0] = sub.Response{
			ID:       subData.ID,
			Name:     subData.Name,
			Enable:   subData.Enable,
			CronExpr: subData.CronExpr,
			Config:   config,
			Result:   result,
			Status:   task.Fetch.Status(subData.ID),
		}
		resp.Success(c, respSub)
	}
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
// @Router /api/v1/sub [put]
func updateSub(c *gin.Context) {
	var req sub.UpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		resp.ErrorBadRequest(c)
		return
	}
	req.Config.SubID = req.ID
	configBytes, err := json.Marshal(req.Config)
	if err != nil {
		resp.ErrorBadRequest(c)
		return
	}
	configStr := string(configBytes)
	subData := &sub.Data{
		ID:       req.ID,
		Name:     req.Name,
		Enable:   req.Enable,
		CronExpr: req.CronExpr,
		Config:   configStr,
	}
	if err := op.UpdateSub(c.Request.Context(), subData); err != nil {
		log.Errorf("failed to update sub: %v", err)
		resp.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	var result taskModel.DBResult
	if err := json.Unmarshal([]byte(subData.Result), &result); err != nil {
		resp.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	var taskData = &taskModel.Config{
		ID:       subData.ID,
		Name:     subData.Name,
		CronExpr: subData.CronExpr,
		LogLevel: "info",
		Timeout:  60,
		Type:     "fetch",
	}
	task.Fetch.Update(taskData, configStr)
	if subData.Enable {
		task.Fetch.Enable(subData.ID)
	} else {
		task.Fetch.Disable(subData.ID)
	}
	resp.Success(c, sub.Response{
		ID:       subData.ID,
		Name:     subData.Name,
		Enable:   subData.Enable,
		CronExpr: subData.CronExpr,
		Config: execer.Fetch{
			SubUrl:    req.Config.SubUrl,
			Proxy:     req.Config.Proxy,
			FailTimes: req.Config.FailTimes,
		},
		Status: task.Fetch.Status(subData.ID),
		Result: result,
	})
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
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 16)
	if err != nil {
		resp.ErrorBadRequest(c)
		return
	}
	if err := op.DeleteSub(c.Request.Context(), uint16(id)); err != nil {
		resp.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	if err := task.Fetch.Remove(uint16(id)); err != nil {
		resp.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	resp.Success(c, nil)
}

// refreshSub 手动刷新订阅
// @Summary 手动刷新订阅
// @Description 根据ID手动刷新单个订阅
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
// @Router /api/v1/sub/refresh/{id} [post]
func refreshSub(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 16)
	if err != nil {
		resp.ErrorBadRequest(c)
		return
	}
	if err := task.Fetch.Run(uint16(id)); err != nil {
		resp.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	resp.Success(c, nil)
}
