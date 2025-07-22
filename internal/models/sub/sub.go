package sub

import (
	"time"
)

type Data struct {
	ID        uint16    `db:"id" json:"id"`
	Enable    bool      `db:"enable" json:"enable"`
	Name      string    `db:"name" json:"name"`
	URL       string    `json:"url"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
}

type CreateRequest struct {
	Enable bool   `db:"enable" json:"enable"`
	Name   string `db:"name" json:"name"`
	URL    string `json:"url"`
}

type UpdateRequest struct {
	ID     uint16 `db:"id" json:"id"`
	Enable bool   `db:"enable" json:"enable"`
	Name   string `db:"name" json:"name"`
	URL    string `json:"url"`
}
