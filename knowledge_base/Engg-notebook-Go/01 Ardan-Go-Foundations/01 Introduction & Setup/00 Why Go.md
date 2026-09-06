#engineers-notebook #golang
### Free lunch is over: 
- Until clock frequencies were getting higher, the lunch was free
- Since when clock frequencies started staying about the same and we now have more CPU cores
- We need to write code that can utilize these cores well.
- Go has it built in the language, **go-routines** and **channels** that allow utilize cores well
### C10K problem:
- How can a single application serve 10k (10-thousand) connections concurrently
- Most languages do this using async-I/O, (the event loop)
- Go does it behind the scenes for us, and comes with a **production ready HTTP server** (TLS, HTTP 2, ...)

### Built for Large Teams
- Small Simple language
- Module system - reusability
- Interfaces - modularity

### Robust & Productive
- Static types
- Garbage collector
- Rich and mature standard library
- Fast compilation
- Forces to check errors

### Great Tooling
- Tools for build, run, test, **benchmark**, install etc ....
- Modules - dependency management
- Built-in **profiler and tracer** with web interface
- Built-in logging and metrics
- Focused on performance and so do many tools around it

### Saves money
- Stable API
- need less servers doing same workload
- Compiles to static executable - easy deployment
- Easy to cross compile - can produce binaries for any platform (mac or binaries) also on a Linux CI system