package system

type ProxyConfig struct {
	Enable   bool   `json:"enable"`
	Type     string `json:"type"`
	Host     string `json:"host"`
	Port     int    `json:"port"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type TaskConfig struct {
	MaxTimeout int `json:"max_timeout"`
	MaxRetry   int `json:"max_retry"`
}
