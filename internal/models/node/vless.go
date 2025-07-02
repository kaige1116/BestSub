package node

type Vless struct {
	Info   Info
	Config VlessConfig
}

type VlessConfig struct {
	BaseConfig        `yaml:",inline"`
	Uuid              string       `yaml:"uuid"`
	Flow              string       `yaml:"flow"`
	PacketEncoding    string       `yaml:"packet-encoding"`
	TLS               bool         `yaml:"tls"`
	Servername        string       `yaml:"servername"`
	Alpn              []string     `yaml:"alpn"`
	Fingerprint       string       `yaml:"fingerprint"`
	ClientFingerprint string       `yaml:"client-fingerprint"`
	SkipCertVerify    bool         `yaml:"skip-cert-verify"`
	RealityOpts       *RealityOpts `yaml:"reality-opts"`
	Network           string       `yaml:"network"`
	Smux              *SmuxOpts    `yaml:"smux"`
}
