package cache

import (
	"fmt"
	"io/ioutil"
	"net/http"
	url2 "net/url"
	"strings"
	"sync"
)

const (
	defaultPath     = "/cache/"
	defaultReplicas = 50
)

type httpGetter struct {
	url string
}

func (g *httpGetter) Get(group, key string) ([]byte, error) {
	url := fmt.Sprintf("%s%s/%s", g.url, url2.QueryEscape(group), url2.QueryEscape(key))
	res, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("server returend: %s", res.Status)
	}
	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("reading response body: %v", err)
	}
	return data, err
}

// 缓存网络服务
type Pool struct {
	host    string                 // 主机：127.0.0.1:8080
	path    string                 // 服务路径：/cache/
	mu      sync.Mutex             // 锁
	peers   *Map                   // 缓存节点
	getters map[string]*httpGetter // 节点对应的数据获取器
}

func NewPool(host string) *Pool {
	return &Pool{
		host: host,
		path: defaultPath,
	}
}

func (p *Pool) Set(peers ...string) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.peers = NewMap(defaultReplicas, nil)
	p.peers.Add(peers...)
	p.getters = make(map[string]*httpGetter, len(peers))
	for _, peer := range peers {
		p.getters[peer] = &httpGetter{url: peer + defaultPath}
	}
}

func (p *Pool) PickPeer(key string) (PeerGetter, bool) {
	p.mu.Lock()
	defer p.mu.Unlock()
	if peer := p.peers.Get(key); peer != "" && peer != p.host {
		fmt.Printf("Pick Peer: %s", peer)
		return p.getters[peer], true
	}
	return nil, false
}

func (p *Pool) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if !strings.HasPrefix(r.URL.Path, p.path) {
		http.Error(w, "HTTPPool serving unexpected path: "+r.URL.Path, http.StatusBadRequest)
		return
	}
	fmt.Printf("%s %s\n", r.Method, r.URL.Path)
	parts := strings.SplitN(r.URL.Path[len(p.path):], "/", 2)
	if len(parts) != 2 {
		http.Error(w, "request like: /cache/{cachegroup}/{key}", http.StatusBadRequest)
		return
	}
	groupName, key := parts[0], parts[1]
	group := GetGroup(groupName)
	if group == nil {
		http.Error(w, fmt.Sprintf("group %s not found", groupName), http.StatusBadRequest)
		return
	}
	data, err := group.Get(key)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/octet-stream")
	w.Write(data.ByteSlice())
}

var _ PeerPicker = (*Pool)(nil)
var _ PeerGetter = (*httpGetter)(nil)
