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
	GetByID(ctx context.Context, id int64) (*sub.Share, error)

	// Update 更新分享链接
	Update(ctx context.Context, shareLink *sub.Share) error

	// Delete 删除分享链接
	Delete(ctx context.Context, id int64) error

	// List 获取分享链接列表
	List(ctx context.Context, offset, limit int) ([]*sub.Share, error)

	// Count 获取分享链接总数
	Count(ctx context.Context) (int64, error)

	// 添加分享链接与输出模板的关联
	AddOutputTemplateRelation(ctx context.Context, shareID, templateID int64) error

	// 添加分享链接与过滤配置的关联
	AddFilterConfigRelation(ctx context.Context, shareID, configID int64) error

	// 添加分享链接与订阅的关联
	AddSubRelation(ctx context.Context, shareID, subID int64) error
}
