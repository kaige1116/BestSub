package channel

import (
	"encoding/json"
	"net/http"
)

type FreeIP struct{}

func (c *FreeIP) Url() string {
	return "https://free.freeipapi.com/api/json"
}

func (c *FreeIP) Header(req *http.Request) {
	UserAgent(req)
}

func (c *FreeIP) CountryCode(body []byte) string {
	var freeip struct {
		CountryCode string `json:"countryCode"`
	}
	if err := json.Unmarshal(body, &freeip); err != nil {
		return ""
	}
	return freeip.CountryCode
}

func init() {
	register(&FreeIP{})
}
