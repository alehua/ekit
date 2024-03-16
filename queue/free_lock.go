package queue

import (
	"sync/atomic"
	"unsafe"
)

// Node 队列每一个节点数据, next指向下一个节点
type Node[T any] struct {
	value T
	next  unsafe.Pointer
}

type LKQueue[T any] struct {
	head, tail unsafe.Pointer
}

func NewLKQueue[T any]() *LKQueue[T] {
	n := unsafe.Pointer(&Node[T]{})
	return &LKQueue[T]{head: n, tail: n}
}

// Enqueue 入队
// 通过CAS操作将一个元素添加到队尾, 并移动队尾指针
func (f *LKQueue[T]) Enqueue(v T) {
	val := &Node[T]{value: v}
	for {
		// 读取末尾的数据
		tail := load[T](&f.tail)
		next := load[T](&tail.next)

		// 确认tail是队尾, (高并发场景tail可能会变化)
		if tail == load[T](&f.tail) {
			if next == nil {
				if cas[T](&tail.next, next, val) { // 将新val放到tail.next
					cas[T](&f.tail, tail, val) // 入队完成, tail移动到下一个
					return
				}
			} else {
				// tail 没有指向最后一个节点, 将tail挪到下一个节点
				cas[T](&f.tail, tail, next)
			}
		}
	}
}

// Dequeue 出队
func (f *LKQueue[T]) Dequeue() T {
	var t T
	for {
		head := load[T](&f.head)
		tail := load[T](&f.tail)
		next := load[T](&head.next)
		// 检查 head、tail、next是否一致
		if head == load[T](&f.head) {
			if head == tail { // 队列为空, 或者tail还没有到队尾
				if next == nil { // 为空
					return t
				}
				// tail往队尾移
				cas[T](&f.tail, tail, next)
			} else {
				v := next.value
				if cas[T](&f.head, head, next) {
					return v
				}
			}
		}
	}
}

// load 读取节点的值
func load[T any](p *unsafe.Pointer) (n *Node[T]) {
	return (*Node[T])(atomic.LoadPointer(p))
}

func cas[T any](p *unsafe.Pointer, old, new *Node[T]) (ok bool) {
	return atomic.CompareAndSwapPointer(
		p, unsafe.Pointer(old), unsafe.Pointer(new))
}
