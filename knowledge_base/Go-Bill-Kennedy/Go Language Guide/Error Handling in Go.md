#go #error-handling 

- Error handling is about 
	- **showing the user** of our API **enough respect** 
	- that we would give them an informed peace of information about the state of error
	- it's about context. if the user of our api have enough context they can take an informed decision


### the built-in `error` interface
```go
type error interface {
	Error() string
}
```
- errors in go are just values, and it can be anything we want them to be
- the error value is decoupled by the `error` interface
### the un-exported `errorString` struct from `"errors"` package
```go
package errors
 
type errorString struct {
	s string
}

func (e *errorString) Error() {
	return e.s
}

func New(text string) error {
	return &errorString{text}
}
```

### `if err nill` statement
- scope of `err` remains within the `if` block and every `if` statement has it's own scope, we can re-use the `err` variable name
```go
package main

import (
	"errors"
	"fmt"
)

func main() {
	if err := webCall(); err != nil {
		fmt.Printf("error: %f", err)
		return
	}
	
	if err := webCall(); err != nil {
		fmt.Printf("error: %f", err)
		return
	}
}

func webCall() error {
	return errors.New("Bad Request")
}
```

## Handling errors
- if a function could return more than one **custom error variable**:
```go
package main

import (
	"errors"
	"fmt"
)

var (
	// ErrBadRequest is returned when there are problems with request
	ErrBadRequest = errors.New("Bad Request")
	
	// ErrPageMoved is returned when a 301/302 is returned
	ErrPageMoved = errors.New("Page Moved")
)

func main() {
	if err := webCall(true); err != nil {
		switch err {
			case ErrBadRequest:
				fmt.Println("A Bad request")
				return
			case ErrPageMoved:
				fmt.Println("The Page moved")
				return
			default:
				fmt.Println(err)
				return
		}
	}
	
	fmt.Println("Life is Good")
}

func webCall(b bool) error {
	if b {
		return ErrBadRequest
	}
	
	return ErrPageMoved
}
```

## Custom Error Types as context
- Can use custom error types when the default error type is not enough to carry enough context
- Idiomatic Go:
	- **custom** error **variables** **start** with `Err` i.e. `ErrBadRequest`
	- **custom** error **types** **end** with `Error` i.e. `UnmarshalTypeError`
- E.g. `UnmarshalTypeError` from "json" package

```go
package json

type UnmarshalTypeError struct {
	Value string
	Type reflect.Type
}

func (e *UnmarshalTypeError) Error() string {
	return "json: cannot unmarshal " + e.Vallue + " into Go type " + e.Type
}

type InvalidUnmarshalError struct {
	Type reflect.Type
}

func (e *InvalidUnmarshalError) Error() string {
	return "json: cannot unmarshal Go type " + e.Type
}

func Unmarshal(b []byte, v interface{}) error {
	rv := reflect.ValueOf(v)
	if rv.Kind() != reflect.Ptr || rv.IsNill() {
		return &InvalidUnmarshalError{reflect.TypeOf(v)}
	}
	return &UnmarshalTypeError{"string", reflect.TypeOf(v)}
}
```
- handling function returning error that can be of more than one **custom error type**
	- Can leverage the Go's **type as context** to **switch** among the possible custom error type
```go
package main

type user struct {
	Name string
}

func main() {
	var u user
	err := Unmarshal([]byte{"name": "vishal"}, u)
	if err != nil {
		switch e := e.(type) {
			case *UnmarshalTypeError:
				// handle
			case *InvalidUnmarshalError:
				// handle
			default:
				fmt.Println(err)
		}
		return
	}
	
	fmt.Println("Life is Good: %s", u.Name)
}
```
- Although type as context can help in dealing with multiple custom errors and provide more context, it brings **tight coupling** with the concrete types
- prefer to rely on `error` interface without having to use type as context always
- to be able to provide more context, prefer wrapping error in Go or use type as behavior

## Moving from *type as context* to *type as behavior* for less coupling

### Type as Context 
```go
package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
)


// client represents a single connection in the room.
type client struct {
	name   string
	reader *bufio.Reader
}

// TypeAsContext shows how to check multiple types of possible custom error
// types that can be returned from the net package.
func (c *client) TypeAsContext() {
	for {
		line, err := c.reader.ReadString('\n')
		if err != nil {
			switch e := err.(type) {
			case *net.OpError:
				if !e.Temporary() {
					log.Println("Temporary: Client leaving chat")
					return
				}

			case *net.AddrError:
				if !e.Temporary() {
					log.Println("Temporary: Client leaving chat")
					return
				}

			case *net.DNSConfigError:
				if !e.Temporary() {
					log.Println("Temporary: Client leaving chat")
					return
				}

			default:
				if err == io.EOF {
					log.Println("EOF: Client leaving chat")
					return
				}

				log.Println("read-routine", err)
			}
		}

		fmt.Println(line)
	}
}

```

### Behavior as Context
- the "net" package has it's own un-exported `temporary` interface defined which is used internally
```go
package net

type temporary interface {
	Temporary() bool
}

func (e *OpError) Temporary() bool {
	// Treat ECONNRESET and ECONNABORTED as temporary errors when
	// they come from calling accept. See issue 6163.
	if e.Op == "accept" && isConnError(e.Err) {
		return true
	}

	if ne, ok := e.Err.(*os.SyscallError); ok {
		t, ok := ne.Err.(temporary)
		return ok && t.Temporary()
	}
	t, ok := e.Err.(temporary)
	return ok && t.Temporary()
}

```
- defining our own `temporary` interface to define behavior that we are bothered about
```go
package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
)


// temporary is declared to test for the existence of the method coming
// from the net package.
type temporary interface {
	Temporary() bool
}

// BehaviorAsContext shows how to check for the behavior of an interface
// that can be returned from the net package.
func (c *client) BehaviorAsContext() {
	for {
		line, err := c.reader.ReadString('\n')
		if err != nil {
			switch e := err.(type) {
			case temporary:
				if !e.Temporary() {
					log.Println("Temporary: Client leaving chat")
					return
				}

			default:
				if err == io.EOF {
					log.Println("EOF: Client leaving chat")
					return
				}

				log.Println("read-routine", err)
			}
		}

		fmt.Println(line)
	}
}
```

## Find the bug
```go
package main

package main

import "fmt"

type customError struct {
	S string
}

func (c customError) Error() string {
	return c.S
}

func foo() ([]int, *customError) {
	return nil, nil
}

func main() {
	var err error
	if _, err = foo(); err != nil { // the go compiler also gives warning in this line saying : "tautological condition: non-nil != nil"
		fmt.Println("failed")
		return
	}
	fmt.Println("no error")
}

// OUTPUT: 
// failed
```
- this happens because, `foo` function returns a type `*customError` instead of `error`
- the `nil` of type `customError` struct (that implements `error` interface) is stored in a variable `err` of type `error` interface
- any variable of interface type has two words in it
	- one word is for the concrete implementation type
	- one word is reference to the actual concrete value of the concrete type
- at line `var err error`, the concrete type has `nil` and concrete value has `nil`, which is the zero value of any interface and thus, `err == nil` will be true
	- ![[Pasted image 20251212103818.png]]
- at line `err = foo()` the `err` now stores a concrete error where the concrete type is not `nil` but of  `*customError` type, and concrete value is `nil`, thus it is non-zero value of `error` interface, and thus `err != nil` will be true
	- ![[Pasted image 20251212103913.png]]
- thus `foo()` function return type being a `*customError` instead of `error` is causing this bug
- fix
```go
package main

package main

import "fmt"

type customError struct {
	S string
}

func (c customError) Error() string {
	return c.S
}

func foo() ([]int, error) {
	return nil, nil
}

func main() {
	var err error
	if _, err = foo(); err != nil {
		fmt.Println("failed")
		return
	}
	fmt.Println("no error")
}

// OUTPUT:
// no error
```


## Wrapping errors
- A design pattern that can help minimize problems in code and help handling errors
- handling an error means
	- if an error occurs, is it going to stop there with the code that's handling it 
		- does the error gets propagated any further
		- should we recovery or shut down from this point
	- logging error
		- log error with full context

```go
package main

import (
	"errors"
	"fmt"
)

// AppError represents a custom error type.
type AppError struct {
	State int
}

// Error implements the error interface.
func (c *AppError) Error() string {
	return fmt.Sprintf("App Error, State: %d", c.State)
}

// Cause iterates through all the wrapped errors
// until the root error value is reached.
func Cause(err error) error {
	root := err
	for {
		if err = errors.Unwrap(root); err == nil {
			return root
		}
		root = err
	}
}

func main() {

	// Make the function call and validate the error.
	if err := firstCall(10); err != nil {

		// How to use the As function.
		var ap *AppError
		if errors.As(err, &ap) {
			fmt.Println("As says it is an AppError")
		}

		// Use type as context to determine cause.
		switch v := Cause(err).(type) {
		case *AppError:

			// We got our custom error type.
			fmt.Println("Custom App Error:", v.State)

		default:

			// We did not get any specific error type.
			fmt.Println("Default Error")
		}

		// Display the error.
		fmt.Println("\n********************************")
		fmt.Printf("%v\n", err)
	}
}

// firstCall makes a call to a second function and wraps any error.
func firstCall(i int) error {
	if err := secondCall(i); err != nil {
		return fmt.Errorf("firstCall->secondCall(%d) : %w", i, err)
	}
	return nil
}

// secondCall makes a call to a third function and wraps any error.
func secondCall(i int) error {
	if err := thirdCall(); err != nil {
		return fmt.Errorf("secondCall->thirdCall() : %w", err)
	}
	return nil
}

// thirdCall create an error value we will validate.
func thirdCall() error {
	return &AppError{99}
}

```