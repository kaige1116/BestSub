package config

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"strconv"
)

var config Config

const DefaultConfigPath = "./data/config/config.json"

func init() {
	configPath := flag.String("c", DefaultConfigPath, "config file path")
	flag.Parse()
	if *configPath == "" {
		*configPath = DefaultConfigPath
	}

	if err := loadFromFile(&config, *configPath); err != nil {
		if os.IsNotExist(err) {
			if err := CreateDefaultConfig(*configPath); err != nil {
				panic(fmt.Errorf("创建默认配置文件失败: %v", err))
			}
			if err := loadFromFile(&config, *configPath); err != nil {
				panic(fmt.Errorf("加载默认配置文件失败: %v", err))
			}
		} else {
			panic(fmt.Errorf("加载配置文件失败: %v", err))
		}
	}

	loadFromEnv(&config)

	if err := ValidateConfig(&config); err != nil {
		panic(fmt.Errorf("配置验证失败: %v", err))
	}

}

func Get() Config {
	return config
}

func loadFromFile(config *Config, filePath string) error {
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return err
	}
	data, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("读取配置文件失败: %v", err)
	}
	if err := json.Unmarshal(data, config); err != nil {
		return fmt.Errorf("解析配置文件失败: %v", err)
	}

	return nil
}

func loadFromEnv(config *Config) {
	if port := os.Getenv("BESTSUB_SERVER_PORT"); port != "" {
		if p, err := parsePort(port); err == nil {
			config.Server.Port = p
		}
	}
	if host := os.Getenv("BESTSUB_SERVER_HOST"); host != "" {
		config.Server.Host = host
	}
	if dbPath := os.Getenv("BESTSUB_DATABASE_PATH"); dbPath != "" {
		config.Database.Path = dbPath
	}
	if dbType := os.Getenv("BESTSUB_DATABASE_TYPE"); dbType != "" {
		config.Database.Type = dbType
	}
	if logLevel := os.Getenv("BESTSUB_LOG_LEVEL"); logLevel != "" {
		config.Log.Level = logLevel
	}
	if logOutput := os.Getenv("BESTSUB_LOG_OUTPUT"); logOutput != "" {
		config.Log.Output = logOutput
	}
	if logDir := os.Getenv("BESTSUB_LOG_DIR"); logDir != "" {
		config.Log.Dir = logDir
	}
	if jwtSecret := os.Getenv("BESTSUB_JWT_SECRET"); jwtSecret != "" {
		config.JWT.Secret = jwtSecret
	}
}

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

func parseInt(intStr string) (int, error) {
	value, err := strconv.Atoi(intStr)
	if err != nil {
		return 0, fmt.Errorf("无效的整数: %s", intStr)
	}
	return value, nil
}
