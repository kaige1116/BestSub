package database

import (
	"context"

	"github.com/bestruirui/bestsub/internal/database/interfaces"
	authModel "github.com/bestruirui/bestsub/internal/models/auth"
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
		authData := authModel.Default()
		hashedBytes, err := bcrypt.GenerateFromPassword([]byte(authData.Password), bcrypt.DefaultCost)
		if err != nil {
			log.Fatalf("failed to hash password: %v", err)
		}
		authData.Password = string(hashedBytes)
		if err := auth.Initialize(ctx, &authData); err != nil {
			log.Fatalf("failed to initialize auth: %v", err)
		}
		log.Info("初始化默认管理员账号 用户名: admin 密码: admin")
	}
	return nil
}
func initSystemConfig(ctx context.Context, systemConfig interfaces.ConfigRepository) error {
	defaultSystemConfig := system.DefaultDbConfig()
	existingSystemConfig, err := systemConfig.GetAll(ctx)
	notExistConfig := make([]system.Data, 0)
	if err != nil {
		log.Fatalf("failed to get existing system config: %v", err)
	}

	existingSystemConfigMap := make(map[string]bool)
	for _, config := range *existingSystemConfig {
		existingSystemConfigMap[config.Key] = true
	}

	for _, group := range defaultSystemConfig {
		for _, config := range group.Data {
			if !existingSystemConfigMap[config.Key] {
				notExistConfig = append(notExistConfig, config)
			}
		}
	}

	if len(notExistConfig) > 0 {
		if err := systemConfig.Create(ctx, &notExistConfig); err != nil {
			log.Fatalf("failed to create missing system config: %v", err)
		}
	}

	return nil
}
func initTask(ctx context.Context, taskRepo interfaces.TaskRepository) error {
	defaultTasks := task.Default()

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
