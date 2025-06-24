package interfaces

import (
	"context"

	"github.com/bestruirui/bestsub/internal/database/models"
)

// 分享链接数据访问接口
type SubShareLinkRepository interface {
	// Create 创建分享链接
	Create(ctx context.Context, shareLink *models.SubShareLink) error

	// GetByID 根据ID获取分享链接
	GetByID(ctx context.Context, id int64) (*models.SubShareLink, error)

	// GetByToken 根据Token获取分享链接
	GetByToken(ctx context.Context, token string) (*models.SubShareLink, error)

	// Update 更新分享链接
	Update(ctx context.Context, shareLink *models.SubShareLink) error

	// Delete 删除分享链接
	Delete(ctx context.Context, id int64) error

	// List 获取分享链接列表
	List(ctx context.Context, offset, limit int) ([]*models.SubShareLink, error)

	// ListEnabled 获取启用的分享链接列表
	ListEnabled(ctx context.Context) ([]*models.SubShareLink, error)

	// Count 获取分享链接总数
	Count(ctx context.Context) (int64, error)

	// IncrementDownloadCount 增加下载次数
	IncrementDownloadCount(ctx context.Context, id int64) error

	// UpdateLastAccess 更新最后访问时间
	UpdateLastAccess(ctx context.Context, id int64) error

	// DeleteExpired 删除过期的分享链接
	DeleteExpired(ctx context.Context) error

	// IsTokenUnique 检查Token是否唯一
	IsTokenUnique(ctx context.Context, token string) (bool, error)
}
