package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/bestruirui/bestsub/internal/database/interfaces"
	"github.com/bestruirui/bestsub/internal/database/sqlite/database"
	"github.com/bestruirui/bestsub/internal/models/auth"
	"github.com/bestruirui/bestsub/internal/utils/local"
	"github.com/bestruirui/bestsub/internal/utils/passwd"
)

// AuthRepository 认证数据访问实现
type AuthRepository struct {
	db *database.Database
}

// NewAuthRepository 创建认证仓库
func newAuthRepository(db *database.Database) interfaces.AuthRepository {
	return &AuthRepository{db: db}
}

// Get 获取认证信息
func (r *AuthRepository) Get(ctx context.Context) (*auth.Data, error) {
	query := `SELECT id, user_name, password, created_at, updated_at FROM auth LIMIT 1`

	var authData auth.Data
	err := r.db.QueryRowContext(ctx, query).Scan(
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

// Update 更新认证信息
func (r *AuthRepository) Update(ctx context.Context, authData *auth.Data) error {

	hashedPassword, err := passwd.Hash(authData.Password)
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}

	query := `UPDATE auth SET password = ?, updated_at = ?`

	result, err := r.db.ExecContext(ctx, query, hashedPassword, local.Time())
	if err != nil {
		return fmt.Errorf("failed to update auth: %w", err)
	}

	// 检查是否有记录被更新
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("用户不存在")
	}

	return nil
}

// UpdateUsername 更新用户名
func (r *AuthRepository) UpdateUsername(ctx context.Context, username string) error {
	if username == "" {
		return fmt.Errorf("用户名不能为空")
	}

	// 单用户系统，直接更新唯一记录
	query := `UPDATE auth SET user_name = ?, updated_at = ?`

	result, err := r.db.ExecContext(ctx, query, username, local.Time())
	if err != nil {
		return fmt.Errorf("failed to update username: %w", err)
	}

	// 检查是否有记录被更新
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("用户不存在")
	}

	return nil
}

// Initialize 初始化认证信息
func (r *AuthRepository) Initialize(ctx context.Context, authData *auth.Data) error {

	hashedPassword, err := passwd.Hash(authData.Password)
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}

	query := `INSERT INTO auth (user_name, password, created_at, updated_at) VALUES (?, ?, ?, ?)`

	now := local.Time()
	_, err = r.db.ExecContext(ctx, query, authData.UserName, hashedPassword, now, now)
	if err != nil {
		return fmt.Errorf("failed to initialize auth: %w", err)
	}

	return nil
}

// IsInitialized 验证是否已初始化
func (r *AuthRepository) IsInitialized(ctx context.Context) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM auth LIMIT 1)`

	var exists bool
	err := r.db.QueryRowContext(ctx, query).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("failed to check auth initialization: %w", err)
	}

	return exists, nil
}
