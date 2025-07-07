package nodepool

import (
	"sync"

	"github.com/bestruirui/bestsub/internal/models/node"
)

// FilterCondition 节点筛选条件（内存对齐优化）
type FilterCondition struct {
	SpeedUpMore   uint32 // 上传速度大于指定值（0表示不筛选，>0表示具体值，KB/s，最大65535KB/s）
	SpeedDownMore uint32 // 下载速度大于指定值（0表示不筛选，>0表示具体值，KB/s，最大65535KB/s）
	Country       uint16 // ISO 3166数字国家代码（0表示不筛选，>0表示具体国家，最大65535）
	DelayLessThan uint16 // 延迟小于指定值（0表示不筛选，>0表示具体值，毫秒，最大65535ms）
	AliveStatus   uint16 // 存活状态位字段筛选（0表示不筛选，其他值表示必须匹配的位）
	RiskLessThan  uint8  // 风险等级小于指定值（0表示不筛选，>0表示具体值，百分比，最大255）
}

// PoolStats 节点池统计信息
type PoolStats struct {
	TotalNodes     int            // 总节点数
	NodesByCountry map[uint16]int // 按国家分组的节点数
	NodesByStatus  map[uint16]int // 按存活状态分组的节点数
	SubLinkCount   int            // 订阅链接数
}

// index 索引结构，支持高效的节点查询（内部使用）
type index struct {
	delay     []uint64 // 按延迟排序的唯一键（升序）
	speedUp   []uint64 // 按上传速度排序的唯一键（降序）
	speedDown []uint64 // 按下载速度排序的唯一键（降序）
	risk      []uint64 // 按风险等级排序的唯一键（升序，低风险在前）

	country     map[uint16][]uint64 // 国家索引（使用ISO 3166数字代码）
	aliveStatus map[uint16][]uint64 // 存活状态位字段索引
}

// pool 节点池管理器（内部使用）
type pool struct {
	mu        sync.RWMutex         // 读写锁
	nodes     map[uint64]node.Data // 唯一键到节点数据的映射
	subs      map[int64][]uint64   // 订阅链接ID到唯一键列表的映射
	iterators map[int64]int        // 订阅链接ID到当前迭代位置的映射
	index     *index               // 索引系统
}

// 全局节点池实例
var globalPool *pool
