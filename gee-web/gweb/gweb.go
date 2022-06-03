package gweb

import (
	"html/template"
	"net/http"
	"strings"
)

// HandlerFunc defines the request handler used by gweb
type HandlerFunc func(c *Context)

// Engine implement the interface of ServeHTTP
type Engine struct {
	*RouterGroup  // Embed type in go: make engine has RouterGroup's power
	router        *router
	groups        []*RouterGroup     // store all groups
	htmlTemplates *template.Template // for html render
	funcMap       template.FuncMap   // for html render
}

// SetFuncMap is to create a func mapping within the tmpl
// SetFuncMap always runs before LoadHTMLGlob
//because the custom tmpl func needs to be prepared on FunMap before tmpl registration
func (engine *Engine) SetFuncMap(funcMap template.FuncMap) {
	engine.funcMap = funcMap
}

// LoadHTMLGlob is to mount the tmpl file under the pattern path to htmlTemplates
func (engine *Engine) LoadHTMLGlob(pattern string) {
	engine.htmlTemplates = template.Must(template.New("").Funcs(engine.funcMap).ParseGlob(pattern))
}

// New Default are the constructor of gweb.Engine
func New() *Engine {
	engine := &Engine{router: newRouter()}
	engine.RouterGroup = &RouterGroup{engine: engine}
	engine.groups = []*RouterGroup{engine.RouterGroup}
	return engine
}

func Default() *Engine {
	r := New()
	// global middleware
	r.Use(Logger())
	return r
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
func (engine *Engine) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	var middlewares []HandlerFunc
	for _, group := range engine.groups {
		if strings.HasPrefix(req.URL.Path, group.prefix) {
			middlewares = append(middlewares, group.middlewares...)
		}
	}
	c := newContext(w, req)
	c.handlers = middlewares
	c.engine = engine
	engine.router.handle(c)
}
