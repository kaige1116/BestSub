package models

import "time"

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
	User         UserInfo  `json:"user"`                                         // 用户信息
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

// UserInfo 用户信息模型
type UserInfo struct {
	Username  string    `json:"username" example:"admin"`                  // 用户名
	CreatedAt time.Time `json:"created_at" example:"2024-01-01T12:00:00Z"` // 创建时间
	UpdatedAt time.Time `json:"updated_at" example:"2024-01-01T12:00:00Z"` // 更新时间
}

// LogoutRequest 登出请求模型
type LogoutRequest struct {
	// 可以为空，通过JWT中间件获取当前用户信息
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

// SessionInfo 会话信息模型
type SessionInfo struct {
	ID        int64     `json:"id" example:"1"`                            // 会话ID
	IPAddress string    `json:"ip_address" example:"192.168.1.1"`          // IP地址
	UserAgent string    `json:"user_agent" example:"Mozilla/5.0..."`       // 用户代理
	IsActive  bool      `json:"is_active" example:"true"`                  // 是否活跃
	ExpiresAt time.Time `json:"expires_at" example:"2024-01-01T12:00:00Z"` // 过期时间
	CreatedAt time.Time `json:"created_at" example:"2024-01-01T12:00:00Z"` // 创建时间
	UpdatedAt time.Time `json:"updated_at" example:"2024-01-01T12:00:00Z"` // 更新时间
}

// SessionListResponse 会话列表响应模型
type SessionListResponse struct {
	Sessions []SessionInfo `json:"sessions"` // 会话列表
	Total    int           `json:"total"`    // 总数
}
