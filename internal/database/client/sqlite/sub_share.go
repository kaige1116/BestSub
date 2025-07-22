package sqlite

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/bestruirui/bestsub/internal/database/interfaces"
	"github.com/bestruirui/bestsub/internal/models/sub"
)

// SubShareRepository 分享链接数据访问实现
type SubShareRepository struct {
	db *DB
}

// newSubShareLinkRepository 创建分享链接仓库
func (db *DB) SubShare() interfaces.SubShareRepository {
	return &SubShareRepository{db: db}
}

// Create 创建分享链接
func (r *SubShareRepository) Create(ctx context.Context, shareLink *sub.Share) error {
	query := `INSERT INTO sub_share_links (enable, name, rename, access_count, max_access_count, token, expires)
	          VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`

	result, err := r.db.db.ExecContext(ctx, query,
		shareLink.Enable,
		shareLink.Name,
		shareLink.Rename,
		shareLink.AccessCount,
		shareLink.MaxAccessCount,
		shareLink.Token,
		shareLink.Expires,
	)

	if err != nil {
		return fmt.Errorf("failed to create share link: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("failed to get share link id: %w", err)
	}

	shareLink.ID = uint16(id)

	return nil
}

// GetByID 根据ID获取分享链接
func (r *SubShareRepository) GetByID(ctx context.Context, id uint16) (*sub.Share, error) {
	query := `SELECT id, enable, name, rename, access_count, max_access_count, token, expires
	          FROM sub_share_links WHERE id = ?`

	var shareLink sub.Share
	err := r.db.db.QueryRowContext(ctx, query, id).Scan(
		&shareLink.ID,
		&shareLink.Enable,
		&shareLink.Name,
		&shareLink.Rename,
		&shareLink.AccessCount,
		&shareLink.MaxAccessCount,
		&shareLink.Token,
		&shareLink.Expires,
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
	query := `UPDATE sub_share_links SET enable = ?, name = ?, rename = ?, access_count = ?, max_access_count = ?, token = ?, expires = ? WHERE id = ?`

	_, err := r.db.db.ExecContext(ctx, query,
		shareLink.Enable,
		shareLink.Name,
		shareLink.Rename,
		shareLink.AccessCount,
		shareLink.MaxAccessCount,
		shareLink.Token,
		shareLink.Expires,
		shareLink.ID,
	)

	if err != nil {
		return fmt.Errorf("failed to update share link: %w", err)
	}

	return nil
}

// Delete 删除分享链接
func (r *SubShareRepository) Delete(ctx context.Context, id uint16) error {
	query := `DELETE FROM sub_share_links WHERE id = ?`

	_, err := r.db.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete share link: %w", err)
	}

	return nil
}

// List 获取分享链接列表
func (r *SubShareRepository) List(ctx context.Context, offset, limit int) ([]*sub.Share, error) {
	query := `SELECT id, enable, name, rename, access_count, max_access_count, token, expires
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
			&shareLink.Rename,
			&shareLink.AccessCount,
			&shareLink.MaxAccessCount,
			&shareLink.Token,
			&shareLink.Expires,
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
func (r *SubShareRepository) Count(ctx context.Context) (uint16, error) {
	query := `SELECT COUNT(*) FROM sub_share_links`

	var count uint16
	err := r.db.db.QueryRowContext(ctx, query).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to count share links: %w", err)
	}

	return count, nil
}

// AddOutputTemplateRelation 添加分享链接与输出模板的关联
func (r *SubShareRepository) AddOutputTemplateRelation(ctx context.Context, shareID, templateID uint16) error {
	query := `INSERT OR IGNORE INTO share_template_relations (share_id, template_id) VALUES (?, ?)`

	_, err := r.db.db.ExecContext(ctx, query, shareID, templateID)
	if err != nil {
		return fmt.Errorf("failed to add output template relation: %w", err)
	}

	return nil
}

// AddFilterConfigRelation 添加分享链接与过滤配置的关联
func (r *SubShareRepository) AddFilterConfigRelation(ctx context.Context, shareID, configID uint16) error {
	query := `INSERT OR IGNORE INTO share_fitter_relations (share_id, fitter_id) VALUES (?, ?)`

	_, err := r.db.db.ExecContext(ctx, query, shareID, configID)
	if err != nil {
		return fmt.Errorf("failed to add filter config relation: %w", err)
	}

	return nil
}

// AddSubRelation 添加分享链接与订阅的关联
func (r *SubShareRepository) AddSubRelation(ctx context.Context, shareID, subID uint16) error {
	query := `INSERT OR IGNORE INTO share_sub_relations (share_id, sub_id) VALUES (?, ?)`

	_, err := r.db.db.ExecContext(ctx, query, shareID, subID)
	if err != nil {
		return fmt.Errorf("failed to add sub relation: %w", err)
	}

	return nil
}
