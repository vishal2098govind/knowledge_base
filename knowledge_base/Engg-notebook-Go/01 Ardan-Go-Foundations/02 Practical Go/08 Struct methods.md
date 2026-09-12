#engineers-notebook #golang #structs #value-semantics #pointer-semantics

## Value vs Pointer semantics with method receivers
- In general use value semantics
- Try to keep same semantics for all methods. If we start with pointer receiver, use it across all methods
- When u **must use pointer receiver**
	- Need to synchronize on a field of the struct, and thus need a lock (like mutex) as a field of the struct
	- If we need to mutate the struct
	- Un-marshaling
```go
func foo() {
	var t Item
	var tjson []byte
	t.Unmarshal(tjson)
	t. // process further with t
}

type Item { ... }
func (i *Item) Unmarshal(b []byte) error {
	// ...
}
```

```go
package main

import "fmt"

func main() {
	i := Item{10, 20}
	fmt.Println(i) // {10, 20}
	i.Move(4, 2)

	// i is still {10, 20} since the Move method uses value semantics in it's receiver and not pointer semantics
	fmt.Println(i)

	fmt.Println(i) // {10, 20}
	i.MovePtr(4, 2)

	// i is now {14, 22} since the MovePtr method uses pointer semantics in it's receiver and not value semantics
	fmt.Println(i)
}

type Item struct {
	X int
	Y int
}

func (i Item) Move(dx, dy int) {
	i.X += dx
	i.Y += dy
}

func (i *Item) MovePtr(dx, dy int) {
	i.X += dx
	i.Y += dy
}
```