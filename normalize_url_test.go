package main

import (
	"strings"
	"testing"
)

func TestNormalizeURL(t *testing.T) {
	tests := []struct {
		name          string
		inputURL      string
		expected      string
		errorContains string
	}{
		{
			name:     "remove https",
			inputURL: "https://www.boot.dev/blog/path",
			expected: "www.boot.dev/blog/path",
		},
		{
			name:     "remove http",
			inputURL: "http://www.boot.dev/blog/path",
			expected: "www.boot.dev/blog/path",
		},
		{
			name:     "remove trailing /",
			inputURL: "www.boot.dev/blog/path/",
			expected: "www.boot.dev/blog/path",
		},
		{
			name:     "remove query /",
			inputURL: "www.boot.dev/blog/path?abc=123",
			expected: "www.boot.dev/blog/path",
		},
		{
			name:     "remove fragment /",
			inputURL: "www.boot.dev/blog/path#PATH",
			expected: "www.boot.dev/blog/path",
		},
		{
			name:     "lower case",
			inputURL: "www.bOoT.dev/BLOG/path",
			expected: "www.boot.dev/blog/path",
		},
		{
			name:          "invalid url",
			inputURL:      ":\\this.Is.Not.A.URL",
			expected:      "",
			errorContains: "unable to parse URL",
		},
	}

	for i, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			actual, err := normalizeURL(tc.inputURL)
			if err != nil {
				if !strings.Contains(err.Error(), tc.errorContains) {
					t.Errorf("Test %v - '%s' FAIL: unexpected error: %v", i, tc.name, err)
				}
			}

			if actual != tc.expected {
				t.Errorf("Test %v - %s FAIL: expected URL: %v, actual: %v", i, tc.name, tc.expected, actual)
			}
		})
	}
}
