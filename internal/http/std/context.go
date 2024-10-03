package std

import (
	"context"
	"github.com/Niromash/niromash-api/internal/http"
	"github.com/Niromash/niromash-api/internal/logger"
	"github.com/goccy/go-json"
	http2 "net/http"
)

type StandardContext struct {
	context.Context
	request  *http2.Request
	response http2.ResponseWriter
	logg     logger.Logger
}

func NewStandardContext(req *http2.Request, resp http2.ResponseWriter) http.HttpContext {
	return &StandardContext{
		request:  req,
		response: resp,
	}
}

func (g *StandardContext) SetHeader(key string, value string) {
	g.request.Header.Set(key, value)
}

func (g *StandardContext) GetMethod() string {
	return g.request.Method
}

func (g *StandardContext) Param(key string) string {
	return g.request.PathValue(key)
}

func (g *StandardContext) JSON(code int, obj interface{}) {
	g.response.WriteHeader(code)

	err := json.NewEncoder(g.response).Encode(obj)
	if err != nil {
		g.logg.Error(err)
	}
}

func (g *StandardContext) Next() {
	return
}

func (g *StandardContext) Abort() {
	return
}

func (g *StandardContext) AbortWithStatus(Code int) {
	g.Abort()
	g.response.WriteHeader(Code)
}

func (g *StandardContext) Status(code int) {
	g.response.WriteHeader(code)
}
