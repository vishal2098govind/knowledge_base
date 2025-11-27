#go #pointer-semantics #value-semantics #escape-analysis


```go
package main

import "fmt"

func main() {
	fmt.Println(NewItem(10, 20))
	fmt.Println(NewItem(200, 20))
}

// Item factory (value semantics)
func NewItem(a int, b int) (Item, error) {
	if a < maxX || a > maxX || b < maxY || b > maxY {
		return Item{}, fmt.Errorf("%d/%d out of bound %d/%d", a, b, maxX, maxY)
	}

	i := Item{X: a, Y: b}

	return i, nil
}

// Item factory (pointer semantics)
func NewItemPtr(a int, b int) (*Item, error) {
	if a < maxX || a > maxX || b < maxY || b > maxY {
		return nil, fmt.Errorf("%d/%d out of bound %d/%d", a, b, maxX, maxY)
	}

	i := Item{X: a, Y: b}
	// go compiler here does `escape analysis` and will allocate memory for i in the heap
	return &i, nil
}

const maxX = 400
const maxY = 600

type Item struct {
	X int
	Y int
}
```
### Escape Analysis
Heap allocation during **pointer-semantics**
```
➜ go build -gcflags=-m main.go
# command-line-arguments
./main.go:6:13: inlining call to fmt.Println
./main.go:12:28: ... argument does not escape
./main.go:12:57: a escapes to heap
./main.go:12:60: b escapes to heap
./main.go:12:63: 400 escapes to heap
./main.go:12:69: 600 escapes to heap
./main.go:6:13: ... argument does not escape
./main.go:6:13: .autotmp_0 escapes to heap
./main.go:26:2: moved to heap: i
./main.go:23:25: ... argument does not escape
./main.go:23:54: a escapes to heap
./main.go:23:57: b escapes to heap
./main.go:23:60: 400 escapes to heap
./main.go:23:66: 600 escapes to heap
```

- Heap allocations take more time than stack allocation
- symbols allocated in heap put more pressure on the garbage collector as it needs to clean them up
- Thus, if we can work with **value-semantics** and only pass things by value, it would be better performant mostly and less pressure on garbage collector
- `time` package in go uses **value-semantics**
```go
package main

import (
	"fmt"
	"time"
)

func main() {
	t := time.Time{}
	fmt.Println(t)
	t = t.Add(1 * time.Hour)
	fmt.Println(t)
}
```

- In value-semantics : everyone has their own copy
- In pointer-semantics : everyone shares the same copy on the heap, and sometimes we need to lock it while synchronisation

- Value semantics in case of receivers:
```go
package main

import (
	"fmt"
)

func main() {
	p := Point{}
	p.Move(1, 2)
	fmt.Println(p) // {0, 0}
}
type Point struct {
	X int
	Y int
}

func (p Point) Move(dx, dy int) {
	// here the receiver is also just a copy in case of value semantics
	p.X += dx // ineffective assignment to field Point.X (SA4005)go-staticcheck
	p.Y += dy // ineffective assignment to field Point.Y (SA4005)go-staticcheck
}
```
- Using pointer semantics for receiver will not complain
```go
package main

import (
	"fmt"
)

func main() {
	p := Point{}
	p.Move(1, 2)
	fmt.Println(p) // {1, 2}
}
type Point struct {
	X int
	Y int
}

func (p *Point) Move(dx, dy int) {
	// using pointer semantics for receiver
	p.X += dx
	p.Y += dy
}
```

### When to use value vs pointer receiver
- In general, use value semantics: easier to debug and understand having an own copy
- Try to keep same semantics on all methods as well
- If we start with value receiver method for a struct, stick to that for other methods as well
- There are some cases where we must use pointer semantics or pointer receivers
	- If we have a lock, like mutex, as a field of a struct, because copying a lock variable is always a BUG
		- gRPC has locks on all the data it generates, thus gRPC code generations always use with pointer semantics
	- If we need to mutate underlying struct
	- While decoding or unmarshalling
		- the `time` package uses pointer semantics when it has to unmarshal json