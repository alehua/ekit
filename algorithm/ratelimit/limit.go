package ratelimit

import (
	"context"
	"errors"
	"github.com/alehua/ekit/internal/queue"
	"sync"
	"sync/atomic"
	"time"
)

type handler func(ctx context.Context, handler func(ctx context.Context) error) error

// CountLimit 计数器算法
type CountLimit struct {
	cnt       atomic.Int32
	threshold int32
}

func (c *CountLimit) BuildCountLimit() handler {
	return func(ctx context.Context, fn func(ctx context.Context) error) error {
		cnt := c.cnt.Add(1)
		defer func() {
			c.cnt.Add(-1)
		}()
		if cnt <= c.threshold {
			return fn(ctx)
		}
		return errors.New("限流")
	}
}

// FixedWindowLimit 固定窗口
type FixedWindowLimit struct {
	window          time.Duration
	lastWindowStart time.Time
	cnt             int32
	threshold       int32
	lock            sync.Mutex
}

func (f *FixedWindowLimit) BuildFixedWindowLimit() handler {
	return func(ctx context.Context, fn func(ctx context.Context) error) error {
		f.lock.Lock()
		now := time.Now()
		// 判断窗口是否过期
		if now.After(f.lastWindowStart.Add(f.window)) {
			f.lastWindowStart = now
			f.cnt = 0
		}
		f.cnt = f.cnt + 1
		f.lock.Unlock()
		if f.cnt <= f.threshold {
			return fn(ctx)
		}
		return errors.New("限流")
	}
}

// SlidingWindowLimit 滑动窗口
type SlidingWindowLimit struct {
	window time.Duration
	queue  queue.PriorityQueue[time.Time]

	lock      sync.Mutex
	threshold int32
}

func (s *SlidingWindowLimit) BuildSlidingWindowLimit() handler {
	return func(ctx context.Context, fn func(ctx context.Context) error) error {
		s.lock.Lock()
		// 看队列里时间差是否咋时间内
		now := time.Now()
		windowStart := now.Add(-s.window)
		for {
			first, _ := s.queue.Peek()
			if first.Before(windowStart) {
				_, _ = s.queue.Dequeue()
			} else {
				break
			}
		}
		if s.queue.Len() < int(s.threshold) {
			_ = s.queue.Enqueue(now)
			s.lock.Unlock()
			return fn(ctx)
		}
		s.lock.Unlock()
		return errors.New("限流")
	}
}

type TokenBucketLimit struct {
	interval  time.Duration
	buckets   chan struct{}
	closeCh   chan struct{}
	closeOnce sync.Once
}

func (t *TokenBucketLimit) BuildTokenBucketLimit() handler {
	return func(ctx context.Context, fn func(ctx context.Context) error) error {
		ticker := time.NewTicker(t.interval)
		defer ticker.Stop()
		go func() {
			for {
				select {
				case <-ticker.C:
					select {
					case t.buckets <- struct{}{}:
					default:
						// 桶满了
					}
				case <-t.closeCh:
					return
				}
			}

		}()
		select {
		case <-t.buckets:
			return fn(ctx)
		default:
			return errors.New("限流")
		}
	}
}

func (t *TokenBucketLimit) Close() {
	t.closeOnce.Do(func() {
		close(t.closeCh)
	})
}

type LeakyBucketLimit struct {
	interval  time.Duration
	closeCh   chan struct{}
	closeOnce sync.Once
}

func (l *LeakyBucketLimit) BuildLeakyBucketLimit() handler {
	return func(ctx context.Context, fn func(ctx context.Context) error) error {
		ticker := time.NewTicker(l.interval)
		defer ticker.Stop()
		select {
		case <-ticker.C:
			return fn(ctx)
		default:
			return errors.New("限流")
		}
	}
}

func (l *LeakyBucketLimit) Close() {
	l.closeOnce.Do(func() {
		close(l.closeCh)
	})
}
