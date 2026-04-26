package main

import (
	"fmt"
	"io"
	"net/http"
	"strings"
)

func getHTML(rawURL string) (string, error) {
	client := http.DefaultClient
	req, err := http.NewRequest("GET", rawURL, nil)
	if err != nil {
		return "", err
	}

	req.Header.Add("User-Agent", "BootCrawler/1.0")
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	if resp.StatusCode >= 400 {
		return "", fmt.Errorf("response returned error type status: %d", resp.StatusCode)
	}
	contentType := resp.Header.Get("Content-Type")
	if !strings.Contains(strings.ToLower(contentType), "text/html") {
		return "", fmt.Errorf("invalid response type: %s", contentType)
	}
	defer resp.Body.Close()

	html, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(html), nil
}
