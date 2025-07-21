package notify

import "bytes"

type Data struct {
	ID     uint16 `db:"id" json:"id"`
	Name   string `db:"name" json:"name"`
	Type   string `db:"type" json:"type"`
	Config string `db:"config" json:"config"`
}

type CreateRequest struct {
	Name   string `json:"name" binding:"required"`                                                                                                                                                                                           // 通知名称
	Type   string `json:"type" binding:"required"`                                                                                                                                                                                           // 通知类型
	Config string `json:"config" binding:"required" example:"{\"server\":\"smtp.example.com\",\"port\":587,\"username\":\"test@example.com\",\"password\":\"test\",\"from\":\"test@example.com\",\"tls\":true,\"to\":\"test@example.com\"}"` // 通知配置
}

type Template struct {
	Type     string `db:"type" json:"type"`           // 模板类型
	Template string `db:"templates" json:"templates"` // 模板内容
}

type Instance interface {
	Init() error
	Send(title string, body *bytes.Buffer) error
}

const (
	TypeLoginSuccess uint16 = 1 << 0 // 登录成功通知
	TypeLoginFailed  uint16 = 1 << 1 // 登录失败通知
)

var TypeMap = map[uint16]string{
	TypeLoginSuccess: "login_success",
	TypeLoginFailed:  "login_failed",
}
