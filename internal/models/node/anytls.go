package node

type AnyTLS struct {
	Info   Info
	Config AnyTLSConfig
}

type AnyTLSConfig struct {
	BaseConfig               `yaml:",inline"`
	Password                 string   `yaml:"password"`
	ClientFingerprint        string   `yaml:"client-fingerprint"`
	IdleSessionCheckInterval int      `yaml:"idle-session-check-interval"`
	IdleSessionTimeout       int      `yaml:"idle-session-timeout"`
	MinIdleSession           int      `yaml:"min-idle-session"`
	Sni                      string   `yaml:"sni"`
	Alpn                     []string `yaml:"alpn"`
	SkipCertVerify           *bool    `yaml:"skip-cert-verify"`
}
