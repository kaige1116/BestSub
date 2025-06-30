package node

type Http struct {
	Info   Info
	Config HttpConfig
}

type HttpConfig struct {
	BaseConfig
	Username       string `yaml:"username"`
	Password       string `yaml:"password"`
	TLS            bool   `yaml:"tls"`
	SkipCertVerify bool   `yaml:"skip-cert-verify"`
	Sni            string `yaml:"sni"`
	Fingerprint    string `yaml:"fingerprint"`
	IpVersion      string `yaml:"ip-version"`
}
