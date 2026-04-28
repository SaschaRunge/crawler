package main

import (
	"fmt"
	"net/url"
	"sync"
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

	if isFirstVisit := cfg.addPageVisit(normalizedURL); !isFirstVisit {
		return
	}

	cfg.concurrencyControl <- struct{}{}
	fmt.Printf("crawling %s ...\n", currentURL.String())
	html, err := getHTML(currentURL.String())
	<-cfg.concurrencyControl
	if err != nil {
		fmt.Printf("error fetching html for %s: %s\n", currentURL.String(), err)
		return
	}

	data := extractPageData(html, currentURL.String())
	cfg.mu.Lock()
	cfg.pages[normalizedURL] = data
	cfg.mu.Unlock()

	for _, link := range data.OutgoingLinks {
		cfg.wg.Go(func() {
			cfg.crawlPage(link)
		})
	}
}

func (cfg *config) addPageVisit(normalizedURL string) (isFirstVisit bool) {
	cfg.mu.Lock()
	defer cfg.mu.Unlock()

	if _, visited := cfg.pages[normalizedURL]; visited {
		return false
	}

	cfg.pages[normalizedURL] = PageData{}
	return true
}
