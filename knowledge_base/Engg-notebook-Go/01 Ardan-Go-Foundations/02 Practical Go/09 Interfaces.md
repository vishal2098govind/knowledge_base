#engineers-notebook #golang #go-interface #comma-ok

```go
var a any
a := "hi"
i := a.(int) // this will panic

// comman-ok will not panic
j, ok := a.(int)
if ok {} else {}
```