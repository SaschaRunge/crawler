package main

import (
	_ "strings"
	"testing"
)

func TestGetHeadingFromHTML(t *testing.T) {
	tests := []struct {
		name          string
		html          string
		expected      string
		errorContains string
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
