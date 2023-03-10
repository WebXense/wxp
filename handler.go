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
	Handler  Handler[T]
}

type Handler[T any] func() *HandlerResponse[T]

func RegisterHandler[T any](handler Handler[T], middleware ...gin.HandlerFunc) {
	setting := handler()

	var requestObj interface{} = new(T)
	var requestName string = reflect.TypeOf(requestObj).Elem().Name()
	if requestName == "" {
		requestObj = nil
	} else {
		requestObj = reflect.New(reflect.TypeOf(requestObj).Elem()).Interface()
	}
	registerApi(setting.Method, setting.Route, requestObj, setting.Response, setting.Handler)

	ginHandler := func(ctx *gin.Context) {
		var err Error
		data, err := setting.Service(&Request[T]{
			Ctx:    ctx,
			Object: ginger.Request[T](ctx),
			Page:   ginger.PaginationRequest(ctx),
			Sort:   ginger.SortRequest(ctx),
		})
		if err != nil {
			Server.ERR(ctx, err)
			return
		}
		Server.OK(ctx, data, nil)
	}
	if middleware == nil {
		Server.Engine.Handle(setting.Method, setting.Route, ginHandler)
		return
	} else {
		middleware = append(middleware, ginHandler)
		Server.Engine.Handle(setting.Method, setting.Route, middleware...)
	}
}
