package sub

import (
	"time"

	taskModel "github.com/bestruirui/bestsub/internal/models/task"
	execer "github.com/bestruirui/bestsub/internal/modules/exec/executor"
)

type Data struct {
	ID        uint16    `db:"id" json:"id"`
	Enable    bool      `db:"enable" json:"enable"`
	CronExpr  string    `db:"cron_expr" json:"cron_expr"`
	Name      string    `db:"name" json:"name"`
	Config    string    `db:"config" json:"config"`
	Result    string    `db:"result" json:"result"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
}

type CreateRequest struct {
	Enable   bool         `json:"enable"`
	CronExpr string       `json:"cron_expr" example:"0 0 * * *" description:"cron表达式"`
	Name     string       `json:"name" example:"测试订阅任务" description:"订阅任务名称"`
	Config   execer.Fetch `json:"config" description:"订阅器配置"`
}

type UpdateRequest struct {
	ID       uint16       `json:"id"`
	Enable   bool         `json:"enable"`
	CronExpr string       `json:"cron_expr" example:"0 0 * * *" description:"cron表达式"`
	Name     string       `json:"name" example:"测试订阅任务" description:"订阅任务名称"`
	Config   execer.Fetch `json:"config" description:"订阅器配置"`
}
type Response struct {
	ID       uint16             `json:"id" description:"订阅任务ID"`
	Name     string             `json:"name" description:"订阅任务名称"`
	Enable   bool               `json:"enable" description:"是否启用"`
	CronExpr string             `json:"cron_expr" example:"0 0 * * *" description:"cron表达式"`
	Config   execer.Fetch       `json:"config" description:"订阅器配置"`
	Status   string             `json:"status" description:"订阅状态"`
	Result   taskModel.DBResult `json:"result" description:"订阅结果"`
}
