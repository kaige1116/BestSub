package database

import (
	"context"
	"fmt"
	"sync"

	"github.com/bestruirui/bestsub/internal/config"
	"github.com/bestruirui/bestsub/internal/database/repository"
	"github.com/bestruirui/bestsub/internal/database/repository/sqlite"
	"github.com/bestruirui/bestsub/internal/utils/log"
)

var (
	manager *Manager
	once    sync.Once
)

// 数据库管理器
type Manager struct {
	repository repository.Repository
	config     config.DatabaseConfig
	mu         sync.RWMutex
	closed     bool
}

// 初始化数据库管理器
func Initialize(cfg config.DatabaseConfig) error {
	var err error
	once.Do(func() {
		manager = &Manager{
			config: cfg,
		}
		err = manager.init()
	})
	return err
}

// Reinitialize 重新初始化数据库管理器（用于测试或配置更改）
func Reinitialize(cfg config.DatabaseConfig) error {
	// 关闭现有连接
	if manager != nil {
		Close()
	}

	// 重置 once，允许重新初始化
	once = sync.Once{}

	return Initialize(cfg)
}

// 获取仓库实例
func GetRepository() repository.Repository {
	if manager == nil {
		log.Fatal("database manager not initialized, call Initialize() first")
	}
	manager.mu.RLock()
	defer manager.mu.RUnlock()

	if manager.closed {
		log.Fatal("database manager has been closed")
	}

	return manager.repository
}

// 检查是否已初始化
func IsInitialized() bool {
	return manager != nil && manager.repository != nil
}

// 健康检查
func HealthCheck(ctx context.Context) error {
	if manager == nil {
		return fmt.Errorf("database manager not initialized")
	}

	manager.mu.RLock()
	defer manager.mu.RUnlock()

	if manager.closed {
		return fmt.Errorf("database manager has been closed")
	}

	// 检查数据库连接
	if sqliteRepo, ok := manager.repository.(*sqlite.Repository); ok {
		return sqliteRepo.Database().Ping(ctx)
	}

	return nil
}

// 关闭数据库连接
func Close() error {
	if manager == nil {
		return nil
	}

	manager.mu.Lock()
	defer manager.mu.Unlock()

	if manager.closed {
		return nil
	}

	var err error
	if manager.repository != nil {
		err = manager.repository.Close()
		log.Debug("数据库连接已关闭")
	}

	manager.closed = true
	return err
}

// 初始化数据库连接
func (m *Manager) init() error {
	m.mu.Lock()
	defer m.mu.Unlock()

	log.Debugf("初始化数据库: 类型 %s, 路径 %s", m.config.Type, m.config.Path)

	switch m.config.Type {
	case "sqlite":
		db, err := sqlite.New(m.config.Path)
		if err != nil {
			return fmt.Errorf("failed to create sqlite database: %w", err)
		}
		m.repository = sqlite.NewRepository(db)

		log.Debugf("SQLite 数据库初始化成功")
		return nil
	case "mysql":
		// TODO: 实现 MySQL 支持
		return fmt.Errorf("mysql support not implemented yet")
	default:
		return fmt.Errorf("unsupported database type: %s", m.config.Type)
	}
}
