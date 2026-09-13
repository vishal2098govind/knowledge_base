#engineers-notebook #golang #error-handling #panics

### We can catch a panic (`recover`) only within `defer`
```go
package main

import "fmt"

func main() {
	fmt.Println(safeDiv(7, 2))
	fmt.Println(safeDiv(7, 0))
}

func safeDiv(a, b int) (_ int, err error) {
	defer func() {
		e := recover()
		if e != nil {
			err = fmt.Errorf("%v", err)
		}
	}()

	return div(a, b), nil
}

func div(a, b int) int {
	return a / b
}
```