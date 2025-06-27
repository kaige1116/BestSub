package notify

type NotifyConfig struct {
	NotifyType NotifyType `json:"type"`
	IsEnabled  *bool      `json:"is_enabled"`
	TemplateId int        `json:"template_id"`
	Config     *string    `json:"config"`
}

type EmailConfig struct {
	Host     string `json:"host"`
	Port     int    `json:"port"`
	Username string `json:"username"`
	Password string `json:"password"`
	From     string `json:"from"`
}
