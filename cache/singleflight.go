package cache

import "sync"

type call struct {
	err error          // 返回错误
	val interface{}    // 返回值
	wg  sync.WaitGroup // 请求组
}

type ReqGroup struct {
	mu sync.Mutex
	m  map[string]*call
}

// 确保并发情况下函数 fn 只被执行一次
func (g *ReqGroup) Do(key string, fn func() (interface{}, error)) (interface{}, error) {
	g.mu.Lock()
	if g.m == nil {
		g.m = make(map[string]*call)
	}
	// 请求中，等待响应
	if c, ok := g.m[key]; ok {
		g.mu.Unlock()
		c.wg.Wait()
		return c.val, c.err
	}
	// 准备请求
	c := new(call)
	c.wg.Add(1)
	g.m[key] = c
	g.mu.Unlock()

	c.val, c.err = fn()
	c.wg.Done()

	g.mu.Lock()
	delete(g.m, key)
	g.mu.Unlock()
	return c.val, c.err
}
