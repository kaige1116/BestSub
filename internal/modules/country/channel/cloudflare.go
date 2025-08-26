package channel

import (
	"bytes"
	"encoding/json"
	"net/http"
)

type CloudflareCDN struct{}

func (c *CloudflareCDN) Url() string {
	return "https://cloudflare.com/cdn-cgi/trace"
}

func (c *CloudflareCDN) Header(req *http.Request) {
}

func (c *CloudflareCDN) CountryCode(body []byte) string {
	prefix := []byte("loc=")
	idx := bytes.Index(body, prefix)
	if idx == -1 {
		return ""
	}
	start := idx + len(prefix)
	endRel := bytes.IndexByte(body[start:], '\n')
	var v []byte
	if endRel == -1 {
		v = body[start:]
	} else {
		v = body[start : start+endRel]
	}
	v = bytes.TrimSpace(v)
	return string(v)
}

type CloudflareSpeed struct{}

func (c *CloudflareSpeed) Url() string {
	return "https://speed.cloudflare.com/meta"
}

func (c *CloudflareSpeed) Header(req *http.Request) {
	UserAgent(req)
}

func (c *CloudflareSpeed) CountryCode(body []byte) string {
	var speed struct {
		CountryCode string `json:"country"`
	}
	if err := json.Unmarshal(body, &speed); err != nil {
		return ""
	}
	return speed.CountryCode
}

func init() {
	register(&CloudflareCDN{})
	register(&CloudflareSpeed{})
}
