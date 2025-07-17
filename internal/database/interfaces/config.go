package interfaces

import (
	"context"

	"github.com/bestruirui/bestsub/internal/models/system"
)

// ConfigRepository 系统配置数据访问接口
type ConfigRepository interface {
	// Create 创建配置
	Create(ctx context.Context, config *[]system.Data) error

	// GetByKey 根据键获取配置
	GetByKey(ctx context.Context, key string) (*system.Data, error)

	// Update 更新配置
	Update(ctx context.Context, config *[]system.Data) error

	// GetConfigsByGroup 获取指定分组下的所有配置
	GetConfigsByGroup(ctx context.Context, group string) ([]system.Data, error)
}
