package node

type Snell struct {
	Info   Info
	Config SnellConfig
}

type SnellConfig struct {
	BaseConfig
	Psk      string    `yaml:"psk"`
	Version  int       `yaml:"version"`
	ObfsOpts *ObfsOpts `yaml:"obfs-opts"`
}

type ObfsOpts struct {
	Mode string `yaml:"mode"`
	Host string `yaml:"host"`
}
