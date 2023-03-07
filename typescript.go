package wxp

import (
	"os"

	"github.com/WebXense/wxp/api"
	"github.com/WebXense/wxp/tf"
)

var apis = make(map[string]api.Api)

func registerApi(method string, route string, request interface{}, response interface{}, handler interface{}) {
	_, ok := apis[method+":"+route]
	if ok {
		panic("api already registered: " + method + ":" + route)
	}
	apis[method+":"+route] = api.Api{
		Method:   method,
		Route:    route,
		Request:  request,
		Response: response,
		Handler:  handler,
	}
}

func generateTypeScript() {
	modelConverter := tf.New()
	apiConverter := api.New()

	for _, a := range apis {
		modelConverter.Add(a.Request)
		modelConverter.Add(a.Response)
		apiConverter.Add(a.Method, a.Route, a.Request, a.Response, a.Handler)
	}

	os.RemoveAll("api")
	err := os.Mkdir("api", os.ModePerm)
	if err != nil {
		panic(err)
	}
	err = os.WriteFile("api/model.ts", []byte(modelConverter.ToString()), os.ModePerm)
	if err != nil {
		panic(err)
	}
	err = os.WriteFile("api/api.ts", []byte(apiConverter.ToString()), os.ModePerm)
	if err != nil {
		panic(err)
	}
}
