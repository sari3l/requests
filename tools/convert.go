package tools

import (
	"fmt"
	"github.com/sari3l/requests"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
	"io/ioutil"
	"reflect"
	"strings"
)

func ConvertGbkToUtf8(str string) string {
	reader := transform.NewReader(strings.NewReader(str), simplifiedchinese.GBK.NewDecoder())
	data, err := ioutil.ReadAll(reader)
	if err != nil {
		fmt.Print(err)
		return ""
	}
	return string(data)
}

func ConvertUtf8ToGbk(str string) string {
	reader := transform.NewReader(strings.NewReader(str), simplifiedchinese.GBK.NewEncoder())
	data, err := ioutil.ReadAll(reader)
	if err != nil {
		fmt.Print(err)
		return ""
	}
	return string(data)
}

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
