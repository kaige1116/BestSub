package node

type SSR struct {
	Info   Info
	Config SSRConfig
}

type SSRConfig struct {
	BaseConfig
	Password       string `yaml:"password"`
	Protocol       string `yaml:"protocol"`
	Obfs           string `yaml:"obfs"`
	ObfsParams     string `yaml:"obfs-params"`
	ProtocolParams string `yaml:"protocol-params"`
	ObfsParam      string `yaml:"obfs-param"`
	ProtocolParam  string `yaml:"protocol-param"`
	Udp            bool   `yaml:"udp"`
}
