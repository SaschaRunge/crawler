package main

import (
	"fmt"
	"net/url"
)

func crawlPage(rawBaseURL, rawCurrentURL string, pages map[string]int) {
	baseURL, err := url.Parse(rawBaseURL)
	if err != nil {
		fmt.Printf("invalid base url: %s, %s\n", rawBaseURL, err)
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
		fmt.Printf("unable to normalize url, shouldn't happen: %s, %s\n", rawCurrentURL, err)
		return
	}

	if _, ok := pages[normalizedURL]; ok {
		pages[normalizedURL]++
		return
	}
	pages[normalizedURL] = 1

	fmt.Printf("crawling %s ...\n", currentURL.String())
	fmt.Println("fetching html ...")
	html, err := getHTML(currentURL.String())
	if err != nil {
		fmt.Printf("error fetching html for %s: %s\n", currentURL.String(), err)
		return
	}

	fmt.Println("extracting urls ...")
	urls, err := getURLsFromHTML(html, baseURL)
	if err != nil {
		fmt.Printf("error extracting urls: %s\n", err)
		return
	}

	//time.Sleep(time.Millisecond * 500)
	for _, u := range urls {
		crawlPage(rawBaseURL, u, pages)
	}
}
