package node

type AnyTLS struct {
	Info   Info
	Config AnyTLSConfig
}

type AnyTLSConfig struct {
	BaseConfig
	Password                 string   `yaml:"password"`
	ClientFingerprint        string   `yaml:"client-fingerprint"`
	Udp                      bool     `yaml:"udp"`
	IdleSessionCheckInterval int      `yaml:"idle-session-check-interval"`
	IdleSessionTimeout       int      `yaml:"idle-session-timeout"`
	MinIdleSession           int      `yaml:"min-idle-session"`
	Sni                      string   `yaml:"sni"`
	Alpn                     []string `yaml:"alpn"`
	SkipCertVerify           bool     `yaml:"skip-cert-verify"`
}
