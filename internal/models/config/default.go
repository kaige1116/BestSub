package config

var defaultBase = Base{
	Server: ServerConfig{
		Port: 8080,
		Host: "0.0.0.0",
	},
	Database: DatabaseConfig{
		Type: "sqlite",
	},
	Log: LogConfig{
		Level:  "debug",
		Output: "console",
	},
}



func DefaultBase() Base {
	return defaultBase
}
