package notify

var templates = []Template{
	{
		Type:     "login_success",
		Template: `{{.Username}}{{.Time}}{{.IP}}{{.UserAgent}}`,
	},
	{
		Type:     "login_failed",
		Template: `{{.Username}}{{.Time}}{{.IP}}{{.UserAgent}}`,
	},
}

func DefaultTemplates() []Template {
	return templates
}
