package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
)

func main() {
	resp, err := http.Get("https://api.github.com/users/vishal2098govind")
	if err != nil {
		fmt.Printf("err: %v", err)
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		fmt.Printf("status: %d\n", resp.StatusCode)
		return
	}
	io.Copy(os.Stdout, resp.Body)
}
