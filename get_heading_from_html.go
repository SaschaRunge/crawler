package main

import (
	"fmt"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func getHeadingFromHTML(html string) string {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		fmt.Printf("warning: unexpected error when reading header: %s", err)
		return ""
	}

	singleSel := doc.FindMatcher(goquery.Single("h1"))
	if singleSel.Text() == "" {
		singleSel = doc.FindMatcher(goquery.Single("h2"))
	}

	return singleSel.Text()
}
