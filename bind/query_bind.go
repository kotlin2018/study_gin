package bind

import "net/http"

// 查询参数绑定
type queryBinding struct {}

// tag `query:""`
func (queryBinding) Name() string {
	return "query"
}

// url查询参数绑定 https://localhost:8080/?a=1&a=2&b=3
func (queryBinding) Bind(r *http.Request,obj interface{}) error {
	values := r.URL.Query() // values = a =[1,2], b=3
	if err := mapForm(obj,values);err !=nil {
		return err
	}
	return validate(obj)
}

