package wxp

import (
	"reflect"

	"github.com/WebXense/ginger/ginger"
	"github.com/gin-gonic/gin"
)

type HandlerResponse[T any] struct {
	Route    string
	Method   string
	Response interface{}
	Service  Service[T]
}

type Handler[T any] func() *HandlerResponse[T]

func RegisterHandler[T any](handler Handler[T], middleware ...gin.HandlerFunc) {
	setting := handler()

	var requestObj interface{} = new(T)
	if reflect.TypeOf(requestObj).Name() == "struct{}" {
		requestObj = nil
	}
	registerApi(setting.Method, setting.Route, requestObj, setting.Response, setting.Service)

	ginHandler := func(ctx *gin.Context) {
		var err Error
		data, err := setting.Service(&Request[T]{
			Ctx:    ctx,
			Object: ginger.Request[T](ctx),
			Page:   ginger.PaginationRequest(ctx),
			Sort:   ginger.SortRequest(ctx),
		})
		if err != nil {
			ERR(ctx, err)
			return
		}
		OK(ctx, data, nil)
	}
	if middleware == nil {
		Engine.Handle(setting.Method, setting.Route, ginHandler)
		return
	} else {
		middleware = append(middleware, ginHandler)
		Engine.Handle(setting.Method, setting.Route, middleware...)
	}
}
