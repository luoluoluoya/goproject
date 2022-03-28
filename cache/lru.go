package cache

import (
	"container/list"
)

type entry struct {
	key   string
	value Value
}

type Value interface {
	Len() int
}

type Cache struct {
	maxBytes  int64
	nBytes    int64
	ll        *list.List
	cache     map[string]*list.Element
	OnEvicted func(string, Value)
}

// 新建缓存
func NewCache(maxBytes int64, onEvicted func(string, Value)) *Cache {
	return &Cache{
		maxBytes:  maxBytes,
		ll:        list.New(),
		cache:     make(map[string]*list.Element),
		OnEvicted: onEvicted,
	}
}

// 新增缓存：已存在则替换
func (c *Cache) Set(key string, value Value) bool {
	elem, exists := c.cache[key]
	if exists { // 已存在，覆写
		c.ll.MoveToFront(elem)
		kv := elem.Value.(*entry)
		c.nBytes += int64(value.Len()) - int64(kv.value.Len())
		kv.value = value
	} else { // 不存在，新增
		elem := c.ll.PushFront(&entry{key: key, value: value})
		c.cache[key] = elem
		c.nBytes += int64(len(key)) + int64(value.Len())
	}
	c.shrinkageCapacity()
	return !exists
}

// 获取缓存
func (c *Cache) Get(key string) (Value, bool) {
	elem, ok := c.cache[key]
	if !ok {
		return nil, false
	}
	c.ll.MoveToFront(elem)
	return elem.Value.(*entry).value, true
}

// 删除缓存并返回旧值
func (c *Cache) Del(key string) Value {
	elem, ok := c.cache[key]
	if !ok {
		return nil
	}
	c.removeElem(elem)
	return elem.Value.(*entry).value
}

// 缓存元素数量
func (c *Cache) len() int {
	return c.ll.Len()
}

// 缩容：保证缓存数据不超过最大最大容量
func (c *Cache) shrinkageCapacity() {
	for c.maxBytes < c.nBytes {
		c.removeElem(c.ll.Back())
	}
}

// 移除元素并更改容量，调用OnEvicted函数
func (c *Cache) removeElem(elem *list.Element) {
	c.ll.Remove(elem)
	kv := elem.Value.(*entry)
	delete(c.cache, kv.key)
	c.nBytes -= int64(len(kv.key)) + int64(kv.value.Len())
	if c.OnEvicted != nil {
		c.OnEvicted(kv.key, kv.value)
	}
}
