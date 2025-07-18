#go #go-template #html/template #contextual-encoding

- Step 1: Parse
- Step 2: Execute
### Example
```go
// ./hello.gohtml

<h1>Hello, {{.Name}}</h1>
{{.Age}}
{{.Meta.Visits}}
{{.Bio}}
```

```go
// ./main.go
package main

import (
	"html/template"
	"os"
)

type User struct {
	Name string
	Age int
	Meta UserMeta
}
type UserMeta struct {
	Visits int
}

func main() {
	t, err := template.ParseFiles("hello.gohtml")
	if err != nil {
		panic(err)
	}
	user := User{ Name: "Vishal", Age: 26, Meta: UserMeta{ Visits: 10 } }
	err = t.Execute(os.Stdout, user)
	if err != nil {
		panic(err)
	}
}

// ---------------------------
// Output:
// ---------------------------
// $ go run .
// <h1>Hello, Vishal</h1>
// 26
// 10
```

#### XSS
The html templating package also takes care of XSS #xss 
```go
// ./main.go
package main

import (
	"html/template"
	"os"
)

type User struct {
	Name string
	Age int
	Meta UserMeta
}
type UserMeta struct {
	Visits int
}

func main() {
	t, err := template.ParseFiles("hello.gohtml")
	if err != nil {
		panic(err)
	}
	user := User{ Name: "<script>alert("Haha, you have been hacked!");</script>", Age: 26, Meta: UserMeta{ Visits: 10 } }
	err = t.Execute(os.Stdout, user)
	if err != nil {
		panic(err)
	}
}

// ---------------------------
// Output:
// ---------------------------
// $ go run .
// <h1>Hello, &lt;script&gt;alert(&#34;Haha, you have been hacked!&#34;);&lt;/script&gt;</h1>%
// 26
// 10
```
- NOTE: XSS is not supported in the `text/template` package and only with `html/template` package

### What if we want to render actual HTML while using `html/template`
- Use `template.HTML` type
```go

// ./main.go
package main

import (
	"html/template"
	"os"
)

type User struct {
	Name string
	Bio template.HTML
	Age int
	Meta UserMeta
}
type UserMeta struct {
	Visits int
}

func main() {
	t, err := template.ParseFiles("hello.gohtml")
	if err != nil {
		panic(err)
	}
	user := User{ Name: "Vishal", Bio: "<script>alert("Haha, you have been hacked!");</script>" }
	err = t.Execute(os.Stdout, user)
	if err != nil {
		panic(err)
	}
}

// ---------------------------
// Output:
// ---------------------------
// $ go run .
// <h1>Hello, Vishal</h1>
// <script>alert("Haha, you have been hacked!");</script>
// 0
// 0
```

### Contextual Encoding by `html/template`
- HTML and JS encoding are done bit differently, thus encoding done is **contextual**  
```gohtml
<h1>Hello, {{.Name}}</h1>
{{.Bio}}

<script>
const user = {
	"name": {{.Name}},
	"bio": {{.Bio}},
}
console.log(user)
</script>
```

```go
// ./main.go
package main

import (
	"html/template"
	"os"
)

type User struct {
	Name string
	Bio string
}

func main() {
	t, err := template.ParseFiles("hello.gohtml")
	if err != nil {
		panic(err)
	}
	user := User{ Name: "Vishal", Bio: "<script>alert("Haha, you have been hacked!");</script>" }
	err = t.Execute(os.Stdout, user)
	if err != nil {
		panic(err)
	}
}
// ---------------------------
// Output: note that the js encoding within the <script> is different than what is in html encoding of the same string
// ---------------------------
// $ go run .
// <h1>Hello, Vishal</h1>
// &lt;script&gt;alert(&#34;Haha, you have been hacked!&#34;);&lt;/script&gt;

// <script>
// const user = {
//    "name": "Vishal",
//    "bio": "\u003cscript\u003ealert(\"Haha, you have been hacked!\");\u003c/script\u003e",
//    "age": 10, // notice age is printed without quotes (i.e. it is not "10") i.e. age is int 
// }
// console.log(user)  
// </script>
```
- ![[Contextual Encoding.png]]

### Using template
- `template.ParseFiles` returns error if there is anything wrong with the html files, like invoking undefined functions
- ideally we want to know about any template parsing errors even before starting up our server
- `template.Execute` returns an error when it is unable to access some of the referred fields in the template
	- like if the template refers to `{{.Name}}`, and the `data` passed to `template.Execute(w, data)` doesn't have `Name` field
	- but by that time, the template might have written some stuff to the writer `w` and while it encounters error, it pauses/stops further writing and returns error. Thus, `w` would have some of the stuff which executed correctly until the first error encountered
```go
// ./templates/home.gohtml

<html>
	<h1>Welcome To Golang</h1>
	<a href="/contacts">Contacts</a>
</html>
{{.Name}}
```

```go
// ./main.go
func homeHandler(w http.ResponseWriter, _ *http.Request) {
	path := filepath.Join("templates", "home.gohtml")
	t, err := template.ParseFiles(path)
	if err != nil {
		log.Printf("parsing template: %v", err)
		http.Error(w, "There was an error parsing template", http.StatusInternalServerError)
		return
	}
	
	err = t.Execute(w, "a string")
	if err != nil {
		http.Error(w, "There was an error executing template", http.StatusInternalServerError)
		return
	}
}
```
- here, the `w` would contain (both error message and correctly executed text so far) 
	- also the status code of this response would also not be changed to `http.StatusInternalServerError`
		- this is because, **once we start writing to a `http.ResponseWriter`, it sets the status code as `http.StatusOK` (200)**
		- and **status code can't be changed once it's been set**
```
# Welcome To Golang

Contacts There was an error executing template
```
##### How to get around this issue?
#io-Writer #strings-Builder #bytes-Buffer
- `bytes.Buffer` and `strings.Builder` implement `io.Writer`
- One way would be execute the entire template to a `bytes.Buffer` or `*strings.Builder` instead of directly to `http.ResponseWriter`, then if it executes successfully, then write it to `http.ResponseWriter`
```go
var buf = bytes.NewBufferString("")
// or var buf = &strings.Builder{}
err = t.Execute(buf, "name")
if err != nil {
	http.Error(w, "There was an error executing template", http.StatusInternalServerError)
	return
} 
fmt.Fprint(w, buf.String())
```

### Named templates
#named-templates
The `html/template` package allows us to create named templates that can be reused in other templates. For instance, let’s imagine we wanted to design our website but we need some filler text to use in the design process. We could start by creating a template with filler text.

```html
<h1>Welcome to my awesome site!</h1>
{{template "lorem-ipsum"}} 
{{template "lorem-ipsum"}} 
{{template "lorem-ipsum"}}

{{define "lorem-ipsum"}}
<p>
  Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor
  incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis
  nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat.
  Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore
  eu fugiat nulla pariatur. Excepteur sint occaecat cupidatat non proident, sunt
  in culpa qui officia deserunt mollit anim id est laborum.
</p>
{{end}}
```

The `{{define "..."}}` code starts a named template block. From this point onwards, anything in the template file until we reach `{{end}}` will be included as part of the named template. The text inside the quotation marks is the name of the template. In our named template, the name is `lorem-ipsum`.
![[Named template.png]]
### Variadic parameters for files in `template.Parse` and `template.ParseFS`
#variadic-paramters

```go
// within html/template package

package template

// ...

func ParseFiles(filenames ...string) (*Template, error) {
	return parseFiles(nil, readFileOS, filenames...)
}

func ParseFS(fs fs.FS, patterns ...string) (*Template, error) {
	return parseFS(nil, fs, patterns)
}
```
- The variadic parameter indicates that all the file names passed in `patterns` or `filenames` are going to be used while creating and returning the final template
```go
tpl, err := template.ParseFS(fs, "home.gohtml", "layout-parts.gohtml")
if err != nil {
	return nil, fmt.Errorf("parsing template: %w", err)
}
```
- The **sequence** of `filenames` or `patterns` **matters**. The files are rendered sequentially.

### Custom Template Functions
`template.Funcs` is used to add in a map of custom functions into a template
```go
func ParseFS(fs fs.FS, patterns ...string) (*Template, error) {
	tpl, err := template.ParseFS(fs, patterns...)
	if err != nil {
		return nil, fmt.Errorf("parsing template: %w", err)
	}
	
	return &Template{htmlTpl: tpl}, nil
}
```