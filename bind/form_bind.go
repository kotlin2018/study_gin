package bind

import (
	"net/http"
)

const defaultMemory = 32 << 20

type formBinding struct {} //解析url查询参数(即http.Header携带的参数)+ POST、PUT、PATCH等方法传递的数据(即:http.Body中携带的数据)
type formPostBinding struct{} //只解析 http.Request.Body携带的数据
type formMultipartBinding struct {} //解析整个请求，将请求体解析为multipart/form-data。

// http.Request 请求的form表单数据绑定，tag标签名 `form:""`
func (formBinding) Name() string {
	return "form"
}

// (请求参数全绑定)
// 绑定:客户端, URL查询参数、PATCH、POST、PUT等方法传过来的表单数据,到自定义数据容器。
func (formBinding) Bind (r *http.Request,obj interface{}) error {
	// 解析通过form(表单)提交的请求参数
	if err := r.ParseForm(); err !=nil {
		return err
	}
	if err := r.ParseMultipartForm(defaultMemory);err !=nil {
		if err != http.ErrNotMultipart {
			return err
		}
	}
	// 将请求参数绑定到tag,然后tag的值赋值给 obj的值;实现 tag值绑定字段值
	if err := mapForm(obj,r.Form); err !=nil {
		return err // Form包含已解析的表单数据，包括URL字段的查询参数和PATCH、POST或PUT表单数据。此字段仅在调用ParseForm之后可用。HTTP客户端忽略形式而使用主体
	}
	return validate(obj)
}

// http.Request.Body中的数据绑定，tag标签名 `form-urlencoded:""`
func (formPostBinding) Name() string {
	return "form-urlencoded"
}

// (请求参数Body绑定,不绑定url查询参数)
// 绑定:客户端通过 PATCH、POST、PUT等方法传过来的参数
func (formPostBinding) Bind(r *http.Request,obj interface{}) error {
	if err := r.ParseForm(); err !=nil {
		return err

	}
	if err := mapForm(obj,r.PostForm);err !=nil {//PostForm包含解析后的表单数据，这些数据来自PATCH、POST或PUT主体参数。
		return err
	}
	return validate(obj)
}

func (formMultipartBinding) Name() string {
	return "multipart/form-data"
}

func (formMultipartBinding) Bind(r *http.Request,obj interface{}) error {
	if err := r.ParseMultipartForm(defaultMemory);err !=nil {
		return err
	}
	if err := mappingByPtr(obj,(*multipartRequest)(r),"form");err !=nil {
		return err
	}
	return validate(obj)
}