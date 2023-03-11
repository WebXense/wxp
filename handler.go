package wxp

import (
	"reflect"

	"github.com/WebXense/ginger/ginger"
	"github.com/WebXense/wxp/errs"
	"github.com/WebXense/wxp/server"
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
	var requestName string = reflect.TypeOf(requestObj).Elem().Name()
	if requestName == "" {
		requestObj = nil
	} else {
		requestObj = reflect.New(reflect.TypeOf(requestObj).Elem()).Interface()
	}
	registerApi(setting.Method, setting.Route, requestObj, setting.Response, handler)

	ginHandler := func(ctx *gin.Context) {
		var err errs.Error
		data, err := setting.Service(&Request[T]{
			Ctx:    ctx,
			Object: ginger.Request[T](ctx),
			Page:   ginger.PaginationRequest(ctx),
			Sort:   ginger.SortRequest(ctx),
		})
		if err != nil {
			server.ERR(ctx, err)
			return
		}
		server.OK(ctx, data, nil)
	}
	if middleware == nil {
		server.Engine.Handle(setting.Method, setting.Route, ginHandler)
		return
	} else {
		middleware = append(middleware, ginHandler)
		server.Engine.Handle(setting.Method, setting.Route, middleware...)
	}
}
