package shutdown

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/bestruirui/bestsub/internal/utils/log"
)

const timeout = 10 * time.Second

var funcs []func() error

func Register(fn func() error) {
	funcs = append(funcs, fn)
}

func execute() {
	if len(funcs) == 0 {
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	done := make(chan struct{})
	go func() {
		defer close(done)
		for i := 0; i < len(funcs); i++ {
			func(fn func() error) {
				defer func() {
					if r := recover(); r != nil {
						fmt.Printf("关闭函数 发生panic: %v", r)
					}
				}()
				if err := fn(); err != nil {
					fmt.Printf("关闭函数 执行失败: %v", err)
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

func Listen() {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	log.Info("程序已启动，按 Ctrl+C 退出")
	sig := <-quit
	log.Warnf("收到退出信号: %v，开始关闭程序", sig)
	execute()
	os.Exit(0)
}
