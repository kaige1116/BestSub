package node

type Tuic struct {
	Info   Info
	Config TuicConfig
}

type TuicConfig struct {
	BaseConfig            `yaml:",inline"`
	Token                 string   `yaml:"token"`
	Uuid                  string   `yaml:"uuid"`
	Password              string   `yaml:"password"`
	Ip                    string   `yaml:"ip"`
	HeartbeatInterval     int      `yaml:"heartbeat-interval"`
	Alpn                  []string `yaml:"alpn"`
	DisableSni            bool     `yaml:"disable-sni"`
	ReduceRtt             bool     `yaml:"reduce-rtt"`
	RequestTimeout        int      `yaml:"request-timeout"`
	UdpRelayMode          string   `yaml:"udp-relay-mode"`
	CongestionController  string   `yaml:"congestion-controller"`
	MaxUdpRelayPacketSize int      `yaml:"max-udp-relay-packet-size"`
	FastOpen              bool     `yaml:"fast-open"`
	SkipCertVerify        bool     `yaml:"skip-cert-verify"`
	MaxOpenStreams        int      `yaml:"max-open-streams"`
	Sni                   string   `yaml:"sni"`
}
