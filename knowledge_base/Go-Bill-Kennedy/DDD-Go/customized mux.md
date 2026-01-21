#go #go-http-mux

### `main.go`
```go
	api := http.Server{
		Addr:         cfg.Web.APIHost,
		Handler:      mux.WebAPI(),
		ReadTimeout:  cfg.Web.ReadTimeout,
		WriteTimeout: cfg.Web.WriteTimeout,
		IdleTimeout:  cfg.Web.IdleTimeout,
		ErrorLog:     logger.NewStdLogger(log, logger.LevelError),
	}
```
### `mux.go`
```go
package mux

import (
	"net/http"

	"github.com/vishal2098govind/service/apis/services/sales/route/sys/checkapi"
)

func WebAPI() *http.ServeMux {
	mux := http.NewServeMux()

	checkapi.Routes(mux)

	return mux
}
```
### `checkapi.go`
```go
package checkapi

import (
	"encoding/json"
	"net/http"
)

func Routes(mux *http.ServeMux) {
	mux.HandleFunc("GET /liveness", liveness)
	mux.HandleFunc("GET /readiness", readiness)
}

func liveness(w http.ResponseWriter, r *http.Request) {
	resp := struct {
		Status string
	}{
		Status: "OK",
	}

	json.NewEncoder(w).Encode(resp)
}

func readiness(w http.ResponseWriter, r *http.Request) {
	resp := struct {
		Status string
	}{
		Status: "OK",
	}

	json.NewEncoder(w).Encode(resp)
}

```
- here, if we want to do a consistent error handling, we would need to do a lot of repetition but, we should try to be DRY
- thus, we can wrap the handlers 
![[wrapping mux|800]]

### `web.go`
```go
package web

import (
	"context"
	"fmt"
	"net/http"
	"os"
)

type Handler func(ctx context.Context, w http.ResponseWriter, r *http.Request) error

type App struct {
	*http.ServeMux
	shutdown chan os.Signal
}

func NewApp(shutdown chan os.Signal) *App {
	return &App{
		ServeMux: http.NewServeMux(),
		shutdown: shutdown,
	}
}

func (a *App) HandleFunc(pattern string, handler Handler) {

	h := func(w http.ResponseWriter, r *http.Request) {

		// CAN PUT SOME CODE HERE

		if err := handler(r.Context(), w, r); err != nil {
			// CAN HANDLER ERROR HERE
			fmt.Println(err)
		}

		// CAN PUT SOME CODE HERE

	}

	a.ServeMux.HandleFunc(pattern, h)
}
```

### `mux.go`
```go
package mux

import (
	"os"

	"github.com/vishal2098govind/service/apis/services/sales/route/sys/checkapi"
	"github.com/vishal2098govind/service/foundations/web"
)

func WebAPI(shutdown chan os.Signal) *web.App {
	mux := web.NewApp(shutdown)

	checkapi.Routes(mux)

	return mux
}

```
### `checkapi.go`
```go
package checkapi

import (
	"context"
	"encoding/json"
	"net/http"
	
	"github.com/vishal2098govind/service/foundations/web"
)

func Routes(app *web.App) {
	app.HandleFunc("GET /liveness", liveness)
	app.HandleFunc("GET /readiness", readiness)
}

func liveness(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	resp := struct {
		Status string
	}{
		Status: "OK",
	}

	return json.NewEncoder(w).Encode(resp)
}

func readiness(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	resp := struct {
		Status string
	}{
		Status: "OK",
	}

	return json.NewEncoder(w).Encode(resp)
}

```