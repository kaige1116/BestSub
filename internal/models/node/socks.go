package node

type Socks struct {
	Info   Info
	Config SocksConfig
}

type SocksConfig struct {
	BaseConfig     `yaml:",inline"`
	Username       string `yaml:"username"`
	Password       string `yaml:"password"`
	TLS            bool   `yaml:"tls"`
	SkipCertVerify bool   `yaml:"skip-cert-verify"`
	Sni            string `yaml:"sni"`
	Fingerprint    string `yaml:"fingerprint"`
	IpVersion      string `yaml:"ip-version"`
}
