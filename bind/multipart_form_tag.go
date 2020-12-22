package bind

import (
	"errors"
	"mime/multipart"
	"net/http"
	"reflect"
)

type multipartRequest http.Request

var _ setter = (*multipartRequest)(nil)

func (r *multipartRequest) TrySet(value reflect.Value, field reflect.StructField, key string, opt setOptions) (isSet bool, err error) {
	// multipart请求的（一个）文件记录的信息;如果存在
	if files := r.MultipartForm.File[key]; len(files) !=0 { // 返回 一个 []*FileHeader
		return setByMultipartFormFile(value, field, files)
	}
	return setByForm(value, field, r.MultipartForm.Value, key, opt)
}

func setByMultipartFormFile(value reflect.Value, field reflect.StructField, files []*multipart.FileHeader) (isSet bool, err error) {
	switch value.Kind() {
	// 如果是指针
	case reflect.Ptr:
		// 如果存在类型
		switch value.Interface().(type) {
		case *multipart.FileHeader:
			value.Set(reflect.ValueOf(files[0]))
			return true,nil
		}
	case reflect.Struct:
		switch value.Interface().(type) {
		case multipart.FileHeader:
			value.Set(reflect.ValueOf(*files[0]))
			return true, nil
		}
	case reflect.Slice:
		slice := reflect.MakeSlice(value.Type(),len(files),len(files))
		isSet,err = setArrayOfMultipartFormFiles(slice, field, files)
		if err != nil || !isSet {
			return isSet, err
		}
		value.Set(slice)
		return true, nil
	case reflect.Array:
		return setArrayOfMultipartFormFiles(value, field, files)
	}
	return false, errors.New("multipart.FileHeader不支持的字段类型")
}

func setArrayOfMultipartFormFiles(value reflect.Value, field reflect.StructField, files []*multipart.FileHeader) (isSet bool, err error) {
	if value.Len() != len(files) {
		return false,errors.New("[]*multipart.FileHeader不支持的数组长度")
	}
	// 遍历file切片
	for i := range files {
		ok,err := setByMultipartFormFile(value.Index(i),field, files[i:i+1])
		if err !=nil || !ok {
			return ok,err
		}
	}
	return true,nil
}


