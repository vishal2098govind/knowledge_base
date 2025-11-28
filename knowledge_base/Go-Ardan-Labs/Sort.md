#go #sort

```go
package main

import (
	"fmt"
	"sort"
)

func main() {
	SortByAge([]Person{
		{Age: 20},
		{Age: 2},
		{Age: 21},
	})
	name := func(p1, p2 *Planet) bool {
		return p1.Name < p2.Name
	}
	mass := func(p1, p2 *Planet) bool {
		return p1.Mass < p2.Mass
	}
	planets := []Planet{
		{"Mercury", 0.055, 0.4},
		{"Venus", 0.815, 0.7},
		{"Earth", 1.0, 1.0},
		{"Mars", 0.107, 1.5},
	}
	By(name).SortPlanets(planets)
	By(mass).SortPlanets(planets)
}

type Person struct {
	Age int
}

type ByAge []Person

func (a ByAge) Len() int           { return len(a) }
func (a ByAge) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByAge) Less(i, j int) bool { return a[i].Age < a[j].Age }

func SortByAge(p []Person) {
	sort.Sort(ByAge(p))
	fmt.Println(p)
}


type Planet struct {
	Name     string
	Mass     float64
	Distance float64
}

type By func(p1, p2 *Planet) bool

type planetSorter struct {
	Planets []Planet
	By      By
}

func (a planetSorter) Len() int           { return len(a.Planets) }
func (a planetSorter) Swap(i, j int)      { a.Planets[i], a.Planets[j] = a.Planets[j], a.Planets[i] }
func (a planetSorter) Less(i, j int) bool { return a.By(&a.Planets[i], &a.Planets[j]) }

func (by By) SortPlanets(p []Planet) {
	sorter := planetSorter{
		Planets: p,
		By:      by,
	}
	sort.Sort(sorter)
	fmt.Println(sorter.Planets)
}

```