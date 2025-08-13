package channel

import (
	"context"
	"net/http"
)

type Channel interface {
	Get(ctx context.Context, client *http.Client) string
}

var Channels = make([]Channel, 0)

func register(channel Channel) {
	Channels = append(Channels, channel)
}
