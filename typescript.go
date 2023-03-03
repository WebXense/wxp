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
		if api.request == nil {
			models[modelName(api.request)] = api.request
		}
		if api.response == nil {
			models[modelName(api.response)] = api.response
		}
	}

	for _, model := range models {
		converter.Add(model)
	}

	converter.WithInterface(true).WithPrefix(`
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
`).ConvertToFile("api.ts")
}

func modelName(model interface{}) string {
	return reflect.ValueOf(model).Type().Name()
}
