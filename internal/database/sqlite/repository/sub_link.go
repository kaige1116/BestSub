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

// SubLinkRepository 链接数据访问实现
type SubLinkRepository struct {
	db *database.Database
}

// newSubLinkRepository 创建链接仓库
func newSubLinkRepository(db *database.Database) interfaces.SubLinkRepository {
	return &SubLinkRepository{db: db}
}

// Create 创建链接
func (r *SubLinkRepository) Create(ctx context.Context, link *models.SubLink) error {
	query := `INSERT INTO sub_links (name, url, user_agent, is_enabled, use_proxy, cron_expr, 
	          last_update, last_status, error_msg, created_at, updated_at) 
	          VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`

	now := timeutils.Now()
	result, err := r.db.ExecContext(ctx, query,
		link.Name,
		link.URL,
		link.UserAgent,
		link.IsEnabled,
		link.UseProxy,
		link.CronExpr,
		link.LastUpdate,
		link.LastStatus,
		link.ErrorMsg,
		now,
		now,
	)

	if err != nil {
		return fmt.Errorf("failed to create sub link: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("failed to get sub link id: %w", err)
	}

	link.ID = id
	link.CreatedAt = now
	link.UpdatedAt = now

	return nil
}

// GetByID 根据ID获取链接
func (r *SubLinkRepository) GetByID(ctx context.Context, id int64) (*models.SubLink, error) {
	query := `SELECT id, name, url, user_agent, is_enabled, use_proxy, cron_expr, 
	          last_update, last_status, error_msg, created_at, updated_at 
	          FROM sub_links WHERE id = ?`

	var link models.SubLink
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&link.ID,
		&link.Name,
		&link.URL,
		&link.UserAgent,
		&link.IsEnabled,
		&link.UseProxy,
		&link.CronExpr,
		&link.LastUpdate,
		&link.LastStatus,
		&link.ErrorMsg,
		&link.CreatedAt,
		&link.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get sub link by id: %w", err)
	}

	return &link, nil
}

// GetByURL 根据URL获取链接
func (r *SubLinkRepository) GetByURL(ctx context.Context, url string) (*models.SubLink, error) {
	query := `SELECT id, name, url, user_agent, is_enabled, use_proxy, cron_expr, 
	          last_update, last_status, error_msg, created_at, updated_at 
	          FROM sub_links WHERE url = ?`

	var link models.SubLink
	err := r.db.QueryRowContext(ctx, query, url).Scan(
		&link.ID,
		&link.Name,
		&link.URL,
		&link.UserAgent,
		&link.IsEnabled,
		&link.UseProxy,
		&link.CronExpr,
		&link.LastUpdate,
		&link.LastStatus,
		&link.ErrorMsg,
		&link.CreatedAt,
		&link.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get sub link by url: %w", err)
	}

	return &link, nil
}

// Update 更新链接
func (r *SubLinkRepository) Update(ctx context.Context, link *models.SubLink) error {
	query := `UPDATE sub_links SET name = ?, url = ?, user_agent = ?, is_enabled = ?, use_proxy = ?, 
	          cron_expr = ?, last_update = ?, last_status = ?, error_msg = ?, updated_at = ? WHERE id = ?`

	_, err := r.db.ExecContext(ctx, query,
		link.Name,
		link.URL,
		link.UserAgent,
		link.IsEnabled,
		link.UseProxy,
		link.CronExpr,
		link.LastUpdate,
		link.LastStatus,
		link.ErrorMsg,
		timeutils.Now(),
		link.ID,
	)

	if err != nil {
		return fmt.Errorf("failed to update sub link: %w", err)
	}

	return nil
}

// Delete 删除链接
func (r *SubLinkRepository) Delete(ctx context.Context, id int64) error {
	query := `DELETE FROM sub_links WHERE id = ?`

	_, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete sub link: %w", err)
	}

	return nil
}

// List 获取链接列表
func (r *SubLinkRepository) List(ctx context.Context, offset, limit int) ([]*models.SubLink, error) {
	query := `SELECT id, name, url, user_agent, is_enabled, use_proxy, cron_expr, 
	          last_update, last_status, error_msg, created_at, updated_at 
	          FROM sub_links ORDER BY created_at DESC LIMIT ? OFFSET ?`

	return r.querySubLinks(ctx, query, limit, offset)
}

// ListEnabled 获取启用的链接列表
func (r *SubLinkRepository) ListEnabled(ctx context.Context) ([]*models.SubLink, error) {
	query := `SELECT id, name, url, user_agent, is_enabled, use_proxy, cron_expr, 
	          last_update, last_status, error_msg, created_at, updated_at 
	          FROM sub_links WHERE is_enabled = true ORDER BY created_at DESC`

	return r.querySubLinks(ctx, query)
}

// Count 获取链接总数
func (r *SubLinkRepository) Count(ctx context.Context) (int64, error) {
	query := `SELECT COUNT(*) FROM sub_links`

	var count int64
	err := r.db.QueryRowContext(ctx, query).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to count sub links: %w", err)
	}

	return count, nil
}

// UpdateStatus 更新链接状态
func (r *SubLinkRepository) UpdateStatus(ctx context.Context, id int64, status, errorMsg string) error {
	query := `UPDATE sub_links SET last_status = ?, error_msg = ?, last_update = ?, updated_at = ? WHERE id = ?`

	now := timeutils.Now()
	_, err := r.db.ExecContext(ctx, query, status, errorMsg, now, now, id)
	if err != nil {
		return fmt.Errorf("failed to update sub link status: %w", err)
	}

	return nil
}

// querySubLinks 通用链接查询方法
func (r *SubLinkRepository) querySubLinks(ctx context.Context, query string, args ...interface{}) ([]*models.SubLink, error) {
	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to query sub links: %w", err)
	}
	defer rows.Close()

	var links []*models.SubLink
	for rows.Next() {
		var link models.SubLink
		err := rows.Scan(
			&link.ID,
			&link.Name,
			&link.URL,
			&link.UserAgent,
			&link.IsEnabled,
			&link.UseProxy,
			&link.CronExpr,
			&link.LastUpdate,
			&link.LastStatus,
			&link.ErrorMsg,
			&link.CreatedAt,
			&link.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan sub link: %w", err)
		}
		links = append(links, &link)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to iterate sub links: %w", err)
	}

	return links, nil
}

// SubLinkModuleConfigRepository 模块配置数据访问实现
type SubLinkModuleConfigRepository struct {
	db *database.Database
}

// newSubLinkModuleConfigRepository 创建模块配置仓库
func newSubLinkModuleConfigRepository(db *database.Database) interfaces.SubLinkModuleConfigRepository {
	return &SubLinkModuleConfigRepository{db: db}
}

// Create 创建模块配置
func (r *SubLinkModuleConfigRepository) Create(ctx context.Context, config *models.SubLinkModuleConfig) error {
	query := `INSERT INTO sub_link_module_configs (sub_link_id, module_type, module_name, is_enabled, priority, config, created_at, updated_at) 
	          VALUES (?, ?, ?, ?, ?, ?, ?, ?)`

	now := timeutils.Now()
	result, err := r.db.ExecContext(ctx, query,
		config.SubLinkID,
		config.ModuleType,
		config.ModuleName,
		config.IsEnabled,
		config.Priority,
		config.Config,
		now,
		now,
	)

	if err != nil {
		return fmt.Errorf("failed to create module config: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("failed to get module config id: %w", err)
	}

	config.ID = id
	config.CreatedAt = now
	config.UpdatedAt = now

	return nil
}

// GetByID 根据ID获取模块配置
func (r *SubLinkModuleConfigRepository) GetByID(ctx context.Context, id int64) (*models.SubLinkModuleConfig, error) {
	query := `SELECT id, sub_link_id, module_type, module_name, is_enabled, priority, config, created_at, updated_at 
	          FROM sub_link_module_configs WHERE id = ?`

	var config models.SubLinkModuleConfig
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&config.ID,
		&config.SubLinkID,
		&config.ModuleType,
		&config.ModuleName,
		&config.IsEnabled,
		&config.Priority,
		&config.Config,
		&config.CreatedAt,
		&config.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get module config by id: %w", err)
	}

	return &config, nil
}

// Update 更新模块配置
func (r *SubLinkModuleConfigRepository) Update(ctx context.Context, config *models.SubLinkModuleConfig) error {
	query := `UPDATE sub_link_module_configs SET sub_link_id = ?, module_type = ?, module_name = ?, 
	          is_enabled = ?, priority = ?, config = ?, updated_at = ? WHERE id = ?`

	_, err := r.db.ExecContext(ctx, query,
		config.SubLinkID,
		config.ModuleType,
		config.ModuleName,
		config.IsEnabled,
		config.Priority,
		config.Config,
		timeutils.Now(),
		config.ID,
	)

	if err != nil {
		return fmt.Errorf("failed to update module config: %w", err)
	}

	return nil
}

// Delete 删除模块配置
func (r *SubLinkModuleConfigRepository) Delete(ctx context.Context, id int64) error {
	query := `DELETE FROM sub_link_module_configs WHERE id = ?`

	_, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete module config: %w", err)
	}

	return nil
}

// DeleteByLinkID 删除链接的所有模块配置
func (r *SubLinkModuleConfigRepository) DeleteByLinkID(ctx context.Context, linkID int64) error {
	query := `DELETE FROM sub_link_module_configs WHERE sub_link_id = ?`

	_, err := r.db.ExecContext(ctx, query, linkID)
	if err != nil {
		return fmt.Errorf("failed to delete module configs by link id: %w", err)
	}

	return nil
}

// List 获取模块配置列表
func (r *SubLinkModuleConfigRepository) List(ctx context.Context, offset, limit int) ([]*models.SubLinkModuleConfig, error) {
	query := `SELECT id, sub_link_id, module_type, module_name, is_enabled, priority, config, created_at, updated_at 
	          FROM sub_link_module_configs ORDER BY sub_link_id, module_type, priority LIMIT ? OFFSET ?`

	return r.queryModuleConfigs(ctx, query, limit, offset)
}

// ListByLinkID 根据链接ID获取模块配置列表
func (r *SubLinkModuleConfigRepository) ListByLinkID(ctx context.Context, linkID int64) ([]*models.SubLinkModuleConfig, error) {
	query := `SELECT id, sub_link_id, module_type, module_name, is_enabled, priority, config, created_at, updated_at 
	          FROM sub_link_module_configs WHERE sub_link_id = ? ORDER BY module_type, priority`

	return r.queryModuleConfigs(ctx, query, linkID)
}

// ListByLinkIDAndType 根据链接ID和模块类型获取配置列表
func (r *SubLinkModuleConfigRepository) ListByLinkIDAndType(ctx context.Context, linkID int64, moduleType string) ([]*models.SubLinkModuleConfig, error) {
	query := `SELECT id, sub_link_id, module_type, module_name, is_enabled, priority, config, created_at, updated_at 
	          FROM sub_link_module_configs WHERE sub_link_id = ? AND module_type = ? ORDER BY priority`

	return r.queryModuleConfigs(ctx, query, linkID, moduleType)
}

// ListEnabledByLinkID 获取链接的启用模块配置列表
func (r *SubLinkModuleConfigRepository) ListEnabledByLinkID(ctx context.Context, linkID int64) ([]*models.SubLinkModuleConfig, error) {
	query := `SELECT id, sub_link_id, module_type, module_name, is_enabled, priority, config, created_at, updated_at 
	          FROM sub_link_module_configs WHERE sub_link_id = ? AND is_enabled = true ORDER BY module_type, priority`

	return r.queryModuleConfigs(ctx, query, linkID)
}

// Count 获取模块配置总数
func (r *SubLinkModuleConfigRepository) Count(ctx context.Context) (int64, error) {
	query := `SELECT COUNT(*) FROM sub_link_module_configs`

	var count int64
	err := r.db.QueryRowContext(ctx, query).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to count module configs: %w", err)
	}

	return count, nil
}

// CountByLinkID 获取链接的模块配置总数
func (r *SubLinkModuleConfigRepository) CountByLinkID(ctx context.Context, linkID int64) (int64, error) {
	query := `SELECT COUNT(*) FROM sub_link_module_configs WHERE sub_link_id = ?`

	var count int64
	err := r.db.QueryRowContext(ctx, query, linkID).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to count module configs by link id: %w", err)
	}

	return count, nil
}

// queryModuleConfigs 通用模块配置查询方法
func (r *SubLinkModuleConfigRepository) queryModuleConfigs(ctx context.Context, query string, args ...interface{}) ([]*models.SubLinkModuleConfig, error) {
	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to query module configs: %w", err)
	}
	defer rows.Close()

	var configs []*models.SubLinkModuleConfig
	for rows.Next() {
		var config models.SubLinkModuleConfig
		err := rows.Scan(
			&config.ID,
			&config.SubLinkID,
			&config.ModuleType,
			&config.ModuleName,
			&config.IsEnabled,
			&config.Priority,
			&config.Config,
			&config.CreatedAt,
			&config.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan module config: %w", err)
		}
		configs = append(configs, &config)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to iterate module configs: %w", err)
	}

	return configs, nil
}
