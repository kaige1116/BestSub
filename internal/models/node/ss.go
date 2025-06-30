package node

type SS struct {
	Info   Info
	Config SSConfig
}

type SSConfig struct {
	BaseConfig
	Cipher            string    `yaml:"cipher"`
	Password          string    `yaml:"password"`
	Udp               bool      `yaml:"udp"`
	UdpOverTcp        bool      `yaml:"udp-over-tcp"`
	UdpOverTcpVersion string    `yaml:"udp-over-tcp-version"`
	IPVersion         string    `yaml:"ip-version"`
	Plugin            string    `yaml:"plugin"`
	ClientFingerprint *string   `yaml:"client-fingerprint"`
	PluginOpts        *any      `yaml:"plugin-opts"`
	Smux              *SmuxOpts `yaml:"smux"`
}

type SSPluginObfs struct {
	Mode string `yaml:"mode"`
	Host string `yaml:"host"`
}

type SSPluginV2ray struct {
	Mode             string            `yaml:"mode"`
	TLS              bool              `yaml:"tls"`
	Fingerprint      string            `yaml:"fingerprint"`
	SkipCertVerify   bool              `yaml:"skip-cert-verify"`
	Host             string            `yaml:"host"`
	Path             string            `yaml:"path"`
	Mux              bool              `yaml:"mux"`
	Headers          map[string]string `yaml:"headers"`
	V2rayHttpUpgrade bool              `yaml:"v2ray-http-upgrade"`
}

type SSPluginGost struct {
	Mode           string            `yaml:"mode"`
	TLS            bool              `yaml:"tls"`
	Fingerprint    string            `yaml:"fingerprint"`
	SkipCertVerify bool              `yaml:"skip-cert-verify"`
	Host           string            `yaml:"host"`
	Path           string            `yaml:"path"`
	Mux            bool              `yaml:"mux"`
	Headers        map[string]string `yaml:"headers"`
}

type SSPluginShadowtls struct {
	Mode     string `yaml:"mode"`
	Password string `yaml:"password"`
	Version  string `yaml:"version"`
}
type SSPluginRestls struct {
	Host         string `yaml:"host"`
	Password     string `yaml:"password"`
	VersionHint  string `yaml:"version-hint"`
	RestlsScript string `yaml:"restls-script"`
}
