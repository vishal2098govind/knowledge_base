package main

import (
	"fmt"
	"slices"
	"sync"
)

func main() {
	p := Player{Name: "Vishal"}
	p.Found("copper")
	p.Found("copper")
	fmt.Println(p.Keys)
}

type Player struct {
	mu   sync.Mutex
	Keys []string
	Name string
	Item
}

func (p *Player) Found(k string) error {
	if k != "copper" && k != "jade" && k != "gold" {
		return fmt.Errorf("invalid key")
	}

	p.mu.Lock()
	// CS
	found := slices.Contains(p.Keys, k)
	if !found {
		p.Keys = append(p.Keys, k)
	}
	p.mu.Unlock()

	return nil
}

type Item struct {
	X int
	Y int
}

func (i *Item) Move(dx, dy int) {
	i.X += dx
	i.Y += dy
}
