package cache

import (
	"fmt"
	"sync"
)

// 分组缓存
type Group struct {
	identify  string
	getter    Getter
	mainCache *cache
}

var (
	mu     sync.RWMutex
	groups = make(map[string]*Group)
)

// 新建缓存分组
func NewGroup(ident string, maxCache int64, getter Getter) *Group {
	if getter == nil {
		panic("getter is nil !")
	}
	mu.Lock()
	defer mu.Unlock()
	if _, ok := groups[ident]; ok {
		panic(fmt.Sprintf("identify %s already used !", ident))
	}
	group := &Group{
		identify:  ident,
		mainCache: &cache{cacheBytes: maxCache},
		getter:    getter,
	}
	groups[ident] = group
	return group
}

func GetGroup(ident string) *Group {
	mu.RLock()
	defer mu.RUnlock()
	return groups[ident]
}

func (g *Group) Get(key string) (ByteView, error) {
	if key == "" {
		return ByteView{}, fmt.Errorf("key is required")
	}
	if v, ok := g.mainCache.Get(key); ok {
		return v, nil
	}
	return g.load(key)
}

func (g *Group) load(key string) (ByteView, error) {
	v, err := g.getter.Get(key)
	if err != nil {
		return ByteView{}, err
	}
	b := ByteView{data: cloneBytes(v)}
	g.mainCache.Set(key, b)
	return b, nil
}
