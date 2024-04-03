package containerx

import "container/list"

/*
List和Element，List 实现了一个双向链表，而 Element 则代表了链表中元素的结构。
*/

// 过度使用范型, 大可不必

type kv[K comparable, T any] struct {
	key   K
	value T
}

// LRUCache 基于container/list实现的LRU缓存
type LRUCache[K comparable, T any] struct {
	limit   int                 //缓存数量
	evicts  *list.List          //双向链表用于淘汰数据
	cache   map[K]*list.Element //记录缓存数据
	onEvict func(key K)         //淘汰数据时回调
}

func NewLRUCache[K comparable, T any](limit int, onEvict func(key K)) *LRUCache[K, T] {
	return &LRUCache[K, T]{
		limit:   limit,
		evicts:  list.New(),
		cache:   make(map[K]*list.Element),
		onEvict: onEvict,
	}
}

func (lru *LRUCache[K, T]) Get(key K) T {
	if ele, ok := lru.cache[key]; ok {
		lru.evicts.MoveToFront(ele)
		return ele.Value.(kv[K, T]).value
	}
	// 这里返回一个空值，因为T是泛型，只能用new来初始化
	// 为啥不可以直接返回T{}
	return *new(T)
}

func (lru *LRUCache[K, T]) Put(key K, value T) {
	if ele, ok := lru.cache[key]; ok {
		lru.evicts.MoveToFront(ele)
		ele.Value = kv[K, T]{key: key, value: value}
		return
	}

	for lru.evicts.Len() >= lru.limit {
		ele := lru.evicts.Back() // 淘汰末尾节点
		if ele != nil {
			lru.evicts.Remove(ele)
			delete(lru.cache, ele.Value.(kv[K, T]).key)
		}
	}

	ele := lru.evicts.PushFront(&kv[K, T]{key: key, value: value})
	lru.cache[key] = ele
}
