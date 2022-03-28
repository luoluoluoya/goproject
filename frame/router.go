package frame

import (
	"log"
	"net/http"
	"strings"
)

type router struct {
	roots    map[string]*node
	handlers map[string]HandlerFunc
}

func newRouter() *router {
	return &router{roots: make(map[string]*node), handlers: make(map[string]HandlerFunc)}
}

func parsePattern(pattern string) (parts []string) {
	items := strings.Split(pattern, "/")
	for _, item := range items {
		if item == "" {
			continue
		}
		parts = append(parts, item)
		if item[0] == '*' {
			break
		}
	}
	return
}

func (r *router) addRoute(method, pattern string, handler HandlerFunc) {
	method = strings.ToUpper(method)
	log.Printf("Route: %4s\t%s\n", method, pattern)
	if _, ok := r.roots[method]; !ok {
		r.roots[method] = &node{}
	}
	parts := parsePattern(pattern)
	r.roots[method].insert(pattern, parts, 0)
	key := method + "-" + pattern
	r.handlers[key] = handler
}

func (r *router) getRoute(method, pattern string) (*node, map[string]string) {
	method = strings.ToUpper(method)
	root, ok := r.roots[method]
	if !ok {
		return nil, nil
	}
	searchParts := parsePattern(pattern)
	n := root.search(searchParts, 0)
	if n == nil {
		return nil, nil
	}
	params := make(map[string]string)
	parts := parsePattern(n.pattern)
	for index, part := range parts {
		if part[0] == ':' {
			params[part[1:]] = searchParts[index]
		}
		if part[0] == '*' && len(part) > 1 {
			params[part[1:]] = strings.Join(searchParts[index:], "/")
			break
		}
	}
	return n, params
}

func (r *router) handle(c *Context) {
	node, params := r.getRoute(strings.ToUpper(c.Method), c.Path)
	var handler HandlerFunc
	if node == nil {
		c.handlers = nil
		handler = func(c *Context) {
			c.String(http.StatusNotFound, "404 NOT FOUND: %s\n", c.Path)
		}
	} else {
		key := strings.ToUpper(c.Method) + "-" + node.pattern
		handler = r.handlers[key]
	}
	c.Params = params
	c.handlers = append(c.handlers, handler)
	c.Next()
}
