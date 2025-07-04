package interfaces

import (
	"context"

	"github.com/bestruirui/bestsub/internal/models/sub"
)

// 存储配置数据访问接口
type SubStorageConfigRepository interface {
	// Create 创建存储配置
	Create(ctx context.Context, config *sub.StorageConfig) error

	// GetByID 根据ID获取存储配置
	GetByID(ctx context.Context, id int64) (*sub.StorageConfig, error)

	// Update 更新存储配置
	Update(ctx context.Context, config *sub.StorageConfig) error

	// Delete 删除存储配置
	Delete(ctx context.Context, id int64) error

	// 根据保存ID获取存储配置列表
	GetBySaveID(ctx context.Context, saveID int64) (*[]sub.StorageConfig, error)

	// 添加存储配置与保存的关联
	AddSaveRelation(ctx context.Context, configID, saveID int64) error
}

// 输出模板数据访问接口
type SubOutputTemplateRepository interface {
	// Create 创建输出模板
	Create(ctx context.Context, template *sub.OutputTemplate) error

	// GetByID 根据ID获取输出模板
	GetByID(ctx context.Context, id int64) (*sub.OutputTemplate, error)

	// Update 更新输出模板
	Update(ctx context.Context, template *sub.OutputTemplate) error

	// Delete 删除输出模板
	Delete(ctx context.Context, id int64) error

	// 根据保存ID获取输出模板列表
	GetBySaveID(ctx context.Context, saveID int64) (*sub.OutputTemplate, error)

	// 根据分享ID获取输出模板列表
	GetByShareID(ctx context.Context, shareID int64) (*sub.OutputTemplate, error)

	// 添加输出模板与分享的关联
	AddShareRelation(ctx context.Context, templateID, shareID int64) error

	// 添加输出模板与保存的关联
	AddSaveRelation(ctx context.Context, templateID, saveID int64) error
}

// 节点筛选规则数据访问接口
type SubNodeFilterRuleRepository interface {
	// Create 创建筛选规则
	Create(ctx context.Context, rule *sub.NodeFilterRule) error

	// Update 更新筛选规则
	Update(ctx context.Context, rule *sub.NodeFilterRule) error

	// GetByID 根据ID获取筛选规则
	GetByID(ctx context.Context, id int64) (*sub.NodeFilterRule, error)

	// Delete 删除筛选规则
	Delete(ctx context.Context, id int64) error

	// 根据保存ID获取筛选规则
	GetBySaveID(ctx context.Context, saveID int64) (*[]sub.NodeFilterRule, error)

	// 根据分享ID获取筛选规则
	GetByShareID(ctx context.Context, shareID int64) (*[]sub.NodeFilterRule, error)

	// 添加筛选规则与分享的关联
	AddShareRelation(ctx context.Context, ruleID, shareID int64) error

	// 添加筛选规则与保存的关联
	AddSaveRelation(ctx context.Context, ruleID, saveID int64) error
}
