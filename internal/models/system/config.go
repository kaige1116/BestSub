package system

// Config 全局配置
type Config struct {
	Server   ServerConfig   `json:"server"`
	Database DatabaseConfig `json:"database"`
	Log      LogConfig      `json:"log"`
	JWT      JWTConfig      `json:"jwt"`
	Session  SessionConfig  `json:"-"`
}

// ServerConfig 服务器配置
type ServerConfig struct {
	Port int    `json:"port"`
	Host string `json:"host"`
}

// DatabaseConfig 数据库配置
type DatabaseConfig struct {
	Type string `json:"type"`
	Path string `json:"-"`
}

// LogConfig 日志配置
type LogConfig struct {
	Level  string `json:"level"`
	Output string `json:"output"`
	Path   string `json:"-"`
}

// JWTConfig JWT配置
type JWTConfig struct {
	Secret string `json:"secret"`
}

// SessionConfig 会话配置
type SessionConfig struct {
	Path string `json:"-"`
}
