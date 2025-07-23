package interfaces

import (
	"context"

	"github.com/bestruirui/bestsub/internal/models/sub"
)

// 输出模板数据访问接口
type SubTemplateRepository interface {
	// Create 创建输出模板
	Create(ctx context.Context, template *sub.Template) error

	// GetByID 根据ID获取输出模板
	GetByID(ctx context.Context, id uint16) (*sub.Template, error)

	// Update 更新输出模板
	Update(ctx context.Context, template *sub.Template) error

	// Delete 删除输出模板
	Delete(ctx context.Context, id uint16) error

	// List 获取输出模板列表
	List(ctx context.Context) (*[]sub.Template, error)
}
