package interfaces

import (
	"context"

	"github.com/bestruirui/bestsub/internal/models/system"
)

// ConfigRepository 系统配置数据访问接口
type ConfigRepository interface {
	Create(ctx context.Context, config *[]system.Data) error

	GetAll(ctx context.Context) (*[]system.Data, error)

	GetByKey(ctx context.Context, key []string) (*[]system.Data, error)

	Update(ctx context.Context, data *[]system.UpdateData) error
}
