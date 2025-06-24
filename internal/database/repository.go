package database

import (
	"sync"

	"github.com/bestruirui/bestsub/internal/database/interfaces"
)

// RepositoryManager 仓库管理器，实现懒加载逻辑
type RepositoryManager struct {
	repo interfaces.Repository
	mu   sync.RWMutex

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

// NewRepositoryManager 创建仓库管理器
func NewRepositoryManager(repo interfaces.Repository) *RepositoryManager {
	return &RepositoryManager{repo: repo}
}

// Auth 获取认证仓库
func (rm *RepositoryManager) Auth() interfaces.AuthRepository {
	rm.mu.Lock()
	defer rm.mu.Unlock()
	if rm.authRepo == nil {
		rm.authRepo = rm.repo.Auth()
	}
	return rm.authRepo
}

// Session 获取会话仓库
func (rm *RepositoryManager) Session() interfaces.SessionRepository {
	rm.mu.Lock()
	defer rm.mu.Unlock()
	if rm.sessionRepo == nil {
		rm.sessionRepo = rm.repo.Session()
	}
	return rm.sessionRepo
}

// SystemConfig 获取系统配置仓库
func (rm *RepositoryManager) SystemConfig() interfaces.SystemConfigRepository {
	rm.mu.Lock()
	defer rm.mu.Unlock()
	if rm.systemConfigRepo == nil {
		rm.systemConfigRepo = rm.repo.SystemConfig()
	}
	return rm.systemConfigRepo
}

// NotificationChannel 获取通知渠道仓库
func (rm *RepositoryManager) NotificationChannel() interfaces.NotificationChannelRepository {
	rm.mu.Lock()
	defer rm.mu.Unlock()
	if rm.notificationChannelRepo == nil {
		rm.notificationChannelRepo = rm.repo.NotificationChannel()
	}
	return rm.notificationChannelRepo
}

// Task 获取任务仓库
func (rm *RepositoryManager) Task() interfaces.TaskRepository {
	rm.mu.Lock()
	defer rm.mu.Unlock()
	if rm.taskRepo == nil {
		rm.taskRepo = rm.repo.Task()
	}
	return rm.taskRepo
}

// SubStorageConfig 获取存储配置仓库
func (rm *RepositoryManager) SubStorageConfig() interfaces.SubStorageConfigRepository {
	rm.mu.Lock()
	defer rm.mu.Unlock()
	if rm.subStorageConfigRepo == nil {
		rm.subStorageConfigRepo = rm.repo.SubStorageConfig()
	}
	return rm.subStorageConfigRepo
}

// SubOutputTemplate 获取输出模板仓库
func (rm *RepositoryManager) SubOutputTemplate() interfaces.SubOutputTemplateRepository {
	rm.mu.Lock()
	defer rm.mu.Unlock()
	if rm.subOutputTemplateRepo == nil {
		rm.subOutputTemplateRepo = rm.repo.SubOutputTemplate()
	}
	return rm.subOutputTemplateRepo
}

// SubNodeFilterRule 获取节点筛选规则仓库
func (rm *RepositoryManager) SubNodeFilterRule() interfaces.SubNodeFilterRuleRepository {
	rm.mu.Lock()
	defer rm.mu.Unlock()
	if rm.subNodeFilterRuleRepo == nil {
		rm.subNodeFilterRuleRepo = rm.repo.SubNodeFilterRule()
	}
	return rm.subNodeFilterRuleRepo
}

// SubLink 获取订阅链接仓库
func (rm *RepositoryManager) SubLink() interfaces.SubLinkRepository {
	rm.mu.Lock()
	defer rm.mu.Unlock()
	if rm.subLinkRepo == nil {
		rm.subLinkRepo = rm.repo.SubLink()
	}
	return rm.subLinkRepo
}

// SubLinkModuleConfig 获取链接模块配置仓库
func (rm *RepositoryManager) SubLinkModuleConfig() interfaces.SubLinkModuleConfigRepository {
	rm.mu.Lock()
	defer rm.mu.Unlock()
	if rm.subLinkModuleConfigRepo == nil {
		rm.subLinkModuleConfigRepo = rm.repo.SubLinkModuleConfig()
	}
	return rm.subLinkModuleConfigRepo
}

// SubSaveConfig 获取保存配置仓库
func (rm *RepositoryManager) SubSaveConfig() interfaces.SubSaveConfigRepository {
	rm.mu.Lock()
	defer rm.mu.Unlock()
	if rm.subSaveConfigRepo == nil {
		rm.subSaveConfigRepo = rm.repo.SubSaveConfig()
	}
	return rm.subSaveConfigRepo
}

// SubShareLink 获取分享链接仓库
func (rm *RepositoryManager) SubShareLink() interfaces.SubShareLinkRepository {
	rm.mu.Lock()
	defer rm.mu.Unlock()
	if rm.subShareLinkRepo == nil {
		rm.subShareLinkRepo = rm.repo.SubShareLink()
	}
	return rm.subShareLinkRepo
}

// Close 关闭数据库连接
func (rm *RepositoryManager) Close() error {
	return rm.repo.Close()
}

// =============================================================================
// 快捷访问方法 - 直接访问仓库实例的便捷函数
// =============================================================================

// Auth 获取认证仓库
func Auth() interfaces.AuthRepository {
	return GetRepository().Auth()
}

// Session 获取会话仓库
func Session() interfaces.SessionRepository {
	return GetRepository().Session()
}

// SystemConfig 获取系统配置仓库
func SystemConfig() interfaces.SystemConfigRepository {
	return GetRepository().SystemConfig()
}

// NotificationChannel 获取通知渠道仓库
func NotificationChannel() interfaces.NotificationChannelRepository {
	return GetRepository().NotificationChannel()
}

// Task 获取任务仓库
func Task() interfaces.TaskRepository {
	return GetRepository().Task()
}

// SubStorageConfig 获取存储配置仓库
func SubStorageConfig() interfaces.SubStorageConfigRepository {
	return GetRepository().SubStorageConfig()
}

// SubOutputTemplate 获取输出模板仓库
func SubOutputTemplate() interfaces.SubOutputTemplateRepository {
	return GetRepository().SubOutputTemplate()
}

// SubNodeFilterRule 获取节点筛选规则仓库
func SubNodeFilterRule() interfaces.SubNodeFilterRuleRepository {
	return GetRepository().SubNodeFilterRule()
}

// SubLink 获取订阅链接仓库
func SubLink() interfaces.SubLinkRepository {
	return GetRepository().SubLink()
}

// SubLinkModuleConfig 获取链接模块配置仓库
func SubLinkModuleConfig() interfaces.SubLinkModuleConfigRepository {
	return GetRepository().SubLinkModuleConfig()
}

// SubSaveConfig 获取保存配置仓库
func SubSaveConfig() interfaces.SubSaveConfigRepository {
	return GetRepository().SubSaveConfig()
}

// SubShareLink 获取分享链接仓库
func SubShareLink() interfaces.SubShareLinkRepository {
	return GetRepository().SubShareLink()
}
