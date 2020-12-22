package gee

import (
	"fmt"
	"reflect"
)

type Cfg struct {

}

type Health struct {
	Url      string
	Cmd      string
	Interval string
	Timeout  string
	Disable  bool
}

func (m model) TableName() string {
	return "table_name"
}

type tabler interface {
	TableName() string
}

type model struct{}

func getTableName(v interface{}) {
	rt := reflect.TypeOf(v)
	rv := reflect.ValueOf(v)
	if tabler, ok := rv.Interface().(tabler); ok {
		fmt.Println(tabler.TableName())
	}

	if tabler, ok := reflect.New(rt).Interface().(tabler); ok {
		fmt.Println(tabler.TableName())
	}
}

type Ser struct {
	name string
	age int
}


func main() {

	//file, err := os.Open("config/config.yaml")
	//
	//if err != nil {
	//	panic(err)
	//}
	//
	//bytes, err := ioutil.ReadAll(file)
	//
	//if err != nil {
	//	panic(err)
	//}

	//cfg := Cfg{}

	//err = yaml.Unmarshal(bytes, &cfg)
	//
	//if err != nil {
	//	panic(err)
	//}
//	c := config.Config()
//
//
//
//	fmt.Println(c)

	//var mod model
	//getTableName(mod)
	//
	//n := 1.2345
	//fmt.Println("old value :", n)
	//
	////通过reflect.ValueOf获取n中的reflect.Value，注意:参数必须是指针才能修改其值
	//v := reflect.ValueOf(&n)
	//e := v.Elem()
	//
	//fmt.Println("type :", e.Type())
	//fmt.Println("can set :", e.CanSet())
	//
	////重新赋值
	//e.SetFloat(2.123)
	//fmt.Println("new value :", n)

	// 声明整型变量a并赋初值
	var a int = 1024
	// 获取变量a的反射值对象
	valueOfA := reflect.ValueOf(a)
	// 获取interface{}类型的值, 通过类型断言转换
	var getA int = valueOfA.Interface().(int)
	// 获取64位的值, 强制类型转换为int类型
	var getA2 int = int(valueOfA.Int())
	fmt.Println(getA, getA2)





}



