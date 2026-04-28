package main

import (
	"fmt"
	"net/url"
	"os"
	"sync"
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
		concurrencyControl: make(chan struct{}),
		wg:                 &sync.WaitGroup{},
	}

	cfg.crawlPage(baseURL.String())
	fmt.Println("done, found the following references:")

	for k, v := range cfg.pages {
		fmt.Printf("%s: %s\n", k, v.Heading)
	}
}
