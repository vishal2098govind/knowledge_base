#engineers-notebook #golang #n
# Go's Netpoller Handles Millions of Connections. But Connections Still Aren't Free.

  

- Date: 2026-05-03

- URL: https://logs.craftedbyvishal.dev/p/gos-netpoller-handles-millions-of

  

---

  

A while back I wrote my first Substack article about building a raw TCP server in Go. You can read it [here](https://vishalageek.substack.com/p/i-thought-go-concurrency-was-about). It ended here:

func main() {

lis, err := net.Listen("tcp", ":8080")

// ...

for {

conn, err := lis.Accept()

// ...

go func(conn net.Conn) {

sc := bufio.NewScanner(conn)

who := conn.RemoteAddr()

for sc.Scan() {

fmt.Printf("[%s]: %s\n", who, sc.Text())

}

// ...

}(conn)

}

}

  

Concurrent. Each client got its own goroutine. The server stopped being unavailable.

  

But I had a question I couldn’t shake.

  

What actually happens inside that goroutine? What is `conn`? What did `net.Listen` do before `Accept()` was even called? What does the kernel know about this connection that my code doesn’t?

  

I knew the textbook answer. In college I drew this:

  

[![](https://substackcdn.com/image/fetch/$s_!TysG!,w_1456,c_limit,f_auto,q_auto:good,fl_progressive:steep/https%3A%2F%2Fsubstack-post-media.s3.amazonaws.com%2Fpublic%2Fimages%2F98c06e5e-899c-48f2-b008-5949c37597ff_8840x8939.png)](https://substackcdn.com/image/fetch/$s_!TysG!,f_auto,q_auto:good,fl_progressive:steep/https%3A%2F%2Fsubstack-post-media.s3.amazonaws.com%2Fpublic%2Fimages%2F98c06e5e-899c-48f2-b008-5949c37597ff_8840x8939.png)

  

SYN, SYN-ACK, ACK. Three-way handshake. Four-way teardown. The state machine labels on the left for the client, on the right for the server. I could reproduce those arrows from memory.

  

What I never understood was what each arrow cost. What got allocated when the SYN arrived. Who actually sent the SYN-ACK. What Go does when you call `conn.Read()` and there’s nothing there yet. What TIME_WAIT actually is and why it exists.

  

I wanted to see this at the kernel level, not in theory, but through code I wrote and understood line by line.

  

So I built a raw TCP echo server and client in Go, no `net/http`, no frameworks, just `net.Dial`, `net.Listen`, `net.Conn`, and drew out every layer until I understood it.

  

This is what I found.

  

* * *

  

## First, What Is a File Descriptor

  

Before any of this makes sense, one concept is worth grounding first.

  

A file descriptor is just an integer. `fd=5`, `fd=10`, `fd=3`. That’s all it is on the surface.

  

Every process has its own per-process file descriptor table, a kernel-maintained mapping from these integers to actual resources in kernel memory. Those resources can be files on disk, pipes, or in our case, sockets. `fd=5` in our Go server process points to one socket. `fd=5` in a completely different process points to something else entirely. The integer means nothing across process boundaries. Within a process, it is unique.

  

When you call `socket()`, the kernel allocates a `tcp_sock` structure in RAM and returns the next available integer from our process’s fd table. That number is your handle. Every syscall you make, `bind()`, `listen()`, `accept()`, `read()`, `write()`, `close()`, takes that integer and the kernel looks up the actual resource behind it.

  

This is why “frees fd=10 from the process fd table” during teardown is significant. The integer `10` is released back to the table. The kernel can hand it to the next `accept()` call for a completely different client. The underlying socket structure may still exist in kernel memory, in TIME_WAIT, but our process has no handle to it anymore. The fd and the socket are two different things that happen to be linked for the duration of a connection.

  

* * *

  

## A Connection Is Not What I Thought It Was

  

Before I wrote any code, I assumed a “connection” was the thing `listener.Accept()` returned.

  

It’s not. Or rather, that’s the end of a much longer process.

  

When you call `net.Listen("tcp", ":4040")`, three syscalls happen in sequence.

  

[![](https://substackcdn.com/image/fetch/$s_!NDFo!,w_2400,c_limit,f_auto,q_auto:good,fl_progressive:steep/https%3A%2F%2Fsubstack-post-media.s3.amazonaws.com%2Fpublic%2Fimages%2Ffa06f129-ea16-4e88-8b5f-891b02afed69_1168x2477.png)](https://substackcdn.com/image/fetch/$s_!NDFo!,f_auto,q_auto:good,fl_progressive:steep/https%3A%2F%2Fsubstack-post-media.s3.amazonaws.com%2Fpublic%2Fimages%2Ffa06f129-ea16-4e88-8b5f-891b02afed69_1168x2477.png)

  

`socket()` asks the kernel to create a TCP socket, a data structure in kernel memory, and hand back a file descriptor. That fd, say `fd=5`, is just a number. A handle. The socket itself lives in the kernel. At this point it has no address, no port. It just exists.

  

`bind(fd, "0.0.0.0", 4040)` attaches the socket to a port. The kernel checks if port 4040 is already in use, checks process permissions, then updates its internal port-to-socket mapping. Now `fd=5` owns `:4040`.

  

`listen(fd)` is where something interesting happens. The socket transitions from a plain socket to a listening socket, and the kernel creates two queues:

tcp_sock

state: LISTEN

local: :4040

syn_queue: []

accept_queue: []

  

These two queues are the kernel’s waiting room for incoming clients. `net.Listen` returns. Our Go code has a listener. Nothing has connected yet.

  

And then our goroutine calls `listener.Accept()`. If the accept queue is empty, it parks. How exactly — that's what the next section is about. But the short version: Go's netpoller watches the fd on our behalf, and wakes the goroutine when something arrives. More on the netpoller in detail further below.

  

* * *

  

## The Handshake Happens Without You

  

When a client sends a SYN packet, our goroutine is not involved.

  

Not partially involved. Not notified. Completely absent.

  

[![](https://substackcdn.com/image/fetch/$s_!PFNi!,w_2400,c_limit,f_auto,q_auto:good,fl_progressive:steep/https%3A%2F%2Fsubstack-post-media.s3.amazonaws.com%2Fpublic%2Fimages%2Fe7575a44-67aa-430e-9565-affb48853985_1820x2604.png)](https://substackcdn.com/image/fetch/$s_!PFNi!,f_auto,q_auto:good,fl_progressive:steep/https%3A%2F%2Fsubstack-post-media.s3.amazonaws.com%2Fpublic%2Fimages%2Fe7575a44-67aa-430e-9565-affb48853985_1820x2604.png)

  

The SYN packet arrives at the NIC. The NIC fires a hardware interrupt. The kernel’s interrupt handler picks it up. The packet walks up the network stack, Ethernet → IP → TCP. The kernel looks up which socket owns the destination port. It finds `fd=5`.

  

Then the kernel creates a new socket for this specific client, `cli:tcp_sock`, with the client’s address already filled in, and places it in the SYN queue. Sends SYN-ACK. When the client’s ACK arrives, the kernel moves `cli:tcp_sock` to the accept queue and marks it ESTABLISHED.

  

Our goroutine is parked the entire time. The three-way handshake completes before `Accept()` returns. The goroutine only wakes up to collect the result.

  

That was the first thing that surprised me. I had always imagined my server code as participating in the handshake. It doesn’t. The kernel does all of it.

  

Look at the diagram above again. The `g_server` column is completely empty throughout the handshake. That emptiness is the point.

  

There’s one more detail worth noting. The kernel creates `cli:tcp_sock` when the SYN arrives, not when the ACK arrives. It pre-allocates the socket structure mid-handshake. The connection only becomes real, both sides agreeing on sequence numbers, when the final ACK lands. But the kernel commits resources earlier.

  

* * *

  

## What “Blocking” Actually Means in Go

  

This is where Go’s netpoller comes in.

  

The netpoller is a component of the Go runtime that sits between our goroutines and the Linux kernel’s epoll mechanism. Its job is simple: watch file descriptors for I/O readiness events, and when one becomes ready, wake up whichever goroutine was waiting on it. It runs on its own OS thread, calling `epoll_wait()` in a loop. Our goroutines never call `epoll_wait()` directly. They park, and the netpoller does the watching.

  

This is what makes Go capable of handling millions of concurrent connections without millions of OS threads. The netpoller converts what looks like blocking network I/O in your code into event-driven I/O underneath — our goroutine writes `conn.Read()` as if it blocks, but the runtime never actually blocks the OS thread. Instead it parks the goroutine and frees the thread to run other goroutines. The OS thread stays available for CPU work or other connections while the kernel watches the fd. When data arrives, the netpoller wakes the right goroutine and the thread picks it up again.

  

Park goroutines. Free threads. Handle more. That’s the model.

  

[![](https://substackcdn.com/image/fetch/$s_!V3Xr!,w_2400,c_limit,f_auto,q_auto:good,fl_progressive:steep/https%3A%2F%2Fsubstack-post-media.s3.amazonaws.com%2Fpublic%2Fimages%2Fbd378160-5d67-467e-a833-abc383ff2173_1549x2043.png)](https://substackcdn.com/image/fetch/$s_!V3Xr!,f_auto,q_auto:good,fl_progressive:steep/https%3A%2F%2Fsubstack-post-media.s3.amazonaws.com%2Fpublic%2Fimages%2Fbd378160-5d67-467e-a833-abc383ff2173_1549x2043.png)

  

The netpoller tracks which goroutine is waiting on which fd using a `pollDesc` — a small struct it maintains for every registered fd:

type pollDesc struct {

fd int

rg *g // goroutine waiting to read

wg *g // goroutine waiting to write

}

  

This is the bridge between a file descriptor and a goroutine. When `g_server` calls `Accept()` and the accept queue is empty, the full sequence is:

  

* The runtime stores `g_server` in `pollDesc[fd=5].rg`

  

* Registers `fd=5` with epoll via `epoll_ctl(ADD, fd=5, EPOLLIN | EPOLLET)`

  

* Calls `gopark()` — `g_server` is removed from the run queue, the scheduler’s list of goroutines ready to execute

  

* The OS thread that was running `g_server` is freed and picks up other goroutines

  
  
  
  

When the handshake completes and something lands in the accept queue, the netpoller’s `epoll_wait()` returns `fd=5` as ready. The netpoller looks up `pollDesc[fd=5].rg`, finds `g_server`, clears the field, and calls `runqueue.add(g_server)` — putting it back on the scheduler’s run queue. An OS thread picks it up and resumes it exactly where it parked, inside `Accept()`, as if nothing happened in between.

  

`Accept()` makes the syscall. The kernel dequeues `cli:tcp_sock` from the accept queue and returns `conn (fd=10)`.

  

`g_server` immediately calls `go handleConn(conn)`, spawning `g_conn_1` with ownership of the connection, and loops back to `listener.Accept()`. The infinite for loop. The server is available again before `g_conn_1` has done anything at all.

  

This is the whole model:

func main() {

l, err := net.Listen("tcp", ":4040")

// ...

for {

conn, err := l.Accept()

// ...

go handleConn(conn)

}

}

  

`g_server` parks at `Accept()`. Kernel does the handshake. Netpoller wakes `g_server`. `g_server` spawns `g_conn_1` and parks again. No OS thread sleeps waiting for a client. No thread is wasted.

  

* * *

  

## Inside the Connection Handler

  

`g_conn_1` takes over `conn (fd=10)` and calls `conn.Read(buf)`. There’s no data yet, the client hasn’t sent anything. Same pattern repeats.

  

[![](https://substackcdn.com/image/fetch/$s_!xCxG!,w_2400,c_limit,f_auto,q_auto:good,fl_progressive:steep/https%3A%2F%2Fsubstack-post-media.s3.amazonaws.com%2Fpublic%2Fimages%2Fd8b61088-29be-40b9-8b88-720bee9dfb83_997x1017.png)](https://substackcdn.com/image/fetch/$s_!xCxG!,f_auto,q_auto:good,fl_progressive:steep/https%3A%2F%2Fsubstack-post-media.s3.amazonaws.com%2Fpublic%2Fimages%2Fd8b61088-29be-40b9-8b88-720bee9dfb83_997x1017.png)

  

The Go runtime registers `fd=10` with epoll, `epoll_ctl(ADD, fd=10, EPOLLIN | EPOLLET)`, stores `g_conn_1` in `pollDesc[fd=10].rg`, and parks it.

  

Now look at the pollDesc table:

{fd:5, rg: g_server, wg: nil}

{fd:10, rg: g_conn_1, wg: nil}

  

Two goroutines. Two fds. One epoll instance. Zero OS threads consumed by either.

  

This is why Go can handle tens of thousands of concurrent connections with a handful of OS threads. Every goroutine waiting for I/O is just a data structure in memory. The OS threads are free.

  

[![](https://substackcdn.com/image/fetch/$s_!LWBG!,w_2400,c_limit,f_auto,q_auto:good,fl_progressive:steep/https%3A%2F%2Fsubstack-post-media.s3.amazonaws.com%2Fpublic%2Fimages%2Fff244cb7-ae1d-41ee-8baa-4ad7d353a204_1765x2606.png)](https://substackcdn.com/image/fetch/$s_!LWBG!,f_auto,q_auto:good,fl_progressive:steep/https%3A%2F%2Fsubstack-post-media.s3.amazonaws.com%2Fpublic%2Fimages%2Fff244cb7-ae1d-41ee-8baa-4ad7d353a204_1765x2606.png)

  

When the client sends data, the NIC fires an interrupt. The kernel places the bytes in `rcv_buff` of `fd=10`, doing TCP reassembly along the way, ordering out-of-order segments, discarding duplicates. The kernel marks `fd=10` as ready in epoll. The netpoller’s `epoll_wait` returns `fd=10`. The netpoller finds `g_conn_1` in `pollDesc[fd=10].rg` and puts it back on the run queue. `g_conn_1` resumes, `conn.Read(buf)` drains from `rcv_buff` into `buf`, and the handler processes the request.

  

* * *

  

## The Client Side — EINPROGRESS

  

Building the client taught me the mirror image of all of this.

  

[![](https://substackcdn.com/image/fetch/$s_!ZAuS!,w_2400,c_limit,f_auto,q_auto:good,fl_progressive:steep/https%3A%2F%2Fsubstack-post-media.s3.amazonaws.com%2Fpublic%2Fimages%2F2fd0ec6d-b05a-451f-a2fa-988486a4661c_1448x1738.png)](https://substackcdn.com/image/fetch/$s_!ZAuS!,f_auto,q_auto:good,fl_progressive:steep/https%3A%2F%2Fsubstack-post-media.s3.amazonaws.com%2Fpublic%2Fimages%2F2fd0ec6d-b05a-451f-a2fa-988486a4661c_1448x1738.png)

  

`net.Dial("tcp", ":4040")` calls `socket()` then `connect()`. No `bind()` call, the kernel assigns an ephemeral port automatically from the range 32768–60999\. So the client’s socket gets `local: :54321, remote: :4040` without our code ever choosing the port number.

  

But here’s where it diverges from the server side. The netpoller registers `fd=3` with epoll for `EPOLLOUT`, not `EPOLLIN`. And the pollDesc entry is:

{fd:3, wg: g_client, rg: nil}

  

`wg` not `rg`. Write goroutine, not read goroutine.

  

`connect()` on a non-blocking socket returns immediately with `EINPROGRESS`, handshake started, not done yet. The goroutine is parked.

  

[![](https://substackcdn.com/image/fetch/$s_!VTEa!,w_2400,c_limit,f_auto,q_auto:good,fl_progressive:steep/https%3A%2F%2Fsubstack-post-media.s3.amazonaws.com%2Fpublic%2Fimages%2Fd5db89b7-4ac5-4550-8b7d-aabd569a1067_1732x2984.png)](https://substackcdn.com/image/fetch/$s_!VTEa!,f_auto,q_auto:good,fl_progressive:steep/https%3A%2F%2Fsubstack-post-media.s3.amazonaws.com%2Fpublic%2Fimages%2Fd5db89b7-4ac5-4550-8b7d-aabd569a1067_1732x2984.png)

  

When the SYN-ACK arrives, the kernel completes the handshake and marks `fd=3` as writable. For a connecting socket, writable means connected. The send buffer is ready, which only becomes true once the connection is established. The kernel reuses “writable” as the signal for “you are now connected.”

  

The netpoller sees `fd=3` is ready via EPOLLOUT, finds `g_client` in `pollDesc[fd=3].wg`, puts it back on the run queue. `g_client` resumes at `connect()`, Go checks `getsockopt(SO_ERROR)` to confirm the connection actually succeeded, and `net.Dial` returns `conn (fd=3)`.

  

The client is now ESTABLISHED. The first `conn.Write([]byte)` puts bytes into `send_buff`. The kernel drains `send_buff` onto the wire.

  

* * *

  

## Building a Raw TCP Echo Server and Client

  

To understand teardown properly, I needed both sides of the connection. So I built a raw TCP echo server and client, no `net/http`, no frameworks, just `net.Listen`, `net.Dial`, and `net.Conn` directly.

  

The server spawns a goroutine per connection, reads lines via `bufio.Scanner`, and echoes each one back:

go func(conn net.Conn) {

defer conn.Close()

sc := bufio.NewScanner(conn)

who := conn.RemoteAddr()

for sc.Scan() {

t := sc.Text()

fmt.Fprintf(conn, "echo: %s\n", t)

}

// ...

}(conn)

  

`defer conn.Close()` is important. Without it, the goroutine exits but the kernel-side socket stays open. The client’s reader blocks forever waiting for data that never comes. The fd leaks. Over time fds accumulate, the process hits its fd limit, and `Accept()` starts failing. `defer` covers all exit paths cleanly.

  

The client reads from stdin and writes to the connection, and simultaneously reads echoes back from the server and prints them to stdout.

  

That “simultaneously” is the problem. Two directions, one connection.

  

Run them sequentially and the first blocks forever.

done := make(chan any, 1)

// writer: stdin → conn

go func() {

io.Copy(conn, WrappedReader{r: os.Stdin})

conn.(*net.TCPConn).CloseWrite()

}()

// reader: conn → stdout

go func() {

io.Copy(os.Stdout, conn)

done <- 1

}()

<-done

conn.Close()

  

The reader signals `done` when it finishes. Main waits on `done`.

  

Why the reader and not the writer? The reader finishing means the server closed its side, nothing more is coming. The writer might finish, user typed the sentinel, while the server is still sending the last echo back. Waiting on the reader catches the true end of the conversation.

  

Here’s what it looks like running:

  

**Client:**

go run client/main.go

connected to server

hi

echo: hi

by

echo: by

bye

echo: bye

EOF

  

**Server:**

go run server/main.go

[127.0.0.1:50472]: connected

[127.0.0.1:50472]: hi

[127.0.0.1:50472]: by

[127.0.0.1:50472]: bye

[127.0.0.1:50472]: disconnected

  

Each line the client types comes back prefixed with `echo:`. When I type `EOF`, the `WrappedReader` returns `io.EOF`, the writer goroutine calls `CloseWrite()`, the server sees the connection close, logs “disconnected”, and shuts down its side. The full teardown plays out exactly as the diagrams show.

  

* * *

  

## CloseWrite Is Not Close

  

When the writer goroutine finishes, I don’t want to close the whole connection. The reader is still waiting for the server’s final echoes.

  

`conn.Close()` would close both directions immediately. The reader goroutine would break mid-flight.

  

`conn.(*net.TCPConn).CloseWrite()` closes only the write side. Underneath it calls `shutdown(fd, SHUT_WR)`. The kernel sends a FIN to the server signalling “I’m done sending”, but `fd=3` stays open for reading. The pollDesc entry stays registered with epoll. The `g_client` reader goroutine keeps running.

  

I had to look this up to understand the type assertion. `net.Conn` is an interface. It doesn’t expose `CloseWrite()`. We have to assert to `*net.TCPConn` first:

conn.(*net.TCPConn).CloseWrite()

  

A half-close. Each direction closed independently, when each side was ready.

  

* * *

  

## The Sentinel Protocol and Its Limits

  

To signal end-of-input from the terminal without Ctrl+D, I built a `WrappedReader` around `os.Stdin`:

type WrappedReader struct {

r io.Reader

}

func (wr WrappedReader) Read(p []byte) (int, error) {

n, err := wr.r.Read(p)

s := string(p[:n])

if s == "EOF\n" {

return 0, io.EOF

}

return n, err

}

  

If it reads `"EOF\n"`, it returns `io.EOF` to the caller instead of passing the bytes through.

  

Three things I got wrong building this before getting it right.

  

`string(p)` not `string(p[:n])`. `p` is a 32KB buffer. Only `n` bytes are valid. The rest is garbage. Reading the full `p` gives you junk after the actual input.

  

Terminal stdin is line-buffered. You get `"EOF\n"` as one chunk, not byte by byte. So the comparison works cleanly.

  

`io.Copy` calls `Read` in a loop. Returning `io.EOF` is the signal to stop. Not an error, just done.

  

It worked. And it immediately showed me its own problem.

  

What if the user actually wants to send the literal string `"EOF"`? The protocol misinterprets it as a control signal. What if the payload is binary, a JPEG, a serialised struct? Any byte sequence could accidentally match a text sentinel.

  

This is the problem real protocols solve with framing.

  

HTTP uses `Content-Length` or `Transfer-Encoding: chunked`. Redis’s RESP uses type-tagged length-prefixed messages, `$6\r\nfoobar\r\n`. gRPC uses binary length-prefixed frames over HTTP/2.

  

The sentinel approach breaks the moment the payload itself contains the sentinel string. Framing works on everything.

  

* * *

  

## What Happens to the Resources

  

The teardown is where I saw the full picture of what a connection actually costs.

  

[![](https://substackcdn.com/image/fetch/$s_!0gZl!,w_2400,c_limit,f_auto,q_auto:good,fl_progressive:steep/https%3A%2F%2Fsubstack-post-media.s3.amazonaws.com%2Fpublic%2Fimages%2F3f36c674-893f-4bc6-811d-6f8fd961f77a_3633x2126.png)](https://substackcdn.com/image/fetch/$s_!0gZl!,f_auto,q_auto:good,fl_progressive:steep/https%3A%2F%2Fsubstack-post-media.s3.amazonaws.com%2Fpublic%2Fimages%2F3f36c674-893f-4bc6-811d-6f8fd961f77a_3633x2126.png)

  

When the writer goroutine finishes and `CloseWrite()` fires, `shutdown(fd=3, SHUT_WR)` is called. The kernel flushes and closes the send buffer. The FIN packet goes out. The socket transitions to `FIN_WAIT_1`. But `fd=3` is still open. The `rcv_buff` is still alive. The reader goroutine is still parked, waiting.

  

`fd=3` is still registered with epoll. Still in the interest list. The netpoller is still watching it.

  

On the server side, the FIN arrives at the NIC. The kernel places an EOF marker into `rcv_buff` of `fd=10`. The kernel marks `fd=10` as ready in epoll, EOF counts as readable.

  

The netpoller unparks `g_conn_1`. `conn.Read` returns the EOF. `sc.Scan()` returns false. The for loop exits.

  

[![](https://substackcdn.com/image/fetch/$s_!H07A!,w_2400,c_limit,f_auto,q_auto:good,fl_progressive:steep/https%3A%2F%2Fsubstack-post-media.s3.amazonaws.com%2Fpublic%2Fimages%2F767c36e5-9b70-4aa0-b447-e7a6fe72fad6_3591x2725.png)](https://substackcdn.com/image/fetch/$s_!H07A!,f_auto,q_auto:good,fl_progressive:steep/https%3A%2F%2Fsubstack-post-media.s3.amazonaws.com%2Fpublic%2Fimages%2F767c36e5-9b70-4aa0-b447-e7a6fe72fad6_3591x2725.png)

  

`defer conn.Close()` fires on the server. Three things happen in order.

  

* `epoll_ctl(DEL, fd=10)` removes `fd=10` from the epoll interest list.

  
  
  

* The pollDesc entry for `g_conn_1` is cleared.

  
  
  

* The netpoller is no longer watching this fd.

  
  
  
  

`close(fd=10)` syscall. The kernel decrements the reference count on the socket. If it hits zero, `fd=10` is freed from the server’s process fd table, and the kernel sends a FIN to the client.

  

The server socket transitions to `LAST_ACK`, waiting for the client’s final ACK.

  

[![](https://substackcdn.com/image/fetch/$s_!C5Ny!,w_2400,c_limit,f_auto,q_auto:good,fl_progressive:steep/https%3A%2F%2Fsubstack-post-media.s3.amazonaws.com%2Fpublic%2Fimages%2F881b2ead-545f-4d02-b752-30e355186494_3418x3642.png)](https://substackcdn.com/image/fetch/$s_!C5Ny!,f_auto,q_auto:good,fl_progressive:steep/https%3A%2F%2Fsubstack-post-media.s3.amazonaws.com%2Fpublic%2Fimages%2F881b2ead-545f-4d02-b752-30e355186494_3418x3642.png)

  

The client’s reader goroutine sees the server’s FIN — the kernel placed an EOF marker in `rcv_buff` of `fd=3`, epoll fired, the netpoller unparked `g_client`. `io.Copy(os.Stdout, conn)` returns. `done <- 1`. Main unblocks. `conn.Close()` is called on the client side.

  

On the client, `conn.Close()` triggers its own three-step sequence. `epoll_ctl(DEL, fd=3)` removes `fd=3` from the epoll interest list. The pollDesc entry for `g_client` is cleared. `close(fd=3)` frees `fd=3` from the client’s process fd table — the Go process can reuse that integer for a new connection immediately. As part of `close(fd=3)`, the kernel sends the final ACK to the server and transitions the client socket to `TIME_WAIT`.

  

The server receives the final ACK. `cli:tcp_sock` transitions to `CLOSED`. The kernel frees everything on the server side — the `tcp_sock` structure, `rcv_buff`, `send_buff`, all TCP state — sequence numbers, timers, everything.

  

On the client side, the socket lingers in `TIME_WAIT` for up to 120 seconds (2 * MSL, Maximum Segment Lifetime). The `rcv_buff` stays alive through this period — the kernel needs it in case the server retransmits its FIN. Only after 2*MSL expires does the kernel free the client’s `tcp_sock`, `rcv_buff`, and all remaining state.

  

Two reasons TIME_WAIT exists. If the final ACK got lost, the server will retransmit its FIN — the client needs to still be alive to respond. And if the same 4-tuple gets reused for a new connection too quickly, a stale packet from the old connection could arrive and corrupt it. `TIME_WAIT` prevents both.

  

The netpoller’s involvement ends at `conn.Close()`. Everything after — `LAST_ACK`, `TIME_WAIT`, final memory free — is the kernel’s TCP state machine running on its own.

  

* * *

  

## The Resource Lifecycle, in Full

  

[![](https://substackcdn.com/image/fetch/$s_!buIB!,w_1456,c_limit,f_auto,q_auto:good,fl_progressive:steep/https%3A%2F%2Fsubstack-post-media.s3.amazonaws.com%2Fpublic%2Fimages%2Fb944e6ed-52ce-4e56-b263-8617a1a44806_1180x735.png)](https://substackcdn.com/image/fetch/$s_!buIB!,f_auto,q_auto:good,fl_progressive:steep/https%3A%2F%2Fsubstack-post-media.s3.amazonaws.com%2Fpublic%2Fimages%2Fb944e6ed-52ce-4e56-b263-8617a1a44806_1180x735.png)

  

The fd dies at `conn.Close()`. The pollDesc dies with it. On the server side, `conn.Close()` sends the FIN, the client’s final ACK arrives, and the kernel frees everything — `tcp_sock`, both buffers, all TCP state — when it reaches CLOSED.

  

On the client side it’s different. `conn.Close()` fires the final ACK as part of `close(fd=3)`, and the socket transitions to TIME_WAIT. But the `rcv_buff` stays alive through TIME_WAIT — the kernel needs it in case the server retransmits its FIN. The client has to be able to receive and respond. Only after 2*MSL expires does the kernel free the `tcp_sock`, `rcv_buff`, and all remaining state.

  

The fd and pollDesc die immediately at `conn.Close()`. The kernel socket outlives them both.

  

Every row in that table is real kernel work. Every connection teardown goes through it.

  

* * *

  

## This Is Why Connection Pools Exist

  

Look at that table again. Every connection teardown fires `epoll_ctl(DEL)`, clears the pollDesc, closes the fd, sends a FIN, waits for ACK, frees the socket structure, frees both buffers, waits out TIME_WAIT. And that’s just the teardown. The setup side is just as expensive — `socket()`, `bind()`, `connect()`, the three-way handshake, buffer allocation, pollDesc registration.

  

Every new connection pays that cost twice. Once to set up, once to tear down.

  

In the teardown diagrams above, I marked every resource cleanup event in green — `epoll_ctl(DEL, fd=10)`, the pollDesc entry being cleared, the fd freed from the process table, the final kernel free of `tcp_sock`, `rcv_buff`, `send_buff`, and all TCP state. Each green box is real kernel work that happens on every single connection close.

  

Connection pooling is the realisation that most of that work is unnecessary if the connection can be reused. Instead of tearing a connection down after each request, a pool keeps it alive and hands it back for the next one. None of those green boxes fire. The `tcp_sock` stays in kernel memory. The buffers stay allocated. The `pollDesc` entry stays registered with epoll. The next request skips the handshake entirely and goes straight to writing bytes into an already-established connection.

  

The netpoller makes Go capable of holding thousands of those alive connections open simultaneously without wasting OS threads. But the connections themselves still cost kernel memory, fd table entries, and buffer allocations. That cost doesn’t disappear just because goroutines are cheap.

  

Go’s netpoller solves the thread problem. Connection pools solve the connection cost problem. They operate at different layers and both matter.

  

Pool them. Reuse them. The teardown is expensive.

  

* * *

  

## What’s Still Open

  

I now understand what happens when you call `net.Listen` and `net.Dial`. The kernel owns more of the process than I realised. The handshake, the queues, the buffers, the teardown, most of it runs without our code ever being called.

  

Go’s contribution is making the interface to all of that feel sequential. `Accept()` looks blocking. `Read()` looks blocking. Neither holds a thread while waiting. The netpoller and epoll handle the gap between “nothing ready” and “resume here.”

  

What I still don’t understand well enough: what happens when thousands of connections are open simultaneously and I don’t want to pay the handshake cost every time. How connection pools work. How HTTP’s request-response framing sits on top of exactly this primitive. How the send and receive buffers interact with TCP’s flow control when one side is slow.

  

That’s where this goes next.

  

* * *

  

_The diagrams in this article are from my Excalidraw canvas where I mapped the full TCP connection lifecycle while building this. The full canvas is in the[repo here](https://github.com/vishal2098govind/http-deep-dive/blob/main/cmd/client-server/tcp-connection-deep-dive.excalidraw)._