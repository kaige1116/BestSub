package sqlite

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/bestruirui/bestsub/internal/database/interfaces"
	"github.com/bestruirui/bestsub/internal/models/task"
	"github.com/bestruirui/bestsub/internal/utils/log"
)

type TaskRepository struct {
	db *DB
}

func (db *DB) Task() interfaces.TaskRepository {
	return &TaskRepository{db: db}
}

func (r *TaskRepository) Create(ctx context.Context, t *task.Data) error {
	log.Debugf("Create task")
	query := `INSERT INTO task (enable, name, config, extra, result)
	          VALUES (?, ?, ?, ?, ?)`

	result, err := r.db.db.ExecContext(ctx, query,
		t.Enable,
		t.Name,
		t.Config,
		t.Extra,
		t.Result,
	)

	if err != nil {
		return fmt.Errorf("failed to create task: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("failed to get task id: %w", err)
	}
	t.ID = uint16(id)
	return nil
}

func (r *TaskRepository) GetByID(ctx context.Context, id uint16) (*task.Data, error) {
	log.Debugf("Get task by id")
	query := `SELECT id, enable, name, config, extra, result
	          FROM task WHERE id = ?`

	var t task.Data
	err := r.db.db.QueryRowContext(ctx, query, id).Scan(
		&t.ID,
		&t.Enable,
		&t.Name,
		&t.Config,
		&t.Extra,
		&t.Result,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get task by id: %w", err)
	}

	return &t, nil
}

func (r *TaskRepository) Update(ctx context.Context, t *task.Data) error {
	log.Debugf("Update task")
	query := `UPDATE task SET enable = ?, name = ?, config = ?, extra = ?, result = ? WHERE id = ?`

	_, err := r.db.db.ExecContext(ctx, query,
		t.Enable,
		t.Name,
		t.Config,
		t.Extra,
		t.Result,
		t.ID,
	)

	if err != nil {
		return fmt.Errorf("failed to update task: %w", err)
	}

	return nil
}

func (r *TaskRepository) Delete(ctx context.Context, id uint16) error {
	log.Debugf("Delete task")
	query := `DELETE FROM task WHERE id = ?`

	_, err := r.db.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete task: %w", err)
	}

	return nil
}

func (r *TaskRepository) List(ctx context.Context) (*[]task.Data, error) {
	log.Debugf("List task")
	query := `SELECT id, enable, name, config, extra, result
	          FROM task ORDER BY id DESC`

	rows, err := r.db.db.QueryContext(ctx, query)
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
			&t.Config,
			&t.Extra,
			&t.Result,
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
