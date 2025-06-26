// Package server provides the entry point for BestSub application.
//
// @title BestSub API
// @version 1.0.0
// @description BestSub - Best Subscription Manager API Documentation
// @description
// @description This is the API documentation for BestSub, a subscription management system.
// @description
// @description ## Authentication
// @description Most endpoints require authentication using JWT tokens.
// @description To authenticate, include the JWT token in the Authorization header:
// @description `Authorization: Bearer <your-jwt-token>`
// @description
// @description ## Error Responses
// @description All error responses follow a consistent format with code, message, and error fields.
// @description
// @description ## Success Responses
// @description All success responses follow a consistent format with code, message, and data fields.

// @contact.name BestSub API Support
// @contact.email support@bestsub.com

// @license.name MIT
// @license.url https://opensource.org/licenses/MIT

// @BasePath /

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.

// @tag.name 认证
// @tag.description 用户认证相关接口

// @tag.name 系统
// @tag.description 系统状态和健康检查接口
package server

import (
	"context"
	"fmt"
	"net/http"
	"sync"
	"time"

	_ "github.com/bestruirui/bestsub/internal/api/handlers"
	"github.com/bestruirui/bestsub/internal/api/middleware"
	"github.com/bestruirui/bestsub/internal/api/router"
	"github.com/bestruirui/bestsub/internal/config"
	"github.com/bestruirui/bestsub/internal/utils/log"
	"github.com/gin-gonic/gin"
)

const (
	defaultReadTimeout     = 30 * time.Second
	defaultWriteTimeout    = 30 * time.Second
	defaultIdleTimeout     = 60 * time.Second
	defaultShutdownTimeout = 30 * time.Second
	defaultMaxHeaderBytes  = 1 << 20 // 1MB
)

var (
	server *Server
	once   sync.Once
)

// Server HTTP服务器
type Server struct {
	httpServer *http.Server
	router     *gin.Engine
	config     config.Config
}

// 初始化HTTP服务器
func Initialize() error {
	var err error
	once.Do(func() {
		cfg := config.Get()

		r, routerErr := setRouter()
		if routerErr != nil {
			err = fmt.Errorf("设置路由失败: %w", routerErr)
			return
		}

		server = &Server{
			httpServer: &http.Server{
				Addr:           fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port),
				Handler:        r,
				ReadTimeout:    defaultReadTimeout,
				WriteTimeout:   defaultWriteTimeout,
				IdleTimeout:    defaultIdleTimeout,
				MaxHeaderBytes: defaultMaxHeaderBytes,
			},
			router: r,
			config: cfg,
		}

		log.Debugf("HTTP 服务器初始化成功 %s", server.httpServer.Addr)
	})

	return err
}

// 启动HTTP服务器
func Start() error {
	if server == nil {
		return fmt.Errorf("HTTP 服务器未初始化, 请先调用 Initialize()")
	}

	log.Infof("启动 HTTP 服务器 %s", server.httpServer.Addr)

	if err := server.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return fmt.Errorf("启动 HTTP 服务器失败: %w", err)
	}

	return nil
}

// 关闭HTTP服务器
func Close() error {
	if server == nil {
		return fmt.Errorf("HTTP 服务器未初始化")
	}

	log.Debug("关闭 HTTP 服务器...")

	ctx, cancel := context.WithTimeout(context.Background(), defaultShutdownTimeout)
	defer cancel()

	if err := server.httpServer.Shutdown(ctx); err != nil {
		log.Errorf("HTTP 服务器强制关闭: %v", err)
		return fmt.Errorf("HTTP 服务器强制关闭: %w", err)
	}

	log.Debug("HTTP 服务器关闭完成")
	return nil
}

// 检查服务器是否已初始化
func IsInitialized() bool {
	return server != nil
}

// 设置路由
func setRouter() (*gin.Engine, error) {
	gin.SetMode(gin.ReleaseMode)
	r := gin.New()

	r.Use(middleware.Logging())
	r.Use(middleware.Recovery())
	r.Use(middleware.Cors())

	if err := router.RegisterAll(r); err != nil {
		return nil, fmt.Errorf("注册路由失败: %w", err)
	}

	log.Debugf("成功注册 %d 个路由", router.GetRouterCount())
	return r, nil
}
