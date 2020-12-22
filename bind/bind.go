// +build !nomsgpack

package bind

import "net/http"

const (
	HttpJSON     			= "application/json"
	HttpYAML       			= "application/x-yaml"
	HttpXML 				= "application/xml"
	HttpXML2    			= "text/xml"
	HttpHTML				= "text/html"
	HttpPlain    			= "text/plain"
	HttpPOSTForm 			= "application/x-www-form-urlencoded"
	HttpMultipartPOSTForm 	= "multipart/form-data"
	HttpProtoBuf   			= "application/x-protobuf"
	HttpMsgPack   			= "application/x-msgpack"
	HttpMsgPack2    		= "application/msgpack"

)

// 绑定接口
type Binding interface {
	Name() string
	Bind(*http.Request,interface{})error
}

// 将绑定方法添加到绑定接口
type BindingBody interface {
	Binding
	BindBody([]byte,interface{})error
}

// 将绑定uri添加到绑定接口
type BindingUri interface {
	Name()string
	BindUri(map[string][]string,interface{})error
}

// 结构体验证器
type StructValidator interface {
	//可以接收任何类型的类型，并且它永远不会报错
	//如果接收到的类型不是struct，任何验证都应该被跳过，并且必须返回nil。
	//如果接收到的类型是struct或指向struct的指针，则应该执行验证。
	//如果结构无效或验证失败，则返回描述性错误。
	//否则必须返回nil。
	ValidateStruct(interface{}) error
	// StructValidator接口的基础实现。
	// Engine返回支持的基础验证器引擎
	Engine()interface{}
}

// 验证结构体的验证器实例
var Validator StructValidator = &defaultValidator{}

// 这些实现了绑定接口，可以用来绑定数据
var (
	JSON          = jsonBinding{}
	XML           = xmlBinding{}
	Form          = formBinding{}
	Query         = queryBinding{}
	FormPost      = formPostBinding{}
	FormMultipart = formMultipartBinding{}
	ProtoBuf      = protobufBinding{}
	MsgPack       = msgPackBinding{}
	YAML          = yamlBinding{}
	Uri           = uriBinding{}
	Header        = headerBinding{}
)

//默认返回基于HTTP方法的适当绑定实例
func Default(method, contentType string) Binding {
	if method == http.MethodGet {
		return Form
	}
	switch contentType {
	case HttpJSON:
		return JSON
	case HttpXML, HttpXML2:
		return XML
	case HttpProtoBuf:
		return ProtoBuf
	case HttpMsgPack, HttpMsgPack2:
		return MsgPack
	case HttpYAML:
		return YAML
	case HttpMultipartPOSTForm:
		return FormMultipart
	default: // case MIMEPOSTForm:
		return Form
	}

}

func validate(obj interface{}) error {
	if Validator == nil {
		return nil
	}
	return Validator.ValidateStruct(obj)
}