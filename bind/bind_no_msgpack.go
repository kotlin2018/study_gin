// +build nomsgpack

package bind

import "net/http"

// 该文件与http.Request.Body有关
// Content-Type 中最常见的数据格式
const (
	HttpJSON 				= "application/json"
	HttpXML 				= "application/xml"
	HttpYAML 				= "application/x-yaml"
	HttpHTML 				= "text/html"
	HttpXML2				= "text/xml"
	HttpPlain				= "text/plain"
	HttpPOSTForm 			= "application/x-www-form-urlencoded"
	HttpMultipartPOSTForm 	= "multipart/form-data"
	HttpPROTOBUF    		= "application/x-protobuf"
)

// 绑定接口
type Binding interface {
	Name() string  // 该绑定的名字
	Bind(*http.Request,interface{}) error // 具体的绑定方法
}

// 能存储数据的绑定接口
type BindingBody interface {
	Binding
	BindBody([]byte, interface{}) error
}

type BindingUri interface {
	Name() string
	BindUri(map[string][]string, interface{}) error
}

type StructValidator interface {
	// ValidateStruct可以接收任何类型的数据，即使配置不正确也不应该报错。
	//如果接收到的类型不是struct，任何验证都应该被跳过，并且必须返回nil。
	//如果接收到的类型是struct或指向struct的指针，则应该执行验证。
	//如果结构无效或验证失败，则返回描述性错误。
	//否则必须返回nil。
	ValidateStruct(interface{}) error

	// Engine返回支持的基础验证器引擎
	// StructValidator实现
	Engine() interface{}
}

//Validator是实现StructValidator接口的默认验证器
var Validator StructValidator = &defaultValidator{}

//这些都实现了绑定接口，可以用来绑定数据
var (
	JSON          = jsonBinding{}
	XML           = xmlBinding{}
	Form          = formBinding{}  // 表单绑定
	Query         = queryBinding{}
	FormPost      = formPostBinding{}
	FormMultipart = formMultipartBinding{}
	ProtoBuf      = protobufBinding{}
	YAML          = yamlBinding{}
	Uri           = uriBinding{}
	Header        = headerBinding{}
)

// 实现了Binding接口的默认函数，返回基于HTTP方法的绑定实例
func Default(method,contentType string) Binding {
	if method == "GET" {
		return Form
	}
	switch contentType {
	case HttpJSON:
		return JSON
	case HttpXML,HttpXML2:
		return XML
	case HttpProtoBuf:
		return ProtoBuf
	case HttpYAML:
		return YAML
	case HttpMultipartPOSTForm:
		return FormMultipart
	default:
		return Form
	}
}

func validate(obj interface{}) error {
	if Validator == nil {
		return nil
	}
	return Validator.ValidateStruct(obj)
}








