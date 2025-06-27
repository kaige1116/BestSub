package interfaces

import (
	"context"

	"github.com/bestruirui/bestsub/internal/database/models"
)

// SubLinkRepository 订阅链接数据访问接口
type SubLinkRepository interface {
	// Create 创建链接
	Create(ctx context.Context, link *models.SubLink) error

	// GetByID 根据ID获取链接
	GetByID(ctx context.Context, id int64) (*models.SubLink, error)

	// GetByURL 根据URL获取链接
	GetByURL(ctx context.Context, url string) (*models.SubLink, error)

	// Update 更新链接
	Update(ctx context.Context, link *models.SubLink) error

	// Delete 删除链接
	Delete(ctx context.Context, id int64) error

	// List 获取链接列表
	List(ctx context.Context, offset, limit int) ([]*models.SubLink, error)

	// Count 获取链接总数
	Count(ctx context.Context) (int64, error)
}
