package bind

import (
	"github.com/go-playground/validator/v10"
	"reflect"
	"sync"
)

// 默认验证器
type defaultValidator struct {
	once 		sync.Once 			//在代码运行中需要的时候执行，且只执行一次
	validate 	*validator.Validate //参数验证器实例
}

// 默认验证器实例
var _ StructValidator = &defaultValidator{}

// 使用反射
// 默认验证器实现bind(绑定接口中的方法)
func (v *defaultValidator) ValidateStruct(obj interface{}) error {
	value := reflect.ValueOf(obj)
	valueType := value.Kind()
	// 如果参数类型是指针
	if valueType == reflect.Ptr {
		valueType = value.Elem().Kind()
	}
	// 如果参数类型是结构体
	if valueType == reflect.Struct {
		v.lazyInit()
		// 判断参数类型是不是结构体
		if err := v.validate.Struct(obj); err !=nil {
			return err
		}
	}
	return nil
}

// Engine返回支持默认值的基础验证器引擎
// 如果你想注册自定义验证，这很有用
func (v *defaultValidator) Engine() interface{} {
	v.lazyInit()
	return v.validate
}

//1、 设置验证器Tag(标签)名 `validator:""`
func (v *defaultValidator) lazyInit() {
	v.once.Do(func() {
		v.validate = validator.New()
		//v.validate.SetTagName("binding") gin源码用法: `binding:"required,"`
		v.validate.SetTagName("validator") // 自定义名
	})
}


