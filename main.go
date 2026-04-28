package main

import (
	"fmt"
	"net/url"
	"os"
	"sync"
)

const (
	maxConcurrentRoutines = 10
)

func main() {
	args := os.Args[1:]

	if len(args) < 1 {
		fmt.Println("no website provided")
		os.Exit(1)
	}
	if len(args) > 1 {
		fmt.Println("too many arguments provided")
		os.Exit(1)
	}

	baseURL, err := url.Parse(args[0])
	if err != nil {
		fmt.Println("invalid base url")
		os.Exit(1)
	}
	cfg := config{
		pages:              map[string]PageData{},
		baseURL:            baseURL,
		mu:                 &sync.Mutex{},
		concurrencyControl: make(chan struct{}, maxConcurrentRoutines),
		wg:                 &sync.WaitGroup{},
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
