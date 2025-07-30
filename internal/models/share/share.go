package share

import (
	"encoding/json"

	nodeModel "github.com/bestruirui/bestsub/internal/models/node"
)

type Data struct {
	ID          uint16 `db:"id" json:"id"`
	Enable      bool   `db:"enable" json:"enable"`
	Name        string `db:"name" json:"name"`
	Token       string `db:"token" json:"token"`
	Config      string `db:"config" json:"config"`
	AccessCount uint32 `db:"access_count" json:"access_count"`
}
type Config struct {
	Template       string           `json:"template" description:"模板，支持：clash, v2ray, surge"`
	MaxAccessCount uint32           `json:"max_access_count" description:"最大访问次数"`
	Expires        uint64           `json:"expires" example:"1722336000" description:"过期时间"`
	SubID          []uint16         `json:"sub_id" description:"订阅ID"`
	Filter         nodeModel.Filter `json:"filter" description:"筛选条件"`
}

type Request struct {
	Enable bool   `json:"enable"`
	Name   string `json:"name"`
	Token  string `json:"token"`
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

func (r *Request) GenData() Data {
	configBytes, err := json.Marshal(r.Config)
	if err != nil {
		return Data{}
	}
	return Data{
		Enable: r.Enable,
		Name:   r.Name,
		Token:  r.Token,
		Config: string(configBytes),
	}
}

func (r *Data) GenResponse() Response {
	var config Config
	if err := json.Unmarshal([]byte(r.Config), &config); err != nil {
		return Response{}
	}
	return Response{
		ID:          r.ID,
		Name:        r.Name,
		Token:       r.Token,
		Enable:      r.Enable,
		AccessCount: r.AccessCount,
		Config:      config,
	}
}
