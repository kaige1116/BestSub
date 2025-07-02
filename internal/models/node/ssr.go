package node

type Ssr struct {
	Info   Info
	Config SsrConfig
}

type SsrConfig struct {
	BaseConfig     `yaml:",inline"`
	Password       string `yaml:"password"`
	Protocol       string `yaml:"protocol"`
	Obfs           string `yaml:"obfs"`
	ObfsParams     string `yaml:"obfs-params"`
	ProtocolParams string `yaml:"protocol-params"`
	ObfsParam      string `yaml:"obfs-param"`
	ProtocolParam  string `yaml:"protocol-param"`
}
