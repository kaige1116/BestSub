package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/bestruirui/bestsub/internal/database/interfaces"
	"github.com/bestruirui/bestsub/internal/database/sqlite/database"
	"github.com/bestruirui/bestsub/internal/models/task"
	"github.com/bestruirui/bestsub/internal/utils/local"
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
func (r *TaskRepository) Create(ctx context.Context, t *task.Data) (int64, error) {
	query := `INSERT INTO tasks (enable, name, description, is_sys_task, cron, type, log_level, timeout, retry, config, last_run_result, last_run_time, last_run_duration, success_count, failed_count, created_at, updated_at)
	          VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`

	now := local.Time()
	result, err := r.db.ExecContext(ctx, query,
		t.Enable,
		t.Name,
		t.Description,
		t.IsSysTask,
		t.Cron,
		t.Type,
		t.LogLevel,
		t.Timeout,
		t.Retry,
		t.Config,
		t.LastRunResult,
		t.LastRunTime,
		t.LastRunDuration,
		t.SuccessCount,
		t.FailedCount,
		now,
		now,
	)

	if err != nil {
		return 0, fmt.Errorf("failed to create task: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("failed to get task id: %w", err)
	}

	return id, nil
}

// GetByID 根据ID获取任务
func (r *TaskRepository) GetByID(ctx context.Context, id int64) (*task.Data, error) {
	query := `SELECT id, enable, name, description, is_sys_task, cron, type, log_level, timeout, retry, config, last_run_result, last_run_time, last_run_duration, success_count, failed_count, created_at, updated_at
	          FROM tasks WHERE id = ?`

	var t task.Data
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&t.ID,
		&t.Enable,
		&t.Name,
		&t.Description,
		&t.IsSysTask,
		&t.Cron,
		&t.Type,
		&t.LogLevel,
		&t.Timeout,
		&t.Retry,
		&t.Config,
		&t.LastRunResult,
		&t.LastRunTime,
		&t.LastRunDuration,
		&t.SuccessCount,
		&t.FailedCount,
		&t.CreatedAt,
		&t.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get task by id: %w", err)
	}

	return &t, nil
}

// Update 更新任务
func (r *TaskRepository) Update(ctx context.Context, t *task.Data) error {
	query := `UPDATE tasks SET enable = ?, name = ?, description = ?, is_sys_task = ?, cron = ?, type = ?, log_level = ?, timeout = ?, retry = ?, config = ?,
	          last_run_result = ?, last_run_time = ?, last_run_duration = ?, success_count = ?, failed_count = ?, updated_at = ? WHERE id = ?`

	_, err := r.db.ExecContext(ctx, query,
		t.Enable,
		t.Name,
		t.Description,
		t.IsSysTask,
		t.Cron,
		t.Type,
		t.LogLevel,
		t.Timeout,
		t.Retry,
		t.Config,
		t.LastRunResult,
		t.LastRunTime,
		t.LastRunDuration,
		t.SuccessCount,
		t.FailedCount,
		local.Time(),
		t.ID,
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
func (r *TaskRepository) List(ctx context.Context, offset, limit int) (*[]task.Data, error) {
	query := `SELECT id, enable, name, description, is_sys_task, cron, type, log_level, timeout, retry, config, last_run_result, last_run_time, last_run_duration, success_count, failed_count, created_at, updated_at
	          FROM tasks ORDER BY created_at DESC LIMIT ? OFFSET ?`

	rows, err := r.db.QueryContext(ctx, query, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to list tasks: %w", err)
	}
	defer rows.Close()

	var tasks []task.Data
	for rows.Next() {
		var t task.Data
		err := rows.Scan(
			&t.ID,
			&t.Enable,
			&t.Name,
			&t.Description,
			&t.IsSysTask,
			&t.Cron,
			&t.Type,
			&t.LogLevel,
			&t.Timeout,
			&t.Retry,
			&t.Config,
			&t.LastRunResult,
			&t.LastRunTime,
			&t.LastRunDuration,
			&t.SuccessCount,
			&t.FailedCount,
			&t.CreatedAt,
			&t.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan task: %w", err)
		}
		tasks = append(tasks, t)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to iterate tasks: %w", err)
	}

	return &tasks, nil
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

// GetSystemTasks 获取所有系统任务
func (r *TaskRepository) GetSystemTasks(ctx context.Context) (*[]task.Data, error) {
	query := `SELECT id, enable, name, description, is_sys_task, cron, type, log_level, timeout, retry, config, last_run_result, last_run_time, last_run_duration, success_count, failed_count, created_at, updated_at
	          FROM tasks WHERE is_sys_task = true ORDER BY created_at DESC`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to get system tasks: %w", err)
	}
	defer rows.Close()

	var tasks []task.Data
	for rows.Next() {
		var t task.Data
		if err := rows.Scan(
			&t.ID,
			&t.Enable,
			&t.Name,
			&t.Description,
			&t.IsSysTask,
			&t.Cron,
			&t.Type,
			&t.LogLevel,
			&t.Timeout,
			&t.Retry,
			&t.Config,
			&t.LastRunResult,
			&t.LastRunTime,
			&t.LastRunDuration,
			&t.SuccessCount,
			&t.FailedCount,
			&t.CreatedAt,
			&t.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("failed to scan system task: %w", err)
		}
		tasks = append(tasks, t)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to iterate system tasks: %w", err)
	}

	return &tasks, nil
}

// GetBySubID 根据订阅ID获取任务列表
func (r *TaskRepository) GetBySubID(ctx context.Context, subID int64) (*[]task.Data, error) {
	query := `SELECT t.id, t.enable, t.name, t.description, t.is_sys_task, t.cron, t.type, t.log_level, t.timeout, t.retry, t.config, t.last_run_result, t.last_run_time, t.last_run_duration, t.success_count, t.failed_count, t.created_at, t.updated_at
	          FROM tasks t
	          INNER JOIN sub_task_relations str ON t.id = str.task_id
	          WHERE str.sub_id = ? ORDER BY t.created_at DESC`

	rows, err := r.db.QueryContext(ctx, query, subID)
	if err != nil {
		return nil, fmt.Errorf("failed to get tasks by sub id: %w", err)
	}
	defer rows.Close()

	var tasks []task.Data
	for rows.Next() {
		var t task.Data
		err := rows.Scan(
			&t.ID,
			&t.Enable,
			&t.Name,
			&t.Description,
			&t.IsSysTask,
			&t.Cron,
			&t.Type,
			&t.LogLevel,
			&t.Timeout,
			&t.Retry,
			&t.Config,
			&t.LastRunResult,
			&t.LastRunTime,
			&t.LastRunDuration,
			&t.SuccessCount,
			&t.FailedCount,
			&t.CreatedAt,
			&t.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan task: %w", err)
		}
		tasks = append(tasks, t)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to iterate tasks: %w", err)
	}

	return &tasks, nil
}

// GetTaskIDsBySubID 根据订阅ID获取任务ID列表
func (r *TaskRepository) GetTaskIDsBySubID(ctx context.Context, subID int64) ([]int64, error) {
	query := `SELECT task_id FROM sub_task_relations WHERE sub_id = ? ORDER BY task_id`

	rows, err := r.db.QueryContext(ctx, query, subID)
	if err != nil {
		return nil, fmt.Errorf("failed to get task ids by sub id: %w", err)
	}
	defer rows.Close()

	var taskIDs []int64
	for rows.Next() {
		var taskID int64
		err := rows.Scan(&taskID)
		if err != nil {
			return nil, fmt.Errorf("failed to scan task id: %w", err)
		}
		taskIDs = append(taskIDs, taskID)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to iterate task ids: %w", err)
	}

	return taskIDs, nil
}

// AddNotifyRelation 添加任务与通知的关联
func (r *TaskRepository) AddNotifyRelation(ctx context.Context, taskID, notifyID int64) error {
	query := `INSERT OR IGNORE INTO task_notify_relations (task_id, notify_id) VALUES (?, ?)`

	_, err := r.db.ExecContext(ctx, query, taskID, notifyID)
	if err != nil {
		return fmt.Errorf("failed to add notify relation: %w", err)
	}

	return nil
}

// AddNotifyTemplateRelation 添加任务与通知模板的关联
func (r *TaskRepository) AddNotifyTemplateRelation(ctx context.Context, taskID, notifyTemplateID int64) error {
	query := `INSERT OR IGNORE INTO task_notify_template_relations (task_id, notify_template_id) VALUES (?, ?)`

	_, err := r.db.ExecContext(ctx, query, taskID, notifyTemplateID)
	if err != nil {
		return fmt.Errorf("failed to add notify template relation: %w", err)
	}

	return nil
}

// DeleteBySubID 根据订阅ID删除所有的任务
func (r *TaskRepository) DeleteBySubID(ctx context.Context, subID int64) error {
	query := `DELETE FROM tasks WHERE id IN (SELECT task_id FROM sub_task_relations WHERE sub_id = ?)`

	_, err := r.db.ExecContext(ctx, query, subID)
	if err != nil {
		return fmt.Errorf("failed to delete tasks by sub id: %w", err)
	}

	return nil
}

// DeleteBySaveID 根据保存ID删除所有的任务
func (r *TaskRepository) DeleteBySaveID(ctx context.Context, saveID int64) error {
	query := `DELETE FROM tasks WHERE id IN (SELECT task_id FROM save_task_relations WHERE save_id = ?)`

	_, err := r.db.ExecContext(ctx, query, saveID)
	if err != nil {
		return fmt.Errorf("failed to delete tasks by save id: %w", err)
	}

	return nil
}
