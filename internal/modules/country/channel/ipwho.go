package channel

import (
	"encoding/json"
	"net/http"
)

type IPWho struct{}

func (c *IPWho) Url() string {
	return "https://api.ip.sb/geoip"
}

func (c *IPWho) Header(req *http.Request) {
	UserAgent(req)
}

func (c *IPWho) CountryCode(body []byte) string {
	var ipwho Common
	if err := json.Unmarshal(body, &ipwho); err != nil {
		return ""
	}
	return ipwho.CountryCode
}

func init() {
	register(&IPWho{})
}
