package tools

import (
	"github.com/sari3l/requests"
	"reflect"
)

func HookResponseGbkToUtf8(response any) (error, any) {
	resp := response.(requests.Response)
	respRef := reflect.ValueOf(&resp).Elem()
	respContent := respRef.FieldByName("Content")
	chineseContent := ConvertGbkToUtf8(respContent.String())
	respContent.SetString(chineseContent)
	return nil, resp
}

func HookResponseUtf8ToGbk(response any) (error, any) {
	resp := response.(requests.Response)
	respRef := reflect.ValueOf(&resp).Elem()
	respContent := respRef.FieldByName("Content")
	chineseContent := ConvertUtf8ToGbk(respContent.String())
	respContent.SetString(chineseContent)
	return nil, resp
}
