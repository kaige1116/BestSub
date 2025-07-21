package sqlite

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/bestruirui/bestsub/internal/database/interfaces"
	"github.com/bestruirui/bestsub/internal/models/sub"
)

// SubSaveConfigRepository 保存配置数据访问实现
type SubSaveConfigRepository struct {
	db *DB
}

// newSubSaveConfigRepository 创建保存配置仓库
func (db *DB) SubSave() interfaces.SubSaveRepository {
	return &SubSaveConfigRepository{db: db}
}

// Create 创建保存配置
func (r *SubSaveConfigRepository) Create(ctx context.Context, config *sub.SaveConfig) error {
	query := `INSERT INTO sub_save (enable, name, description, rename, file_name, created_at, updated_at)
	          VALUES (?, ?, ?, ?, ?, ?, ?)`

	now := time.Now()
	result, err := r.db.db.ExecContext(ctx, query,
		config.Enable,
		config.Name,
		config.Description,
		config.Rename,
		config.FileName,
		now,
		now,
	)

	if err != nil {
		return fmt.Errorf("failed to create save config: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("failed to get save config id: %w", err)
	}

	config.ID = uint16(id)
	config.CreatedAt = now
	config.UpdatedAt = now

	return nil
}

// GetByID 根据ID获取保存配置
func (r *SubSaveConfigRepository) GetByID(ctx context.Context, id uint16) (*sub.SaveConfig, error) {
	query := `SELECT id, enable, name, description, rename, file_name, created_at, updated_at
	          FROM sub_save WHERE id = ?`

	var config sub.SaveConfig
	err := r.db.db.QueryRowContext(ctx, query, id).Scan(
		&config.ID,
		&config.Enable,
		&config.Name,
		&config.Description,
		&config.Rename,
		&config.FileName,
		&config.CreatedAt,
		&config.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get save config by id: %w", err)
	}

	return &config, nil
}

// Update 更新保存配置
func (r *SubSaveConfigRepository) Update(ctx context.Context, config *sub.SaveConfig) error {
	query := `UPDATE sub_save SET enable = ?, name = ?, description = ?, rename = ?, file_name = ?, updated_at = ? WHERE id = ?`

	_, err := r.db.db.ExecContext(ctx, query,
		config.Enable,
		config.Name,
		config.Description,
		config.Rename,
		config.FileName,
		time.Now(),
		config.ID,
	)

	if err != nil {
		return fmt.Errorf("failed to update save config: %w", err)
	}

	return nil
}

// Delete 删除保存配置
func (r *SubSaveConfigRepository) Delete(ctx context.Context, id uint16) error {
	query := `DELETE FROM sub_save WHERE id = ?`

	_, err := r.db.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete save config: %w", err)
	}

	return nil
}

// GetByTaskID 根据任务ID获取保存配置列表
func (r *SubSaveConfigRepository) GetByTaskID(ctx context.Context, taskID uint16) (*[]sub.SaveConfig, error) {
	query := `SELECT sc.id, sc.enable, sc.name, sc.description, sc.rename, sc.file_name, sc.created_at, sc.updated_at
	          FROM sub_save sc
	          INNER JOIN save_task_relations str ON sc.id = str.save_id
	          WHERE str.task_id = ?`

	rows, err := r.db.db.QueryContext(ctx, query, taskID)
	if err != nil {
		return nil, fmt.Errorf("failed to get save configs by task id: %w", err)
	}
	defer rows.Close()

	var configs []sub.SaveConfig
	for rows.Next() {
		var config sub.SaveConfig
		err := rows.Scan(
			&config.ID,
			&config.Enable,
			&config.Name,
			&config.Description,
			&config.Rename,
			&config.FileName,
			&config.CreatedAt,
			&config.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan save config: %w", err)
		}
		configs = append(configs, config)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to iterate save configs: %w", err)
	}

	return &configs, nil
}

// AddTaskRelation 添加保存配置与任务的关联
func (r *SubSaveConfigRepository) AddTaskRelation(ctx context.Context, saveID, taskID uint16) error {
	query := `INSERT OR IGNORE INTO save_task_relations (save_id, task_id) VALUES (?, ?)`

	_, err := r.db.db.ExecContext(ctx, query, saveID, taskID)
	if err != nil {
		return fmt.Errorf("failed to add task relation: %w", err)
	}

	return nil
}

// AddOutputTemplateRelation 添加保存配置与输出模板的关联
func (r *SubSaveConfigRepository) AddOutputTemplateRelation(ctx context.Context, saveID, templateID uint16) error {
	query := `INSERT OR IGNORE INTO save_template_relations (save_id, template_id) VALUES (?, ?)`

	_, err := r.db.db.ExecContext(ctx, query, saveID, templateID)
	if err != nil {
		return fmt.Errorf("failed to add output template relation: %w", err)
	}

	return nil
}

// AddFilterConfigRelation 添加保存配置与过滤配置的关联
func (r *SubSaveConfigRepository) AddFilterConfigRelation(ctx context.Context, saveID, configID uint16) error {
	query := `INSERT OR IGNORE INTO save_fitter_relations (save_id, fitter_id) VALUES (?, ?)`

	_, err := r.db.db.ExecContext(ctx, query, saveID, configID)
	if err != nil {
		return fmt.Errorf("failed to add filter config relation: %w", err)
	}

	return nil
}

// AddSubRelation 添加保存配置与订阅的关联
func (r *SubSaveConfigRepository) AddSubRelation(ctx context.Context, saveID, subID uint16) error {
	query := `INSERT OR IGNORE INTO save_sub_relations (save_id, sub_id) VALUES (?, ?)`

	_, err := r.db.db.ExecContext(ctx, query, saveID, subID)
	if err != nil {
		return fmt.Errorf("failed to add sub relation: %w", err)
	}

	return nil
}

// AddStorageConfigRelation 添加保存配置与存储配置的关联
func (r *SubSaveConfigRepository) AddStorageConfigRelation(ctx context.Context, saveID, configID uint16) error {
	query := `INSERT OR IGNORE INTO save_storage_relations (save_id, storage_id) VALUES (?, ?)`

	_, err := r.db.db.ExecContext(ctx, query, saveID, configID)
	if err != nil {
		return fmt.Errorf("failed to add storage config relation: %w", err)
	}

	return nil
}
