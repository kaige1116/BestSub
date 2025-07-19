package auth

import "time"

// 用户认证信息
type Data struct {
	ID       uint8  `db:"id" json:"-"`              // 主键ID
	UserName string `db:"username" json:"username"` // 用户名
	Password string `db:"password" json:"-"`
}

// 认证会话
type Session struct {
	IsActive     bool   `json:"is_active"`
	ClientIP     uint32 `json:"client_ip"`
	UserAgent    string `json:"user_agent"`
	ExpiresAt    uint32 `json:"expires_at"`
	CreatedAt    uint32 `json:"created_at"`
	LastAccessAt uint32 `json:"last_access_at"`
	HashRToken   uint64 `json:"-"`
	HashAToken   uint64 `json:"-"`
}

type SessionResponse struct {
	ID           uint8     `json:"id"`
	IsActive     bool      `json:"is_active"`
	ClientIP     string    `json:"client_ip"`
	UserAgent    string    `json:"user_agent"`
	ExpiresAt    time.Time `json:"expires_at"`
	CreatedAt    time.Time `json:"created_at"`
	LastAccessAt time.Time `json:"last_access_at"`
}

// LoginRequest 登录请求模型
type LoginRequest struct {
	Username string `json:"username" binding:"required" example:"admin"` // 用户名
	Password string `json:"password" binding:"required" example:"admin"` // 密码
}

// LoginResponse 登录响应模型
type LoginResponse struct {
	AccessToken      string    `json:"access_token" example:"access_token_string"`        // JWT访问令牌
	RefreshToken     string    `json:"refresh_token" example:"refresh_token_string"`      // 刷新令牌
	AccessExpiresAt  time.Time `json:"access_expires_at" example:"2024-01-01T12:00:00Z"`  // 令牌过期时间
	RefreshExpiresAt time.Time `json:"refresh_expires_at" example:"2024-01-01T12:00:00Z"` // 刷新令牌过期时间
}

// ChangePasswordRequest 修改密码请求模型
type ChangePasswordRequest struct {
	Username    string `json:"username" binding:"required" example:"admin"`            // 用户名
	OldPassword string `json:"old_password" binding:"required" example:"old_password"` // 旧密码
	NewPassword string `json:"new_password" binding:"required" example:"new_password"` // 新密码
}

// UpdateUserInfoRequest 更新用户信息请求模型
type UpdateUserInfoRequest struct {
	Username string `json:"username" binding:"required" example:"admin"` // 新用户名
}

// SessionListResponse 会话列表响应模型
type SessionListResponse struct {
	Sessions []SessionResponse `json:"sessions"` // 会话列表
	Total    uint8             `json:"total"`    // 总数
}

// RefreshTokenRequest 刷新令牌请求模型
type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required" example:"refresh_token_string"` // 刷新令牌
}
