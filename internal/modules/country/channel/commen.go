package channel

import (
	"net/http"

	"github.com/bestruirui/bestsub/internal/utils/ua"
)

type Common struct {
	CountryCode string `json:"country_code"`
}

func UserAgent(req *http.Request) {
	ua.SetHeader(req)
}
