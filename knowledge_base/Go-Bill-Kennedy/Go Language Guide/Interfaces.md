#go #go-interface

> Polymorphism means that you write a certain program and it behaves differently depending on the data it operates on.

```go
type reader interface {
	read_good(b []byte) (int, error)
	read_bad(n int) ([]byte, error)
}
```

- why is `read_bad` bad and `read_good` good for Go?
	- every time `read_bad` is invoked, there has to be a new allocation that has to be made on heap for the bytes to be returned (determined during escape analysis), there's no other way


## Polymorphism
```go
type reader interface {
	read(b []byte) (int, error)
}

type file struct {
	name string
}
func (file) read(b []byte) (int, error) {
	
}

type pipe struct {

}
func (pipe) read(b []byte) (int, error) {

}

func retrieve(r reader) error {
	data := make([]byte, 100)
	len, err := r.read(data)
	if err != nil {
		return err
	}
	
	fmt.Println(string(data[:len]))
	return nil
}

func main() {
	f := file{"data.json"}
	p := pipe("cfg_service")
	
	retrieve(f)
	retrieve(p)
}
```

## Method Sets
#go-method-sets
```go
package main

import "fmt"

type notifier interface {
	notify()
}

type user struct {
	name  string
	email string
}

func (u *user) notify() {
	fmt.Printf("Sending user email to:%s<%s>", u.name, u.email)
}

func sendNotification(n notifier) {
	n.notify()
}

func main() {
	u := user{name: "Vishal", email: "vishal@gmail.com"}
	
	// cannot use u (variable of struct type user) as notifier value in argument to
	// sendNotification: user does not implement notifier (method notify has pointer receiver)
	// sendNotification(u)
	sendNotification(&u)
}
```

### Method sets rules
- The Go compiler attaches methods of a struct differently to pointers and values based on the whether the method is defined using pointer or value semantics
	- To a value of type struct (e.g. T) (not pointer of type struct e.g. \*T), the Go compiler only attaches the methods using value semantics and not pointer semantics (e.g. \*T) (i.e. methods with value receivers e.g. T, and not pointer receivers)
		- This helps in achieving **data integrity**
		- If `sendNotification(u)` would have been allowed by the Go compiler, 
			- then if `notify` implementation in `user` would change some fields on the receiving `u`, 
			- which is not something the caller of `sendNotification(u)` would be expecting while **passing a copy** of `u` to `sendNotification(u)`, 
			- then the lines of code, below the function call to `sendNotification(u)`, in the caller would be seeing new values of fields of `u` which `notify` method would have changed
			- Thus, only allowing `sendNotification(&u)` lets the caller of the `sendNotification` function know that they are **sharing `u`** and **not passing a copy of `u`** and thus any changes to the fields of `u` would not be un-expected
	- To a value of type pointer to the struct (e.g. \*T) (not value of type struct e.g. T), the Go compiler attaches all the method, be it using value or pointer semantics
![[Pasted image 20251209002720.png]]

### Storage by value
- When we store a value, the interface value has its own copy of the value. Changes to the original value will not be seen.
- When we store a pointer, the interface value has its own copy of the address. Changes to the original value will be seen.
```go
package main

import "fmt"

// printer displays information.
type printer interface {
	print()
}

// cannon defines a cannon printer.
type cannon struct {
	name string
}

// print displays the printer's name.
func (c cannon) print() {
	fmt.Printf("Printer Name: %s\n", c.name)
}

// epson defines a epson printer.
type epson struct {
	name string
}

// print displays the printer's name.
func (e *epson) print() {
	fmt.Printf("Printer Name: %s\n", e.name)
}

func main() {

	// Create a cannon and epson printer.
	c := cannon{"PIXMA TR4520"}
	e := epson{"WorkForce Pro WF-3720"}

	// Add the printers to the collection using both
	// value and pointer semantics.
	printers := []printer{

		// Store a copy of the cannon printer value.
		c,

		// Store a copy of the epson printer value's address.
		&e,
	}

	// Change the name field for both printers.
	c.name = "PROGRAF PRO-1000"
	e.name = "Home XP-4100"

	// Iterate over the slice of printers and call
	// print against the copied interface value.
	for _, p := range printers {
		p.print()
	}

	// When we store a value, the interface value has its own
	// copy of the value. Changes to the original value will
	// not be seen.

	// When we store a pointer, the interface value has its own
	// copy of the address. Changes to the original value will
	// be seen.
}
```
![[Pasted image 20251209011648.png]]