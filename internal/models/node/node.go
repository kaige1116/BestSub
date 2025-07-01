package node

type Collection struct {
	AnyTLS    []AnyTLS
	Http      []Http
	Hysteria  []Hysteria
	Hysteria2 []Hysteria2
	Mieru     []Mieru
	Snell     []Snell
	Socks     []Socks
	SS        []SS
	SSH       []SSH
	SSR       []SSR
	Trojan    []Trojan
	TUIC      []TUIC
	VLESS     []VLESS
	VMess     []VMess
	WireGuard []WireGuard
}

type Info struct {
	Id          int64
	UniqueKey   uint64
	Alive       bool
	Delay       [5]int
	SpeedUp     [5]int
	SpeedDown   [5]int
	Timestamp   int64
	Country     string
	CountryCode string
}

type BaseConfig struct {
	Name          string  `yaml:"name"`
	Type          string  `yaml:"type"`
	Server        string  `yaml:"server"`
	Port          string  `yaml:"port"`
	IpVersion     *string `yaml:"ip-version"`
	Udp           *bool   `yaml:"udp"`
	InterfaceName *string `yaml:"interface-name"`
	RoutingMark   *int    `yaml:"routing-mark"`
	Tfo           *bool   `yaml:"tfo"`
	Mptcp         *bool   `yaml:"mptcp"`
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
