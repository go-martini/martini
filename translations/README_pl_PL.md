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
* [Domyślna konfiguracja (Martini Classic)](#domyślna-konfiguracja-martini-classic))
  * [Handlery](#handlery)
  * [Routing](#routing)
  * [Usługi](#usługi)
  * [Serwowanie plików statycznych](#serwowanie-plików-statycznych)
* [Handlery middleware'ów](#handlery-middlewareów)
  * [Next()](#next)
* [Zmienne środowiskowe Martini](#zmienne-środowiskowe-martini)
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
  * [*log.Logger](http://godoc.org/log#Logger) - globalny logger dla Martini.
  * [martini.Context](http://godoc.org/github.com/go-martini/martini#Context) - kontekst żądania HTTP.
  * [martini.Params](http://godoc.org/github.com/go-martini/martini#Params) - `map[string]string` przechowująca nazwane parametry, znalezione podczas dopasowywania _routes_.
  * [martini.Routes](http://godoc.org/github.com/go-martini/martini#Routes) - usługa wspierająca _route'y_.
  * [martini.Route](http://godoc.org/github.com/go-martini/martini#Route) - bieżacy aktywny _route_.
  * [http.ResponseWriter](http://godoc.org/net/http/#ResponseWriter) - interfejs zapisu odpowiedzi HTTP.
  * [*http.Request](http://godoc.org/net/http/#Request) - żądanie HTTP.

### Routing
W Martini, jako _route_ należy rozumieć metodę HTTP skojarzoną ze wzorcem dopasowującym adres URL.
Każdy wzorzec może być skojarzony z jedną lub więcej metod handlera:
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

_Route'y_ są dopasowywane w kolejności ich definiowania. Pierwszy dopasowany _route_ zostanie wywołany. 

Wzorce ścieżek _route'ów_ mogą zawierać nazwane paremetry, dostępne poprzez usługę  [martini.Params](http://godoc.org/github.com/go-martini/martini#Params):
~~~ go
m.Get("/hello/:name", func(params martini.Params) string {
  return "Hello " + params["name"]
})
~~~

_Route'y_ mogą zostać dopasowane z wartościami globalnymi:
~~~ go
m.Get("/hello/**", func(params martini.Params) string {
  return "Hello " + params["_1"]
})
~~~

Również wyrażenia regularne mogą zostać użyte:
~~~go
m.Get("/hello/(?P<name>[a-zA-Z]+)", func(params martini.Params) string {
  return fmt.Sprintf ("Hello %s", params["name"])
})
~~~
Więcej informacji o budowie wyrażeń regularnych znajdziesz w [dokumentacji Go](http://golang.org/pkg/regexp/syntax/).

Handlery można organizować w stosy wywołań, co przydaje się przy mechanizmach takich jak uwierzytelnianie i autoryzacja:
~~~ go
m.Get("/secret", authorize, func() {
  // funkcja będzie wywoływana dopóty, dopóki authorize nie zwróci odpowiedzi
})
~~~

Grupy _route'ów_ mogą zostać dodane przy pomocy metody Group.
~~~ go
m.Group("/books", func(r martini.Router) {
    r.Get("/:id", GetBooks)
    r.Post("/new", NewBook)
    r.Put("/update/:id", UpdateBook)
    r.Delete("/delete/:id", DeleteBook)
})
~~~

W taki sam sposób jak przekazujesz middleware'y do handlerów, to możesz przekazywać middleware'y do grup.
~~~ go
m.Group("/books", func(r martini.Router) {
    r.Get("/:id", GetBooks)
    r.Post("/new", NewBook)
    r.Put("/update/:id", UpdateBook)
    r.Delete("/delete/:id", DeleteBook)
}, MyMiddleware1, MyMiddleware2)
~~~

### Usługi
Usługi są obiektami możliwymi do wstrzyknięcia poprzez listę argumentów danego handlera i mogą być mapowane na poziomie *globalnym* lub *żądania*.

#### Mapowanie globalne
Instancja Martini implementuje interfejs inject.Injector interface, więc mapowanie jest bardzo proste:
~~~ go
db := &MyDatabase{}
m := martini.Classic()
m.Map(db) // usługa będzie dostępna dla wszystkich handlerów jako *MyDatabase
// ...
m.Run()
~~~

#### Mapowanie na poziomie żądania
Mapowanie na poziomie żądania może być wykonane w handlerze poprzez [martini.Context](http://godoc.org/github.com/go-martini/martini#Context):
~~~ go
func MyCustomLoggerHandler(c martini.Context, req *http.Request) {
  logger := &MyCustomLogger{req}
  c.Map(logger) // zmapowany jako *MyCustomLogger
}
~~~

#### Mapowanie wartości na interfejsy
Jedną z mocnych stron usług jest możliwość zmapowania konkretnej usługi na interfejs. Dla przykładu, jeśli chcesz nadpisać [http.ResponseWriter](http://godoc.org/net/http#ResponseWriter) obiektem, który go opakowuje i wykonuje dodatkowe operacje, to możesz napisać następujący handler:
~~~ go
func WrapResponseWriter(res http.ResponseWriter, c martini.Context) {
  rw := NewSpecialResponseWriter(res)
  c.MapTo(rw, (*http.ResponseWriter)(nil)) // nadpisz oryginalny ResponseWriter naszym ResponseWriterem
}
~~~

### Serwowanie plików statycznych
Instancja [martini.Classic()](http://godoc.org/github.com/go-martini/martini#Classic) automatycznie serwuje statyczne pliki z katalogu "public" znajdującym się bezpośrednio w głównym katalogu serwera. Możliwe jest serwowanie dodatkowych katalogów poprzez dodanie handlerów [martini.Static](http://godoc.org/github.com/go-martini/martini#Static).
~~~ go
m.Use(martini.Static("assets")) // serwuj zasoby z katalogu "assets"
~~~

#### Serwowanie domyślnej strony
Możesz zdefiniować adres URL lokalnego pliku, który będzie serwowany gdy żądany adres URL nie zostanie znaleziony. Dodatkowo możesz zdefiniować prefiks wykluczający, który spowoduje, że niektórze adresy URL zostaną zignorowane. Jest to przydatna opcja dla serwerów, które jednocześnie serwują statyczne pliki i mają zdefiniowane handlery (np. REST API). Warto także rozważyć zdefiniowanie statycznych handlerów jako części łańcucha NotFound.

W poniższym przykładzie aplikacja serwuje plik `/index.html`, gdy tylko adres URL nie zostanie dopasowany do istniejącego lokalnego pliku i nie zaczyna się prefiksem `/api/v`:
~~~ go
static := martini.Static("assets", martini.StaticOptions{Fallback: "/index.html", Exclude: "/api/v"})
m.NotFound(static, http.NotFound)
~~~

## Handlery middleware'ów
Handlery middleware'ów są uruchamiane po otrzymaniu żądania HTTP a przed przekazaniem go do routera. W zasadzie nie ma różnicy między nimi a handlerami Martini. Handler middleware'a można dodać do stosu wywołań w następujący sposób:
~~~ go
m.Use(func() {
  // wykonaj operacje zdefiniowane przez middleware
})
~~~

Pełną kontrolę na stosem middleware'owym zapewnia funkcja `Handlers`. Poniższy przykład prezentuje, jak można zamienić poprzednio skonfigurowane handlery:
~~~ go
m.Handlers(
  Middleware1,
  Middleware2,
  Middleware3,
)
~~~

Handlery middleware'ów sprawdzają się doskonale dla mechanizmów takich jak logowanie, autoryzacja, uwierzytelnianie, obsługa sesji, kompresja odpowiedzi, strony błędów i innych, których operacje muszą zostać wykonane przed i po obsłudze żądania HTTP:
~~~ go
// validate an api key
m.Use(func(res http.ResponseWriter, req *http.Request) {
  if req.Header.Get("X-API-KEY") != "secret123" {
    res.WriteHeader(http.StatusUnauthorized)
  }
})
~~~

### Next()
[Context.Next()](http://godoc.org/github.com/go-martini/martini#Context) jest opcjonalną funkcją, którą handlery middleware'ów wywołują, żeby przekazać tymczasowo obsługę żadania do kolejnych handlerów, a później do niej wrócić. Mechanizm sprawdza się doskonale w przypadku wykonywania operacji po obsłudze żądania HTTP:
~~~ go
// zaloguj przed i po żądaniu
m.Use(func(c martini.Context, log *log.Logger){
  log.Println("before a request")

  c.Next()

  log.Println("after a request")
})
~~~

## Zmienne środowiskowe Martini

Niektóre handlery Martini wykorzystują globalną zmienną `martini.Env` by dostarczać specjalne funkcje dla środowisk deweloperskich i produkcyjnych. Zaleca się ustawienie zmiennej `MARTINI_ENV=production` w środowisku produkcyjnym.

## FAQ

### Gdzie mam szukać middleware'u X?

Proponujemy zacząć poszukiwania od projektów należących do [martini-contrib](https://github.com/martini-contrib). Jeśli dany middleware się tam nie znajduje, skontaktuj się z członkiem zespołu martini-contrib i poproś go o dodanie nowego repozytorium do organizacji.

* [acceptlang](https://github.com/martini-contrib/acceptlang) - Handler umożliwiający parsowanie nagłówka HTTP `Accept-Language`.
* [accessflags](https://github.com/martini-contrib/accessflags) - Handler dołączający obsługę kontroli dostępu.
* [auth](https://github.com/martini-contrib/auth) - Handlery uwierzytelniające.
* [binding](https://github.com/martini-contrib/binding) - Handler mapujący/walidujący żądanie na strukturę.
* [cors](https://github.com/martini-contrib/cors) - Handler dostarcza wsparcie dla CORS.
* [csrf](https://github.com/martini-contrib/csrf) - Ochrona CSRF dla aplikacji.
* [encoder](https://github.com/martini-contrib/encoder) - Usługa enkodująca treść odpowiedzi w różnych formatach, wspiera negocjacje formatu.
* [gzip](https://github.com/martini-contrib/gzip) - Handler dla kompresji GZIP żądań.
* [gorelic](https://github.com/martini-contrib/gorelic) - NewRelic middleware.
* [logstasher](https://github.com/martini-contrib/logstasher) - Middleware zwracający odpowiedź formacie kompatybilnym z logstash JSONem.
* [method](https://github.com/martini-contrib/method) - Nadpisywanie metod HTTP poprzez nagłówek.
* [oauth2](https://github.com/martini-contrib/oauth2) - Handler dostarczający logowanie OAuth 2.0 dla aplikacji Martini. Logowanie Google Sign-in, Facebook Connect i Github wspierane.
* [permissions2](https://github.com/xyproto/permissions2) - Handler śledzący użytkowników, ich logowania i uprawnienia.
* [render](https://github.com/martini-contrib/render) - Handler dostarczający usługę łatwo renderującą odpowiedź do formatu JSON i szablonów HTML.
* [secure](https://github.com/martini-contrib/secure) - Implementuje kilka szybkich "quick-wins" związanych z bezpieczeństwem.
* [sessions](https://github.com/martini-contrib/sessions) - Handler dostarcza usługę sesji.
* [sessionauth](https://github.com/martini-contrib/sessionauth) - Handler, który umożliwia w prosty sposób nałożenie reguły wymagania logowania dla konkretnych adresów oraz obsługę zalogowanych użytkowników w sesji. 
* [strict](https://github.com/martini-contrib/strict) - Strict Mode 
* [strip](https://github.com/martini-contrib/strip) - Pomijanie prefiksu URL.
* [staticbin](https://github.com/martini-contrib/staticbin) - Handler umożliwia serwowanie statycznych plików z zasobów binarnych.
* [throttle](https://github.com/martini-contrib/throttle) - Middleware kontrolujący przepustowość handlerów.
* [vauth](https://github.com/rafecolton/vauth) - Handlery wspierające vendorowe uwierzytelnianie (obecnie GitHub i TravisCI).
* [web](https://github.com/martini-contrib/web) - Kontekst znany z web.go.

### Jak mogę zintegrować Martini z istniejącymi serwerami?

Instacja Martini implementuje `http.Handler`, więc może być łatwo wykorzystana do serwowania całych drzew zasobów na istniejących serwerach Go. Przykład przedstawia działającą aplikację Martini dla Google App Engine:

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

### Jak mogę zmienić host/port?

Funkcja `Run` sprawdza, czy są zdefiniowane zmienne środowiskowe HOST i PORT, i jeśli są to ich używa. W przeciwnym wypadku Martini uruchomi się z domyślnymi ustawieniami localhost:3000.
W celu uzyskania większej kontroli nad hostem i portem, skorzystaj z funkcji `martini.RunOnAddr`.

~~~ go
  m := martini.Classic()
  // ...
  log.Fatal(m.RunOnAddr(":8080"))
~~~

### Automatyczne przeładowywanie kodu aplikacji (Live code reload)

[gin](https://github.com/codegangsta/gin) i [fresh](https://github.com/pilu/fresh) wspierają przeładowywanie kodu aplikacji.

## Rozwijanie
Martini w założeniu ma pozostać czysty i uporządkowany. Większość kontrybucji powinna trafić jako repozytorium organizacji [martini-contrib](https://github.com/martini-contrib). Jeśli masz kontrybucję do core'a projektu Martini, zgłoś Pull Requesta.

## O projekcie

Inspirowany [expressem](https://github.com/visionmedia/express) i [sinatrą](https://github.com/sinatra/sinatra)

Martini został obsesyjnie zaprojektowany przez nikogo innego jak przez [Code Gangsta](http://codegangsta.io/)
