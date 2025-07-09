package config

import (
	"context"
	"strconv"

	"github.com/bestruirui/bestsub/internal/database"
	"github.com/bestruirui/bestsub/internal/models/system"
)

func GetTaskConfig() *system.TaskConfig {
	ctx := context.Background()
	configRepo := database.SystemConfig()

	// 使用 GetConfigsByGroup 批量获取 task 组的所有配置
	configs, err := configRepo.GetConfigsByGroup(ctx, "task")
	if err != nil {
		return nil
	}

	config := &system.TaskConfig{}

	// 遍历配置项，根据键设置对应的字段值
	for _, cfg := range configs {
		switch cfg.Key {
		case "task.max_timeout":
			if cfg.Value != "" {
				maxTimeout, err := strconv.Atoi(cfg.Value)
				if err != nil {
					return nil
				}
				config.MaxTimeout = maxTimeout
			}
		case "task.max_retry":
			if cfg.Value != "" {
				maxRetry, err := strconv.Atoi(cfg.Value)
				if err != nil {
					return nil
				}
				config.MaxRetry = maxRetry
			}
		}
	}

	return config
}
