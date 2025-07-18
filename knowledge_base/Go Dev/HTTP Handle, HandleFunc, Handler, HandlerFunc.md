#go #http-handler-func #net/http #io-Writer #http-ResponseWriter #fmt-Fprint #http-ServeMux #http-handle-func #http-Handler

In go, we handle any incoming web request by reflecting to the reality that there are two things **request** and **response**
```
func handlerFunc(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "<h1>Welcome</h1>")
}
```

- **`http.ResponseWriter`** is an `interface`, where as `http.Request` is a `struct`
	- **`ResponseWriter`** needs to be or could be substituted by multiple implementations
	- thus has diff advantages:
		- different types of connections might allow for different type of responses
			- some responses might be streamable, where as some might be something where we need to write the entire response before sending all at once
		- helps to test easily
			- **`httptest.ResponseRecorder`** implements `http.ResponseWriter`, allowing us to easily test our code

- functions that match the signature where they take in `http.ResponseWriter` and `http.Request` and used quite often because they are default way of handling any sort of incoming web request in go
- thus, in `http` package there's also a type defined for this function signature called [`type http.HandlerFunc`](https://pkg.go.dev/net/http#HandlerFunc)
```
type HandlerFunc func(ResponseWriter, *Request)
```

#### Why are we able to pass an instance of ResponseWriter into `Fprintf`
- `Fprint` takes in `io.Writer` as it's first argument and anything to be written on to this writer
- `io.Writer`
```
type Writer interface {
	Write(p []byte) (n int, err error)
}
```
- `http.ResponseWriter` interface also has `Write` method, and thus qualifies to be of a type of `io.Writer`, thus `http.ResponseWriter` implements `io.Writer`
```
type ResponseWriter interface {
	Header() Header
	Write([]byte) (int, error)
	WriteHeader(statusCode int)
}
```

#### `HandleFunc` Vs `HandlerFunc`
- `http.HandleFunc` is used to register incoming request route handlers by accepting functions of type`http.HandlerFunc`
```
package http

...

func HandleFunc(pattern string, handler func(ResponseWriter, *Request)) {
 ...
}

...
```

#### `http.HandleFunc`, `ServeMux.HandleFunc`, `http.HandlerFunc`, `http.Handler`
- `http` defines a default `ServeMux` called `DefaultServeMux` defined at package level, which is used by **package level** `http.HandleFunc`
- `ServeMux` defines it's own method `HandleFunc`.
- The body of both `HandleFunc` is same, other than just the instance of `ServeMux` being used
	- `http.HandleFunc` uses `DefaultServeMux`
	- `func (*ServeMux) HandleFunc` uses `ServeMux`
	
- at the end of the day, we need to invoke `http.ListenAndServe` function which takes in `http.Handler` interface
```

package http
...
type Handler interface {
	ServeHTTP(ResponseWriter, *Request) // same as of type http.HandlerFunc
}
...

func ListenAndServe(addr string, handler Handler) error {
 ...
}
```
- The `http.ServeMux` implements `http.Handler` interface
```
package http

...

func (mux *ServeMux) ServeHTTP(w ResponseWriter, r *Request) {
	...
}
```
- Thus, if we are using `http.DefaultServeMux` which has package level state, we can use
```
package main

...
	http.HandleFunc("/", handlerFunc)
	http.ListenAndServe(":3000", nil) // passing nil to use http.DefaultServeMux
...

```


#### All confusing this together:
- **`HandlerFunc`**: 
	- `type http.HandlerFunc func(ResponseWriter, *Request)`
	- it also implements `http.Handler`
- **`HandleFunc`** (**`DefaultServeMux`**) **& `ServeMux.HandleFunc`**
	- `http.DefaultServeMux` uses package level state which is used by package level `http.HandleFunc` `func`
	- `ServeMux` is a struct type
	- `HandleFunc` is used to **register** `HandlerFunc` for a particular request route
- **`Handler.ServeHTTP` & `ServeMux.ServeHTTP`**
	- `ServeMux` implements `Handler` interface
	- `http.ListenAndServe(addr string, handler Handler)` expects implementation of `Handler`. If nil, uses `DefaultServeMux`
	- `http.ListenAndServe` actually keeps the go program away from getting exited
```
type Handler interface {
	ServeHTTP(ResponseWriter, *Request) // same as of type http.HandlerFunc
}
```

#### Using `DefaultServeMux`
```go
package main

import (
	"fmt"
	"net/http"
)

func handlerFunc(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "<h1>Welcome</h1>")
}

func main() {
	http.HandleFunc("/", handlerFunc)
	fmt.Println("Starting the server on :3000")
	http.ListenAndServe(":3000", nil)
}
```
#### Using `ServeMux`
```go
package main

import (
	"fmt"
	"net/http"
)

func handlerFunc(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "<h1>Welcome</h1>")
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", handlerFunc)
	fmt.Println("Starting the server on :3000")
	http.ListenAndServe(":3000", mux)
}
```
#### Using `http.HandlerFunc` type
- The `http.HandlerFunc` type inside the `http` package also has a method `ServeHTTP`, and thus `http.HandlerFunc` also implements `http.Handler` interface
```go
package http

// ...

type HandlerFunc func(ResponseWriter, *Request)

// ServeHTTP calls f(w, r).
func (f HandlerFunc) ServeHTTP(w ResponseWriter, r *Request) {
	f(w, r)
}
```
- Thus we can **convert** any function with parameters `(w http.ResponseWriter, r *http.Request)`  into `http.HandlerFunc` and use it to pass to `http.ListenAndServe` function
```go
package main

import (
	"fmt"
	"net/http"
)

func handlerFunc(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "<h1>Welcome</h1>")
}

func main() {
	fmt.Println("Starting the server on :3000")
	http.ListenAndServe(":3000", http.HandlerFunc(handlerFunc))
}
```
- This is useful in cases where we want to include #middleware in our router
```go

func main() {
	r := chi.NewRouter()
	r.Use(IPLoggerMiddleware)
	
	r.Get("/", usersC.Home)
	// ...
	
	http.ListenAndServe(":3000", r)
}

func IPLoggerMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ip := r.RemoteAddr
		fmt.Printf("IP: %v", ip)
		h.ServeHTTP(w, r)
	})
}

```
#### Using Custom router
- Implements `http.Handler`
```go
package main

import (
	"fmt"
	"net/http"
)

func homeHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, `<h1>Welcome</h1><a href="/contacts">Contacts</a>`)
}

func contactHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println(r.Method)
	fmt.Fprint(w, `<h1>contacts</h1><p>To get in touch, email me at <a href="mailto:vishal.govind2098@gmail.com">vishal.govind2098@gmail.com</a></p>`)
}

type Router struct{} // Implements http.Handler

func (router Router) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.URL.Path {
	case "/":
		homeHandler(w, r)
	case "/contacts":
		contactHandler(w, r)
	default:
		http.NotFound(w, r)	
	}
}

func main() {
	router := Router{}
	fmt.Println("Starting the server on :3000")
	http.ListenAndServe(":3000", router)
}
```
#### Why would we want to have such a setup of `Router` type rather than plain `HandlerFunc`s?
 It's very common to need different information
```go
package main

import (
	"database/sql"
	"fmt"
	"net/http"
)

type Server struct {
	DB *sql.DB
}

func (s *Server) homeHandler(w http.ResponseWriter, r *http.Request) {
	// we can make use of s.DB here
	fmt.Fprint(w, `<h1>Welcome</h1><a href="/contacts">Contacts</a>`)
}

func (s *Server) contactHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, `<h1>contacts</h1><p>To get in touch, email me at <a href="mailto:vishal.govind2098@gmail.com">vishal.govind2098@gmail.com</a></p>`)
}

type Router struct{
	Server *Server
} // Implements http.Handler

func (router Router) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.URL.Path {
	case "/":
		router.s.homeHandler(w, r)
	case "/contacts":
		router.s.contactHandler(w, r)
	default:
		http.NotFound(w, r)	
	}
}

func main() {
	s := Server{ DB: db }
	router := Router{ Server: &s }
	fmt.Println("Starting the server on :3000")
	http.ListenAndServe(":3000", router)
}
```


#### Why are there so many ways of doing this same thing in Go?
- one can simply convert a function to `Handler` type
- one can simply choose to use `http.DefaultServeMux` and simply register handlers using `http.HandleFunc`
- Behind the scenes, every approach is at the end converted to `http.Handler`  type and further use `ServeHTTP` method inside of it