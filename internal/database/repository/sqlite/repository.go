package sqlite

import (
	"github.com/bestruirui/bestsub/internal/database/repository"
	"github.com/bestruirui/bestsub/internal/database/repository/interfaces"
)

// Repository SQLite仓库工厂
type Repository struct {
	db *Database

	// 认证相关
	authRepo    interfaces.AuthRepository
	sessionRepo interfaces.SessionRepository

	// 配置相关
	systemConfigRepo        interfaces.SystemConfigRepository
	notificationChannelRepo interfaces.NotificationChannelRepository

	// 任务相关
	taskRepo interfaces.TaskRepository

	// 订阅配置相关
	subStorageConfigRepo  interfaces.SubStorageConfigRepository
	subOutputTemplateRepo interfaces.SubOutputTemplateRepository
	subNodeFilterRuleRepo interfaces.SubNodeFilterRuleRepository

	// 订阅链接相关
	subLinkRepo             interfaces.SubLinkRepository
	subLinkModuleConfigRepo interfaces.SubLinkModuleConfigRepository

	// 订阅保存相关
	subSaveConfigRepo interfaces.SubSaveConfigRepository

	// 订阅分享相关
	subShareLinkRepo interfaces.SubShareLinkRepository
}

// 编译时检查是否实现了接口
var _ repository.Repository = (*Repository)(nil)

// NewRepository 创建新的SQLite仓库工厂
func NewRepository(db *Database) repository.Repository {
	return &Repository{
		db: db,
	}
}

// Auth 获取认证仓库
func (r *Repository) Auth() interfaces.AuthRepository {
	if r.authRepo == nil {
		r.authRepo = NewAuthRepository(r.db)
	}
	return r.authRepo
}

// Session 获取会话仓库
func (r *Repository) Session() interfaces.SessionRepository {
	if r.sessionRepo == nil {
		r.sessionRepo = NewSessionRepository(r.db)
	}
	return r.sessionRepo
}

// SystemConfig 获取系统配置仓库
func (r *Repository) SystemConfig() interfaces.SystemConfigRepository {
	if r.systemConfigRepo == nil {
		r.systemConfigRepo = NewSystemConfigRepository(r.db)
	}
	return r.systemConfigRepo
}

// NotificationChannel 获取通知渠道仓库
func (r *Repository) NotificationChannel() interfaces.NotificationChannelRepository {
	if r.notificationChannelRepo == nil {
		r.notificationChannelRepo = NewNotificationChannelRepository(r.db)
	}
	return r.notificationChannelRepo
}

// Task 获取任务仓库
func (r *Repository) Task() interfaces.TaskRepository {
	if r.taskRepo == nil {
		r.taskRepo = NewTaskRepository(r.db)
	}
	return r.taskRepo
}

// SubStorageConfig 获取存储配置仓库
func (r *Repository) SubStorageConfig() interfaces.SubStorageConfigRepository {
	if r.subStorageConfigRepo == nil {
		r.subStorageConfigRepo = NewSubStorageConfigRepository(r.db)
	}
	return r.subStorageConfigRepo
}

// SubOutputTemplate 获取输出模板仓库
func (r *Repository) SubOutputTemplate() interfaces.SubOutputTemplateRepository {
	if r.subOutputTemplateRepo == nil {
		r.subOutputTemplateRepo = NewSubOutputTemplateRepository(r.db)
	}
	return r.subOutputTemplateRepo
}

// SubNodeFilterRule 获取节点筛选规则仓库
func (r *Repository) SubNodeFilterRule() interfaces.SubNodeFilterRuleRepository {
	if r.subNodeFilterRuleRepo == nil {
		r.subNodeFilterRuleRepo = NewSubNodeFilterRuleRepository(r.db)
	}
	return r.subNodeFilterRuleRepo
}

// SubLink 获取订阅链接仓库
func (r *Repository) SubLink() interfaces.SubLinkRepository {
	if r.subLinkRepo == nil {
		r.subLinkRepo = NewSubLinkRepository(r.db)
	}
	return r.subLinkRepo
}

// SubLinkModuleConfig 获取链接模块配置仓库
func (r *Repository) SubLinkModuleConfig() interfaces.SubLinkModuleConfigRepository {
	if r.subLinkModuleConfigRepo == nil {
		r.subLinkModuleConfigRepo = NewSubLinkModuleConfigRepository(r.db)
	}
	return r.subLinkModuleConfigRepo
}

// SubSaveConfig 获取保存配置仓库
func (r *Repository) SubSaveConfig() interfaces.SubSaveConfigRepository {
	if r.subSaveConfigRepo == nil {
		r.subSaveConfigRepo = NewSubSaveConfigRepository(r.db)
	}
	return r.subSaveConfigRepo
}

// SubShareLink 获取分享链接仓库
func (r *Repository) SubShareLink() interfaces.SubShareLinkRepository {
	if r.subShareLinkRepo == nil {
		r.subShareLinkRepo = NewSubShareLinkRepository(r.db)
	}
	return r.subShareLinkRepo
}

// Database 获取数据库连接
func (r *Repository) Database() *Database {
	return r.db
}

// Close 关闭数据库连接
func (r *Repository) Close() error {
	return r.db.Close()
}
