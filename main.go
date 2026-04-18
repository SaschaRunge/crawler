package main

import (
	"fmt"
)

func main() {
	fmt.Println("Hello, World!")

	html := "<html><body><h1>Welcome to Boot.dev</h1><p>outside of main</p><main><p>Learn to code by building real projects.</p><p>This is the second paragraph.</p></main></body></html>"
	getFirstParagraphFromHTML(html)

}
