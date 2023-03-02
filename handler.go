package wxp

import (
	"github.com/WebXense/ginger/ginger"
	"github.com/gin-gonic/gin"
)

func GET[T any](path string, service Service[T], response interface{}, middleware ...gin.HandlerFunc) {
	handler := func(ctx *gin.Context) {
		var err Error
		data, err := service(ctx, ginger.Request[T](ctx))
		if err != nil {
			ERR(ctx, err)
			return
		}
		OK(ctx, data, nil)
	}
	if middleware == nil {
		Engine.GET(path, handler)
		return
	} else {
		middleware = append(middleware, handler)
		Engine.GET(path, middleware...)
	}
}

func POST[T any](path string, service Service[T], response interface{}, middleware ...gin.HandlerFunc) {
	handler := func(ctx *gin.Context) {
		var err Error
		data, err := service(ctx, ginger.Request[T](ctx))
		if err != nil {
			ERR(ctx, err)
			return
		}
		OK(ctx, data, nil)
	}
	if middleware == nil {
		Engine.POST(path, handler)
		return
	} else {
		middleware = append(middleware, handler)
		Engine.POST(path, middleware...)
	}
}

func PUT[T any](path string, service Service[T], response interface{}, middleware ...gin.HandlerFunc) {
	handler := func(ctx *gin.Context) {
		var err Error
		data, err := service(ctx, ginger.Request[T](ctx))
		if err != nil {
			ERR(ctx, err)
			return
		}
		OK(ctx, data, nil)
	}
	if middleware == nil {
		Engine.PUT(path, handler)
		return
	} else {
		middleware = append(middleware, handler)
		Engine.PUT(path, middleware...)
	}
}

func DELETE[T any](path string, service Service[T], response interface{}, middleware ...gin.HandlerFunc) {
	handler := func(ctx *gin.Context) {
		var err Error
		data, err := service(ctx, ginger.Request[T](ctx))
		if err != nil {
			ERR(ctx, err)
			return
		}
		OK(ctx, data, nil)
	}
	if middleware == nil {
		Engine.DELETE(path, handler)
		return
	} else {
		middleware = append(middleware, handler)
		Engine.DELETE(path, middleware...)
	}
}
