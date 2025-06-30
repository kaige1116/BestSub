package node

type WireGuard struct {
	Info   Info
	Config WireGuardConfig
}

type WireGuardConfig struct {
	BaseConfig
	PrivateKey       string         `yaml:"private-key"`
	Ip               string         `yaml:"ip"`
	Ipv6             string         `yaml:"ipv6"`
	PublicKey        string         `yaml:"public-key"`
	AllowedIps       []string       `yaml:"allowed-ips"`
	PreSharedKey     string         `yaml:"pre-shared-key"`
	Reserved         any            `yaml:"reserved"` // Can be array of int or string
	Udp              bool           `yaml:"udp"`
	Mtu              int            `yaml:"mtu"`
	DialerProxy      string         `yaml:"dialer-proxy"`
	RemoteDnsResolve bool           `yaml:"remote-dns-resolve"`
	Dns              []string       `yaml:"dns"`
	Peers            []WgPeer       `yaml:"peers"`
	AmneziaWgOption  *AmneziaWgOpts `yaml:"amnezia-wg-option"`
}

type WgPeer struct {
	Server       string   `yaml:"server"`
	Port         string   `yaml:"port"`
	PublicKey    string   `yaml:"public-key"`
	AllowedIps   []string `yaml:"allowed-ips"`
	PreSharedKey string   `yaml:"pre-shared-key"`
	Reserved     any      `yaml:"reserved"` // Can be array of int or string
}

type AmneziaWgOpts struct {
	Jc   int `yaml:"jc"`
	Jmin int `yaml:"jmin"`
	Jmax int `yaml:"jmax"`
	S1   int `yaml:"s1"`
	S2   int `yaml:"s2"`
	H1   int `yaml:"h1"`
	H2   int `yaml:"h2"`
	H3   int `yaml:"h3"`
	H4   int `yaml:"h4"`
}
