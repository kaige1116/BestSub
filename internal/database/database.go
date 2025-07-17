package database

import (
	"context"

	"github.com/bestruirui/bestsub/internal/database/client/sqlite"
	"github.com/bestruirui/bestsub/internal/database/interfaces"
	"github.com/bestruirui/bestsub/internal/utils/log"
)

var repo interfaces.Repository

func Initialize(sqltype, path string) error {
	var err error
	switch sqltype {
	case "sqlite":
		repo, err = sqlite.New(path)
		if err != nil {
			log.Fatalf("failed to create sqlite database: %v", err)
		}
	default:
		log.Fatalf("unsupported database type: %s", sqltype)
	}
	if err := repo.Migrate(); err != nil {
		log.Fatalf("failed to migrate database: %v", err)
	}
	if err := initAuth(context.Background(), AuthRepo()); err != nil {
		log.Fatalf("failed to initialize auth: %v", err)
	}
	if err := initSystemConfig(context.Background(), ConfigRepo()); err != nil {
		log.Fatalf("failed to initialize system config: %v", err)
	}
	if err := initTask(context.Background(), TaskRepo()); err != nil {
		log.Fatalf("failed to initialize tasks: %v", err)
	}
	return nil
}
func Close() error {
	return repo.Close()
}
