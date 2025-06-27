package sublink

import (
	"time"

	"github.com/bestruirui/bestsub/internal/models/detector"
	"github.com/bestruirui/bestsub/internal/models/notify"
	"github.com/bestruirui/bestsub/internal/models/parser"
)

// BaseData 订阅链接数据模型
type BaseData struct {
	Name      string                    `db:"name" json:"name" example:"测试订阅"`                        // 链接名称
	URL       string                    `db:"url" json:"url" example:"https://example.com/subscribe"` // 订阅链接地址
	Type      parser.ParserType         `db:"type" json:"type" example:"mihomo"`                      // 链接类型
	UserAgent string                    `db:"user_agent" json:"user_agent" example:"Mozilla/5.0"`     // 自定义 User-Agent
	IsEnabled bool                      `db:"is_enabled" json:"is_enabled" example:"true"`            // 启用/禁用状态
	UseProxy  bool                      `db:"use_proxy" json:"use_proxy" example:"false"`             // 是否启用代理更新订阅
	Detector  []detector.DetectorConfig `db:"detector" json:"detector"`                               // 检测器配置
	Notify    []*notify.NotifyConfig    `db:"notify" json:"notify"`                                   // 通知配置
	CronExpr  string                    `db:"cron_expr" json:"cron_expr" example:"0 */6 * * *"`       // 定时更新设置（cron表达式）
}

type Data struct {
	ID int64 `db:"id" json:"id" example:"1"`
	BaseData
	LastUpdate time.Time `db:"last_update" json:"last_update" example:"2024-01-01T12:00:00Z"` // 最后更新时间
	LastStatus string    `db:"last_status" json:"last_status" example:"success"`              // 最后更新状态
	ErrorMsg   string    `db:"error_msg" json:"error_msg" example:""`                         // 错误信息
	CreatedAt  time.Time `db:"created_at" json:"created_at" example:"2024-01-01T12:00:00Z"`   // 创建时间
	UpdatedAt  time.Time `db:"updated_at" json:"updated_at" example:"2024-01-01T12:00:00Z"`   // 更新时间
}
