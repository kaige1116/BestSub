// Package server 提供 BestSub 应用程序的入口点。
//
// @title BestSub API
// @version 1.0.0
// @description BestSub -  API 文档
// @description
// @description 这是 BestSub 的 API 文档
// @description
// @description ## 认证
// @description 大多数接口需要使用 JWT 令牌进行认证。
// @description 认证时，请在 Authorization 头中包含 JWT 令牌：
// @description `Authorization: Bearer <your-jwt-token>`
// @description
// @description ## 错误响应
// @description 所有错误响应都遵循统一格式，包含 code、message 和 error 字段。
// @description
// @description ## 成功响应
// @description 所有成功响应都遵循统一格式，包含 code、message 和 data 字段。
//
// @contact.name BestSub API 支持
// @contact.email support@bestsub.com
//
// @license.name GPL-3.0
// @license.url https://opensource.org/license/gpl-3-0
//
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description 类型为 "Bearer"，后跟空格和 JWT 令牌。
//
// @tag.name 认证
// @tag.description 用户认证相关接口
//
// @tag.name 系统
// @tag.description 系统状态和健康检查接口
package server

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/bestruirui/bestsub/internal/config"
	"github.com/bestruirui/bestsub/internal/models/system"
	_ "github.com/bestruirui/bestsub/internal/server/handlers"
	"github.com/bestruirui/bestsub/internal/server/middleware"
	"github.com/bestruirui/bestsub/internal/server/router"
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

var server *Server

// Server HTTP服务器
type Server struct {
	httpServer *http.Server
	router     *gin.Engine
	config     system.Config
}

// 初始化HTTP服务器
func Initialize() error {
	cfg := config.Base()

	r, routerErr := setRouter()
	if routerErr != nil {
		return fmt.Errorf("设置路由失败: %w", routerErr)
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
	return nil
}

// 启动HTTP服务器
func Start() error {
	if server == nil {
		return fmt.Errorf("HTTP 服务器未初始化, 请先调用 Initialize()")
	}

	log.Infof("启动 HTTP 服务器 %s", server.httpServer.Addr)

	go func() {
		if err := server.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Errorf("启动 HTTP 服务器失败: %v", err)
		}
	}()

	return nil
}

// 关闭HTTP服务器
func Close() error {
	if server == nil {
		return fmt.Errorf("HTTP 服务器未初始化")
	}

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
