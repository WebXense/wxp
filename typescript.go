package wxp

import (
	"os"

	"github.com/WebXense/wxp/tf"
)

type api struct {
	method   string
	route    string
	request  interface{}
	response interface{}
}

var apis = make(map[string]api)

func registerApi(method string, route string, request interface{}, response interface{}) {
	apis[route] = api{
		method:   method,
		route:    route,
		request:  request,
		response: response,
	}
}

func generateTypeScript() {
	converter := tf.New()

	for _, api := range apis {
		converter.Add(api.request)
		converter.Add(api.response)
	}

	os.RemoveAll("api")
	err := os.Mkdir("api", os.ModeAppend)
	if err != nil {
		panic(err)
	}
	err = os.WriteFile("api/model.ts", []byte(converter.ToString()), os.ModeAppend)
	if err != nil {
		panic(err)
	}
}
