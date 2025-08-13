package country

import (
	"context"
	"net/http"

	"github.com/bestruirui/bestsub/internal/modules/conutry/channel"
)

func GetCode(ctx context.Context, client *http.Client) string {
	for _, channel := range channel.Channels {
		country := channel.Get(ctx, client)
		if country != "" {
			return country
		}
	}
	return ""
}
