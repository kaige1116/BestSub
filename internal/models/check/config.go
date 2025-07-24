package check

import taskModel "github.com/bestruirui/bestsub/internal/models/task"

type Data struct {
	ID     uint16 `db:"id" json:"id"`
	Name   string `db:"name" json:"name" description:"检测任务名称"`
	Enable bool   `db:"enable" json:"enable" description:"是否启用"`
	Task   string `db:"task" json:"task" description:"任务配置"`
	Config string `db:"config" json:"config" description:"检测器配置"`
	Result string `db:"result" json:"result" description:"检测结果"`
}

type CreateRequest struct {
	Name   string           `db:"name" json:"name" example:"测试检测任务" description:"检测任务名称"`
	Enable bool             `db:"enable" json:"enable" description:"是否启用"`
	Task   taskModel.Config `db:"task" json:"task" description:"任务配置"`
	Config any              `db:"config" json:"config" description:"检测器配置"`
}

type UpdateRequest struct {
	ID     uint16           `db:"id" json:"id" description:"检测任务ID"`
	Name   string           `db:"name" json:"name" description:"检测任务名称"`
	Enable bool             `db:"enable" json:"enable" description:"是否启用"`
	Task   taskModel.Config `db:"task" json:"task" description:"任务配置"`
	Config any              `db:"config" json:"config" description:"检测器配置"`
}

type Response struct {
	ID     uint16             `db:"id" json:"id" description:"检测任务ID"`
	Name   string             `db:"name" json:"name" description:"检测任务名称"`
	Enable bool               `db:"enable" json:"enable" description:"是否启用"`
	Task   taskModel.Config   `db:"task" json:"task" description:"任务配置"`
	Config any                `db:"config" json:"config" description:"检测器配置"`
	Status string             `db:"-" json:"status" description:"检测状态"`
	Result taskModel.DBResult `db:"result" json:"result" description:"检测结果"`
}
