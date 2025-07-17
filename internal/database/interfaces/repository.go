package interfaces

// Repository 统一的仓库接口
type Repository interface {
	Auth() AuthRepository

	Config() ConfigRepository

	Notify() NotifyRepository
	NotifyTemplate() NotifyTemplateRepository

	Task() TaskRepository

	Sub() SubRepository
	SubSaveConfig() SubSaveRepository
	SubShareLink() SubShareRepository
	SubStorageConfig() SubStorageConfigRepository
	SubOutputTemplate() SubOutputTemplateRepository
	SubNodeFilterRule() SubNodeFilterRuleRepository

	Close() error
	Migrate() error
}
