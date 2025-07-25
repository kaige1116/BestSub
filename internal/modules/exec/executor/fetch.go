package execer

import (
	"context"
	"time"

	"github.com/bestruirui/bestsub/internal/models/task"
	"github.com/bestruirui/bestsub/internal/modules/register"
	"github.com/bestruirui/bestsub/internal/utils/log"
)

type Fetch struct {
	SubID     uint16 `json:"sub_id,omitempty"`
	SubUrl    string `json:"sub_url" description:"订阅链接"`
	Proxy     bool   `json:"proxy" example:"false" description:"是否使用代理"`
	FailTimes uint16 `json:"fail_times" example:"0" description:"失败多少次自动禁用订阅,0为不自动禁用"`
}

func (e *Fetch) Init() error {
	return nil
}

func (e *Fetch) Run(ctx context.Context, log *log.Logger) task.ReturnResult {
	for {
		select {
		case <-ctx.Done():
			log.Infof("fetch 任务执行完成 %d", e.SubID)
			return task.ReturnResult{
				Status:          true,
				LastRunResult:   "fetch 任务执行完成",
				LastRunTime:     time.Now(),
				LastRunDuration: 0,
			}
		default:
			time.Sleep(1 * time.Second)
			log.Infof("fetch 任务执行中 %d", e.SubID)
		}
	}
}

func init() {
	register.Exec(&Fetch{})
}
