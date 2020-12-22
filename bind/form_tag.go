package bind

import (
	"errors"
	"fmt"
	"gee/internal/bytesconv"
	"gee/internal/json"
	"reflect"
	"strconv"
	"strings"
	"time"
)

// 这个文件就是为请求源是 form(表单) 设置tag值 (即: `form:""` 的用法)

var errUnknownType = errors.New("unknown type")
var emptyField = reflect.StructField{}
//4、
type setter interface {
 	TrySet(value reflect.Value,field reflect.StructField,key string,opt setOptions) (isSet bool,err error)
}

//为字段赋值的结构体
type setOptions struct {
	isDefaultExists bool
	defaultValue 	string
}
//
type formSource map[string][]string
//
var _ setter = formSource(nil)

// 5、通过请求的表单源来设置一个值(如map[string][]string)
func (form formSource) TrySet(value reflect.Value, field reflect.StructField, tagValue string, opt setOptions) (isSetted bool, err error) {
	return setByForm(value, field, form, tagValue, opt)
}

//11、 为Int类型字段设置值
func setIntField(val string,bitSize int,field reflect.Value) error {
	if val == "" {
		val = "0"
	}
	intVal ,err := strconv.ParseInt(val,10,bitSize)
	if err == nil {
		field.SetInt(intVal)
	}
	return err
}

//18、 为持续时间类型设置值
func setTimeDuration(val string,value reflect.Value,field reflect.StructField) error {
	d,err := time.ParseDuration(val) // 这个"10h" 转换为 10h0m0s
	if err !=nil {
		return err
	}
	value.Set(reflect.ValueOf(d))
	return nil
}

//12、
func setUintField(val string, bitSize int, field reflect.Value) error {
	if val == "" {
		val = "0"
	}
	uintVal, err := strconv.ParseUint(val, 10, bitSize)
	if err == nil {
		field.SetUint(uintVal)
	}
	return err
}

//13、
func setBoolField(val string, field reflect.Value) error {
	if val == "" {
		val = "false"
	}
	boolVal, err := strconv.ParseBool(val)
	if err == nil {
		field.SetBool(boolVal)
	}
	return err
}

//14、
func setFloatField(val string, bitSize int, field reflect.Value) error {
	if val == "" {
		val = "0.0"
	}
	floatVal, err := strconv.ParseFloat(val, bitSize)
	if err == nil {
		field.SetFloat(floatVal)
	}
	return err
}

//15、 为结构体时间类型字段设置值
func setTimeField(val string, structField reflect.StructField, value reflect.Value) error {
	timeFormt := structField.Tag.Get("time_format")
	if timeFormt == "" {
		timeFormt = time.RFC3339
	}
	switch tf := strings.ToLower(timeFormt); tf { //将大写字母全部转换成小写字母
	case "unix","unixnano":
		tv,err := strconv.ParseInt(val,10,64)
		if err !=nil {
			return err
		}

		d := time.Duration(1) //持续时间是 1
		if tf == "unixnano" {
			d = time.Second // 1秒钟
		}

		t := time.Unix(tv/int64(d),tv%int64(d))
		value.Set(reflect.ValueOf(t))
		return nil
	}

	if val == "" {
		value.Set(reflect.ValueOf(time.Time{}))
		return nil
	}

	l := time.Local
	if isUTC,_ := strconv.ParseBool(structField.Tag.Get("time_utc"));isUTC {
		l = time.UTC
	}

	if locTag := structField.Tag.Get("time_location");locTag != "" {
		loc,err := time.LoadLocation(locTag)
		if err !=nil {
			return err
		}
		l = loc
	}

	t, err := time.ParseInLocation(timeFormt,val,l)
	if err !=nil {
		return err
	}
	value.Set(reflect.ValueOf(t))
	return nil
}

//10、 根据字段类型:(map,slice,array,chan等聚合类型; int,string等基本数据类型),将tag的值绑定到字段上(即: 字段值 = tag值)
func setWithProperType(val string,value reflect.Value,field reflect.StructField) error {
	switch value.Kind() {
	case reflect.Int:
		return setIntField(val,0,value)
	case reflect.Int8:
		return setIntField(val,8,value)
	case reflect.Int16:
		return setIntField(val,16,value)
	case reflect.Int32:
		return setIntField(val,32,value)
	case reflect.Int64:
		// 类型断言
		switch value.Interface().(type) {
		case time.Duration:
			return setTimeDuration(val,value,field)
		}
		return setIntField(val,64,value)
	case reflect.Uint:
		return setUintField(val, 0, value)
	case reflect.Uint8:
		return setUintField(val, 8, value)
	case reflect.Uint16:
		return setUintField(val, 16, value)
	case reflect.Uint32:
		return setUintField(val, 32, value)
	case reflect.Uint64:
		return setUintField(val, 64, value)
	case reflect.Bool:
		return setBoolField(val, value)
	case reflect.Float32:
		return setFloatField(val, 32, value)
	case reflect.Float64:
		return setFloatField(val, 64, value)
	case reflect.String:
		value.SetString(val)
	case reflect.Struct:
		switch value.Interface().(type) {
		case time.Time:
			return setTimeField(val, field, value)
		}
		return json.Unmarshal(bytesconv.StringToBytes(val), value.Addr().Interface()) //将json反序列化成结构体
	case reflect.Map:
		return json.Unmarshal(bytesconv.StringToBytes(val), value.Addr().Interface())
	default:
		return errUnknownType
	}
	return nil
}

//16、
func setArray(val []string,value reflect.Value,field reflect.StructField) error {
	for i,s := range val {
		err := setWithProperType(s,value.Index(i),field)
		if err != nil {
			return err
		}
	}
	return nil
}

//17、
func setSlice(val []string,value reflect.Value,field reflect.StructField) error {
	slice := reflect.MakeSlice(value.Type(),len(val),len(val))
	err := setArray(val,slice,field)
	if err != nil {
		return err
	}
	value.Set(slice)
	return nil
}


// 9、根据字段的类型，将tag的值绑定到字段上(即: 字段值 = tag的值)
func setByForm(value reflect.Value,field reflect.StructField,form map[string][]string,tagValue string,opt setOptions) (isSet bool,err error) {
	vs,ok := form[tagValue] // tagValue中以 ";"分割的值
	if !ok && !opt.isDefaultExists {
		return false,nil
	}
	switch value.Kind() {
	// 如果字段类型是切片
	case reflect.Slice:
		if !ok {
			vs = []string{opt.defaultValue}
		}
		return true,setSlice(vs,value,field)
	// 如果字段类型是数组
	case reflect.Array:
		if !ok {
			vs = []string{opt.defaultValue}
		}
		if len(vs) != value.Len() {
			return false, fmt.Errorf("%q 这是一个无效值 %s", vs, value.Type().String())
		}
		return true, setArray(vs, value, field)
	default:
		var val string
		if !ok {
			val = opt.defaultValue
		}
		if len(vs) > 0 {
			val = vs[0]
		}
		return true, setWithProperType(val, value, field)
	}
}

//19、 例如:str ="chicken" , sep="ken"; 则返回 head = "chic" ,tail = ""
// 例如:str ="chicken" , sep="dmr"; 则返回 head = "chicken", "")
func head(str,sep string)(head,tail string){
	idx := strings.Index(str,sep)
	if idx < 0 {
		return str,""
	}
	return str[:idx],str[idx+len(sep):]
}

//1、通过 `uri` tag 将 tag的值绑定到 字段
func mapUri(ptr interface{}, m map[string][]string) error {
	return mapFormByTag(ptr, m, "uri")
}

//2、通过 `form` tag 将 tag的值绑定到 字段
func mapForm(ptr interface{}, form map[string][]string) error {
	return mapFormByTag(ptr, form, "form")
}
// 3、同过tag为类型是map的字段赋值
func mapFormByTag(ptr interface{}, form map[string][]string, tag string) error {
	ptrVal := reflect.ValueOf(ptr)
	var pointed interface{}
	// 如果ptrVal的类型是指针
	if ptrVal.Kind() == reflect.Ptr {
		ptrVal = ptrVal.Elem() // 获取reflect.Value的指针，该指针能操作反射对象
		pointed = ptrVal.Interface() //当前值作为接口
	}

	// 如果类型是Map 并且 Map的key是string类型
	if ptrVal.Kind() == reflect.Map &&
		// 如果 map的key类型是 string (即:map[string])
		ptrVal.Type().Key().Kind() == reflect.String {
		if pointed !=nil {
			ptr = pointed
		}
		return setFormMap(ptr, form)
	}
	return mappingByPtr(ptr, formSource(form), tag)
}


//6、
func mappingByPtr(ptr interface{}, setter setter, tag string) error {
	_, err := mapping(reflect.ValueOf(ptr), emptyField, setter, tag)
	return err
}

//7、通过 `form` 这个tag标签实现: 将tag的值赋值给字段 (即: 字段的值 = tag的值)
func mapping(value reflect.Value, field reflect.StructField, setter setter, tag string) (bool, error) {
	// 忽略这个"-"字段
	if field.Tag.Get(tag) == "-" {
		return false,nil
	}

	var vKind = value.Kind()
	// 如果value是指针
	if vKind == reflect.Ptr {
		var isNew bool
		vPtr := value
		//1、如果value的值 == 0,获取反射对象value的类型
		// 如果value的值为nil,并且value是: chan、func、接口、map、指针或slice值
		if value.IsNil() {
			isNew = true
			vPtr = reflect.New(value.Type().Elem()) //返回map类型的返回值类型: map slice ptr chan array 例如: map[interface{}]map, map[interface{}][]
		}

		//2、如果value的值 != 0,递归获取反射对象value的类型
		isSet,err := mapping(vPtr.Elem(),field,setter,tag)
		if err !=nil {
			return false, err
		}

		//如果 1、2 步都执行了
		if isNew && isSet {
			// value重新赋值
			value.Set(vPtr)
		}
		return isSet,nil
	}

	// 如果类型是结构体，或者该字段是嵌入字段(即是引用类型)
	if vKind != reflect.Struct || !field.Anonymous {
		ok, err := tryToSetValue(value, field, setter, tag)
		if err != nil {
			return false, err
		}
		if ok {
			return true, nil
		}
	}

	if vKind == reflect.Struct {
		tValue := value.Type()
		var isSet bool // isSet = true 这个结构体有引用类型的首字母大写(导出)字段
		// 遍历结构体字段数,获取每一个字段的信息
		for i :=0; i < value.NumField(); i++ {
			sf := tValue.Field(i) // 具体的字段 (封装了该字段的所有信息)
			// 如果是未导出字段,并且字段不是引用类型(map struct chan slice 等聚合类型)
			if sf.PkgPath != "" && !sf.Anonymous {
				continue  // 中断当前循环，继续后面的循环 (即跳过去)
			}
			// 如果是导出字段，并且 是聚合类型字段(即 struct,map,chan,slice等容器)
			// 则递归执行
			ok, err := mapping(value.Field(i),tValue.Field(i),setter,tag)
			if err !=nil {
				return false, err
			}
			isSet = isSet || ok //isSet = true
		}
		return isSet,nil
	}
	return false,nil
}

//8、 为 tag赋值, 例如:  `form:";default=9"` `form:"slice;default=9"` , " "内表示tag的值,多个值用";"或者","隔开
func tryToSetValue(value reflect.Value, field reflect.StructField, setter setter, tag string) (bool, error) {
	var tagValue string
	var setOpt setOptions

	tagValue = field.Tag.Get(tag) // 字段 tag的值
	tagValue,opts := head(tagValue,";") //tag中的多个值之间用";"隔开

	if tagValue == ""{
		tagValue = field.Name // 将字段名赋值给tag值
	}

	if tagValue == ""{ // 当字段为空时,(即该字段不存在)
		return false,nil
	}

	var opt string
	// 递归，并设置默认tag值，即 "default"表示tag的默认值
	for len(opts) > 0 {
		opt,opts = head(opts,";")
		if k,v := head(opt,"="); k == "default" {//`form` 默认值
			setOpt.isDefaultExists = true
			setOpt.defaultValue = v
		}
	}
	return setter.TrySet(value,field,tagValue,setOpt)
}

//20、 任意类型转换成 map[string][]string
func setFormMap(ptr interface{}, form map[string][]string) error {
	//返回类型的类型
	//如果类型的类型不是Array, Chan, Map, Ptr或Slice，它会报错
	el := reflect.TypeOf(ptr).Elem()

	// 如果类型是切片，例如: map[interface{}][]
	if el.Kind() == reflect.Slice {
		ptrMap,ok := ptr.(map[string][]string)
		if !ok {
			return errors.New("不能转换为map[string][]string") //convert 转换
		}

		for k, v := range form {
			ptrMap[k] = v
		}
		return nil
	}
	// 如果类型不是切片
	ptrMap, ok := ptr.(map[string]string)
	if !ok {
		return errors.New("不能转换为map[string]string")
	}
	for k,v := range form {
		ptrMap[k] = v[len(v)-1]
	}
	return nil
}

