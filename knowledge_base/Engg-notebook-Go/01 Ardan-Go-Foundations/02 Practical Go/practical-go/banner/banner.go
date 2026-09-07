package main

import (
	"fmt"
	"strings"
	"unicode/utf8"
)

// banner("Go", 6) -> output
//   Go
// ------

func main() {
	banner("Go", 6)
	banner("G♡", 6)

	fmt.Println("len:", len("G♡"))
	fmt.Println("len:", len("Go"))
	fmt.Println("s[1]:", "G♡"[1])
	fmt.Printf("s[1]:%c\n", "G♡"[1])
	for i, c := range "G♡" {
		fmt.Printf("%c at %d\n", c, i)
	}
}

func banner(text string, width int) {
	padding := (width - utf8.RuneCount([]byte(text))) / 2
	padStr := strings.Repeat(" ", padding)
	fmt.Printf("%s%s\n", padStr, text)
	fmt.Printf("%s\n", strings.Repeat("-", width))
}
