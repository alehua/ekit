package contanerx

// 对 "container/heap"的范型封装, 实现优先级队列
// 堆使用的数据结构是最小二叉树，即根节点比左边子树和右边子树的所有值都小

// Cmp function
type Cmp[T any] func(T, T) int

// HeapX to implements the interface of "heap.Interface"
type HeapX[T any] struct {
	Data []T
	Cmp  Cmp[T]
}

func (h *HeapX[T]) Len() int { return len(h.Data) }

func (h *HeapX[T]) Less(i, j int) bool {
	v := h.Cmp(h.Data[i], h.Data[j])
	return v < 0
}
func (h *HeapX[T]) Swap(i, j int) { h.Data[i], h.Data[j] = h.Data[j], h.Data[i] }

func (h *HeapX[T]) Push(x any) {
	// Push and Pop use pointer receivers because they modify the slice's length,
	// not just its contents.
	v := append(h.Data, x.(T))
	h.Data = v
}

func (h *HeapX[T]) Pop() any {
	old := h.Data
	n := len(old)
	x := old[n-1]
	h.Data = old[0 : n-1]
	return x
}
