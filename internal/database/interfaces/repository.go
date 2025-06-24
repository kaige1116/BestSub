package interfaces

// Repository 统一的仓库接口
type Repository interface {
	// 认证相关
	Auth() AuthRepository
	Session() SessionRepository

	// 配置相关
	SystemConfig() SystemConfigRepository
	NotificationChannel() NotificationChannelRepository

	// 任务相关
	Task() TaskRepository

	// 订阅配置相关
	SubStorageConfig() SubStorageConfigRepository
	SubOutputTemplate() SubOutputTemplateRepository
	SubNodeFilterRule() SubNodeFilterRuleRepository

	// 订阅链接相关
	SubLink() SubLinkRepository
	SubLinkModuleConfig() SubLinkModuleConfigRepository

	// 订阅保存相关
	SubSaveConfig() SubSaveConfigRepository

	// 订阅分享相关
	SubShareLink() SubShareLinkRepository

	// 数据库管理
	Close() error
}
