package repository

import (
	"sync"

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

// RepositoryFactory 仓库工厂接口，用于创建具体的仓库实例
type RepositoryFactory interface {
	CreateAuthRepository() interfaces.AuthRepository
	CreateSessionRepository() interfaces.SessionRepository
	CreateSystemConfigRepository() interfaces.SystemConfigRepository
	CreateNotificationChannelRepository() interfaces.NotificationChannelRepository
	CreateTaskRepository() interfaces.TaskRepository
	CreateSubStorageConfigRepository() interfaces.SubStorageConfigRepository
	CreateSubOutputTemplateRepository() interfaces.SubOutputTemplateRepository
	CreateSubNodeFilterRuleRepository() interfaces.SubNodeFilterRuleRepository
	CreateSubLinkRepository() interfaces.SubLinkRepository
	CreateSubLinkModuleConfigRepository() interfaces.SubLinkModuleConfigRepository
	CreateSubSaveConfigRepository() interfaces.SubSaveConfigRepository
	CreateSubShareLinkRepository() interfaces.SubShareLinkRepository
}

// BaseRepository 基础仓库实现，提供通用的懒加载逻辑
type BaseRepository struct {
	factory RepositoryFactory
	mu      sync.RWMutex

	// 缓存的仓库实例
	authRepo                interfaces.AuthRepository
	sessionRepo             interfaces.SessionRepository
	systemConfigRepo        interfaces.SystemConfigRepository
	notificationChannelRepo interfaces.NotificationChannelRepository
	taskRepo                interfaces.TaskRepository
	subStorageConfigRepo    interfaces.SubStorageConfigRepository
	subOutputTemplateRepo   interfaces.SubOutputTemplateRepository
	subNodeFilterRuleRepo   interfaces.SubNodeFilterRuleRepository
	subLinkRepo             interfaces.SubLinkRepository
	subLinkModuleConfigRepo interfaces.SubLinkModuleConfigRepository
	subSaveConfigRepo       interfaces.SubSaveConfigRepository
	subShareLinkRepo        interfaces.SubShareLinkRepository
}

// NewBaseRepository 创建基础仓库
func NewBaseRepository(factory RepositoryFactory) *BaseRepository {
	return &BaseRepository{
		factory: factory,
	}
}

// Auth 获取认证仓库
func (r *BaseRepository) Auth() interfaces.AuthRepository {
	r.mu.Lock()
	defer r.mu.Unlock()
	if r.authRepo == nil {
		r.authRepo = r.factory.CreateAuthRepository()
	}
	return r.authRepo
}

// Session 获取会话仓库
func (r *BaseRepository) Session() interfaces.SessionRepository {
	r.mu.Lock()
	defer r.mu.Unlock()
	if r.sessionRepo == nil {
		r.sessionRepo = r.factory.CreateSessionRepository()
	}
	return r.sessionRepo
}

// SystemConfig 获取系统配置仓库
func (r *BaseRepository) SystemConfig() interfaces.SystemConfigRepository {
	r.mu.Lock()
	defer r.mu.Unlock()
	if r.systemConfigRepo == nil {
		r.systemConfigRepo = r.factory.CreateSystemConfigRepository()
	}
	return r.systemConfigRepo
}

// NotificationChannel 获取通知渠道仓库
func (r *BaseRepository) NotificationChannel() interfaces.NotificationChannelRepository {
	r.mu.Lock()
	defer r.mu.Unlock()
	if r.notificationChannelRepo == nil {
		r.notificationChannelRepo = r.factory.CreateNotificationChannelRepository()
	}
	return r.notificationChannelRepo
}

// Task 获取任务仓库
func (r *BaseRepository) Task() interfaces.TaskRepository {
	r.mu.Lock()
	defer r.mu.Unlock()
	if r.taskRepo == nil {
		r.taskRepo = r.factory.CreateTaskRepository()
	}
	return r.taskRepo
}

// SubStorageConfig 获取存储配置仓库
func (r *BaseRepository) SubStorageConfig() interfaces.SubStorageConfigRepository {
	r.mu.Lock()
	defer r.mu.Unlock()
	if r.subStorageConfigRepo == nil {
		r.subStorageConfigRepo = r.factory.CreateSubStorageConfigRepository()
	}
	return r.subStorageConfigRepo
}

// SubOutputTemplate 获取输出模板仓库
func (r *BaseRepository) SubOutputTemplate() interfaces.SubOutputTemplateRepository {
	r.mu.Lock()
	defer r.mu.Unlock()
	if r.subOutputTemplateRepo == nil {
		r.subOutputTemplateRepo = r.factory.CreateSubOutputTemplateRepository()
	}
	return r.subOutputTemplateRepo
}

// SubNodeFilterRule 获取节点筛选规则仓库
func (r *BaseRepository) SubNodeFilterRule() interfaces.SubNodeFilterRuleRepository {
	r.mu.Lock()
	defer r.mu.Unlock()
	if r.subNodeFilterRuleRepo == nil {
		r.subNodeFilterRuleRepo = r.factory.CreateSubNodeFilterRuleRepository()
	}
	return r.subNodeFilterRuleRepo
}

// SubLink 获取订阅链接仓库
func (r *BaseRepository) SubLink() interfaces.SubLinkRepository {
	r.mu.Lock()
	defer r.mu.Unlock()
	if r.subLinkRepo == nil {
		r.subLinkRepo = r.factory.CreateSubLinkRepository()
	}
	return r.subLinkRepo
}

// SubLinkModuleConfig 获取链接模块配置仓库
func (r *BaseRepository) SubLinkModuleConfig() interfaces.SubLinkModuleConfigRepository {
	r.mu.Lock()
	defer r.mu.Unlock()
	if r.subLinkModuleConfigRepo == nil {
		r.subLinkModuleConfigRepo = r.factory.CreateSubLinkModuleConfigRepository()
	}
	return r.subLinkModuleConfigRepo
}

// SubSaveConfig 获取保存配置仓库
func (r *BaseRepository) SubSaveConfig() interfaces.SubSaveConfigRepository {
	r.mu.Lock()
	defer r.mu.Unlock()
	if r.subSaveConfigRepo == nil {
		r.subSaveConfigRepo = r.factory.CreateSubSaveConfigRepository()
	}
	return r.subSaveConfigRepo
}

// SubShareLink 获取分享链接仓库
func (r *BaseRepository) SubShareLink() interfaces.SubShareLinkRepository {
	r.mu.Lock()
	defer r.mu.Unlock()
	if r.subShareLinkRepo == nil {
		r.subShareLinkRepo = r.factory.CreateSubShareLinkRepository()
	}
	return r.subShareLinkRepo
}
