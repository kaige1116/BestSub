package sqlite

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/bestruirui/bestsub/internal/database/interfaces"
	"github.com/bestruirui/bestsub/internal/models/sub"
)

type SubRepository struct {
	db *DB
}

func (db *DB) Sub() interfaces.SubRepository {
	return &SubRepository{db: db}
}

func (r *SubRepository) Create(ctx context.Context, link *sub.Data) error {
	query := `INSERT INTO subs (enable, name, url, created_at, updated_at)
	          VALUES (?, ?, ?, ?, ?, ?)`

	now := time.Now()
	result, err := r.db.db.ExecContext(ctx, query,
		link.Enable,
		link.Name,
		link.URL,
		now,
		now,
	)

	if err != nil {
		return fmt.Errorf("failed to create sub: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("failed to get sub id: %w", err)
	}

	link.ID = uint16(id)
	link.CreatedAt = now
	link.UpdatedAt = now

	return nil
}

func (r *SubRepository) GetByID(ctx context.Context, id uint16) (*sub.Data, error) {
	query := `SELECT id, enable, name, url, created_at, updated_at
	          FROM subs WHERE id = ?`

	var s sub.Data
	var enable bool
	err := r.db.db.QueryRowContext(ctx, query, id).Scan(
		&s.ID,
		&enable,
		&s.Name,
		&s.URL,
		&s.CreatedAt,
		&s.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get sub by id: %w", err)
	}

	return &s, nil
}

// Update 更新订阅链接
func (r *SubRepository) Update(ctx context.Context, link *sub.Data) error {
	query := `UPDATE subs SET enable = ?, name = ?, url = ?, updated_at = ? WHERE id = ?`

	_, err := r.db.db.ExecContext(ctx, query,
		link.Enable,
		link.Name,
		link.URL,
		time.Now(),
		link.ID,
	)

	if err != nil {
		return fmt.Errorf("failed to update sub: %w", err)
	}

	return nil
}

func (r *SubRepository) Delete(ctx context.Context, id uint16) error {
	query := `DELETE FROM subs WHERE id = ?`

	_, err := r.db.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete sub: %w", err)
	}

	return nil
}

func (r *SubRepository) List(ctx context.Context, offset, limit int) (*[]sub.Data, error) {
	query := `SELECT id, enable, name, url, created_at, updated_at
	          FROM subs ORDER BY created_at DESC LIMIT ? OFFSET ?`

	rows, err := r.db.db.QueryContext(ctx, query, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to list subs: %w", err)
	}
	defer rows.Close()

	var subs []sub.Data
	for rows.Next() {
		var s sub.Data
		var enable bool
		err := rows.Scan(
			&s.ID,
			&enable,
			&s.Name,
			&s.URL,
			&s.CreatedAt,
			&s.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan sub: %w", err)
		}
		subs = append(subs, s)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to iterate subs: %w", err)
	}

	return &subs, nil
}

func (r *SubRepository) Count(ctx context.Context) (uint16, error) {
	query := `SELECT COUNT(*) FROM subs`

	var count uint16
	err := r.db.db.QueryRowContext(ctx, query).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to count subs: %w", err)
	}

	return count, nil
}

func (r *SubRepository) GetByTaskID(ctx context.Context, taskID uint16) (uint16, error) {
	query := `SELECT sub_id FROM sub_task_relations WHERE task_id = ?`

	var subID uint16
	err := r.db.db.QueryRowContext(ctx, query, taskID).Scan(&subID)
	if err != nil {
		if err == sql.ErrNoRows {
			return 0, fmt.Errorf("no sub found for task id %d", taskID)
		}
		return 0, fmt.Errorf("failed to get sub by task id: %w", err)
	}

	return subID, nil
}

func (r *SubRepository) GetByShareID(ctx context.Context, shareID uint16) ([]uint16, error) {
	query := `SELECT sub_id FROM share_sub_relations WHERE share_id = ?`

	rows, err := r.db.db.QueryContext(ctx, query, shareID)
	if err != nil {
		return nil, fmt.Errorf("failed to get subs by share id: %w", err)
	}
	defer rows.Close()

	var subIDs []uint16
	for rows.Next() {
		var subID uint16
		err := rows.Scan(&subID)
		if err != nil {
			return nil, fmt.Errorf("failed to scan sub id: %w", err)
		}
		subIDs = append(subIDs, subID)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to iterate sub ids: %w", err)
	}

	return subIDs, nil
}

func (r *SubRepository) GetBySaveID(ctx context.Context, saveID uint16) ([]uint16, error) {
	query := `SELECT sub_id FROM save_sub_relations WHERE save_id = ?`

	rows, err := r.db.db.QueryContext(ctx, query, saveID)
	if err != nil {
		return nil, fmt.Errorf("failed to get subs by save id: %w", err)
	}
	defer rows.Close()

	var subIDs []uint16
	for rows.Next() {
		var subID uint16
		err := rows.Scan(&subID)
		if err != nil {
			return nil, fmt.Errorf("failed to scan sub id: %w", err)
		}
		subIDs = append(subIDs, subID)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to iterate sub ids: %w", err)
	}

	return subIDs, nil
}

func (r *SubRepository) AddTaskRelation(ctx context.Context, subID, taskID uint16) error {
	query := `INSERT OR IGNORE INTO sub_task_relations (sub_id, task_id) VALUES (?, ?)`

	_, err := r.db.db.ExecContext(ctx, query, subID, taskID)
	if err != nil {
		return fmt.Errorf("failed to add task relation: %w", err)
	}

	return nil
}

func (r *SubRepository) AddSaveRelation(ctx context.Context, subID, saveID uint16) error {
	query := `INSERT OR IGNORE INTO save_sub_relations (sub_id, save_id) VALUES (?, ?)`

	_, err := r.db.db.ExecContext(ctx, query, subID, saveID)
	if err != nil {
		return fmt.Errorf("failed to add save relation: %w", err)
	}

	return nil
}

func (r *SubRepository) AddShareRelation(ctx context.Context, subID, shareID uint16) error {
	query := `INSERT OR IGNORE INTO share_sub_relations (sub_id, share_id) VALUES (?, ?)`

	_, err := r.db.db.ExecContext(ctx, query, subID, shareID)
	if err != nil {
		return fmt.Errorf("failed to add share relation: %w", err)
	}

	return nil
}
