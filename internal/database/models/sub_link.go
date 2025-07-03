package models

import (
	"time"
)

// 订阅链接
type SubLink struct {
	ID          int64     `db:"id" json:"id"`
	Enable      bool      `db:"enable" json:"enable"`
	Name        string    `db:"name" json:"name"`
	URL         string    `db:"url" json:"url"`
	Description string    `db:"description" json:"description"` // 配置描述
	CreatedAt   time.Time `db:"created_at" json:"created_at"`
	UpdatedAt   time.Time `db:"updated_at" json:"updated_at"`
}
