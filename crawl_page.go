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
	maxPages           int
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

	if stop := cfg.tryReservePage(normalizedURL); stop {
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

func (cfg *config) tryReservePage(normalizedURL string) (stop bool) {
	cfg.mu.Lock()
	defer cfg.mu.Unlock()

	if _, visited := cfg.pages[normalizedURL]; visited {
		return true
	}

	if cfg.stopCrawlingLocked() {
		return true
	}

	cfg.pages[normalizedURL] = PageData{}
	return false
}

// caller needs to lock
func (cfg *config) stopCrawlingLocked() bool {
	return len(cfg.pages) >= cfg.maxPages
}
