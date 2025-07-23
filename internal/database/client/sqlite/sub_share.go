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
	query := `INSERT INTO sub_share (enable, name, config)
	          VALUES (?, ?, ?)`

	result, err := r.db.db.ExecContext(ctx, query,
		shareLink.Enable,
		shareLink.Name,
		shareLink.Config,
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
	query := `SELECT id, enable, name, config
	          FROM sub_share WHERE id = ?`

	var shareLink sub.Share
	err := r.db.db.QueryRowContext(ctx, query, id).Scan(
		&shareLink.ID,
		&shareLink.Enable,
		&shareLink.Name,
		&shareLink.Config,
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
	query := `UPDATE sub_share SET enable = ?, name = ?, config = ? WHERE id = ?`

	_, err := r.db.db.ExecContext(ctx, query,
		shareLink.Enable,
		shareLink.Name,
		shareLink.Config,
		shareLink.ID,
	)

	if err != nil {
		return fmt.Errorf("failed to update share link: %w", err)
	}

	return nil
}

// Delete 删除分享链接
func (r *SubShareRepository) Delete(ctx context.Context, id uint16) error {
	query := `DELETE FROM sub_share WHERE id = ?`

	_, err := r.db.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete share link: %w", err)
	}

	return nil
}

// List 获取分享链接列表
func (r *SubShareRepository) List(ctx context.Context) ([]*sub.Share, error) {
	query := `SELECT id, enable, name, config
	          FROM sub_share`

	rows, err := r.db.db.QueryContext(ctx, query)
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
			&shareLink.Config,
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
