package notify

import (
	"time"
)

// 全局管理通知渠道配置
type Data struct {
	ID          int64     `db:"id" json:"id"`
	Enable      bool      `db:"enable" json:"enable"`           // 是否启用
	Name        string    `db:"name" json:"name"`               // 通知渠道名称
	Type        string    `db:"type" json:"type"`               // 通知类型
	Config      string    `db:"config" json:"config"`           // 通知配置
	TestResult  string    `db:"test_result" json:"test_result"` // 测试连接结果
	LastTest    time.Time `db:"last_test" json:"last_test"`     // 最后测试时间
	Description string    `db:"description" json:"description"` // 配置描述
	CreatedAt   time.Time `db:"created_at" json:"created_at"`
	UpdatedAt   time.Time `db:"updated_at" json:"updated_at"`
}
