package frame

import (
	"testing"
)

func newTestRouterGroup() *router {
	r := New()
	gv1 := r.Group("/v1")
	gv1.GET("/user/:id", func(c *Context) { // GET /v1/user/:id
	})
	gv1.GET("/users/list", func(c *Context) { // GET /v1/user/list
	})

	gv2 := r.Group("/v2")
	gv2.GET("/user/:userId", func(c *Context) { // GET /v2/user/L:id
	})
	gv2.GET("/users/list", func(c *Context) { // GET /v2/user/list
	})

	gv3 := gv1.Group("/v3")
	gv3.GET("/users/list", func(c *Context) { // GET /v1/v3/user/list
	})
	return r.router
}

func TestRouterGroup_Group(t *testing.T) {
	r := newTestRouterGroup()
	// not found
	n, ps := r.getRoute("GET", "/user/10")
	if n != nil || ps != nil {
		t.Fatal("nil shouldn be returned")
	}

	// gv1
	n, ps = r.getRoute("GET", "/v1/user/10")
	if n == nil {
		t.Fatal("nil shouldn't be returned")
	}
	if v, ok := ps["id"]; v != "10" || !ok {
		t.Fatal("ps['id'] shouldn't be 10")
	}

	// gv3
	n, ps = r.getRoute("GET", "/v1/v3/users/list")
	if n == nil {
		t.Fatal("nil shouldn't be returned")
	}

	//gv2
	n, ps = r.getRoute("GET", "/v2/user/10")
	if n == nil {
		t.Fatal("nil shouldn't be returned")
	}
	if v, ok := ps["userId"]; v != "10" || !ok {
		t.Fatal("ps['id'] shouldn't be 10")
	}
}
