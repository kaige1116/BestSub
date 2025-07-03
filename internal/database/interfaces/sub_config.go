package interfaces

import (
	"context"

	"github.com/bestruirui/bestsub/internal/database/models"
)

// 存储配置数据访问接口
type SubStorageConfigRepository interface {
	// Create 创建存储配置
	Create(ctx context.Context, config *models.SubStorageConfig) error

	// GetByID 根据ID获取存储配置
	GetByID(ctx context.Context, id int64) (*models.SubStorageConfig, error)

	// Update 更新存储配置
	Update(ctx context.Context, config *models.SubStorageConfig) error

	// Delete 删除存储配置
	Delete(ctx context.Context, id int64) error

	// 根据保存ID获取存储配置
	GetBySaveID(ctx context.Context, saveID int64) (*models.SubStorageConfig, error)
}

// 输出模板数据访问接口
type SubOutputTemplateRepository interface {
	// Create 创建输出模板
	Create(ctx context.Context, template *models.SubOutputTemplate) error

	// GetByID 根据ID获取输出模板
	GetByID(ctx context.Context, id int64) (*models.SubOutputTemplate, error)

	// Update 更新输出模板
	Update(ctx context.Context, template *models.SubOutputTemplate) error

	// Delete 删除输出模板
	Delete(ctx context.Context, id int64) error

	// 根据任务ID获取输出模板
	GetByShareID(ctx context.Context, shareID int64) (*models.SubOutputTemplate, error)

	// 根据保存ID获取输出模板
	GetBySaveID(ctx context.Context, saveID int64) (*models.SubOutputTemplate, error)
}

// 节点筛选规则数据访问接口
type SubNodeFilterRuleRepository interface {
	// Create 创建筛选规则
	Create(ctx context.Context, rule *models.SubNodeFilterRule) error

	// GetByID 根据ID获取筛选规则
	GetByID(ctx context.Context, id int64) (*models.SubNodeFilterRule, error)

	// Update 更新筛选规则
	Update(ctx context.Context, rule *models.SubNodeFilterRule) error

	// Delete 删除筛选规则
	Delete(ctx context.Context, id int64) error

	// 根据保存ID获取筛选规则
	GetBySaveID(ctx context.Context, saveID int64) (*models.SubNodeFilterRule, error)

	// 根据分享ID获取筛选规则
	GetByShareID(ctx context.Context, shareID int64) (*models.SubNodeFilterRule, error)
}
