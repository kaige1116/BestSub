package interfaces

// Repository 统一的仓库接口
type Repository interface {
	Auth() AuthRepository

	Config() ConfigRepository

	Notify() NotifyRepository
	NotifyTemplate() NotifyTemplateRepository

	Task() TaskRepository

	Sub() SubRepository
	SubSave() SubSaveRepository
	SubShare() SubShareRepository
	SubStorage() SubStorageConfigRepository
	SubOutputTemplate() SubOutputTemplateRepository
	SubNodeFilterRule() SubNodeFilterRuleRepository

	Close() error
	Migrate() error
}
