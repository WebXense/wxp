package api

import (
	"log"
	"reflect"
	"runtime"
	"strings"

	"github.com/iancoleman/strcase"
)

func New() *converter {
	return &converter{
		apis: make(map[string]Api),
	}
}

type Api struct {
	Method   string
	Route    string
	Request  interface{}
	Response interface{}
	Service  interface{}
}

type converter struct {
	apis map[string]Api
}

func (c *converter) Add(method string, route string, request interface{}, response interface{}, service interface{}) {
	c.apis[route] = Api{
		Method:   method,
		Route:    route,
		Request:  request,
		Response: response,
		Service:  service,
	}
}

func (c *converter) ToString() string {
	output := ""
	for _, api := range c.apis {
		output += c.convertToApi(api)
	}
	return prefix + output
}

func (c *converter) convertToApi(a Api) string {
	switch a.Method {
	case "GET":
		return c.convertToGet(a)
	case "POST":
		return c.convertToNonGet(a, "post")
	case "PUT":
		return c.convertToNonGet(a, "put")
	case "DELETE":
		return c.convertToNonGet(a, "del")
	default:
		log.Println("[WARNING] api: unknown method")
	}
	return ""
}

func (c *converter) convertToGet(a Api) string {
	if a.Request != nil {
		uriList := c.getUriList(a.Request)
		if len(uriList) > 0 {
			a.Route = c.replaceUri(a.Route, uriList)
		} else {
			a.Route += "\""
		}
	} else {
		a.Route += "\""
	}
	output := "export const " + c.nameOfFunc(a.Service) + " = async ("
	if a.Request != nil {
		output += "req: model." + c.nameOfModel(a.Request)
	}
	output += "): Promise<Response<"
	if a.Response != nil {
		output += "model." + c.nameOfModel(a.Response)
	} else {
		output += "null"
	}
	output += "> | null> => {\n"
	output += "    return get<"
	if a.Response != nil {
		output += "model." + c.nameOfModel(a.Response)
	} else {
		output += "null"
	}
	output += ">(\"" + a.Route

	if a.Request != nil {
		formList := c.getQueryList(a.Request)
		if len(formList) > 0 {
			output += ", [\n"
			for _, form := range formList {
				output += "        [\"" + form[0] + "\", req." + form[1] + "],\n"
			}
			output += "    ]"
		}
	}
	output += ")\n"
	output += "}\n\n"
	return output
}

func (c *converter) convertToNonGet(a Api, method string) string {
	if a.Request != nil {
		uriList := c.getUriList(a.Request)
		if len(uriList) > 0 {
			a.Route = c.replaceUri(a.Route, uriList)
		} else {
			a.Route += "\""
		}
	} else {
		a.Route += "\""
	}
	output := "export const " + c.nameOfFunc(a.Service) + " = async ("
	if a.Request != nil {
		output += "req: model." + c.nameOfModel(a.Request)
	}
	output += "): Promise<Response<"
	if a.Response != nil {
		output += "model." + c.nameOfModel(a.Response)
	} else {
		output += "null"
	}
	output += "> | null> => {\n"
	output += "    return " + method + "<"
	if a.Response != nil {
		output += "model." + c.nameOfModel(a.Response)
	} else {
		output += "null"
	}
	output += ">(\"" + a.Route

	if a.Request != nil {
		output += ", req"
	}
	output += ")\n"
	output += "}\n\n"
	return output
}

func (c *converter) nameOfModel(model interface{}) string {
	if reflect.TypeOf(model).Kind() == reflect.Ptr {
		model = reflect.ValueOf(model).Elem().Interface()
	}
	return reflect.TypeOf(model).Name()
}

func (c *converter) nameOfFunc(f interface{}) string {
	xs := strings.Split(runtime.FuncForPC(reflect.ValueOf(f).Pointer()).Name(), ".")
	name := strcase.ToLowerCamel(xs[len(xs)-1])
	if name[:2] == "Fm" {
		name = name[:len(name)-2]
	}
	return name
}

func (c *converter) getQueryList(model interface{}) [][]string {
	output := make([][]string, 0)
	if reflect.TypeOf(model).Kind() == reflect.Ptr {
		model = reflect.ValueOf(model).Elem().Interface()
	}
	if reflect.TypeOf(model).Kind() != reflect.Struct {
		log.Println("[WARNING] api: model must be a struct")
		return output
	}

	for i := 0; i < reflect.TypeOf(model).NumField(); i++ {
		tag := reflect.TypeOf(model).Field(i).Tag.Get("form")
		if tag != "" {
			output = append(output, []string{tag, tag})
		}
	}
	return output
}

func (c *converter) getUriList(model interface{}) [][]string {
	output := make([][]string, 0)
	if reflect.TypeOf(model).Kind() == reflect.Ptr {
		model = reflect.ValueOf(model).Elem().Interface()
	}
	if reflect.TypeOf(model).Kind() != reflect.Struct {
		log.Println("[WARNING] api: model must be a struct")
		return output
	}

	for i := 0; i < reflect.TypeOf(model).NumField(); i++ {
		tag := reflect.TypeOf(model).Field(i).Tag.Get("uri")
		if tag != "" {
			output = append(output, []string{tag, tag})
		}
	}
	return output
}

func (c *converter) replaceUri(route string, uriList [][]string) string {
	for j, uri := range uriList {
		xs := strings.Split(route, "/")

		temp := "/"
		for i, x := range xs {
			if x == "" {
				continue
			}
			if x == ":"+uri[0] {
				temp += "\" + req." + uri[1]
				if i != len(xs)-1 {
					temp += " + \"/"
				}
			} else {
				temp += x
				if i != len(xs)-1 {
					temp += "/"
				} else if temp[len(temp)-1] != '"' && j == len(uriList)-1 {
					temp += "\""
				}
			}
		}
		route = temp
	}
	return route
}

const prefix = `
import * as model from './model';

export const getAPIHost = (): string => {
    return process.env.NEXT_PUBLIC_API_HOST || 'http://127.0.0.1:5000';
}

export interface Response<T> {
    success: boolean;
    error?: Error;
    pagination?: Pagination;
    data?: T;
}

export interface Pagination {
    page: number;
    size: number;
    total: number;
}

export interface Error {
    code: string;
    message: string;
}

export const get = async <T>(url: string, params?: string[][]): Promise<Response<T> | null> => {
    try {
        url = getAPIHost() + url;
        if (params) {
            url += '?' + params.map(([key, value]) => key + "=" + value).join('&');
        }
        const response = await fetch(url, {
            method: 'GET',
            headers: {
                "Authorization": "Bearer " + localStorage.getItem("token"),
            },
        });
        return _handleResponse(response);
    } catch (err) {
        console.error(err);
        window.alert("Network error");
        return null;
    }
}

export const post = async <T>(url: string, body?: any): Promise<Response<T> | null> => {
    return await _nonGet('POST', url, body);
}

export const put = async <T>(url: string, body?: any): Promise<Response<T> | null> => {
    return await _nonGet('PUT', url, body);
}

export const del = async <T>(url: string, body?: any): Promise<Response<T> | null> => {
    return await _nonGet('DELETE', url, body);
}

export const upload = async <T>(url: string, file: File): Promise<Response<T> | null> => {
    try {
        url = getAPIHost() + url;
        const formData = new FormData();
        formData.append('file', file);
        const response = await fetch(url, {
            method: 'POST',
            headers: {
                "Authorization": "Bearer " + localStorage.getItem("token"),
            },
            body: formData,
        });
        return _handleResponse(response);
    } catch (err) {
        console.error(err);
        window.alert("Network error");
        return null;
    }
}

const _nonGet = async <T>(method: string, url: string, body?: any): Promise<Response<T> | null> => {
    try {
        url = getAPIHost() + url;
        const response = await fetch(url, {
            method: method,
            headers: {
                "Authorization": "Bearer " + localStorage.getItem("token"),
                "Content-Type": "application/json",
            },
            body: JSON.stringify(body),
        });
        return _handleResponse(response);
    } catch (err) {
        console.error(err);
        window.alert("Network error");
        return null;
    }
}

const _handleResponse = async <T>(resp: globalThis.Response): Promise<T> => {
    if (resp.status === 401) {
        window.location.href = '/';
        window.alert("Unauthorized");
    } else if (resp.status === 403) {
        window.alert("Forbidden");
    }
    return await resp.json();
}

`
