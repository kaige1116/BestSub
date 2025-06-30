package node

type Hysteria struct {
	Info   Info
	Config HysteriaConfig
}

type HysteriaConfig struct {
	BaseConfig
	Ports               string   `yaml:"ports"`
	AuthStr             string   `yaml:"auth-str"`
	Obfs                string   `yaml:"obfs"`
	Alpn                []string `yaml:"alpn"`
	Protocol            string   `yaml:"protocol"`
	Up                  string   `yaml:"up"`
	Down                string   `yaml:"down"`
	Sni                 string   `yaml:"sni"`
	SkipCertVerify      bool     `yaml:"skip-cert-verify"`
	RecvWindowConn      int      `yaml:"recv-window-conn"`
	RecvWindow          int      `yaml:"recv-window"`
	Ca                  string   `yaml:"ca"`
	CaStr               string   `yaml:"ca-str"`
	DisableMtuDiscovery bool     `yaml:"disable_mtu_discovery"`
	Fingerprint         string   `yaml:"fingerprint"`
	FastOpen            bool     `yaml:"fast-open"`
}
