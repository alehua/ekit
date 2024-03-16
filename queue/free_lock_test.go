package queue

import (
	"sync"
	"testing"
)

func TestFreeLockQueue(t *testing.T) {
	q := NewLKQueue[int]()
	wg := &sync.WaitGroup{}
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			q.Enqueue(i)
		}(i)
	}
	wg.Wait()
	for i := 0; i < 10; i++ {
		println(q.Dequeue())
	}
}
