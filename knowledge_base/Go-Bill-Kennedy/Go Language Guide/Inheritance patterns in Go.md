#go #go-embed #composition-over-inheritance 

### Pattern A: “Base class” for shared fields & methods

**In OOP**

```java
class Animal {
    protected String name;
    public void Speak() { System.out.println("..."); }
}

class Dog extends Animal {
    @Override
    public void Speak() { System.out.println("Woof"); }
}
```


**In Go: use composition (optionally embedding)**

```go
type Animal struct {
    Name string
}

func (a Animal) Speak() {
    fmt.Println("...")
}

type Dog struct {
    Animal // embed for convenience
    Breed string
}

// Optionally shadow Speak:
func (d Dog) Speak() {
    fmt.Println("Woof, I'm", d.Name)
}
```

Key idea: **Dog has an Animal**, not “is-an Animal”.  
You get reuse of fields and methods (via promotion) but no type hierarchy.

### Pattern B: “Abstract base class” + concrete subclasses

**In OOP**

```java
abstract class Shape {
    abstract double area();
}

class Circle extends Shape { ... }
class Rect extends Shape { ... }
```

**In Go: use interfaces for behavior**

```go
type Shape interface {
    Area() float64
}

type Circle struct {
    Radius float64
}

func (c Circle) Area() float64 {
    return math.Pi * c.Radius * c.Radius
}

type Rect struct {
    W, H float64
}

func (r Rect) Area() float64 {
    return r.W * r.H
}

func TotalArea(shapes []Shape) float64 {
    var sum float64
    for _, s := range shapes {
        sum += s.Area()
    }
    return sum
}
```

No base class, no inheritance.  
Just a **contract (interface)** + multiple types that satisfy it.

### Pattern C: “Overriding methods” in subclasses

In inheritance, you often:

- Put logic in the base class
- Override in child classes for custom behavior (Template Method Pattern, etc.)

**In Go, there are two typical approaches:**

#### 1. Use *strategy objects* (functions or interfaces)

Instead of overriding methods, you **inject behavior**.

```go
type Storage interface {
    Save(ctx context.Context, data []byte) error
}

type FileStorage struct { /* ... */ }
func (f FileStorage) Save(ctx context.Context, data []byte) error { /* ... */ return nil }

type S3Storage struct { /* ... */ }
func (s S3Storage) Save(ctx context.Context, data []byte) error { /* ... */ return nil }

type Service struct {
    Store Storage  // behavior injected
}

func (s Service) Process(ctx context.Context, data []byte) error {
    // common processing logic
    return s.Store.Save(ctx, data) // polymorphic call
}
```

Here in OOP you might have:

- `BaseService` with an overridable `Save()`
    
- `FileService`, `S3Service` overriding `Save()`
    

In Go:

- `Service` is composed with a `Storage` interface implementation.
    
- **No inheritance; same flexibility.**
    

#### 2. Use embedding + method shadowing (carefully)

You can embed a “base” type and shadow methods when needed:

```go
type BaseService struct {}

func (BaseService) Handle(req Request) {
    fmt.Println("base handling")
}

type CustomService struct {
    BaseService
}

func (CustomService) Handle(req Request) {
    fmt.Println("custom handling")
}
```

Call sites:

```go
var s CustomService
s.Handle(req)        // custom
s.BaseService.Handle(req) // base
```

This looks like overriding, but it’s **just shadowing**; both methods coexist. It’s useful, but you use it **deliberately**, not as the default design pattern.

### Pattern D: “Mixins” or trait-like reuse

In some OOP languages you do:

- Mixin for logging
- Mixin for metrics
- Mixin for retry behavior

In Go, you do **small focused types** embedded or held as fields:

```go
type Logger struct {
    *log.Logger
}

type Metrics struct {
    // counters, histograms, etc.
}

type Service struct {
    Logger
    Metrics
}
```

Now `Service` gets:
- `Service.Println(...)` (promoted from Logger)
- Methods from `Metrics` (possibly promoted)

Again: no inheritance. Just **building blocks**.