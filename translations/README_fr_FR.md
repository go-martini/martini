# Martini  [![wercker status](https://app.wercker.com/status/9b7dbc6e2654b604cd694d191c3d5487/s/master "wercker status")](https://app.wercker.com/project/bykey/9b7dbc6e2654b604cd694d191c3d5487)[![GoDoc](https://godoc.org/github.com/go-martini/martini?status.png)](http://godoc.org/github.com/go-martini/martini)

Martini est un puissant paquet pour développer rapidement des applications web ou services en Golang.


## Pour commencer

Après avoir installé Go et configuré le chemin d'accès pour [GOPATH](http://golang.org/doc/code.html#GOPATH), créer votre premier fichier '.go'. Nous l'appellerons 'server.go'.

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

Ensuite installer le paquet Martini (**go 1.1** ou supérieur est requis):

~~~
go get github.com/go-martini/martini
~~~

Ensuite démarrer votre serveur:
~~~
go run server.go
~~~

Vous aurez maintenant un serveur web Martini disponible à l'adresse `localhost:3000`.

## Besoin d'aide

Souscrivez à la [Liste d'emails](https://groups.google.com/forum/#!forum/martini-go)

Regarder les vidéos [Demo en vidéo](http://martini.codegangsta.io/#demo)

Poser vos questions sur Stackoverflow en utilisant le tag [martini](http://stackoverflow.com/questions/tagged/martini)

La documentation GoDoc [documentation](http://godoc.org/github.com/go-martini/martini)


## Caractéristiques
* Extrêmement simple à utiliser.
* Pas de design intrusif.
* Peut s'utiliser avec les autres paquets Golang.
* Chemin d'accès et routeur disponibles.
* Complètement modulable, ajout/retrait de fonctionnalité.
* Beaucoup de bon handlers/middlewares à utiliser
* Prêt à utiliser.
* **Entièrement compatible avec l'interface [http.HandlerFunc](http://godoc.org/net/http#HandlerFunc).**

## Plus de Middleware
Pour plus de middleware et de fonctionnalités, consulter le dépôt [martini-contrib](https://github.com/martini-contrib).

## Table des matières (index)
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
Pour commencer plus rapidement, [martini.Classic()](http://godoc.org/github.com/go-martini/martini#Classic) donne beaucoup d'informations par default qui fonctionne pour la plupart des applications web:

~~~ go
  m := martini.Classic()
  // ... middleware and routing goes here
  m.Run()
~~~

En dessous quelques fonctionnalités [martini.Classic()](http://godoc.org/github.com/go-martini/martini#Classic) récupèrera automatiquement:
  * Request/Response Logging - [martini.Logger](http://godoc.org/github.com/go-martini/martini#Logger)
  * Panic Recovery - [martini.Recovery](http://godoc.org/github.com/go-martini/martini#Recovery)
  * Static File serving - [martini.Static](http://godoc.org/github.com/go-martini/martini#Static)
  * Routing - [martini.Router](http://godoc.org/github.com/go-martini/martini#Router)

### Handlers
Les Handlers sont le coeur et l'âme de Martini. Un handler est basiquement n'importe quel fonction utilisable.
~~~ go
m.Get("/", func() {
  println("hello world")
})
~~~

#### Return Values
Si un handler retourne quelque chose, Martini écrira le résultat courant [http.ResponseWriter](http://godoc.org/net/http#ResponseWriter) comme une chaîne de caractères (string):
~~~ go
m.Get("/", func() string {
  return "hello world" // HTTP 200 : "hello world"
})
~~~
Vous pouvez aussi renvoyer un code statut optionnellement:
~~~ go
m.Get("/", func() (int, string) {
  return 418, "i'm a teapot" // HTTP 418 : "i'm a teapot"
})
~~~

#### Service Injection
Les Handlers sont appelés via réflexion. Martini utilise "l'injection par dépendance" pour résoudre les dépendances des Handlers dans la liste d'arguments. **Cela permet à Martini d'être complètement compatible avec l'interface golang `http.HandlerFunc`.**

Si vous ajoutez un argument à votre Handler, Martini ira chercher la liste des services et essayera de résoudre les dépendances selon le type d'assertion:
~~~ go
m.Get("/", func(res http.ResponseWriter, req *http.Request) { // res and req are injected by Martini
  res.WriteHeader(200) // HTTP 200
})
~~~
Les services suivant sont inclus avec [martini.Classic()](http://godoc.org/github.com/go-martini/martini#Classic):
  * [*log.Logger](http://godoc.org/log#Logger) - Global logger for Martini.
  * [martini.Context](http://godoc.org/github.com/go-martini/martini#Context) - http request context.
  * [martini.Params](http://godoc.org/github.com/go-martini/martini#Params) - `map[string]string` of named params found by route matching.
  * [martini.Routes](http://godoc.org/github.com/go-martini/martini#Routes) - Route helper service.
  * [http.ResponseWriter](http://godoc.org/net/http/#ResponseWriter) - http Response writer interface.
  * [*http.Request](http://godoc.org/net/http/#Request) - http Request.

### Routeur
Dans Martini, un chemin est une méthode HTTP lié à un modèle d'adresse URL.
Chaque chemin peut avoir un seul ou plusieurs méthodes d'handler:
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
Les chemins sont assortis dans l'ordre qu'ils sont définis. Le premier chemin qui correspond à la requête sera invoqué.


Les chemins peuvent inclure des paramètres nommés, accessibles avec le service [martini.Params](http://godoc.org/github.com/go-martini/martini#Params):
~~~ go
m.Get("/hello/:name", func(params martini.Params) string {
  return "Hello " + params["name"]
})
~~~

Les chemins peuvent correspondre à des globs:
~~~ go
m.Get("/hello/**", func(params martini.Params) string {
  return "Hello " + params["_1"]
})
~~~
Les expressions régulières peuvent aussi être utilisées:
~~~go
m.Get("/hello/(?P<name>[a-zA-Z]+)", func(params martini.Params) string {
  return fmt.Sprintf ("Hello %s", params["name"])
})
~~~
Jetez un oeil à la documentation [Go documentation](http://golang.org/pkg/regexp/syntax/) pour plus d'informations sur la syntaxe des expressions régulières.

Un handlers de chemin peut s'ajouter l'un par dessus l'autre, ce qui est pratique pour les choses comme l'authentification et l'autorisation:
~~~ go
m.Get("/secret", authorize, func() {
  // this will execute as long as authorize doesn't write a response
})
~~~

Un groupe de chemins peut aussi être ajouté en utilisant la méthode de Group.
~~~ go
m.Group("/books", func(r martini.Router) {
    r.Get("/:id", GetBooks)
    r.Post("/new", NewBook)
    r.Put("/update/:id", UpdateBook)
    r.Delete("/delete/:id", DeleteBook)
})
~~~

Comme vous pouvez passer des middlewares à un handler, vous pouvez passer des middlewares à des groupes.
~~~ go
m.Group("/books", func(r martini.Router) {
    r.Get("/:id", GetBooks)
    r.Post("/new", NewBook)
    r.Put("/update/:id", UpdateBook)
    r.Delete("/delete/:id", DeleteBook)
}, MyMiddleware1, MyMiddleware2)
~~~

### Services
Services sont des objets qui sont disponibles pour être injecté dans la liste d'arguments d'un handler. Vous pouvez pointer un service selon différent niveau *Global* ou *Requête*.


#### Global Mapping
Une instance Martini implémente l'inject.Injector interface, alors pointer un service est facile:
~~~ go
db := &MyDatabase{}
m := martini.Classic()
m.Map(db) // the service will be available to all handlers as *MyDatabase
// ...
m.Run()
~~~

#### Requête-Level Mapping
Pour utiliser un pointeur au niveau d'une requête cela peut être fait avec un handler [martini.Context](http://godoc.org/github.com/go-martini/martini#Context):
~~~ go
func MyCustomLoggerHandler(c martini.Context, req *http.Request) {
  logger := &MyCustomLogger{req}
  c.Map(logger) // mapped as *MyCustomLogger
}
~~~

#### Mapping valeurs a des Interfaces
L'une des plus puissantes parties à propos des services c'est l'habilité à pointer à une interface. Par exemple, si vous voulez surcharger le [http.ResponseWriter](http://godoc.org/net/http#ResponseWriter) avec un objet qui est englober et réalise des opérations, vous pouvez écrire le handler suivant:
~~~ go
func WrapResponseWriter(res http.ResponseWriter, c martini.Context) {
  rw := NewSpecialResponseWriter(res)
  c.MapTo(rw, (*http.ResponseWriter)(nil)) // override ResponseWriter with our wrapper ResponseWriter
}
~~~

### Serving fichiers statiques
Instancie automatiquement des fichiers statiques du dossier "public" à la racine de votre serveur [martini.Classic()](http://godoc.org/github.com/go-martini/martini#Classic).
Vous pouvez instancier plus de dossiers en ajoutant plus de [martini.Static](http://godoc.org/github.com/go-martini/martini#Static) handlers.
~~~ go
m.Use(martini.Static("assets")) // serve from the "assets" directory as well
~~~

## Middleware Handlers
Middleware handlers sont placés entre la requête http entrante et le routeur. Ils ne sont aucunement différents des autres handlers présent dans Martini. Vous pouvez ajouter un middleware handler à la structure comme ceci:
~~~ go
m.Use(func() {
  // do some middleware stuff
})
~~~
Vous avez un contrôle total sur la structure middleware avec la fonction 'Handlers'. Cela remplacera n'importe quels handlers défini précédemment:
~~~ go
m.Handlers(
  Middleware1,
  Middleware2,
  Middleware3,
)
~~~
Middleware Handlers fonctionne très bien pour des automatisations comme le logging, l'autorisation, l'authentification, sessions, gzipping, pages d'erreur et n'importe quels opérations qui se font avant ou après chaque requête http:
~~~ go
// validate an api key
m.Use(func(res http.ResponseWriter, req *http.Request) {
  if req.Header.Get("X-API-KEY") != "secret123" {
    res.WriteHeader(http.StatusUnauthorized)
  }
})
~~~

### Next() (Suivant)
[Context.Next()](http://godoc.org/github.com/go-martini/martini#Context) est une fonction optionnelle que peut appeler le Middleware Handlers permettant d'attendre la réponse avant que les autres Handlers soient exécutés. Cela fonctionne très bien pour toutes opérations qui intervient après une requête http:
~~~ go
// log before and after a request
m.Use(func(c martini.Context, log *log.Logger){
  log.Println("before a request")

  c.Next()

  log.Println("after a request")
})
~~~

## Martini Env
Plusieurs Martini handlers utilise 'martini.Env' comme variable globale pour fournir des fonctionnalités spéciales entre l'environnement de développement et l'environnement de production. Il est recommandé que la variable 'MARTINI_ENV=production' soit défini pour déployer un serveur Martini en environnement de production.

## FAQ (Foire aux questions)

### Ou je peux trouver des middleware ?
Commencer par regarder dans le [martini-contrib](https://github.com/martini-contrib) projet. S'il n'est pas dedans n'hésitez pas à contacter un membre de l'équipe martini-contrib pour ajouter un nouveau dépôt à l'organisation.

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

### Comment j'intègre avec des serveurs existants ?
Une instance Martini implémente `http.Handler`, alors il peut être facilement utilisé pour alimenter des sous arbres sur des serveurs Go existants. Par exemple ceci est une application Martini pour Google App Engine:

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

### Comment je change le port/adresse?

La fonction `Run` de Martini regarde le PORT et l'adresse de variables d'environnement et les utilises. Sinon Martini utilisera par default le localhost:3000.
Pour avoir plus de flexibilité sur le port et l'adresse, utilisé la fonction `http.ListenAndServe` à la place.

~~~ go
  m := martini.Classic()
  // ...
  log.Fatal(http.ListenAndServe(":8080", m))
~~~

### Live code reload?

[gin](https://github.com/codegangsta/gin) and [fresh](https://github.com/pilu/fresh) tous les deux recharge en direct(live) le code des applications martini.

## Contribuer
Martini est destiné à rester petit et épuré. Toutes les contributions doivent finir dans un dépot dans l'origanisation [martini-contrib](https://github.com/martini-contrib). Si vous avez une contribution pour le noyau de Martini, n'hésitez pas à envoyer une Pull Request.

## A propos

Inspiré par [express](https://github.com/visionmedia/express) et [Sinatra](https://github.com/sinatra/sinatra), Martini est obsessionnellement désigné par personne autre que le [Code Gangsta](http://codegangsta.io/)

