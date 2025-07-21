package channel

import (
	"bytes"

	"github.com/bestruirui/bestsub/internal/modules/register"
)

type WebHook struct {
	Url string `desc:"url" type:"string" required:"true" description:"WebHook地址"`
}

func (e *WebHook) Init() error {
	return nil
}

func (e *WebHook) Send(title string, body *bytes.Buffer) error {
	return nil
}

func init() {
	register.Notify(&WebHook{})
}
