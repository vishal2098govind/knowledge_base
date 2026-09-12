package main

import "fmt"

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
}
