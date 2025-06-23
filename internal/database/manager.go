package database

import (
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

// Manager 数据库管理器
type Manager struct {
	repository repository.Repository
	mu         sync.RWMutex
}

// Initialize 初始化数据库管理器
func Initialize(config config.DatabaseConfig) error {
	var err error
	once.Do(func() {
		manager = &Manager{}
		err = manager.init(config.Type, config.Path)
	})
	return err
}

// GetRepository 获取仓库实例
func GetRepository() repository.Repository {
	if manager == nil {
		log.Fatal("database manager not initialized, call Initialize() first")
	}
	manager.mu.RLock()
	defer manager.mu.RUnlock()
	return manager.repository
}

// Close 关闭数据库连接
func Close() error {
	if manager == nil {
		return nil
	}
	manager.mu.Lock()
	defer manager.mu.Unlock()

	if manager.repository != nil {
		return manager.repository.Close()
	}
	return nil
}

// init 初始化数据库连接
func (m *Manager) init(dbType, dbPath string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	switch dbType {
	case "sqlite":
		db, err := sqlite.New(dbPath)
		if err != nil {
			return fmt.Errorf("failed to create sqlite database: %w", err)
		}
		m.repository = sqlite.NewRepository(db)
		return nil
	case "mysql":
		// TODO: 实现 MySQL 支持
		return fmt.Errorf("mysql support not implemented yet")
	default:
		return fmt.Errorf("unsupported database type: %s", dbType)
	}
}
