<!--

TODO:
  - Proof reading
  - Check, if specific code parts like URL's are NOT translated
  - Check, if only the "Du"-form is used
  - Check the table of contents with the titles of each section

-->





# Martini  [![wercker status](https://app.wercker.com/status/9b7dbc6e2654b604cd694d191c3d5487/s/master "wercker status")](https://app.wercker.com/project/bykey/9b7dbc6e2654b604cd694d191c3d5487)[![GoDoc](https://godoc.org/github.com/go-martini/martini?status.png)](http://godoc.org/github.com/go-martini/martini)

Martini ist eine mächtiges Package zur schnellen Entwicklung von Webanwendungen/services in Golang. 

## Ein Projekt starten

Nach der Installation von Go und dem Einrichten des [GOPATH](http://golang.org/doc/code.html#GOPATH), erstelle deine erste `.go`-Datei. Speichere sie unter `server.go`.

~~~ go
package main

import "github.com/go-martini/martini"

func main() {
  m := martini.Classic()
  m.Get("/", func() string {
    return "Hallo Welt!"
  })
  m.Run()
}
~~~

Installiere anschließend das Martini Package (**Go 1.1** oder höher ist vorausgesetzt):
~~~
go get github.com/go-martini/martini
~~~

Starte den Server:
~~~
go run server.go
~~~

Der Martini Webserver ist nun unter `localhost:3000` erreichbar.

## Hilfe

Abboniere die [Mailing list](https://groups.google.com/forum/#!forum/martini-go)

Schaue das [Demovideo](http://martini.codegangsta.io/#demo)

Stelle Fragen auf Stackoverflow mit dem [Martini-Tag](http://stackoverflow.com/questions/tagged/martini)

GoDoc [Dokumentation](http://godoc.org/github.com/go-martini/martini)


## Eigenschaften
* Sehr einfach nutzbar
* Systemunabhängiges Design
* Einfach anwendbar mit Anderen Golang Packages
* TODO: Awesome path matching and routing.
* Modulares Design - einfaches Hinzufügen und Entfernen von Funktionen
* Eine Vielzahl von guten Handlern/Middlewares nutzbar
* TODO: Großer Funktionsumfang mitgeliefert // Great 'out of the box' feature set.
* **Voll kompatibel mit dem [http.HandlerFunc](http://godoc.org/net/http#HandlerFunc) Interface.**
* TODO: Standardmäßge Seitenübertragung //Default document serving (e.g., for serving AngularJS apps in HTML5 mode).

## Mehr Middleware
Mehr Informationen zur Middleware und Funktionalität finden Du in der [martini-contrib](https://github.com/martini-contrib) Repository.

## Inhaltsverzeichnis
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
Einen schnellen Start in ein Projekt ermöglicht [martini.Classic()](http://godoc.org/github.com/go-martini/martini#Classic), dessen Voreinstellungen sich für die meisten Webanwendungen eignen:
~~~ go
  m := martini.Classic()
  // ... Middleware und Routing hier einfügen
  m.Run()
~~~

TODO:
Below is some of the functionality [martini.Classic()](http://godoc.org/github.com/go-martini/martini#Classic) pulls in automatically:
  * Request/Response Logging - [martini.Logger](http://godoc.org/github.com/go-martini/martini#Logger)
  * Panic Recovery - [martini.Recovery](http://godoc.org/github.com/go-martini/martini#Recovery)
  * Static File serving - [martini.Static](http://godoc.org/github.com/go-martini/martini#Static)
  * Routing - [martini.Router](http://godoc.org/github.com/go-martini/martini#Router)

### Handlers
Handlers sind das Herz und die Seele von Martini. Ein Handler ist grundsätzlich jede Art von aufrufbaren Funktionen:
~~~ go
m.Get("/", func() {
  println("Hallo Welt")
})
~~~

#### Rückgabewerte
Wenn ein Handerl etwas zurückgibt, übergibt Martini den Wert an den aktuellen [http.ResponseWriter](http://godoc.org/net/http#ResponseWriter) in Form einen String:
~~~ go
m.Get("/", func() string {
  return "Hallo Welt" // HTTP 200 : "Hallo Welt"
})
~~~

Die Rückgabe eines Statuscode ist optional:
~~~ go
m.Get("/", func() (int, string) {
  return 418, "Ich bin eine Teekanne" // HTTP 418 : "Ich bin eine Teekanne"
})
~~~

#### Service Injection
Handler werden per Reflection aufgerufen. Martini macht Gebrauch von *Dependency Injection*, um Abhängigkeiten in der Argumentliste von Handlern aufzulösen. **Dies macht Martini komplett inkompatibel mit Golangs `http.HandlerFunc` Interface.**

Fügst Du einem Handler ein Argument hinzu, sucht Martini in seiner Liste von Services und versucht, die Abhängigkeiten via Type Assertion aufzulösen. 
~~~ go
m.Get("/", func(res http.ResponseWriter, req *http.Request) { // res und req wurden von Martini injiziert
  res.WriteHeader(200) // HTTP 200
})
~~~

Die Folgenden Services sind Bestandteil von [martini.Classic()](http://godoc.org/github.com/go-martini/martini#Classic):
  * [*log.Logger](http://godoc.org/log#Logger) - Globaler Logger für Martini.
  * [martini.Context](http://godoc.org/github.com/go-martini/martini#Context) - http request context.
  * [martini.Params](http://godoc.org/github.com/go-martini/martini#Params) - `map[string]string` von benannten Parametern, welche durch route matching gefunden wurden.
  * [martini.Routes](http://godoc.org/github.com/go-martini/martini#Routes) - Route helper service.
  * [http.ResponseWriter](http://godoc.org/net/http/#ResponseWriter) - http Response writer interface.
  * [*http.Request](http://godoc.org/net/http/#Request) - http Request.

### Routing
Eine Route ist in Martini eine HTTP-Methode gepaart mit einem TODO: URL-matching pattern. Jede Route kann ein oder mehrere Handler-Methoden übernehmen:
~~~ go
m.Get("/", func() {
  // zeige etwas
})

m.Patch("/", func() {
  // aktualisiere etwas
})

m.Post("/", func() {
  // erstelle etwas
})

m.Put("/", func() {
  // ersetzte etwas
})

m.Delete("/", func() {
  // Lösche etwas
})

m.Options("/", func() {
  // HTTP Optionen
})

m.NotFound(func() {
  // behandle 404 Statuscode
})
~~~

Routen werden in der Reihenfolge, in welcher sie definiert wurden, zugeordnet. Die erste zugeordnete Route wird daraufhin aufgerufen.  

Routen-Muster enhalten gegebenenfalls benannte Parameter, die über den [martini.Params](http://godoc.org/github.com/go-martini/martini#Params) Service aufrufbar sind:
~~~ go
m.Get("/hello/:name", func(params martini.Params) string {
  return "Hallo " + params["name"]
})
~~~

Routen können mit Globs versehen werden:
~~~ go
m.Get("/hello/**", func(params martini.Params) string {
  return "Hallo " + params["_1"]
})
~~~

Reguläre Ausdrücke sind ebenfalls möglich:
~~~go
m.Get("/hello/(?P<name>[a-zA-Z]+)", func(params martini.Params) string {
  return fmt.Sprintf ("Hallo %s", params["name"])
})
~~~
Weitere Informationen zum Syntax regulärer Ausdrücke findest Du in der [Go Dokumentation](http://golang.org/pkg/regexp/syntax/).

Routen-Handler können auch in einander verschachtelt werden. Dies ist bei der Authentifizierung und Berechtigungen nützlich.
~~~ go
m.Get("/secret", authorize, func() {
  // dies wird ausgeführt, solange authorize nichts zurückgibt
})
~~~

Routen-Gruppen können durch die Group-Methode hinzugefügt werden.
~~~ go
m.Group("/books", func(r martini.Router) {
    r.Get("/:id", GetBooks)
    r.Post("/new", NewBook)
    r.Put("/update/:id", UpdateBook)
    r.Delete("/delete/:id", DeleteBook)
})
~~~

Wie bei Handlern können Middlewares auch Gruppen übergeben werden.
~~~ go
m.Group("/books", func(r martini.Router) {
    r.Get("/:id", GetBooks)
    r.Post("/new", NewBook)
    r.Put("/update/:id", UpdateBook)
    r.Delete("/delete/:id", DeleteBook)
}, MyMiddleware1, MyMiddleware2)
~~~

### Services
Services sind Okjekte, welche der Argumentliste von Handlern beigefügt werden können.
Du kannst einen Service der *Global* oder *Request* Ebene zuordnen.

#### Global Mapping
Eine Martini-Instanz implementiert das inject.Injector interface, sodass ein Service leicht zugeordnet werden kann:
~~~ go
db := &MyDatabase{}
m := martini.Classic()
m.Map(db) // Der Service ist allen Handlern unter *MyDatabase verfügbar
// ...
m.Run()
~~~

#### Request-Level Mapping
Das Zuordnen auf der Request-Ebene kann in einem Handler via  [martini.Context](http://godoc.org/github.com/go-martini/martini#Context) realisiert werden:
~~~ go
func MyCustomLoggerHandler(c martini.Context, req *http.Request) {
  logger := &MyCustomLogger{req}
  c.Map(logger) // Zugeordnet als *MyCustomLogger
}
~~~

#### Mapping values to Interfaces
Einer der mächtigsten Aspekte von Services ist dessen Fähigkeit, einen Service einem Interface zuzuordnen. TODO: For instance, if you wanted to override the [http.ResponseWriter](http://godoc.org/net/http#ResponseWriter) with an object that wrapped it and performed extra operations, you can write the following handler:
~~~ go
func WrapResponseWriter(res http.ResponseWriter, c martini.Context) {
  rw := NewSpecialResponseWriter(res)
  c.MapTo(rw, (*http.ResponseWriter)(nil)) // override ResponseWriter with our wrapper ResponseWriter
}
~~~

### Serving Static Files
Eine [martini.Classic()](http://godoc.org/github.com/go-martini/martini#Classic) Instanz übertragt automatisch statische Dateien aus dem "pulbic"-Ordner im Stammverzeichnis Deines Servers. Dieses Verhalten lässt sirch durch weitere [martini.Static](http://godoc.org/github.com/go-martini/martini#Static) Handler auf andere Verzeichnisse übertragen.
~~~ go
m.Use(martini.Static("assets")) // Überträgt auch vom "assets"-Verzeichnis
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
* [vauth](https://github.com/rafecolton/vauth) - Handlers for vender webhook authentication (currently GitHub and TravisCI)

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

### Wie ändere ich den Port/Host?

Martinis `Run` Funktion sucht automatisch nach den PORT und HOST Umgebungsvariablen, um diese zu nutzen. Andernfalls ist localhost:3000 voreingestellt.
Für mehr Flexibilität über den Port und den Host nutze stattdessen die `martini.RunOnAddr` Funktion.

~~~ go
  m := martini.Classic()
  // ...
  log.Fatal(m.RunOnAddr(":8080"))
~~~

### Automatisches Aktualisieren?

[Gin](https://github.com/codegangsta/gin) und [Fresh](https://github.com/pilu/fresh) aktualisieren Martini-Apps live.

## Beitragen
Martinis Grundsatz ist Minimalismus und sauberer Code. Die meisten Beiträge sollten sich in der [martini-contrib](https://github.com/martini-contrib) Repository wiederfinden. Beinhaltet Dein Beitrag Veränderungen am Kern von Martini, zögere nicht, einen Pull Request zu machen.

## Über das Projekt

Inspiriert von [Express](https://github.com/visionmedia/express) und [Sinatra](https://github.com/sinatra/sinatra)

Martini wird leidenschaftlich entwickelt von Niemand gerigeren als dem [Code Gangsta](http://codegangsta.io/)
