package repository

import (
	"github.com/bestruirui/bestsub/internal/database/repository/interfaces"
)

// Repository 统一的仓库接口
type Repository interface {
	// 认证相关
	Auth() interfaces.AuthRepository
	Session() interfaces.SessionRepository

	// 配置相关
	SystemConfig() interfaces.SystemConfigRepository
	NotificationChannel() interfaces.NotificationChannelRepository

	// 任务相关
	Task() interfaces.TaskRepository

	// 订阅配置相关
	SubStorageConfig() interfaces.SubStorageConfigRepository
	SubOutputTemplate() interfaces.SubOutputTemplateRepository
	SubNodeFilterRule() interfaces.SubNodeFilterRuleRepository

	// 订阅链接相关
	SubLink() interfaces.SubLinkRepository
	SubLinkModuleConfig() interfaces.SubLinkModuleConfigRepository

	// 订阅保存相关
	SubSaveConfig() interfaces.SubSaveConfigRepository

	// 订阅分享相关
	SubShareLink() interfaces.SubShareLinkRepository

	// 数据库管理
	Close() error
}
