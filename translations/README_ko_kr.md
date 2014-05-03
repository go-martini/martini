# Martini  [![wercker status](https://app.wercker.com/status/9b7dbc6e2654b604cd694d191c3d5487/s/master "wercker status")](https://app.wercker.com/project/bykey/9b7dbc6e2654b604cd694d191c3d5487)[![GoDoc](https://godoc.org/github.com/go-martini/martini?status.png)](http://godoc.org/github.com/go-martini/martini)

Martini is a powerful package for quickly writing modular web applications/services in Golang.
마티니(Martini)는 강력하고 손쉬운 웹애플리캐이션 / 웹서비스개발을 위한 Golang 모듈 패키지입니다.

## 시작하기

After installing Go and setting up your [GOPATH](http://golang.org/doc/code.html#GOPATH), create your first `.go` file. We'll call it `server.go`.

Go 인스톨 및 [GOPATH](http://golang.org/doc/code.html#GOPATH) 환경변수 설정 이후에, `.go` 파일 하나를 만들어 보죠..흠... 일단 `server.go`라고 부르겠습니다.
~~~ go
package main

import "github.com/go-martini/martini"

func main() {
  m := martini.Classic()
  m.Get("/", func() string {
    return "Hello, 세계!"
  })
  m.Run()
}
~~~

Then install the Martini package (**go 1.1** and greater is required):
마티니 패키지를 인스톨 합니다. (**go 1.1** 혹은 그 이상 버젼 필요):
~~~
go get github.com/go-martini/martini
~~~

Then run your server:
이제 서버를 돌려 봅시다:
~~~
go run server.go
~~~

You will now have a Martini webserver running on `localhost:3000`.
마티니 웹서버가 `localhost:3000`에서 돌아가고 있는 것을 확인하실 수 있을 겁니다.

## 도움이 필요하다면?

Join the [Mailing list](https://groups.google.com/forum/#!forum/martini-go)
[메일링 리스트](https://groups.google.com/forum/#!forum/martini-go)에 가입해 주세요

Watch the [Demo Video](http://martini.codegangsta.io/#demo)
[데모 비디오](http://martini.codegangsta.io/#demo)도 있어요.

Ask questions on Stackoverflow using the [martini tag](http://stackoverflow.com/questions/tagged/martini)
혹은 Stackoverflow에 [마티니 태크](http://stackoverflow.com/questions/tagged/martini)를 이용해서 물어봐 주세요

문제는 전부다 영어로 되어 있다는 건데요 -_-;;;
나는 한글 아니면 보기다 싫어! 이런 분들은 아래 링크를 참조하세요
- [golang-korea](https://code.google.com/p/golang-korea/)
- 이 문서 번역가([RexK](http://github.com/RexK))의 이메일로 연락주세요.

## 주요기능
* Extremely simple to use.
* 사용의 간편함
* Non-intrusive design.
* 비간섭(Non-intrusive) 디자인
* Plays nice with other Golang packages.
* 다른 Golang 패키지들과 잘 어울립니다.
* Awesome path matching and routing.
* 끝내주는 경로 매칭과 라우팅.
* Modular design - Easy to add functionality, easy to rip stuff out.
* 모듈 형 디자인 - 기능추가 쉽고, 코드 꺼내오기도 쉬움.
* Lots of good handlers/middlewares to use.
* 쓸모있는 핸들러와 미들웨어가 많음.
* Great 'out of the box' feature set.
* 훌률한 패키지화(out of the box) 기능들
* **Fully compatible with the [http.HandlerFunc](http://godoc.org/net/http#HandlerFunc) interface.**
* **[http.HandlerFunc](http://godoc.org/net/http#HandlerFunc) 인터페이스와 호환율 100%**

## 미들웨어(Middleware)
For more middleware and functionality, check out the repositories in the  [martini-contrib](https://github.com/martini-contrib) organization.
미들웨어들과 추가기능들은 [martini-contrib](https://github.com/martini-contrib)에서 확인해 주세요.

## 목차
* [Classic Martini](#classic-martini)
  * [Handlers](#handlers)
  * [Routing](#routing)
  * [Services](#services)
  * [Serving Static Files](#serving-static-files)
* [Middleware Handlers](#middleware-handlers)
  * [Next()](#next)
* [Martini Env](#martini-env)
* [FAQ](#faq)

## Classic Martini
To get up and running quickly, [martini.Classic()](http://godoc.org/github.com/go-martini/martini#Classic) provides some reasonable defaults that work well for most web applications:
마티니를 쉽고 빠르게 이용하시려면, [martini.Classic()](http://godoc.org/github.com/go-martini/martini#Classic)를 이용해 보세요. 보통 웹애플리케이션에서 사용하는 설정들이 이미 포함되어 있습니다.
~~~ go
  m := martini.Classic()
  // ... 미들웨어와 라우팅 설정은 이곳에 오면 작성하면 됩니다.
  m.Run()
~~~

Below is some of the functionality [martini.Classic()](http://godoc.org/github.com/go-martini/martini#Classic) pulls in automatically:
아래는 [martini.Classic()](http://godoc.org/github.com/go-martini/martini#Classic)의 자동으로 장착하는 기본 기능들입니다.
  * Request/Response Logging - [martini.Logger](http://godoc.org/github.com/go-martini/martini#Logger)
  * Request/Response 로그 기능 - [martini.Logger](http://godoc.org/github.com/go-martini/martini#Logger)
  * Panic Recovery - [martini.Recovery](http://godoc.org/github.com/go-martini/martini#Recovery)
  * 패닉 리커버리 (Panic Recovery) - [martini.Recovery](http://godoc.org/github.com/go-martini/martini#Recovery)
  * Static File serving - [martini.Static](http://godoc.org/github.com/go-martini/martini#Static)
  * 정적 파일 서빙
  * Routing - [martini.Router](http://godoc.org/github.com/go-martini/martini#Router)

### 핸들러(Handlers)
Handlers are the heart and soul of Martini. A handler is basically any kind of callable function:
핸들러(Handlers)는 마티니의 핵심입니다. 핸들러는 기본적으로 실행 가능한 모든형태의 함수들입니다.
~~~ go
m.Get("/", func() {
  println("hello 세계")
})
~~~

#### 반환 값
If a handler returns something, Martini will write the result to the current [http.ResponseWriter](http://godoc.org/net/http#ResponseWriter) as a string:
핸들러가 반환을 하는 함수라면, 마티니는 반환 값을 [http.ResponseWriter](http://godoc.org/net/http#ResponseWriter)에 입력 할 것입니다.
~~~ go
m.Get("/", func() string {
  return "hello 세계" // HTTP 200 : "hello 세계"
})
~~~

You can also optionally return a status code:
원하신다면, 상태코드도 함께 반화 할 수 있습니다.
~~~ go
m.Get("/", func() (int, string) {
  return 418, "난 주전자야!" // HTTP 418 : "난 주전자야!"
})
~~~

#### 서비스 주입(Service Injection)
Handlers are invoked via reflection. Martini makes use of *Dependency Injection* to resolve dependencies in a Handlers argument list. **This makes Martini completely  compatible with golang's `http.HandlerFunc` interface.**
핸들러들은 리플렉션을 통해 호출됩니다. 마티니는 *의존성 주입*을 이용해서 핸들러의 인수들을 주입합니다. **이것이 마티니를 `http.HandlerFunc` 인터페이스와 100% 호환할 수 있게 해줍니다.**

If you add an argument to your Handler, Martini will search its list of services and attempt to resolve the dependency via type assertion:
핸들러의 인수를 입력했다면, 마티니가 서비스 리스트를 살펴본 후 타입확인(type assertion)을 통해 의존성을 해결을 시도 할 것입니다.
~~~ go
m.Get("/", func(res http.ResponseWriter, req *http.Request) { // res와 req는 마티니에 의해 주입되었다.
  res.WriteHeader(200) // HTTP 200
})
~~~

The following services are included with [martini.Classic()](http://godoc.org/github.com/go-martini/martini#Classic):
아래 서비스들은 [martini.Classic()](http://godoc.org/github.com/go-martini/martini#Classic):에 포함되어 있습니다.
  * [*log.Logger](http://godoc.org/log#Logger) - Global logger for Martini.
  * [*log.Logger](http://godoc.org/log#Logger) - 마티니의 글러벌(전역) 로그.
  * [martini.Context](http://godoc.org/github.com/go-martini/martini#Context) - http request context.
  * [martini.Context](http://godoc.org/github.com/go-martini/martini#Context) - http 요청 컨텍스트.
  * [martini.Params](http://godoc.org/github.com/go-martini/martini#Params) - `map[string]string` of named params found by route matching.
  * [martini.Params](http://godoc.org/github.com/go-martini/martini#Params) - 루트 매칭으로 찾은 인자를 `map[string]string`으로 변형.
  * [martini.Routes](http://godoc.org/github.com/go-martini/martini#Routes) - 루트 도우미 서미스.
  * [http.ResponseWriter](http://godoc.org/net/http/#ResponseWriter) - http Response writer interface.
  * [http.ResponseWriter](http://godoc.org/net/http/#ResponseWriter) - http Response writer 인터페이스.
  * [*http.Request](http://godoc.org/net/http/#Request) - http 요구.

### 라우팅(Routing)
In Martini, a route is an HTTP method paired with a URL-matching pattern.
Each route can take one or more handler methods:
마티니에서 루트는 HTTP 메소드와 URL매칭 패턴의 패어이다. 각 루트는 하나 혹은 그 이상의 핸들러 메소드를 가질 수 있다.
~~~ go
m.Get("/", func() {
  // 뭘 좀 보여줘 봐
})

m.Patch("/", func() {
  // 업데이트 좀 해
})

m.Post("/", func() {
  // 뭘 좀 만들어봐
})

m.Put("/", func() {
  // 뭘 좀 교환해봐
})

m.Delete("/", func() {
  // 없애버려!
})

m.Options("/", func() {
  // http 옵션 메소드
})

m.NotFound(func() {
  // 404 해결하기
})
~~~

Routes are matched in the order they are defined. The first route that
matches the request is invoked.
루트들은 정의된 순서대로 매칭된다. 들어온 요그에 첫번째 매칭된 루트가 호출된다.

Route patterns may include named parameters, accessible via the [martini.Params](http://godoc.org/github.com/go-martini/martini#Params) service:
루트 패턴은 [martini.Params](http://godoc.org/github.com/go-martini/martini#Params) service로 액세스 가능한 인자들을 포함하기도 한다:
~~~ go
m.Get("/hello/:name", func(params martini.Params) string {
  return "Hello " + params["name"]
})
~~~

Routes can be matched with regular expressions and globs as well:
~~~ go
m.Get("/hello/**", func(params martini.Params) string {
  return "Hello " + params["_1"]
})
~~~

Route handlers can be stacked on top of each other, which is useful for things like authentication and authorization:
~~~ go
m.Get("/secret", authorize, func() {
  // this will execute as long as authorize doesn't write a response
})
~~~

Route groups can be added too using the Group method.
~~~ go
m.Group("/books", func(r martini.Router) {
    r.Get("/:id", GetBooks)
    r.Post("/new", NewBook)
    r.Put("/update/:id", UpdateBook)
    r.Delete("/delete/:id", DeleteBook)
})
~~~

Just like you can pass middlewares to a handler you can pass middlewares to groups.
~~~ go
m.Group("/books", func(r martini.Router) {
    r.Get("/:id", GetBooks)
    r.Post("/new", NewBook)
    r.Put("/update/:id", UpdateBook)
    r.Delete("/delete/:id", DeleteBook)
}, MyMiddleware1, MyMiddleware2)
~~~

### Services
Services are objects that are available to be injected into a Handler's argument list. You can map a service on a *Global* or *Request* level.

#### Global Mapping
A Martini instance implements the inject.Injector interface, so mapping a service is easy:
~~~ go
db := &MyDatabase{}
m := martini.Classic()
m.Map(db) // the service will be available to all handlers as *MyDatabase
// ...
m.Run()
~~~

#### Request-Level Mapping
Mapping on the request level can be done in a handler via [martini.Context](http://godoc.org/github.com/go-martini/martini#Context):
~~~ go
func MyCustomLoggerHandler(c martini.Context, req *http.Request) {
  logger := &MyCustomLogger{req}
  c.Map(logger) // mapped as *MyCustomLogger
}
~~~

#### Mapping values to Interfaces
One of the most powerful parts about services is the ability to map a service to an interface. For instance, if you wanted to override the [http.ResponseWriter](http://godoc.org/net/http#ResponseWriter) with an object that wrapped it and performed extra operations, you can write the following handler:
~~~ go
func WrapResponseWriter(res http.ResponseWriter, c martini.Context) {
  rw := NewSpecialResponseWriter(res)
  c.MapTo(rw, (*http.ResponseWriter)(nil)) // override ResponseWriter with our wrapper ResponseWriter
}
~~~

### Serving Static Files
A [martini.Classic()](http://godoc.org/github.com/go-martini/martini#Classic) instance automatically serves static files from the "public" directory in the root of your server.
You can serve from more directories by adding more [martini.Static](http://godoc.org/github.com/go-martini/martini#Static) handlers.
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
[Context.Next()](http://godoc.org/github.com/go-martini/martini#Context) is an optional function that Middleware Handlers can call to yield the until after the other Handlers have been executed. This works really well for any operations that must happen after an http request:
~~~ go
// log before and after a request
m.Use(func(c martini.Context, log *log.Logger){
  log.Println("before a request")

  c.Next()
  
  log.Println("after a request")
})
~~~

## Martini Env

Some Martini handlers make use of the `martini.Env` global variable to provide special functionality for development environments vs production environments. It is reccomended that the `MARTINI_ENV=production` environment variable to be set when deploying a Martini server into a production environment.

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
* [cors](https://github.com/martini-contrib/cors) - Handler that enables CORS support.
* [oauth2](https://github.com/martini-contrib/oauth2) - Handler that provides OAuth 2.0 login for Martini apps. Google Sign-in, Facebook Connect and Github login is supported.

### How do I integrate with existing servers?

A Martini instance implements `http.Handler`, so it can easily be used to serve subtrees 
on existing Go servers. For example this is a working Martini app for Google App Engine:

~~~ go
package hello

import (
  "net/http"
  "github.com/go-martini/martini"
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
  log.Fatal(http.ListenAndServe(":8080", m))
~~~

### Live code reload?

[gin](https://github.com/codegangsta/gin) and [fresh](https://github.com/pilu/fresh) both live reload martini apps.

## Contributing
Martini is meant to be kept tiny and clean. Most contributions should end up in a repository in the [martini-contrib](https://github.com/martini-contrib) organization. If you do have a contribution for the core of Martini feel free to put up a Pull Request.

## About

Inspired by [express](https://github.com/visionmedia/express) and [sinatra](https://github.com/sinatra/sinatra)

Martini is obsessively designed by none other than the [Code Gangsta](http://codegangsta.io/)
