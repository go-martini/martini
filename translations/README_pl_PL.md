# Martini  [![wercker status](https://app.wercker.com/status/9b7dbc6e2654b604cd694d191c3d5487/s/master "wercker status")](https://app.wercker.com/project/bykey/9b7dbc6e2654b604cd694d191c3d5487)[![GoDoc](https://godoc.org/github.com/go-martini/martini?status.png)](http://godoc.org/github.com/go-martini/martini)

Martini to solidny framework umożliwiający sprawne tworzenie modularnych aplikacji internetowych i usług sieciowych w języku Go.

## Pierwsze kroki

Po zakończonej instalacji Go i ustawieniu zmiennej [GOPATH](http://golang.org/doc/code.html#GOPATH), utwórz swój pierwszy plik `.go`. Nazwijmy go `server.go`.

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

Następnie zainstaluj pakiet Martini (środowisko **go** w wersji **1.1** lub nowszej jest wymagane):
~~~
go get github.com/go-martini/martini
~~~

Uruchom serwer:
~~~
go run server.go
~~~

W tym momencie webserwer Martini jest uruchomiony na `localhost:3000`.

## Uzyskiwanie pomocy

Dołącz do [grup dyskusyjnych](https://groups.google.com/forum/#!forum/martini-go)

Obejrzyj przygotowane [demo](http://martini.codegangsta.io/#demo)

Zadawaj pytania na Stackoverflow dodając [tag martini](http://stackoverflow.com/questions/tagged/martini)

GoDoc [dokumentacja](http://godoc.org/github.com/go-martini/martini)


## Cechy frameworka
* Bardzo prosty w użyciu.
* Posiada niewymagającą ingerencji budowę.
* Łatwo integruje się z innymi pakietami w języku Go.
* Sprawnie dopasowuje ścieżki i routing.
* Modularny projekt - łatwo dodać funkcję i łatwo usunąć.
* Bogate zasoby handlerów i middleware'ów do wykorzystania.
* Spora część funkcji działa 'z paczki'.
* **W pełni kompatybilny z interfejsem [http.HandlerFunc](http://godoc.org/net/http#HandlerFunc).**
* Umożliwia serwowanie domyślnych stron (np. przy serwowaniu aplikacji napisanych w AngularJS w trybie HTML5).

## Więcej middleware'ów
W celu uzyskania więcej informacji o middleware'ach i ich możliwościach, przejrzyj repozytoria należące do organizacji [martini-contrib](https://github.com/martini-contrib).

## Spis treści
* [Domyślna konfiguracja (Martini Classic)](#classic-martini)
  * [Handlery](#handlers)
  * [Routing](#routing)
  * [Usługi](#services)
  * [Serwowanie plików statycznych](#serving-static-files)
* [Handlery middleware'ów](#middleware-handlers)
  * [Next()](#next)
* [Zmienne środowiskowe Martini](#martini-env)
* [FAQ](#faq)

## Domyślna konfiguracja (Martini Classic)
Martini pozwala bardzo szybko uruchomić webserver korzystając przy tym z [martini.Classic()](http://godoc.org/github.com/go-martini/martini#Classic). Standardowo Classic dostarcza domyślne ustawienia, które z powodzeniem pozwolą nam uruchomić wiele aplikacji internetowych:
~~~ go
  m := martini.Classic()
  // ... miejsce na middleware'y i routing
  m.Run()
~~~

Poniżej wymieniono kilka funkcji [martini.Classic()](http://godoc.org/github.com/go-martini/martini#Classic) dostarczanych automatycznie:
  * Logowanie żądań/odpowiedzi - [martini.Logger](http://godoc.org/github.com/go-martini/martini#Logger)
  * Panic Recovery - [martini.Recovery](http://godoc.org/github.com/go-martini/martini#Recovery)
  * Serwowanie plików statycznych - [martini.Static](http://godoc.org/github.com/go-martini/martini#Static)
  * Routing - [martini.Router](http://godoc.org/github.com/go-martini/martini#Router)

### Handlery
Handlery to serce i dusza Martini. Handlerem można nazwać każdą funkcję postaci:
~~~ go
m.Get("/", func() {
  println("hello world")
})
~~~

#### Wartości zwracane
Jeśli handler zwróci wartość, Martini przekaże ją do bieżącego [http.ResponseWriter](http://godoc.org/net/http#ResponseWriter) jako łańcuch znaków:
~~~ go
m.Get("/", func() string {
  return "hello world" // HTTP 200 : "hello world"
})
~~~

Opcjonalnie można zwrócić także status HTTP:
~~~ go
m.Get("/", func() (int, string) {
  return 418, "i'm a teapot" // HTTP 418 : "i'm a teapot"
})
~~~

#### Wstrzykiwanie usług
Handlery są wywoływane przez refleksję. Martini korzysta z *wstrzykiwania zależności* w celu rozwiązania tych, które występują na liście argumentów handlera. **To sprawia, że Martini jest w pełni zgodny z interfejsem `http.HandlerFunc`.**

Jeśli dodasz argument do handlera, Martini przeszuka swoja listę usług i spróbuje dopasować zależność na podstawie asercji typów:
~~~ go
m.Get("/", func(res http.ResponseWriter, req *http.Request) { // res i req są wstrzykiwane przez Martini
  res.WriteHeader(200) // HTTP 200
})
~~~

Następujące usługi są dostarczane razem z [martini.Classic()](http://godoc.org/github.com/go-martini/martini#Classic):
  * [*log.Logger](http://godoc.org/log#Logger) - Globalny logger dla Martini.
  * [martini.Context](http://godoc.org/github.com/go-martini/martini#Context) - kontekst żądania HTTP.
  * [martini.Params](http://godoc.org/github.com/go-martini/martini#Params) - `map[string]string` przechowująca nazwane parametry, znalezione podczas dopasowywania _routes_.
  * [martini.Routes](http://godoc.org/github.com/go-martini/martini#Routes) - usługa wspierająca _routes_.
  * [martini.Route](http://godoc.org/github.com/go-martini/martini#Route) - bieżacy aktywny _route_.
  * [http.ResponseWriter](http://godoc.org/net/http/#ResponseWriter) - interfejs zapisu odpowiedzi HTTP.
  * [*http.Request](http://godoc.org/net/http/#Request) - żądanie HTTP.

### Routing
W Martini, jako _route_ należy rozumieć metodę HTTP skojarzoną ze wzorcem dopasowującym adres URL.
Każdy wzorzec może być skojarzony z jedną lub wiecęj metodą handlera:
~~~ go
m.Get("/", func() {
  // wyświetl coś
})

m.Patch("/", func() {
  // zaaktualizuj coś
})

m.Post("/", func() {
  // utwórz coś
})

m.Put("/", func() {
  // zamień coś
})

m.Delete("/", func() {
  // zniszcz coś
})

m.Options("/", func() {
  // opcje HTTP
})

m.NotFound(func() {
  // obsłuż 404
})
~~~

Routes are matched in the order they are defined. The first route that
matches the request is invoked.

Route patterns may include named parameters, accessible via the [martini.Params](http://godoc.org/github.com/go-martini/martini#Params) service:
~~~ go
m.Get("/hello/:name", func(params martini.Params) string {
  return "Hello " + params["name"]
})
~~~

Routes can be matched with globs:
~~~ go
m.Get("/hello/**", func(params martini.Params) string {
  return "Hello " + params["_1"]
})
~~~

Regular expressions can be used as well:
~~~go
m.Get("/hello/(?P<name>[a-zA-Z]+)", func(params martini.Params) string {
  return fmt.Sprintf ("Hello %s", params["name"])
})
~~~
Take a look at the [Go documentation](http://golang.org/pkg/regexp/syntax/) for more info about regular expressions syntax .

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

#### Serving a Default Document
You can specify the URL of a local file to serve when the requested URL is not
found. You can also specify an exclusion prefix so that certain URLs are ignored.
This is useful for servers that serve both static files and have additional
handlers defined (e.g., REST API). When doing so, it's useful to define the
static handler as a part of the NotFound chain.

The following example serves the `/index.html` file whenever any URL is
requested that does not match any local file and does not start with `/api/v`:
~~~ go
static := martini.Static("assets", martini.StaticOptions{Fallback: "/index.html", Exclude: "/api/v"})
m.NotFound(static, http.NotFound)
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

Some Martini handlers make use of the `martini.Env` global variable to provide special functionality for development environments vs production environments. It is recommended that the `MARTINI_ENV=production` environment variable to be set when deploying a Martini server into a production environment.

## FAQ

### Where do I find middleware X?

Start by looking in the [martini-contrib](https://github.com/martini-contrib) projects. If it is not there feel free to contact a martini-contrib team member about adding a new repo to the organization.

* [acceptlang](https://github.com/martini-contrib/acceptlang) - Handler for parsing the `Accept-Language` HTTP header.
* [accessflags](https://github.com/martini-contrib/accessflags) - Handler to enable Access Control.
* [auth](https://github.com/martini-contrib/auth) - Handlers for authentication.
* [binding](https://github.com/martini-contrib/binding) - Handler for mapping/validating a raw request into a structure.
* [cors](https://github.com/martini-contrib/cors) - Handler that enables CORS support.
* [csrf](https://github.com/martini-contrib/csrf) - CSRF protection for applications
* [encoder](https://github.com/martini-contrib/encoder) - Encoder service for rendering data in several formats and content negotiation.
* [gzip](https://github.com/martini-contrib/gzip) - Handler for adding gzip compress to requests
* [gorelic](https://github.com/martini-contrib/gorelic) - NewRelic middleware
* [logstasher](https://github.com/martini-contrib/logstasher) - Middleware that prints logstash-compatiable JSON 
* [method](https://github.com/martini-contrib/method) - HTTP method overriding via Header or form fields.
* [oauth2](https://github.com/martini-contrib/oauth2) - Handler that provides OAuth 2.0 login for Martini apps. Google Sign-in, Facebook Connect and Github login is supported.
* [permissions2](https://github.com/xyproto/permissions2) - Handler for keeping track of users, login states and permissions.
* [render](https://github.com/martini-contrib/render) - Handler that provides a service for easily rendering JSON and HTML templates.
* [secure](https://github.com/martini-contrib/secure) - Implements a few quick security wins.
* [sessions](https://github.com/martini-contrib/sessions) - Handler that provides a Session service.
* [sessionauth](https://github.com/martini-contrib/sessionauth) - Handler that provides a simple way to make routes require a login, and to handle user logins in the session
* [strict](https://github.com/martini-contrib/strict) - Strict Mode 
* [strip](https://github.com/martini-contrib/strip) - URL Prefix stripping.
* [staticbin](https://github.com/martini-contrib/staticbin) - Handler for serving static files from binary data
* [throttle](https://github.com/martini-contrib/throttle) - Request rate throttling middleware.
* [vauth](https://github.com/rafecolton/vauth) - Handlers for vender webhook authentication (currently GitHub and TravisCI)
* [web](https://github.com/martini-contrib/web) - hoisie web.go's Context

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
To have more flexibility over port and host, use the `martini.RunOnAddr` function instead.

~~~ go
  m := martini.Classic()
  // ...
  log.Fatal(m.RunOnAddr(":8080"))
~~~

### Live code reload?

[gin](https://github.com/codegangsta/gin) and [fresh](https://github.com/pilu/fresh) both live reload martini apps.

## Contributing
Martini is meant to be kept tiny and clean. Most contributions should end up in a repository in the [martini-contrib](https://github.com/martini-contrib) organization. If you do have a contribution for the core of Martini feel free to put up a Pull Request.

## About

Inspired by [express](https://github.com/visionmedia/express) and [sinatra](https://github.com/sinatra/sinatra)

Martini is obsessively designed by none other than the [Code Gangsta](http://codegangsta.io/)
