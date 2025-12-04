#go #maps

- Zero value of a map is nil
- the line : `var users map[string]user` sets `users` to it's zero value i.e. nil
```go
package main

import "fmt"

type user struct {
	FirstName string
	LastName  string
}

func main() {
	var users map[string]user
	users["Ford"] = user{"Henry", "Ford"}
	users["Mouse"] = user{"Micky", "Mouse"}
	users["Jackson"] = user{"Michael", "Jackson"}

	fmt.Printf("%+v\n", users)
}

// OUTPUT:
// panic: assignment to entry in nil map
//
// goroutine 1 [running]:
// main.main()
//        /Users/govind/Dev/personel/go-ardan-labs/practical-go/maps/main.go:12 +0x2c
// exit status 2
```
- Maps must be initialised using either 
	- `make` function: `users := make(map[string]user)`
	- or empty literal construction: `users := map[string]user{}` 
- Since we avoid creating slices using empty literal construction, for consistency for maps, prefer using `make` function over empty literal construction

## Difference among initialisations of reference types (slices, maps, channels)
#go-make #go-reference-types

```
type user struct {FirstName string; LastName string}
var usernames map[string]user
usernames := make(map[string]user)
usernames := map[string]user{}

var users []user
users := make([]user, len, cap)
users := []user{}
```
### 1. Maps

#### A. Declaration — no memory allocated

`var usernames map[string]user`

- This **declares** a variable.
- `usernames` is **nil**.
- You **cannot write** into this map yet.
- Len = 0, but also **no backing hash table exists**.

If you try:
`usernames["vishal"] = user{}`
You get:
`panic: assignment to entry in nil map`

---

#### B. Using `make` — allocates an actual map

`usernames := make(map[string]user)`
- Memory allocated
- Backing hash table created
- Safe to insert key/value pairs
- len = 0 initially (empty)

This is the **_most common and correct_ way**.

---

#### C. Using map literal

`usernames := map[string]user{}`

Equivalent to `make(map[string]user)`  
except that you also can pre-initialize:

```go
usernames := map[string]user{
	"vishal": {FirstName: "Vishal", LastName: "Govind"},
}
```

---

### Summary for maps

|Syntax|Allocated?|Insert Safe?|Use Case|
|---|---|---|---|
|`var m map[K]V`|❌ No (nil map)|❌ No|Declare only (assign later)|
|`m := make(map[K]V)`|✅ Yes|✅ Yes|Standard construction|
|`m := map[K]V{}`|✅ Yes|✅ Yes|Standard, literal syntax|

---

### 2. Slices

#### A. Declaration — nil slice
`var users []user`

- Declares a slice variable
- `users == nil`
- Len = 0, Cap = 0
- You **can append** to a nil slice → Go will allocate as needed:
`users = append(users, user{})`
Works fine.

---

#### B. Using make

`users := make([]user, len, cap)`

If you do:
`users := make([]user, 5, 10)`
Then:
- `len` = 5 → slice already contains 5 **zero-value user structs**
- `cap` = 10 → room for 5 more before resize

Memory layout created immediately.

You can mutate the first 5 elements directly:

`users[0].FirstName = "Vishal"`

---

#### C. Slice literal

`users := []user{}`

Equivalent to:
`users := make([]user, 0)`
This is:
- Not nil (length 0, cap 0 or small)
- You can append, Go grows it.

---

#### Summary for slices

|Syntax|Nil?|Len|Cap|Append Safe?|Use Case|
|---|---|---|---|---|---|
|`var s []T`|✔️ nil|0|0|✔️ append ok|Declare only|
|`s := make([]T, n, m)`|❌ not nil|n|m|✔️ append ok|Preallocate space|
|`s := []T{}`|❌ not nil|0|~0–small|✔️ append ok|Literal empty slice|

---

### Short, Practical Mental Model

#### Maps
- `var m map[...]...` → **nil**, no table → can't write
- `make` or `{}` → **ready to use**, can write

#### Slices
- `var s []T` → **nil slice**, but append is fine
- `s := make([]T, n)` → pre-filled slice
- `s := []T{}` → empty but ready

---

### Example of all together

```go
type user struct {
    FirstName string
    LastName  string
}

// MAPS
var m1 map[string]user       // nil
m2 := make(map[string]user)  // empty but allocated
m3 := map[string]user{}      // empty but allocated

// SLICES
var s1 []user                // nil slice
s2 := make([]user, 5, 10)    // len 5, cap 10
s3 := []user{}               // empty slice
```


## Why is the `make` function way is more standard and correct over using empty literal

Short answer: **both are correct**, and **there is no functional difference** between:
`usernames := make(map[string]user)`
and
```go
usernames := map[string]user{}
```

So why do Go developers say **the `make` version is “more standard”**?
It’s not a rule — it’s a **convention** rooted in Go’s design philosophy.

### Reason 1 — Consistency with slices: `make` is the constructor
In Go, `make` is _the_ way you construct reference types:
- map
- slice
- channel
Example:
```go
s := make([]int, 0) 
ch := make(chan int) 
m := make(map[string]user)
```
Using `make` gives readers an immediate signal:
> “This creates an initialised, ready-to-use reference type.”

This improves readability when skimming code.

### Reason 2 — Literal `{}` is more commonly used when initial values exist
Example:
```go
usernames := map[string]user{
     "vishal": {FirstName: "Vishal", LastName: "Govind"},
}
```

So this version:
```go
usernames := map[string]user{}
```

looks like “initialising with no elements” rather than semantically creating a map.
That’s why some prefer `make` — it matches intent more clearly.

### Reason 3 — Some codebases enforce *stylistic consistency*
Large Go codebases (Google, Uber, etc.) often choose one style so that:
- code is uniform
- linters are simple
- diffs don’t randomly switch styles

Many linters prefer:
```go
make(map[string]user)
make([]T, 0)
make(chan T)
```
Because it prevents style fragmentation.
### Reason 4 — `make` allows optional capacity hints (future-proofing)
Even though you aren’t using it here, `make` has this capability:
```go
make(map[string]user, 100)
```

This pre-allocates space for ~100 map entries → reducing rehash operations.
The literal cannot do this.
So some developers always use `make` so they can easily add a capacity later.

### Reason 5 — Historic Go examples/documentation used `make`
Early Go examples and standard library code used:
```go
m := make(map[string]int)
```
So many devs simply followed that style.

## Basic Initialise, Write, Read and Delete
#value-semantics 
- `users := make(map[string]user)` := here, the key is a `string` and value is `user` which uses value-semantics
- `map` always **stores it's own copy** of whatever that data, and also **returns a copy** of the original data
- `mouse := users["Mouse"]` here, mouse is a copy of the original `user` struct in memory
Example:
```go
package main

import "fmt"

// user represents someone using the program.
type user struct {
	name    string
	surname string
}

func main() {

	// Declare and make a map that stores values
	// of type user with a key of type string.
	users := make(map[string]user)

	// Add key/value pairs to the map.
	users["Roy"] = user{"Rob", "Roy"}
	users["Ford"] = user{"Henry", "Ford"}
	users["Mouse"] = user{"Mickey", "Mouse"}
	users["Jackson"] = user{"Michael", "Jackson"}

	// Read the value at a specific key.
	mouse := users["Mouse"]

	fmt.Printf("%+v\n", mouse)

	// Replace the value at the Mouse key.
	users["Mouse"] = user{"Jerry", "Mouse"}

	// Read the Mouse key again.
	fmt.Printf("%+v\n", users["Mouse"])

	// Delete the value at a specific key.
	delete(users, "Roy")

	// Check the length of the map. There are only 3 elements.
	fmt.Println(len(users))

	// It is safe to delete an absent key.
	delete(users, "Roy")

	fmt.Println("Goodbye.")
}
```

## Reading an absent key (comma-ok)
```go
package main

import "fmt"

func main() {

	// Create a map to track scores for players in a game.
	scores := make(map[string]int)

	// Read the element at key "anna". It is absent so we get
	// the zero-value for this map's value type.
	score := scores["anna"]

	fmt.Println("Score:", score)

	// If we need to check for the presence of a key we use
	// a 2 variable assignment. The 2nd variable is a bool.
	score, ok := scores["anna"]

	fmt.Println("Score:", score, "Present:", ok)

	// We can leverage the zero-value behavior to write
	// convenient code like this:
	scores["anna"]++

	// Without this behavior we would have to code in a
	// defensive way like this:
	if n, ok := scores["anna"]; ok {
		scores["anna"] = n + 1
	} else {
		scores["anna"] = 1
	}

	score, ok = scores["anna"]
	fmt.Println("Score:", score, "Present:", ok)
}
```

## Only types that can have equality defined on them can be a map key.
```go
package main

import "fmt"

// user represents someone using the program.
type user struct {
	name    string
	surname string
}

// users defines a set of users.
type users []user

func main() {

	// Declare and make a map that uses a slice as the key.
	u := make(map[users]int)

	// ./example3.go:22: invalid map key type users

	// Iterate over the map.
	for key, value := range u {
		fmt.Println(key, value)
	}
}
```

## Declare, initialise and iterate over a map
- iterating over a map is random
```go
package main

import "fmt"

// user represents someone using the program.
type user struct {
	name    string
	surname string
}

func main() {

	// Declare and initialize the map with values.
	users := map[string]user{
		"Roy":     {"Rob", "Roy"},
		"Ford":    {"Henry", "Ford"},
		"Mouse":   {"Mickey", "Mouse"},
		"Jackson": {"Michael", "Jackson"},
	}

	// Iterate over the map printing each key and value.
	for key, value := range users {
		fmt.Println(key, value)
	}

	fmt.Println()

	// Iterate over the map printing just the keys.
	// Notice the results are different.
	for key := range users {
		fmt.Println(key)
	}
}
```
- Iterating over a alphabetically sorted keys
```go
package main

import (
	"fmt"
	"sort"
)

// user represents someone using the program.
type user struct {
	name    string
	surname string
}

func main() {

	// Declare and initialize the map with values.
	users := map[string]user{
		"Roy":     {"Rob", "Roy"},
		"Ford":    {"Henry", "Ford"},
		"Mouse":   {"Mickey", "Mouse"},
		"Jackson": {"Michael", "Jackson"},
	}

	// Pull the keys from the map.
	var keys []string
	for key := range users {
		keys = append(keys, key)
	}

	// Sort the keys alphabetically.
	sort.Strings(keys)

	// Walk through the keys and pull each value from the map.
	for _, key := range keys {
		fmt.Println(key, users[key])
	}
}
```

## Cannot take the address of an element in a map
```go
package main

// player represents someone playing our game.
type player struct {
	name  string
	score int
}

func main() {

	// Declare a map with initial values using a map literal.
	players := map[string]player{
		"anna":  {"Anna", 42},
		"jacob": {"Jacob", 21},
	}

	// Trying to take the address of a map element fails.
	anna := &players["anna"]
	anna.score++

	// ./example4.go:23:10: cannot take the address of players["anna"]

	// Instead take the element, modify it, and put it back.
	player := players["anna"]
	player.score++
	players["anna"] = player
}
```

## Maps are reference types
#go-reference-types #go-reference-types-dual-behaviour
- reference types in go have **dual behaviour**.
	- when copying a reference type (slice or map or channels or interfaces or functions) while moving around the program, we use **value semantics**
		- e.g. while passing to a function : `double(scores, "anna")`
	- but while reading and writing to these reference types, we use **pointer semantics**
		- e.g. `scores[player] = scores[player] * 2` -> doing both read and write here
```go
package main

import "fmt"

func main() {

	// Initialize a map with values.
	scores := map[string]int{
		"anna":  21,
		"jacob": 12,
	}

	// Pass the map to a function to perform some mutation.
	double(scores, "anna")

	// See the change is visible in our map.
	fmt.Println("Score:", scores["anna"])
}

// double finds the score for a specific player and
// multiplies it by 2.
func double(scores map[string]int, player string) {
	scores[player] = scores[player] * 2
}
```
- Thus, accessing a map is not safe in a multi-threaded environment, unless the map is synchronised and orchestrated