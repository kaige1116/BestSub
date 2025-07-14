package database

import (
	"context"
	"fmt"
	"sync"

	"github.com/bestruirui/bestsub/internal/database/interfaces"
	"github.com/bestruirui/bestsub/internal/utils/log"
)

var manager *Manager

// 数据库管理器
type Manager struct {
	repository *interfaces.RepositoryManager
	sqltype    string
	path       string
	mu         sync.RWMutex
	closed     bool
}

// 初始化数据库管理器
func Initialize(sqltype, path string) error {
	var err error
	manager = &Manager{
		sqltype: sqltype,
		path:    path,
	}
	err = manager.init()
	return err
}

// 获取仓库实例
func GetRepository() *interfaces.RepositoryManager {
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
	// 这里我们需要通过其他方式来检查连接，因为现在是通过 RepositoryManager
	// 暂时返回 nil，后续可以添加专门的健康检查方法
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

	log.Debugf("初始化数据库: 类型 %s, 路径 %s", m.sqltype, m.path)

	repository, err := Init(context.Background(), m.sqltype, m.path)
	if err != nil {
		return fmt.Errorf("failed to initialize database: %w", err)
	}
	m.repository = repository
	return nil
}
