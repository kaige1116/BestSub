package system

type ProxyConfig struct {
	Enable   bool
	Type     string
	Host     string
	Port     int
	Username string
	Password string
}
