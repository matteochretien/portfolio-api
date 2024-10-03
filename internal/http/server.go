package http

type HttpServer interface {
	HttpRouter
	Start() error
}

type HttpHandler func(c HttpContext)
