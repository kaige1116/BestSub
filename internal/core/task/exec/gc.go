package execer

import (
	"context"
	"runtime"
	"runtime/debug"
	"time"

	"github.com/bestruirui/bestsub/internal/models/task"
	"github.com/bestruirui/bestsub/internal/modules/register"
	"github.com/bestruirui/bestsub/internal/utils/log"
)

const (
	bytesToMB = 1024 * 1024
)

type GCExec struct {
	ForceGC bool `desc:"force" default:"false" description:"是否强制执行垃圾回收，true时会调用debug.FreeOSMemory()"`
}

func (e *GCExec) Init() error {
	return nil
}

func (e *GCExec) Exec(ctx context.Context, log *log.Logger, args ...any) error {
	startTime := time.Now()

	var memBefore runtime.MemStats
	runtime.ReadMemStats(&memBefore)

	if e.ForceGC {
		debug.FreeOSMemory()
	} else {
		runtime.GC()
	}

	var memAfter runtime.MemStats
	runtime.ReadMemStats(&memAfter)

	duration := time.Since(startTime)
	memFreed := float64(memBefore.Alloc-memAfter.Alloc) / bytesToMB

	log.Infof("垃圾回收任务执行完成，耗时: %d ms, 释放内存: %.2f MB", duration.Milliseconds(), memFreed)
	return nil
}

func init() {
	register.Exec(task.TypeGC, &GCExec{})
}
