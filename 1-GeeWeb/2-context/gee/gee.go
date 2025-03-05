package gee

type HandlerFunc func(*Context)

type Engine struct {
	router *router
}
