package setting

import "github.com/bestruirui/bestsub/internal/utils/desc"

type SettingAdvance = desc.Data

type GroupSettingAdvance struct {
	GroupName   string           `json:"group_name"`
	Description string           `json:"description"`
	Data        []SettingAdvance `json:"data"`
}

type Setting struct {
	Key   string `json:"key" example:"proxy.enabled"`
	Value string `json:"value" example:"true"`
}
