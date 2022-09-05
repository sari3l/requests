package requests

import (
	"github.com/sari3l/requests/types"
)

func Get(url string, ext ...types.Ext) *Response {
	return initRequest("GET", url, &ext).request()
}

func Post(url string, ext ...types.Ext) *Response {
	return initRequest("POST", url, &ext).request()
}

func Put(url string, ext ...types.Ext) *Response {
	return initRequest("PUT", url, &ext).request()
}

func Delete(url string, ext ...types.Ext) *Response {
	return initRequest("DELETE", url, &ext).request()
}

func Head(url string, ext ...types.Ext) *Response {
	return initRequest("HEAD", url, &ext).request()
}

func Options(url string, ext ...types.Ext) *Response {
	return initRequest("OPTIONS", url, &ext).request()
}

func initRequest(method string, url string, exts *[]types.Ext) *Session {
	s := HTMLSession(5, "", true, true)
	return s.init(method, url, exts)
}
