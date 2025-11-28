#go #bufio #bufio-reader #bufio-writer #bufio-scanner

## 1. What `bufio` _really_ is

`bufio` gives you **buffered readers and writers** that sit **between** your code and some underlying `io.Reader` / `io.Writer` (like `net.Conn`, `os.File`, `resp.Body`, etc.).

Why?
- To **avoid lots of small syscalls** (`Read` / `Write` on file/conn).
- To give you **nice higher-level primitives** (read a line, read until delimiter, peek, etc.).
- To allow **parsers** (HTTP, JSON, text protocols) to work with efficient buffered data instead of calling the OS every time.

Mental model (read side):
```
[TCP recv buffer / file kernel buffer]
          ↓ (syscall: read)
  net.Conn / *os.File  (io.Reader)
          ↓ (single big read)
       bufio.Reader
          ↓ (many logical reads)
       your code / decoders

```

On write side:
```
your code / encoders
          ↓ (many small writes)
       bufio.Writer
          ↓ (occasional big writes)
  net.Conn / *os.File  (io.Writer)
          ↓ (syscall: write)
[TCP send buffer / file kernel buffer]
```

## 2. `bufio.Reader` – what’s inside & how it behaves

### 2.1 Basic idea

You wrap an `io.Reader`:
```go
r := bufio.NewReader(resp.Body) // or os.File, net.Conn, etc.
```
Internally, `bufio.Reader` keeps:
- a `[]byte` buffer (`buf`)    
- two indices into that buffer:
    - `r` = **read index** (start of unread data)
    - `w` = **write index** (end of unread data)
- an underlying `io.Reader` (`rd`)

At any moment:
```
buf: [ ... consumed ... | unread data | free space .... ]
      0                 r             w
```
- Data valid for you: `buf[r:w]`.
- Free space for next refill: `buf[w:cap(buf)]`.

### 2.2 The core: `Read(p []byte)`

Pseudo-logic (simplified):
```go
func (b *Reader) Read(p []byte) (n int, err error) {
    // 1. If buffer has data, serve from there.
    if b.r < b.w {
        n = copy(p, b.buf[b.r:b.w])
        b.r += n
        return n, nil
    }

    // 2. Buffer is empty: if user's p is larger than our buffer,
    //    read directly from underlying reader into p (bypass buf).
    if len(p) >= len(b.buf) {
        return b.rd.Read(p)
    }

    // 3. Otherwise, refill internal buffer, then copy to p.
    if err := b.fill(); err != nil {
        // fill might return io.EOF but with some data.
    }
    n = copy(p, b.buf[b.r:b.w])
    b.r += n
    // err from fill is returned when buffer consumed
    return n, nil
}
```
Key points:
- If **data is already in the buffer**, `Read` is just a `copy` + index bump → no syscall.
- If **buffer empty**, it calls `rd.Read` once to refill **a big chunk**.
- If **your `p` is huge**, it might bypass the buffer and read directly into `p` (optimization).

So one large kernel read can feed many small logical `Read` calls from your code.

### 2.3 Relation to TCP receive buffer (tying to your earlier mental model)

For something like:
```go
conn := dialTCP(...)
r := bufio.NewReader(conn)
n, err := r.Read(p)
```
The flow:

1. **Kernel** has data in TCP receive buffer.
2. `conn.Read` syscall copies a chunk from kernel buffer → Go heap (`bufio.Reader.buf`).
3. `bufio.Reader.Read` copies from `bufio`’s buffer → your `p`.
    
Once `bufio.Reader` has pulled data into its `buf`, your next reads are:
- Just memory copies within Go.
- Until the buffer empties → next `conn.Read` from kernel buffer.

### 2.4 Higher-level methods (`ReadByte`, `Peek`, `ReadBytes`, `ReadString`, etc.)

This is where `bufio.Reader` becomes _really_ useful.
#### 2.4.1 `ReadByte`
```go
b, err := r.ReadByte()
```
- If `r < w`, it just returns `buf[r]` and increments `r`.    
- If buffer empty, it refills then returns one byte.

This is ideal for byte-by-byte parsers without paying syscall per byte.

#### 2.4.2 `Peek(n int)`
```go
peek, err := r.Peek(4) // look ahead 4 bytes without consuming
```
- Ensures buffer has at least `n` bytes (may call `fill`).
- Returns a slice **referencing the internal buffer**.    
- Important: **That slice becomes invalid after the next read/peek/fill**.
Great for prefix checks:
```go
peek, _ := r.Peek(3)
if string(peek) == "GET" { ... }
```
#### 2.4.3 `ReadSlice(delim byte)`
```go
line, err := r.ReadSlice('\n')
```
- Reads until it finds `delim` (e.g., `\n`) **within the buffer**.
- May call `fill` repeatedly.
- Returns a slice that is **backed by `buf`** – same invalidation rule as `Peek`.
This is what many line-based protocols use under the hood (HTTP request lines, headers, etc.).

#### 2.4.4 `ReadBytes` / `ReadString`
These are like `ReadSlice` but **return a newly allocated slice/string** that remains valid:
```go
b, err := r.ReadBytes('\n')
s, err := r.ReadString('\n')
```
They may do multiple `fill`s and `append` until the delimiter is found.

### 2.5 `ReadLine` (special mention)
```go
line, isPrefix, err := r.ReadLine()
```
- Deals with `\r\n` line endings, etc.
- `isPrefix == true` means the line didn’t fit in the buffer; you must call again and append.
- The docs basically say: for most cases, **use `Scanner` instead**.
## 3. `bufio.Writer` – the write side

`bufio.Writer` collects small writes and flushes them as fewer, larger writes to the underlying writer.
### 3.1 Basic idea
```go
w := bufio.NewWriter(conn)
defer w.Flush()
```
Internally:
- Same `buf` and index, but now:
    - buffer holds data **to be written**.
    - when buffer is full (or `Flush` called), it pushes data to underlying `io.Writer`.
### 3.2 `Write(p []byte)`
Pseudo-logic:
```go
func (b *Writer) Write(p []byte) (n int, err error) {
    for len(p) > 0 {
        if b.available() == 0 {
            if err := b.Flush(); err != nil {
                return n, err
            }
        }

        // If p is larger than buffer and buffer is empty, write p directly.
        if len(p) >= len(b.buf) && b.isEmpty() {
            k, err := b.wr.Write(p)
            n += k
            return n, err
        }

        // Otherwise, copy into buffer.
        k := copy(b.buf[b.n:], p)
        b.n += k
        n += k
        p = p[k:]
    }
    return n, nil
}
```
So again:
- Your code can call `Write` with lots of small chunks.
- Those get aggregated in `bufio` until:
    - buffer gets full → `Flush` automatically,
    - or you call `Flush` manually.
### 3.3 Critical rule: **always flush**
If you wrap your writer, you must eventually:

```
goif err := w.Flush(); err != nil { ... }
```

Common trap:

```go
w := bufio.NewWriter(file) w.WriteString("hello") // program exits, some data may still be in buffer → not written to file
```

Same on TCP: data might be sitting in user-space buffer, not yet written to conn.
## 4. `bufio.Scanner` – high-level token reader
`Scanner` is a helper built on `bufio.Reader` (with its own internal buffer) to **scan tokens**, usually lines.
### 4.1 Typical usage
```go
scanner := bufio.NewScanner(r)
for scanner.Scan() {
    line := scanner.Text()
    // process line
}
if err := scanner.Err(); err != nil {
    // handle error
}
```
Default `SplitFunc`: `ScanLines` (splits on `\n`, handling `\r\n` etc.).
### 4.2 Split functions
You can change how it splits:
```go
scanner.Split(bufio.ScanWords)
```
Or define your own:
```go
scanner.Split(func(data []byte, atEOF bool) (advance int, token []byte, err error) {
    // return how many bytes to advance, the token, and error if any
})
```
Use this for custom protocols / CSV-like formats / etc.
### 4.3 Limitations to know
- Default max token size ~64K.
- If token exceeds this, you get `bufio.ErrTooLong`.
- You can increase it:
```go
buf := make([]byte, 0, 64*1024)
scanner.Buffer(buf, 1024*1024) // max token size 1MB
```
But for **very large records** or complex parsers, `Scanner` isn’t ideal—you’d want your own loop using `Reader`.


## 5. How it fits into HTTP & your earlier mental model
Given your recent questions about:
- kernel recv buffer
- `net.Conn.Read`
- `resp.Body`
- JSON decoder buffer (`dec.buf`)

Let’s stitch it together for HTTP client reading a JSON response.

### 5.1 Possible path (simplified)
1. **Kernel**: response chunks arrive → TCP receive buffer.
2. `net.Conn.Read` (inside `net/http` transport) copies from kernel → Go’s `bufio.Reader.buf`.
3. `net/http` parses status line, headers using buffered methods like `ReadSlice('\n')`.
4. For the body (`resp.Body`), `net/http` wraps the conn into a `*body` which **reads from that same buffered reader**.
5. Your code: `json.NewDecoder(resp.Body)`; the decoder maintains its own `dec.buf` slice:
    - `Decoder` repeatedly calls `resp.Body.Read(dec.buf[...])`.
    - That goes back to `*body` → underlying `bufio.Reader` → underlying `net.Conn.Read` only when `bufio` buffer empties.

So your JSON decoder usually doesn’t touch the kernel directly; it’s reading from a buffered reader sitting on top of `net.Conn`.

On the response side (server):

- You might have:
```go
bw := bufio.NewWriter(conn)
// write status line, headers, body
bw.Flush()
```
- `bufio.Writer` groups small header writes + body chunks into fewer `conn.Write` calls.

## 6. Gotchas & best practices
### 6.1 Don’t mix buffered and unbuffered reads/writes on the same underlying object

Bad:
```go
r := bufio.NewReader(conn)
b1 := make([]byte, 10)
r.Read(b1)         // buffered
conn.Read(b2)      // direct: might read leftover data that bufio "owns"
```
Why: once `bufio.Reader` read from `conn`, some bytes are in the internal buffer. Reading directly from `conn` skips over those bytes; the stream gets out of sync.

Same for writer:
```go
w := bufio.NewWriter(conn)
w.Write([]byte("GET / HTTP/1.1\r\n"))
conn.Write([]byte("\r\n")) // mixing, unsafe wrt ordering and buffering
```
Stick to **all-buffered** or **all-unbuffered** for a given underlying stream.

### 6.2 Mind slices from `Peek` / `ReadSlice` / `ReadLine`
Those slices **alias the internal buffer**:
- Valid only until:
    - next read/fill,
    - or next call that might reuse/shift the buffer.
If you need data later, copy it:
```go
line, err := r.ReadSlice('\n')
if err != nil { ... }
lineCopy := append([]byte(nil), line...)
```
### 6.3 Buffer size tuning

Defaults:
```go
func NewReader(rd io.Reader) *Reader    // default size (4096) 
func NewWriter(wr io.Writer) *Writer
```
For heavy network/file IO, you might want bigger buffers:
```go
r := bufio.NewReaderSize(conn, 32*1024) 
w := bufio.NewWriterSize(conn, 32*1024)
```
Trade-off:
- Larger buffer → fewer syscalls, possibly better throughput.
- But:
    - More memory per connection.
    - Data might sit in user buffer longer before being flushed (if you rely only on auto flush).

### 6.4 Reuse with `Reset`
Useful with pools:
```go
br := bufio.NewReaderSize(nil, 32*1024) 
br.Reset(conn)  

// later: 
br.Reset(anotherConn)`
```
Same for writers / scanners. Reduces allocations for high-throughput servers.

## 7. Some concrete usage patterns
### 7.1 Line reader (classic use)
```go
f, err := os.Open("log.txt")
if err != nil { /* ... */ }
defer f.Close()

r := bufio.NewReader(f)
for {
    line, err := r.ReadString('\n')
    if err == io.EOF {
        if len(line) > 0 {
            fmt.Print(line) // last partial line
        }
        break
    }
    if err != nil { /* ... */ }
    fmt.Print(line)
}
```
Here, `ReadString('\n')`:
- Internally uses buffering and multiple `fill`s if needed.
- Avoids you doing manual loops.
### 7.2 Buffered writer for logs
```go
f, _ := os.Create("out.log")
w := bufio.NewWriterSize(f, 64*1024)
defer w.Flush()

for i := 0; i < 100000; i++ {
    fmt.Fprintf(w, "line %d\n", i) // many small writes, but flushed in big chunks
}
```
### 7.3 Scanner with a custom split
Example: read comma-separated tokens:
```go
scanner := bufio.NewScanner(r)
scanner.Split(func(data []byte, atEOF bool) (advance int, token []byte, err error) {
    if atEOF && len(data) == 0 {
        return 0, nil, nil
    }
    if i := bytes.IndexByte(data, ','); i >= 0 {
        return i+1, data[:i], nil
    }
    if atEOF {
        return len(data), data, nil
    }
    return 0, nil, nil
})
for scanner.Scan() {
    fmt.Println("token:", scanner.Text())
}
if err := scanner.Err(); err != nil {
    log.Fatal(err)
}
```

## 8. When should _you_ explicitly use `bufio`?
Some concrete rules tailored to what you’re building:
- **Network protocols / microservices in Go**:
    - For line-based protocols or custom binary protocols: always wrap `net.Conn` in `bufio.Reader` & `bufio.Writer`.
    - For HTTP server: `net/http` already uses buffered IO internally, so you usually don’t need your own wrapper.
- **File reading**:
    - If reading whole file via `os.ReadFile` – no need.
    - If streaming parse / line by line / chunk by chunk: `bufio.NewReader` / `Scanner` is ideal.
- **Performance-critical paths**:
    - When profiling shows many small reads/writes (e.g., JSON over TCP, high Queries Per Second i.e. high QPS): introduce `bufio` at the edges to smooth them out.

## 🔍 Why QPS matters in backend engineering

When I said:

> “in high QPS systems, buffered IO like `bufio` helps reduce syscalls…”

I meant:
- If your service handles **a lot of requests per second**, even micro-inefficiencies become expensive.
- Unbuffered `Read`/`Write` → many tiny syscalls → kernel overhead → slower performance under load.
- `bufio` batches smaller operations → fewer syscalls → better throughput.

So in backend engineering (your future target → scalable Go services):
- QPS is one of the most important metrics for tuning performance.
- It defines how much traffic your service can sustain before:
    - latency increases,
    - CPU usage spikes,
    - eventually the system saturates.