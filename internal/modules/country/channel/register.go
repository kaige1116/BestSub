package channel

import (
	"net/http"
)

type Channel interface {
	Url() string
	Header(req *http.Request)
	CountryCode(body []byte) string
}

var Channels = make([]Channel, 0)

func register(channel Channel) {
	Channels = append(Channels, channel)
}
