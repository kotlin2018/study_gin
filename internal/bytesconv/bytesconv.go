package bytesconv // 字节切片转换包

import (
	"reflect"
	"unsafe" //直接读写内存，包用于 Go 编译器 能操作未导出结构体
)

// 字符串转换成字节切片
func StringToBytes(s string)(b []byte) {
	sh := *(*reflect.StringHeader)(unsafe.Pointer(&s)) // 反射获取字符串指针
	bh := (*reflect.SliceHeader)(unsafe.Pointer(&b))   // 反射获取字节切片指针
	bh.Data,bh.Len,bh.Cap = sh.Data,sh.Len,sh.Len
	return b
}

// 字节切片转换成字符串
func BytesToString(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}
