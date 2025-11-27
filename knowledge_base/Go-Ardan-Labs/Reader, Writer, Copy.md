#go #io-Reader #io-Writer 

`io.Copy` efficiently streams data from an `io.Reader` to an `io.Writer` using an internal buffer, copying until **EOF** and returning the total bytes moved and any error, making it the **idiomatic** way to pipe input to output in Go.

```go
package main

import (
	"fmt"
	"io"
	"os"
)

func main() {
	fmt.Println(ReadWrite("io.go"))
}

func ReadWrite(fname string) (string, error) {
	f, err := os.Open(fname)
	if err != nil {
		return "", err
	}
	defer f.Close()

	wcf, err := os.Create(fmt.Sprintf("write_%s", fname))
	if err != nil {
		return "", err
	}
	defer wcf.Close()

	wf, err := os.Create(fmt.Sprintf("copy_%s", fname))
	if err != nil {
		return "", err
	}
	defer wf.Close()

	if _, err := io.Copy(wf, f); err != nil {
		return "", err
	}

	if _, err := f.Seek(0, 0); err != nil {
		return "", err
	}

	buf := make([]byte, 8) // initialize buffer of size 20
	content := []byte{}

	i := 0
	for {
		fmt.Println("---- calling Read() ----")
		fmt.Printf("buffer capacity = %d\n", cap(buf))

		n, err := f.Read(buf) // n has no.of bytes read
		content = append(content, buf[:n]...)
		fmt.Printf("i:%d\n", i)
		// fmt.Printf("content read in this iteration: %s\n", string(content))

		if err != nil {
			if err == io.EOF {
				// file ended
				break
			}
			return "", err
		}

		i++
	}

	if _, err := wcf.Write(content); err != nil {
		return "", err
	}

	return string(content), nil
}

```