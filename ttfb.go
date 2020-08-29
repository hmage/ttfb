// TTFB
// Inspired by ttfb.sh, written by @jaygooby and @sandeepraju
package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	urlpkg "net/url"
	"sort"
	"strconv"
	"strings"
	"time"
)

func main() {
	numRequests := 1
	flag.IntVar(&numRequests, "n", 1, "number of times to fetch the url")
	flag.Parse()

	client := http.Client{}

	longest := 0
	results := map[string][]time.Duration{}
	for _, url := range flag.Args() {
		if len(url) > longest {
			longest = len(url)
		}
		// we use +1 because first request is always slow -- so use it as a warmup
		for i := 0; i < numRequests+1; i++ {
			if numRequests > 1 && i != 0 {
				fmt.Print(".")
			}
			parsed, err := urlpkg.Parse(url)
			if err != nil {
				log.Printf("Skipping because failed to parse address %s: %s", url, err)
				continue
			}
			if parsed.Scheme == "" {
				parsed.Scheme = "http"
			}
			req := http.Request{
				Method: "GET",
				URL:    parsed,
				Header: make(http.Header),
			}

			// that's a hint, servers are free to ignore the range request
			req.Header.Set("Range", "bytes=0-0")

			start := time.Now()
			resp, err := client.Do(&req)
			if err != nil {
				log.Printf("Failed to request %s: %s", url, err)
				continue
			}
			// read one byte
			buf := make([]byte, 1)
			n, err := resp.Body.Read(buf)
			if err != nil {
				log.Printf("Failed to read body from HTTP request to %s: %s", url, err)
				continue
			}
			if n != 1 {
				log.Printf("Failed to read body from HTTP request to %s: Wanted to read 1 byte, got %d", url, n)
				continue
			}
			elapsed := time.Since(start)
			// don't read rest of the body, we're not interested in it
			resp.Body.Close()
			// first run is warmup, don't record results
			if i == 0 {
				continue
			}
			results[url] = append(results[url], elapsed)
		}
		if numRequests > 1 {
			fmt.Println()
		}
	}

	// match output format with ttfb.sh
	f := func(t time.Duration) string {
		return strings.TrimLeft(strconv.FormatFloat(t.Seconds(), 'f', 6, 64), "0")
	}

	// output results in same order they were given in command line
	for _, url := range flag.Args() {
		elapsed := results[url]
		// sort result timings
		sort.Slice(elapsed, func(left, right int) bool {
			return elapsed[left] < elapsed[right]
		})
		median := elapsed[0]
		if len(elapsed)%2 == 1 {
			// odd number of elements
			median = elapsed[len(elapsed)/2]
		} else {
			// even number of elements
			right := len(elapsed) / 2
			left := right - 1
			median = (elapsed[left] + elapsed[right]) / 2
		}
		fastest := median
		slowest := median
		for _, val := range elapsed {
			if fastest > val {
				fastest = val
			}
			if slowest < val {
				slowest = val
			}
		}
		fmt.Printf("%-*s \x1b[32mfastest \x1b[39m%s \x1b[91mslowest \x1b[39m%s \x1b[95mmedian \x1b[39m%s\n", longest+1, url, f(fastest), f(slowest), f(median))
	}
}
