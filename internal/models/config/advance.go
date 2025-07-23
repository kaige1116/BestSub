package config

type Advance struct {
	Key         string `db:"key" json:"key"`
	Value       string `db:"value" json:"value"`
	Type        string `db:"-" json:"type"`
	Options     string `db:"-" json:"options,omitempty"`
	Description string `db:"-" json:"description,omitempty"`
}

type GroupAdvance struct {
	GroupName   string    `json:"group_name"`
	Description string    `json:"description"`
	Data        []Advance `json:"data"`
}

type UpdateAdvance struct {
	Key   string `json:"key" example:"proxy.enabled"`
	Value string `json:"value" example:"true"`
}
