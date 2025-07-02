package node

type Trojan struct {
	Info   Info
	Config TrojanConfig
}

type TrojanConfig struct {
	BaseConfig        `yaml:",inline"`
	Password          string       `yaml:"password"`
	Sni               string       `yaml:"sni"`
	Alpn              []string     `yaml:"alpn"`
	ClientFingerprint string       `yaml:"client-fingerprint"`
	Fingerprint       string       `yaml:"fingerprint"`
	SkipCertVerify    bool         `yaml:"skip-cert-verify"`
	SsOpts            *SsOpts      `yaml:"ss-opts"`
	RealityOpts       *RealityOpts `yaml:"reality-opts"`
	Network           string       `yaml:"network"`
	Smux              *SmuxOpts    `yaml:"smux"`
}

type SsOpts struct {
	Enabled  bool   `yaml:"enabled"`
	Method   string `yaml:"method"`
	Password string `yaml:"password"`
}
