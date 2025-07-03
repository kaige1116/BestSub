package interfaces

import (
	"context"

	"github.com/bestruirui/bestsub/internal/database/models"
)

// TaskRepository 任务数据访问接口
type TaskRepository interface {
	// Create 创建任务
	Create(ctx context.Context, task *models.Task) error

	// Update 更新任务
	Update(ctx context.Context, id int64, task *models.Task) error

	// Delete 删除任务
	Delete(ctx context.Context, id int64) error

	// 根据任务ID获取任务
	GetByID(ctx context.Context, id int64) (*models.Task, error)

	// 根据订阅ID获取任务
	GetBySubID(ctx context.Context, subID int64) (*models.Task, error)

	// 根据保存ID获取任务
	GetBySaveID(ctx context.Context, saveID int64) (*models.Task, error)
}
