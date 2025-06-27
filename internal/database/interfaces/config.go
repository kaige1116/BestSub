package interfaces

import (
	"context"

	"github.com/bestruirui/bestsub/internal/database/models"
)

// SystemConfigRepository 系统配置数据访问接口
type SystemConfigRepository interface {
	// Create 创建配置
	Create(ctx context.Context, config *models.SystemConfig) error

	// GetByKey 根据键获取配置
	GetByKey(ctx context.Context, key string) (*models.SystemConfig, error)

	// Update 更新配置
	Update(ctx context.Context, config *models.SystemConfig) error

	// DeleteByKey 根据键删除配置
	DeleteByKey(ctx context.Context, key string) error

	// SetValue 设置配置值（如果不存在则创建，存在则更新）
	SetValue(ctx context.Context, key, value, configType, group, description string) error

	// GetAllKeys 获取所有配置键
	GetAllKeys(ctx context.Context) ([]string, error)

	// GetAllGroups 获取所有配置分组
	GetAllGroups(ctx context.Context) ([]string, error)

	// GetConfigsByGroup 获取指定分组下的所有配置
	GetConfigsByGroup(ctx context.Context, group string) ([]models.SystemConfig, error)
}
