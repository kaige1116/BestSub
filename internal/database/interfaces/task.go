package interfaces

import (
	"context"

	"github.com/bestruirui/bestsub/internal/models/task"
)

// TaskRepository 任务数据访问接口
type TaskRepository interface {
	// Create 创建任务
	Create(ctx context.Context, task *task.Data) (id uint16, err error)

	// Update 更新任务
	Update(ctx context.Context, task *task.Data) error

	// Delete 删除任务
	Delete(ctx context.Context, id uint16) error

	// 根据任务ID获取任务
	GetByID(ctx context.Context, id uint16) (*task.Data, error)

	// List 获取任务列表
	List(ctx context.Context) (*[]task.Data, error)
}
