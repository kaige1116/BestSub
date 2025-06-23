package shutdown

import (
	"context"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/bestruirui/bestsub/internal/utils/log"
)

// 关闭函数类型
type ShutdownFunc func(ctx context.Context) error

// 退出管理器
type Shutdown struct {
	funcs []ShutdownFunc
	mu    sync.Mutex
}

var instance *Shutdown

func init() {
	instance = &Shutdown{
		funcs: make([]ShutdownFunc, 0),
	}
}

// 注册清理函数
func (s *Shutdown) register(fn ShutdownFunc) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.funcs = append(s.funcs, fn)
	log.Debug("已注册退出清理函数")
}

// 监听退出信号
func (s *Shutdown) listen() {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	log.Debug("退出监听器已启动")

	sig := <-quit
	log.Warnf("程序即将退出: %v", sig)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	s.mu.Lock()
	funcs := make([]ShutdownFunc, len(s.funcs))
	copy(funcs, s.funcs)
	s.mu.Unlock()

	for i := len(funcs) - 1; i >= 0; i-- {
		if err := funcs[i](ctx); err != nil {
			log.Errorf("清理函数执行失败: %v", err)
		}
	}

	log.Debug("退出完成")
	os.Exit(0)
}

// 注册清理函数
func Register(fn ShutdownFunc) {
	instance.register(fn)
}

// 监听退出信号
func Listen() {
	instance.listen()
}
