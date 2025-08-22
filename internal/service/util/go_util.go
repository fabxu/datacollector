package util

import (
	"context"
	"sync"
	"time"

	cmlog "github.com/fabxu/log"
)

func DeferRecover(ctx context.Context) {
	logger := cmlog.Extract(ctx)
	err := recover()
	if err == nil {
		return
	}
	logger.Errorf("【PANIC】: %s", err)
}

// StartTick 启动定时器
func StartTick(ctx context.Context, intervalTime time.Duration, t func() error) {
	logger := cmlog.Extract(ctx)
	f := func() {
		defer DeferRecover(ctx)
		if err := t(); err != nil {
			logger.Errorf("tick err: %s", err.Error())
		}
	}

	go func() {
		ticker := time.NewTicker(intervalTime)
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				f()
			case <-ctx.Done():
				return
			}
		}
	}()
}
func GoWithErr(ctx context.Context, wg *sync.WaitGroup, f func() error) {
	logger := cmlog.Extract(ctx)
	if wg != nil {
		wg.Add(1)
	}
	go func(handler func() error) {
		// defer 是倒序触发执, 先入的后执行
		defer func() {
			if wg != nil {
				wg.Done()
			}
		}()
		defer DeferRecover(ctx)

		if e := handler(); e != nil {
			logger.Errorf("go run fail: %v", e)
		}
	}(f)
}
