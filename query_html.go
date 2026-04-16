package main

import (
	"errors"
	"fmt"
	"net/url"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func getHeadingFromHTML(html string) string {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		fmt.Printf("warning: unexpected error when reading header: %s", err) // maybe remove, see how much log clutter
		return ""
	}

	singleSel := doc.FindMatcher(goquery.Single("h1"))
	if singleSel.Text() == "" {
		singleSel = doc.FindMatcher(goquery.Single("h2"))
	}

	return singleSel.Text()
}

func getFirstParagraphFromHTML(html string) string {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		fmt.Printf("warning: unexpected error when reading paragraph: %s", err)
		return ""
	}

	sel := doc.Find("main").Find("p").First()
	if sel.Text() == "" {
		sel = doc.Find("p").First()
	}
	return sel.Text()
}

func getURLsFromHTML(html string, baseURL *url.URL) ([]string, error) {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		fmt.Printf("warning: unexpected error retrieving urls: %s", err)
		return []string{}, err
	}

	urls := []string{}
	containsParseError := false
	doc.Find("a[href]").Each(func(_ int, s *goquery.Selection) {
		attr, exists := s.Attr("href")
		if !exists {
			return
		}

		link, err := url.Parse(attr)
		if err != nil {
			containsParseError = true
			return
		}

		absURL := baseURL.ResolveReference(link)
		urls = append(urls, absURL.String())
	})

	err = nil
	if containsParseError {
		err = errors.New("unable to parse all links")
	}

	return urls, err
}
