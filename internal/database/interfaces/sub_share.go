package interfaces

import (
	"context"

	"github.com/bestruirui/bestsub/internal/models/sub"
)

// 分享链接数据访问接口
type SubShareRepository interface {
	// Create 创建分享链接
	Create(ctx context.Context, shareLink *sub.Share) error

	// GetByID 根据ID获取分享链接
	GetByID(ctx context.Context, id uint16) (*sub.Share, error)

	// Update 更新分享链接
	Update(ctx context.Context, shareLink *sub.Share) error

	// Delete 删除分享链接
	Delete(ctx context.Context, id uint16) error

	// List 获取分享链接列表
	List(ctx context.Context) ([]*sub.Share, error)
}
