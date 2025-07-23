package storage

type Data struct {
	ID     uint16 `db:"id" json:"id"`
	Name   string `db:"name" json:"name"`     // 存储配置名称
	Type   string `db:"type" json:"type"`     // 存储类型：webdav, local, ftp, sftp, s3, oss
	Config string `db:"config" json:"config"` // 存储配置（JSON格式）
}
