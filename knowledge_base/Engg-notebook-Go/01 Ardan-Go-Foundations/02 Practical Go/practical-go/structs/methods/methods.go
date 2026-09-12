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
