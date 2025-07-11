package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/bestruirui/bestsub/internal/database/interfaces"
	"github.com/bestruirui/bestsub/internal/database/sqlite/database"
	"github.com/bestruirui/bestsub/internal/models/sub"
	"github.com/bestruirui/bestsub/internal/utils/local"
)

// SubRepository 订阅链接数据访问实现
type SubRepository struct {
	db *database.Database
}

// newSubRepository 创建订阅链接仓库
func newSubRepository(db *database.Database) interfaces.SubRepository {
	return &SubRepository{db: db}
}

// Create 创建订阅链接
func (r *SubRepository) Create(ctx context.Context, link *sub.Data) error {
	query := `INSERT INTO subs (enable, name, description, url, created_at, updated_at)
	          VALUES (?, ?, ?, ?, ?, ?)`

	now := local.Time()
	result, err := r.db.ExecContext(ctx, query,
		true, // enable默认为true
		link.Name,
		link.Description,
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

	link.ID = id
	link.CreatedAt = now
	link.UpdatedAt = now

	return nil
}

// GetByID 根据ID获取订阅链接
func (r *SubRepository) GetByID(ctx context.Context, id int64) (*sub.Data, error) {
	query := `SELECT id, enable, name, description, url, created_at, updated_at
	          FROM subs WHERE id = ?`

	var s sub.Data
	var enable bool
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&s.ID,
		&enable,
		&s.Name,
		&s.Description,
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
	query := `UPDATE subs SET enable = ?, name = ?, description = ?, url = ?, updated_at = ? WHERE id = ?`

	_, err := r.db.ExecContext(ctx, query,
		true, // enable默认为true
		link.Name,
		link.Description,
		link.URL,
		local.Time(),
		link.ID,
	)

	if err != nil {
		return fmt.Errorf("failed to update sub: %w", err)
	}

	return nil
}

// Delete 删除订阅链接
func (r *SubRepository) Delete(ctx context.Context, id int64) error {
	query := `DELETE FROM subs WHERE id = ?`

	_, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete sub: %w", err)
	}

	return nil
}

// List 获取订阅链接列表
func (r *SubRepository) List(ctx context.Context, offset, limit int) (*[]sub.Data, error) {
	query := `SELECT id, enable, name, description, url, created_at, updated_at
	          FROM subs ORDER BY created_at DESC LIMIT ? OFFSET ?`

	rows, err := r.db.QueryContext(ctx, query, limit, offset)
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
			&s.Description,
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

// Count 获取订阅链接总数
func (r *SubRepository) Count(ctx context.Context) (int64, error) {
	query := `SELECT COUNT(*) FROM subs`

	var count int64
	err := r.db.QueryRowContext(ctx, query).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to count subs: %w", err)
	}

	return count, nil
}

// GetByTaskID 根据任务ID获取订阅ID
func (r *SubRepository) GetByTaskID(ctx context.Context, taskID int64) (int64, error) {
	query := `SELECT sub_id FROM sub_task_relations WHERE task_id = ?`

	var subID int64
	err := r.db.QueryRowContext(ctx, query, taskID).Scan(&subID)
	if err != nil {
		if err == sql.ErrNoRows {
			return 0, fmt.Errorf("no sub found for task id %d", taskID)
		}
		return 0, fmt.Errorf("failed to get sub by task id: %w", err)
	}

	return subID, nil
}

// GetByShareID 根据分享ID获取订阅ID列表
func (r *SubRepository) GetByShareID(ctx context.Context, shareID int64) ([]int64, error) {
	query := `SELECT sub_id FROM share_sub_relations WHERE share_id = ?`

	rows, err := r.db.QueryContext(ctx, query, shareID)
	if err != nil {
		return nil, fmt.Errorf("failed to get subs by share id: %w", err)
	}
	defer rows.Close()

	var subIDs []int64
	for rows.Next() {
		var subID int64
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

// GetBySaveID 根据保存ID获取订阅ID列表
func (r *SubRepository) GetBySaveID(ctx context.Context, saveID int64) ([]int64, error) {
	query := `SELECT sub_id FROM save_sub_relations WHERE save_id = ?`

	rows, err := r.db.QueryContext(ctx, query, saveID)
	if err != nil {
		return nil, fmt.Errorf("failed to get subs by save id: %w", err)
	}
	defer rows.Close()

	var subIDs []int64
	for rows.Next() {
		var subID int64
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

// AddTaskRelation 添加任务与订阅的关联
func (r *SubRepository) AddTaskRelation(ctx context.Context, subID, taskID int64) error {
	query := `INSERT OR IGNORE INTO sub_task_relations (sub_id, task_id) VALUES (?, ?)`

	_, err := r.db.ExecContext(ctx, query, subID, taskID)
	if err != nil {
		return fmt.Errorf("failed to add task relation: %w", err)
	}

	return nil
}

// AddSaveRelation 添加保存配置与订阅的关联
func (r *SubRepository) AddSaveRelation(ctx context.Context, subID, saveID int64) error {
	query := `INSERT OR IGNORE INTO save_sub_relations (sub_id, save_id) VALUES (?, ?)`

	_, err := r.db.ExecContext(ctx, query, subID, saveID)
	if err != nil {
		return fmt.Errorf("failed to add save relation: %w", err)
	}

	return nil
}

// AddShareRelation 添加分享链接与订阅的关联
func (r *SubRepository) AddShareRelation(ctx context.Context, subID, shareID int64) error {
	query := `INSERT OR IGNORE INTO share_sub_relations (sub_id, share_id) VALUES (?, ?)`

	_, err := r.db.ExecContext(ctx, query, subID, shareID)
	if err != nil {
		return fmt.Errorf("failed to add share relation: %w", err)
	}

	return nil
}
