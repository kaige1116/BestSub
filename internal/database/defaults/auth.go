package defaults

import "github.com/bestruirui/bestsub/internal/models/auth"

// Auth 获取默认的认证配置
func Auth() *auth.Data {
	return &auth.Data{
		UserName: "admin",
		Password: "admin",
	}
}
