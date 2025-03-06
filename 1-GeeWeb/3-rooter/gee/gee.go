package gee

import "net/http"

type HandlerFunc func(*Context)

type Engine struct {
	router *router
}

func New() *Engine {
	return &Engine{newRouter()}
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
	c := newContext(w, req)
	engine.router.handle(c)
}
