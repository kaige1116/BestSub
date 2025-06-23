package models

import (
	"time"
)

// 用户认证信息
type Auth struct {
	UserName  string    `db:"user_name" json:"user_name"` // 用户名
	Password  string    `db:"password" json:"-"`          // 密码加密存储（bcrypt），不在JSON中返回
	CreatedAt time.Time `db:"created_at" json:"created_at"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
}

// 认证会话
type Session struct {
	ID           int64     `db:"id" json:"id"`
	TokenHash    string    `db:"token_hash" json:"-"` // JWT Token的哈希值，不在JSON中返回
	ExpiresAt    time.Time `db:"expires_at" json:"expires_at"`
	RefreshToken string    `db:"refresh_token" json:"-"` // 刷新Token，不在JSON中返回
	IPAddress    string    `db:"ip_address" json:"ip_address"`
	UserAgent    string    `db:"user_agent" json:"user_agent"`
	IsActive     bool      `db:"is_active" json:"is_active"`
	CreatedAt    time.Time `db:"created_at" json:"created_at"`
	UpdatedAt    time.Time `db:"updated_at" json:"updated_at"`
}
