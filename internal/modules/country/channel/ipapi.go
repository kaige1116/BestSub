package channel

import (
	"encoding/json"
	"net/http"
)

type IPAPI struct{}

func (c *IPAPI) Url() string {
	return "https://ipapi.co/json"
}

func (c *IPAPI) Header(req *http.Request) {

}

func (c *IPAPI) CountryCode(body []byte) string {
	var ipapi Common
	if err := json.Unmarshal(body, &ipapi); err != nil {
		return ""
	}
	return ipapi.CountryCode
}

func init() {
	register(&IPAPI{})
}
