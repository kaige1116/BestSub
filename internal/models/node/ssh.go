package node

type SSH struct {
	Info   Info
	Config SSHConfig
}

type SSHConfig struct {
	BaseConfig
	Username             string   `yaml:"username"`
	Password             string   `yaml:"password"`
	PrivateKey           string   `yaml:"private-key"`
	PrivateKeyPassphrase string   `yaml:"private-key-passphrase"`
	HostKey              []string `yaml:"host-key"`
	HostKeyAlgorithms    []string `yaml:"host-key-algorithms"`
}
