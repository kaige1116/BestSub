package storage

import (
	"context"
	"encoding/json"
)

type Data struct {
	ID     uint16 `db:"id" json:"id"`
	Name   string `db:"name" json:"name"`
	Type   string `db:"type" json:"type"`
	Config string `db:"config" json:"config"`
}

type Request struct {
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

func (r *Request) GenData(id uint16) Data {
	configBytes, _ := json.Marshal(r.Config)
	return Data{
		ID:     id,
		Name:   r.Name,
		Type:   r.Type,
		Config: string(configBytes),
	}
}

func (d *Data) GenResponse() Response {
	var config any
	json.Unmarshal([]byte(d.Config), &config)
	return Response{
		ID:     d.ID,
		Name:   d.Name,
		Type:   d.Type,
		Config: config,
	}
}
