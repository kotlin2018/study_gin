1、bind包照顾gin binding 包写的，里面大量函数我都照自己的理解写了中文注释；
2、gin 的请求参数绑定，即 `binding:""` 就是这个包实现的。
3、gin 的ShouldBindJSON()、ShoulBindYAML()等方法是基于 ShoulBindWith()实现的，而ShouldBindWith(）函数是基于binding.Binding.Bind()实现的
最底层的就是实现就是以下这个接口
````
// gin源码
package binding
type Binding interface {
	Name() string
	Bind(*http.Request, interface{}) error
}
````


