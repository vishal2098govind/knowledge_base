#go #go-embed #inheritance #composition-over-inheritance 

understanding **upcasting** is key to seeing why Go embedding ≠ inheritance.

# What “Upcasting” Means (in OOP languages)
In classical OOP (Java, C#, C++), if you have:
```
Parent
↑
Child  (inherits Parent)
```
Then **a Child _is-a_ Parent**, so you can _upcast_:
```
Child c = new Child();
Parent p = c;   // <-- UPCAST: child → parent
```
This is allowed because the language has _inheritance_ and _type hierarchies_.  
The child is considered a subtype of its parent.

### Upcasting definition:

> **Assigning a more specific type (child) to a more general type (parent).**

Examples:
- Dog → Animal
- SavingsAccount → BankAccount
- Square → Shape

# Go does NOT allow upcasting with embedding
Take this Go code:
```go
type A struct {
    X int
}

type B struct {
    A   // embedded
    Y int
}
```
If Go embedding were real inheritance, you could do:
```go
var b B
var a A = b   // ❌ DOES NOT COMPILE
```
Because **B is NOT a subtype of A**.

`B` _contains_ an `A`, but it **is not an A**.

No inheritance → no upcasting.

---

# Why doesn’t Go allow it?

Because embedding is just **composition + field/method promotion**.

There is _no_ type hierarchy like:
```
A
↑
B
```
Instead, it's like:
```
B
 ├── A (as a field)
 └── Y
```

# Quick visual comparison
In OOP inheritance (real upcasting)
```java
Dog isa Animal

Animal a = new Dog();   // UPCAST
```
### **In Go embedding**

```
Dog has an Animal
```

So:
```go
var d Dog var a Animal = d  // ❌ impossible
```

You _must_ access the embedded field explicitly:
```go
var a Animal = d.Animal  // ✔ valid
```

### Upcast = treat a child as a parent.

### Go has no parent-child relationships.

### Therefore, no upcasting with embedding.

# Go’s Philosophy: Why No Inheritance?
Go _intentionally_ does **not** support inheritance, not because they couldn’t add it, but because it contradicts the core philosophy of the language.

This is one of the most important things to internalize when becoming an expert Go backend engineer.  
Let’s go into the **why**, based on Go’s design principles and what Bill Kennedy, Rob Pike, and the Go team repeatedly emphasize.
## 1. Inheritance creates _tight coupling_
In OOP languages, inheritance ties two types together _forever_:
```
Child → Parent
```
This creates problems:
- Changing the parent breaks all children.
- Deep hierarchies become rigid and unmaintainable.
- You get "fragile base class" issues.
- You unintentionally inherit behaviour you did NOT want.

Go _hates_ tight coupling.  
Go wants **types to be independent and explicit**, not bound by hierarchy.

## 2. Inheritance encourages deep, complex type trees

Go's creators observed that OOP inheritance leads to:

- **Class hierarchies that grow uncontrollably**
- Hard-to-understand relationships
- Hidden behaviour coming from far up the tree
- "Too much machinery" in the type system

Rob Pike famously said:

> _“The big mistake in object-oriented languages is the type hierarchy.  
> Go takes a different approach.”_

Go tries to keep the mental model **flat and simple**.

## 3. Inheritance brings _implicit_ behaviour

If `B` inherits from `A`, it suddenly gets:
- A’s fields
- A’s methods
- Overridable methods
- Virtual dispatch

This automatic, magical behaviour is the opposite of Go’s values:

> **Go values explicit, simple, predictable code.**

Embedding _looks_ like inheritance, but it is not — because Go wants you to see exactly what is happening.

No hidden behaviour.

## 4. Go prefers _composition over inheritance_

This is the #1 reason.

Instead of saying:

> “B _is_ an A”

Go wants you to model:

> “B _has_ an A”

This enables:

- flexible design
- no type hierarchy
- reusable components
- decoupling

This aligns with modern system design (micro-services, cloud-native, etc.).

## 5. Go uses _interfaces_ for polymorphism, not inheritance

Go’s interfaces are **implicit**, meaning types automatically satisfy them.

Example:
```go
type Reader interface {
    Read(p []byte) (int, error)
}

type File struct {}

func (File) Read(p []byte) (int, error) { ... }
```

`File` automatically implements `Reader`.  
No need to write:
```java
class File implements Reader
```

This leads to:
- zero boilerplate
- high flexibility
- clean abstractions without brittle hierarchies

## 6. Go avoids the “diamond problem”

In multiple inheritance languages:
```
     A
   /   \
  B     C
   \   /
     D
```
Which `A` method does `D` inherit?  
This leads to complicated rules.

Go wants NONE of that complexity.

Interfaces + embedding solve these problems cleanly.

## 7. Go’s goal is _simplicity + readability_ above all

Rob Pike:

> _“One of Go’s successes is its simplicity.  
> We rejected many features that introduce cognitive load.”_

Inheritance was rejected because:

- It complicates the type system
- It encourages design patterns that lead to rigid software
- It results in surprise behaviour
- It makes reasoning about code harder

Go wants the developer to read a file and understand everything _right there_.
## Go’s Philosophy = Simplicity + Clarity + Composition + Interfaces

- **Simple syntax**
- **Predictable behaviour**
- **No hierarchy**
- **Small surface area**
- **Explicit code**
- **Composition over inheritance**
- **Interfaces for polymorphism**
- **Encapsulation without magic**

This is why Go codebases stay **clean**, **stable**, and **maintainable**, even at very large scale (Kubernetes, Docker, AWS CDK, etc.).

>**Go doesn’t support inheritance because inheritance hides behavior, creates fragile hierarchies, encourages tight coupling, and complicates the type system.**  
>**Go instead chooses composition and interfaces to keep software simple, explicit, and maintainable.**

