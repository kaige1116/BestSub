package node

// MihomoConfig 统一的Mihomo配置结构体，包含所有协议的配置项
type MihomoConfig struct {
	Name          string `yaml:"name" json:"name,omitempty"`
	Type          string `yaml:"type" json:"type,omitempty"`
	Server        string `yaml:"server" json:"server,omitempty"`
	Port          any    `yaml:"port" json:"port,omitempty"`
	IpVersion     string `yaml:"ip-version" json:"ip-version,omitempty"`
	Udp           bool   `yaml:"udp" json:"udp,omitempty"`
	InterfaceName string `yaml:"interface-name" json:"interface-name,omitempty"`
	RoutingMark   int    `yaml:"routing-mark" json:"routing-mark,omitempty"`
	Tfo           bool   `yaml:"tfo" json:"tfo,omitempty"`
	Mptcp         bool   `yaml:"mptcp" json:"mptcp,omitempty"`
	DialerProxy   string `yaml:"dialer-proxy" json:"dialer-proxy,omitempty"`

	// VMess 配置
	Uuid                string    `yaml:"uuid" json:"uuid,omitempty"`
	AlterId             int       `yaml:"alterId" json:"alterId,omitempty"`
	Cipher              string    `yaml:"cipher" json:"cipher,omitempty"`
	PacketEncoding      string    `yaml:"packet-encoding" json:"packet-encoding,omitempty"`
	GlobalPadding       bool      `yaml:"global-padding" json:"global-padding,omitempty"`
	AuthenticatedLength bool      `yaml:"authenticated-length" json:"authenticated-length,omitempty"`
	TLS                 bool      `yaml:"tls" json:"tls,omitempty"`
	Servername          string    `yaml:"servername" json:"servername,omitempty"`
	Alpn                *[]string `yaml:"alpn" json:"alpn,omitempty"`
	Fingerprint         string    `yaml:"fingerprint" json:"fingerprint,omitempty"`
	ClientFingerprint   string    `yaml:"client-fingerprint" json:"client-fingerprint,omitempty"`
	SkipCertVerify      bool      `yaml:"skip-cert-verify" json:"skip-cert-verify,omitempty"`

	// VLESS 配置
	Flow string `yaml:"flow" json:"flow,omitempty"`

	// Trojan 配置
	Password string `yaml:"password" json:"password,omitempty"`
	Sni      string `yaml:"sni" json:"sni,omitempty"`

	// Shadowsocks 配置
	UdpOverTcp        bool   `yaml:"udp-over-tcp" json:"udp-over-tcp,omitempty"`
	UdpOverTcpVersion string `yaml:"udp-over-tcp-version" json:"udp-over-tcp-version,omitempty"`
	Plugin            string `yaml:"plugin" json:"plugin,omitempty"`
	PluginOpts        any    `yaml:"plugin-opts" json:"plugin-opts,omitempty"`

	// ShadowsocksR 配置
	Protocol       string `yaml:"protocol" json:"protocol,omitempty"`
	Obfs           string `yaml:"obfs" json:"obfs,omitempty"`
	ObfsParams     string `yaml:"obfs-params" json:"obfs-params,omitempty"`
	ProtocolParams string `yaml:"protocol-params" json:"protocol-params,omitempty"`
	ObfsParam      string `yaml:"obfs-param" json:"obfs-param,omitempty"`
	ProtocolParam  string `yaml:"protocol-param" json:"protocol-param,omitempty"`

	// HTTP/SOCKS 配置
	Username string `yaml:"username" json:"username,omitempty"`

	// Hysteria 配置
	Ports               string `yaml:"ports" json:"ports,omitempty"`
	AuthStr             string `yaml:"auth-str" json:"auth-str,omitempty"`
	Up                  string `yaml:"up" json:"up,omitempty"`
	Down                string `yaml:"down" json:"down,omitempty"`
	RecvWindowConn      int    `yaml:"recv-window-conn" json:"recv-window-conn,omitempty"`
	RecvWindow          int    `yaml:"recv-window" json:"recv-window,omitempty"`
	Ca                  string `yaml:"ca" json:"ca,omitempty"`
	CaStr               string `yaml:"ca-str" json:"ca-str,omitempty"`
	DisableMtuDiscovery bool   `yaml:"disable_mtu_discovery" json:"disable_mtu_discovery,omitempty"`
	FastOpen            bool   `yaml:"fast-open" json:"fast-open,omitempty"`

	// Hysteria2 配置
	ObfsPassword                   string `yaml:"obfs-password" json:"obfs-password,omitempty"`
	InitialStreamReceiveWindow     int    `yaml:"initial-stream-receive-window" json:"initial-stream-receive-window,omitempty"`
	MaxStreamReceiveWindow         int    `yaml:"max-stream-receive-window" json:"max-stream-receive-window,omitempty"`
	InitialConnectionReceiveWindow int    `yaml:"initial-connection-receive-window" json:"initial-connection-receive-window,omitempty"`
	MaxConnectionReceiveWindow     int    `yaml:"max-connection-receive-window" json:"max-connection-receive-window,omitempty"`

	// TUIC 配置
	Token                 string `yaml:"token" json:"token,omitempty"`
	Ip                    string `yaml:"ip" json:"ip,omitempty"`
	HeartbeatInterval     int    `yaml:"heartbeat-interval" json:"heartbeat-interval,omitempty"`
	DisableSni            bool   `yaml:"disable-sni" json:"disable-sni,omitempty"`
	ReduceRtt             bool   `yaml:"reduce-rtt" json:"reduce-rtt,omitempty"`
	RequestTimeout        int    `yaml:"request-timeout" json:"request-timeout,omitempty"`
	UdpRelayMode          string `yaml:"udp-relay-mode" json:"udp-relay-mode,omitempty"`
	CongestionController  string `yaml:"congestion-controller" json:"congestion-controller,omitempty"`
	MaxUdpRelayPacketSize int    `yaml:"max-udp-relay-packet-size" json:"max-udp-relay-packet-size,omitempty"`
	MaxOpenStreams        int    `yaml:"max-open-streams" json:"max-open-streams,omitempty"`

	// WireGuard 配置
	PrivateKey       string         `yaml:"private-key" json:"private-key,omitempty"`
	Ipv6             string         `yaml:"ipv6" json:"ipv6,omitempty"`
	PublicKey        string         `yaml:"public-key" json:"public-key,omitempty"`
	AllowedIps       *[]string      `yaml:"allowed-ips" json:"allowed-ips,omitempty"`
	PreSharedKey     string         `yaml:"pre-shared-key" json:"pre-shared-key,omitempty"`
	Reserved         any            `yaml:"reserved" json:"reserved,omitempty"`
	Mtu              int            `yaml:"mtu" json:"mtu,omitempty"`
	RemoteDnsResolve bool           `yaml:"remote-dns-resolve" json:"remote-dns-resolve,omitempty"`
	Dns              *[]string      `yaml:"dns" json:"dns,omitempty"`
	Peers            *[]WgPeer      `yaml:"peers" json:"peers,omitempty"`
	AmneziaWgOption  *AmneziaWgOpts `yaml:"amnezia-wg-option" json:"amnezia-wg-option,omitempty"`

	// SSH 配置
	PrivateKeyPassphrase string    `yaml:"private-key-passphrase" json:"private-key-passphrase,omitempty"`
	HostKey              *[]string `yaml:"host-key" json:"host-key,omitempty"`
	HostKeyAlgorithms    *[]string `yaml:"host-key-algorithms" json:"host-key-algorithms,omitempty"`

	// Snell 配置
	Psk      string    `yaml:"psk" json:"psk,omitempty"`
	Version  int       `yaml:"version" json:"version,omitempty"`
	ObfsOpts *ObfsOpts `yaml:"obfs-opts" json:"obfs-opts,omitempty"`

	// Mieru 配置
	PortRange    string `yaml:"port-range" json:"port-range,omitempty"`
	Transport    string `yaml:"transport" json:"transport,omitempty"`
	Multiplexing string `yaml:"multiplexing" json:"multiplexing,omitempty"`

	// AnyTLS 配置
	IdleSessionCheckInterval int `yaml:"idle-session-check-interval" json:"idle-session-check-interval,omitempty"`
	IdleSessionTimeout       int `yaml:"idle-session-timeout" json:"idle-session-timeout,omitempty"`
	MinIdleSession           int `yaml:"min-idle-session" json:"min-idle-session,omitempty"`

	// 通用网络配置
	Network string `yaml:"network" json:"network,omitempty"`

	// 复合配置选项
	RealityOpts *RealityOpts `yaml:"reality-opts" json:"reality-opts,omitempty"`
	Smux        *SmuxOpts    `yaml:"smux" json:"smux,omitempty"`
	SsOpts      *SsOpts      `yaml:"ss-opts" json:"ss-opts,omitempty"`
}

// SmuxOpts SMUX配置选项
type SmuxOpts struct {
	Enabled        bool        `yaml:"enabled" json:"enabled,omitempty"`
	Protocol       string      `yaml:"protocol" json:"protocol,omitempty"`
	MaxConnections int         `yaml:"max-connections" json:"max-connections,omitempty"`
	MinStreams     int         `yaml:"min-streams" json:"min-streams,omitempty"`
	MaxStreams     int         `yaml:"max-streams" json:"max-streams,omitempty"`
	Statistic      bool        `yaml:"statistic" json:"statistic,omitempty"`
	OnlyTcp        bool        `yaml:"only-tcp" json:"only-tcp,omitempty"`
	Padding        bool        `yaml:"padding" json:"padding,omitempty"`
	BrutalOpts     *BrutalOpts `yaml:"brutal-opts" json:"brutal-opts,omitempty"`
}

// BrutalOpts Brutal配置选项
type BrutalOpts struct {
	Enabled bool `yaml:"enabled" json:"enabled,omitempty"`
	Up      int  `yaml:"up" json:"up,omitempty"`
	Down    int  `yaml:"down" json:"down,omitempty"`
}

// RealityOpts Reality配置选项
type RealityOpts struct {
	PublicKey             string `yaml:"public-key" json:"public-key,omitempty"`
	ShortId               string `yaml:"short-id" json:"short-id,omitempty"`
	SupportX25519Mlkem768 bool   `yaml:"support-x25519mlkem768" json:"support-x25519mlkem768,omitempty"`
}

// WgPeer WireGuard对等节点配置
type WgPeer struct {
	Server       string    `yaml:"server" json:"server,omitempty"`
	Port         string    `yaml:"port" json:"port,omitempty"`
	PublicKey    string    `yaml:"public-key" json:"public-key,omitempty"`
	AllowedIps   *[]string `yaml:"allowed-ips" json:"allowed-ips,omitempty"`
	PreSharedKey string    `yaml:"pre-shared-key" json:"pre-shared-key,omitempty"`
	Reserved     any       `yaml:"reserved" json:"reserved,omitempty"` // Can be array of int or string
}

// AmneziaWgOpts AmneziaWG配置选项
type AmneziaWgOpts struct {
	Jc   int `yaml:"jc" json:"jc,omitempty"`
	Jmin int `yaml:"jmin" json:"jmin,omitempty"`
	Jmax int `yaml:"jmax" json:"jmax,omitempty"`
	S1   int `yaml:"s1" json:"s1,omitempty"`
	S2   int `yaml:"s2" json:"s2,omitempty"`
	H1   int `yaml:"h1" json:"h1,omitempty"`
	H2   int `yaml:"h2" json:"h2,omitempty"`
	H3   int `yaml:"h3" json:"h3,omitempty"`
	H4   int `yaml:"h4" json:"h4,omitempty"`
}

// ObfsOpts 混淆配置选项
type ObfsOpts struct {
	Mode string `yaml:"mode" json:"mode,omitempty"`
	Host string `yaml:"host" json:"host,omitempty"`
}

// SsOpts Shadowsocks配置选项
type SsOpts struct {
	Enabled  bool   `yaml:"enabled" json:"enabled,omitempty"`
	Method   string `yaml:"method" json:"method,omitempty"`
	Password string `yaml:"password" json:"password,omitempty"`
}
