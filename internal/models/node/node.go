package node

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
	Config []byte // 节点配置的JSON格式数据
	Info   Info   // 节点Info结构体（值类型）
}

type Info struct {
	Timestamp [5]int64  // 单位 ms
	SpeedUp   [5]uint32 // 单位 KB/s
	SpeedDown [5]uint32 // 单位 KB/s
	Delay     [5]uint16 // 单位 ms

	AddTime     int64  // 入库时间戳
	UniqueKey   uint64 // 唯一键
	AliveStatus uint16
	Country     uint16

	Risk uint8
}

type BaseConfig struct {
	Name          string  `yaml:"name"`
	Type          string  `yaml:"type"`
	Server        string  `yaml:"server"`
	Port          any     `yaml:"port"`
	IpVersion     *string `yaml:"ip-version"`
	Udp           bool    `yaml:"udp"`
	InterfaceName *string `yaml:"interface-name"`
	RoutingMark   *int    `yaml:"routing-mark"`
	Tfo           bool    `yaml:"tfo"`
	Mptcp         bool    `yaml:"mptcp"`
	DialerProxy   *string `yaml:"dialer-proxy"`
}

type SmuxOpts struct {
	Enabled        bool        `yaml:"enabled"`
	Protocol       string      `yaml:"protocol"`
	MaxConnections int         `yaml:"max-connections"`
	MinStreams     int         `yaml:"min-streams"`
	MaxStreams     int         `yaml:"max-streams"`
	Statistic      bool        `yaml:"statistic"`
	OnlyTcp        bool        `yaml:"only-tcp"`
	Padding        bool        `yaml:"padding"`
	BrutalOpts     *BrutalOpts `yaml:"brutal-opts"`
}

type BrutalOpts struct {
	Enabled bool `yaml:"enabled"`
	Up      int  `yaml:"up"`
	Down    int  `yaml:"down"`
}

type RealityOpts struct {
	PublicKey             string `yaml:"public-key"`
	ShortId               string `yaml:"short-id"`
	SupportX25519Mlkem768 bool   `yaml:"support-x25519mlkem768"`
}
