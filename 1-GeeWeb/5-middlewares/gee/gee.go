package gee

import (
	"log"
	"net/http"
	"strings"
)

type HandlerFunc func(*Context)

type Engine struct {
	router *router
	*RouterGroup
	groups []*RouterGroup //store all groups
}

type RouterGroup struct {
	prefix      string
	middlewares []HandlerFunc //support middleware
	parent      *RouterGroup  //support nesting
	engine      *Engine       //all groups share an Engine instance
}

func New() *Engine {
	engine := &Engine{router: newRouter()}
	engine.RouterGroup = &RouterGroup{engine: engine}
	engine.groups = []*RouterGroup{engine.RouterGroup}
	return engine
}

func (engine *Engine) GET(pattern string, handlerFunc HandlerFunc) {
	engine.router.addRoute("GET", pattern, handlerFunc)
}

func (engine *Engine) POST(pattern string, handlerFunc HandlerFunc) {
	engine.router.addRoute("POST", pattern, handlerFunc)
}

func (engine *Engine) Run(port string) (err error) {
	return http.ListenAndServe(port, engine)
}

func (engine *Engine) ServeHTTP(w http.ResponseWriter, req *http.Request) {
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

func (group *RouterGroup) Group(prefix string) *RouterGroup {
	engine := group.engine
	newGroup := &RouterGroup{
		prefix: group.prefix + prefix,
		parent: group,
		engine: engine,
	}
	engine.groups = append(engine.groups, newGroup)
	return newGroup
}

func (group *RouterGroup) addRoute(method string, comp string, handleFunc HandlerFunc) {
	pattern := group.prefix + comp
	log.Printf("Route %4s - %s", method, pattern)
	group.engine.router.addRoute(method, pattern, handleFunc)
}

func (group *RouterGroup) GET(pattern string, handleFunc HandlerFunc) {
	group.addRoute("GET", pattern, handleFunc)
}

func (group *RouterGroup) POST(pattern string, handleFunc HandlerFunc) {
	group.addRoute("POST", pattern, handleFunc)
}

func (group *RouterGroup) Use(middlewares ...HandlerFunc) {
	group.middlewares = append(group.middlewares, middlewares...)
}
