#go #tcp #net #net-Conn 
### TCP Server
```go
package main

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"net"
)

type client chan<- string

var (
	entries  = make(chan client) // to signal when a client gets connected
	leaves   = make(chan client) // to signal when a client disconnects
	messages = make(chan string)
)

func main() {
	lis, err := net.Listen("tcp", ":8080")
	if err != nil {
		fmt.Printf("failed to listen for tcp traffic: %s\n", err)
		return
	}
	fmt.Println("listening for tcp connections")

	ctx := context.Background()

	go broadcaster(ctx)

	for {
		fmt.Println("accepting connectionss..")
		conn, err := lis.Accept() // blocking call until someone requests to connect
		if err != nil {
			fmt.Printf("failed to accept incoming connection request: %s\n", err)
			return
		}
		fmt.Printf("[%s] connected\n", conn.RemoteAddr())

		go handleConn(conn)
	}
}

func broadcaster(ctx context.Context) {
	clients := make(map[client]bool)

	for {
		select {
		case <-ctx.Done():
			for cli := range clients {
				close(cli)
			}
			return
		case msg := <-messages:
			fmt.Println(msg)
			for ch := range clients {
				ch <- msg
			}

		case cli := <-entries:
			clients[cli] = true

		case cli := <-leaves:
			delete(clients, cli)
			close(cli)
		}
	}
}

func handleConn(conn net.Conn) {
	defer conn.Close()

	ch := make(chan string)
	sc := bufio.NewScanner(conn)
	who := conn.RemoteAddr()

	entries <- ch

	go func(w io.Writer, ch <-chan string) {
		for msg := range ch {
			_, err := fmt.Fprintln(w, msg)
			if err != nil {
				fmt.Println(err)
			}
		}
	}(conn, ch)

	messages <- fmt.Sprintf("[%s] joined the chat", who)

	for sc.Scan() {
		msg := fmt.Sprintf("[%s]: %s", who, sc.Text())
		messages <- msg
	}
	err := sc.Err()
	if err != nil {
		fmt.Printf("error reading from connection: %s\n", err)
		return
	}

	// fmt.Fprintln(conn, "thank you for connecting us")

	messages <- fmt.Sprintf("[%s] left the chat", who)
	leaves <- ch
	fmt.Printf("[%s]: client disconnected\n", who)
}
```

### TCP Client (`nc` ~ `netcat`)
```go
package main

import (
	"fmt"
	"io"
	"net"
	"os"
)

func main() {
	conn, err := net.Dial("tcp", ":8080")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer func() {
		fmt.Println("closing connection")
		conn.Close()
	}()

	go func() {
		io.Copy(conn, os.Stdin)
		fmt.Println("copy done")
		conn.(*net.TCPConn).CloseWrite()
	}()

	io.Copy(os.Stdout, conn)
}
```

## ChatGPT chat around this:
### Question: trying to implement netcat
help me understand logical difference (ignoring error handling and all) between the following two tcp client implementations: the first one is taken from the gopl.io netcat

```go
package main

import (
	"io"
	"log"
	"net"
	"os"
)

func main() {
	conn, err := net.Dial("tcp", "localhost:8000")
	if err != nil {
		log.Fatal(err)
	}
	done := make(chan struct{})
	go func() {
		io.Copy(os.Stdout, conn) // NOTE: ignoring errors
		log.Println("done")
		done <- struct{}{} // signal the main goroutine
	}()
	mustCopy(conn, os.Stdin)
	conn.Close()
	<-done // wait for background goroutine to finish
}

func mustCopy(dst io.Writer, src io.Reader) {
	if _, err := io.Copy(dst, src); err != nil {
		log.Fatal(err)
	}
}
```

the second one is my implementation:

```go
package main

import (
	"fmt"
	"io"
	"net"
	"os"
)

func main() {
	conn, err := net.Dial("tcp", ":8080")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer func() {
		fmt.Println("closing connection")
		conn.Close()
	}()

	go func() {
		for {
			io.Copy(conn, os.Stdout)
		}
	}()

	for {
		io.Copy(os.Stdout, conn)
	}
}
```
### Answer:

They _look_ similar (“two goroutines copying in opposite directions”), but logically they’re doing very different things.

##### 1) Direction of data flow (the big bug)

###### GoPL netcat

- **Background goroutine:** `io.Copy(os.Stdout, conn)`  
    ✅ reads **from conn** (server → client) and writes **to stdout**
    
- **Main goroutine:** `io.Copy(conn, os.Stdin)`  
    ✅ reads **from stdin** (client input) and writes **to conn** (client → server)
    

So it’s a proper bidirectional pipe:

- user input → server
    
- server output → terminal
    

###### Your program

You flipped one of the directions:

```go
go func() {
    for {
        io.Copy(conn, os.Stdout)
    }
}()
```

`io.Copy(conn, os.Stdout)` means:
- read from **stdout**
- write to **conn**
But `os.Stdout` is a _writer_, not a meaningful _reader_ (reading from it doesn’t represent “what you typed”). So this goroutine is conceptually wrong. You probably meant:
- `io.Copy(conn, os.Stdin)` (stdin → conn)

And your other loop:

`for {     io.Copy(os.Stdout, conn) }`

This one is the correct direction (conn → stdout), but see the next issues.

##### 2) Termination behavior (the other huge difference)

###### GoPL netcat has a clean shutdown protocol

- When stdin hits EOF, `io.Copy(conn, os.Stdin)` returns.
- Then it **closes the connection**: `conn.Close()`
- Closing the conn causes the other goroutine’s `io.Copy(os.Stdout, conn)` to return too (read unblocks / returns EOF).
- `done` channel is used so main waits until the background goroutine finishes printing any last bytes.

That means it shuts down deterministically.

###### Your program never stops (and “defer Close” never runs)

You have:

```go
for {
    io.Copy(os.Stdout, conn)
}
```

Even if `io.Copy` returns (EOF, error, etc.), you immediately call it again forever. So:
- main never returns → deferred close never runs
- you keep spinning / blocking pointlessly after the connection is done

Same for the goroutine: you retry forever too.

In other words: GoPL version runs “until input ends, then close and exit”.  
Yours runs “forever, regardless of EOF”.
##### 3) Concurrency / correctness: double-reads and weird states

Your main loop repeatedly starts new `io.Copy(os.Stdout, conn)` operations. Each call tries to read from the same `conn`. You currently only have one such loop in main, so you aren’t _simultaneously_ double-reading—**but** the structure invites accidental double readers if you later add anything else. And after EOF, it becomes a busy loop (copy returns immediately, repeats, etc.).

The GoPL version has exactly **one reader** from conn (the background goroutine) and exactly **one writer** to conn (main goroutine). That’s the intended full-duplex pattern.
##### 4) Practical “what you intended” version (minimal fix)
If you want the same logic as GoPL, your structure should be basically:
```go
done := make(chan struct{})
go func() {
    io.Copy(os.Stdout, conn) // conn -> stdout
    close(done)
}()

io.Copy(conn, os.Stdin) // stdin -> conn
conn.Close()
<-done
```
That’s the _logical_ equivalent.

##### Summary in one line
- **GoPL netcat:** full-duplex + graceful shutdown coordination.
- **Yours:** one direction is backwards (`Stdout` used as a source), and both sides loop forever so the program doesn’t terminate cleanly and can spin after EOF.
If you want, I can also explain _why_ netcat closes the connection only after stdin completes (and what “half-close” with `CloseWrite()` changes for TCP).
### Learning: `io.Copy()` function blocks the go-routine until `EOF` is read by the `io.Reader.Read()`
 no need an in-definite for loop over `io.Copy(io.Writer, io.Reader)`

---
### Question: when is EOF reached (^C vs ^D) in `os.Stdin`
### Answer: EOF is _not_ reached when you press `Ctrl+C`.
##### Short answer (then we’ll go deep)
> **EOF is _not_ reached when you press `Ctrl+C`.**

- **`Ctrl+C`** sends **SIGINT** → your process is killed → no EOF is involved.
- **EOF is reached when the _writer on the other end closes_**:
    - For `os.Stdin`: when the terminal signals end-of-input (`Ctrl+D` on Unix)
    - For `conn`: when the **remote peer closes its write side** of the TCP connection
##### Let’s walk through your program _timeline-style_

```go
go func() {
    io.Copy(conn, os.Stdin)
    cancel()
}()

io.Copy(os.Stdout, conn)
<-ctx.Done()
```
##### Goroutine 1: `io.Copy(conn, os.Stdin)`

**Reads from:** `os.Stdin`  
**Writes to:** TCP socket

This `io.Copy` blocks in a loop calling:

`read(os.Stdin) → write(conn)`

###### When does it return?

Only when **`os.Stdin.Read()` returns EOF**.

###### What causes EOF on `os.Stdin`?

- On Unix/macOS:
    - Press **`Ctrl+D`** on an empty line
- Or:
    - stdin is piped from a file and the file ends  
- Or:
    - stdin is closed by the parent process  

📌 **Not Ctrl+C**

##### Goroutine 2 (main): `io.Copy(os.Stdout, conn)`

**Reads from:** TCP socket  
**Writes to:** stdout

This blocks in:

`read(conn) → write(os.Stdout)`

###### When does _this_ return?

When **the server closes its write side** of the connection.
That happens if:
- server calls `conn.Close()`
- server calls `CloseWrite()` (half-close)
- server process exits/crashes

##### What happens on Ctrl+C?

|Event|What actually happens|
|---|---|
|Press Ctrl+C|Kernel sends SIGINT|
|Go runtime|Default behavior → terminate process|
|`defer conn.Close()`|❌ **Does NOT run**|
|`cancel()`|❌ irrelevant|
|EOF|❌ never happens|

So:

> **Ctrl+C ≠ EOF**  
> It’s a **signal**, not a stream condition.


### Learning: ^D is used to indicate EOF in `os.Stdin`, no need to indefinitely loop over `io.Copy`

---
### Question: Is there anything missing to be equivalent to `nc`
```go
package main

import (
	"context"
	"fmt"
	"io"
	"net"
	"os"
)

func main() {
	conn, err := net.Dial("tcp", ":8080")
	ctx, cancel := context.WithCancel(context.Background())
	if err != nil {
		fmt.Println(err)
		return
	}
	defer func() {
		fmt.Println("closing connection")
		conn.Close()
	}()

	go func() {
		io.Copy(conn, os.Stdin)
		fmt.Println("copy done")
		cancel()
	}()

	io.Copy(os.Stdout, conn)
	<-ctx.Done()
}
```
### Answer: “Yes in spirit, but not in TCP semantics.”
##### What `nc` / GoPL netcat is _actually_ doing
Conceptually, `nc` does **full-duplex streaming with half-close**:
- stdin → socket **until stdin EOF**
- socket → stdout **until peer closes**
- when stdin finishes, **signal the peer** that no more data will be sent
- still keep reading until the peer is done

This is the key detail.
##### Your version vs real `nc`
###### Your code (simplified behavior)
```go
go func() {
    io.Copy(conn, os.Stdin)
    cancel()
}()

io.Copy(os.Stdout, conn)
<-ctx.Done()
```
###### What’s correct
✅ One `io.Copy` per direction  
✅ No manual `for` loops  
✅ Concurrent bidirectional streaming

###### What’s **not** netcat-correct

###### ❌ 1. You never half-close the TCP connection

When stdin reaches EOF:

`io.Copy(conn, os.Stdin) // returns cancel()`

But the server **does not know** you’re done sending.

In TCP terms:
- no FIN is sent
- server may keep waiting for more input
- deadlock is possible
###### What `nc` actually does here
It performs a **half-close**:
```txt
stdin EOF
→ shutdown(SHUT_WR)
→ FIN sent
→ server sees EOF on read
→ server finishes response
→ server closes
→ client read returns EOF
```

In Go this is:

`conn.(*net.TCPConn).CloseWrite()`
###### ❌ 2. `context` doesn’t control anything here

This line:

`<-ctx.Done()`

does **not** coordinate shutdown.

Why?
- `io.Copy(os.Stdout, conn)` blocks first
- context cancellation doesn’t interrupt `io.Copy`
- the context is only a notification, not a control mechanism

The GoPL version instead uses **connection closure** as the synchronization primitive — which is the _correct_ abstraction for TCP.
### Learning: Need to close (atleast half-close or send FIN) TCP connection on reaching EOF in `os.Stdin`

---
### Question: `conn.Close()` vs `conn.(*net.TCPConn).CloseWrite()`
### Answer: It's pure TCP semantics
##### TL;DR (keep this forever)

|Call|TCP meaning|Effect|
|---|---|---|
|`conn.Close()`|**Full close**|Sends **FIN** (write side) + closes read side|
|`CloseWrite()`|**Half-close (shutdown WR)**|Sends **FIN** for **write only**, read stays open|

In other words:

> **TCP is full-duplex.  
> You can stop sending while still receiving.**

##### TCP background (must understand)

A TCP connection has **two independent streams**:

`Client → Server Server → Client`

Each direction can be closed **independently**.
##### `conn.Close()` — full close

`conn.Close()`

###### What happens

- Your process:
    
    - closes write side
        
    - closes read side
        
- Kernel:
    
    - sends **FIN**
        
    - discards unread inbound data
        
- Future:
    
    - reads & writes fail
        
    - peer eventually sees EOF or RST
        

###### Meaning

> “I’m completely done with this connection.”

###### Effect on `io.Copy`
- Any `io.Copy` reading **from** this conn:
    - unblocks immediately
- Any `io.Copy` writing **to** this conn:
    - errors
###### Typical use
- Server shutdown
- Error paths
- “I don’t care about the response anymore”

##### `CloseWrite()` — half-close (the subtle one)

`tcpConn := conn.(*net.TCPConn) tcpConn.CloseWrite()`

###### What happens

- Your process:    
    - closes **write side only**
    - read side stays open
- Kernel:
    - sends **FIN**
- Peer:
    - sees EOF **when reading**
    - can continue writing back
###### Meaning
> “I’m done sending, but I still want to read.”
###### Effect on `io.Copy`
- `io.Copy(conn, os.Stdin)`:
    - finishes
- `io.Copy(os.Stdout, conn)`:
    - **continues normally**
    - until peer closes
##### Why `nc` _must_ use `CloseWrite`
Consider this interaction:

```
Client sends request
Client is done sending
Server computes response
Server sends response
Server closes
```
If client does:
`conn.Close()`
👉 server may get **RST** or client drops response prematurely.
If client does:

`CloseWrite()`
👉 server:
- sees EOF on request
- finishes work
- sends response
- closes cleanly
##### Visual timeline (this matters)

###### With `CloseWrite()`

```
Client → Server : data
Client → Server : FIN  (CloseWrite)
Server → Client : data
Server → Client : FIN
```

###### With `Close()`

```
Client → Server : data
Client → Server : FIN
Client ← Server : ❌ discarded
```
##### Why `io.Copy` pairs perfectly with `CloseWrite`

```go
go func() {
    io.Copy(conn, os.Stdin)
    tcpConn.CloseWrite()
}()

io.Copy(os.Stdout, conn)
```

- stdin EOF → send FIN
- server sees EOF
- server finishes
- client reads until server FIN
- everything shuts down naturally

This is **textbook full-duplex streaming**.

##### One more subtlety (advanced but important)
If you call:
`conn.Close()`
while another goroutine is blocked on:
`io.Copy(os.Stdout, conn)`
👉 the read **unblocks immediately** (EOF or error)

If you call:
`CloseWrite()`
👉 read **does not unblock**  
👉 only peer closing will unblock it

This distinction is critical for graceful shutdown.
##### Mental model (lock this in)
> **`Close()` kills the socket.  
> `CloseWrite()` politely says “I’m done talking, you go ahead.”**
##### When to use what

| Scenario                     | Correct call                     |
| ---------------------------- | -------------------------------- |
| Netcat / streaming protocols | `CloseWrite()`                   |
| HTTP request body finished   | `CloseWrite()` (done internally) |
| Fatal error                  | `Close()`                        |
| App shutdown                 | `Close()`                        |
| Fire-and-forget              | `Close()`                        |

### Learning: `conn.Close()` is full close, `conn(*net.TCPConn).CloseWrite()` is half-close
Mental model (lock this in)
> **`Close()` kills the socket.  
> `CloseWrite()` politely says “I’m done talking, you go ahead.”**
---