package sub

import "time"

type Share struct {
	ID     uint16 `db:"id" json:"id"`
	Enable bool   `db:"enable" json:"enable"` // 是否启用
	Name   string `db:"name" json:"name"`     // 分享链接名称
	Config string `db:"config" json:"config"` // 以json格式存储
}
type Config struct {
	AccessCount    int       `db:"access_count" json:"access_count"`         // 已访问次数
	MaxAccessCount int       `db:"max_access_count" json:"max_access_count"` // 最大访问次数（0表示无限制）
	Token          string    `db:"token" json:"token"`                       // 随机复杂字符串（用于URL）
	Expires        time.Time `db:"expires" json:"expires"`                   // 过期时间
}
