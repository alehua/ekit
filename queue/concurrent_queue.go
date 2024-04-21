package queue

import (
	"github.com/alehua/ekit/internal/queue"
	"sync/atomic"
	"unsafe"
)

// ConcurrentQueue 并发队列， 基于链表实现
type ConcurrentQueue[T any] struct {
	// *node[T]
	head unsafe.Pointer
	// *node[T]
	tail unsafe.Pointer
}

type node[T any] struct {
	val T
	// *node[T]
	next unsafe.Pointer
}

func NewConcurrentQueue[T any]() *ConcurrentQueue[T] {
	head := &node[T]{}
	ptr := unsafe.Pointer(head)
	return &ConcurrentQueue[T]{
		head: ptr,
		tail: ptr,
	}
}

func (c *ConcurrentQueue[T]) Enqueue(t T) error {
	newNode := &node[T]{val: t}
	newPtr := unsafe.Pointer(newNode)
	for {
		tailPtr := atomic.LoadPointer(&c.tail)
		tail := (*node[T])(tailPtr)
		tailNext := atomic.LoadPointer(&tail.next)
		if tailNext != nil {
			continue
		}
		if atomic.CompareAndSwapPointer(&tail.next, tailNext, newPtr) {
			// 如果失败也不用担心，说明有人抢先一步了
			atomic.CompareAndSwapPointer(&c.tail, tailPtr, newPtr)
			return nil
		}
	}
}

func (c *ConcurrentQueue[T]) Dequeue() (T, error) {
	for {
		headPtr := atomic.LoadPointer(&c.head)
		head := (*node[T])(headPtr)
		tailPtr := atomic.LoadPointer(&c.tail)
		tail := (*node[T])(tailPtr)
		if head == tail {
			var t T
			return t, queue.ErrEmptyQueue
		}
		headNextPtr := atomic.LoadPointer(&head.next)
		if atomic.CompareAndSwapPointer(&c.head, headPtr, headNextPtr) {
			headNext := (*node[T])(headNextPtr)
			return headNext.val, nil
		}
	}
}
