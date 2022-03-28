package frame

import (
	"fmt"
	"time"
)

func Logger() HandlerFunc {
	return func(c *Context) {
		s := time.Now()
		c.Next()
		fmt.Printf("[%s] %5s %s %v\n", s.Format("2006-01-02 15:04:05"), c.Method, c.Request.RequestURI, time.Since(s))
	}
}
