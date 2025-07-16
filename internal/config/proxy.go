package config

import (
	"context"
	"fmt"
	"strconv"

	"github.com/bestruirui/bestsub/internal/database"
	"github.com/bestruirui/bestsub/internal/models/system"
)

func Proxy() (*system.ProxyConfig, error) {
	ctx := context.Background()
	configRepo := database.SystemConfig()

	// 使用 GetConfigsByGroup 批量获取 proxy 组的所有配置
	configs, err := configRepo.GetConfigsByGroup(ctx, "proxy")
	if err != nil {
		return nil, fmt.Errorf("failed to get proxy configs: %w", err)
	}

	config := &system.ProxyConfig{}

	// 遍历配置项，根据键设置对应的字段值
	for _, cfg := range configs {
		switch cfg.Key {
		case "proxy.enable":
			config.Enable = cfg.Value == "true"
		case "proxy.type":
			config.Type = cfg.Value
		case "proxy.host":
			config.Host = cfg.Value
		case "proxy.port":
			if cfg.Value != "" {
				port, err := strconv.Atoi(cfg.Value)
				if err != nil {
					return nil, fmt.Errorf("invalid proxy port: %w", err)
				}
				config.Port = port
			}
		case "proxy.username":
			config.Username = cfg.Value
		case "proxy.password":
			config.Password = cfg.Value
		}
	}

	return config, nil
}
