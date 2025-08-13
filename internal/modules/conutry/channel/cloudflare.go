package channel

import (
	"bytes"
	"context"
	"io"
	"net/http"
)

type Cloudflare struct{}

func (c *Cloudflare) Get(ctx context.Context, client *http.Client) string {
	request, err := http.NewRequestWithContext(ctx, "GET", "https://cloudflare.com/cdn-cgi/trace", nil)
	if err != nil {
		return ""
	}
	response, err := client.Do(request)
	if err != nil {
		return ""
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		return ""
	}

	data, err := io.ReadAll(response.Body)
	if err != nil {
		return ""
	}

	prefix := []byte("loc=")
	idx := bytes.Index(data, prefix)
	if idx == -1 {
		return ""
	}
	start := idx + len(prefix)
	endRel := bytes.IndexByte(data[start:], '\n')
	var v []byte
	if endRel == -1 {
		v = data[start:]
	} else {
		v = data[start : start+endRel]
	}
	v = bytes.TrimSpace(v)
	return string(v)
}

func init() {
	register(&Cloudflare{})
}
