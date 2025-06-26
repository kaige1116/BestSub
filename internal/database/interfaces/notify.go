package interfaces

import (
	"context"

	"github.com/bestruirui/bestsub/internal/database/models"
)

// NotificationChannelRepository 通知渠道数据访问接口
type NotifyRepository interface {
	// Create 创建通知渠道
	Create(ctx context.Context, channel *models.Notify) error

	// GetByID 根据ID获取通知渠道
	GetByID(ctx context.Context, id int64) (*models.Notify, error)

	// Update 更新通知渠道
	Update(ctx context.Context, channel *models.Notify) error

	// Delete 删除通知渠道
	Delete(ctx context.Context, id int64) error

	// List 获取通知渠道列表
	List(ctx context.Context, offset, limit int) ([]*models.Notify, error)

	// ListActive 获取活跃的通知渠道列表
	ListActive(ctx context.Context) ([]*models.Notify, error)

	// ListByType 根据类型获取通知渠道列表
	ListByType(ctx context.Context, channelType string) ([]*models.Notify, error)

	// Count 获取通知渠道总数
	Count(ctx context.Context) (int64, error)

	// UpdateTestResult 更新测试结果
	UpdateTestResult(ctx context.Context, id int64, testResult string) error
}
