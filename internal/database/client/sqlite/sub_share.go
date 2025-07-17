package sqlite

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/bestruirui/bestsub/internal/database/interfaces"
	"github.com/bestruirui/bestsub/internal/models/sub"
	"github.com/bestruirui/bestsub/internal/utils/local"
)

// SubShareRepository 分享链接数据访问实现
type SubShareRepository struct {
	db *DB
}

// newSubShareLinkRepository 创建分享链接仓库
func (db *DB) SubShareLink() interfaces.SubShareRepository {
	return &SubShareRepository{db: db}
}

// Create 创建分享链接
func (r *SubShareRepository) Create(ctx context.Context, shareLink *sub.Share) error {
	query := `INSERT INTO sub_share_links (enable, name, description, rename, access_count, max_access_count, token, expires, created_at, updated_at)
	          VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`

	now := local.Time()
	result, err := r.db.db.ExecContext(ctx, query,
		shareLink.Enable,
		shareLink.Name,
		shareLink.Description,
		shareLink.Rename,
		shareLink.AccessCount,
		shareLink.MaxAccessCount,
		shareLink.Token,
		shareLink.Expires,
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
func (r *SubShareRepository) GetByID(ctx context.Context, id int64) (*sub.Share, error) {
	query := `SELECT id, enable, name, description, rename, access_count, max_access_count, token, expires, created_at, updated_at
	          FROM sub_share_links WHERE id = ?`

	var shareLink sub.Share
	err := r.db.db.QueryRowContext(ctx, query, id).Scan(
		&shareLink.ID,
		&shareLink.Enable,
		&shareLink.Name,
		&shareLink.Description,
		&shareLink.Rename,
		&shareLink.AccessCount,
		&shareLink.MaxAccessCount,
		&shareLink.Token,
		&shareLink.Expires,
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

// Update 更新分享链接
func (r *SubShareRepository) Update(ctx context.Context, shareLink *sub.Share) error {
	query := `UPDATE sub_share_links SET enable = ?, name = ?, description = ?, rename = ?, access_count = ?, max_access_count = ?, token = ?, expires = ?, updated_at = ? WHERE id = ?`

	_, err := r.db.db.ExecContext(ctx, query,
		shareLink.Enable,
		shareLink.Name,
		shareLink.Description,
		shareLink.Rename,
		shareLink.AccessCount,
		shareLink.MaxAccessCount,
		shareLink.Token,
		shareLink.Expires,
		local.Time(),
		shareLink.ID,
	)

	if err != nil {
		return fmt.Errorf("failed to update share link: %w", err)
	}

	return nil
}

// Delete 删除分享链接
func (r *SubShareRepository) Delete(ctx context.Context, id int64) error {
	query := `DELETE FROM sub_share_links WHERE id = ?`

	_, err := r.db.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete share link: %w", err)
	}

	return nil
}

// List 获取分享链接列表
func (r *SubShareRepository) List(ctx context.Context, offset, limit int) ([]*sub.Share, error) {
	query := `SELECT id, enable, name, description, rename, access_count, max_access_count, token, expires, created_at, updated_at
	          FROM sub_share_links ORDER BY created_at DESC LIMIT ? OFFSET ?`

	rows, err := r.db.db.QueryContext(ctx, query, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to list share links: %w", err)
	}
	defer rows.Close()

	var shareLinks []*sub.Share
	for rows.Next() {
		var shareLink sub.Share
		err := rows.Scan(
			&shareLink.ID,
			&shareLink.Enable,
			&shareLink.Name,
			&shareLink.Description,
			&shareLink.Rename,
			&shareLink.AccessCount,
			&shareLink.MaxAccessCount,
			&shareLink.Token,
			&shareLink.Expires,
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

// Count 获取分享链接总数
func (r *SubShareRepository) Count(ctx context.Context) (int64, error) {
	query := `SELECT COUNT(*) FROM sub_share_links`

	var count int64
	err := r.db.db.QueryRowContext(ctx, query).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to count share links: %w", err)
	}

	return count, nil
}

// AddOutputTemplateRelation 添加分享链接与输出模板的关联
func (r *SubShareRepository) AddOutputTemplateRelation(ctx context.Context, shareID, templateID int64) error {
	query := `INSERT OR IGNORE INTO share_template_relations (share_id, template_id) VALUES (?, ?)`

	_, err := r.db.db.ExecContext(ctx, query, shareID, templateID)
	if err != nil {
		return fmt.Errorf("failed to add output template relation: %w", err)
	}

	return nil
}

// AddFilterConfigRelation 添加分享链接与过滤配置的关联
func (r *SubShareRepository) AddFilterConfigRelation(ctx context.Context, shareID, configID int64) error {
	query := `INSERT OR IGNORE INTO share_fitter_relations (share_id, fitter_id) VALUES (?, ?)`

	_, err := r.db.db.ExecContext(ctx, query, shareID, configID)
	if err != nil {
		return fmt.Errorf("failed to add filter config relation: %w", err)
	}

	return nil
}

// AddSubRelation 添加分享链接与订阅的关联
func (r *SubShareRepository) AddSubRelation(ctx context.Context, shareID, subID int64) error {
	query := `INSERT OR IGNORE INTO share_sub_relations (share_id, sub_id) VALUES (?, ?)`

	_, err := r.db.db.ExecContext(ctx, query, shareID, subID)
	if err != nil {
		return fmt.Errorf("failed to add sub relation: %w", err)
	}

	return nil
}
