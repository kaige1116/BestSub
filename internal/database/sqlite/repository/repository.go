package repository

import (
	"github.com/bestruirui/bestsub/internal/database/interfaces"
	"github.com/bestruirui/bestsub/internal/database/sqlite/database"
)

// Repository SQLite仓库实现
type Repository struct {
	db *database.Database
}

// NewRepository 创建新的SQLite仓库
func NewRepository(db *database.Database) *Repository {
	return &Repository{db: db}
}

// Auth 获取认证仓库
func (r *Repository) Auth() interfaces.AuthRepository {
	return newAuthRepository(r.db)
}

// SystemConfig 获取系统配置仓库
func (r *Repository) SystemConfig() interfaces.SystemConfigRepository {
	return newSystemConfigRepository(r.db)
}

// Notify 获取通知渠道仓库
func (r *Repository) Notify() interfaces.NotifyRepository {
	return newNotifyRepository(r.db)
}

// NotifyTemplate 获取通知模板仓库
func (r *Repository) NotifyTemplate() interfaces.NotifyTemplateRepository {
	return newNotifyTemplateRepository(r.db)
}

// Task 获取任务仓库
func (r *Repository) Task() interfaces.TaskRepository {
	return newTaskRepository(r.db)
}

// SubStorageConfig 获取存储配置仓库
func (r *Repository) SubStorageConfig() interfaces.SubStorageConfigRepository {
	return newSubStorageConfigRepository(r.db)
}

// SubOutputTemplate 获取输出模板仓库
func (r *Repository) SubOutputTemplate() interfaces.SubOutputTemplateRepository {
	return newSubOutputTemplateRepository(r.db)
}

// SubNodeFilterRule 获取节点筛选规则仓库
func (r *Repository) SubNodeFilterRule() interfaces.SubNodeFilterRuleRepository {
	return newSubNodeFilterRuleRepository(r.db)
}

// SubLink 获取订阅链接仓库
func (r *Repository) SubLink() interfaces.SubRepository {
	return newSubRepository(r.db)
}

// SubSaveConfig 获取保存配置仓库
func (r *Repository) SubSaveConfig() interfaces.SubSaveRepository {
	return newSubSaveConfigRepository(r.db)
}

// SubShareLink 获取分享链接仓库
func (r *Repository) SubShareLink() interfaces.SubShareRepository {
	return newSubShareLinkRepository(r.db)
}

// Close 关闭数据库连接
func (r *Repository) Close() error {
	return r.db.Close()
}

// Database 获取数据库实例（用于健康检查等）
func (r *Repository) Database() *database.Database {
	return r.db
}
