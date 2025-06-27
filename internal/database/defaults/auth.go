package defaults

import (
	"github.com/bestruirui/bestsub/internal/database/models"
)

// Auth 获取默认的认证配置
func Auth() *models.Auth {
	return &models.Auth{
		UserName: "admin",
		Password: "admin",
	}
}
