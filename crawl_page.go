package main

import (
	"fmt"
	"net/url"
)

func crawlPage(rawBaseURL, rawCurrentURL string, pages map[string]int) {
	baseURL, err := url.Parse(rawBaseURL)
	if err != nil {
		return
	}
	currentURL, err := url.Parse(rawCurrentURL)
	if err != nil {
		return
	}
	if baseURL.Hostname() != currentURL.Hostname() {
		return
	}

	normalizedURL, err := normalizeURL(rawCurrentURL)
	if err != nil {
		return
	}

	if _, ok := pages[normalizedURL]; ok {
		pages[normalizedURL]++
		return
	}
	pages[normalizedURL] = 1

	fmt.Printf("crawling %s ...\n", normalizedURL)
	fmt.Println("fetching html ...")
	html, err := getHTML(normalizedURL)
	if err != nil {
		fmt.Printf("error fetching html for %s: %s\n", normalizedURL, err)
		return
	}

	fmt.Println("extracting urls ...")
	urls, err := getURLsFromHTML(html, baseURL)
	if err != nil {
		fmt.Printf("error extracting urls: %s\n", err)
		return
	}

	for _, u := range urls {
		crawlPage(rawBaseURL, u, pages)
	}
}
