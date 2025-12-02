#go #slices #arrays #value-semantics #pointer-semantics 

```go
var fruits [5]string // array
fruits[0] = "apple"
fruits[1] = "banana"
...
fruits[4] = "orange"

// value semantic version of for loop
for i, fruit range fruits {
	fmt.Println(i, fruit)
}
```

![[Pasted image 20251202114238.png]]

## Range Mechanics
### for-range loop using pointer semantics
```go
package main

import fmt

func main() {
	friends := [5]string{"Annie", "Betty", "Charley", "Dog", "Edward"}
	fmt.Printf("Bfr[%s] : ", friends[1])
	
	for i := range friends {
		friends[1] = "Jack"
		if i == 1 {
			fmt.printf("Aft[%s]\n", friends[1])
		}
	}
}


// Output:
// Bfr[Betty] : Aft[Jack]
```
### for-range loop using value semantics
- this for-rage loop is not iterating over the original underlying friends array, but over a COPY, thus `v` is taken from that COPY.
```go
package main

import fmt

func main() {
	friends := [5]string{"Annie", "Betty", "Charley", "Dog", "Edward"}
	fmt.Printf("Bfr[%s] : ", friends[1])
	
	// this for-rage loop is not iterating over the original underlying friends array, but over a COPY
	for i, v := range friends {
		friends[1] = "Jack"
		if i == 1 {
			fmt.printf("Aft[%s]\n", v)
		}
	}
}


// Output:
// Bfr[Betty] : Aft[Betty]
```
### Another example
```go
package main

import "fmt"

func main() {
	// Using value-semantic form of for-range
	friends := []string{"A", "B", "C", "D", "E"}
	for _, v := range friends {
		friends = friends[:2]
		fmt.Printf("v[%s]\n", v)
	}

	// Using pointer-semantic form of for-range
	friends = []string{"A", "B", "C", "D", "E"}
	for i := range friends {
		friends = friends[:2]
		fmt.Printf("v[%s]\n", friends[i])
	}
}

// OUTPUT:
// v[A]
// v[B]
// v[C]
// v[D]
// v[E]
// v[A]
// v[B]
// panic: runtime error: index out of range [2] with length 2
```
![[Pasted image 20251203011322.png]]
## Slices
- can create `slices` using built-in function `make`
- `make` is specifically used to create/pre-allocate these data-structures only:
	- `slices`
	- `maps`
	- and `channels`
- empty struct type in Go: `struct{}`
	- this is a zero allocation type
	- a value or a global variable that is embedded into the runtime
```go
package main

import fmt

func main() {
	
	fruits := make([]string, 5, 8) // allocates a slice of len = 5, cap = 8. underlying array is of size 8 (cap)
	
	fruits[0] = "apple"
	
	var data []string // declare a nil slice of strings with underlying array pointing to nil
	u := []string{} // empty slice of strings with underlying array pointing to empty struct `struct{}`
	
	data = append(data, "apple") // append function takes in a value-copy of data slice (value-semantics mutation)
	// if underlying array is full before appending (len == cap), 
	// then it allocates fresh underlying array of double the capacity, 
	// copies all items there and updates the internal references of slices to the new array, 
	// leaving the old array for GC
	
}
```

### 3-index slice
- use a 3-index slice to cut out a slice from an existing slice to be able to mention the capacity of the underlying slice as well
```go
package main

import "fmt"

func main() {
	slice1 := make([]string, 6, 8)
	slice1[0] = "apple"
	slice1[1] = "orange"
	slice1[2] = "grapes"
	slice1[3] = "banana"
	slice1[4] = "mango"
	slice1[5] = "watermelon"
	fmt.Printf("slice1: %v, cap(slice1): %d\n", slice1, cap(slice1))

	slice2 := slice1[2:4]
	fmt.Printf("slice2: %v, cap(slice2): %d\n", slice2, cap(slice2))

	slice2 = append(slice2, "CHANGE1")
	fmt.Printf("slice2: %v, cap(slice2): %d\n", slice2, cap(slice2))
	fmt.Printf("slice1: %v, cap(slice1): %d\n", slice1, cap(slice1))

	slice3 := slice1[2:4:4]
	fmt.Printf("slice3: %v, cap(slice3): %d\n", slice3, cap(slice3))
	
	slice3 = append(slice3, "CHANGE2") // now, slice3 is no longer referring to that backing array of slice1 anymore since it's len=cap was 2 and appending led to allocation of a fresh underlying array
	fmt.Printf("slice3: %v, cap(slice3): %d\n", slice3, cap(slice3))
	fmt.Printf("slice1: %v, cap(slice1): %d\n", slice1, cap(slice1))
}
// OUTPUT:
// slice1: [apple orange grapes banana mango watermelon], cap(slice1): 8
// slice2: [grapes banana], cap(slice2): 6
// slice2: [grapes banana CHANGE1], cap(slice2): 6
// slice1: [apple orange grapes banana CHANGE1 watermelon], cap(slice1): 8
// slice3: [grapes banana], cap(slice3): 2
// slice3: [grapes banana CHANGE2], cap(slice3): 4
// slice1: [apple orange grapes banana CHANGE1 watermelon], cap(slice1): 8
```

## Strings and Slices
#strings 
- Strings are represented in UTF-8 in Go
- Go source code is always UTF-8
- `for i, c range s {}` loops over each **rune**
	- a **rune** in Go is `int32` i.e. of **4 bytes**
- `len()` and `s[i:j]`(slicing a string)  deal with each **byte**.
	- a byte in Go is `uint8`
> A Go quote: **"Every array in Go is just a slice waiting to happen!!!"**
- slicing over a string is done over with individual underlying bytes
```go
package main

import (
	"fmt"
	"unicode/utf8"
)

func main() {
	s := "G♡!"
	var buf [utf8.UTFMax]byte // array, not slice (since utf8.UTFMax is a const value 4)

	for i, r := range s {

		// capture the number of bytes associated with this rune
		rl := utf8.RuneLen(r)

		// calculate the slice offset for the bytes associated with this rune
		soff := i + rl
		
		// copy of rune bytes from string to buffer
		// NOTE: copy needs both src and dst to be of slice, not arrays
		// copy(buf, s[i:soff]) -> ERROR
		// ERROR Message: invalid argument: copy expects slice arguments; found buf (variable of type [4]byte) and s[i:soff] (value of type string)
		// Thus, by doing `buf[:]` we are passing a slice constructed basing on the buf array to the copy
		// A Go quote: "Every array in Go is just a slice waiting to happen!!!"
		copy(buf[:], s[i:soff])

		fmt.Printf("%2d: %q, codepoint: %#6x, encoded bytes: %#v\n", i, r, r, buf[:rl])
	}
}

// OUTPUT:
//  0: 'G', codepoint:   0x47, encoded bytes: []byte{0x47}
//  1: '♡', codepoint: 0x2661, encoded bytes: []byte{0xe2, 0x99, 0xa1}
//  4: '!', codepoint:   0x21, encoded bytes: []byte{0x21}
```
