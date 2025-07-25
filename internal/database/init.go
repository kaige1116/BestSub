package database

import (
	"context"

	"github.com/bestruirui/bestsub/internal/database/interfaces"
	authModel "github.com/bestruirui/bestsub/internal/models/auth"
	"github.com/bestruirui/bestsub/internal/models/config"
	"github.com/bestruirui/bestsub/internal/models/notify"
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
	defaultSystemConfig := config.DefaultAdvance()
	existingSystemConfig, err := systemConfig.GetAll(ctx)
	notExistConfig := make([]config.Advance, 0)
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
func initNotifyTemplate(ctx context.Context, notifyTemplateRepo interfaces.NotifyTemplateRepository) error {
	defaultNotifyTemplates := notify.DefaultTemplates()
	existingNotifyTemplates, err := notifyTemplateRepo.List(ctx)
	if err != nil {
		log.Fatalf("failed to get existing notify templates: %v", err)
	}
	existingNotifyTemplatesMap := make(map[string]bool)
	for _, template := range *existingNotifyTemplates {
		existingNotifyTemplatesMap[template.Type] = true
	}
	for _, template := range defaultNotifyTemplates {
		if !existingNotifyTemplatesMap[template.Type] {
			if err := notifyTemplateRepo.Create(ctx, &template); err != nil {
				log.Fatalf("failed to create missing notify template %s: %v", template.Type, err)
			}
		}
	}
	return nil
}
