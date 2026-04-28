package main

import (
	"fmt"
	"net/url"
	"sync"
	"time"
)

type config struct {
	pages              map[string]PageData
	baseURL            *url.URL
	mu                 *sync.Mutex
	concurrencyControl chan struct{}
	wg                 *sync.WaitGroup
}

func (cfg *config) crawlPage(rawCurrentURL string) {
	currentURL, err := url.Parse(rawCurrentURL)
	if err != nil {
		return
	}
	if cfg.baseURL.Hostname() != currentURL.Hostname() {
		return
	}

	normalizedURL, err := normalizeURL(rawCurrentURL)
	if err != nil {
		fmt.Printf("unable to normalize url, shouldn't happen: %s, %s\n", rawCurrentURL, err)
		return
	}

	if _, ok := cfg.pages[normalizedURL]; ok {
		return
	}

	fmt.Printf("crawling %s ...\n", currentURL.String())
	fmt.Println("fetching html ...")
	html, err := getHTML(currentURL.String())
	if err != nil {
		fmt.Printf("error fetching html for %s: %s\n", currentURL.String(), err)
		return
	}

	cfg.pages[normalizedURL] = extractPageData(html, currentURL.String())

	fmt.Println("extracting urls ...")
	urls, err := getURLsFromHTML(html, cfg.baseURL)
	if err != nil {
		fmt.Printf("error extracting urls: %s\n", err)
		return
	}

	time.Sleep(time.Millisecond * 500)
	for _, u := range urls {
		cfg.crawlPage(u)
	}
}
