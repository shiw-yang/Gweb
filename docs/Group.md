# Group Control

分组控制（Group Control）是Web框架应该提供的基础功能。

经常需要对路由中的一组进行相似的处理和控制，例如GeekTutu所举的例子：

- 以`/post`开头的路由匿名可访问。
- 以`/admin`开头的路由需要鉴权。
- 以`/api`开头的路由是 RESTful 接口，可以对接第三方平台，需要三方平台鉴权。



## 分组嵌套

路由分组通过前缀来划分，还需要支持分组的嵌套，结合上中间件的支持，分组控制使得中间件在使用中的效率和收益更加明显。这里还是结合GeekTutu举得例子：

> 例如`/admin`的分组，可以应用鉴权中间件；`/`分组应用日志中间件，`/`是默认的最顶层的分组，也就意味着给所有的路由，即整个框架增加了记录日志的能力

### 实现思路

首先，要先思考一个Group对象需要具有那些属性：

- 前缀`prefix` ：支持找到父组件的Group，也能通过profix + pattern的方法注册完整路由，对相同prefix的路由执行应用相同的middlewares
- 中间件数组`middlewares`：用于记录该分组下需要执行的所有中间件，方便遍历执行

在没有引入分组控制之前，主要是通过`addRoute`方法来添加路由，因此，Group对象也应该有这样的能力，这里举一个使用分组的一个例子，主要的模仿gin的注册路由的使用

- 引擎Engine：获得整个系统基础资源调度的能力

```go
r := gee.New()
v1 := r.Group("/v1")
v1.GET("/", func(c *gee.Context) {
	c.HTML(http.StatusOK, "<h1>Hello Gee</h1>")
})
```

一通逆推猛如虎啊，终于能够理解了为什么代码里面要这样定义Group结构体，顺带一提，在geektutu的实现里面多了一个parent对象，估计是为了方便获取父组件的middlewares，但是其实perfix就够用了，因此我自己最终的实现是这样的：

#### gweb/routerGroup.go

```go
// RouterGroup packets into the Engine
	RouterGroup struct {
		prefix      string
		middlewares []HandlerFunc // support middleware
		engine      *Engine       // all groups share an Engine instance
	}
```

在geektutu包括gin的实现里面，都是把`Engine`作为最顶层的分组，通过Go语言的Embedded特性，使得`Engine`拥有`Group`所拥有的能力

```go
// Engine implement the interface of ServeHTTP
	Engine struct {
   	*RouterGroup  // Embed type in go: make engine has RouterGroup's power
   	router        *router
   	groups        []*RouterGroup     // store all groups
   	htmlTemplates *template.Template // for html render
   	funcMap       template.FuncMap   // for html render
	}
```

在这样之后，就可以将原先由引擎`Engine`才能实现的路由创建功能全部交给分组`Group`实现。

> 这里对于`Engine`的设计和定义，Gin和Geektutu的实现不一样的的点是：
>
> geektutu用的是*RouterGroup，而Gin用的是RouterGroup。
>
> 两者的差别不大，仅仅只是在引擎构造的时候需要处理好地址和对象的区别传参

## 出现的问题

最后，我根据Gin的结构，重构了gweb的代码。在调整的过程中，我发现了Go语言中一个叫做循环引用的问题，原先我将中间件单独整理为一个包`middlewares`，所以文件目录结构如图所示：

```BASH
.                 
|-- go.mod        
|-- gweb          
|   |-- context.go
|   |-- go.mod    
|   |-- gweb.go
|   |-- gweb_test.go
|   |-- middlewares
|   		|-- logger.go
|   |-- router.go
|   |-- routerGroup.go
|   |-- router_test.go
|   `-- trie.go
|-- main.go
|-- pojo
|-- static
`-- templates
```

可以看到的是`logger.go`是在`middlewares`包下的，因为我本身是希望在做完简单的基础框架之后，在制作一系列支持框架的中间件，因此我希望把他们整理在一个目录下，但是这样呢就会在操作的过程中导致`import cycle`的问题。

因为在`Logger.go`里面，需要用到`gweb`包下的东西，在import里面有这样的：

```go
// Logger.go
package middle

import (
	"gweb"
  "log"
  "time"
)
```

而在`gweb.go`里面又需要再次将middleware包引用，就变成了

```go
package gweb

import (
	// other packages
  "middlewares"
)
```

这样就会产生一个循环引用的问题。

Gin的解决方案：

由于logger属于是gin提供的一个简单的日志中间件，因此，只需要简单的把他放到和Gin同级目录就好。

如果希望加入更多的中间件，可以通过新建不同的package来提供类似插件级别的使用



