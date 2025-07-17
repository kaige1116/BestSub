package database

import (
	"context"

	"github.com/bestruirui/bestsub/internal/database/defaults"
	"github.com/bestruirui/bestsub/internal/database/interfaces"
	"github.com/bestruirui/bestsub/internal/models/system"
	"github.com/bestruirui/bestsub/internal/models/task"
	"github.com/bestruirui/bestsub/internal/utils/log"
	"golang.org/x/crypto/bcrypt"
)

func initAuth(ctx context.Context, auth interfaces.AuthRepository) error {
	isInitialized, err := auth.IsInitialized(ctx)
	if err != nil {
		log.Fatalf("failed to check if database is initialized: %v", err)
	}
	if !isInitialized {
		authData := defaults.Auth()
		hashedBytes, err := bcrypt.GenerateFromPassword([]byte(authData.Password), bcrypt.DefaultCost)
		if err != nil {
			log.Fatalf("failed to hash password: %v", err)
		}
		authData.Password = string(hashedBytes)
		if err := auth.Initialize(ctx, authData); err != nil {
			log.Fatalf("failed to initialize auth: %v", err)
		}
		log.Info("初始化默认管理员账号 用户名: admin 密码: admin")
	}
	return nil
}
func initSystemConfig(ctx context.Context, systemConfig interfaces.ConfigRepository) error {
	defaultSystemConfig := defaults.Configs()

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
func initTask(ctx context.Context, taskRepo interfaces.TaskRepository) error {
	defaultTasks := defaults.Tasks()

	existingTasks, err := taskRepo.GetSystemTasks(ctx)
	if err != nil {
		log.Fatalf("failed to get existing system tasks: %v", err)
	}

	existingTasksMap := make(map[string]bool)
	for _, task := range *existingTasks {
		existingTasksMap[task.Type] = true
	}

	defaultTasksMap := make(map[string]task.Data)
	for _, task := range defaultTasks {
		defaultTasksMap[task.Type] = task
	}

	for taskType, defaultTask := range defaultTasksMap {
		if !existingTasksMap[taskType] {
			if _, err := taskRepo.Create(ctx, &defaultTask); err != nil {
				log.Fatalf("failed to create missing system task %s: %v", taskType, err)
			}
		}
	}

	for _, existingTask := range *existingTasks {
		if _, exists := defaultTasksMap[existingTask.Type]; !exists {
			if err := taskRepo.Delete(ctx, existingTask.ID); err != nil {
				log.Fatalf("failed to delete extra system task %s: %v", existingTask.Type, err)
			}
		}
	}

	return nil
}
