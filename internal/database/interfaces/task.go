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

	// 根据订阅ID获取任务列表（一个订阅可以有多个任务）
	GetBySubID(ctx context.Context, subID uint16) (*[]task.Data, error)

	// 根据订阅ID获取任务ID列表
	GetTaskIDsBySubID(ctx context.Context, subID uint16) ([]uint16, error)

	// List 获取任务列表
	List(ctx context.Context, offset, limit int) (*[]task.Data, error)

	// Count 获取任务总数
	Count(ctx context.Context) (uint16, error)

	// GetSystemTasks 获取所有系统任务
	GetSystemTasks(ctx context.Context) (*[]task.Data, error)

	// 添加任务与通知的关联
	AddNotifyRelation(ctx context.Context, taskID, notifyID uint16) error

	// 添加任务与通知模板的关联
	AddNotifyTemplateRelation(ctx context.Context, taskID, notifyTemplateID uint16) error
}
