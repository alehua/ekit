package time_wheel

import (
	"container/list"
	"context"
	"sync"
	"time"
)

type TimeWheel struct {
	once         sync.Once     // 时间轮全局单例控制
	Interval     time.Duration // 时间轮间隔
	ticker       *time.Ticker  // 定时器
	stop         chan struct{}
	addTaskCh    chan *TaskElement // 添加任务通道
	removeTaskCh chan string       // 移除任务通道

	slots []*list.List // 时间轮槽
	cur   int          // 当前时间轮槽

	keyToTask map[string]*list.Element
}

type TaskElement struct {
	task  func(ctx context.Context) error
	pos   int    // 任务在时间轮中的位置
	cycle int    // 还需要在时间轮转多少圈
	key   string // 任务唯一标识
}

type ITaskElement interface {
	Name() string
}
