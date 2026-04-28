package main

import (
	"fmt"
	"net/url"
	"os"
	"strconv"
	"sync"
)

func main() {
	args := os.Args[1:]

	if len(args) < 3 {
		fmt.Println("no website provided")
		os.Exit(1)
	}
	if len(args) > 3 {
		fmt.Println("too many arguments provided")
		os.Exit(1)
	}

	baseURL, err := url.Parse(args[0])
	if err != nil {
		fmt.Println("invalid base url in argument 1")
		os.Exit(1)
	}

	maxConcurrency, err := strconv.ParseInt(args[1], 10, 32)
	if err != nil {
		fmt.Println("invalid value for max concurrency in argument 2")
		os.Exit(1)
	}

	maxPages, err := strconv.ParseInt(args[2], 10, 32)
	if err != nil {
		fmt.Println("invalid value for max pages in argument 3")
		os.Exit(1)
	}

	cfg := config{
		pages:              map[string]PageData{},
		baseURL:            baseURL,
		mu:                 &sync.Mutex{},
		concurrencyControl: make(chan struct{}, maxConcurrency),
		wg:                 &sync.WaitGroup{},
		maxPages:           int(maxPages),
	}

	cfg.wg.Go(func() {
		cfg.crawlPage(baseURL.String())
	})
	cfg.wg.Wait()
	fmt.Println("done, found the following references:")

	for k, v := range cfg.pages {
		fmt.Printf("%s: %s\n", k, v.Heading)
	}
}
