package interfaces

import (
	"context"

	"github.com/bestruirui/bestsub/internal/models/sub"
)

// SubRepository 订阅链接数据访问接口
type SubRepository interface {
	// Create 创建链接
	Create(ctx context.Context, link *sub.Data) error

	// GetByID 根据ID获取链接
	GetByID(ctx context.Context, id uint16) (*sub.Data, error)

	// Update 更新链接
	Update(ctx context.Context, link *sub.Data) error

	// Delete 删除链接
	Delete(ctx context.Context, id uint16) error

	// List 获取订阅链接列表
	List(ctx context.Context, offset, limit int) (*[]sub.Data, error)

	// Count 获取订阅链接总数
	Count(ctx context.Context) (uint16, error)

	// 根据任务ID获取订阅ID
	GetByTaskID(ctx context.Context, taskID uint16) (uint16, error)

	// 根据分享ID获取订阅ID
	GetByShareID(ctx context.Context, shareID uint16) ([]uint16, error)

	// 根据保存ID获取订阅ID
	GetBySaveID(ctx context.Context, saveID uint16) ([]uint16, error)

	// 添加任务与订阅的关联
	AddTaskRelation(ctx context.Context, subID, taskID uint16) error

	// 添加保存配置与订阅的关联
	AddSaveRelation(ctx context.Context, subID, saveID uint16) error

	// 添加分享链接与订阅的关联
	AddShareRelation(ctx context.Context, subID, shareID uint16) error
}
