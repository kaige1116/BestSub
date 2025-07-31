package nodepool

import (
	"sync/atomic"

	nodeModel "github.com/bestruirui/bestsub/internal/models/node"
	subModel "github.com/bestruirui/bestsub/internal/models/sub"
	"github.com/bestruirui/bestsub/internal/utils/generic"
)

var Global = &GlobalPool{
	subMap: make(map[uint16]*SubStorage),
	stats: Stats{
		Country: make(map[uint16]*CountryStats),
	},
}

func GetPoolBySubID(subID uint16, size int) *SubStorage {
	subStorage, ok := Global.subMap[subID]
	if !ok {
		subStorage = &SubStorage{
			data:  make([]*nodeModel.Data, 0, size),
			index: make(map[uint64]int, size),
			Info:  subModel.NodeInfo{},
		}
		Global.subMap[subID] = subStorage
	}
	return subStorage
}

func DeletePool(subID uint16) {
	atomic.AddInt32(&Global.stats.TotalNodes, -int32(len(Global.subMap[subID].data)))
	delete(Global.subMap, subID)
}

func UpdateStats() {
	var sumNodeCount, sumAliveCount, sumSubCount, sumDelay, sumSpeedUp, sumSpeedDown, sumRisk uint64
	for _, sub := range Global.subMap {
		sub.UpdateInfo()
		sumNodeCount += uint64(sub.Info.RawCount)
		sumAliveCount += uint64(sub.Info.AliveCount)
		sumSubCount++
		sumDelay += uint64(sub.Info.Delay)
		sumSpeedUp += uint64(sub.Info.SpeedUp)
		sumSpeedDown += uint64(sub.Info.SpeedDown)
		sumRisk += uint64(sub.Info.Risk)
	}
	Global.stats.SubLinkCount = int32(sumSubCount)
	Global.stats.TotalNodes = int32(sumNodeCount)
	Global.stats.AliveNodes = int32(sumAliveCount)
	if sumAliveCount > 0 {
		Global.stats.AvgDelay = uint16(sumDelay / sumAliveCount)
		Global.stats.AvgSpeedUp = uint32(sumSpeedUp / sumAliveCount)
		Global.stats.AvgSpeedDown = uint32(sumSpeedDown / sumAliveCount)
		Global.stats.AvgRisk = uint8(sumRisk / sumAliveCount)
	} else {
		Global.stats.AvgDelay = 0
		Global.stats.AvgSpeedUp = 0
		Global.stats.AvgSpeedDown = 0
		Global.stats.AvgRisk = 0
	}
}

func (p *SubStorage) AddNode(data *[]nodeModel.Data) uint32 {
	for _, node := range *data {
		if p.index[node.Info.UniqueKey] != 0 {
			continue
		}
		node.Info.Delay = *generic.NewQueue[uint16](5)
		node.Info.SpeedUp = *generic.NewQueue[uint32](5)
		node.Info.SpeedDown = *generic.NewQueue[uint32](5)
		p.data = append(p.data, &node)
		p.index[node.Info.UniqueKey] = len(p.data) - 1
	}
	return uint32(len(p.data))
}

func (p *SubStorage) DeleteNode(uniqueKey uint64) {
	index, ok := p.index[uniqueKey]
	if !ok {
		return
	}
	if p.data[index].Info.Country != 0 {
		Global.stats.Country[p.data[index].Info.Country].Count--
	}
	Global.stats.TotalNodes--
	p.data[index] = p.data[len(p.data)-1]
	p.index[p.data[index].Info.UniqueKey] = index
	p.data = p.data[:len(p.data)-1]

	delete(p.index, uniqueKey)
}

func (p *SubStorage) UpdateInfo() {
	p.Info.RawCount = int32(len(p.data))
	var sumDelay, sumSpeedUp, sumSpeedDown, sumRisk uint64
	var country = make(map[uint16]*struct {
		Count     int32
		Delay     uint64
		SpeedUp   uint64
		SpeedDown uint64
		Risk      uint64
	}, 256)
	for _, node := range p.data {
		if node.Info.AliveStatus != 0 {
			p.Info.AliveCount++
			sumDelay += uint64(node.Info.Delay.Average())
			sumSpeedUp += uint64(node.Info.SpeedUp.Average())
			sumSpeedDown += uint64(node.Info.SpeedDown.Average())
			sumRisk += uint64(node.Info.Risk)
			if node.Info.Country != 0 {
				country[node.Info.Country].Count++
				country[node.Info.Country].Delay += uint64(node.Info.Delay.Average())
				country[node.Info.Country].SpeedUp += uint64(node.Info.SpeedUp.Average())
				country[node.Info.Country].SpeedDown += uint64(node.Info.SpeedDown.Average())
				country[node.Info.Country].Risk += uint64(node.Info.Risk)
			}
		}
	}
	for country, stats := range country {
		old, ok := Global.stats.Country[country]
		if !ok {
			Global.stats.Country[country] = &CountryStats{
				Count:     stats.Count,
				Delay:     uint16(stats.Delay / uint64(stats.Count)),
				SpeedUp:   uint32(stats.SpeedUp / uint64(stats.Count)),
				SpeedDown: uint32(stats.SpeedDown / uint64(stats.Count)),
				Risk:      uint8(stats.Risk / uint64(stats.Count)),
			}
			continue
		}
		old.Count += stats.Count
		old.Delay = uint16((uint64(old.Delay) + stats.Delay) / uint64(old.Count))
		old.SpeedUp = uint32((uint64(old.SpeedUp) + stats.SpeedUp) / uint64(old.Count))
		old.SpeedDown = uint32((uint64(old.SpeedDown) + stats.SpeedDown) / uint64(old.Count))
		old.Risk = uint8((uint64(old.Risk) + stats.Risk) / uint64(old.Count))
	}
	p.Info.Delay = uint16(sumDelay / uint64(p.Info.AliveCount))
	p.Info.SpeedUp = uint32(sumSpeedUp / uint64(p.Info.AliveCount))
	p.Info.SpeedDown = uint32(sumSpeedDown / uint64(p.Info.AliveCount))
	p.Info.Risk = uint8(sumRisk / uint64(p.Info.AliveCount))
}

func (p *SubStorage) GetAllNode() []*nodeModel.Data {
	return p.data
}

func (p *SubStorage) FilterNode(filter nodeModel.Filter, fn func([]byte)) {
	for _, node := range p.data {
		if node.Info.SpeedUp.Average() <= filter.SpeedUpMore && filter.SpeedUpMore != 0 {
			continue
		}
		if node.Info.SpeedDown.Average() <= filter.SpeedDownMore && filter.SpeedDownMore != 0 {
			continue
		}
		if node.Info.Delay.Average() >= filter.DelayLessThan && filter.DelayLessThan != 0 {
			continue
		}
		if node.Info.Risk >= filter.RiskLessThan && filter.RiskLessThan != 0 {
			continue
		}
		fn(node.Raw)
	}
}
