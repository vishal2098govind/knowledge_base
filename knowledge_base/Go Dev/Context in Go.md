#go #context

## Context package in go

- `type Context` is just an object or interface with couple of different methods in it.
- It is mostly used for 
	- storing values that are very specific to each individual request 
	- or it can also be used to do things like set **timeouts** and the cancel processes.
	- Article: [How to use contexts in Go](https://www.digitalocean.com/community/tutorials/how-to-use-contexts-in-go)

When developing a large application, especially in server software, sometimes it’s helpful for a function to know more about the environment it’s being executed in aside from the information needed for a function to work on its own. For example, if a web server function is handling an HTTP request for a specific client, the function may only need to know which URL the client is requesting to serve the response. The function might only take that URL as a parameter. However, things can always happen when serving a response, such as the client disconnecting before receiving the response. If the function serving the response doesn’t know the client disconnected, the server software may end up spending more computing time than it needs calculating a response that will never be used.

In this case, being aware of the context of the request, such as the client’s connection status, allows the server to stop processing the request once the client disconnects. This saves valuable compute resources on a busy server and frees them up to handle another client’s request. This type of information can also be helpful in other contexts where functions take time to execute, such as making database calls. To enable ubiquitous access to this type of information, Go has included a `context` package in its standard library.

### `context.WithValue(parent Context, key, val any)` method
Returns a copy of parent in which the value associated with `key` is `val`

`context.Value(key any) any` is used to get the value of a `key`

Use context Values only for request-scoped data that transits processes and APIs, not for passing optional parameters to functions.

The provided key must be comparable and should not be of type string or any other built-in type to avoid collisions between packages using context. Users of WithValue should define their own types for keys. To avoid allocating when assigning to an `interface{}`, context keys often have concrete type struct{}. Alternatively, exported context key variables' static type should be a pointer or interface.

### `context.Background` and `context.TODO`:
`context` package has two ways of creating instance of context:
	- `context.Background()`
		- used when we know we're kind of at the **top level**, 
		- and we need to generate a context that's going to be **used throughout the application**
		- typically used by main function, initialisation and tests, top-level context for incoming requests
	- `context.TODO()`
		- designed to be a place-holder
		- supposed to clearly show that we are intending on using a different context 
		- or we are not sure where we're going to get a context from.
		- should be used when it's unclear which context to use or it is not yet available

```go
ctx := context.Background()
```

### Custom and Un exported type for context keys
- Using a custom type, specifically an un-exported custom type is a better idea when we're storing context values.

```go
package main

import (
	"context"
	"fmt"
)

type ctxKey string

const (
	favouriteColorKey ctxKey = "favourite-color"
)

func main() {
	ctx := context.Backgroun()
	ctx := context.WithValue(ctx, favouriteColorKey, "blue")
	value := ctx.Value("favourite-color)
	fmt.Println(value) // nil
	value := ctx.Value(favouriteColorKey)
	fmt.Println(value) // blue
}
```

- Why we want to do this? - i.e. Why custom and un exported type?
	- One of the things that could potentially happen, while looking up for a context value, is that other code whether it's outside of our package or somewhere else, might also be storing values inside of the context
	- And what we don't want to happen is to two different packages to override the same key causing some sort of issue

#### Converting `context.Value` to desired type from any
#type-assertions

```go
// using both values of type-assertion will not lead to any run-time error if type conversion fails
value, ok := ctx.Value(favouriteColorKey).(int)
if !ok {
	// handle
}

// not using the second value returned by the type-assertion will lead to a run-time error, resulting in panic
value := ctx.Value(favouriteColorKey).(string)
```


### Custom context package to wrap the std context package
- This helps us import our custom context package throughout our codebase, and not actually using the standard library context package
```go
package context

import (
	"context"

	userM "github.com/vishal2098govind/lenslocked/models/user"
)

type key string

const (
	userKey key = "user"
)

func WithUser(ctx context.Context, user *userM.User) context.Context {
	return context.WithValue(ctx, userKey, user)
}

func User(ctx context.Context) *userM.User {
	user, ok := ctx.Value(userKey).(*userM.User)
	if !ok {
		// this block of code is reached in two cases:
		// 1. if the key is not set ever, it will have a nil which is not of type of *userM.User (in go, even nil values have type)
		// 2. if the key is set to an invalid value, which is not of type *userM.User
		return nil
	}
	return user
}
```
In go, [When nil Isn't Equal to nil](https://www.calhoun.io/when-nil-isnt-equal-to-nil)

## Using context with http.Handler
#middleware 

```go
// main.go
package main

import (
	//...
	"github.com/vishal2098govind/lenslocked/controllers"
	
)

func main() {
	// ...
	usersC := controllers.Users{UserService: &userS, SessionService: &sessionS}
	
	r.Use(usersC.UserMW)
	r.Get("/users/me", usersC.CurrentUser)
	// ...
}
```

```go
// context/users.go
package context

import (
	"context"

	userM "github.com/vishal2098govind/lenslocked/models/user"
)

type key string

const (
	userKey key = "user"
)

func WithUser(ctx context.Context, user *userM.User) context.Context {
	return context.WithValue(ctx, userKey, user)
}

func User(ctx context.Context) *userM.User {
	user, ok := ctx.Value(userKey).(*userM.User)
	if !ok {
		return nil
	}
	return user
}

```

```go
// controllers/users.go
package controllers

import (
	"github.com/vishal2098govind/lenslocked/context"
)

type Users struct {
	UserService *userM.UserService

	SessionService *sessionM.SessionService

	Templates struct {
		New         Template // signup
		SignIn      Template // signin
		CurrentUser Template // current user
	}
}

// user middleware includes a context with currently signed in user based on the incoming request
func (u Users) UserMW(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		sessionC, err := readCookie(r, CookieSession)
		if err != nil {
			h.ServeHTTP(w, r)
			return
		}

		res, err := u.SessionService.User(sessionM.GetUserIdRequest{
			Token: sessionC,
		})
		if err != nil || res.User == nil {
			fmt.Println(err)
			h.ServeHTTP(w, r)
			return
		}

		ctx := r.Context()
		// it's important to note that context is coming from our custom context package
		ctx = context.WithUser(ctx, res.User)
		r = r.WithContext(ctx)
		h.ServeHTTP(w, r)
	})
}

func (u Users) CurrentUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	// it's important to note that context is coming from our custom context package
	user := context.User(ctx)
	if user == nil {
		http.Redirect(w, r, "/signin", http.StatusFound)
		return
	}

	u.Templates.CurrentUser.Execute(w, r, user)
}
```