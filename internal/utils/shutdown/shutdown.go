package shutdown

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/bestruirui/bestsub/internal/utils/log"
)

// 关闭函数定义
type ShutdownFunc struct {
	Name string
	Fn   func() error
}

// 关闭管理器
type Manager struct {
	funcs   []ShutdownFunc
	mu      sync.RWMutex
	timeout time.Duration
}

var (
	defaultManager *Manager
	once           sync.Once
)

// 获取默认管理器
func getInstance() *Manager {
	once.Do(func() {
		defaultManager = &Manager{
			funcs:   make([]ShutdownFunc, 0),
			timeout: 10 * time.Second, // 默认10秒超时
		}
	})
	return defaultManager
}

// 注册关闭函数
func (m *Manager) register(fn ShutdownFunc) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.funcs = append(m.funcs, fn)
	log.Debugf("已注册退出清理函数: %s", fn.Name)
}

// 执行关闭函数
func (m *Manager) execute() {
	m.mu.RLock()
	funcs := make([]ShutdownFunc, len(m.funcs))
	copy(funcs, m.funcs)
	timeout := m.timeout
	m.mu.RUnlock()

	if len(funcs) == 0 {
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	done := make(chan struct{})
	go func() {
		defer close(done)
		for i := len(funcs) - 1; i >= 0; i-- {
			func(fn ShutdownFunc) {
				defer func() {
					if r := recover(); r != nil {
						fmt.Printf("关闭函数 %s 发生panic: %v", fn.Name, r)
					}
				}()
				if err := fn.Fn(); err != nil {
					fmt.Printf("关闭函数 %s 执行失败: %v", fn.Name, err)
				}
			}(funcs[i])
		}
	}()

	select {
	case <-done:
	case <-ctx.Done():
		fmt.Printf("关闭函数执行超时: %v", ctx.Err())
	}
}

// 监听退出信号
func (m *Manager) listen() {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	log.Info("程序已启动，按 Ctrl+C 退出")
	sig := <-quit
	log.Warnf("收到退出信号: %v，开始关闭程序", sig)
	m.execute()
	os.Exit(0)
}

func Register(name string, fn func() error) {
	getInstance().register(ShutdownFunc{Name: name, Fn: fn})
}

func Listen() {
	getInstance().listen()
}
