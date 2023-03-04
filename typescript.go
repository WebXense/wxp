package wxp

import (
	"reflect"
	"time"

	tf "github.com/tkrajina/typescriptify-golang-structs/typescriptify"
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
	converter.ManageType(time.Time{}, tf.TypeOptions{TSType: "Date", TSTransform: "new Date(__VALUE__)"})
	converter.BackupDir = "" // don't backup

	models := make(map[string]interface{})
	for _, api := range apis {
		if api.request != nil {
			models[modelName(api.request)] = api.request
		}
		if api.response != nil {
			models[modelName(api.response)] = api.response
		}
	}

	LogDebug("models: %v", models)

	for _, model := range models {
		converter.Add(model)
	}

	converter.WithInterface(true).ConvertToFile("api.ts")
}

func modelName(model interface{}) string {
	return reflect.ValueOf(model).String()
}
