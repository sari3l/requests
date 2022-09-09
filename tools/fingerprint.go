package tools

import (
	"github.com/CUCyber/ja3transport"
	"github.com/sari3l/requests/types"
	"net/http"
	"reflect"
	"strconv"
	"strings"
)

// HookClientJA3Func 传入JA3 Client指纹
//例如 "771,4865-4866-4867-49195-49199-49196-49200-52393-52392-49171-49172-156-157-47-53-10,0-23-65281-10-11-35-16-5-13-18-51-45-43-27-21,29-23-24,0"
func HookClientJA3Func(fingerprint string) types.Hook {
	return func(client any) (error, any) {
		c := client.(http.Client)
		tr, _ := ja3transport.NewTransport(fingerprint)
		reflect.ValueOf(&c).Elem().FieldByName("Transport").Set(reflect.ValueOf(tr))
		return nil, c
	}
}

// HookClientMitmFunc 传入MitmEngine指纹，自动转换
// 例如 "303:4,5,a,13,2f,32,33,35,38,39,c009,c00a,c013,c014:0,a,b,17,ff01:17,18:0"
// 目前不是很优雅，封装了JA3的函数，后面自己写个TP生成器
func HookClientMitmFunc(fingerprint string) types.Hook {
	parts := strings.Split(fingerprint, ":")
	if len(parts) < 5 {
		panic("错误长度")
	}
	var newParts = make([]string, 5)
	for i, part := range parts[:5] {
		items := strings.Split(part, ",")
		newPart := make([]string, len(items))
		for i2, item := range items {
			num, _ := strconv.ParseInt(item, 16, 0)
			newPart[i2] = strconv.FormatInt(num, 10)
		}
		newParts[i] = strings.Join(newPart, "-")
	}
	fingerprint = strings.Join(newParts, ",")
	return HookClientJA3Func(fingerprint)
}
