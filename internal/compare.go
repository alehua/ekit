package internal

// Comparator 用于比较两个对象的大小 src < dst, 返回-1，src = dst, 返回0，src > dst, 返回1
type Comparator[T any] func(src T, dst T) int
