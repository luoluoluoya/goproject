package cache

import (
	"fmt"
	"testing"
)

type String string

func (s String) Len() int {
	return len(s)
}

func TestCache_Set(t *testing.T) {
	c := NewCache(100, nil)

	key, value := "key1", String("hello world")
	if !c.Set(key, value) {
		t.Fatalf("new item cacahe set return false\n")
	}

	if v, ok := c.Get(key); !ok || v.(String) != value {
		t.Fatalf("cache hit key1=‘%s’ failed", string(value))
	}
	if _, ok := c.Get("key2"); ok {
		t.Fatalf("cache miss key2 failed")
	}

	if c.nBytes != int64(len(key)+value.Len()) {
		t.Fatalf("cache used size error\n")
	}

	value = "hello go"
	c.Set(key, value)
	if v, ok := c.Get(key); !ok || v.(String) != value {
		fmt.Println(v)
		t.Fatalf("cache hit key1=‘%s’ failed", string(value))
	}
	if c.nBytes != int64(len(key)+value.Len()) {
		t.Fatalf("cache used size error\n")
	}
}

func TestCache_Del(t *testing.T) {
	c := NewCache(100, nil)
	key1, value1 := "key1", String("hello world")
	key2, value2 := "key2", String("hello go")
	c.Set(key1, value1)
	c.Set(key2, value2)
	if v, ok := c.Get(key1); !ok || v.(String) != value1 {
		t.Fatalf("cache hit key1=‘%s’ failed", string(value1))
	}
	if v, ok := c.Get(key2); !ok || v.(String) != value2 {
		t.Fatalf("cache hit key2=‘%s’ failed", string(value2))
	}

	if c.nBytes != int64(len(key1)+len(key2)+value1.Len()+value2.Len()) {
		t.Fatalf("cache used size error\n")
	}

	value := c.Del("key3")
	if value != nil {
		t.Fatalf("delete uncached value is not nil\n")
	}

	value = c.Del(key1)
	if value == nil {
		t.Fatalf("delete %s, value is not %s\n", key1, value1)
	}
	if _, ok := c.Get(key1); ok {
		t.Fatalf("cache hit key1=‘%s’ success", string(value1))
	}
	if c.nBytes != int64(len(key2)+value2.Len()) {
		t.Fatalf("cache used size error\n")
	}
}

func TestCache_Capacity(t *testing.T) {
	c := NewCache(20, nil)
	key1, value1 := "key1", String("hello world")
	key2, value2 := "key2", String("hello go")
	c.Set(key1, value1)
	c.Set(key2, value2)
	if _, ok := c.Get(key1); ok {
		t.Fatalf("cache hit key1=‘%s’", string(value1))
	}
	if c.nBytes != int64(len(key2)+value2.Len()) {
		t.Fatalf("cache used size error\n")
	}
}
