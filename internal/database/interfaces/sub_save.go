package interfaces

import (
	"context"

	"github.com/bestruirui/bestsub/internal/models/sub"
)

// 保存配置数据访问接口
type SubSaveRepository interface {
	// Create 创建保存配置
	Create(ctx context.Context, config *sub.SaveConfig) error

	// GetByID 根据ID获取保存配置
	GetByID(ctx context.Context, id uint16) (*sub.SaveConfig, error)

	// Update 更新保存配置
	Update(ctx context.Context, config *sub.SaveConfig) error

	// Delete 删除保存配置
	Delete(ctx context.Context, id uint16) error

	// 根据任务ID获取保存配置列表
	GetByTaskID(ctx context.Context, taskID uint16) (*[]sub.SaveConfig, error)

	// 添加保存配置与任务的关联
	AddTaskRelation(ctx context.Context, saveID, taskID uint16) error

	// 添加保存配置与输出模板的关联
	AddOutputTemplateRelation(ctx context.Context, saveID, templateID uint16) error

	// 添加保存配置与过滤配置的关联
	AddFilterConfigRelation(ctx context.Context, saveID, configID uint16) error

	// 添加保存配置与订阅的关联
	AddSubRelation(ctx context.Context, saveID, subID uint16) error

	// 添加保存配置与存储配置的关联
	AddStorageConfigRelation(ctx context.Context, saveID, configID uint16) error
}
