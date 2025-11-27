#go #strings #rune #utf-8 #byte #uint8 #int32

- Strings in Go are UTF-8 encoded
	- In UTF-8, all the ASCII characters take a single byte, remaining characters can take at max **4 bytes**
	- Thus UTF-8 characters are dealt with runes in Go. each character can fit in a rune since a rune is of 4 bytes (`int32`)
- `len()` and `s[]`  deal with each **byte**.
	- a byte in Go is `uint8`
- `for i, c range s {}` loops over each **rune**
	- a **rune** in Go is `int32` i.e. of **4 bytes**

In Go, a `rune` is an alias for the `int32` data type, meaning it always occupies ==**4 bytes**== (32 bits) in memory. It is used to represent a single Unicode code point. 

However, when a rune is part of a string encoded in UTF-8 (which is the default in Go), the number of bytes it occupies within that string can vary: 

- **1 byte** for ASCII characters (U+0000 to U+007F).
- **2 bytes** for certain characters like some accented letters.
- **3 bytes** for many other symbols and characters, like the Euro sign (€) or some Chinese characters.
- **4 bytes** for a wide range of characters, including emojis.

To work with characters in a way that respects the variable-length UTF-8 encoding of a string, we should iterate over the string using a `for range` loop, which automatically decodes each rune and provides its byte index and value. 

We can also use functions from the standard `unicode/utf8` package for explicit handling of runes and their byte lengths within a string.

```go
package main

import (
	"fmt"
	"unicode/utf8"
)

func main() {
	s := "G♡!"
	fmt.Println("len:", len(s))
	fmt.Println("s[1]:", s[1])
	fmt.Printf("s[1]: %c\n", s[1])

	for i, c := range s {
		fmt.Printf("%c at %d\n", c, i)
	}

	fmt.Printf("# runes or # characters in s: %d\n", utf8.RuneCountInString(s))
}

/* Output:
len: 5
s[1]: 226
s[1]: â
G at 0 => G is the 1st rune (0th index)
♡ at 1 => ♡ is the second rune (2st index) and spans over 3 bytes (1st, 2nd and 3rd index) - each index represents a byte
! at 4 => thus, ♡ takes 3 bytes
# runes or # characters in s: 3
*/
```