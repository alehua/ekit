package queue

import (
	"container/heap"
	"fmt"
	"github.com/alehua/ekit/containerx"
)

// 通过heapX实现优先级队列

type Heap[T any] struct {
	data *containerx.HeapX[T]
}

// Push pushes the element x onto the heap.
func (h *Heap[T]) Push(v T) {
	heap.Push(h.data, v)
}

func (h *Heap[T]) Pop() T {
	return heap.Pop(h.data).(T)
}

func (h *Heap[T]) Peek() (T, error) {
	if h.data.Len() == 0 {
		return *new(T), fmt.Errorf("heap empty")
	}
	return h.data.Data[0], nil
}

func (h *Heap[T]) Element(index int) (t T, err error) {
	if index < 0 || index >= h.data.Len() {
		return t, fmt.Errorf("out of index")
	}
	return h.data.Data[index], nil
}

func (h *Heap[T]) Len() int {
	return h.data.Len()
}

func (h *Heap[T]) Remove(index int) T {
	return heap.Remove(h.data, index).(T)
}

// NewHeap return Heap pointer and init the heap tree
func NewHeap[T any](t []T, cmp containerx.Cmp[T]) *Heap[T] {
	ret := containerx.HeapX[T]{Data: t, Cmp: cmp}
	heap.Init(&ret)
	return &Heap[T]{&ret}
}

func NewHeapPriorityQueue[T any](cmp containerx.Cmp[T]) *Heap[T] {
	return &Heap[T]{data: containerx.NewHeapX(cmp)}
}
