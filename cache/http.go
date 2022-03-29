package cache

import (
	"fmt"
	"net/http"
	"strings"
)

const defaultPath = "/cache"

// 缓存网络服务
type Pool struct {
	host string
	path string
}

func NewPool(host string) *Pool {
	return &Pool{
		host: host,
		path: defaultPath,
	}
}

func (p *Pool) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if !strings.HasPrefix(r.URL.Path, p.path) {
		http.Error(w, "HTTPPool serving unexpected path: "+r.URL.Path, http.StatusBadRequest)
		return
	}
	fmt.Printf("%s %s\n", r.Method, r.URL.Path)
	fmt.Println(r.Header)
	parts := strings.SplitN(r.URL.Path[len(p.path)+1:], "/", 2)
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
