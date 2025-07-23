package sub

type Template struct {
	ID       uint16 `db:"id" json:"id"`
	Name     string `db:"name" json:"name"`
	Type     string `db:"type" json:"type"`         // 模板类型：mihomo, singbox, v2ray, clash
	Template string `db:"template" json:"template"` // 模板内容
}
