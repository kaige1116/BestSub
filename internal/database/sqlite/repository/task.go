package repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/bestruirui/bestsub/internal/database/interfaces"
	"github.com/bestruirui/bestsub/internal/database/models"
	"github.com/bestruirui/bestsub/internal/database/sqlite/database"
	timeutils "github.com/bestruirui/bestsub/internal/utils/time"
)

// TaskRepository 任务数据访问实现
type TaskRepository struct {
	db *database.Database
}

// newTaskRepository 创建任务仓库
func newTaskRepository(db *database.Database) interfaces.TaskRepository {
	return &TaskRepository{db: db}
}

// Create 创建任务
func (r *TaskRepository) Create(ctx context.Context, task *models.Task) error {
	query := `INSERT INTO tasks (type, name, description, status, priority, link_id, config, result, 
	          error_msg, start_time, end_time, duration, retry_count, max_retries, next_run, cron_expr, created_at, updated_at) 
	          VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`

	now := timeutils.Now()
	result, err := r.db.ExecContext(ctx, query,
		task.Type,
		task.Name,
		task.Description,
		task.Status,
		task.Priority,
		task.LinkID,
		task.Config,
		task.Result,
		task.ErrorMsg,
		task.StartTime,
		task.EndTime,
		task.Duration,
		task.RetryCount,
		task.MaxRetries,
		task.NextRun,
		task.CronExpr,
		now,
		now,
	)

	if err != nil {
		return fmt.Errorf("failed to create task: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("failed to get task id: %w", err)
	}

	task.ID = id
	task.CreatedAt = now
	task.UpdatedAt = now

	return nil
}

// GetByID 根据ID获取任务
func (r *TaskRepository) GetByID(ctx context.Context, id int64) (*models.Task, error) {
	query := `SELECT id, type, name, description, status, priority, link_id, config, result, 
	          error_msg, start_time, end_time, duration, retry_count, max_retries, next_run, cron_expr, created_at, updated_at 
	          FROM tasks WHERE id = ?`

	var task models.Task
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&task.ID,
		&task.Type,
		&task.Name,
		&task.Description,
		&task.Status,
		&task.Priority,
		&task.LinkID,
		&task.Config,
		&task.Result,
		&task.ErrorMsg,
		&task.StartTime,
		&task.EndTime,
		&task.Duration,
		&task.RetryCount,
		&task.MaxRetries,
		&task.NextRun,
		&task.CronExpr,
		&task.CreatedAt,
		&task.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get task by id: %w", err)
	}

	return &task, nil
}

// Update 更新任务
func (r *TaskRepository) Update(ctx context.Context, task *models.Task) error {
	query := `UPDATE tasks SET type = ?, name = ?, description = ?, status = ?, priority = ?, 
	          link_id = ?, config = ?, result = ?, error_msg = ?, start_time = ?, end_time = ?, 
	          duration = ?, retry_count = ?, max_retries = ?, next_run = ?, cron_expr = ?, updated_at = ? WHERE id = ?`

	_, err := r.db.ExecContext(ctx, query,
		task.Type,
		task.Name,
		task.Description,
		task.Status,
		task.Priority,
		task.LinkID,
		task.Config,
		task.Result,
		task.ErrorMsg,
		task.StartTime,
		task.EndTime,
		task.Duration,
		task.RetryCount,
		task.MaxRetries,
		task.NextRun,
		task.CronExpr,
		timeutils.Now(),
		task.ID,
	)

	if err != nil {
		return fmt.Errorf("failed to update task: %w", err)
	}

	return nil
}

// Delete 删除任务
func (r *TaskRepository) Delete(ctx context.Context, id int64) error {
	query := `DELETE FROM tasks WHERE id = ?`

	_, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete task: %w", err)
	}

	return nil
}

// List 获取任务列表
func (r *TaskRepository) List(ctx context.Context, offset, limit int) ([]*models.Task, error) {
	query := `SELECT id, type, name, description, status, priority, link_id, config, result, 
	          error_msg, start_time, end_time, duration, retry_count, max_retries, next_run, cron_expr, created_at, updated_at 
	          FROM tasks ORDER BY created_at DESC LIMIT ? OFFSET ?`

	return r.queryTasks(ctx, query, limit, offset)
}

// ListByStatus 根据状态获取任务列表
func (r *TaskRepository) ListByStatus(ctx context.Context, status string) ([]*models.Task, error) {
	query := `SELECT id, type, name, description, status, priority, link_id, config, result, 
	          error_msg, start_time, end_time, duration, retry_count, max_retries, next_run, cron_expr, created_at, updated_at 
	          FROM tasks WHERE status = ? ORDER BY created_at DESC`

	return r.queryTasks(ctx, query, status)
}

// ListByType 根据类型获取任务列表
func (r *TaskRepository) ListByType(ctx context.Context, taskType string) ([]*models.Task, error) {
	query := `SELECT id, type, name, description, status, priority, link_id, config, result, 
	          error_msg, start_time, end_time, duration, retry_count, max_retries, next_run, cron_expr, created_at, updated_at 
	          FROM tasks WHERE type = ? ORDER BY created_at DESC`

	return r.queryTasks(ctx, query, taskType)
}

// ListByLinkID 根据链接ID获取任务列表
func (r *TaskRepository) ListByLinkID(ctx context.Context, linkID int64) ([]*models.Task, error) {
	query := `SELECT id, type, name, description, status, priority, link_id, config, result, 
	          error_msg, start_time, end_time, duration, retry_count, max_retries, next_run, cron_expr, created_at, updated_at 
	          FROM tasks WHERE link_id = ? ORDER BY created_at DESC`

	return r.queryTasks(ctx, query, linkID)
}

// ListPending 获取待执行的任务列表
func (r *TaskRepository) ListPending(ctx context.Context) ([]*models.Task, error) {
	query := `SELECT id, type, name, description, status, priority, link_id, config, result, 
	          error_msg, start_time, end_time, duration, retry_count, max_retries, next_run, cron_expr, created_at, updated_at 
	          FROM tasks WHERE status = ? ORDER BY priority ASC, created_at ASC`

	return r.queryTasks(ctx, query, models.TaskStatusPending)
}

// ListScheduled 获取定时任务列表
func (r *TaskRepository) ListScheduled(ctx context.Context, before time.Time) ([]*models.Task, error) {
	query := `SELECT id, type, name, description, status, priority, link_id, config, result, 
	          error_msg, start_time, end_time, duration, retry_count, max_retries, next_run, cron_expr, created_at, updated_at 
	          FROM tasks WHERE next_run <= ? AND status = ? ORDER BY next_run ASC`

	return r.queryTasks(ctx, query, before, models.TaskStatusPending)
}

// Count 获取任务总数
func (r *TaskRepository) Count(ctx context.Context) (int64, error) {
	query := `SELECT COUNT(*) FROM tasks`

	var count int64
	err := r.db.QueryRowContext(ctx, query).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to count tasks: %w", err)
	}

	return count, nil
}

// CountByStatus 根据状态获取任务总数
func (r *TaskRepository) CountByStatus(ctx context.Context, status string) (int64, error) {
	query := `SELECT COUNT(*) FROM tasks WHERE status = ?`

	var count int64
	err := r.db.QueryRowContext(ctx, query, status).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to count tasks by status: %w", err)
	}

	return count, nil
}

// UpdateStatus 更新任务状态
func (r *TaskRepository) UpdateStatus(ctx context.Context, id int64, status string) error {
	query := `UPDATE tasks SET status = ?, updated_at = ? WHERE id = ?`

	_, err := r.db.ExecContext(ctx, query, status, timeutils.Now(), id)
	if err != nil {
		return fmt.Errorf("failed to update task status: %w", err)
	}

	return nil
}

// UpdateResult 更新任务结果
func (r *TaskRepository) UpdateResult(ctx context.Context, id int64, result, errorMsg string) error {
	query := `UPDATE tasks SET result = ?, error_msg = ?, updated_at = ? WHERE id = ?`

	_, err := r.db.ExecContext(ctx, query, result, errorMsg, timeutils.Now(), id)
	if err != nil {
		return fmt.Errorf("failed to update task result: %w", err)
	}

	return nil
}

// UpdateTiming 更新任务时间信息
func (r *TaskRepository) UpdateTiming(ctx context.Context, id int64, startTime, endTime time.Time, duration int) error {
	query := `UPDATE tasks SET start_time = ?, end_time = ?, duration = ?, updated_at = ? WHERE id = ?`

	_, err := r.db.ExecContext(ctx, query, startTime, endTime, duration, timeutils.Now(), id)
	if err != nil {
		return fmt.Errorf("failed to update task timing: %w", err)
	}

	return nil
}

// IncrementRetryCount 增加重试次数
func (r *TaskRepository) IncrementRetryCount(ctx context.Context, id int64) error {
	query := `UPDATE tasks SET retry_count = retry_count + 1, updated_at = ? WHERE id = ?`

	_, err := r.db.ExecContext(ctx, query, timeutils.Now(), id)
	if err != nil {
		return fmt.Errorf("failed to increment retry count: %w", err)
	}

	return nil
}

// UpdateNextRun 更新下次执行时间
func (r *TaskRepository) UpdateNextRun(ctx context.Context, id int64, nextRun time.Time) error {
	query := `UPDATE tasks SET next_run = ?, updated_at = ? WHERE id = ?`

	_, err := r.db.ExecContext(ctx, query, nextRun, timeutils.Now(), id)
	if err != nil {
		return fmt.Errorf("failed to update next run: %w", err)
	}

	return nil
}

// DeleteCompleted 删除已完成的任务
func (r *TaskRepository) DeleteCompleted(ctx context.Context, before time.Time) error {
	query := `DELETE FROM tasks WHERE status = ? AND end_time < ?`

	_, err := r.db.ExecContext(ctx, query, models.TaskStatusCompleted, before)
	if err != nil {
		return fmt.Errorf("failed to delete completed tasks: %w", err)
	}

	return nil
}

// DeleteFailed 删除失败的任务
func (r *TaskRepository) DeleteFailed(ctx context.Context, before time.Time) error {
	query := `DELETE FROM tasks WHERE status = ? AND end_time < ?`

	_, err := r.db.ExecContext(ctx, query, models.TaskStatusFailed, before)
	if err != nil {
		return fmt.Errorf("failed to delete failed tasks: %w", err)
	}

	return nil
}

// queryTasks 通用任务查询方法
func (r *TaskRepository) queryTasks(ctx context.Context, query string, args ...interface{}) ([]*models.Task, error) {
	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to query tasks: %w", err)
	}
	defer rows.Close()

	var tasks []*models.Task
	for rows.Next() {
		var task models.Task
		err := rows.Scan(
			&task.ID,
			&task.Type,
			&task.Name,
			&task.Description,
			&task.Status,
			&task.Priority,
			&task.LinkID,
			&task.Config,
			&task.Result,
			&task.ErrorMsg,
			&task.StartTime,
			&task.EndTime,
			&task.Duration,
			&task.RetryCount,
			&task.MaxRetries,
			&task.NextRun,
			&task.CronExpr,
			&task.CreatedAt,
			&task.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan task: %w", err)
		}
		tasks = append(tasks, &task)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to iterate tasks: %w", err)
	}

	return tasks, nil
}
