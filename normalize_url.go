package main

import (
	"fmt"
	"net/url"
	"strings"
)

func normalizeURL(input string) (string, error) {
	u, err := url.Parse(input)
	if err != nil {
		return "", fmt.Errorf("unable to parse URL: %w", err)
	}

	path, _ := strings.CutSuffix(u.Path, "/")
	return strings.ToLower(u.Host + path), nil
}
