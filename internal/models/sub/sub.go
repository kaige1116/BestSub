package sub

import (
	"encoding/json"
	"time"
)

type Data struct {
	ID        uint16    `db:"id" json:"id"`
	Enable    bool      `db:"enable" json:"enable"`
	Name      string    `db:"name" json:"name"`
	CronExpr  string    `db:"cron_expr" json:"cron_expr"`
	Config    string    `db:"config" json:"config"`
	Result    string    `db:"result" json:"result"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
}

type NameAndID struct {
	ID   uint16 `json:"id"`
	Name string `json:"name"`
}

type Config struct {
	Url     string `json:"url" required:"true" description:"订阅地址"`
	Proxy   bool   `json:"proxy" required:"false" description:"是否启用代理" example:"false"`
	Timeout int    `json:"timeout" required:"false" description:"超时时间单位:秒" example:"10"`
	Type    string `json:"type" required:"false" description:"订阅类型" example:"clash" optional:"clash,singbox,base64,v2ray,auto"`
}

type Result struct {
	Success  uint16    `json:"success,omitempty" description:"成功次数"`
	Fail     uint16    `json:"fail,omitempty" description:"失败次数"`
	Msg      string    `json:"msg,omitempty" description:"消息"`
	RawCount uint32    `json:"raw_count,omitempty" description:"原始节点数量"`
	Count    uint32    `json:"count,omitempty" description:"去重后的节点数量"`
	LastRun  time.Time `json:"last_run,omitempty" description:"上次运行时间"`
	Duration uint16    `json:"duration,omitempty" description:"运行时长(单位:毫秒)"`
}

type NodeInfo struct {
	RawCount   int32  `json:"raw_count" description:"原始节点数量"`
	AliveCount int32  `json:"alive_count" description:"存活节点数量"`
	SpeedUp    uint32 `json:"speed_up" description:"平均上行速度"`
	SpeedDown  uint32 `json:"speed_down" description:"平均下行速度"`
	Delay      uint16 `json:"delay" description:"平均延迟"`
	Risk       uint8  `json:"risk" description:"平均风险"`
}

type Request struct {
	Name     string `json:"name" description:"订阅任务名称"`
	Enable   bool   `json:"enable" description:"是否启用"`
	CronExpr string `json:"cron_expr" example:"0 0 * * *" description:"cron表达式"`
	Config   Config `json:"config"`
}

type Response struct {
	ID        uint16    `json:"id" description:"订阅任务ID"`
	Name      string    `json:"name" description:"订阅任务名称"`
	Enable    bool      `json:"enable" description:"是否启用"`
	CronExpr  string    `json:"cron_expr" description:"cron表达式"`
	Config    Config    `json:"config" description:"订阅器配置"`
	Status    string    `json:"status" description:"订阅状态"`
	Result    Result    `json:"result" description:"订阅结果"`
	NodeInfo  NodeInfo  `json:"node_info" description:"节点信息"`
	CreatedAt time.Time `json:"created_at" description:"创建时间"`
	UpdatedAt time.Time `json:"updated_at" description:"更新时间"`
}

func (c *Request) GenData(id uint16) Data {
	configBytes, err := json.Marshal(c.Config)
	if err != nil {
		return Data{}
	}
	return Data{
		ID:       id,
		Name:     c.Name,
		Enable:   c.Enable,
		CronExpr: c.CronExpr,
		Config:   string(configBytes),
	}
}
func (d *Data) GenResponse(status string, nodeInfo NodeInfo) Response {
	var config Config
	if err := json.Unmarshal([]byte(d.Config), &config); err != nil {
		return Response{}
	}
	var result Result
	if err := json.Unmarshal([]byte(d.Result), &result); err != nil {
		return Response{}
	}
	return Response{
		ID:        d.ID,
		Name:      d.Name,
		Enable:    d.Enable,
		CronExpr:  d.CronExpr,
		Config:    config,
		Status:    status,
		Result:    result,
		NodeInfo:  nodeInfo,
		CreatedAt: d.CreatedAt,
		UpdatedAt: d.UpdatedAt,
	}
}
