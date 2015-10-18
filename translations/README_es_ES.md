# Martini  [![wercker status](https://app.wercker.com/status/9b7dbc6e2654b604cd694d191c3d5487/s/master "wercker status")](https://app.wercker.com/project/bykey/9b7dbc6e2654b604cd694d191c3d5487)[![GoDoc](https://godoc.org/github.com/go-martini/martini?status.png)](http://godoc.org/github.com/go-martini/martini)

Martini es un poderoso paquete para escribir rápidamente aplicaciones/servicios web modulares en Golang.


## Vamos a iniciar

Después de instalar Go y de configurar su [GOPATH](http://golang.org/doc/code.html#GOPATH), cree su primer archivo `.go`. Vamos a llamar a este `server.go`.

~~~ go
package main

import "github.com/go-martini/martini"

func main() {
  m := martini.Classic()
  m.Get("/", func() string {
    return "Hola Mundo!"
  })
  m.Run()
}
~~~

Luego instale el paquete Martini (Es necesario **go 1.1** o superior):
~~~
go get github.com/go-martini/martini
~~~

Después corra su servidor:
~~~
go run server.go
~~~

Ahora tendrá un webserver Martini corriendo en el puerto `localhost:3000`.

## Obtenga ayuda

Suscribase a la [Lista de email](https://groups.google.com/forum/#!forum/martini-go)

Observe el [Video demostrativo](http://martini.codegangsta.io/#demo)

Use la etiqueta [martini](http://stackoverflow.com/questions/tagged/martini) para preguntas en Stackoverflow

Documentación [GoDoc](http://godoc.org/github.com/go-martini/martini)


## Caracteríticas
* Extremadamente simple de usar.
* Diseño no intrusivo.
* Buena integración con otros paquetes Golang.
* Enrutamiento impresionante.
* Diseño modular - Fácil de añadir y remover funcionalidades.
* Muy buen uso de handlers/middlewares.
* Grandes características innovadoras.
* **Compatibilidad total con la interface [http.HandlerFunc](http://godoc.org/net/http#HandlerFunc).**
* Sirviendo documentos por defecto (e.g. para servir aplicaciones AngularJS en modo HTML5).

## Más Middlewares
Para más middlewares y funcionalidades, revisar los repositorios en [martini-contrib](https://github.com/martini-contrib).

## Lista de contenidos
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
Para iniciar rápidamente, [martini.Classic()](http://godoc.org/github.com/go-martini/martini#Classic) prevee algunas herramientas que funcionan bien para la mayoría de aplicaciones web:
~~~ go
  m := martini.Classic()
  // middlewares y rutas aquí
  m.Run()
~~~

Algunas funcionalidades que [martini.Classic()](http://godoc.org/github.com/go-martini/martini#Classic) ofrece automáticamente son:
  * Request/Response Logging - [martini.Logger](http://godoc.org/github.com/go-martini/martini#Logger)
  * Panic Recovery - [martini.Recovery](http://godoc.org/github.com/go-martini/martini#Recovery)
  * Static File serving - [martini.Static](http://godoc.org/github.com/go-martini/martini#Static)
  * Routing - [martini.Router](http://godoc.org/github.com/go-martini/martini#Router)

### Handlers
Handlers son el corazón y el alma de Martini. Un handler es básicamente cualquier tipo de función que puede ser llamada.
~~~ go
m.Get("/", func() {
  println("hola mundo")
})
~~~

#### Retorno de Valores
Si un handler retorna cualquier cosa, Martini escribirá el valor retornado como una cadena [http.ResponseWriter](http://godoc.org/net/http#ResponseWriter):
~~~ go
m.Get("/", func() string {
  return "hola mundo" // HTTP 200 : "hola mundo"
})
~~~

Usted también puede retornar un código de estado:
~~~ go
m.Get("/", func() (int, string) {
  return 418, "soy una tetera" // HTTP 418 : "soy una tetera"
})
~~~

#### Inyección de Servicios
Handlers son invocados vía reflexión. Martini utiliza *Inyección de Dependencia* para resolver dependencias en la lista de argumentos Handlers. **Esto hace que Martini sea completamente compatible con  la interface `http.HandlerFunc` de golang.**

Si agrega un argumento a su Handler, Martini buscará en su lista de servicios e intentará resolver su dependencia vía su tipo de aserción:
~~~ go
m.Get("/", func(res http.ResponseWriter, req *http.Request) { // res y req son inyectados por Martini
  res.WriteHeader(200) // HTTP 200
})
~~~

Los siguientes servicios son incluidos con [martini.Classic()](http://godoc.org/github.com/go-martini/martini#Classic):
  * [*log.Logger](http://godoc.org/log#Logger) - Log Global para Martini.
  * [martini.Context](http://godoc.org/github.com/go-martini/martini#Context) - http request context.
  * [martini.Params](http://godoc.org/github.com/go-martini/martini#Params) - `map[string]string` de nombres de los parámetros buscados por la ruta.
  * [martini.Routes](http://godoc.org/github.com/go-martini/martini#Routes) - Servicio de ayuda para las Rutas.
  * [http.ResponseWriter](http://godoc.org/net/http/#ResponseWriter) - http Response escribe la interfaz.
  * [*http.Request](http://godoc.org/net/http/#Request) - http Request.

### Rutas
En Martini, una ruta es un método HTTP emparejado con un patrón URL. Cada ruta puede tener uno o más métodos handler:
~~~ go
m.Get("/", func() {
  // mostrar algo
})

m.Patch("/", func() {
  // actualizar algo
})

m.Post("/", func() {
  // crear algo
})

m.Put("/", func() {
  // reemplazar algo
})

m.Delete("/", func() {
  // destruir algo
})

m.Options("/", func() {
  // opciones HTTP
})

m.NotFound(func() {
  // manipula 404
})
~~~

Las rutas son emparejadas en el orden en que son definidas. La primera ruta que coincide con la solicitud es invocada.

Los patrones de rutas puede incluir nombres como parámetros accesibles vía el servicio [martini.Params](http://godoc.org/github.com/go-martini/martini#Params):
~~~ go
m.Get("/hello/:name", func(params martini.Params) string {
  return "Hello " + params["name"]
})
~~~

Las rutas se pueden combinar con globs:
~~~ go
m.Get("/hello/**", func(params martini.Params) string {
  return "Hello " + params["_1"]
})
~~~

Las expresiones regulares pueden ser usadas también:
~~~go
m.Get("/hello/(?P<name>[a-zA-Z]+)", func(params martini.Params) string {
  return fmt.Sprintf ("Hello %s", params["name"])
})
~~~
Observe la [documentación](http://golang.org/pkg/regexp/syntax/) para mayor información sobre la sintaxis de expresiones regulares.


Handlers de ruta pueden ser apilados encima de otros, lo cual es útil para cosas como autenticación y autorización:
~~~ go
m.Get("/secret", authorize, func() {
  // será ejecutado cuando autorice escribir una respuesta
})
~~~

Grupos de rutas puede ser añadidas usando el método Group.
~~~ go
m.Group("/books", func(r martini.Router) {
    r.Get("/:id", GetBooks)
    r.Post("/new", NewBook)
    r.Put("/update/:id", UpdateBook)
    r.Delete("/delete/:id", DeleteBook)
})
~~~

Al igual que usted puede pasar middlewares a un handler, puede pasar middlewares a grupos.
~~~ go
m.Group("/books", func(r martini.Router) {
    r.Get("/:id", GetBooks)
    r.Post("/new", NewBook)
    r.Put("/update/:id", UpdateBook)
    r.Delete("/delete/:id", DeleteBook)
}, MyMiddleware1, MyMiddleware2)
~~~

### Servicios
Lo servicios son objetos que están disponibles para ser inyectados en una lista de argumentos Handler. Usted puede mapear un servicio a nivel *Global* o *Request*.

#### Mapeo Global
Una instancia de Martini implementa la interface `inject.Injector`, asi que el mapeo de un servicio es sencillo:
~~~ go
db := &MyDatabase{}
m := martini.Classic()
m.Map(db) // el servicio estará disponible para todos los handlers como *MyDatabase.
// ...
m.Run()
~~~

#### Mapeo por Request
El mapeo a nivel de request puede ser hecho en un handler via  [martini.Context](http://godoc.org/github.com/go-martini/martini#Context):
~~~ go
func MyCustomLoggerHandler(c martini.Context, req *http.Request) {
  logger := &MyCustomLogger{req}
  c.Map(logger) // mapeado como *MyCustomLogger
}
~~~

#### Valores de Mapeo para Interfaces
Una de las partes más poderosas sobre servicios es la capacidad de mapear un servicio para una interface. Por ejemplo, si desea sobreescribir [http.ResponseWriter](http://godoc.org/net/http#ResponseWriter) con un objeto que envuelva y realice operaciones extra, puede escribir el siguiente handler:
~~~ go
func WrapResponseWriter(res http.ResponseWriter, c martini.Context) {
  rw := NewSpecialResponseWriter(res)
  c.MapTo(rw, (*http.ResponseWriter)(nil)) // sobreescribir ResponseWriter con nuestro ResponseWriter
}
~~~

### Sirviendo Archivos Estáticos
Una instancia de [martini.Classic()](http://godoc.org/github.com/go-martini/martini#Classic) sirve automáticamente archivos estáticos del directorio "public" en la raíz de su servidor.
Usted puede servir más directorios, añadiendo más [martini.Static](http://godoc.org/github.com/go-martini/martini#Static) handlers.
~~~ go
m.Use(martini.Static("assets")) // sirviendo los archivos del directorio "assets"
~~~

## Middleware Handlers
Los Middleware Handlers se sitúan entre una solicitud HTTP y un router. En esencia, ellos no son diferentes de cualquier otro Handler en Martini. Usted puede añadir un handler de middleware para la pila de la siguiente forma:
~~~ go
m.Use(func() {
  // Hacer algo con middleware
})
~~~

Puede tener el control total sobre la pila del Middleware con la función `Handlers`. Esto reemplazará a cualquier handler que haya sido establecido previamente:
~~~ go
m.Handlers(
  Middleware1,
  Middleware2,
  Middleware3,
)
~~~

Middleware Handlers trabaja realmente bien como logging, autorización, autenticación, sesión, gzipping, páginas de errores y una serie de otras operaciones que deben suceder antes o después de una solicitud http:
~~~ go
// Valida una llave de api
m.Use(func(res http.ResponseWriter, req *http.Request) {
  if req.Header.Get("X-API-KEY") != "secret123" {
    res.WriteHeader(http.StatusUnauthorized)
  }
})
~~~

### Next()
[Context.Next()](http://godoc.org/github.com/go-martini/martini#Context) es una función opcional que Middleware Handlers puede llamar para entregar hasta después de que una solicitud http haya sido ejecutada. Esto trabaja muy bien para calquier operación que debe suceder luego de una solicitud http:
~~~ go
// log antes y después de una solicitud
m.Use(func(c martini.Context, log *log.Logger){
  log.Println("antes de una solicitud")

  c.Next()

  log.Println("luego de una solicitud")
})
~~~

## Martini Env

Martini handlers hace uso de `martini.Env`, una variable global para proveer funcionalidad especial en ambientes de desarrollo y ambientes de producción. Es recomendado que una variable `MARTINI_ENV=production` sea definida cuando se despliegue en un ambiente de producción.

## FAQ

### ¿Dónde puedo encontrar una middleware X?

Inicie su búsqueda en los proyectos [martini-contrib](https://github.com/martini-contrib). Si no esta allí, no dude en contactar a algún miembro del equipo martini-contrib para adicionar un nuevo repositorio para la organización.

* [acceptlang](https://github.com/martini-contrib/acceptlang) - Handler para analizar el `Accept-Language` HTTP header.
* [accessflags](https://github.com/martini-contrib/accessflags) - Handler para habilitar Access Control.
* [auth](https://github.com/martini-contrib/auth) - Handlers para autenticación.
* [binding](https://github.com/martini-contrib/binding) - Handler para mapeo/validación de un request en una estructura.
* [cors](https://github.com/martini-contrib/cors) - Handler habilita soporte para CORS.
* [csrf](https://github.com/martini-contrib/csrf) - CSRF protección para aplicaciones.
* [encoder](https://github.com/martini-contrib/encoder) - Servicio de codificador para renderizado de datos en varios formatos y gestión de contenidos.
* [gzip](https://github.com/martini-contrib/gzip) - Handler para agregar compresión gzip para los requests.
* [gorelic](https://github.com/martini-contrib/gorelic) - NewRelic middleware
* [logstasher](https://github.com/martini-contrib/logstasher) - Middleware que imprime logstash-compatiable JSON.
* [method](https://github.com/martini-contrib/method) - HTTP method overriding via cabecera o campos de formulario.
* [oauth2](https://github.com/martini-contrib/oauth2) - Handler que provee OAuth 2.0 login para Martini apps. Google Sign-in, Facebook Connect y Github login son soportados.
* [permissions2](https://github.com/xyproto/permissions2) - Handler para hacer seguimiento de usuarios, estados de login y permisos.
* [render](https://github.com/martini-contrib/render) - Handler que provee un servicio para fácil renderizado de plantillas HTML y JSON.
* [secure](https://github.com/martini-contrib/secure) - Implementa un par de rápidos triunfos de seguridad.
* [sessions](https://github.com/martini-contrib/sessions) - Handler que provee un servicio de sesiones.
* [sessionauth](https://github.com/martini-contrib/sessionauth) - Handler que provee una simple forma de hacer routes, requerir un login y manejar login de usuarios en la sesión.
* [strict](https://github.com/martini-contrib/strict) - Modo estricto (Strict Mode)
* [strip](https://github.com/martini-contrib/strip) - Prefijo de extracción de URL.
* [staticbin](https://github.com/martini-contrib/staticbin) - Handler para servir archivos estáticos desde datos binarios.
* [throttle](https://github.com/martini-contrib/throttle) - Request rate throttling middleware.
* [vauth](https://github.com/rafecolton/vauth) - Handlers para expender autenticación con webhooks (actualmente GitHub y TravisCI)
* [web](https://github.com/martini-contrib/web) - Contexto hoisie web.go's

### ¿Cómo se integra con los servidores existentes?

Una instancia de Martini implementa `http.Handler`, de modo que puede ser fácilmente utilizado para servir sub-rutas y directorios en servidores Go existentes. Por ejemplo, este es un aplicativo Martini trabajando para Google App Engine:

~~~ go
package hello

import (
  "net/http"
  "github.com/go-martini/martini"
)

func init() {
  m := martini.Classic()
  m.Get("/", func() string {
    return "Hola Mundo!"
  })
  http.Handle("/", m)
}
~~~

### ¿Cómo cambiar el puerto/host?

La función `Run` de Martini observa las variables de entorno `PORT` y `HOST` para utilizarlas. De lo contrário, Martini asume por defecto `localhost:3000`. Para tener mayor flexibilidad sobre el puerto y host, use la función `martini.RunOnAddr`.

~~~ go
  m := martini.Classic()
  // ...
  log.Fatal(m.RunOnAddr(":8080"))
~~~

### ¿Servidor con autoreload?

[gin](https://github.com/codegangsta/gin) y [fresh](https://github.com/pilu/fresh) son ambas aplicaciones para autorecarga de Martini.

## Contribuyendo
Martini se desea mantener pequeño y limpio. La mayoría de contribuciones deben realizarse en el repositorio [martini-contrib](https://github.com/martini-contrib). Si desea hacer una contribución al core de Martini es libre de realizar un Pull Request.

## Sobre

Inspirado por [Express](https://github.com/visionmedia/express) y [Sinatra](https://github.com/sinatra/sinatra)

Martini está diseñado obsesivamente por nada menos que [Code Gangsta](http://codegangsta.io/)

