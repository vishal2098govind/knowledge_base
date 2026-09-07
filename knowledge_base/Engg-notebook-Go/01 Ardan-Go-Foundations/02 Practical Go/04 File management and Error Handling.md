#engineers-notebook #golang #file #defer #wrapping-errors 

Defer happens when function exits, **no matter what (including panic)**
If multiple defers within a function, they are executed in **reverse order** (stack, LIFO)

Idiom: try to acquire a resource, check for error, defer release

```go
err := Kill("server.pid")
if err != nil {
	fmt.Println("ERROR", err)
	if errors.Is(err, fs.ErrNotExist) {
		fmt.Println("not found")
	}
	// chaining wrapped err
	for e := err; e != nil; e = errors.Unwrap(err) {
		fmt.Println("error: ", err)
	}
}
```

```sh
➜  kill-server git:(vishal-kb) ✗ go run .
ERROR "server.pid" - bad pid, expected integer
error:  "server.pid" - bad pid, expected integer
error:  expected integer
```