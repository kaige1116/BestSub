package auth

import "time"

type Data struct {
	ID       uint8  `db:"id" json:"-"`
	UserName string `db:"username" json:"username"`
	Password string `db:"password" json:"-"`
}

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

type LoginRequest struct {
	Username string `json:"username" binding:"required" example:"admin"`
	Password string `json:"password" binding:"required" example:"admin"`
}

type LoginResponse struct {
	AccessToken      string    `json:"access_token" example:"access_token_string"`
	RefreshToken     string    `json:"refresh_token" example:"refresh_token_string"`
	AccessExpiresAt  time.Time `json:"access_expires_at" example:"2024-01-01T12:00:00Z"`
	RefreshExpiresAt time.Time `json:"refresh_expires_at" example:"2024-01-01T12:00:00Z"`
}

type ChangePasswordRequest struct {
	Username    string `json:"username" binding:"required" example:"admin"`
	OldPassword string `json:"old_password" binding:"required" example:"old_password"`
	NewPassword string `json:"new_password" binding:"required" example:"new_password"`
}

type UpdateUserInfoRequest struct {
	Username string `json:"username" binding:"required" example:"admin"`
}

type SessionListResponse struct {
	Sessions []SessionResponse `json:"sessions"`
	Total    uint8             `json:"total"`
}

type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required" example:"refresh_token_string"`
}
type LoginNotify struct {
	Username  string `json:"username"`
	IP        string `json:"ip"`
	Time      string `json:"time"`
	Msg       string `json:"msg"`
	UserAgent string `json:"user_agent"`
}
