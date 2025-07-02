package node

type Hysteria2 struct {
	Info   Info
	Config Hysteria2Config
}

type Hysteria2Config struct {
	BaseConfig                     `yaml:",inline"`
	Ports                          string   `yaml:"ports"`
	Password                       string   `yaml:"password"`
	Up                             string   `yaml:"up"`
	Down                           string   `yaml:"down"`
	Obfs                           string   `yaml:"obfs"`
	ObfsPassword                   string   `yaml:"obfs-password"`
	Sni                            string   `yaml:"sni"`
	SkipCertVerify                 bool    `yaml:"skip-cert-verify"`
	Fingerprint                    string   `yaml:"fingerprint"`
	Alpn                           []string `yaml:"alpn"`
	Ca                             string   `yaml:"ca"`
	CaStr                          string   `yaml:"ca-str"`
	InitialStreamReceiveWindow     int      `yaml:"initial-stream-receive-window"`
	MaxStreamReceiveWindow         int      `yaml:"max-stream-receive-window"`
	InitialConnectionReceiveWindow int      `yaml:"initial-connection-receive-window"`
	MaxConnectionReceiveWindow     int      `yaml:"max-connection-receive-window"`
}
