package poolx

import (
	"context"
	"fmt"
	"runtime"
)

// TaskFunc 一个可执行的任务
type TaskFunc func(ctx context.Context) error

// Run 执行任务
func (t TaskFunc) Run(ctx context.Context) error { return t(ctx) }

// taskWrapper 是Task的装饰器
// 可以加耗时监控、重试等策略
type taskWrapper struct {
	t Task
}

func (tw *taskWrapper) Run(ctx context.Context) error {
	defer func() {
		// 处理 panic
		const panicBuffLen = 2048
		if err := recover(); err != nil {
			buf := make([]byte, panicBuffLen)
			buf = buf[:runtime.Stack(buf, false)]
			err = fmt.Errorf("任务运行异常：%s", fmt.Sprintf("\t%+v\n%s\n", err, buf))
		}
	}()
	return tw.t.Run(ctx)
}

// Task 任务
type Task interface {
	// Run 执行任务
	Run(ctx context.Context) error
}
