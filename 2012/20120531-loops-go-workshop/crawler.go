package main

import (
	"fmt"
)

type Fetcher interface {
	// Fetch returns the body of URL and
	// a slice of URLs found on that page.
	Fetch(url string) (body string, urls []string, err error)
}

type response struct { 
        url string 
        depth int 
        body string 
        urls []string 
        err error 
} 

func Fetch(fetcher Fetcher, url string, depth int, responses chan 
	*response) { 
        body, urls, err := fetcher.Fetch(url) 
        responses <- &response{url, depth, body, urls, err} 
} 

// Crawl uses fetcher to recursively crawl 
// pages starting with url, to a maximum of depth. 
func Crawl(url string, depth int, fetcher Fetcher) { 
        if depth == 0 { 
                return 
        } 
        responses := make(chan *response, 1) 
        go Fetch(fetcher, url, depth, responses) 
        // Using maps as sets is inefficient, though convenient.  See http://blog.golang.org/2011/06/profiling-go-programs.html 
        seen := map[string]bool{ url: true } 
        pending := 1 
        for pending > 0 { 
                r := <- responses 
                pending-- 
                if r.err != nil { 
                        fmt.Println(r.err) 
                        continue 
                } 
                fmt.Printf("found: %s %q\n", r.url, r.body) 
                if r.depth == 0 { 
                        continue 
                } 
                for _, u := range r.urls { 
                        if !seen[u] { 
                                go Fetch(fetcher, u, r.depth-1, responses) 
                                seen[u] = true 
                                pending++ 
                        } 
                } 
        } 
} 

func main() {
	Crawl("http://golang.org/", 4, fetcher)
}


// fakeFetcher is Fetcher that returns canned results.
type fakeFetcher map[string]*fakeResult

type fakeResult struct {
	body string
	urls     []string
}

func (f *fakeFetcher) Fetch(url string) (string, []string, error) {
	if res, ok := (*f)[url]; ok {
		return res.body, res.urls, nil
	}
	return "", nil, fmt.Errorf("not found: %s", url)
}

// fetcher is a populated fakeFetcher.
var fetcher = &fakeFetcher{
	"http://golang.org/": &fakeResult{
		"The Go Programming Language",
		[]string{
			"http://golang.org/pkg/",
			"http://golang.org/cmd/",
		},
	},
	"http://golang.org/pkg/": &fakeResult{
		"Packages",
		[]string{
			"http://golang.org/",
			"http://golang.org/cmd/",
			"http://golang.org/pkg/fmt/",
			"http://golang.org/pkg/os/",
		},
	},
	"http://golang.org/pkg/fmt/": &fakeResult{
		"Package fmt",
		[]string{
			"http://golang.org/",
			"http://golang.org/pkg/",
		},
	},
	"http://golang.org/pkg/os/": &fakeResult{
		"Package os",
		[]string{
			"http://golang.org/",
			"http://golang.org/pkg/",
		},
	},
}

