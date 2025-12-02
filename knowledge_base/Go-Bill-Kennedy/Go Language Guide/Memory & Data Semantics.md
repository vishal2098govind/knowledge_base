#go #memory-semantics #data-semantics

## Variables
- `int` vs `int8`/`int16`/ `int32`/ `int64`
	- `int` let's the compiler choose the size of `int` used by the underlying hardware architecture
		- for a 32-bit architecture int32 is used
		- for a 64-bit architecture int64 is used

### Conversion over casting
- Conversion allocates fresh memory for the new value, casting uses existing allocated memory of the value
- Conversion is much safer than casting
- Go has type conversion, no type casting
```go
aa := 10 // int
aaa := int32(aa)
```