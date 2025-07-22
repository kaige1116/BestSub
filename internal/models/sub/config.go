package sub

// 全局管理保存方式配置，供订阅输出系统引用
type StorageConfig struct {
	ID     uint16 `db:"id" json:"id"`
	Enable bool   `db:"enable" json:"enable"` // 是否启用
	Name   string `db:"name" json:"name"`     // 存储配置名称
	Type   string `db:"type" json:"type"`     // 存储类型：webdav, local, ftp, sftp, s3, oss
	Config string `db:"config" json:"config"` // 存储配置（JSON格式）
}

// 全局输出模板配置
type OutputTemplate struct {
	ID       uint16 `db:"id" json:"id"`
	Name     string `db:"name" json:"name"`         // 模板名称
	Type     string `db:"type" json:"type"`         // 模板类型：mihomo, singbox, v2ray, clash
	Template string `db:"template" json:"template"` // 模板内容
}

// 全局节点输出筛选规则
type NodeFilterRule struct {
	ID       uint16 `db:"id" json:"id"`
	Name     string `db:"name" json:"name"`         // 规则名称
	Field    string `db:"field" json:"field"`       // 规则字段
	Operator string `db:"operator" json:"operator"` // 规则操作符
	Value    string `db:"value" json:"value"`       // 规则值
}
