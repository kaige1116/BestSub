package notify

func DefaultTemplates() []Template {
	return []Template{
		{"login_success", "登录成功", "{{.Username}}{{.Time}}{{.IP}}{{.UserAgent}}"},
		{"login_failed", "登录失败", "{{.Username}}{{.Time}}{{.IP}}{{.UserAgent}}"},
	}
}
