package sub

import "time"

type SaveConfig struct {
	ID          int64     `db:"id" json:"id"`
	Enable      bool      `db:"enable" json:"enable"`           // 是否启用
	Name        string    `db:"name" json:"name"`               // 保存配置名称
	Rename      string    `db:"rename" json:"rename"`           // 保存文件名
	FileName    string    `db:"file_name" json:"file_name"`     // 保存文件名
	Description string    `db:"description" json:"description"` // 配置描述
	CreatedAt   time.Time `db:"created_at" json:"created_at"`
	UpdatedAt   time.Time `db:"updated_at" json:"updated_at"`
}
