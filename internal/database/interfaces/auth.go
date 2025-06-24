package interfaces

import (
	"context"

	"github.com/bestruirui/bestsub/internal/database/models"
)

// 单用户认证数据访问接口
type AuthRepository interface {
	// 获取认证信息
	Get(ctx context.Context) (*models.Auth, error)

	// 更新认证信息
	Update(ctx context.Context, auth *models.Auth) error

	// 初始化认证信息（首次创建密码）
	Initialize(ctx context.Context, auth *models.Auth) error

	// 验证是否已初始化
	IsInitialized(ctx context.Context) (bool, error)

	// 验证密码
	VerifyPassword(ctx context.Context, username, password string) error
}

// 会话数据访问接口
type SessionRepository interface {
	// 创建会话
	Create(ctx context.Context, session *models.Session) error

	// 根据ID获取会话
	GetByID(ctx context.Context, id int64) (*models.Session, error)

	// 根据Token哈希获取会话
	GetByTokenHash(ctx context.Context, tokenHash string) (*models.Session, error)

	// 根据刷新Token获取会话
	GetByRefreshToken(ctx context.Context, refreshToken string) (*models.Session, error)

	// 更新会话
	Update(ctx context.Context, session *models.Session) error

	// 删除会话
	Delete(ctx context.Context, id int64) error

	// 删除所有会话
	DeleteAll(ctx context.Context) error

	// 删除过期会话
	DeleteExpired(ctx context.Context) error

	// 获取所有活跃会话
	GetAllActive(ctx context.Context) ([]*models.Session, error)

	// 停用所有会话
	DeactivateAll(ctx context.Context) error
}
