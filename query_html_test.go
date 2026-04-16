package main

import (
	"net/url"
	"reflect"
	_ "strings"
	"testing"
)

func TestGetHeadingFromHTML(t *testing.T) {
	tests := []struct {
		name     string
		html     string
		expected string
	}{
		{
			name: "get h1 header",
			html: `<html>
  <body>
    <h1>Welcome to Boot.dev</h1>
    <main>
      <p>Learn to code by building real projects.</p>
      <p>This is the second paragraph.</p>
    </main>
  </body>
</html>`,
			expected: "Welcome to Boot.dev",
		},
		{
			name:     "get h2 header when h1 empty",
			html:     "<html><body><h1></h1><main><h2>This is a secondary header.</h2><p>Learn to code by building real projects.</p><p>This is the second paragraph.</p></main></body></html>",
			expected: "This is a secondary header.",
		},
		{
			name:     "get h2 header when h1 missing",
			html:     "<html><body><main><h2>This is a secondary header.</h2><p>Learn to code by building real projects.</p><p>This is the second paragraph.</p></main></body></html>",
			expected: "This is a secondary header.",
		},
		{
			name:     "no header found",
			html:     "<html><body><main><p>Learn to code by building real projects.</p><p>This is the second paragraph.</p></main></body></html>",
			expected: "",
		},
	}

	for i, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			actual := getHeadingFromHTML(tc.html)
			if actual != tc.expected {
				t.Errorf("Test %v - %s FAIL: expected: %s, actual: %v", i, tc.name, tc.expected, actual)
			}
		})
	}
}

func TestGetFirstParagraphFromHTML(t *testing.T) {
	tests := []struct {
		name     string
		html     string
		expected string
	}{
		{
			name:     "get main paragraph",
			html:     "<html><body><h1>Welcome to Boot.dev</h1><p>Fallback.</p><main><p>Learn to code by building real projects.</p><p>This is the second paragraph.</p></main></body></html>",
			expected: "Learn to code by building real projects.",
		},
		{
			name:     "get fallback paragraph",
			html:     "<html><body><h1>Welcome to Boot.dev</h1><p>Fallback.</p><main></main></body></html>",
			expected: "Fallback.",
		},
		{
			name:     "no paragraphs",
			html:     "<html><body><h1>Welcome to Boot.dev</h1><main></main></body></html>",
			expected: "",
		},
	}

	for i, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			actual := getFirstParagraphFromHTML(tc.html)
			if actual != tc.expected {
				t.Errorf("Test %v - %s FAIL: expected: %s, actual: %v", i, tc.name, tc.expected, actual)
			}
		})
	}
}

func TestGetURLsFromHTML(t *testing.T) {
	tests := []struct {
		name          string
		inputURL      string
		html          string
		expected      []string
		errorContains string
	}{
		{
			name:     "extract url",
			inputURL: "https://crawler-test.com",
			html:     `<html><body><a href="https://crawler-test.com"><span>Boot.dev</span></a></body></html>`,
			expected: []string{"https://crawler-test.com"},
		},
		{
			name:     "extract multiple urls",
			inputURL: "https://crawler-test.com",
			html:     `<html><body><a href="https://crawler-test.com"><span>Boot.dev</span></a><a href="https://somewhereelse.com/test">test</a></body></html>`,
			expected: []string{"https://crawler-test.com", "https://somewhereelse.com/test"},
		},
		{
			name:     "extract relative",
			inputURL: "https://crawler-test.com",
			html:     `<html><body><a href="/nextSite"><span>Boot.dev</span></a></body></html>`,
			expected: []string{"https://crawler-test.com/nextSite"},
		},
		{
			name:          "parse error",
			inputURL:      "https://crawler-test.com",
			html:          `<html><body><a href=":\\nextSite"><span>Boot.dev</span></a></body></html>`,
			expected:      []string{},
			errorContains: "unable to parse all links",
		},
	}

	for _, test := range tests {
		baseURL, err := url.Parse(test.inputURL)
		if err != nil {
			t.Errorf("couldn't parse input URL: %v", err)
			return
		}

		actual, err := getURLsFromHTML(test.html, baseURL)
		if err != nil && test.errorContains != err.Error() {
			t.Errorf("get url returned error: %s", err)
			return
		}

		if !reflect.DeepEqual(actual, test.expected) {
			t.Errorf("expected %v, got %v", test.expected, actual)
		}
	}
}
