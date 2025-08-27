package node

import (
	"encoding/json"

	"github.com/bestruirui/bestsub/internal/utils/generic"
	"github.com/cespare/xxhash/v2"
)

const (
	Alive        uint16 = 1 << 0
	Country      uint16 = 1 << 1
	AliveCustom2 uint16 = 1 << 2
	AliveCustom3 uint16 = 1 << 3
	AliveCustom4 uint16 = 1 << 4
	AliveCustom5 uint16 = 1 << 5
	AliveCustom6 uint16 = 1 << 6
	AliveCustom7 uint16 = 1 << 7
)

type Data struct {
	Base
	Info *Info
}

type Base struct {
	Raw       []byte
	SubId     uint16
	UniqueKey uint64
}

type UniqueKey struct {
	Server     string `yaml:"server"`
	Servername string `yaml:"servername"`
	Port       string `yaml:"port"`
	Type       string `yaml:"type"`
	Uuid       string `yaml:"uuid"`
	Username   string `yaml:"username"`
	Password   string `yaml:"password"`
}

type Info struct {
	SpeedUp     generic.Queue[uint32]
	SpeedDown   generic.Queue[uint32]
	Delay       generic.Queue[uint16]
	Risk        uint8
	AliveStatus uint16
	IP          uint32
	Country     string
}

type SimpleInfo struct {
	SpeedUp   uint32 `json:"speed_up"`
	SpeedDown uint32 `json:"speed_down"`
	Delay     uint16 `json:"delay"`
	Risk      uint8  `json:"risk"`
	Count     uint32 `json:"count"`
}

type Filter struct {
	SubId         []uint16 `json:"sub_id"`
	SpeedUpMore   uint32   `json:"speed_up_more"`
	SpeedDownMore uint32   `json:"speed_down_more"`
	Country       []string `json:"country"`
	DelayLessThan uint16   `json:"delay_less_than"`
	AliveStatus   uint16   `json:"alive_status"`
	RiskLessThan  uint8    `json:"risk_less_than"`
}

func (i *Info) SetAliveStatus(AliveStatus uint16, status bool) {
	if status {
		i.AliveStatus |= AliveStatus
	} else {
		i.AliveStatus &= ^AliveStatus
	}
}

func (u *UniqueKey) Gen() uint64 {
	bytes, _ := json.Marshal(u)
	return xxhash.Sum64(bytes)
}
