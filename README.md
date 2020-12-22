前言: 

我很喜欢: gin的 RESTful路由、路由分组、ShouldBindJSON()、ShouldBind()方法，这些方法直接将请求参数绑定到结构体。想弄明白gin是怎么做到的这些的，因此有了这个项目。后续我将做一个类似beego、gin的大而全的GolangWeb脚手架
		

1、这个项目的bind包是照着 gin binding 包写的，里面大量函数我都照自己的理解写了中文注释。


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


