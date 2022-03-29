package cache

// 节点选择器
type PeerPicker interface {
	PickPeer(key string) (PeerGetter, bool)
}

// 节点数据获取器
type PeerGetter interface {
	Get(group, key string) ([]byte, error)
}
