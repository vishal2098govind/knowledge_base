#engineers-notebook #golang #structs #go-docs

```go
type Item {X int; Y int}

// all possible "new" functions or factory functions of a struct:
func New(x, y int) Item
func New(x, y int) *Item
func NewItem(x, y int) Item
func NewItem(x, y int) *Item
func NewItem(x, y int) (*Item, error)
func NewItem(x, y int) (Item, error)
```

even though these have different function names, go docs tool groups all of them under `type Item` as it's constructors or factories
check https://pkg.go.dev/time#Time `type Time` in time package.
It has all different functions under type Time
- `func Date(year int, month Month, day, hour, min, sec, nsec int, loc *Location) Time
- `func Parse(layout, value string) (Time, error)`
- `func Now() Time`
- `func ParseInLocation(layout, value string, loc *Location) (Time, error)`
- `func Unix(sec int64, nsec int64) Time`
- `func UnixMicro(usec int64) Time`
- `func UnixMilli(msec int64) Time`
![[Pasted image 20260912114625.png]]


## Value semantics factories less pressure on Garbage collection
Go compiler does escape analysis in a factory function while returning a pointer and thus does heap allocation instead of stack. Thus returning a value helps staying performant and creating less pressure on garbage collector to clean orphan heap allocations.
For **pointer semantics**, we need to think about **locking** to **synchronize** access if need.

```go
package main

import "fmt"

func main() {
	i, err := NewItem(30, 20)
	fmt.Println(i, &i.x, err)
	i2, err := NewItemPtr(40, 20)
	fmt.Println(i2, &i2.x, err)
}

type Item struct {
	x int
	y int
}

// value semantics: returns copy of Item up the stack instead of returning pointer to Item, thus no escape analysis of i to heap instead of stack
func NewItem(x, y int) (Item, error) {
	if x == 10 {
		return Item{}, fmt.Errorf("invalid x: %v", x)
	}

	i := Item{x, y}
	fmt.Println(&i.x)
	return i, nil
}

// pointer semantics: returns pointer to Item up the stack instead of copy of Item, thus escape analysis of i to heap instead of stack
func NewItemPtr(x, y int) (*Item, error) {
	if x == 10 {
		return nil, fmt.Errorf("invalid x: %v", x)
	}

	i := Item{x, y}
	fmt.Println(&i.x)
	return &i, nil
}

```
#### output:
```sh
➜  structs git:(vishal-kb) ✗ go run .
0xc000092030
{30 20} 0xc000092020 <nil>
0xc000092060
&{40 20} 0xc000092060 <nil>
```
### escape analysis by garbage collector:
```go
package main

import "fmt"

func main() {
	i, err := NewItem(30, 20)
	fmt.Println(i, &i.x, err)
	i2, err := NewItemPtr(40, 20)
	fmt.Println(i2, &i2.x, err)
}

type Item struct {
	x int
	y int
}

// value semantics: returns copy of Item up the stack instead of returning pointer to Item, thus no escape analysis of i to heap instead of stack
func NewItem(x, y int) (Item, error) {
	if x == 10 {
		return Item{}, fmt.Errorf("invalid x: %v", x)
	}

	i := Item{x, y}
	// fmt.Println(&i.x) # commenting for avoiding heap allocation due to passing to println
	return i, nil
}

// pointer semantics: returns pointer to Item up the stack instead of copy of Item, thus escape analysis of i to heap instead of stack
func NewItemPtr(x, y int) (*Item, error) {
	if x == 10 {
		return nil, fmt.Errorf("invalid x: %v", x)
	}

	i := Item{x, y}
	fmt.Println(&i.x)
	return &i, nil
}

```

```sh
go build -gcflags=-m .
# practical-go/structs
./structs.go:35:13: inlining call to fmt.Println
./structs.go:7:13: inlining call to fmt.Println
./structs.go:9:13: inlining call to fmt.Println
./structs.go:20:28: ... argument does not escape
./structs.go:20:46: x escapes to heap
./structs.go:34:2: moved to heap: i # i in NewItemPtr
./structs.go:31:25: ... argument does not escape
./structs.go:31:43: x escapes to heap
./structs.go:35:13: ... argument does not escape
./structs.go:6:2: moved to heap: i
./structs.go:7:13: ... argument does not escape
./structs.go:7:14: i escapes to heap
./structs.go:9:13: ... argument does not escape
```