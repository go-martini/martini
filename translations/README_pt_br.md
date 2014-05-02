# Martini  [![wercker status](https://app.wercker.com/status/9b7dbc6e2654b604cd694d191c3d5487/s/master "wercker status")](https://app.wercker.com/project/bykey/9b7dbc6e2654b604cd694d191c3d5487)[![GoDoc](https://godoc.org/github.com/go-martini/martini?status.png)](http://godoc.org/github.com/go-martini/martini)

Martini é um poderoso pacote para escrever aplicações/serviços modulares em Golang..


## Vamos Começar

Após a instalação do Go e de configurar o [GOPATH](http://golang.org/doc/code.html#GOPATH), crie seu primeiro arquivo `.go`. Vamos chamá-lo de `server.go`.

~~~ go
package main

import "github.com/go-martini/martini"

func main() {
  m := martini.Classic()
  m.Get("/", func() string {
    return "Hello world!"
  })
  m.Run()
}
~~~

Então instale o pacote do Martini (É necessário **go 1.1** ou superior):
~~~
go get github.com/go-martini/martini
~~~

Então rode o servidor:
~~~
go run server.go
~~~

Agora você tem um webserver Martini rodando na porta `localhost:3000`.

## Pegue uma Ajuda

Assine a [Lista de email](https://groups.google.com/forum/#!forum/martini-go)

Veja o [Vídeo demonstrativo](http://martini.codegangsta.io/#demo)

Use a tag [martini tag](http://stackoverflow.com/questions/tagged/martini) para perguntas no Stackoverflow



## Caracteríticas
* Extrema simplicidade no uso.
* Não intrusivo ao desenho da app.
* Funciona bem com outros pacotes Golang.
* Impressionante caminho de correspondência e rotas.
* Design modular - Fácil para adicionar e remover funcionalidades.
* Muito bom no uso handlers/middlewares.
* Grandes caracteríticas inovadoras.
* **Completa compatibilidade com a interface [http.HandlerFunc](http://godoc.org/net/http#HandlerFunc).**

## Mais Middleware
Para mais middleware e funcionalidades, veja os repositórios em [martini-contrib](https://github.com/martini-contrib).

## Tabela de Conteudos
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
Para iniciar e rodar facilmente, [martini.Classic()](http://godoc.org/github.com/go-martini/martini#Classic) prove algumas ferramentas rasoáveis para maioria das aplicações web:
~~~ go
  m := martini.Classic()
  // ... middleware e rota aqui
  m.Run()
~~~

Algumas das funcionalidade que o [martini.Classic()](http://godoc.org/github.com/go-martini/martini#Classic) oferece automaticamente são:
  * Request/Response Logging - [martini.Logger](http://godoc.org/github.com/go-martini/martini#Logger)
  * Panic Recovery - [martini.Recovery](http://godoc.org/github.com/go-martini/martini#Recovery)
  * Servidor de arquivos státicos - [martini.Static](http://godoc.org/github.com/go-martini/martini#Static)
  * Rotas - [martini.Router](http://godoc.org/github.com/go-martini/martini#Router)

### Handlers
Handlers são o coração e a alma do Martini. Um handler é basicamente qualquer função que pode ser chamada:
~~~ go
m.Get("/", func() {
  println("hello world")
})
~~~

#### Retorno de Valores
Se um handler retornar alguma coisa, Martini ira escrever o resultado atual [http.ResponseWriter](http://godoc.org/net/http#ResponseWriter) que é uma string:
~~~ go
m.Get("/", func() string {
  return "hello world" // HTTP 200 : "hello world"
})
~~~

Você também pode retornar o código de status:
~~~ go
m.Get("/", func() (int, string) {
  return 418, "Eu sou um bule" // HTTP 418 : "Eu sou um bule"
})
~~~

#### Injeção de Serviços
Handlers são chamados via reflexão. Martini utiliza *Injeção de Dependencia* para resolver as dependencias nas listas de argumentos dos Handlers . **Isso faz Martini ser completamente compatível com a interface `http.HandlerFunc` do golang.**

Se você adicionar um argumento ao seu Handler, Martini ira procurar na sua lista de serviço e tentar resolver sua dependencia via tipo de afirmação:
~~~ go
m.Get("/", func(res http.ResponseWriter, req *http.Request) { // res e req são injetados pelo Martini
  res.WriteHeader(200) // HTTP 200
})
~~~

Os seguintes serviços são incluídos com [martini.Classic()](http://godoc.org/github.com/go-martini/martini#Classic):
  * [*log.Logger](http://godoc.org/log#Logger) - Log Global para Martini.
  * [martini.Context](http://godoc.org/github.com/go-martini/martini#Context) - http request context.
  * [martini.Params](http://godoc.org/github.com/go-martini/martini#Params) - `map[string]string` de nomes dos parâmetros buscados pela rota.
  * [martini.Routes](http://godoc.org/github.com/go-martini/martini#Routes) - Serviço de auxílio as rotas.
  * [http.ResponseWriter](http://godoc.org/net/http/#ResponseWriter) - http Response escreve a interface.
  * [*http.Request](http://godoc.org/net/http/#Request) - http Request.

### Rotas
No Martini, uma rota é um método HTTP emparelhado com um padrão de URL de correspondência.
Cada rota pode ter um ou mais métodos handler:
~~~ go
m.Get("/", func() {
  // mostra alguma coisa
})

m.Patch("/", func() {
  // altera alguma coisa
})

m.Post("/", func() {
  // cria alguma coisa
})

m.Put("/", func() {
  // sobrescreve alguma coisa
})

m.Delete("/", func() {
  // destrói alguma coisa
})

m.Options("/", func() {
  // opções do HTTP
})

m.NotFound(func() {
  // manipula 404
})
~~~

As rotas são combinadas na ordem em que são definidas. A primeira rota que corresponde a solicitação é chamada.

O padrão de rotas pode incluir parâmetros que podem ser acessados via [martini.Params](http://godoc.org/github.com/go-martini/martini#Params):
~~~ go
m.Get("/hello/:name", func(params martini.Params) string {
  return "Hello " + params["name"]
})
~~~

As rotas podem ser combinados com expressões regulares e globs:
~~~ go
m.Get("/hello/**", func(params martini.Params) string {
  return "Hello " + params["_1"]
})
~~~

Handlers de rota podem ser empilhados em cima uns dos outros, o que é útil para coisas como autenticação e autorização:
~~~ go
m.Get("/secret", authorize, func() {
  // Será executado quando authorize não escrever uma resposta
})
~~~

Grupos de rota podem ser adicionados usando o método Group.
~~~ go
m.Group("/books", func(r martini.Router) {
    r.Get("/:id", GetBooks)
    r.Post("/new", NewBook)
    r.Put("/update/:id", UpdateBook)
    r.Delete("/delete/:id", DeleteBook)
})
~~~

Assim como você pode passar middlewares para um manipulador você pode passar middlewares para grupos.
~~~ go
m.Group("/books", func(r martini.Router) {
    r.Get("/:id", GetBooks)
    r.Post("/new", NewBook)
    r.Put("/update/:id", UpdateBook)
    r.Delete("/delete/:id", DeleteBook)
}, MyMiddleware1, MyMiddleware2)
~~~

### Serviços
Serviços são objetos que estão disponíveis para ser injetado em uma lista de argumentos de Handler. Você pode mapear um serviço num nível *Global* ou *Request*.

#### Mapeamento Global
Um exemplo onde o Martini implementa a interface inject.Injector, então o mapeamento de um serviço é fácil:
~~~ go
db := &MyDatabase{}
m := martini.Classic()
m.Map(db) // o serviço estará disponível para todos os handlers *MyDatabase.
// ...
m.Run()
~~~

#### Mapeamento Request
Mapeamento do nível de request pode ser feito via handler através [martini.Context](http://godoc.org/github.com/go-martini/martini#Context):
~~~ go
func MyCustomLoggerHandler(c martini.Context, req *http.Request) {
  logger := &MyCustomLogger{req}
  c.Map(logger) // mapeamento é *MyCustomLogger
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
