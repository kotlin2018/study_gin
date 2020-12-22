package gee

import (
	"gee/bind"
	"math"
	"net/http"
)

const (
	HttpJSON 				= bind.HttpJSON
	HttpYAML 				= bind.HttpYAML
	HttpXML 				= bind.HttpXML
	HttpXML2				= bind.HttpXML2
	HttpHTML 				= bind.HttpHTML
	HttpPlain 				= bind.HttpPlain
	HttpPOSTForm 			= bind.HttpPOSTForm
	HttpMultipartPOSTForm 	= bind.HttpMultipartPOSTForm
)

const BodyByTestKey = "BodyByTestKey"
const abortIndex int8 = math.MaxInt8/2

// 上下文结构体,封装了 *http.Request，http.ResponseWriter等等
type Context struct {
	writer http.ResponseWriter
}
