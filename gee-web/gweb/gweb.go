package gweb

import (
	"log"
	"net/http"
	"strings"
)

// HandlerFunc defines the request handler used by gweb
type HandlerFunc func(c *Context)

type (

	// RouterGroup packets into the Engine
	RouterGroup struct {
		prefix      string
		middlewares []HandlerFunc // support middleware
		engine      *Engine       // all groups share an Engine instance
	}

	// Engine implement the interface of ServeHTTP
	Engine struct {
		*RouterGroup // Embed type in go: make engine has RouterGroup's power
		router       *router
		groups       []*RouterGroup // store all groups
	}
)

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
	group.addRoute("GET", pattern, handler)
}

// POST defines the method to add POST request
func (group *RouterGroup) POST(pattern string, handler HandlerFunc) {
	group.addRoute("POST", pattern, handler)
}

// Run defines the method to start a http server
func (engine Engine) Run(addr string) (err error) {
	return http.ListenAndServe(addr, engine)
}

// Use is defined to add middleware to the group
func (group *RouterGroup) Use(middlewares ...HandlerFunc) {
	group.middlewares = append(group.middlewares, middlewares...)
}

// engine as the http/net is Handler, you can enter http.ListenAndServe
// you can see as an interface, must implement ServerHttp func
// copy from http.ListenAndServe :
//		The handler is typically nil, in which case the DefaultServeMux is used.
// copy from Handler interface:
/*	A Handler responds to an HTTP request.
 *
 *	ServeHTTP should write reply headers and data to the ResponseWriter
 *	and then return. Returning signals that the request is finished; it
 *	is not valid to use the ResponseWriter or read from the
 *	Request.Body after or concurrently with the completion of the
 *	ServeHTTP call.
 */
func (engine Engine) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	var middlewares []HandlerFunc
	for _, group := range engine.groups {
		if strings.HasPrefix(req.URL.Path, group.prefix) {
			middlewares = append(middlewares, group.middlewares...)
		}
	}
	c := newContext(w, req)
	c.handlers = middlewares
	engine.router.handle(c)
}
