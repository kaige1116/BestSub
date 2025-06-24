package sqlite

import (
	"github.com/bestruirui/bestsub/internal/database/repository"
	"github.com/bestruirui/bestsub/internal/database/repository/interfaces"
)

// Repository SQLite仓库实现
type Repository struct {
	*repository.BaseRepository
	db *Database
}

// 编译时检查是否实现了接口
var _ repository.Repository = (*Repository)(nil)
var _ repository.RepositoryFactory = (*Repository)(nil)

// NewRepository 创建新的SQLite仓库
func NewRepository(db *Database) repository.Repository {
	repo := &Repository{db: db}
	repo.BaseRepository = repository.NewBaseRepository(repo)
	return repo
}

// 实现 RepositoryFactory 接口的工厂方法
func (r *Repository) CreateAuthRepository() interfaces.AuthRepository {
	return NewAuthRepository(r.db)
}

func (r *Repository) CreateSessionRepository() interfaces.SessionRepository {
	return NewSessionRepository(r.db)
}

func (r *Repository) CreateSystemConfigRepository() interfaces.SystemConfigRepository {
	return NewSystemConfigRepository(r.db)
}

func (r *Repository) CreateNotificationChannelRepository() interfaces.NotificationChannelRepository {
	return NewNotificationChannelRepository(r.db)
}

func (r *Repository) CreateTaskRepository() interfaces.TaskRepository {
	return NewTaskRepository(r.db)
}

func (r *Repository) CreateSubStorageConfigRepository() interfaces.SubStorageConfigRepository {
	return NewSubStorageConfigRepository(r.db)
}

func (r *Repository) CreateSubOutputTemplateRepository() interfaces.SubOutputTemplateRepository {
	return NewSubOutputTemplateRepository(r.db)
}

func (r *Repository) CreateSubNodeFilterRuleRepository() interfaces.SubNodeFilterRuleRepository {
	return NewSubNodeFilterRuleRepository(r.db)
}

func (r *Repository) CreateSubLinkRepository() interfaces.SubLinkRepository {
	return NewSubLinkRepository(r.db)
}

func (r *Repository) CreateSubLinkModuleConfigRepository() interfaces.SubLinkModuleConfigRepository {
	return NewSubLinkModuleConfigRepository(r.db)
}

func (r *Repository) CreateSubSaveConfigRepository() interfaces.SubSaveConfigRepository {
	return NewSubSaveConfigRepository(r.db)
}

func (r *Repository) CreateSubShareLinkRepository() interfaces.SubShareLinkRepository {
	return NewSubShareLinkRepository(r.db)
}

// Database 获取数据库连接
func (r *Repository) Database() *Database {
	return r.db
}

// Close 关闭数据库连接
func (r *Repository) Close() error {
	return r.db.Close()
}
