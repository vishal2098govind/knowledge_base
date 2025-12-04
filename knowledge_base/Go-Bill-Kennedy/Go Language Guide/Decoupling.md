#go #decoupling

## `IP` and `IPMask` Example
- In the `net` package, there are two custom types defined `IP` and `IPMask`
```go
type IP []byte
type IPMask []byte
```
- We can also say that `IP` and `IPMask` types are **based on** a built-in type `[]byte`
- So, now since they are basing on a slice of bytes, and since slice is a reference type in itself, thus `IP` and `IPMask` are also **reference types**, which also thus means that `IP` and `IPMask` now also have #go-reference-types-dual-behaviour 
	- i.e. now we can use **value-semantics** to move them around 
	- and also use **pointer-semantics** to read and write

```go
// Mask is using a value reciever of type IP and returning a value of IP. 
// This method is using value semantics for type IP
func (ip IP) Mask(mask IPMask) IP {
	// method body or logic ...
}

```

Once we know the data inputs and input data types and data outputs and output data types, we should never really be confused about how to design an API. The data itself is defining the design. 

So, here, since the based type for `IP` and `IPMask` is a slice which uses **value semantics** while moving data around, and **pointer semantics** while reading and writing, the methods on this type *should* **(there are exceptions anyways)** also using **value semantics** while accepting such data and returning such data (i.e. while moving around such data)

## `Time` Example
Consider a struct type `Time` from the `time` package
```go
type Time struct {
	sec int64
	nsec int32
	loc *Location
}
```
Sometimes it's obvious about which semantics to use while moving around data of a certain type. But sometimes it's not. The decision of semantics is not made based on performance or possibility of being able to use a certain semantics for certain type, but rather the decision is made based on 
- correctness
- what is the reasonable thing to do for this type of data
- what is the expectation
- what are people going to assume when they are going to use the API built around this type or look at data of such type

Usually the factory functions of a type (`New()`) dictate the semantics that will be used by the methods of that type
```go
// Here the factory returns a value, this means that we should be using value semantics while passing the objects of type Time around
func New() Time {
	sec, nsec := now()
	return Time{sec + unixToInternal, nsec, Local}
}
```
- Value semantics help in not over polluting the heap, unlike the pointer semantics, and thus reducing GC overhead
```go
// uses value semantics
func (t Time) Add(d Duration) Time {
	// code here
}

func div(t Time, d Duration) (qmod2 int, r Duration) {
	// code here
}
```
- In case of exceptions like 
	- where we have to deal with pointers for representing nil values 
	- or like where we need to unmarshal/marshal
	- in such cases where we have to change semantics from value to pointer, of a certain type, it would be done in a **very small scope** and would be normally around ideas of decoding and encoding
```go
// The only places where pointer semantics is being used for `Time` api are these unmarshal related functions
func (t *Time) UnmarshalBinary(data []byte) error {/* code here */}
func (t *Time) GobDecode(data []byte) error {/*code here*/}
func (t *Time) UnmarshalJSON(data []byte) error {/* code here */}
func (t *Time) UnmarshalTEXT(data []byte) error {/* code here */}
```

- We should never **change** from **pointer semantics to value semantics**. **It is never ok**. Buggy.
- We can change from value to pointer semantics within a small scope while dealing with nils or unmarshalling. That's ok only the change is in such small ok, not beyond.

## `os.Open` Example
```go
func Open(name string) (file *File, err error) {
	return OpenFile(name, O_RDONLY, 0)
}
```
Here, this a factory function for file returns a pointer to File. 
Thus, factory function dictates the semantics being used
and thus, **pointer semantics** is being used for File and so should we and should also should **share** File
Thus, regardless of the operation on File, they should **always be shared and never copied**

```go
func (f *File) Chdir() error {
	// code here
}
```
`Chdir` is using pointer semantics on the receiver. This is because, we must be sharing file value and never copy, even if the method doesn't mutate anything on the underlying type 