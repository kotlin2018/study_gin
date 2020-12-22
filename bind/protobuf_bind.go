package bind

import (
	"io/ioutil"
	"net/http"
	"github.com/golang/protobuf/proto" //数据编码格式，类似于xml和json [使用链接](https://studygolang.com/articles/21434?fr=sidebar)
)

type protobufBinding struct {}

func (protobufBinding) Name () string {
	return "protobuf"
}

func (b protobufBinding) Bind(r *http.Request,obj interface{}) error {
	buf , err := ioutil.ReadAll(r.Body)
	if err != nil {
		return err
	}
	return b.BindBody(buf,obj)
}

func (protobufBinding) BindBody(body []byte, obj interface{}) error {
	if err := proto.Unmarshal(body, obj.(proto.Message)); err != nil {
		return err
	}
	// Here it's same to return validate(obj), but util now we can't add
	// `binding:""` to the struct which automatically generate by gen-proto
	return nil
	// return validate(obj)
}