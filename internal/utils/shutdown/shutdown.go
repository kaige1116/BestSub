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
						fmt.Printf("Closing functions panic: %v", r)
					}
				}()
				if err := fn(); err != nil {
					fmt.Printf("Closing functions execution failed: %v", err)
				}
			}(funcs[i])
		}
	}()

	select {
	case <-done:
	case <-ctx.Done():
		fmt.Printf("Closing functions execution timeout: %v", ctx.Err())
	}
}

func Listen() {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	log.Info("Program started, press Ctrl+C to exit")
	sig := <-quit
	log.Warnf("Received exit signal: %v, starting to close program", sig)
	execute()
	os.Exit(0)
}
