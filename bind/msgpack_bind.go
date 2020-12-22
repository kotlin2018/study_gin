// +build !nomsgpack

package bind

import (
	"bytes"
	"github.com/ugorji/go/codec"
	"io"
	"net/http"
)

// 编解码文本格式
type msgPackBinding struct {}

// 返回 tag标记的名称: `msg_pack`
func (msgPackBinding) Name() string {
	return "msg_pack"
}

// 将请求参数 解析成 `msg` tag 标记的 go 结构体
func (msgPackBinding) Bind(req *http.Request, obj interface{}) error {
	return decodeMsgPack(req.Body, obj)
}

// 将 文本文件中的值 绑定到 使用`msg` tag 标记的 go结构体中
func (msgPackBinding) BindBody(body []byte, obj interface{}) error {
	return decodeMsgPack(bytes.NewReader(body), obj)
}

// 解析出文本文件中的值
func decodeMsgPack(r io.Reader, obj interface{}) error {
	cdc := new(codec.MsgpackHandle)
	if err := codec.NewDecoder(r, cdc).Decode(&obj); err != nil {
		return err
	}
	return validate(obj)
}

