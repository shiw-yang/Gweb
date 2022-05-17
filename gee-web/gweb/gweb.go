package gweb

import (
	"log"
	"net/http"
)

// HandlerFunc defines the request handler used by gweb
type HandlerFunc func(http.ResponseWriter, *http.Request)

// Engine implement the interface of ServeHTTP
type Engine struct {
	router map[string]HandlerFunc
}

// New is the constructor of gweb.Engine 构造函数
func New() *Engine {
	return &Engine{router: make(map[string]HandlerFunc)}
}

// addRoute is setting router
func (engine Engine) addRoute(method, pattern string, handler HandlerFunc) {
	key := method + "-" + pattern
	log.Printf("Route %4s - %s", method, pattern)
	engine.router[key] = handler
}

// GET defines the method to add GET request
func (engine Engine) GET(pattern string, handler HandlerFunc) {
	engine.addRoute("GET", pattern, handler)
}

// POST defines the method to add POST request
func (engine Engine) POST(pattern string, handler HandlerFunc) {
	engine.addRoute("POST", pattern, handler)
}

// Run defines the method to start a http server
func (engine Engine) Run(addr string) (err error) {
	return http.ListenAndServe(addr, engine)
}

func (engine Engine) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	key := req.Method + "-" + req.URL.Path
	if handler, ok := engine.router[key]; ok {
		handler(w, req)
	} else {
		w.WriteHeader(http.StatusNotFound)
	}
}
