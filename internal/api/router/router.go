// Package router provides HTTP routing configuration
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

package router

import (
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	_ "github.com/bestruirui/bestsub/docs" // 导入生成的swagger文档
	"github.com/bestruirui/bestsub/internal/api/handlers"
	"github.com/bestruirui/bestsub/internal/api/middleware"
	"github.com/bestruirui/bestsub/internal/config"
)

// SetupRouter 设置路由
func SetupRouter() *gin.Engine {
	// 根据配置设置Gin模式
	cfg := config.Get()
	if cfg.Log.Level == "debug" {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.New()

	// 添加中间件
	r.Use(middleware.Logging())
	r.Use(middleware.Recovery())
	r.Use(middleware.Cors())

	// 创建处理器
	authHandler := handlers.NewAuthHandler()
	healthHandler := handlers.NewHealthHandler()

	// 健康检查路由
	r.GET("/health", healthHandler.HealthCheck)
	r.GET("/ready", healthHandler.ReadinessCheck)
	r.GET("/live", healthHandler.LivenessCheck)

	// Swagger文档路由
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, ginSwagger.URL("/swagger/doc.json")))

	// API v1 路由组
	v1 := r.Group("/api/v1")
	{
		// 认证相关路由（无需认证）
		authGroup := v1.Group("/auth")
		{
			authGroup.POST("/login", authHandler.Login)
			authGroup.POST("/refresh", authHandler.RefreshToken)
		}

		// 需要认证的路由
		protectedGroup := v1.Group("/auth")
		protectedGroup.Use(middleware.Auth())
		{
			protectedGroup.POST("/logout", authHandler.Logout)
			protectedGroup.POST("/change-password", authHandler.ChangePassword)
			protectedGroup.POST("/update-username", authHandler.UpdateUsername)
			protectedGroup.GET("/user", authHandler.GetUserInfo)
			protectedGroup.GET("/sessions", authHandler.GetSessions)
			protectedGroup.DELETE("/sessions/:id", authHandler.DeleteSession)
		}

		// 其他API路由组可以在这里添加
		// 例如：
		// configGroup := v1.Group("/config")
		// configGroup.Use(auth.AuthMiddleware())
		// {
		//     // 配置相关路由
		// }
	}

	return r
}
