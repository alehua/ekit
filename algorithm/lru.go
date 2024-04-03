package algorithm

import "container/list"

type keyLru struct {
	limit    int                      //缓存数量
	evicts   *list.List               //双向链表用于淘汰数据
	elements map[string]*list.Element //记录缓存数据
}
