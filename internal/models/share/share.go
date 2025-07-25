package share

import "time"

type Data struct {
	ID          uint16 `db:"id" json:"id"`
	Enable      bool   `db:"enable" json:"enable"`
	Name        string `db:"name" json:"name"`
	Token       string `db:"token" json:"token"`
	Config      string `db:"config" json:"config"`
	AccessCount uint32 `db:"access_count" json:"access_count"`
}
type Config struct {
	MaxAccessCount uint32    `json:"max_access_count"`
	Expires        time.Time `json:"expires" example:"2030-07-25T10:00:00Z"`
	SubID          []uint16  `json:"sub_id"`
}
type CreateRequest struct {
	Enable bool   `json:"enable"`
	Name   string `json:"name"`
	Token  string `json:"token"`
	Config Config `json:"config"`
}
type UpdateRequest struct {
	ID     uint16 `json:"id"`
	Name   string `json:"name"`
	Token  string `json:"token"`
	Enable bool   `json:"enable"`
	Config Config `json:"config"`
}
type Response struct {
	ID          uint16 `json:"id"`
	Name        string `json:"name"`
	Token       string `json:"token"`
	Enable      bool   `json:"enable"`
	AccessCount uint32 `json:"access_count"`
	Config      Config `json:"config"`
}
type UpdateAccessCountDB struct {
	ID          uint16 `db:"id"`
	AccessCount uint32 `db:"access_count"`
}
