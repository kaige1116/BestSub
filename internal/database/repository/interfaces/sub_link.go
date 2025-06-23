package interfaces

import (
	"context"

	"github.com/bestruirui/bestsub/internal/database/models"
)

// 链接数据访问接口
type SubLinkRepository interface {
	// Create 创建链接
	Create(ctx context.Context, link *models.SubLink) error

	// GetByID 根据ID获取链接
	GetByID(ctx context.Context, id int64) (*models.SubLink, error)

	// Update 更新链接
	Update(ctx context.Context, link *models.SubLink) error

	// Delete 删除链接
	Delete(ctx context.Context, id int64) error

	// List 获取链接列表
	List(ctx context.Context, offset, limit int) ([]*models.SubLink, error)

	// ListEnabled 获取启用的链接列表
	ListEnabled(ctx context.Context) ([]*models.SubLink, error)

	// Count 获取链接总数
	Count(ctx context.Context) (int64, error)

	// UpdateStatus 更新链接状态
	UpdateStatus(ctx context.Context, id int64, status, errorMsg string) error

	// GetByURL 根据URL获取链接
	GetByURL(ctx context.Context, url string) (*models.SubLink, error)
}

// 模块配置数据访问接口
type SubLinkModuleConfigRepository interface {
	// Create 创建模块配置
	Create(ctx context.Context, config *models.SubLinkModuleConfig) error

	// GetByID 根据ID获取模块配置
	GetByID(ctx context.Context, id int64) (*models.SubLinkModuleConfig, error)

	// Update 更新模块配置
	Update(ctx context.Context, config *models.SubLinkModuleConfig) error

	// Delete 删除模块配置
	Delete(ctx context.Context, id int64) error

	// DeleteByLinkID 删除链接的所有模块配置
	DeleteByLinkID(ctx context.Context, linkID int64) error

	// List 获取模块配置列表
	List(ctx context.Context, offset, limit int) ([]*models.SubLinkModuleConfig, error)

	// ListByLinkID 根据链接ID获取模块配置列表
	ListByLinkID(ctx context.Context, linkID int64) ([]*models.SubLinkModuleConfig, error)

	// ListByLinkIDAndType 根据链接ID和模块类型获取配置列表
	ListByLinkIDAndType(ctx context.Context, linkID int64, moduleType string) ([]*models.SubLinkModuleConfig, error)

	// ListEnabledByLinkID 获取链接的启用模块配置列表
	ListEnabledByLinkID(ctx context.Context, linkID int64) ([]*models.SubLinkModuleConfig, error)

	// Count 获取模块配置总数
	Count(ctx context.Context) (int64, error)

	// CountByLinkID 获取链接的模块配置总数
	CountByLinkID(ctx context.Context, linkID int64) (int64, error)
}
