package interfaces

import (
	"context"

	"github.com/bestruirui/bestsub/internal/models/config"
)

// ConfigRepository 系统配置数据访问接口
type ConfigRepository interface {
	Create(ctx context.Context, config *[]config.Advance) error

	GetAll(ctx context.Context) (*[]config.Advance, error)

	GetByKey(ctx context.Context, key []string) (*[]config.Advance, error)

	Update(ctx context.Context, data *[]config.UpdateAdvance) error
}
