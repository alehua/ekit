package queue

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestIntHeap(t *testing.T) {
	// 大顶堆
	queuex := NewHeap[int]([]int{1, 2, 3}, func(i int, j int) int {
		return i - j
	})
	queuex.Push(4)
	queuex.Push(5)
	queuex.Push(0)
	n := queuex.Len()
	for i := 0; i < n; i++ {
		val := queuex.Pop()
		t.Log(val)
		assert.Equal(t, i, val)
	}
}

func TestIntHeapUp(t *testing.T) {
	// 小顶堆
	queuex := NewHeap[int]([]int{1, 2, 3}, func(i int, j int) int {
		return j - i
	})
	queuex.Push(4)
	queuex.Push(5)
	queuex.Push(0)
	n := queuex.Len()
	for i := 0; i < n; i++ {
		val := queuex.Pop()
		t.Log(val)
		assert.Equal(t, n-i-1, val)
	}
}

func TestStructHeap(t *testing.T) {
	// 小顶堆
	type cases struct {
		name string
		age  int
	}
	queuex := NewHeap[cases]([]cases{
		cases{"a", 1},
		cases{"b", 2},
	}, func(i, j cases) int {
		return i.age - j.age
	})
	queuex.Push(cases{"c", 3})
	queuex.Push(cases{"d", 0})
	n := queuex.Len()
	for i := 0; i < n; i++ {
		val := queuex.Pop()
		t.Log(val)
		assert.Equal(t, i, val.age)
	}
}
