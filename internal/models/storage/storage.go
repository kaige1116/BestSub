package storage

import (
	"context"
)

type Data struct {
	ID     uint16 `db:"id" json:"id"`
	Name   string `db:"name" json:"name"`
	Type   string `db:"type" json:"type"`
	Config string `db:"config" json:"config"`
}

type CreateRequest struct {
	Name   string `json:"name" example:"webdav"`
	Type   string `json:"type" example:"webdav"`
	Config any    `json:"config"`
}

type UpdateRequest struct {
	ID     uint16 `json:"id"`
	Name   string `json:"name" example:"webdav"`
	Type   string `json:"type" example:"webdav"`
	Config any    `json:"config"`
}

type Response struct {
	ID     uint16 `json:"id"`
	Name   string `json:"name"`
	Type   string `json:"type"`
	Config any    `json:"config"`
}

type Instance interface {
	Init() error
	Upload(ctx context.Context) error
}
