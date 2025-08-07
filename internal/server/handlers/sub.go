package handlers

import (
	"net/http"
	"strconv"

	"github.com/bestruirui/bestsub/internal/core/cron"
	"github.com/bestruirui/bestsub/internal/core/nodepool"
	"github.com/bestruirui/bestsub/internal/database/op"
	"github.com/bestruirui/bestsub/internal/models/sub"
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
			router.NewRoute("/:id", router.PUT).
				Handle(updateSub),
		).
		AddRoute(
			router.NewRoute("/:id", router.DELETE).
				Handle(deleteSub),
		).
		AddRoute(
			router.NewRoute("/refresh/:id", router.POST).
				Handle(refreshSub),
		).
		AddRoute(
			router.NewRoute("/name", router.GET).
				Handle(getSubNameAndID),
		)
}

// createSub 创建订阅链接
// @Summary 创建订阅链接
// @Description 创建单个订阅链接
// @Tags 订阅
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body sub.Request true "创建订阅链接请求"
// @Success 200 {object} resp.SuccessStruct{data=sub.Response} "创建成功"
// @Failure 400 {object} resp.ErrorStruct "请求参数错误"
// @Failure 401 {object} resp.ErrorStruct "未授权"
// @Failure 500 {object} resp.ErrorStruct "服务器内部错误"
// @Router /api/v1/sub [post]
func createSub(c *gin.Context) {
	var req sub.Request
	if err := c.ShouldBindJSON(&req); err != nil {
		resp.ErrorBadRequest(c)
		return
	}
	subData := req.GenData(0)
	if err := op.CreateSub(c.Request.Context(), &subData); err != nil {
		log.Errorf("failed to create sub: %v", err)
		resp.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	cron.FetchAdd(&subData)
	respData := subData.GenResponse(cron.FetchStatus(subData.ID), sub.NodeInfo{})
	resp.Success(c, respData)
}

// getSubs 获取订阅链接
// @Summary 获取订阅链接
// @Tags 订阅
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
			respSubList[i] = subList[i].GenResponse(cron.FetchStatus(subList[i].ID), nodepool.GetPoolBySubID(subList[i].ID, 0).Info)
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
		respSub[0] = subData.GenResponse(cron.FetchStatus(subData.ID), nodepool.GetPoolBySubID(subData.ID, 0).Info)
		resp.Success(c, respSub)
	}
}

// getSubNameAndID 获取订阅链接名称和ID
// @Summary 获取订阅链接名称和ID
// @Tags 订阅
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} resp.SuccessStruct{data=[]sub.NameAndID} "获取成功"
// @Failure 400 {object} resp.ErrorStruct "请求参数错误"
// @Failure 401 {object} resp.ErrorStruct "未授权"
// @Failure 500 {object} resp.ErrorStruct "服务器内部错误"
// @Router /api/v1/sub/name [get]
func getSubNameAndID(c *gin.Context) {
	subList, err := op.GetSubList(c.Request.Context())
	if err != nil {
		resp.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	var respSubList = make([]sub.NameAndID, len(subList))
	for i, sub := range subList {
		respSubList[i].ID = sub.ID
		respSubList[i].Name = sub.Name
	}
	resp.Success(c, respSubList)
}

// updateSub 更新订阅链接
// @Summary 更新订阅链接
// @Description 根据请求体中的ID更新订阅链接信息
// @Tags 订阅
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body sub.Request true "更新订阅链接请求"
// @Success 200 {object} resp.SuccessStruct{data=sub.Response} "更新成功"
// @Failure 400 {object} resp.ErrorStruct "请求参数错误"
// @Failure 401 {object} resp.ErrorStruct "未授权"
// @Failure 404 {object} resp.ErrorStruct "订阅链接不存在"
// @Failure 500 {object} resp.ErrorStruct "服务器内部错误"
// @Router /api/v1/sub [put]
func updateSub(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 16)
	if err != nil {
		resp.ErrorBadRequest(c)
		return
	}
	var req sub.Request
	if err := c.ShouldBindJSON(&req); err != nil {
		resp.ErrorBadRequest(c)
		return
	}
	subData := req.GenData(uint16(id))
	if err := op.UpdateSub(c.Request.Context(), &subData); err != nil {
		log.Errorf("failed to update sub: %v", err)
		resp.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	if err := cron.FetchUpdate(&subData); err != nil {
		log.Errorf("failed to update sub: %v", err)
		resp.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	respData := subData.GenResponse(cron.FetchStatus(subData.ID), nodepool.GetPoolBySubID(subData.ID, 0).Info)
	resp.Success(c, respData)
}

// deleteSub 删除订阅链接
// @Summary 删除订阅链接
// @Description 根据ID删除单个订阅链接
// @Tags 订阅
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
	if err := cron.FetchRemove(uint16(id)); err != nil {
		resp.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	nodepool.DeletePool(uint16(id))
	resp.Success(c, nil)
}

// refreshSub 手动刷新订阅
// @Summary 手动刷新订阅
// @Description 根据ID手动刷新单个订阅
// @Tags 订阅
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "订阅链接ID"
// @Success 200 {object} resp.SuccessStruct{data=sub.Result} "刷新成功"
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
	result := cron.FetchRun(uint16(id))
	resp.Success(c, result)
}
