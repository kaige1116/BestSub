package nodepool

import (
	nodeModel "github.com/bestruirui/bestsub/internal/models/node"
	subModel "github.com/bestruirui/bestsub/internal/models/sub"
)

type GlobalPool struct {
	subMap map[uint16]*SubStorage
	stats  Stats
}

type SubStorage struct {
	data  []*nodeModel.Data
	Info  subModel.NodeInfo
	index map[uint64]int
}

type Stats struct {
	SubLinkCount int32
	TotalNodes   int32
	AliveNodes   int32
	AvgDelay     uint16
	AvgSpeedUp   uint32
	AvgSpeedDown uint32
	AvgRisk      uint8
	Country      map[uint16]*CountryStats
}

type CountryStats struct {
	Count     int32
	Delay     uint16
	SpeedUp   uint32
	SpeedDown uint32
	Risk      uint8
}
