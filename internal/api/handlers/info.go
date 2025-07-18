package handlers

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/bestruirui/bestsub/internal/api/common"
	"github.com/bestruirui/bestsub/internal/api/middleware"
	"github.com/bestruirui/bestsub/internal/api/router"
	sys "github.com/bestruirui/bestsub/internal/core/system"
	"github.com/bestruirui/bestsub/internal/database/op"
	"github.com/bestruirui/bestsub/internal/models/system"
	"github.com/bestruirui/bestsub/internal/utils/info"
	"github.com/bestruirui/bestsub/internal/utils/local"
	"github.com/bestruirui/bestsub/internal/utils/log"
	"github.com/gin-gonic/gin"
)

func init() {
	router.NewGroupRouter("/api/v1/system").
		AddRoute(
			router.NewRoute("/health", router.GET).
				Handle(healthCheck).
				WithDescription("Health check endpoint"),
		).
		AddRoute(
			router.NewRoute("/ready", router.GET).
				Handle(readinessCheck).
				WithDescription("Readiness check endpoint"),
		).
		AddRoute(
			router.NewRoute("/live", router.GET).
				Handle(livenessCheck).
				WithDescription("Liveness check endpoint"),
		)

	router.NewGroupRouter("/api/v1/system").
		Use(middleware.Auth()).
		AddRoute(
			router.NewRoute("/info", router.GET).
				Handle(systemInfo).
				WithDescription("Get system information"),
		)
}

// healthCheck 健康检查
// @Summary 健康检查
// @Description 检查服务健康状态，包括数据库连接状态
// @Tags 系统
// @Accept json
// @Produce json
// @Success 200 {object} common.ResponseSuccessStruct{data=system.HealthResponse} "服务正常"
// @Failure 503 {object} common.ResponseErrorStruct "服务不可用"
// @Router /api/v1/system/health [get]
func healthCheck(c *gin.Context) {
	// 检查数据库连接状态
	opStatus := "connected"

	// 尝试执行一个简单的数据库查询来检查连接
	authRepo := op.AuthRepo()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := authRepo.IsInitialized(ctx)
	if err != nil {
		log.Errorf("Database health check failed: %v", err)
		opStatus = "disconnected"
	}

	response := system.HealthResponse{
		Status:    "ok",
		Timestamp: local.Time().Format(time.RFC3339),
		Version:   info.Version,
		Database:  opStatus,
	}

	// 如果数据库连接失败，返回503状态码
	if opStatus == "disconnected" {
		response.Status = "error"
		common.ResponseError(c, http.StatusServiceUnavailable, errors.New("database connection failed"))
		return
	}

	common.ResponseSuccess(c, response)
}

// readinessCheck 就绪检查
// @Summary 就绪检查
// @Description 检查服务是否准备好接收请求
// @Tags 系统
// @Accept json
// @Produce json
// @Success 200 {object} common.ResponseSuccessStruct{data=system.HealthResponse} "服务就绪"
// @Failure 503 {object} common.ResponseErrorStruct "服务未就绪"
// @Router /api/v1/system/ready [get]
func readinessCheck(c *gin.Context) {
	// 检查关键组件是否就绪
	ready := true
	var errorMsg string

	// 检查数据库是否已初始化
	authRepo := op.AuthRepo()
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	isInitialized, err := authRepo.IsInitialized(ctx)
	if err != nil || !isInitialized {
		ready = false
		errorMsg = "Database not initialized"
		log.Errorf("Readiness check failed: op not initialized, error: %v", err)
	}

	response := system.HealthResponse{
		Status:    "ready",
		Timestamp: local.Time().Format(time.RFC3339),
		Version:   info.Version,
		Database:  "initialized",
	}

	if !ready {
		response.Status = "not ready"
		response.Database = "not initialized"
		common.ResponseError(c, http.StatusServiceUnavailable, errors.New(errorMsg))
		return
	}

	common.ResponseSuccess(c, response)
}

// livenessCheck 存活检查
// @Summary 存活检查
// @Description 检查服务是否存活（简单的ping检查）
// @Tags 系统
// @Accept json
// @Produce json
// @Success 200 {object} common.ResponseSuccessStruct "服务存活"
// @Router /api/v1/system/live [get]
func livenessCheck(c *gin.Context) {
	common.ResponseSuccess(c, map[string]interface{}{
		"status":    "alive",
		"timestamp": local.Time().Format(time.RFC3339),
	})
}

// systemInfo 系统信息
// @Summary 系统信息
// @Description 获取程序运行相关信息，包括内存使用、运行时长、网络流量、CPU信息等
// @Tags 系统
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} common.ResponseSuccessStruct{data=system.Info} "获取成功"
// @Failure 401 {object} common.ResponseErrorStruct "未授权"
// @Failure 500 {object} common.ResponseErrorStruct "服务器内部错误"
// @Router /api/v1/system/info [get]
func systemInfo(c *gin.Context) {
	common.ResponseSuccess(c, sys.GetSystemInfo())
}
