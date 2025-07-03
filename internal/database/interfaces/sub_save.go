package interfaces

import (
	"context"

	"github.com/bestruirui/bestsub/internal/database/models"
)

// 保存配置数据访问接口
type SubSaveConfigRepository interface {
	// Create 创建保存配置
	Create(ctx context.Context, config *models.SubSaveConfig) error

	// GetByID 根据ID获取保存配置
	GetByID(ctx context.Context, id int64) (*models.SubSaveConfig, error)

	// Update 更新保存配置
	Update(ctx context.Context, config *models.SubSaveConfig) error

	// Delete 删除保存配置
	Delete(ctx context.Context, id int64) error
}
