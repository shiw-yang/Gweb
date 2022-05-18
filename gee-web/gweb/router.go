package gweb

import "net/http"

type router struct {
	handlers map[string]HandlerFunc
}

// newRouter is an init func
func newRouter() *router {
	return &router{handlers: make(map[string]HandlerFunc)}
}

// addRoute is setting router
func (r *router) addRoute(method, pattern string, handler HandlerFunc) {
	key := method + "-" + pattern
	r.handlers[key] = handler
}

func (r *router) handle(c *Context) {
	key := c.Method + "-" + c.Path
	if handler, ok := r.handlers[key]; ok {
		handler(c)
	} else {
		c.String(http.StatusNotFound, "404 NOT FOUND: %s\n", c.Path)
	}
}
