#go #concurrency #go-context #go-channels

### `package crawler`
```go
package crawler

import (
	"context"
	"net/http"

	"golang.org/x/net/html"
)

const (
	MAX_DEPTH = 100
)

type Link struct {
	url   string
	depth int
}

func Crawl(ctx context.Context, url string) ([]Link, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}

	res, err := http.DefaultClient.Do(req)
	// res, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	doc, err := html.Parse(res.Body)
	if err != nil {
		return nil, err
	}
	links := bfsCrawl(doc, res, 0)

	return links, nil
}

func bfsCrawl(node *html.Node, res *http.Response, depth int) []Link {
	if depth > MAX_DEPTH {
		return []Link{}
	}

	links := []Link{}
	links = append(links, extract(node, res, depth)...)

	for c := node.FirstChild; c != nil; c = c.NextSibling {
		links = append(links, bfsCrawl(c, res, depth+1)...)
	}
	return links
}

func extract(n *html.Node, res *http.Response, depth int) []Link {
	links := []Link{}
	if n.Type == html.ElementNode && n.Data == "a" {
		for _, attr := range n.Attr {
			if attr.Key != "href" {
				continue
			}
			l, err := res.Request.URL.Parse(attr.Val)
			if err != nil {
				continue
			}
			links = append(links, Link{l.String(), depth})
		}
	}
	return links
}

// func main() {
// 	el, _ := Crawl("http://gopl.io/")
// 	// ls := []string{}
// 	for _, v := range el {
// 		// ls = append(ls, v.url)
// 		fmt.Println(v)
// 	}

// 	// fmt.Println("--------------------")

// 	// l, _ := links.Extract("http://gopl.io/")

// 	// for i, v := range l {
// 	// 	if i < len(ls) {
// 	// 		fmt.Println(v == ls[i])
// 	// 	}
// 	// 	fmt.Println(v)
// 	// }

// }

```

### `package main`
```go
package main

import (
	"context"
	"fmt"
	"os"
	"practical-go/crawler"
	"sync"
	"time"
)

func main() {
	urls := []string{
		"http://gopl.io",
		"http://gopl.io",
		"http://gopl.io",
		"https://golang.org/help",
		"https://golang.org/help",
		"https://golang.org/help",
		"https://golang.org/help",
		"https://golang.org/help",
		"https://golang.org/help",
		"https://golang.org/help",
	}
	urls = append(urls, urls...)
	urls = append(urls, urls...)
	urls = append(urls, urls...)
	urls = append(urls, urls...)

	perfSerial := len(os.Args) == 2 && os.Args[1] == "-s"

	start := time.Now()
	if perfSerial {
		serial(urls)
	} else {
		parallel(urls)
	}
	fmt.Println("took ", time.Since(start))
}

func parallel(urls []string) {
	linksCh := make(chan []crawler.Link)
	urlCh := make(chan string)
	var wg sync.WaitGroup
	// wg.Add(len(urls))

	ctx, cancel := context.WithCancel(context.Background())
	// defer cancel()

	go func() {
		os.Stdin.Read(make([]byte, 1))
		cancel()
	}()

	sem := make(chan struct{}, 96) // for not allowing more than 20 go-routines to crawl at a time

	// pooling
	g := 100

	for i := range g {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()

			for {
				select {
				case <-ctx.Done():
					return
				case v, ok := <-urlCh:
					if !ok {
						return
					}

					fmt.Printf("[gr-%d] crawling [url=%s]\n", i, v)
					sem <- struct{}{}
					func() {
						defer func() { <-sem }()
						l, _ := crawler.Crawl(ctx, v)
						select {
						case <-ctx.Done():
						case linksCh <- l:
						}
					}()

				}
			}

			// for v := range urlCh {
			// 	// select {
			// 	// case <-ctx.Done():
			// 	// 	return
			// 	// default:
			// 	// 	// continue
			// 	// }

			// 	fmt.Printf("[gr-%d] crawling [url=%s]\n", i, v)
			// 	sem <- struct{}{}
			// 	func() {
			// 		defer func() { <-sem }()
			// 		l, _ := crawler.Crawl(ctx, v)
			// 		select {
			// 		case <-ctx.Done():
			// 		case linksCh <- l:
			// 		}
			// 	}()
			// 	wg.Done()
			// }
		}(i)
	}

	go func() {
	loop:
		for _, v := range urls {
			select {
			case <-ctx.Done():
				break loop
			case urlCh <- v:
				// wg.Add(1)
			}
		}
		close(urlCh)
	}()

	go func() {
		wg.Wait()
		close(linksCh)
	}()

	count := 0
	for v := range linksCh {
		count += len(v)
		fmt.Printf("received %d more from linksCh \n", len(v))
		// for _, l := range v {
		// 	fmt.Println(l)
		// }
	}
	fmt.Println("len(links)=", count)
}

func serial(urls []string) {
	links := []crawler.Link{}

	for _, v := range urls {
		l, _ := crawler.Crawl(context.Background(), v)
		links = append(links, l...)
	}

	for _, v := range links {
		fmt.Println(v)
	}
	fmt.Println("len(links)=", len(links))
}

```