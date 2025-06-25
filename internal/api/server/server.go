package server

import (
	"context"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/bestruirui/bestsub/internal/api/router"
	"github.com/bestruirui/bestsub/internal/config"
	"github.com/bestruirui/bestsub/internal/utils/log"
	"github.com/gin-gonic/gin"
)

var (
	server *Server
	once   sync.Once
)

// Server HTTP服务器
type Server struct {
	httpServer *http.Server
	router     *gin.Engine
}

// New 创建新的HTTP服务器
func Initialize() error {
	cfg := config.Get()

	// 设置路由
	r := router.SetupRouter()

	// 创建HTTP服务器
	httpServer := &http.Server{
		Addr:           fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port),
		Handler:        r,
		ReadTimeout:    30 * time.Second,
		WriteTimeout:   30 * time.Second,
		IdleTimeout:    60 * time.Second,
		MaxHeaderBytes: 1 << 20, // 1MB
	}

	server = &Server{
		httpServer: httpServer,
		router:     r,
	}

	return nil
}

// Start 启动HTTP服务器
func Start() error {
	log.Infof("Starting HTTP server on %s", server.httpServer.Addr)

	// 启动服务器
	if err := server.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return fmt.Errorf("failed to start HTTP server: %w", err)
	}

	return nil
}

// 关闭HTTP服务器
func Close() error {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := server.httpServer.Shutdown(ctx); err != nil {
		log.Errorf("HTTP server forced to shutdown: %v", err)
		return err
	}
	return nil
}
