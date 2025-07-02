package node

type Ssh struct {
	Info   Info
	Config SshConfig
}

type SshConfig struct {
	BaseConfig           `yaml:",inline"`
	Username             string   `yaml:"username"`
	Password             string   `yaml:"password"`
	PrivateKey           string   `yaml:"private-key"`
	PrivateKeyPassphrase string   `yaml:"private-key-passphrase"`
	HostKey              []string `yaml:"host-key"`
	HostKeyAlgorithms    []string `yaml:"host-key-algorithms"`
}
