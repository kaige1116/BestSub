package channel

import (
	"encoding/json"
	"net/http"
)

type MYIP struct {
	CC string `json:"cc"`
}

func (c *MYIP) Url() string {
	return "https://api.myip.com"
}

func (c *MYIP) Header(req *http.Request) {
}

func (c *MYIP) CountryCode(body []byte) string {
	var myip MYIP
	if err := json.Unmarshal(body, &myip); err != nil {
		return ""
	}
	return myip.CC
}

func init() {
	register(&MYIP{})
}
