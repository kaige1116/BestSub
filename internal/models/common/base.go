package common

import "time"

// BaseDbModel 基础模型，包含所有实体的公共字段
type BaseDbModel struct {
	ID          int64     `db:"id" json:"id"`
	Enable      bool      `db:"enable" json:"enable"`
	Name        string    `db:"name" json:"name"`
	Description string    `db:"description" json:"description"`
	CreatedAt   time.Time `db:"created_at" json:"created_at"`
	UpdatedAt   time.Time `db:"updated_at" json:"updated_at"`
}

type BaseRequestModel struct {
	Enable      *bool  `json:"enable"`
	Name        string `json:"name"`
	Description string `json:"description"`
}
type BaseUpdateRequestModel struct {
	ID          int64  `json:"id"`
	Enable      *bool  `json:"enable"`
	Name        string `json:"name"`
	Description string `json:"description"`
}
