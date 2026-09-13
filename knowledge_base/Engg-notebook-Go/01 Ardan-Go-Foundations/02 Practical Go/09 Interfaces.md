#engineers-notebook #golang #go-interface #comma-ok #dot-type

### Empty interface `interface{}` or `any`
```go
var a any
a := "hi"
i := a.(int) // this will panic

// comman-ok will not panic
j, ok := a.(int)
if ok {} else {}
```

### `.(type)` operator
```go
switch a.(type) {
	case int:
		fmt.Printf("%d", a)
	case string:
		fmt.Printf("%s", a)
}
```