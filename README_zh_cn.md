# Martini  [![wercker status](https://app.wercker.com/status/174bef7e3c999e103cacfe2770102266 "wercker status")](https://app.wercker.com/project/bykey/174bef7e3c999e103cacfe2770102266) [![GoDoc](https://godoc.org/github.com/codegangsta/martini?status.png)](http://godoc.org/github.com/codegangsta/martini)

Martini是一个强大为了编写模块化Web应用而生的GO语言框架.

## 第一个应用

在你安装了GO语言和设置了你的[GOPATH](http://golang.org/doc/code.html#GOPATH)之后, 创建你的自己的'.go'文件, 这里我们假设它的名字叫做 'server.go'.

~~~ go
package main

import "github.com/codegangsta/martini"

func main() {
  m := martini.Classic()
  m.Get("/", func() string {
    return "Hello world!"
  })
  m.Run()
}
~~~

然后安装Martini的包. (注意Martini需要Go语言1.1或者以上的版本支持):
~~~
go get github.com/codegangsta/martini
~~~

最后运行你的服务:
~~~
go run server.go
~~~

这时你将会有一个Martini的服务监听了, 地址是: `localhost:3000`.

## 获得帮助

请加入: [邮件列表](https://groups.google.com/forum/#!forum/martini-go)

或者可以查看在线演示地址: [演示视频](http://martini.codegangsta.io/#demo)

## 功能列表
* 使用极其简单.
* 无侵入式的设计.
* 很好的与其他的Go语言包协同使用.
* 超赞的路径匹配和路由请求.
* 模块化的设计 - 容易插入功能件，也容易将其拔出来.
* 已有很多的中间件可以直接使用.
* 框架内已拥有很好的开箱即用的功能支持.
* **完全兼容[http.HandlerFunc](http://godoc.org/net/http#HandlerFunc)接口.**

## 更多中间件
更多的中间件和功能组件, 请查看代码仓库: [martini-contrib](https://github.com/martini-contrib).

## 目录
* [核心 Martini](#classic-martini)
  * [处理器](#handlers)
  * [路由请求](#routing)
  * [服务](#services)
  * [服务静态文件](#serving-static-files)
* [中间件处理器](#middleware-handlers)
  * [Next()](#next)
* [常见问答](#faq)

## 核心 Martini
为了更快速的启用Martini, [martini.Classic()](http://godoc.org/github.com/codegangsta/martini#Classic) 提供了一些默认的方便Web开发的工具:
~~~ go
  m := martini.Classic()
  // ... middleware and routing goes here
  m.Run()
~~~

下面是Martini核心已经包含的功能 [martini.Classic()](http://godoc.org/github.com/codegangsta/martini#Classic):
  * Request/Response Logging （请求/相应日志） - [martini.Logger](http://godoc.org/github.com/codegangsta/martini#Logger)
  * Panic Recovery （容错） - [martini.Recovery](http://godoc.org/github.com/codegangsta/martini#Recovery)
  * Static File serving （静态文件服务） - [martini.Static](http://godoc.org/github.com/codegangsta/martini#Static)
  * Routing （路由请求） - [martini.Router](http://godoc.org/github.com/codegangsta/martini#Router)

### 处理器
处理器是Martini的灵魂和核心所在. 一个处理器基本上可以是任何的函数:
~~~ go
m.Get("/", func() {
  println("hello world")
})
~~~

#### 返回值
当一个处理器返回结果的时候, Martini将会把返回值作为字符串写入到当前的[http.ResponseWriter](http://godoc.org/net/http#ResponseWriter)里面:
~~~ go
m.Get("/", func() string {
  return "hello world" // HTTP 200 : "hello world"
})
~~~

另外你也可以选择性的返回多一个状态码:
~~~ go
m.Get("/", func() (int, string) {
  return 418, "i'm a teapot" // HTTP 418 : "i'm a teapot"
})
~~~

#### 服务的注入 
处理器是通过反射来调用的. Martini 通过*Dependency Injection* *（依赖注入）* 来为处理器注入参数列表. **这样使得Martini与Go语言的`http.HandlerFunc`接口完全兼容.** 

如果你加入一个参数到你的处理器, Martini将会搜索它参数列表中的服务，并且通过类型判断来解决依赖关系:
~~~ go
m.Get("/", func(res http.ResponseWriter, req *http.Request) { // res 和 req 是通过Martini注入的
  res.WriteHeader(200) // HTTP 200
})
~~~

下面的这些服务已经被包含在核心Martini中: [martini.Classic()](http://godoc.org/github.com/codegangsta/martini#Classic):
  * [*log.Logger](http://godoc.org/log#Logger) - Martini的全局日志.
  * [martini.Context](http://godoc.org/github.com/codegangsta/martini#Context) - http request context （请求上下文）.
  * [martini.Params](http://godoc.org/github.com/codegangsta/martini#Params) - `map[string]string` of named params found by route matching. （名字和参数键值对的参数列表）
  * [martini.Routes](http://godoc.org/github.com/codegangsta/martini#Routes) - Route helper service. （路由协助处理）
  * [http.ResponseWriter](http://godoc.org/net/http/#ResponseWriter) - http Response writer interface. (响应结果的流借口)
  * [*http.Request](http://godoc.org/net/http/#Request) - http Request. （http请求)

### 路由请求
在Martini中, 路由是一个HTTP方法配对一个URL匹配模型. 每一个路由可以对应一个或多个处理器方法:
Each route can take one or more handler methods:
~~~ go
m.Get("/", func() {
  // show something
})

m.Patch("/", func() {
  // update something
})

m.Post("/", func() {
  // create something
})

m.Put("/", func() {
  // replace something
})

m.Delete("/", func() {
  // destroy something
})

m.Options("/", func() {
  // http options
})

m.NotFound(func() {
  // handle 404
})
~~~

路由匹配的顺序是按照他们被定义的顺序执行的. 最先被定义的路由将会首先被用户请求匹配并调用.

路由模型可能包含参数列表, 可以通过[martini.Params](http://godoc.org/github.com/codegangsta/martini#Params)服务来获取:
~~~ go
m.Get("/hello/:name", func(params martini.Params) string {
  return "Hello " + params["name"]
})
~~~

路由匹配可以通过正则表达式或者glob的形式:
~~~ go
m.Get("/hello/**", func(params martini.Params) string {
  return "Hello " + params["_1"]
})
~~~

路由处理器可以被相互叠加使用, 例如很有用的地方可以是在验证和授权的时候:
~~~ go
m.Get("/secret", authorize, func() {
  // 该方法将会在authorize方法没有输出结果的时候执行.
})
~~~

### 服务
服务即是被注入到处理器中的参数. 你可以映射一个服务到*全局*或者*请求*的级别.


#### 全局映射
如果一个Martini实现了inject.Injector的接口, 那么映射成为一个服务就非常简单:
~~~ go
db := &MyDatabase{}
m := martini.Classic()
m.Map(db) // *MyDatabase 这个服务将可以在所有的处理器中被使用到.
// ...
m.Run()
~~~

#### 请求级别的影射
映射成为请求级别的服务可以用[martini.Context](http://godoc.org/github.com/codegangsta/martini#Context)来完成:
~~~ go
func MyCustomLoggerHandler(c martini.Context, req *http.Request) {
  logger := &MyCustomLogger{req}
  c.Map(logger) // 映射成为了 *MyCustomLogger
}
~~~

#### 映射值到接口
关于服务最强悍的地方之一就是它能够映射服务到接口. 例如说, 假设你想要覆盖[http.ResponseWriter](http://godoc.org/net/http#ResponseWriter)成为一个对象, 那么你可以封装它并包含你自己的额外操作, 你可以如下这样来编写你的处理器:
~~~ go
func WrapResponseWriter(res http.ResponseWriter, c martini.Context) {
  rw := NewSpecialResponseWriter(res)
  c.MapTo(rw, (*http.ResponseWriter)(nil)) // override ResponseWriter with our wrapper ResponseWriter
}
~~~

### Serving Static Files
A [martini.Classic()](http://godoc.org/github.com/codegangsta/martini#Classic) instance automatically serves static files from the "public" directory in the root of your server.
You can serve from more directories by adding more [martini.Static](http://godoc.org/github.com/codegangsta/martini#Static) handlers.
~~~ go
m.Use(martini.Static("assets")) // serve from the "assets" directory as well
~~~

## Middleware Handlers
Middleware Handlers sit between the incoming http request and the router. In essence they are no different than any other Handler in Martini. You can add a middleware handler to the stack like so:
~~~ go
m.Use(func() {
  // do some middleware stuff
})
~~~

You can have full control over the middleware stack with the `Handlers` function. This will replace any handlers that have been previously set:
~~~ go
m.Handlers(
  Middleware1,
  Middleware2,
  Middleware3,
)
~~~

Middleware Handlers work really well for things like logging, authorization, authentication, sessions, gzipping, error pages and any other operations that must happen before or after an http request:
~~~ go
// validate an api key
m.Use(func(res http.ResponseWriter, req *http.Request) {
  if req.Header.Get("X-API-KEY") != "secret123" {
    res.WriteHeader(http.StatusUnauthorized)
  }
})
~~~

### Next()
[Context.Next()](http://godoc.org/github.com/codegangsta/martini#Context) is an optional function that Middleware Handlers can call to yield the until after the other Handlers have been executed. This works really well for any operations that must happen after an http request:
~~~ go
// log before and after a request
m.Use(func(c martini.Context, log *log.Logger){
  log.Println("before a request")

  c.Next()
  
  log.Println("after a request")
})
~~~

## FAQ

### Where do I find middleware X?

Start by looking in the [martini-contrib](https://github.com/martini-contrib) projects. If it is not there feel free to contact a martini-contrib team member about adding a new repo to the organization.

* [auth](https://github.com/martini-contrib/auth) - Handlers for authentication.
* [binding](https://github.com/martini-contrib/binding) - Handler for mapping/validating a raw request into a structure.
* [gzip](https://github.com/martini-contrib/gzip) - Handler for adding gzip compress to requests
* [render](https://github.com/martini-contrib/render) - Handler that provides a service for easily rendering JSON and HTML templates.
* [acceptlang](https://github.com/martini-contrib/acceptlang) - Handler for parsing the `Accept-Language` HTTP header.
* [sessions](https://github.com/martini-contrib/sessions) - Handler that provides a Session service.
* [strip](https://github.com/martini-contrib/strip) - URL Prefix stripping.
* [method](https://github.com/martini-contrib/method) - HTTP method overriding via Header or form fields.
* [secure](https://github.com/martini-contrib/secure) - Implements a few quick security wins.
* [encoder](https://github.com/martini-contrib/encoder) - Encoder service for rendering data in several formats and content negotiation.

### How do I integrate with existing servers?

A Martini instance implements `http.Handler`, so it can easily be used to serve subtrees 
on existing Go servers. For example this is a working Martini app for Google App Engine:

~~~ go
package hello

import (
  "net/http"
  "github.com/codegangsta/martini"
)

func init() {
  m := martini.Classic()
  m.Get("/", func() string {
    return "Hello world!"
  })
  http.Handle("/", m)
}
~~~

### How do I change the port/host?

Martini's `Run` function looks for the PORT and HOST environment variables and uses those. Otherwise Martini will default to localhost:3000.
To have more flexibility over port and host, use the `http.ListenAndServe` function instead.

~~~ go
  m := martini.Classic()
  // ...
  http.ListenAndServe(":8080", m)
~~~

## Contributing
Martini is meant to be kept tiny and clean. Most contributions should end up in a repository in the [martini-contrib](https://github.com/martini-contrib) organization. If you do have a contribution for the core of Martini feel free to put up a Pull Request.

## About

Inspired by [express](https://github.com/visionmedia/express) and [sinatra](https://github.com/sinatra/sinatra)

Martini is obsessively designed by none other than the [Code Gangsta](http://codegangsta.io/)
