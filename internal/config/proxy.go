package config

import (
	"context"
	"fmt"
	"strconv"

	"github.com/bestruirui/bestsub/internal/database"
	"github.com/bestruirui/bestsub/internal/models/system"
)

func GetProxyConfig() (*system.Proxy, error) {
	ctx := context.Background()
	configRepo := database.SystemConfig()

	config := &system.Proxy{}

	enableConfig, err := configRepo.GetByKey(ctx, "proxy.enable")
	if err != nil {
		return nil, fmt.Errorf("failed to get proxy.enable: %w", err)
	}
	if enableConfig != nil {
		config.Enable = enableConfig.Value == "true"
	}

	typeConfig, err := configRepo.GetByKey(ctx, "proxy.type")
	if err != nil {
		return nil, fmt.Errorf("failed to get proxy.type: %w", err)
	}
	if typeConfig != nil {
		config.Type = typeConfig.Value
	}

	hostConfig, err := configRepo.GetByKey(ctx, "proxy.host")
	if err != nil {
		return nil, fmt.Errorf("failed to get proxy.host: %w", err)
	}
	if hostConfig != nil {
		config.Host = hostConfig.Value
	}
	portConfig, err := configRepo.GetByKey(ctx, "proxy.port")
	if err != nil {
		return nil, fmt.Errorf("failed to get proxy.port: %w", err)
	}
	if portConfig != nil && portConfig.Value != "" {
		port, err := strconv.Atoi(portConfig.Value)
		if err != nil {
			return nil, fmt.Errorf("invalid proxy port: %w", err)
		}
		config.Port = port
	}

	usernameConfig, err := configRepo.GetByKey(ctx, "proxy.username")
	if err != nil {
		return nil, fmt.Errorf("failed to get proxy.username: %w", err)
	}
	if usernameConfig != nil {
		config.Username = usernameConfig.Value
	}

	passwordConfig, err := configRepo.GetByKey(ctx, "proxy.password")
	if err != nil {
		return nil, fmt.Errorf("failed to get proxy.password: %w", err)
	}
	if passwordConfig != nil {
		config.Password = passwordConfig.Value
	}

	return config, nil
}
