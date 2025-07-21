package notify

type Data struct {
	ID     uint16 `db:"id" json:"id"`
	Type   string `db:"type" json:"type"`
	Config string `db:"config" json:"config"`
}

type CreateRequest struct {
	Type   string `json:"type" binding:"required"`                                                                                                                                                                                           // 通知类型
	Config string `json:"config" binding:"required" example:"{\"server\":\"smtp.example.com\",\"port\":587,\"username\":\"test@example.com\",\"password\":\"test\",\"from\":\"test@example.com\",\"tls\":true,\"to\":\"test@example.com\"}"` // 通知配置
}

type Template struct {
	ID       uint16 `db:"id" json:"id"`
	Name     string `db:"name" json:"name"`           // 模板名称
	Template string `db:"templates" json:"templates"` // 模板内容
}

type TemplateCreateRequest struct {
	Name     string `json:"name" binding:"required"`      // 模板名称
	Template string `json:"templates" binding:"required"` // 模板内容
}
