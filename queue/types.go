package queue

type Queue[T any] interface {
	// Enqueue 将元素放入队列，如果此时队列已经满了，那么返回错误
	Enqueue(t T) error
	// Dequeue 从队首获得一个元素
	// 如果此时队列里面没有元素，那么返回错误
	Dequeue() (T, error)
}
