package nodepool

import (
	"sync"

	"github.com/bestruirui/bestsub/internal/models/node"
)

type Collection struct {
	AnyTLS    []node.AnyTLS
	Http      []node.Http
	Hysteria  []node.Hysteria
	Hysteria2 []node.Hysteria2
	Mieru     []node.Mieru
	Snell     []node.Snell
	Socks     []node.Socks
	Ss        []node.Ss
	Ssh       []node.Ssh
	Ssr       []node.Ssr
	Trojan    []node.Trojan
	Tuic      []node.Tuic
	Vless     []node.Vless
	Vmess     []node.Vmess
	WireGuard []node.WireGuard
}

// NodePool 内存优化的节点池
type NodePool struct {
	collection *Collection
	indexes    *NodeIndexes
	mu         sync.RWMutex
	totalNodes int64
}

// NodeIndexes 多维度索引结构
type NodeIndexes struct {
	uniqueKey map[uint64]*NodeInfo   // 唯一键索引，用于去重
	subLinkID map[int64][]*NodeInfo  // 按订阅链接ID索引
	nodeType  map[string][]*NodeInfo // 按节点类型索引
}

// NodeInfo 节点信息
type NodeInfo struct {
	ArrayIndex int // 在数组中的索引位置
}

// reflectCache 反射缓存结构
type reflectCache struct {
	collectionFields map[string]int         // Collection字段名到索引的映射
	nodeConfigFields map[string][]fieldInfo // 节点类型到Config字段信息的映射
}

type fieldInfo struct {
	index         int    // 字段在结构体中的索引
	name          string // 字段名称
	isEmbedded    bool   // 是否为内嵌字段
	embeddedIndex int    // 如果是内嵌字段，这是内嵌字段在父结构体中的索引
}

// NodeIterator 节点迭代器，用于多线程安全的顺序访问
type NodeIterator struct {
	mu           sync.Mutex
	currentType  int      // 当前节点类型索引
	currentIndex int      // 当前类型中的节点索引
	typeNames    []string // 节点类型名称列表
	finished     bool     // 是否已遍历完成
	subLinkID    int64    // 关联的订阅链接ID
}

// IteratorManager 迭代器管理器，管理多个订阅链接的迭代器
type IteratorManager struct {
	mu        sync.RWMutex
	iterators map[int64]*NodeIterator // subLinkID -> NodeIterator
}

const (
	// 默认容量配置
	defaultUniqueKeyCapacity = 10000
	defaultSubLinkCapacity   = 100
	defaultNodeTypeCapacity  = 10
	maxStringSliceCapacity   = 16
)

var (
	reflectCacheInstance *reflectCache
	globalIteratorMgr    *IteratorManager
	iteratorMgrOnce      sync.Once
	cachedTypeNames      []string

	// 用于生成唯一键的字段名列表
	uniqueKeyFields = map[string]bool{
		"Server":     true,
		"Port":       true,
		"Username":   true,
		"Password":   true,
		"AuthStr":    true,
		"Uuid":       true,
		"Servername": true,
	}
)

// createNodeIndexes 创建新的节点索引
func createNodeIndexes() *NodeIndexes {
	return &NodeIndexes{
		uniqueKey: make(map[uint64]*NodeInfo, defaultUniqueKeyCapacity),
		subLinkID: make(map[int64][]*NodeInfo, defaultSubLinkCapacity),
		nodeType:  make(map[string][]*NodeInfo, defaultNodeTypeCapacity),
	}
}
