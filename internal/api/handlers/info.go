package handlers

import (
	"context"
	"net/http"
	"time"

	"github.com/bestruirui/bestsub/internal/api/middleware"
	"github.com/bestruirui/bestsub/internal/api/models"
	"github.com/bestruirui/bestsub/internal/api/router"
	"github.com/bestruirui/bestsub/internal/database"
	"github.com/bestruirui/bestsub/internal/utils/info"
	"github.com/bestruirui/bestsub/internal/utils/log"
	"github.com/bestruirui/bestsub/internal/utils/system"
	timeutils "github.com/bestruirui/bestsub/internal/utils/time"
	"github.com/gin-gonic/gin"
)

// healthHandler 健康检查处理器
type healthHandler struct{}

// init 函数用于自动注册路由
func init() {
	h := newHealthHandler()

	router.NewGroupRouter("/api/v1/system").
		AddRoute(
			router.NewRoute("/health", router.GET).
				Handle(h.healthCheck).
				WithDescription("Health check endpoint"),
		).
		AddRoute(
			router.NewRoute("/ready", router.GET).
				Handle(h.readinessCheck).
				WithDescription("Readiness check endpoint"),
		).
		AddRoute(
			router.NewRoute("/live", router.GET).
				Handle(h.livenessCheck).
				WithDescription("Liveness check endpoint"),
		)

	router.NewGroupRouter("/api/v1/system").
		Use(middleware.Auth()).
		AddRoute(
			router.NewRoute("/info", router.GET).
				Handle(h.systemInfo).
				WithDescription("Get system information"),
		)
}

// newHealthHandler 创建健康检查处理器
func newHealthHandler() *healthHandler {
	return &healthHandler{}
}

// healthCheck 健康检查
// @Summary 健康检查
// @Description 检查服务健康状态，包括数据库连接状态
// @Tags 系统
// @Accept json
// @Produce json
// @Success 200 {object} models.SuccessResponse{data=models.HealthResponse} "服务正常"
// @Failure 503 {object} models.ErrorResponse "服务不可用"
// @Router /api/v1/system/health [get]
func (h *healthHandler) healthCheck(c *gin.Context) {
	// 检查数据库连接状态
	databaseStatus := "connected"

	// 尝试执行一个简单的数据库查询来检查连接
	authRepo := database.Auth()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := authRepo.IsInitialized(ctx)
	if err != nil {
		log.Errorf("Database health check failed: %v", err)
		databaseStatus = "disconnected"
	}

	response := models.HealthResponse{
		Status:    "ok",
		Timestamp: timeutils.Now().Format(time.RFC3339),
		Version:   info.Version,
		Database:  databaseStatus,
	}

	// 如果数据库连接失败，返回503状态码
	if databaseStatus == "disconnected" {
		response.Status = "error"
		c.JSON(http.StatusServiceUnavailable, models.ErrorResponse{
			Code:    http.StatusServiceUnavailable,
			Message: "Service Unavailable",
			Error:   "Database connection failed",
		})
		return
	}

	c.JSON(http.StatusOK, models.SuccessResponse{
		Code:    http.StatusOK,
		Message: "Service is healthy",
		Data:    response,
	})
}

// readinessCheck 就绪检查
// @Summary 就绪检查
// @Description 检查服务是否准备好接收请求
// @Tags 系统
// @Accept json
// @Produce json
// @Success 200 {object} models.SuccessResponse{data=models.HealthResponse} "服务就绪"
// @Failure 503 {object} models.ErrorResponse "服务未就绪"
// @Router /api/v1/system/ready [get]
func (h *healthHandler) readinessCheck(c *gin.Context) {
	// 检查关键组件是否就绪
	ready := true
	var errorMsg string

	// 检查数据库是否已初始化
	authRepo := database.Auth()
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	isInitialized, err := authRepo.IsInitialized(ctx)
	if err != nil || !isInitialized {
		ready = false
		errorMsg = "Database not initialized"
		log.Errorf("Readiness check failed: database not initialized, error: %v", err)
	}

	response := models.HealthResponse{
		Status:    "ready",
		Timestamp: timeutils.Now().Format(time.RFC3339),
		Version:   info.Version,
		Database:  "initialized",
	}

	if !ready {
		response.Status = "not ready"
		response.Database = "not initialized"
		c.JSON(http.StatusServiceUnavailable, models.ErrorResponse{
			Code:    http.StatusServiceUnavailable,
			Message: "Service Not Ready",
			Error:   errorMsg,
		})
		return
	}

	c.JSON(http.StatusOK, models.SuccessResponse{
		Code:    http.StatusOK,
		Message: "Service is ready",
		Data:    response,
	})
}

// livenessCheck 存活检查
// @Summary 存活检查
// @Description 检查服务是否存活（简单的ping检查）
// @Tags 系统
// @Accept json
// @Produce json
// @Success 200 {object} models.SuccessResponse "服务存活"
// @Router /api/v1/system/live [get]
func (h *healthHandler) livenessCheck(c *gin.Context) {
	c.JSON(http.StatusOK, models.SuccessResponse{
		Code:    http.StatusOK,
		Message: "Service is alive",
		Data: map[string]interface{}{
			"status":    "alive",
			"timestamp": timeutils.Now().Format(time.RFC3339),
		},
	})
}

// systemInfo 系统信息
// @Summary 系统信息
// @Description 获取程序运行相关信息，包括内存使用、运行时长、网络流量、CPU信息等
// @Tags 系统
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} models.SuccessResponse{data=system.Info} "获取成功"
// @Failure 401 {object} models.ErrorResponse "未授权"
// @Failure 500 {object} models.ErrorResponse "服务器内部错误"
// @Router /api/v1/system/info [get]
func (h *healthHandler) systemInfo(c *gin.Context) {
	// 获取系统监控器实例
	sysInfo := system.GetSystemInfo()

	// 获取系统信息

	log.Debug("System information retrieved successfully")

	c.JSON(http.StatusOK, models.SuccessResponse{
		Code:    http.StatusOK,
		Message: "System information retrieved successfully",
		Data:    sysInfo,
	})
}
