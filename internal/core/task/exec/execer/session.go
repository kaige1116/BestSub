package execer

import (
	"context"

	"github.com/bestruirui/bestsub/internal/core/session"
	"github.com/bestruirui/bestsub/internal/core/task/exec"
	"github.com/bestruirui/bestsub/internal/models/task"
	"github.com/bestruirui/bestsub/internal/utils/local"
	"github.com/bestruirui/bestsub/internal/utils/log"
)

// init 自动注册Session处理器
func init() {
	exec.Register(&exec.RegisterInfo{
		Type:    task.TypeSessionClean,
		Handler: &SessionExec{},
		Config:  &SessionConfig{},
	})
}

// SessionExec 会话清理任务处理器
type SessionExec struct{}

// SessionConfig 会话清理任务配置
type SessionConfig struct{}

// Do 执行会话清理任务
func (h *SessionExec) Do(ctx context.Context, logger *log.Logger, task *exec.TaskInfo) error {
	startTime := local.Time()
	session.Cleanup()
	logger.Infof("会话清理任务执行完成，耗时: %d ms", local.Time().Sub(startTime).Milliseconds())
	return nil
}
