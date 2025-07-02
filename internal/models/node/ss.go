package node

type Ss struct {
	Info   Info
	Config SsConfig
}

type SsConfig struct {
	BaseConfig        `yaml:",inline"`
	Cipher            string    `yaml:"cipher"`
	Password          string    `yaml:"password"`
	UdpOverTcp        bool      `yaml:"udp-over-tcp"`
	UdpOverTcpVersion string    `yaml:"udp-over-tcp-version"`
	Plugin            string    `yaml:"plugin"`
	ClientFingerprint *string   `yaml:"client-fingerprint"`
	PluginOpts        *any      `yaml:"plugin-opts"`
	Smux              *SmuxOpts `yaml:"smux"`
}

type SsPluginObfs struct {
	Mode string `yaml:"mode"`
	Host string `yaml:"host"`
}

type SsPluginV2ray struct {
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

type SsPluginGost struct {
	Mode           string            `yaml:"mode"`
	TLS            bool              `yaml:"tls"`
	Fingerprint    string            `yaml:"fingerprint"`
	SkipCertVerify bool              `yaml:"skip-cert-verify"`
	Host           string            `yaml:"host"`
	Path           string            `yaml:"path"`
	Mux            bool              `yaml:"mux"`
	Headers        map[string]string `yaml:"headers"`
}

type SsPluginShadowtls struct {
	Mode     string `yaml:"mode"`
	Password string `yaml:"password"`
	Version  string `yaml:"version"`
}
type SsPluginRestls struct {
	Host         string `yaml:"host"`
	Password     string `yaml:"password"`
	VersionHint  string `yaml:"version-hint"`
	RestlsScript string `yaml:"restls-script"`
}
