package bind

import (
	"bytes"
	"gopkg.in/yaml.v2"
	"io"
	"net/http"
)

type yamlBinding struct{}

// 返回 tag标记的名称: `yaml`
func (yamlBinding) Name() string {
	return "yaml"
}

// 将请求参数转换成 `yaml` tag 标记的 go 结构体
func (yamlBinding) Bind(req *http.Request, obj interface{}) error {
	return decodeYAML(req.Body,obj)
}

// 将 .yaml文件中的值 绑定到使用 `yaml` tag 标记的 go结构体中
func (yamlBinding) BindBody (body []byte, obj interface{}) error {
	return decodeYAML(bytes.NewReader(body),obj)
}

// 从输入流中解析出yam文件中的值
// 将输入转换成 `yaml`Tag 标记的 结构体
func decodeYAML(r io.Reader, obj interface{}) error {
	decoder := yaml.NewDecoder(r)
	if err := decoder.Decode(obj);err !=nil {
		return err
	}
	return validate(obj)
}