package models

import (
	"time"
)

type User struct {
	ID        int64     `db:"id" json:"id"`
	Username  string    `db:"username" json:"username"`
	Password  string    `db:"password" json:"-"` // 密码加密存储（bcrypt），不在JSON中返回
	CreatedAt time.Time `db:"created_at" json:"created_at"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
}

type Session struct {
	ID           int64     `db:"id" json:"id"`
	UserID       int64     `db:"user_id" json:"user_id"`
	TokenHash    string    `db:"token_hash" json:"-"` // JWT Token的哈希值，不在JSON中返回
	ExpiresAt    time.Time `db:"expires_at" json:"expires_at"`
	RefreshToken string    `db:"refresh_token" json:"-"` // 刷新Token，不在JSON中返回
	IPAddress    string    `db:"ip_address" json:"ip_address"`
	UserAgent    string    `db:"user_agent" json:"user_agent"`
	IsActive     bool      `db:"is_active" json:"is_active"`
	CreatedAt    time.Time `db:"created_at" json:"created_at"`
	UpdatedAt    time.Time `db:"updated_at" json:"updated_at"`
}
