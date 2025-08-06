package config

import "github.com/bestruirui/bestsub/internal/utils/desc"

type Advance = desc.Data

type GroupAdvance struct {
	GroupName   string    `json:"group_name"`
	Description string    `json:"description"`
	Data        []Advance `json:"data"`
}

type UpdateAdvance struct {
	Key   string `json:"key" example:"proxy.enabled"`
	Value string `json:"value" example:"true"`
}
