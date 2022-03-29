package cache

import (
	"fmt"
	"log"
	"sync"
)

// 分组缓存
type Group struct {
	identify  string
	getter    Getter
	mainCache *cache
	peers     PeerPicker
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

// 注册节点
func (g *Group) RegisterPeers(peers PeerPicker) {
	if g.peers != nil {
		panic("RegisterPeerPicker called more than once")
	}
	g.peers = peers
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

// key不存在于本地缓存，从其他地方加载
func (g *Group) load(key string) (bytes ByteView, err error) {
	if g.peers != nil {
		bytes, err = g.loadPeer(key)
	}
	if g.peers == nil || err != nil {
		bytes, err = g.loadLocally(key)
	}
	if err != nil {
		return ByteView{}, err
	}
	g.mainCache.Set(key, bytes)
	return bytes, nil
}

// 本地加载
func (g *Group) loadLocally(key string) (ByteView, error) {
	v, err := g.getter.Get(key)
	if err != nil {
		return ByteView{}, err
	}
	return ByteView{data: cloneBytes(v)}, nil
}

// 节点加载
func (g *Group) loadPeer(key string) (ByteView, error) {
	getter, ok := g.peers.PickPeer(key)
	if !ok {
		return ByteView{}, fmt.Errorf("peer of key [%s] noy found", key)
	}
	data, err := getter.Get(g.identify, key)
	if err != nil {
		log.Println("[GeeCache] Failed to get from peer", err)
	}
	return ByteView{data: cloneBytes(data)}, err
}
