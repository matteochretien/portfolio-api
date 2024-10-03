package http

type HttpRouter interface {
	GET(path string, handlers ...HttpHandler)
	POST(path string, handlers ...HttpHandler)
	PUT(path string, handlers ...HttpHandler)
	DELETE(path string, handlers ...HttpHandler)
	Handle(method, path string, handlers ...HttpHandler)
	Group(path string, handlers ...HttpHandler) HttpRouter
	Use(handlers ...HttpHandler)
}
