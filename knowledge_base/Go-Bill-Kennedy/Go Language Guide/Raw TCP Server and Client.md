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