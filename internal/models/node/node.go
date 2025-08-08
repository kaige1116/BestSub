package node

import "github.com/bestruirui/bestsub/internal/utils/generic"

// AliveStatus 位字段常量定义
const (
	Alive        uint16 = 1 << 0 // 第0位：节点存活
	AliveCustom1 uint16 = 1 << 1 // 第1位：自定义测试1存活
	AliveCustom2 uint16 = 1 << 2 // 第2位：自定义测试2存活
	AliveCustom3 uint16 = 1 << 3 // 第3位：自定义测试3存活
	AliveCustom4 uint16 = 1 << 4 // 第4位：自定义测试4存活
	AliveCustom5 uint16 = 1 << 5 // 第5位：自定义测试5存活
	AliveCustom6 uint16 = 1 << 6 // 第6位：自定义测试6存活
	AliveCustom7 uint16 = 1 << 7 // 第7位：自定义测试7存活
)

type Data struct {
	Base
	Info *Info // 节点Info结构体（值类型）
}

type Base struct {
	Raw       []byte // 节点配置的JSON格式数据
	SubId     uint16 // 订阅ID
	UniqueKey uint64 // 唯一键
}

type Info struct {
	SpeedUp     generic.Queue[uint32] // 单位 KB/s
	SpeedDown   generic.Queue[uint32] // 单位 KB/s
	Delay       generic.Queue[uint16] // 单位 ms
	Risk        uint8
	AliveStatus uint16
	IP          uint32
	Country     uint16
}

type SimpleInfo struct {
	SpeedUp   uint32 `json:"speed_up"`
	SpeedDown uint32 `json:"speed_down"`
	Delay     uint16 `json:"delay"`
	Risk      uint8  `json:"risk"`
	Count     uint32 `json:"count"`
}

type Filter struct {
	SubId         uint16   // 订阅ID
	SpeedUpMore   uint32   // 上传速度大于指定值（0表示不筛选，>0表示具体值，KB/s，最大65535KB/s）
	SpeedDownMore uint32   // 下载速度大于指定值（0表示不筛选，>0表示具体值，KB/s，最大65535KB/s）
	Country       []uint16 // ISO 3166数字国家代码（0表示不筛选，>0表示具体国家，最大65535）
	DelayLessThan uint16   // 延迟小于指定值（0表示不筛选，>0表示具体值，毫秒，最大65535ms）
	AliveStatus   uint16   // 存活状态位字段筛选（0表示不筛选，其他值表示必须匹配的位）
	RiskLessThan  uint8    // 风险等级小于指定值（0表示不筛选，>0表示具体值，百分比，最大255）
}

func (i *Info) SetAliveStatus(AliveStatus uint16, status bool) {
	if status {
		i.AliveStatus |= AliveStatus
	} else {
		i.AliveStatus &= ^AliveStatus
	}
}
