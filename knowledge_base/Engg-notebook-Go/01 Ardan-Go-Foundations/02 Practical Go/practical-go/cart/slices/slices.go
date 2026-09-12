package main

import (
	"fmt"
	"sort"
)

func main() {
	s1 := make([]int, 10)
	s2 := s1[3:7]
	fmt.Println(s1)
	fmt.Println(s2)
	s2 = append(s2, 10)
	fmt.Println(&s1[3], &s2[0])
	fmt.Println(s1)
	fmt.Println(s2)
	s2 = append(s2, 20)
	fmt.Println(&s1[3], &s2[0])
	fmt.Println(s1)
	fmt.Println(s2)
	s2 = append(s2, 30)
	fmt.Println(&s1[3], &s2[0])
	fmt.Println(s1)
	fmt.Println(s2)
	s2 = append(s2, 40)
	fmt.Println(&s1[3], &s2[0])
	fmt.Println(s1)
	fmt.Println(s2)

	fmt.Println(concat([]string{"A", "B", "C"}, []string{"D", "E"}))
	fmt.Println(median([]float64{10, 8, 9}))
	v := []float64{10, 8, 9, 6}
	fmt.Println(median(v))
	fmt.Println(v)

	// value semantics in for range loop without index:
	players := []Player{{"John", 10}, {"Doe", 20}}
	for _, p := range players {
		p.Score += 100
	}
	fmt.Println(players)
	for i := range players {
		players[i].Score += 100
	}
	fmt.Println(players)
}

type Player struct {
	Name  string
	Score int
}

func concat(s1, s2 []string) []string {
	s3 := make([]string, len(s1)+len(s2))
	copy(s3, s1)
	copy(s3[len(s1):], s2)
	return s3
}

// sort and find middle value
// if odd length -> return middle
// if even length -> return average of two middle elements
func median(values []float64) float64 {
	v := make([]float64, len(values))
	copy(v, values)
	sort.Slice(v, func(i, j int) bool {
		return v[i] < v[j]
	})
	len := len(v)
	mid := len / 2
	fmt.Println(v, v[mid])
	if len%2 == 0 {
		return (v[mid] + v[mid-1]) / 2
	}
	return v[mid]
}
