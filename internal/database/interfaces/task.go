package interfaces

import (
	"context"
	"time"
	"github.com/bestruirui/bestsub/internal/database/models"
)

// TaskRepository 任务数据访问接口
type TaskRepository interface {
	// Create 创建任务
	Create(ctx context.Context, task *models.Task) error
	
	// GetByID 根据ID获取任务
	GetByID(ctx context.Context, id int64) (*models.Task, error)
	
	// Update 更新任务
	Update(ctx context.Context, task *models.Task) error
	
	// Delete 删除任务
	Delete(ctx context.Context, id int64) error
	
	// List 获取任务列表
	List(ctx context.Context, offset, limit int) ([]*models.Task, error)
	
	// ListByStatus 根据状态获取任务列表
	ListByStatus(ctx context.Context, status string) ([]*models.Task, error)
	
	// ListByType 根据类型获取任务列表
	ListByType(ctx context.Context, taskType string) ([]*models.Task, error)
	
	// ListByLinkID 根据链接ID获取任务列表
	ListByLinkID(ctx context.Context, linkID int64) ([]*models.Task, error)
	
	// ListPending 获取待执行的任务列表
	ListPending(ctx context.Context) ([]*models.Task, error)
	
	// ListScheduled 获取定时任务列表
	ListScheduled(ctx context.Context, before time.Time) ([]*models.Task, error)
	
	// Count 获取任务总数
	Count(ctx context.Context) (int64, error)
	
	// CountByStatus 根据状态获取任务总数
	CountByStatus(ctx context.Context, status string) (int64, error)
	
	// UpdateStatus 更新任务状态
	UpdateStatus(ctx context.Context, id int64, status string) error
	
	// UpdateResult 更新任务结果
	UpdateResult(ctx context.Context, id int64, result, errorMsg string) error
	
	// UpdateTiming 更新任务时间信息
	UpdateTiming(ctx context.Context, id int64, startTime, endTime time.Time, duration int) error
	
	// IncrementRetryCount 增加重试次数
	IncrementRetryCount(ctx context.Context, id int64) error
	
	// UpdateNextRun 更新下次执行时间
	UpdateNextRun(ctx context.Context, id int64, nextRun time.Time) error
	
	// DeleteCompleted 删除已完成的任务
	DeleteCompleted(ctx context.Context, before time.Time) error
	
	// DeleteFailed 删除失败的任务
	DeleteFailed(ctx context.Context, before time.Time) error
}
