package auth

import "time"

// 用户认证信息
type Data struct {
	ID        int64     `db:"id" json:"id"`               // 主键ID
	UserName  string    `db:"user_name" json:"user_name"` // 用户名
	Password  string    `db:"password" json:"-"`          // 密码加密存储（bcrypt），不在JSON中返回
	CreatedAt time.Time `db:"created_at" json:"created_at"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
}

// 认证会话
type Session struct {
	ID           int64     `db:"id" json:"id"`
	IsActive     bool      `db:"is_active" json:"is_active"`
	IPAddress    string    `db:"ip_address" json:"ip_address"`
	UserAgent    string    `db:"user_agent" json:"user_agent"`
	ExpiresAt    time.Time `db:"expires_at" json:"expires_at"`
	TokenHash    string    `db:"token_hash" json:"-"`    // JWT Token的哈希值，不在JSON中返回
	RefreshToken string    `db:"refresh_token" json:"-"` // 刷新Token，不在JSON中返回
	CreatedAt    time.Time `db:"created_at" json:"created_at"`
	UpdatedAt    time.Time `db:"updated_at" json:"updated_at"`
}

// LoginRequest 登录请求模型
type LoginRequest struct {
	Username string `json:"username" binding:"required" example:"admin"` // 用户名
	Password string `json:"password" binding:"required" example:"admin"` // 密码
}

// LoginResponse 登录响应模型
type LoginResponse struct {
	AccessToken  string    `json:"access_token" example:"access_token_string"`   // JWT访问令牌
	RefreshToken string    `json:"refresh_token" example:"refresh_token_string"` // 刷新令牌
	ExpiresAt    time.Time `json:"expires_at" example:"2024-01-01T12:00:00Z"`    // 令牌过期时间
	User         Data      `json:"user"`                                         // 用户信息
}

// ChangePasswordRequest 修改密码请求模型
type ChangePasswordRequest struct {
	OldPassword string `json:"old_password" binding:"required" example:"old_password"` // 旧密码
	NewPassword string `json:"new_password" binding:"required" example:"new_password"` // 新密码
}

// UpdateUserInfoRequest 更新用户信息请求模型
type UpdateUserInfoRequest struct {
	Username string `json:"username" binding:"required" example:"admin"` // 新用户名
}

// SessionListResponse 会话列表响应模型
type SessionListResponse struct {
	Sessions []Session `json:"sessions"` // 会话列表
	Total    int       `json:"total"`    // 总数
}

// RefreshTokenRequest 刷新令牌请求模型
type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required" example:"refresh_token_string"` // 刷新令牌
}

// RefreshTokenResponse 刷新令牌响应模型
type RefreshTokenResponse struct {
	AccessToken  string    `json:"access_token" example:"new_access_token_string"`   // 新的JWT访问令牌
	RefreshToken string    `json:"refresh_token" example:"new_refresh_token_string"` // 新的刷新令牌
	ExpiresAt    time.Time `json:"expires_at" example:"2024-01-01T12:00:00Z"`        // 新令牌过期时间
}
