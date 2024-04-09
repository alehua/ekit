package time_wheel

import (
	"container/list"
	"context"
	"fmt"
	"time"
)

func NewTimeWheel(ctx context.Context,
	slotsNum int, interval time.Duration) *TimeWheel {
	tw := &TimeWheel{
		Interval:     interval,
		ticker:       time.NewTicker(interval),
		stop:         make(chan struct{}),
		addTaskCh:    make(chan *TaskElement), //无缓存
		removeTaskCh: make(chan string),
		keyToTask:    make(map[string]*list.Element),
		slots:        make([]*list.List, slotsNum),
	}
	for i := 0; i < slotsNum; i++ {
		tw.slots[i] = list.New()
	}

	go tw.run(ctx)

	return tw
}

func (tw *TimeWheel) AddTask(key string,
	task func(ctx context.Context) error,
	executeAt time.Time) {
	pos, cycle := tw.getPositionAndCycle(executeAt)
	tw.addTaskCh <- &TaskElement{key: key, pos: pos, cycle: cycle, task: task}
}

func (tw *TimeWheel) RemoveTask(key string) {
	tw.removeTaskCh <- key
}

func (tw *TimeWheel) Stop() {
	tw.once.Do(func() {
		close(tw.stop)
	})
}

func (tw *TimeWheel) run(ctx context.Context) {
	defer func() {
		er := recover()
		if err, ok := er.(error); ok {
			fmt.Println("时间轮运行出错, error=", err.Error())
		}
		tw.ticker.Stop()
	}()
	for {
		select {
		case <-ctx.Done():
			return
		case <-tw.stop:
			return
		case <-tw.ticker.C:
			tw.tick()
		case task := <-tw.addTaskCh:
			tw.add(task)
		case key := <-tw.removeTaskCh:
			tw.remove(key)
		}
	}
}

func (tw *TimeWheel) tick() {
	lst := tw.slots[tw.cur]
	defer func() {
		tw.cur = (tw.cur + 1) % len(tw.slots)
	}()
	for e := lst.Front(); e != nil; {
		task, _ := e.Value.(*TaskElement)
		if task.cycle > 0 {
			task.cycle--
			e = e.Next()
			continue
		}
		go func(key string, task func(ctx context.Context) error) {
			// 执行任务
			ctx, cancel := context.WithTimeout(context.Background(), time.Second)
			defer cancel()
			defer func() {
				if err := recover(); err != nil {
					fmt.Println("任务执行异常, key=", key)
				}
			}()
			err := task(ctx)
			if err != nil {
				fmt.Println("任务执行异常, key=", key)
			}
		}(task.key, task.task)
		e = e.Next()
		tw.RemoveTask(task.key) // 执行后移除
	}
}

func (tw *TimeWheel) add(task *TaskElement) {
	lst := tw.slots[task.pos]
	// 优化这里需要加锁
	if _, ok := tw.keyToTask[task.key]; ok {
		tw.remove(task.key)
	}
	// 将任务添加到节点尾部
	newTask := lst.PushBack(task)
	tw.keyToTask[task.key] = newTask
}

func (tw *TimeWheel) remove(key string) {
	oldTask, ok := tw.keyToTask[key]
	if !ok {
		return
	}
	delete(tw.keyToTask, key)
	task, _ := oldTask.Value.(*TaskElement)
	_ = tw.slots[task.pos].Remove(oldTask)
}

// getPositionAndCycle 获取任务在时间轮中的位置和周期 核心方法
func (tw *TimeWheel) getPositionAndCycle(executeAt time.Time) (int, int) {
	delay := int(executeAt.Sub(time.Now()).Seconds())
	// delay := time.Until(executeAt)
	cycle := delay / (len(tw.slots) * int(tw.Interval.Seconds()))
	pos := (tw.cur + delay/int(tw.Interval.Seconds())) % len(tw.slots)
	return pos, cycle
}
