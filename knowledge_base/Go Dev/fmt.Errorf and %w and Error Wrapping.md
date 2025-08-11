#go #error-handling #fmt-Errorf #error-chaing #error-verb

### Wrapping errors with `fmt.Errorf`

```go
package main

func Connect() error {
	return errors.New("connection failed")
}

func CreateUser() error {
	err := Connect()
	if err != nil {
		return err
	}
	
	// ... continue on
	return nil
}

func main() {
	err := CreateUser()
	if err != nil {
		fmt.Println(err)
	}
}
// Output:
// connection failed
// NOTE: simply returning err would be confusing to know where the error actually came from
```
Wrapping #wrapping-errors  with `fmt.Errorf` using `%w` ( #error-verb) helps to add bit of context to make the errors easier to comprehend when we're looking at the logs or anything like that, that we might need to actually figure out well why this error actually occurred
```go
package main

func Connect() error {
	return errors.New("connection failed")
}

func CreateUser() error {
	err := Connect()
	if err != nil {
		return fmt.Errorf("create user: %w", err)
	}
	
	// ... continue on
	return nil
}

func CreateOrg() error {
	err := CreateUser()
	if err != nil {
		return fmt.Errorf("create org: %w", err)
	}
}

func main() {
	err := CreateOrg()
	if err != nil {
		fmt.Println(err)
	}
}
// Output:
// create org: create user: connection failed
// NOTE: Now we have this nice daisy chain of errors, so we know that while creating an org, there was a user being created and the connection failed at that point of time
```

### What is special in `fmt.Errorf` and `%w` 
- `%w` is kind of a special way of telling where the error is going to be placed in the message
- `w` in `%w` stands for **wrapping** (i think!)
- `fmt.Errorf("%w", err)` is also like `fmt.Printf("%v", err)` 
- In `fmt.Printf` there's no `%w` and only the `fmt.Errorf` has this `%w` which is a special one for errors
- but the `fmt.Errorf` is also going to **retain some information** about the error. This is achieved via **Error Wrapping** feature
- Sometimes additional information or context is needed to be added to an error **without creating a custom error**
- According to go docs, an `error e`  **wraps** another `error` if the type of  `e` has one of the methods:
	- `Unwrap() error`
	- `Unwrap() []error`
- **error chain** refers to series of errors that are wrapped together
- Since [Go 1.13](https://go.dev/blog/go1.13-errors), the `fmt.Errorf` function supports a new **`%w` verb**. When this verb is present, the error returned by `fmt.Errorf` will have an `Unwrap` method returning the argument of `%w`, which must be an error. In all other ways, `%w` is identical to `%v`.
```
if err != nil {
    // Return an error which unwraps to err.
    return fmt.Errorf("decompress %v: %w", name, err)
}
```
- Wrapping an error with `%w` makes it available to `errors.Is` and `errors.As`:
```
err := fmt.Errorf("access denied: %w", ErrPermission)
...
if errors.Is(err, ErrPermission) ...
```
- Use `errors.As` when dealing with errors which are wrapped with a type/interface with some specific methods that we might want to use while handling the error
Example:
```go
package main

import (
	"database/sql"
	"errors"
	"fmt"
)

func GetStudents() error {
	return sql.ErrNoRows
}

type temporary interface {
	Temporary() bool
}

func GetSchool() error {
	return fmt.Errorf("error getting students: %w", GetStudents())
}

type temp struct {
	Err error
}

func (*temp) Temporary() bool {
	return true;
}

func (t *temp) Error() string {
	return "temp error"
}

func GetTempSchool() error {
	retur fmt.Errorf("error getting temp school: %w", &temp{Err: fmt.Errorf("temp school error")})
}

func main() {
	err := GetSchool()
	if errors.Is(err, sql.ErrNoRows) {
		fmt.Println("not found")
	}

	err := GetTempSchool()
	var t temporary
	if errors.As(err, &t) {
		if t.Temporary() {
			fmt.Println("this is temporary")
		}
	}
}

// Output:
// not found
```

### Whether to Wrap
When adding additional context to an error, either with `fmt.Errorf` or by implementing a custom type, you need to decide whether the new error should wrap the original. There is no single answer to this question; it depends on the context in which the new error is created. Wrap an error to expose it to callers. Do not wrap an error when doing so would expose implementation details.

As one example, imagine a `Parse` function which reads a complex data structure from an `io.Reader`. If an error occurs, we wish to report the line and column number at which it occurred. If the error occurs while reading from the `io.Reader`, we will want to wrap that error to allow inspection of the underlying problem. Since the caller provided the `io.Reader` to the function, it makes sense to expose the error produced by it.

In contrast, a function which makes several calls to a database probably should not return an error which unwraps to the result of one of those calls. If the database used by the function is an implementation detail, then exposing these errors is a violation of abstraction. For example, if the `LookupUser` function of your package `pkg` uses Go’s `database/sql` package, then it may encounter a `sql.ErrNoRows` error. If you return that error with `fmt.Errorf("accessing DB: %v", err)` then a caller cannot look inside to find the `sql.ErrNoRows`. But if the function instead returns `fmt.Errorf("accessing DB: %w", err)`, then a caller could reasonably write

```go
err := pkg.LookupUser(...)
if errors.Is(err, sql.ErrNoRows) …
```

At that point, the function must always return `sql.ErrNoRows` if you don’t want to break your clients, even if you switch to a different database package. In other words, wrapping an error makes that error part of your API. If you don’t want to commit to supporting that error as part of your API in the future, you shouldn’t wrap the error.

It’s important to remember that whether you wrap or not, the error text will be the same. A _person_ trying to understand the error will have the same information either way; the choice to wrap is about whether to give _programs_ additional information so they can make more informed decisions, or to withhold that information to **preserve an abstraction layer**.

- when you return an error from another package you should convert the error to a form that does not expose the underlying error, unless you are willing to commit to returning that specific error in the future.

```
f, err := os.Open(filename)
if err != nil {
    // The *os.PathError returned by os.Open is an internal detail.
    // To avoid exposing it to the caller, repackage it as a new
    // error with the same text. We use the %v formatting verb, since
    // %w would permit the caller to unwrap the original *os.PathError.
    return fmt.Errorf("%v", err)
}
```

- If a function is defined as returning an error wrapping some sentinel ( #sentinel-errors) l or type, do not return the underlying error directly.

```
var ErrPermission = errors.New("permission denied")

// DoSomething returns an error wrapping ErrPermission if the user
// does not have permission to do something.
func DoSomething() error {
    if !userHasPermission() {
        // If we return ErrPermission directly, callers might come
        // to depend on the exact error value, writing code like this:
        //
        //     if err := pkg.DoSomething(); err == pkg.ErrPermission { … }
        //
        // This will cause problems if we want to add additional
        // context to the error in the future. To avoid this, we
        // return an error wrapping the sentinel so that users must
        // always unwrap it:
        //
        //     if err := pkg.DoSomething(); errors.Is(err, pkg.ErrPermission) { ... }
        return fmt.Errorf("%w", ErrPermission)
    }
    // ...
}
```
