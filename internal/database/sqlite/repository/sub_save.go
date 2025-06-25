package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/bestruirui/bestsub/internal/database/interfaces"
	"github.com/bestruirui/bestsub/internal/database/models"
	"github.com/bestruirui/bestsub/internal/database/sqlite/database"
	timeutils "github.com/bestruirui/bestsub/internal/utils/time"
)

// SubSaveConfigRepository 保存配置数据访问实现
type SubSaveConfigRepository struct {
	db *database.Database
}

// newSubSaveConfigRepository 创建保存配置仓库
func newSubSaveConfigRepository(db *database.Database) interfaces.SubSaveConfigRepository {
	return &SubSaveConfigRepository{db: db}
}

// Create 创建保存配置
func (r *SubSaveConfigRepository) Create(ctx context.Context, config *models.SubSaveConfig) error {
	query := `INSERT INTO sub_save_configs (name, description, is_enabled, storage_id, output_template_id, 
	          node_filter_id, file_name, save_interval, last_save, last_status, error_msg, save_count, created_at, updated_at) 
	          VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`

	now := timeutils.Now()
	result, err := r.db.ExecContext(ctx, query,
		config.Name,
		config.Description,
		config.IsEnabled,
		config.StorageID,
		config.OutputTemplateID,
		config.NodeFilterID,
		config.FileName,
		config.SaveInterval,
		config.LastSave,
		config.LastStatus,
		config.ErrorMsg,
		config.SaveCount,
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

	config.ID = id
	config.CreatedAt = now
	config.UpdatedAt = now

	return nil
}

// GetByID 根据ID获取保存配置
func (r *SubSaveConfigRepository) GetByID(ctx context.Context, id int64) (*models.SubSaveConfig, error) {
	query := `SELECT id, name, description, is_enabled, storage_id, output_template_id, 
	          node_filter_id, file_name, save_interval, last_save, last_status, error_msg, save_count, created_at, updated_at 
	          FROM sub_save_configs WHERE id = ?`

	var config models.SubSaveConfig
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&config.ID,
		&config.Name,
		&config.Description,
		&config.IsEnabled,
		&config.StorageID,
		&config.OutputTemplateID,
		&config.NodeFilterID,
		&config.FileName,
		&config.SaveInterval,
		&config.LastSave,
		&config.LastStatus,
		&config.ErrorMsg,
		&config.SaveCount,
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
func (r *SubSaveConfigRepository) Update(ctx context.Context, config *models.SubSaveConfig) error {
	query := `UPDATE sub_save_configs SET name = ?, description = ?, is_enabled = ?, storage_id = ?, 
	          output_template_id = ?, node_filter_id = ?, file_name = ?, save_interval = ?, 
	          last_save = ?, last_status = ?, error_msg = ?, save_count = ?, updated_at = ? WHERE id = ?`

	_, err := r.db.ExecContext(ctx, query,
		config.Name,
		config.Description,
		config.IsEnabled,
		config.StorageID,
		config.OutputTemplateID,
		config.NodeFilterID,
		config.FileName,
		config.SaveInterval,
		config.LastSave,
		config.LastStatus,
		config.ErrorMsg,
		config.SaveCount,
		timeutils.Now(),
		config.ID,
	)

	if err != nil {
		return fmt.Errorf("failed to update save config: %w", err)
	}

	return nil
}

// Delete 删除保存配置
func (r *SubSaveConfigRepository) Delete(ctx context.Context, id int64) error {
	query := `DELETE FROM sub_save_configs WHERE id = ?`

	_, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete save config: %w", err)
	}

	return nil
}

// List 获取保存配置列表
func (r *SubSaveConfigRepository) List(ctx context.Context, offset, limit int) ([]*models.SubSaveConfig, error) {
	query := `SELECT id, name, description, is_enabled, storage_id, output_template_id, 
	          node_filter_id, file_name, save_interval, last_save, last_status, error_msg, save_count, created_at, updated_at 
	          FROM sub_save_configs ORDER BY created_at DESC LIMIT ? OFFSET ?`

	return r.querySaveConfigs(ctx, query, limit, offset)
}

// ListEnabled 获取启用的保存配置列表
func (r *SubSaveConfigRepository) ListEnabled(ctx context.Context) ([]*models.SubSaveConfig, error) {
	query := `SELECT id, name, description, is_enabled, storage_id, output_template_id, 
	          node_filter_id, file_name, save_interval, last_save, last_status, error_msg, save_count, created_at, updated_at 
	          FROM sub_save_configs WHERE is_enabled = true ORDER BY created_at DESC`

	return r.querySaveConfigs(ctx, query)
}

// Count 获取保存配置总数
func (r *SubSaveConfigRepository) Count(ctx context.Context) (int64, error) {
	query := `SELECT COUNT(*) FROM sub_save_configs`

	var count int64
	err := r.db.QueryRowContext(ctx, query).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to count save configs: %w", err)
	}

	return count, nil
}

// UpdateStatus 更新保存状态
func (r *SubSaveConfigRepository) UpdateStatus(ctx context.Context, id int64, status, errorMsg string) error {
	query := `UPDATE sub_save_configs SET last_status = ?, error_msg = ?, last_save = ?, updated_at = ? WHERE id = ?`

	now := timeutils.Now()
	_, err := r.db.ExecContext(ctx, query, status, errorMsg, now, now, id)
	if err != nil {
		return fmt.Errorf("failed to update save config status: %w", err)
	}

	return nil
}

// IncrementSaveCount 增加保存次数
func (r *SubSaveConfigRepository) IncrementSaveCount(ctx context.Context, id int64) error {
	query := `UPDATE sub_save_configs SET save_count = save_count + 1, updated_at = ? WHERE id = ?`

	_, err := r.db.ExecContext(ctx, query, timeutils.Now(), id)
	if err != nil {
		return fmt.Errorf("failed to increment save count: %w", err)
	}

	return nil
}

// querySaveConfigs 通用保存配置查询方法
func (r *SubSaveConfigRepository) querySaveConfigs(ctx context.Context, query string, args ...interface{}) ([]*models.SubSaveConfig, error) {
	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to query save configs: %w", err)
	}
	defer rows.Close()

	var configs []*models.SubSaveConfig
	for rows.Next() {
		var config models.SubSaveConfig
		err := rows.Scan(
			&config.ID,
			&config.Name,
			&config.Description,
			&config.IsEnabled,
			&config.StorageID,
			&config.OutputTemplateID,
			&config.NodeFilterID,
			&config.FileName,
			&config.SaveInterval,
			&config.LastSave,
			&config.LastStatus,
			&config.ErrorMsg,
			&config.SaveCount,
			&config.CreatedAt,
			&config.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan save config: %w", err)
		}
		configs = append(configs, &config)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to iterate save configs: %w", err)
	}

	return configs, nil
}
