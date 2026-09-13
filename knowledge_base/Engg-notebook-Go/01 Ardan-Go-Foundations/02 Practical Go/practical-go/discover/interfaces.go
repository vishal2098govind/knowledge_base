package main

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
	i := a.(string)

}
