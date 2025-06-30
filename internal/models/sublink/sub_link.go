package sublink

import (
	"time"

	"github.com/bestruirui/bestsub/internal/models/detector"
	"github.com/bestruirui/bestsub/internal/models/notify"
	"github.com/bestruirui/bestsub/internal/models/parser"
)

type FetchConfig struct {
	URL         string            `db:"url" json:"url" example:"https://example.com/subscribe"`                 // 订阅链接地址
	Type        parser.ParserType `db:"type" json:"type" example:"mihomo"`                                      // 链接类型
	UserAgent   string            `db:"user_agent" json:"user_agent" default:"clash.meta" example:"clash.meta"` // 自定义 User-Agent
	ProxyEnable bool              `db:"proxy_enable" json:"proxy_enable" example:"false"`                       // 是否启用代理更新订阅
	Timeout     int               `db:"timeout" json:"timeout" default:"5" example:"5"`                         // 超时时间（秒）
	Retries     int               `db:"retries" json:"retries" default:"3" example:"3"`                         // 重试次数
}

type FetchResult struct {
	StatusCode int    `json:"status_code"`              // HTTP状态码
	SubType    string `json:"sub_type"`                 // 订阅类型
	NodeCount  int    `json:"node_count"`               // 节点数量
	Size       int64  `json:"size"`                     // 内容大小
	Duration   string `json:"duration" example:"100ms"` // 请求耗时
}

type BaseData struct {
	Name        string                    `db:"name" json:"name" example:"测试订阅"` // 链接名称
	FetchConfig FetchConfig               `db:"fetch_config" json:"fetch_config"`
	IsEnabled   bool                      `db:"is_enabled" json:"is_enabled" example:"true"`      // 启用/禁用状态
	Detector    []detector.DetectorConfig `db:"detector" json:"detector"`                         // 检测器配置
	Notify      []*notify.NotifyConfig    `db:"notify" json:"notify"`                             // 通知配置
	CronExpr    string                    `db:"cron_expr" json:"cron_expr" example:"0 */6 * * *"` // 定时更新设置（cron表达式）
}

type Data struct {
	ID int64 `db:"id" json:"id" example:"1"`
	BaseData
	LastStatus string    `db:"last_status" json:"last_status" example:"success"`            // 最后更新状态
	ErrorMsg   string    `db:"error_msg" json:"error_msg" example:""`                       // 错误信息
	CreatedAt  time.Time `db:"created_at" json:"created_at" example:"2024-01-01T12:00:00Z"` // 创建时间
	UpdatedAt  time.Time `db:"updated_at" json:"updated_at" example:"2024-01-01T12:00:00Z"` // 更新时间
}
