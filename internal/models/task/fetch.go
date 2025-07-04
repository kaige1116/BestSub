package task

import "github.com/bestruirui/bestsub/internal/models/parser"

type FetchConfig struct {
	Type        parser.ParserType `json:"type"`
	UserAgent   string            `json:"user_agent"`
	URL         string            `json:"url"`
	Retries     int               `json:"retries"`
	Timeout     int               `json:"timeout"`
	ProxyEnable bool              `json:"proxy_enable"`
}
