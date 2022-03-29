package cache

import (
	"hash/crc32"
	"sort"
	"strconv"
)

type Hash func([]byte) uint32

type Map struct {
	hash     Hash           // hash函数
	replicas int            // 虚拟节点倍数
	keys     []int          // hash环
	hashMap  map[int]string // 虚拟节点和真实节点映射
}

func NewMap(replicas int, hash Hash) *Map {
	if hash == nil {
		hash = crc32.ChecksumIEEE
	}
	return &Map{
		hash:     hash,
		replicas: replicas,
		hashMap:  make(map[int]string),
	}
}

func (m *Map) Add(keys ...string) {
	for _, key := range keys {
		for i := 0; i < m.replicas; i++ {
			h := int(m.hash([]byte(strconv.Itoa(i) + key)))
			m.keys = append(m.keys, h)
			m.hashMap[h] = key
		}
	}
	sort.Ints(m.keys)
}

func (m *Map) Get(key string) string {
	if len(m.keys) == 0 {
		return ""
	}
	hash := int(m.hash([]byte(key)))
	idx := sort.Search(len(m.keys), func(i int) bool {
		return m.keys[i] >= hash
	})
	return m.hashMap[m.keys[idx%len(m.keys)]]
}
