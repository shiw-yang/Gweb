package gweb

import (
	"log"
	"net/http"
)

// HandlerFunc defines the request handler used by gweb
type HandlerFunc func(c *Context)

// Engine implement the interface of ServeHTTP
type Engine struct {
	router       *router
	*RouterGroup                // make engine has RouterGroup's power
	groups       []*RouterGroup // store all groups
}

// RouterGroup packets into the Engine
type RouterGroup struct {
	prefix      string
	middlewares []HandlerFunc // support middleware
	parent      *RouterGroup  // support nesting
	engine      *Engine       // all groups share an Engine instance
}

// New is the constructor of gweb.Engine
func New() *Engine {
	engine := &Engine{router: newRouter()}
	engine.RouterGroup = &RouterGroup{engine: engine}
	engine.groups = []*RouterGroup{engine.RouterGroup}
	return engine
}

func (group *RouterGroup) Group(prefix string) *RouterGroup {
	engine := group.engine
	newGroup := &RouterGroup{
		prefix: engine.prefix + prefix,
		parent: group,
		engine: engine,
	}
	engine.groups = append(engine.groups, newGroup)
	return newGroup
}

func (group *RouterGroup) addRoute(method, comp string, handler HandlerFunc) {
	pattern := group.prefix + comp
	log.Printf("Route %4s - %s", method, pattern)
	group.engine.router.addRoute(method, pattern, handler)
}

// GET defines the method to add GET request
func (group *RouterGroup) GET(pattern string, handler HandlerFunc) {
	group.engine.addRoute("GET", pattern, handler)
}

// POST defines the method to add POST request
func (group *RouterGroup) POST(pattern string, handler HandlerFunc) {
	group.engine.addRoute("POST", pattern, handler)
}

// Run defines the method to start a http server
func (engine Engine) Run(addr string) (err error) {
	return http.ListenAndServe(addr, engine)
}

func (engine Engine) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	c := newContext(w, req)
	engine.router.handle(c)
}
