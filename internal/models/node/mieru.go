package node

type Mieru struct {
	Info   Info
	Config MieruConfig
}

type MieruConfig struct {
	BaseConfig
	PortRange    string `yaml:"port-range"`
	Transport    string `yaml:"transport"`
	Username     string `yaml:"username"`
	Password     string `yaml:"password"`
	Multiplexing string `yaml:"multiplexing"`
}
