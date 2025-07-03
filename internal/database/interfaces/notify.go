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

	// 根据任务ID获取通知渠道
	GetByTaskID(ctx context.Context, taskID int64) (*models.Notify, error)
}
