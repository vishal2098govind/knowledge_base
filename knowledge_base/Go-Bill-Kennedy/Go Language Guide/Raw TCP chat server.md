#go #tcp #tcp-server #concurrency #go-channels #go-context

#### Simply accept a raw tcp connection and exit server
```go
package main

import (
	"fmt"
	"net"
	"bufio"
)

func main() {
	lis, err := net.Listen("tcp", ":8080")
	if err != nil {
		fmt.Printf("failed to listen for tcp traffic: %s\n", err)
		return
	}
	fmt.Println("listening for tcp connections")
	
	conn, err := lis.Accept() // blocking call until someone requests to connect
	if err != nil {
		fmt.Printf("failed to accept incoming connection request: %s\n", err)
		return
	}
	fmt.Printf("[%s] connected\n", conn.RemoteAddr())
}
```
##### Output:
```bash
// OUTPUT
# terminal-1                           |  # terminal-2
$ go run .                             |
listening for tcp connections          |
                                       |  $ nc localhost 8080
[[::1]:58656] connected                |
$                                      |  $
```

#### Simply accept a raw tcp connection and exit server after closing first connection
```go
package main

import (
	"bufio"
	"fmt"
	"net"
)

func main() {
	lis, err := net.Listen("tcp", ":8080")
	if err != nil {
		fmt.Printf("failed to listen for tcp traffic: %s\n", err)
		return
	}
	fmt.Println("listening for tcp connections")

	conn, err := lis.Accept() // blocking call until someone requests to connect
	if err != nil {
		fmt.Printf("failed to accept incoming connection request: %s\n", err)
		return
	}
	fmt.Printf("[%s] connected\n", conn.RemoteAddr())

	sc := bufio.NewScanner(conn)
	for sc.Scan() {
		fmt.Println(sc.Text())
	}
	err = sc.Err()
	if err != nil {
		fmt.Printf("error reading from connection: %s\n", err)
		return
	}
	
	fmt.Println("client disconnected")
}
```
##### Output
```bash
# terminal-1                           |# terminal-2
$ go run .                             |
listening for tcp connections          |
                                       |$ nc localhost 8080
[[::1]:58656] connected                |
                                       |hi
hi                                     |
                                       |hello
hello                                  |
                                       |I love smruti
I love smruti                          |
                                       |ok, bye
ok, bye                                |
                                       |^C
client disconnected                    |$
$                                      |
```

#### Accept multiple raw tcp connections serially
- sequentially means accepting a new connection only after closing current connection
```go
package main

import (
	"bufio"
	"fmt"
	"net"
)

func main() {
	lis, err := net.Listen("tcp", ":8080")
	if err != nil {
		fmt.Printf("failed to listen for tcp traffic: %s\n", err)
		return
	}
	fmt.Println("listening for tcp connections")

	for {
		fmt.Println("accepting connectionss..")
		conn, err := lis.Accept() // blocking call until someone requests to connect
		if err != nil {
			fmt.Printf("failed to accept incoming connection request: %s\n", err)
			return
		}
		fmt.Printf("[%s] connected\n", conn.RemoteAddr())

		sc := bufio.NewScanner(conn)

		for sc.Scan() {
			fmt.Println(sc.Text())
		}
		err = sc.Err()
		if err != nil {
			fmt.Printf("error reading from connection: %s\n", err)
			return
		}
		fmt.Println("client disconnected")
	}
}

```

```bash
# terminal-1 (server)                  |# terminal-2 (58656)          |# terminal-3 (59089)
$ go run .                             |                              |
listening for tcp connections          |                              |
                                       |$ nc localhost 8080           |
[[::1]:58656] connected                |                              |
accepting connectionss..               |                              |
                                       |                              |$ nc localhost 8080
                                       |hi from t-2                   |
hi from t-2                            |                              |
                                       |bye from t-2                  |
bye from t-2                           |                              |
                                       |^C                            |
client disconnected                    |$                             |
accepting connectionss..               |                              |
[[::1]:59089] connected                |                              |
                                       |                              |hi from t-3
hi from t-3                            |                              |
                                       |                              |bye from t-3
bye from t-2                           |                              |
                                       |                              |^C
client disconnected                    |                              |
accepting connectionss..               |                              |
```

#### Accept multiple raw tcp connections concurrently
```go
package main

import (
	"bufio"
	"fmt"
	"net"
)

func main() {
	lis, err := net.Listen("tcp", ":8080")
	if err != nil {
		fmt.Printf("failed to listen for tcp traffic: %s\n", err)
		return
	}
	fmt.Println("listening for tcp connections")

	for {
		fmt.Println("accepting connectionss..")
		conn, err := lis.Accept() // blocking call until someone requests to connect
		if err != nil {
			fmt.Printf("failed to accept incoming connection request: %s\n", err)
			return
		}
		fmt.Printf("[%s] connected\n", conn.RemoteAddr())

		go func(conn net.Conn) {
			sc := bufio.NewScanner(conn)
			who := conn.RemoteAddr()

			for sc.Scan() {
				fmt.Printf("[%s]: %s\n", who, sc.Text())
			}
			err = sc.Err()
			if err != nil {
				fmt.Printf("error reading from connection: %s\n", err)
				return
			}
			fmt.Printf("[%s]: client disconnected\n", who)
		}(conn)
	}

}

```
##### Output
```bash
# terminal-1 (server)                   |# terminal-2 (59490)         |# terminal-2 (59491) 
$ go run .                              |                             |
listening for tcp connections           |                             |
accepting connectionss..                |                             |
                                        |$ nc localhost 8080          |
[[::1]:59490] connected                 |                             |
accepting connectionss..                |                             |
                                        |                             |$ nc localhost 8080  
[[::1]:59491] connected                 |                             |
accepting connectionss..                |                             |
                                        |hi from t-2                  |
[[::1]:59490]: hi from t-2              |                             |
                                        |                             |hi from t-3
[[::1]:59491]: hi from t-3              |                             |
                                        |bye from t-2                 |
[[::1]:59490]: bye from t-2             |                             |
                                        |^C                           |
[[::1]:59490]: client disconnected      |                             |
                                        |                             |bye from t-3
[[::1]:59491]: bye from t-3             |                             |
                                        |                             |^C
[[::1]:59491]: client disconnected      |                             |
```

#### Broadcast chat server
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
			fmt.Fprintln(w, msg)
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

	messages <- fmt.Sprintf("[%s] left the chat", who)
	leaves <- ch
	fmt.Printf("[%s]: client disconnected\n", who)
}
```
#### Output
```bash
# terminal-1 (server)                   |# terminal-2 (:60717)        |# terminal-2 (:60719) 
$ go run .                              |                             |
listening for tcp connections           |                             |
accepting connectionss..                |                             |
                                        |$ nc localhost 8080          |
[[::1]:60717] connected                 |                             |
                                        |[[::1]:60717] joined the chat|
accepting connectionss..                |                             |
[[::1]:60717] joined the chat           |                             |
                                        |                             |$ nc localhost 8080
[[::1]:60719] connected                 |                             |
accepting connectionss..                |                             |
[[::1]:60719] joined the chat           |                             |
                                        |[[::1]:60719] joined the chat|[[::1]:60719] joined the chat
                                        |hi from t-1                  |
[[::1]:60717]: hi from t-1              |                             |
                                        |[[::1]:60717] hi from t-1    |
                                        |                             |[[::1]:60717]: hi from t-1
                                        |                             |hi from t-2
[[::1]:60719]: hi from t-2              |                             |
                                        |[[::1]:60719] hi from t-2    |[[::1]:60719]: hi from t-2
                                        |bye from t-1                 |
[[::1]:60717]: bye from t-1             |                             |
                                        |[[::1]:60717] bye from t-1   |[[::1]:60717]: bye from t-1
                                        |^C                           |
[[::1]:60717] left the chat             |                             |
                                        |                             |[[::1]:60717] left the chat
[[::1]:60717]: client disconnected      |                             |
                                        |                             |^C
[[::1]:60719] left the chat             |                             |
[[::1]:60719]: client disconnected      |                             |
```