#go #go-embed #go-interface #polymorphism

### A Reuse: embedding a common implementation

Suppose you have multiple HTTP handlers that all need:
- logging
- metrics
- some common pre-checks

You can build a reusable “base”:

```go
type BaseHandler struct {
    Logger *log.Logger
}

func (h BaseHandler) logRequest(r *http.Request) {
    h.Logger.Println("request:", r.URL.Path)
}
```

Then embed it:

```go
type UserHandler struct {
    BaseHandler
}

func (h UserHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    h.logRequest(r)
    // user-specific logic
}

type OrderHandler struct {
    BaseHandler
}

func (h OrderHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    h.logRequest(r)
    // order-specific logic
}
```

Here:

- `BaseHandler` gives common logic
- `UserHandler` / `OrderHandler` each have their own specific logic

---

### B. Polymorphism: both types implement an interface

In Go, **interfaces are satisfied implicitly**.

```go
type Handler interface {
    ServeHTTP(http.ResponseWriter, *http.Request)
}
```

Because both `UserHandler` and `OrderHandler` define `ServeHTTP`, they automatically satisfy `Handler`:

```go
func RegisterRoutes(mux *http.ServeMux, handlers map[string]Handler) {
    for path, h := range handlers {
        mux.Handle(path, h)
    }
}

func main() {
    logger := log.New(os.Stdout, "[svc] ", log.LstdFlags)

    mux := http.NewServeMux()

    handlers := map[string]Handler{
        "/users":  UserHandler{BaseHandler{Logger: logger}},
        "/orders": OrderHandler{BaseHandler{Logger: logger}},
    }

    RegisterRoutes(mux, handlers)
    http.ListenAndServe(":8080", mux)
}
```

Here you get:
- **Reuse**: via `BaseHandler` embedding
- **Polymorphism**: via `Handler` interface

This is the Go way to express what OOP would do with:
- `abstract class BaseHandler implements IHandler`
- concrete classes overriding `ServeHTTP`

---

### C. Embedding types that already implement interfaces

Another nice pattern: your embedded type already satisfies some interface, and your outer type **inherits that behavior** (not the type).

Example:

```go
type JobRunner interface {
    Run(ctx context.Context) error
}

type JobFunc func(ctx context.Context) error

func (f JobFunc) Run(ctx context.Context) error {
    return f(ctx)
}

```

Now you can embed `JobFunc`:

```go
type LoggedJob struct {
    JobFunc
    Logger *log.Logger
}

func NewLoggedJob(logger *log.Logger, fn func(ctx context.Context) error) LoggedJob {
    return LoggedJob{
        JobFunc: JobFunc(fn), // embedded
        Logger:  logger,
    }
}

// Optionally shadow Run:
func (j LoggedJob) Run(ctx context.Context) error {
    j.Logger.Println("job starting")
    err := j.JobFunc.Run(ctx) // delegate
    j.Logger.Println("job finished, err:", err)
    return err
}
```

Usage:

```go
func RunAll(ctx context.Context, jobs []JobRunner) error {
    for _, j := range jobs {
        if err := j.Run(ctx); err != nil {
            return err
        }
    }
    return nil
}
```

- `JobFunc` implements `JobRunner`.
- `LoggedJob` embeds `JobFunc`, and also implements `JobRunner` (because it has `Run`).
- You can plug in either `JobFunc` or `LoggedJob` into something that expects `JobRunner`.

Again: **behavior-based polymorphism**, not inheritance.

---

## Mental model to carry forward

**For “inheritance-ish” needs:**

- **Shared state / logic?** → use a separate struct and **embed** it (or use as a field).
- **Want polymorphism?** → define an **interface**, and have multiple concrete types implement it.
- **Want different behavior in different “subclasses”?**  
    → inject behavior via interfaces/functions (strategy pattern), or shadow methods only where really needed.

A nice Bill-style mantra you can keep in your head:
> - Use **structs** to model _what you have_ (state).
> - Use **methods** to model _what you can do_ (behavior).
> - Use **interfaces** to model _what is required_ (contracts).
> - Use **embedding** to _reuse_ state/behavior, not to build hierarchies.