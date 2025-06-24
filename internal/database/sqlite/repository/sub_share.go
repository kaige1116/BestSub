package repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/bestruirui/bestsub/internal/database/interfaces"
	"github.com/bestruirui/bestsub/internal/database/models"
	"github.com/bestruirui/bestsub/internal/database/sqlite/database"
)

// SubShareLinkRepository 分享链接数据访问实现
type SubShareLinkRepository struct {
	db *database.Database
}

// newSubShareLinkRepository 创建分享链接仓库
func newSubShareLinkRepository(db *database.Database) interfaces.SubShareLinkRepository {
	return &SubShareLinkRepository{db: db}
}

// Create 创建分享链接
func (r *SubShareLinkRepository) Create(ctx context.Context, shareLink *models.SubShareLink) error {
	query := `INSERT INTO sub_share_links (name, description, token, is_enabled, output_template_id, 
	          node_filter_id, expires_at, max_downloads, download_count, last_access, created_at, updated_at) 
	          VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`

	now := time.Now()
	result, err := r.db.ExecContext(ctx, query,
		shareLink.Name,
		shareLink.Description,
		shareLink.Token,
		shareLink.IsEnabled,
		shareLink.OutputTemplateID,
		shareLink.NodeFilterID,
		shareLink.ExpiresAt,
		shareLink.MaxDownloads,
		shareLink.DownloadCount,
		shareLink.LastAccess,
		now,
		now,
	)

	if err != nil {
		return fmt.Errorf("failed to create share link: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("failed to get share link id: %w", err)
	}

	shareLink.ID = id
	shareLink.CreatedAt = now
	shareLink.UpdatedAt = now

	return nil
}

// GetByID 根据ID获取分享链接
func (r *SubShareLinkRepository) GetByID(ctx context.Context, id int64) (*models.SubShareLink, error) {
	query := `SELECT id, name, description, token, is_enabled, output_template_id, 
	          node_filter_id, expires_at, max_downloads, download_count, last_access, created_at, updated_at 
	          FROM sub_share_links WHERE id = ?`

	var shareLink models.SubShareLink
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&shareLink.ID,
		&shareLink.Name,
		&shareLink.Description,
		&shareLink.Token,
		&shareLink.IsEnabled,
		&shareLink.OutputTemplateID,
		&shareLink.NodeFilterID,
		&shareLink.ExpiresAt,
		&shareLink.MaxDownloads,
		&shareLink.DownloadCount,
		&shareLink.LastAccess,
		&shareLink.CreatedAt,
		&shareLink.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get share link by id: %w", err)
	}

	return &shareLink, nil
}

// GetByToken 根据Token获取分享链接
func (r *SubShareLinkRepository) GetByToken(ctx context.Context, token string) (*models.SubShareLink, error) {
	query := `SELECT id, name, description, token, is_enabled, output_template_id, 
	          node_filter_id, expires_at, max_downloads, download_count, last_access, created_at, updated_at 
	          FROM sub_share_links WHERE token = ?`

	var shareLink models.SubShareLink
	err := r.db.QueryRowContext(ctx, query, token).Scan(
		&shareLink.ID,
		&shareLink.Name,
		&shareLink.Description,
		&shareLink.Token,
		&shareLink.IsEnabled,
		&shareLink.OutputTemplateID,
		&shareLink.NodeFilterID,
		&shareLink.ExpiresAt,
		&shareLink.MaxDownloads,
		&shareLink.DownloadCount,
		&shareLink.LastAccess,
		&shareLink.CreatedAt,
		&shareLink.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get share link by token: %w", err)
	}

	return &shareLink, nil
}

// Update 更新分享链接
func (r *SubShareLinkRepository) Update(ctx context.Context, shareLink *models.SubShareLink) error {
	query := `UPDATE sub_share_links SET name = ?, description = ?, token = ?, is_enabled = ?, 
	          output_template_id = ?, node_filter_id = ?, expires_at = ?, max_downloads = ?, 
	          download_count = ?, last_access = ?, updated_at = ? WHERE id = ?`

	_, err := r.db.ExecContext(ctx, query,
		shareLink.Name,
		shareLink.Description,
		shareLink.Token,
		shareLink.IsEnabled,
		shareLink.OutputTemplateID,
		shareLink.NodeFilterID,
		shareLink.ExpiresAt,
		shareLink.MaxDownloads,
		shareLink.DownloadCount,
		shareLink.LastAccess,
		time.Now(),
		shareLink.ID,
	)

	if err != nil {
		return fmt.Errorf("failed to update share link: %w", err)
	}

	return nil
}

// Delete 删除分享链接
func (r *SubShareLinkRepository) Delete(ctx context.Context, id int64) error {
	query := `DELETE FROM sub_share_links WHERE id = ?`

	_, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete share link: %w", err)
	}

	return nil
}

// List 获取分享链接列表
func (r *SubShareLinkRepository) List(ctx context.Context, offset, limit int) ([]*models.SubShareLink, error) {
	query := `SELECT id, name, description, token, is_enabled, output_template_id, 
	          node_filter_id, expires_at, max_downloads, download_count, last_access, created_at, updated_at 
	          FROM sub_share_links ORDER BY created_at DESC LIMIT ? OFFSET ?`

	return r.queryShareLinks(ctx, query, limit, offset)
}

// ListEnabled 获取启用的分享链接列表
func (r *SubShareLinkRepository) ListEnabled(ctx context.Context) ([]*models.SubShareLink, error) {
	query := `SELECT id, name, description, token, is_enabled, output_template_id, 
	          node_filter_id, expires_at, max_downloads, download_count, last_access, created_at, updated_at 
	          FROM sub_share_links WHERE is_enabled = true AND (expires_at IS NULL OR expires_at > ?) 
	          ORDER BY created_at DESC`

	return r.queryShareLinks(ctx, query, time.Now())
}

// Count 获取分享链接总数
func (r *SubShareLinkRepository) Count(ctx context.Context) (int64, error) {
	query := `SELECT COUNT(*) FROM sub_share_links`

	var count int64
	err := r.db.QueryRowContext(ctx, query).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to count share links: %w", err)
	}

	return count, nil
}

// IncrementDownloadCount 增加下载次数
func (r *SubShareLinkRepository) IncrementDownloadCount(ctx context.Context, id int64) error {
	query := `UPDATE sub_share_links SET download_count = download_count + 1, updated_at = ? WHERE id = ?`

	_, err := r.db.ExecContext(ctx, query, time.Now(), id)
	if err != nil {
		return fmt.Errorf("failed to increment download count: %w", err)
	}

	return nil
}

// UpdateLastAccess 更新最后访问时间
func (r *SubShareLinkRepository) UpdateLastAccess(ctx context.Context, id int64) error {
	query := `UPDATE sub_share_links SET last_access = ?, updated_at = ? WHERE id = ?`

	now := time.Now()
	_, err := r.db.ExecContext(ctx, query, now, now, id)
	if err != nil {
		return fmt.Errorf("failed to update last access: %w", err)
	}

	return nil
}

// DeleteExpired 删除过期的分享链接
func (r *SubShareLinkRepository) DeleteExpired(ctx context.Context) error {
	query := `DELETE FROM sub_share_links WHERE expires_at IS NOT NULL AND expires_at < ?`

	_, err := r.db.ExecContext(ctx, query, time.Now())
	if err != nil {
		return fmt.Errorf("failed to delete expired share links: %w", err)
	}

	return nil
}

// IsTokenUnique 检查Token是否唯一
func (r *SubShareLinkRepository) IsTokenUnique(ctx context.Context, token string) (bool, error) {
	query := `SELECT COUNT(*) FROM sub_share_links WHERE token = ?`

	var count int
	err := r.db.QueryRowContext(ctx, query, token).Scan(&count)
	if err != nil {
		return false, fmt.Errorf("failed to check token uniqueness: %w", err)
	}

	return count == 0, nil
}

// queryShareLinks 通用分享链接查询方法
func (r *SubShareLinkRepository) queryShareLinks(ctx context.Context, query string, args ...interface{}) ([]*models.SubShareLink, error) {
	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to query share links: %w", err)
	}
	defer rows.Close()

	var shareLinks []*models.SubShareLink
	for rows.Next() {
		var shareLink models.SubShareLink
		err := rows.Scan(
			&shareLink.ID,
			&shareLink.Name,
			&shareLink.Description,
			&shareLink.Token,
			&shareLink.IsEnabled,
			&shareLink.OutputTemplateID,
			&shareLink.NodeFilterID,
			&shareLink.ExpiresAt,
			&shareLink.MaxDownloads,
			&shareLink.DownloadCount,
			&shareLink.LastAccess,
			&shareLink.CreatedAt,
			&shareLink.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan share link: %w", err)
		}
		shareLinks = append(shareLinks, &shareLink)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to iterate share links: %w", err)
	}

	return shareLinks, nil
}
