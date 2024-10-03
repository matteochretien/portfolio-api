package http

import "context"

type HttpContext interface {
	context.Context
	Param(key string) string
	JSON(code int, obj interface{})
	Next()
	Abort()
	AbortWithStatus(Code int)
	SetHeader(key string, value string)
	Status(code int)
	GetMethod() string
}
