package sqlite

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/bestruirui/bestsub/internal/database/interfaces"
	"github.com/bestruirui/bestsub/internal/models/auth"
	"github.com/bestruirui/bestsub/internal/utils/local"
)

// Get 获取认证信息
func (db *DB) Auth() interfaces.AuthRepository {
	return &AuthRepository{db: db}
}

// AuthRepository 认证数据访问实现
type AuthRepository struct {
	db *DB
}

// Get 获取认证信息
func (db *AuthRepository) Get(ctx context.Context) (*auth.Data, error) {
	query := `SELECT id, user_name, password, created_at, updated_at FROM auth LIMIT 1`

	var authData auth.Data
	err := db.db.db.QueryRowContext(ctx, query).Scan(
		&authData.ID,
		&authData.UserName,
		&authData.Password,
		&authData.CreatedAt,
		&authData.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get auth: %w", err)
	}

	return &authData, nil
}

// UpdateName 更新用户名
func (r *AuthRepository) UpdateName(ctx context.Context, name string) (time.Time, error) {
	query := `UPDATE auth SET user_name = ?, updated_at = ? WHERE id = (SELECT id FROM auth LIMIT 1)`

	now := local.Time()
	result, err := r.db.db.ExecContext(ctx, query, name, now)
	if err != nil {
		return time.Time{}, fmt.Errorf("failed to update username: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return time.Time{}, fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return time.Time{}, fmt.Errorf("no auth record found to update")
	}

	return now, nil
}

// UpdatePassword 更新密码
func (r *AuthRepository) UpdatePassword(ctx context.Context, hashPassword string) (time.Time, error) {
	query := `UPDATE auth SET password = ?, updated_at = ? WHERE id = (SELECT id FROM auth LIMIT 1)`

	now := local.Time()
	result, err := r.db.db.ExecContext(ctx, query, hashPassword, now)
	if err != nil {
		return time.Time{}, fmt.Errorf("failed to update password: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return time.Time{}, fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return time.Time{}, fmt.Errorf("no auth record found to update")
	}

	return now, nil
}

// Initialize 初始化认证信息
func (r *AuthRepository) Initialize(ctx context.Context, authData *auth.Data) error {
	query := `INSERT INTO auth (user_name, password, created_at, updated_at) VALUES (?, ?, ?, ?)`
	now := local.Time()
	_, err := r.db.db.ExecContext(ctx, query, authData.UserName, authData.Password, now, now)
	if err != nil {
		return fmt.Errorf("failed to initialize auth: %w", err)
	}

	return nil
}

// IsInitialized 验证是否已初始化
func (r *AuthRepository) IsInitialized(ctx context.Context) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM auth LIMIT 1)`

	var exists bool
	err := r.db.db.QueryRowContext(ctx, query).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("failed to check auth initialization: %w", err)
	}

	return exists, nil
}
