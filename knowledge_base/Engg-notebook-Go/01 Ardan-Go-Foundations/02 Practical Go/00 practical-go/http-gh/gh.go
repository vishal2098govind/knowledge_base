package main

import (
	"encoding/json"
	"fmt"
	"net/http"
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

	var reply struct {
		Name         string
		Public_Repos int
		PublicGists  int32 `json:"public_gists"`
	}

	dec := json.NewDecoder(resp.Body)
	dec.Decode(&reply)
	fmt.Println(reply.Name, reply.Public_Repos, reply.PublicGists)
}
