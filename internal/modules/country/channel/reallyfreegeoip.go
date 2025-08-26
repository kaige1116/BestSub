package channel

import (
	"encoding/json"
	"net/http"
)

type ReallyFreeGeoIP struct{}

func (c *ReallyFreeGeoIP) Url() string {
	return "https://reallyfreegeoip.org/json"
}

func (c *ReallyFreeGeoIP) Header(req *http.Request) {
	UserAgent(req)
}

func (c *ReallyFreeGeoIP) CountryCode(body []byte) string {
	var reallyfreegeoip struct {
		CountryCode string `json:"country_code"`
	}
	if err := json.Unmarshal(body, &reallyfreegeoip); err != nil {
		return ""
	}
	return reallyfreegeoip.CountryCode
}

func init() {
	register(&ReallyFreeGeoIP{})
}
