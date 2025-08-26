package channel

import (
	"encoding/json"
	"net/http"
)

type IPSB struct{}

func (c *IPSB) Url() string {
	return "https://api.ip.sb/geoip"
}

func (c *IPSB) Header(req *http.Request) {
	UserAgent(req)
}

func (c *IPSB) CountryCode(body []byte) string {
	var ip_sb Common
	if err := json.Unmarshal(body, &ip_sb); err != nil {
		return ""
	}
	return ip_sb.CountryCode
}

func init() {
	register(&IPSB{})
}
