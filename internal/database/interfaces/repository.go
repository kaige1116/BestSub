package interfaces

import (
	"sync"
)

// Repository 统一的仓库接口
type Repository interface {
	// 认证相关
	Auth() AuthRepository
	Session() SessionRepository

	// 配置相关
	SystemConfig() SystemConfigRepository
	Notify() NotifyRepository

	// 任务相关
	Task() TaskRepository

	// 订阅配置相关
	SubStorageConfig() SubStorageConfigRepository
	SubOutputTemplate() SubOutputTemplateRepository
	SubNodeFilterRule() SubNodeFilterRuleRepository

	// 订阅链接相关
	SubLink() SubLinkRepository

	// 订阅保存相关
	SubSaveConfig() SubSaveConfigRepository

	// 订阅分享相关
	SubShareLink() SubShareLinkRepository

	// 数据库管理
	Close() error
}

// RepositoryManager 仓库管理器，实现懒加载逻辑
type RepositoryManager struct {
	repo Repository
	mu   sync.RWMutex

	// 缓存的仓库实例
	authRepo              AuthRepository
	sessionRepo           SessionRepository
	systemConfigRepo      SystemConfigRepository
	notifyRepo            NotifyRepository
	taskRepo              TaskRepository
	subStorageConfigRepo  SubStorageConfigRepository
	subOutputTemplateRepo SubOutputTemplateRepository
	subNodeFilterRuleRepo SubNodeFilterRuleRepository
	subLinkRepo           SubLinkRepository
	subSaveConfigRepo     SubSaveConfigRepository
	subShareLinkRepo      SubShareLinkRepository
}

// NewRepositoryManager 创建仓库管理器
func NewRepositoryManager(repo Repository) *RepositoryManager {
	return &RepositoryManager{repo: repo}
}

// Auth 获取认证仓库
func (rm *RepositoryManager) Auth() AuthRepository {
	rm.mu.Lock()
	defer rm.mu.Unlock()
	if rm.authRepo == nil {
		rm.authRepo = rm.repo.Auth()
	}
	return rm.authRepo
}

// Session 获取会话仓库
func (rm *RepositoryManager) Session() SessionRepository {
	rm.mu.Lock()
	defer rm.mu.Unlock()
	if rm.sessionRepo == nil {
		rm.sessionRepo = rm.repo.Session()
	}
	return rm.sessionRepo
}

// SystemConfig 获取系统配置仓库
func (rm *RepositoryManager) SystemConfig() SystemConfigRepository {
	rm.mu.Lock()
	defer rm.mu.Unlock()
	if rm.systemConfigRepo == nil {
		rm.systemConfigRepo = rm.repo.SystemConfig()
	}
	return rm.systemConfigRepo
}

// Notify 获取通知渠道仓库
func (rm *RepositoryManager) Notify() NotifyRepository {
	rm.mu.Lock()
	defer rm.mu.Unlock()
	if rm.notifyRepo == nil {
		rm.notifyRepo = rm.repo.Notify()
	}
	return rm.notifyRepo
}

// Task 获取任务仓库
func (rm *RepositoryManager) Task() TaskRepository {
	rm.mu.Lock()
	defer rm.mu.Unlock()
	if rm.taskRepo == nil {
		rm.taskRepo = rm.repo.Task()
	}
	return rm.taskRepo
}

// SubStorageConfig 获取存储配置仓库
func (rm *RepositoryManager) SubStorageConfig() SubStorageConfigRepository {
	rm.mu.Lock()
	defer rm.mu.Unlock()
	if rm.subStorageConfigRepo == nil {
		rm.subStorageConfigRepo = rm.repo.SubStorageConfig()
	}
	return rm.subStorageConfigRepo
}

// SubOutputTemplate 获取输出模板仓库
func (rm *RepositoryManager) SubOutputTemplate() SubOutputTemplateRepository {
	rm.mu.Lock()
	defer rm.mu.Unlock()
	if rm.subOutputTemplateRepo == nil {
		rm.subOutputTemplateRepo = rm.repo.SubOutputTemplate()
	}
	return rm.subOutputTemplateRepo
}

// SubNodeFilterRule 获取节点筛选规则仓库
func (rm *RepositoryManager) SubNodeFilterRule() SubNodeFilterRuleRepository {
	rm.mu.Lock()
	defer rm.mu.Unlock()
	if rm.subNodeFilterRuleRepo == nil {
		rm.subNodeFilterRuleRepo = rm.repo.SubNodeFilterRule()
	}
	return rm.subNodeFilterRuleRepo
}

// SubLink 获取订阅链接仓库
func (rm *RepositoryManager) SubLink() SubLinkRepository {
	rm.mu.Lock()
	defer rm.mu.Unlock()
	if rm.subLinkRepo == nil {
		rm.subLinkRepo = rm.repo.SubLink()
	}
	return rm.subLinkRepo
}

// SubSaveConfig 获取保存配置仓库
func (rm *RepositoryManager) SubSaveConfig() SubSaveConfigRepository {
	rm.mu.Lock()
	defer rm.mu.Unlock()
	if rm.subSaveConfigRepo == nil {
		rm.subSaveConfigRepo = rm.repo.SubSaveConfig()
	}
	return rm.subSaveConfigRepo
}

// SubShareLink 获取分享链接仓库
func (rm *RepositoryManager) SubShareLink() SubShareLinkRepository {
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
