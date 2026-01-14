#go #logging #program-counter

# help me understand in detail what is the pc in slog package. why is it used.

In Go’s **`log/slog`** package, **`pc`** stands for **Program Counter**.  
It’s a low-level runtime concept that `slog` uses to **precisely identify where a log call originated**.

Let’s unpack this carefully—from CPU basics → Go runtime → why `slog` needs it.

---

## 1️⃣ What is a Program Counter (PC)?

At the machine level, **every executing program has a program counter**:

- It’s a CPU register
    
- It holds the **memory address of the next instruction to execute**
    
- As code runs, the PC moves forward
    

In Go, the **runtime exposes PC values** so tooling can answer questions like:

> “Which function, file, and line number is executing right now?”

---

## 2️⃣ How Go exposes PC information

Go’s runtime lets you retrieve the PC using functions like:

```go
runtime.Caller(skip)
```

This returns:

```go
pc       // program counter
file     // source file
line     // line number
ok
```

Internally:

- `pc` is just a number (an instruction address)
    
- Go can **map PC → function name → file → line** using debug symbols
    

---

## 3️⃣ Where does `pc` appear in `log/slog`?

In `slog`, you’ll see `pc` in places like:

```go
type Record struct {
    Time    time.Time
    Level   Level
    Message string
    PC      uintptr   // ← this
    Attrs   []Attr
}
```

### What does `Record.PC` represent?

> **The program counter of the call site that created the log record**

In plain English:

> “The exact instruction address where `logger.Info(...)` was called.”

---

## 4️⃣ Why doesn’t `slog` just store file & line directly?

Because **PC is more powerful and more flexible**.

### ❌ If slog stored file + line directly

- Every log call would eagerly:
    
    - Walk the stack
        
    - Resolve symbols
        
    - Convert PC → file/line
        
- This is **slow**, especially in hot paths
    

### ✅ By storing only PC

- `slog` captures a **cheap, opaque identifier**
    
- Resolution is **deferred**
    
- Handlers decide:
    
    - Whether they even need source info
        
    - How expensive resolution should be
        

This matches Go’s philosophy:

> _Pay only for what you use_

---

## 5️⃣ How PC is used in practice

### Step 1: Log call

```go
logger.Info("db connected")
```

Internally:

- `slog` captures the **PC of this call**
    
- Stores it in `Record.PC`
    

### Step 2: Handler decides what to do

For example, `TextHandler` or `JSONHandler`:

```go
HandlerOptions{
    AddSource: true,
}
```

When `AddSource` is enabled:

- Handler resolves:
    
```go
runtime.CallersFrames([]uintptr{record.PC})
```
    
- Extracts:
    
    - function name
        
    - file
        
    - line
        

### Step 3: Output

```json
{
  "msg": "db connected",
  "source": {
    "file": "db/conn.go",
    "line": 42,
    "function": "Connect"
  }
}
```

If `AddSource = false`:

- PC is **never resolved**
    
- Zero overhead
    

---

## 6️⃣ Why PC is especially important for structured logging

`slog` is designed for:

- JSON logs
    
- Observability pipelines
    
- OpenTelemetry / Loki / ELK
    
- High-volume production logging
    

PC enables:

### 🔹 Accurate source attribution

Even when logs go through:

- helper functions
    
- wrappers
    
- shared logging utilities
    

### 🔹 Correct stack depth handling

`slog` can do:

```go
runtime.Callers(skip + handlerDepth)
```

to avoid blaming logging wrappers instead of real callers.

### 🔹 Zero-allocation logging paths

PC is:

- just a `uintptr`
    
- no strings
    
- no heap allocation
    

---

## 7️⃣ Why this matters for YOU as a backend engineer

This design shows **production-grade thinking**:

|Design Choice|Reason|
|---|---|
|Store PC instead of file/line|Defer expensive work|
|Resolve in handler|Customizable cost|
|Use runtime metadata|Precise source mapping|
|Make it optional|Performance first|

This same pattern appears in:

- tracing systems
    
- profilers
    
- debuggers
    
- OpenTelemetry SDKs
    

---

## 8️⃣ Mental model (simple analogy)

Think of **PC as a GPS coordinate**, not an address string.

- PC = latitude/longitude
    
- File + line = street name + house number
    
- Resolve only if someone asks for directions
    

---

## 9️⃣ TL;DR

- **PC = Program Counter**
    
- In `slog`, it stores **where the log call originated**
    
- It is:
    
    - cheap to capture
        
    - expensive to resolve
        
- `slog` defers resolution to:
    
    - improve performance
        
    - allow flexible handlers
        
- This is why `slog` scales well in production systems

# trace exactly how runtime.Callers works and inspect how slog calculates stack depth

I’ll walk through (1) what `runtime.Callers` actually captures and how frames are resolved, then (2) how `slog` picks the right “skip” so your wrapper functions don’t become the reported source.

### 1) `runtime.Callers` — what it _really_ does

At a high level, `runtime.Callers(skip, pc []uintptr) int`:

1. **Walks the current goroutine’s call stack**
    
2. **Collects return PCs** (program counters) for each stack frame
    
3. Writes them into the provided `pc` slice
    
4. Returns how many it wrote
    

#### Important details that matter in practice

**A) It captures “return PCs,” not exact call-site PCs**

- The PCs you get are typically _the return address_ for a call instruction.
    
- When you later map PC → file/line, Go adjusts things so you still get the correct line most of the time (but this explains some “off by one line” weirdness you may have seen).
    

**B) `skip` is counted in stack frames**

- `skip=0` includes `runtime.Callers` itself (not useful).
    
- You usually do `skip=1` or more.
    
- `skip` is “how many frames above this function to start from.”
    

Example mental model:

```
yourFunc -> helper -> runtime.Callers -> (collect frames)
```

- `skip=0` starts at `runtime.Callers`
    
- `skip=1` starts at `helper`
    
- `skip=2` starts at `yourFunc`
    

**C) It’s intentionally low-level and fast**

- It does _not_ do symbolization (function/file/line) itself.
    
- It just grabs PCs from the runtime’s stack metadata.
    

---

### 1.1) Turning PCs into readable frames: `runtime.CallersFrames`

The second step is:

```go
frames := runtime.CallersFrames(pcs)
for {
    fr, more := frames.Next()
    // fr.Function, fr.File, fr.Line, fr.Entry
    if !more { break }
}
```

This is where:

- PC → function name
    
- PC → file/line
    
- inlining is handled (it may return multiple logical frames for inlined calls)
    

**Inlining is a big reason `CallersFrames` exists.**  
Modern Go aggressively inlines. The runtime can still report “logical frames” even if no physical stack frame exists for that call.

So:

- `runtime.Callers` gives you raw PCs
    
- `runtime.CallersFrames` interprets them, respecting inlining info
    

---

### 1.2) A small “trace” you can run mentally

If you write:

```go
func A() { B() }
func B() { C() }
func C() {
    pcs := make([]uintptr, 10)
    n := runtime.Callers(0, pcs)
    _ = pcs[:n]
}
```

The first few “physical” frames will include:

- `runtime.Callers`
    
- `C`
    
- `B`
    
- `A`
    
- `main.main`
    
- runtime startup frames…
    

If you instead do `runtime.Callers(2, pcs)` you’ll usually start at `B` (because you skipped `runtime.Callers` and `C`).

---

## 2) How `slog` calculates stack depth (and why it’s tricky)

`slog` wants the `Record.PC` to point at the **user’s call site** like:

```go
logger.Info("hello") // <-- source should be here
```

Not inside:

- `(*Logger).Info`
    
- `(*Logger).log`
    
- the handler chain
    

### 2.1) The core idea: fixed “skip” + optional user-supplied adjustment

Conceptually `slog` does something like:

- Start with a **baseline skip** value that jumps over:
    
    - runtime frames
        
    - internal slog frames (`Logger.Info`, `Logger.log`, etc.)
        
- Then apply an **extra offset** to handle wrappers:
    
    - `HandlerOptions.ReplaceAttr`? not relevant
        
    - **`HandlerOptions.AddSource`** triggers source resolution later
        
    - **Wrapper loggers** are the big reason you need adjustment
        

In many logging libraries you’ll see something like `AddCallerSkip(n)`.  
In `slog`, the equivalent mechanism is typically done by exposing ways to create a logger that accounts for wrappers (in stdlib this shows up as “skip” being part of internal record creation logic; many people implement their wrapper by capturing PC themselves or by using `runtime.Callers` with a higher skip and logging via `LogAttrs`).

### 2.2) Why stack depth isn’t just “skip 2”

Because there are multiple layers:

**User code**  
→ maybe your wrapper `LogInfo(...)`  
→ `slog.Logger.Info`  
→ `slog.Logger.log`  
→ record creation  
→ handler

If you always skip the same amount, you’ll often blame:

- your wrapper function (bad), or
    
- slog internals (worse)
    

So slog needs a “known good baseline” for its own internals, and then you (or slog) may need to add _more_ skip for wrappers.

---

## 3) A practical way to _inspect_ both in your own code

If you want to actually see the stack frames and compare them with what `slog` chooseseems to pick, here’s a minimal snippet:

```go
package main

import (
	"fmt"
	"log/slog"
	"os"
	"runtime"
)

func dumpStack(tag string, skip int) {
	pcs := make([]uintptr, 32)
	n := runtime.Callers(skip, pcs)
	frames := runtime.CallersFrames(pcs[:n])

	fmt.Println("----", tag, "skip =", skip, "----")
	for i := 0; i < 8; i++ { // print first few
		fr, more := frames.Next()
		fmt.Printf("%2d %s\n   %s:%d\n", i, fr.Function, fr.File, fr.Line)
		if !more {
			break
		}
	}
}

func wrapper(l *slog.Logger, msg string) {
	dumpStack("inside wrapper", 0)
	l.Info(msg)
}

func main() {
	l := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{AddSource: true}))

	dumpStack("in main before wrapper", 0)
	wrapper(l, "hello")
}
```

What you’ll learn from this:

- Which frame corresponds to **your actual call site**
    
- How many frames you must skip to land exactly there
    
- Whether your wrapper adds an extra layer that `slog` will report as the source
    

Once you see the stack, “stack depth” stops being abstract.

---

## 4) What to do when you have your own wrapper

If you have something like:

```go
func Info(msg string, args ...any) { logger.Info(msg, args...) }
```

Then **source will point to `Info(...)`**, not the real caller.

Fix patterns (common in Go):

### Option A: Don’t wrap; pass logger around

Call `logger.Info` directly from call sites.

### Option B: Capture caller PC yourself in wrapper (advanced)

You can grab the right PC using `runtime.Callers` with a skip that jumps over wrapper frames, then log with `LogAttrs` in a way that preserves/uses that PC (this is where people sometimes customize handlers, because stdlib keeps PC inside the record creation path).

### Option C: Provide wrapper APIs that accept context-like metadata

So you don’t need wrappers everywhere (less “skip” pain).