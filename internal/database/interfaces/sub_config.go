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

	// List 获取存储配置列表
	List(ctx context.Context, offset, limit int) ([]*models.SubStorageConfig, error)

	// ListActive 获取活跃的存储配置列表
	ListActive(ctx context.Context) ([]*models.SubStorageConfig, error)

	// ListByType 根据类型获取存储配置列表
	ListByType(ctx context.Context, storageType string) ([]*models.SubStorageConfig, error)

	// Count 获取存储配置总数
	Count(ctx context.Context) (int64, error)

	// UpdateTestResult 更新测试结果
	UpdateTestResult(ctx context.Context, id int64, testResult string) error
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

	// List 获取输出模板列表
	List(ctx context.Context, offset, limit int) ([]*models.SubOutputTemplate, error)

	// ListActive 获取活跃的输出模板列表
	ListActive(ctx context.Context) ([]*models.SubOutputTemplate, error)

	// ListByFormat 根据格式获取输出模板列表
	ListByFormat(ctx context.Context, format string) ([]*models.SubOutputTemplate, error)

	// GetDefault 获取默认模板
	GetDefault(ctx context.Context, format string) (*models.SubOutputTemplate, error)

	// SetDefault 设置默认模板
	SetDefault(ctx context.Context, id int64, format string) error

	// Count 获取输出模板总数
	Count(ctx context.Context) (int64, error)
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

	// List 获取筛选规则列表
	List(ctx context.Context, offset, limit int) ([]*models.SubNodeFilterRule, error)

	// ListEnabled 获取启用的筛选规则列表
	ListEnabled(ctx context.Context) ([]*models.SubNodeFilterRule, error)

	// ListByType 根据类型获取筛选规则列表
	ListByType(ctx context.Context, ruleType string) ([]*models.SubNodeFilterRule, error)

	// Count 获取筛选规则总数
	Count(ctx context.Context) (int64, error)

	// UpdatePriority 更新优先级
	UpdatePriority(ctx context.Context, id int64, priority int) error
}
