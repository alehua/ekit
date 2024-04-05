package queue

import (
	"fmt"
	"sync"
	"time"
)

type DelayQueue[T Delayable] struct {
	mu    sync.Mutex // 保护优先级队列
	cond  *sync.Cond
	items *Heap[T]
}

func NewDelayQueue[T Delayable]() *DelayQueue[T] {
	return &DelayQueue[T]{
		cond: sync.NewCond(&sync.Mutex{}),
		items: NewHeapPriorityQueue[T](func(t1, t2 T) int {
			srcDelay := t1.Delay()
			dstDelay := t2.Delay()
			if srcDelay > dstDelay {
				return 1
			}
			if srcDelay == dstDelay {
				return 0
			}
			return -1
		}),
	}
}

// Add 向延迟队列中添加元素
func (dq *DelayQueue[T]) Add(value T) {
	dq.mu.Lock()
	defer dq.mu.Unlock()

	dq.items.Push(value)
	if dq.items.Len() == 1 {
		dq.cond.Signal() // 如果有第一个元素，则唤醒等待的goroutine
	}
}

// Take 从延迟队列中获取到期的元素
func (dq *DelayQueue[T]) Take() (T, bool) {
	dq.mu.Lock()
	defer dq.mu.Unlock()
	for {
		val, err := dq.items.Peek()
		switch err {
		case nil:
			delay := val.Delay()
			if delay <= 0 {
				t := dq.items.Pop()
				return t, true
			}
			dq.cond.Wait()
		case fmt.Errorf("heap empty"):
			dq.cond.Wait()
		default:
			var t T
			return t, false
		}
	}
}

type Delayable interface {
	Delay() time.Duration
}
