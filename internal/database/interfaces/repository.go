package interfaces

// Repository 统一的仓库接口
type Repository interface {
	Auth() AuthRepository

	Config() ConfigRepository

	Notify() NotifyRepository
	NotifyTemplate() NotifyTemplateRepository

	Task() TaskRepository

	Sub() SubRepository
	SubShare() SubShareRepository
	SubTemplate() SubTemplateRepository

	Storage() StorageRepository

	Close() error
	Migrate() error
}
