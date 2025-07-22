package sqlite

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/bestruirui/bestsub/internal/database/interfaces"
	"github.com/bestruirui/bestsub/internal/models/task"
)

// TaskRepository 任务数据访问实现
type TaskRepository struct {
	db *DB
}

// newTaskRepository 创建任务仓库
func (db *DB) Task() interfaces.TaskRepository {
	return &TaskRepository{db: db}
}

// Create 创建任务
func (r *TaskRepository) Create(ctx context.Context, t *task.Data) (uint16, error) {
	query := `INSERT INTO tasks (enable, name, system, cron, type, log_level, timeout, config, last_run_result, last_run_time, last_run_duration, success_count, failed_count)
	          VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`

	result, err := r.db.db.ExecContext(ctx, query,
		t.Enable,
		t.Name,
		t.System,
		t.Cron,
		t.Type,
		t.LogLevel,
		t.Timeout,
		t.Config,
		t.LastRunResult,
		t.LastRunTime,
		t.LastRunDuration,
		t.SuccessCount,
		t.FailedCount,
	)

	if err != nil {
		return 0, fmt.Errorf("failed to create task: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("failed to get task id: %w", err)
	}
	t.ID = uint16(id)
	return t.ID, nil
}

// GetByID 根据ID获取任务
func (r *TaskRepository) GetByID(ctx context.Context, id uint16) (*task.Data, error) {
	query := `SELECT id, enable, name, system, cron, type, log_level, timeout, config, last_run_result, last_run_time, last_run_duration, success_count, failed_count
	          FROM tasks WHERE id = ?`

	var t task.Data
	err := r.db.db.QueryRowContext(ctx, query, id).Scan(
		&t.ID,
		&t.Enable,
		&t.Name,
		&t.System,
		&t.Cron,
		&t.Type,
		&t.LogLevel,
		&t.Timeout,
		&t.Config,
		&t.LastRunResult,
		&t.LastRunTime,
		&t.LastRunDuration,
		&t.SuccessCount,
		&t.FailedCount,
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
	query := `UPDATE tasks SET enable = ?, name = ?, system = ?, cron = ?, type = ?, log_level = ?, timeout = ?, config = ?,
	          last_run_result = ?, last_run_time = ?, last_run_duration = ?, success_count = ?, failed_count = ? WHERE id = ?`

	_, err := r.db.db.ExecContext(ctx, query,
		t.Enable,
		t.Name,
		t.System,
		t.Cron,
		t.Type,
		t.LogLevel,
		t.Timeout,
		t.Config,
		t.LastRunResult,
		t.LastRunTime,
		t.LastRunDuration,
		t.SuccessCount,
		t.FailedCount,
		t.ID,
	)

	if err != nil {
		return fmt.Errorf("failed to update task: %w", err)
	}

	return nil
}

// Delete 删除任务
func (r *TaskRepository) Delete(ctx context.Context, id uint16) error {
	query := `DELETE FROM tasks WHERE id = ?`

	_, err := r.db.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete task: %w", err)
	}

	return nil
}

// List 获取任务列表
func (r *TaskRepository) List(ctx context.Context, offset, limit int) (*[]task.Data, error) {
	query := `SELECT id, enable, name, system, cron, type, log_level, timeout, config, last_run_result, last_run_time, last_run_duration, success_count, failed_count
	          FROM tasks ORDER BY created_at DESC LIMIT ? OFFSET ?`

	rows, err := r.db.db.QueryContext(ctx, query, limit, offset)
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
			&t.System,
			&t.Cron,
			&t.Type,
			&t.LogLevel,
			&t.Timeout,
			&t.Config,
			&t.LastRunResult,
			&t.LastRunTime,
			&t.LastRunDuration,
			&t.SuccessCount,
			&t.FailedCount,
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
func (r *TaskRepository) Count(ctx context.Context) (uint16, error) {
	query := `SELECT COUNT(*) FROM tasks`

	var count uint16
	err := r.db.db.QueryRowContext(ctx, query).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to count tasks: %w", err)
	}

	return count, nil
}

// GetSystemTasks 获取所有系统任务
func (r *TaskRepository) GetSystemTasks(ctx context.Context) (*[]task.Data, error) {
	query := `SELECT id, enable, name, system, cron, type, log_level, timeout, config, last_run_result, last_run_time, last_run_duration, success_count, failed_count
	          FROM tasks WHERE system = true ORDER BY created_at DESC`

	rows, err := r.db.db.QueryContext(ctx, query)
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
			&t.System,
			&t.Cron,
			&t.Type,
			&t.LogLevel,
			&t.Timeout,
			&t.Config,
			&t.LastRunResult,
			&t.LastRunTime,
			&t.LastRunDuration,
			&t.SuccessCount,
			&t.FailedCount,
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
func (r *TaskRepository) GetBySubID(ctx context.Context, subID uint16) (*[]task.Data, error) {
	query := `SELECT t.id, t.enable, t.name, t.system, t.cron, t.type, t.log_level, t.timeout, t.config, t.last_run_result, t.last_run_time, t.last_run_duration, t.success_count, t.failed_count
	          FROM tasks t
	          INNER JOIN sub_task_relations str ON t.id = str.task_id
	          WHERE str.sub_id = ? ORDER BY t.created_at DESC`

	rows, err := r.db.db.QueryContext(ctx, query, subID)
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
			&t.System,
			&t.Cron,
			&t.Type,
			&t.LogLevel,
			&t.Timeout,
			&t.Config,
			&t.LastRunResult,
			&t.LastRunTime,
			&t.LastRunDuration,
			&t.SuccessCount,
			&t.FailedCount,
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
func (r *TaskRepository) GetTaskIDsBySubID(ctx context.Context, subID uint16) ([]uint16, error) {
	query := `SELECT task_id FROM sub_task_relations WHERE sub_id = ? ORDER BY task_id`

	rows, err := r.db.db.QueryContext(ctx, query, subID)
	if err != nil {
		return nil, fmt.Errorf("failed to get task ids by sub id: %w", err)
	}
	defer rows.Close()

	var taskIDs []uint16
	for rows.Next() {
		var taskID uint16
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
func (r *TaskRepository) AddNotifyRelation(ctx context.Context, taskID, notifyID uint16) error {
	query := `INSERT OR IGNORE INTO task_notify_relations (task_id, notify_id) VALUES (?, ?)`

	_, err := r.db.db.ExecContext(ctx, query, taskID, notifyID)
	if err != nil {
		return fmt.Errorf("failed to add notify relation: %w", err)
	}

	return nil
}

// AddNotifyTemplateRelation 添加任务与通知模板的关联
func (r *TaskRepository) AddNotifyTemplateRelation(ctx context.Context, taskID, notifyTemplateID uint16) error {
	query := `INSERT OR IGNORE INTO task_notify_template_relations (task_id, notify_template_id) VALUES (?, ?)`

	_, err := r.db.db.ExecContext(ctx, query, taskID, notifyTemplateID)
	if err != nil {
		return fmt.Errorf("failed to add notify template relation: %w", err)
	}

	return nil
}

// DeleteBySubID 根据订阅ID删除所有的任务
func (r *TaskRepository) DeleteBySubID(ctx context.Context, subID uint16) error {
	query := `DELETE FROM tasks WHERE id IN (SELECT task_id FROM sub_task_relations WHERE sub_id = ?)`

	_, err := r.db.db.ExecContext(ctx, query, subID)
	if err != nil {
		return fmt.Errorf("failed to delete tasks by sub id: %w", err)
	}

	return nil
}

// DeleteBySaveID 根据保存ID删除所有的任务
func (r *TaskRepository) DeleteBySaveID(ctx context.Context, saveID uint16) error {
	query := `DELETE FROM tasks WHERE id IN (SELECT task_id FROM save_task_relations WHERE save_id = ?)`

	_, err := r.db.db.ExecContext(ctx, query, saveID)
	if err != nil {
		return fmt.Errorf("failed to delete tasks by save id: %w", err)
	}

	return nil
}
