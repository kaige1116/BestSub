package system

type Data struct {
	Key         string `db:"key" json:"key"`
	Value       string `db:"value" json:"value"`
	Type        string `db:"-" json:"type"`
	Description string `db:"-" json:"description"`
}

type GroupData struct {
	GroupName   string `json:"group_name"`
	Description string `json:"description"`
	Data        []Data `json:"data"`
}

type UpdateData struct {
	Key   string `json:"key" example:"proxy.enabled"`
	Value string `json:"value" example:"true"`
}
