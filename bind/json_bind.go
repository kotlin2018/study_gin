package bind

import (
	"bytes"
	"fmt"
	"gee/internal/json"
	"io"
	"net/http"
)

// 将一个数字解编码为一个接口或者数字
var EnableDecoderUseNumber = false

//
var EnableDecoderDisallowUnknownFields = false

// json绑定结构体
type jsonBinding struct {}

// 返回 tag标记的名称: `json`
func (jsonBinding) Name() string {
	return "json"
}

// 将请求参数绑定成json结构体
func (jsonBinding) Bind(req *http.Request, obj interface{}) error {
	if req == nil || req.Body == nil {
		return fmt.Errorf("无效的请求")
	}
	return decodeJSON(req.Body,obj)
}

// // 将 .json文件中的值 绑定到使用 `json` tag 标记的 go结构体中
func (jsonBinding) BindBody(body []byte, obj interface{}) error {
	return decodeJSON(bytes.NewReader(body),obj)
}

// 将输入编码成json结构体 (json Tag 标记的结构体) `json`
func decodeJSON (r io.Reader,obj interface{}) error {
	decoder := json.NewDecoder(r)
	if EnableDecoderUseNumber { //将一个数字分解为接口{}中的一个数字，而不是浮点64
		decoder.UseNumber()
	}
	if EnableDecoderDisallowUnknownFields { // 导致解码器返回一个错误，当目标是一个结构并且输入包含的对象键与目标中任何非忽略的导出字段不匹配时
		decoder.DisallowUnknownFields()
	}
	// 从它的输入中读取下一个json编码的值，并将其存储在v指向的值中。
	if err := decoder.Decode(obj); err != nil {
		return err
	}
	return validate(obj)
}
