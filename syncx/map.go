package syncx

import "sync"

// Map 对sync.Map 的泛型封装
type Map[K comparable, V any] struct {
	m sync.Map
}

// Load 加载键值对
func (m *Map[K, V]) Load(key K) (value V, ok bool) {
	var anyVal any
	anyVal, ok = m.m.Load(key)
	if anyVal != nil {
		value = anyVal.(V)
	}
	return
}

// Store 存储键值对
func (m *Map[K, V]) Store(key K, value V) {
	m.m.Store(key, value)
}

// LoadOrStore 加载或者存储一个键值对
func (m *Map[K, V]) LoadOrStore(key K, value V) (actual V, loaded bool) {
	var anyVal any
	anyVal, loaded = m.m.LoadOrStore(key, value)
	if anyVal != nil {
		actual = anyVal.(V)
	}
	return
}

// LoadAndDelete 加载并且删除一个键值对
func (m *Map[K, V]) LoadAndDelete(key K) (value V, loaded bool) {
	var anyVal any
	anyVal, loaded = m.m.LoadAndDelete(key)
	if anyVal != nil {
		value = anyVal.(V)
	}
	return
}

// Delete 删除键值对
func (m *Map[K, V]) Delete(key K) {
	m.m.Delete(key)
}

// Range 遍历, f 不能为 nil
// 传入 f 的时候，K 和 V 直接使用对应的类型，如果 f 返回 false，那么就会中断遍历
func (m *Map[K, V]) Range(f func(key K, value V) bool) {
	m.m.Range(func(key, value any) bool {
		var (
			k K
			v V
		)
		if value != nil {
			v = value.(V)
		}
		if key != nil {
			k = key.(K)
		}
		return f(k, v)
	})
}
