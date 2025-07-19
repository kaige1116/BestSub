package sub

import "time"

// 订阅输出系统的分享链接功能
type Share struct {
	ID             uint16    `db:"id" json:"id"`
	Enable         bool      `db:"enable" json:"enable"`                     // 是否启用
	Name           string    `db:"name" json:"name"`                         // 分享链接名称
	Rename         string    `db:"rename" json:"rename"`                     // 保存文件名
	AccessCount    int       `db:"access_count" json:"access_count"`         // 已访问次数
	MaxAccessCount int       `db:"max_access_count" json:"max_access_count"` // 最大访问次数（0表示无限制）
	Token          string    `db:"token" json:"token"`                       // 随机复杂字符串（用于URL）
	Expires        time.Time `db:"expires" json:"expires"`                   // 过期时间
	Description    string    `db:"description" json:"description"`           // 配置描述
	CreatedAt      time.Time `db:"created_at" json:"created_at"`
	UpdatedAt      time.Time `db:"updated_at" json:"updated_at"`
}
