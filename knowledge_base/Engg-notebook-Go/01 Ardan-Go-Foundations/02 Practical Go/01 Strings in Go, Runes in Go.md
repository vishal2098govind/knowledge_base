```sh
➜  hw git:(vishal-kb) ✗ cd ../banner 
➜  banner git:(vishal-kb) ✗ go run banner.go
  Go
------
  G@
------
➜  banner git:(vishal-kb) ✗ 
➜  banner git:(vishal-kb) ✗ go run banner.go
  Go
------
G♡
------
➜  banner git:(vishal-kb) ✗ 
➜  banner git:(vishal-kb) ✗ 
```

## How text is represented in computers
- **ASCII table**
	- In beginning, computers were mainly developed in english speaking countries, when a **single byte** was used to represent all of the characters via encoding.
- **UTF-8**
	- Then other nations came in and wanted to use computers as well and support other languages as well. Then, **one byte was not enough** for encoding
	- There was needed some kind of **database** of giving a number to every possible character in every human language
	- Unicode is a big database which starts from ASCII as base.
	- The database is kept updated every couple of years
	- There are several **encoding schemes** to represent unicode character to sequence of bytes. Most popular one is **UTF-8**
	- UTF-8 is a **variable length encoding scheme**
		- variable length meaning, a character from UTF-8 encoding can be from a single byte and **up to 4 bytes**
	- Refer [[unicode.pdf]]
	- in UTF-8, all of ASCII characters still use single byte, but non ASCII may take more than one byte
	- For example `<<` this unicode character **takes 2 bytes** and there's also something called as **control character** in UTF-8
		- thus, `<<` has UTF-8 encoding of `C2AB` for example, this entire single encoding is called **code-point** in UTF-8 or unicode
		- in Go it's called a **rune**

```go
fmt.Println("len:", len("Go")) // 2
fmt.Println("len:", len("G♡")) // 4
```
- Go is built for performance
- A go's string is a **struct** which has 2 fields
	- length
	- pointer to underlying bytes
- when we ask for length of a string, it refers to the length field and returns without having to peek memory every time and calculate number of **code-points** or **runes** since they can be of any number of bytes given go uses UTF-8 encoding
- Similarly concept is used while indexing a string and looping over a string
```go
fmt.Println("s[1]:", "G♡"[1]) // 226
fmt.Printf("s[1]:%c\n", "G♡"[1]) // s[1]:â

// looping over string iterates over runes (code-points) in a string
for i, c := range "G♡" {
	fmt.Printf("%c at %d\n", c, i)
}
// G at 0
// ♡ at 1
```
- length and indexing use bytes
- for looping using range uses runes
- byte = uint8 (1 byte)
- rune = int32 (4 bytes)
```go
type byte = uint8 // a line from go/src/builtin/builtin.go
type rune = int32 // a line from go/src/builtin/builtin.go
```
- to get rune count in a string, can use `utf8.RuneCount([]byte(str))` instead of `len(str)`
```go
text := "G♡"
utf8.RuneCount([]byte(text)) // 2
len(text) // 4
```