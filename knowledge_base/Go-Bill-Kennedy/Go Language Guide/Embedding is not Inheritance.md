#go #go-embed

**Struct embedding in Go is _not inheritance_.**

# 1. What Embedding _Actually_ Is

Struct embedding is simply **a form of field inclusion + promotion**.

When you write:
```go
type A struct {
    X int
}

type B struct {
    A     // embedding
    Y int
}
```

Go copies the fields of `A` **into** `B` (conceptually), and also gives you automatic **field/method promotion**, meaning:
```go
var b B
b.X = 10   // allowed! (promoted field)
b.Y = 20
```

It _feels_ like inheritance, but **nothing about the type system thinks B is an A**.

# 2. Why Embedding Is Not Inheritance

In inheritance:
- A subclass _is-a_ parent type
- You can upcast: subclass → parent
- Virtual methods allow overrides
- You have polymorphism by default

In Go embedding:

- **B is _not_ an A**
- You **cannot upcast** a B to an A
- You **cannot override methods**
- There is **no polymorphism** unless you use interfaces
- Embedding is purely **composition**

Try this:
```go
var b B
var a A

a = b   // ❌ DOES NOT COMPILE
```

This single fact destroys any idea of inheritance.

# 3. Method "Overriding" Does NOT Exist

If both A and B define a method called `Foo()`, B does NOT override A.

It simply shadows the promoted method.

```go
func (A) Foo() { fmt.Println("A Foo") }
func (B) Foo() { fmt.Println("B Foo") }

var b B
b.Foo()   // prints "B Foo"
b.A.Foo() // prints "A Foo" (explicit)
```
If this were real inheritance → override would replace the parent's method.  
Here → both methods still exist independently.

# 4. What Bill Kennedy Wants You To Understand

Bill Kennedy’s main teaching is the Go mantra:

> **Composition over inheritance.**

Embedding is simply:

✔ Reuse of state  
✔ Reuse of behavior  
✔ Field/method promotion for convenience

It is NOT:

✘ Type hierarchy  
✘ Polymorphism by default  
✘ Behavioral inheritance  
✘ Subclassing

You _choose_ the relationships explicitly using **interfaces**, not through inheritance trees.

---

# 5. So What’s the Real Mental Model?

### Think of embedding like this:

> **You are _copying_ another type’s fields and methods into your type's namespace.**  
> But you are **not** forming a parent-child relationship.

```bash
B
 ├── X (from embedded A)
 ├── Foo() (maybe from A)
 └── fields/methods of B
```
This lets you "build" types by combining smaller types instead of inheriting.

### Realistic Example: Common Use Case
```go
type Logger struct {
    *log.Logger
}

type Service struct {
    Logger   // embed logger for convenience
    Name string
}
```
You can now do:
```go
s.Println("service started")  // promoted Logger method
```

>**Embedding = reuse with promotion; inheritance = type hierarchy.**  
>**Go only supports the first.**

