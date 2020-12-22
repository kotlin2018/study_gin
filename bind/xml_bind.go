package bind

import (
	"bytes"
	"encoding/xml"
	"io"
	"net/http"
)

type xmlBinding struct {}

// 返回 tag标记的名称: `xml`
func (xmlBinding) Name() string {
	return "xml"
}

// 将请求参数 解析成 `xml` tag 标记的 go 结构体
func (xmlBinding) Bind(req *http.Request, obj interface{}) error {
	return decodeXML(req.Body,obj)
}

// 将 .xml文件中的值 绑定到 使用`xml` tag 标记的 go结构体中
func (xmlBinding) BindBody(body []byte, obj interface{}) error {
	return decodeXML(bytes.NewReader(body),obj)
}

// 解析出xml文件中的值
func decodeXML (r io.Reader, obj interface{}) error {
	decoder := xml.NewDecoder(r)
	if err := decoder.Decode(obj); err !=nil {
		return err
	}
	return validate(obj)
}