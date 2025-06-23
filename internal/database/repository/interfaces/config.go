package interfaces

import (
	"context"

	"github.com/bestruirui/bestsub/internal/database/models"
)

// SystemConfigRepository 系统配置数据访问接口
type SystemConfigRepository interface {
	// Create 创建配置
	Create(ctx context.Context, config *models.SystemConfig) error

	// GetByID 根据ID获取配置
	GetByID(ctx context.Context, id int64) (*models.SystemConfig, error)

	// GetByKey 根据键获取配置
	GetByKey(ctx context.Context, key string) (*models.SystemConfig, error)

	// Update 更新配置
	Update(ctx context.Context, config *models.SystemConfig) error

	// Delete 删除配置
	Delete(ctx context.Context, id int64) error

	// DeleteByKey 根据键删除配置
	DeleteByKey(ctx context.Context, key string) error

	// List 获取配置列表
	List(ctx context.Context, offset, limit int) ([]*models.SystemConfig, error)

	// ListByGroup 根据分组获取配置列表
	ListByGroup(ctx context.Context, group string) ([]*models.SystemConfig, error)

	// Count 获取配置总数
	Count(ctx context.Context) (int64, error)

	// SetValue 设置配置值
	SetValue(ctx context.Context, key, value, configType, group, description string) error

	// GetValue 获取配置值
	GetValue(ctx context.Context, key string) (string, error)
}

// NotificationChannelRepository 通知渠道数据访问接口
type NotificationChannelRepository interface {
	// Create 创建通知渠道
	Create(ctx context.Context, channel *models.NotificationChannel) error

	// GetByID 根据ID获取通知渠道
	GetByID(ctx context.Context, id int64) (*models.NotificationChannel, error)

	// Update 更新通知渠道
	Update(ctx context.Context, channel *models.NotificationChannel) error

	// Delete 删除通知渠道
	Delete(ctx context.Context, id int64) error

	// List 获取通知渠道列表
	List(ctx context.Context, offset, limit int) ([]*models.NotificationChannel, error)

	// ListActive 获取活跃的通知渠道列表
	ListActive(ctx context.Context) ([]*models.NotificationChannel, error)

	// ListByType 根据类型获取通知渠道列表
	ListByType(ctx context.Context, channelType string) ([]*models.NotificationChannel, error)

	// Count 获取通知渠道总数
	Count(ctx context.Context) (int64, error)

	// UpdateTestResult 更新测试结果
	UpdateTestResult(ctx context.Context, id int64, testResult string) error
}
