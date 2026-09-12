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
