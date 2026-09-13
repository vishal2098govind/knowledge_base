package main

import "fmt"

// Though exercise: Sorting

func Sort(s Sortable) {

}

type Sortable interface {
	Less(i, j int) bool
	Swap(i, j int)
	Len() int
}

func main() {
	var a any
	a = "Hi"
	i := a.(string)
	fmt.Println(i)

	// j := a.(int)
	j, ok := a.(int)
	if ok {
		fmt.Println(j)
	} else {
		fmt.Println("")
	}
}
