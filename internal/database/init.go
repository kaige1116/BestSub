package database

import (
	"context"

	"github.com/bestruirui/bestsub/internal/database/defaults"
	"github.com/bestruirui/bestsub/internal/database/interfaces"
	"github.com/bestruirui/bestsub/internal/database/sqlite"
	"github.com/bestruirui/bestsub/internal/models/system"
	"github.com/bestruirui/bestsub/internal/utils/log"
)

func Init(ctx context.Context, sqltype, path string) (repository *interfaces.RepositoryManager, err error) {
	var repo interfaces.Repository
	var migrator interfaces.Migrator

	switch sqltype {
	case "sqlite":
		db, err := sqlite.New(path)
		if err != nil {
			log.Fatalf("failed to create sqlite database: %v", err)
		}
		repo = sqlite.NewRepo(db)
		migrator = sqlite.NewMigrator(db)
	default:
		log.Fatalf("unsupported database type: %s", sqltype)
	}
	repository = interfaces.NewRepositoryManager(repo)
	log.Debugf("数据库初始化开始")
	if err := migrator.Apply(ctx); err != nil {
		log.Fatalf("failed to apply migrations: %v", err)
	}
	log.Debugf("数据库迁移成功")
	if err := initAuth(ctx, repository); err != nil {
		log.Fatalf("failed to initialize auth: %v", err)
	}
	log.Debugf("数据库初始化成功")
	if err := initSystemConfig(ctx, repository); err != nil {
		log.Fatalf("failed to initialize system config: %v", err)
	}
	log.Debugf("系统配置初始化成功")
	return repository, nil
}

func initAuth(ctx context.Context, repo *interfaces.RepositoryManager) error {
	auth := repo.Auth()
	isInitialized, err := auth.IsInitialized(ctx)
	if err != nil {
		log.Fatalf("failed to check if database is initialized: %v", err)
	}
	if !isInitialized {
		if err := auth.Initialize(ctx, defaults.Auth()); err != nil {
			log.Fatalf("failed to initialize auth: %v", err)
		}
		log.Info("初始化默认管理员账号 用户名: admin 密码: admin")
	}
	return nil
}
func initSystemConfig(ctx context.Context, repo *interfaces.RepositoryManager) error {
	defaultSystemConfig := defaults.Configs()
	systemConfig := repo.SystemConfig()

	existingKeys, err := systemConfig.GetAllKeys(ctx)
	if err != nil {
		log.Fatalf("failed to get existing config keys: %v", err)
	}

	existingKeysMap := make(map[string]bool)
	for _, key := range existingKeys {
		existingKeysMap[key] = true
	}

	defaultKeysMap := make(map[string]system.Data)
	for _, config := range defaultSystemConfig {
		defaultKeysMap[config.Key] = config
	}

	for key, config := range defaultKeysMap {
		if !existingKeysMap[key] {
			if err := systemConfig.Create(ctx, &config); err != nil {
				log.Fatalf("failed to create missing config %s: %v", key, err)
			}
		}
	}

	for _, key := range existingKeys {
		if _, exists := defaultKeysMap[key]; !exists {
			if err := systemConfig.DeleteByKey(ctx, key); err != nil {
				log.Fatalf("failed to delete extra config %s: %v", key, err)
			}
		}
	}

	return nil
}
