#go #http #cookies #net/http #csrf

#### Set Cookie into response
We can set cookies using `http.SetCookie(w ResponseWriter, cookie *Cookie)` which adds a `Set-Cookie` header to the provided `ResponseWriter's` header
```go
cookie := http.Cookie{
	Name: "cookie key",
	Value: "cookie's value",
	Path: "/" // dictates which paths on our server have access to this cookie
	// Path: "/dashboard" will let any paths starting with /dashboard to access this cookie
}
http.SetCookie(w, &cookie)
```
- Cookies have several use cases. We can use it to set user authentication token whenever a user signs in for subsequent protected requests to be successful
 
#### Reading cookie from request:
We can use `http.Request.Cookie(name string) (*http.Cookie, error)` or `http.Request.Cookies() []*http.Cookie`
```go
emailC, err := r.Cookie("email")
if err != nil {
	fmt.Fprint(w, "The email cookie could not be read")
}
fmt.Fprintf(w, "Email cookie: %s\n", emailC.Value)
fmt.Fprintf(w, "Headers: %+v\n", r.Header)
```
Output:
```md
Email cookie: bob.govind@gmail.com
Headers: map[Accept:[text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7] Accept-Encoding:[gzip, deflate, br, zstd] Accept-Language:[en-GB,en-US;q=0.9,en;q=0.8] Cache-Control:[no-cache] Connection:[keep-alive] **Cookie:[email=bob.govind@gmail.com]** Pragma:[no-cache] Sec-Ch-Ua:["Not)A;Brand";v="8", "Chromium";v="138", "Google Chrome";v="138"] Sec-Ch-Ua-Mobile:[?0] Sec-Ch-Ua-Platform:["macOS"] Sec-Fetch-Dest:[document] Sec-Fetch-Mode:[navigate] Sec-Fetch-Site:[none] Sec-Fetch-User:[?1] Upgrade-Insecure-Requests:[1] User-Agent:[Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/138.0.0.0 Safari/537.36]]
```

#### Securing cookies from XSS
#xss 
Cross site scripting were used to execute some arbitrary javascript that any attacker decides to write
The way it's most frequently being used is that inside of a `<script>` tag, they would instead try to get our website to load javascript from their server.
This becomes dangerous whenever we have things like cookies and having them leaked could lead to user credentials and other things being stolen.
`http.Template` ensures to escape or encode user inputs to avoid `XSS`
Cookies are by default accessible to JS
```js
> document.cookie;
'email=bob.govind@gmail.com'
```
We can make it more secured by limiting the access of Cookies by JS so that JS cannot access our cookies
```go
cookie := http.Cookie{
	Name: "cookie key",
	Value: "cookie's value",
	Path: "/",
	HttpOnly: true, // this disables access of cookie by JS
}
http.SetCookie(w, &cookie)
```

#### CSRF Attack
**C**ross-**s**ite **r**equest **f**orgery
The browser will include cookies for a website even when we click the link to that website and that link is present in another website
This opens the door for a CSRF attack especially if any of the actions use a GET request instead of a POST/PUT/DELETE method in the API

##### CSRF Middleware
**Why using a library vs writing our own CSRF middleware?**
- CSRF libs tend to work out of the box
- Minimal footprint, so easy to swap out later

Using `gorilla/csrf` package
- `csrf.Protect` middleware provides CSRF protection on routes attached to a router or a sub-router
	- sets `csrf-token` into cookie
	- make sure any POST/PUT requests that are coming into our server have a valid `csrf-token` associated with them
- we need to make sure our forms have the `csrf-token` embedded into them which is obtained by 
	- `csrf.Token()` 
	- `csrf.TemplateField` gives snippet of HTML that is meant to be stuck inside of a form which is already a hidden input field

### Cookie tampering - Session Tokens
#session-token

**Session tokens** prevent cookie tampering through obfuscation
Each user after login gets assigned a random string which is added to the cookie
This cookie is also added to as a mapping in the database table corresponding to each user

```
user_id | session token
==================================
	 1  | aogeqobgqwegwWQOfea
	 2  | oiaenfoq1jnaege
```

Cookies can be changed, but attackers won't be able to predict or easily guess a valid session token if they are generated with **sufficient randomness**

#### `crypto/rand.Read()` 
#crypto/rand 
Using `rand.Read()` to get random bytes of given length
```go
package main

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
)

func main() {
	n := 8
	b := make([]byte, n)
	fmt.Println(b)
	nRead, err := rand.Read(b)
	if err != nil {
		panic(err)
	}
	if nRead < n {
		panic("didn't read enough bytes")
	}
	fmt.Println(b)

	// encoding []byte to string
	fmt.Println(base64.URLEncoding.EncodeToString(b))
}

// OUTPUT:
// [0, 0, 0, 0, 0, 0, 0, 0]
// [100, 210, 132, 2, 92, 13, 49, 45] - random bytes
// p05Ty-jk4mw=
```

#### `math/rand`
```go
package main

import (
	"fmt"
	"math/rand"
)

func main() {
	s := make([]int, 100)
	fmt.Println(rand.Intn(len(s)))
	fmt.Println(rand.Intn(len(s)))
	fmt.Println(rand.Intn(len(s)))
}

// OUTPUT Run-1:
// 81
// 87
// 47

// OUTPUT Run-2:
// 81
// 87
// 47

// It would give same output on every run
```

- `math/rand` has a `Seed` which is a number we can pass in that kind of starts off the random number (pseudo random number) generator
- Given a specific seed, the numbers that it gets from that point on will all be predictable or same
- So to get different outputs, we can change the seed on every run
```go
package main

import (
	"fmt"
	"math/rand"
	"time"
)

func main() {
	s := make([]int, 100)
	rand.Seed(time.Now().UnixNano())
	fmt.Println(rand.Intn(len(s)))
	fmt.Println(rand.Intn(len(s)))
	fmt.Println(rand.Intn(len(s)))
}

// OUTPUT Run-1:
// 81
// 87
// 47

// OUTPUT Run-2:
// 43
// 17
// 42
```
Using `math/rand` is not a good fit for generating random session tokens because if an attacker would somehow figure out the seed, they would then have a way of reliably predicting any future session tokens.

#### Storing Session in Database
Store session token **hashes** instead of actual raw session tokens to prevent access in case of database leaks
Session token hashes per user needs to be **unique** because if a hash is given to multiple users, we cannot predict which user
```sql
CREATE TABLE IF NOT EXISTS
	sessions (
		id SERIAL PRIMARY KEY,
		user_id INT UNIQUE,
		token_hash TEXT UNIQUE NOT NULL -- a hash not actual token, similar reasons as password
); 
```

#### Deleting cookies
To override a cookie, we can just set a new cookie with that name and send the response back. The browser would know to override that cookie.
To delete a cookie, simply set the `MaxAge` of that cookie to -1
```go
func deleteCookie(w http.ResponseWriter, name string) {
	cookie := newCookie(name, "")
	cookie.MaxAge = -1
	http.SetCookie(w, cookie)
}
```