package config

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
)

var config Config

const DefaultConfigPath = "./data/config/config.json"

// 从指定路径加载配置文件
func Initialize(configPath string) error {

	if configPath == "" {
		configPath = DefaultConfigPath
	}

	// 尝试从配置文件加载
	if err := loadFromFile(&config, configPath); err != nil {
		// 如果配置文件不存在，创建默认配置文件
		if os.IsNotExist(err) {
			if err := CreateDefaultConfig(configPath); err != nil {
				return fmt.Errorf("创建默认配置文件失败: %v", err)
			}
			// 重新加载配置文件
			if err := loadFromFile(&config, configPath); err != nil {
				return fmt.Errorf("加载默认配置文件失败: %v", err)
			}
		} else {
			return fmt.Errorf("加载配置文件失败: %v", err)
		}
	}

	// 从环境变量覆盖配置
	loadFromEnv(&config)

	// 验证配置
	if err := ValidateConfig(&config); err != nil {
		return fmt.Errorf("配置验证失败: %v", err)
	}

	return nil
}

func Get() Config {
	return config
}

// 从文件加载配置
func loadFromFile(config *Config, filePath string) error {
	// 检查文件是否存在
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return err
	}

	// 读取文件内容
	data, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("读取配置文件失败: %v", err)
	}

	// 解析JSON
	if err := json.Unmarshal(data, config); err != nil {
		return fmt.Errorf("解析配置文件失败: %v", err)
	}

	return nil
}

// 从环境变量加载配置
func loadFromEnv(config *Config) {
	// 服务器配置
	if port := os.Getenv("BESTSUB_SERVER_PORT"); port != "" {
		if p, err := parsePort(port); err == nil {
			config.Server.Port = p
		}
	}
	if host := os.Getenv("BESTSUB_SERVER_HOST"); host != "" {
		config.Server.Host = host
	}

	// 数据库配置
	if dbPath := os.Getenv("BESTSUB_DATABASE_PATH"); dbPath != "" {
		config.Database.Path = dbPath
	}
	if dbType := os.Getenv("BESTSUB_DATABASE_TYPE"); dbType != "" {
		config.Database.Type = dbType
	}

	// 日志配置
	if logLevel := os.Getenv("BESTSUB_LOG_LEVEL"); logLevel != "" {
		config.Log.Level = logLevel
	}
	if logOutput := os.Getenv("BESTSUB_LOG_OUTPUT"); logOutput != "" {
		config.Log.Output = logOutput
	}
	if logDir := os.Getenv("BESTSUB_LOG_DIR"); logDir != "" {
		config.Log.Dir = logDir
	}

	// JWT配置
	if jwtSecret := os.Getenv("BESTSUB_JWT_SECRET"); jwtSecret != "" {
		config.JWT.Secret = jwtSecret
	}
	if jwtExpires := os.Getenv("BESTSUB_JWT_EXPIRES_IN"); jwtExpires != "" {
		if exp, err := parseInt(jwtExpires); err == nil {
			config.JWT.ExpiresIn = exp
		}
	}
}

// 解析端口号
func parsePort(portStr string) (int, error) {
	port, err := strconv.Atoi(portStr)
	if err != nil {
		return 0, fmt.Errorf("无效的端口号: %s", portStr)
	}
	if port <= 0 || port > 65535 {
		return 0, fmt.Errorf("端口号超出范围: %d", port)
	}
	return port, nil
}

// parseInt 解析整数
func parseInt(intStr string) (int, error) {
	value, err := strconv.Atoi(intStr)
	if err != nil {
		return 0, fmt.Errorf("无效的整数: %s", intStr)
	}
	return value, nil
}
