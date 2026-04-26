package main

import (
	"fmt"
	"os"
)

func main() {
	args := os.Args[1:]

	if len(args) < 1 {
		fmt.Println("no website provided")
		os.Exit(1)
	}
	if len(args) > 1 {
		fmt.Println("too many arguments provided")
		os.Exit(1)
	}

	baseURL := args[0]
	fmt.Printf("starting crawl of: %s\n", baseURL)
	html, err := getHTML(baseURL)
	if err != nil {
		fmt.Printf("failed to fetch with error: %s", err)
		os.Exit(1)
	}

	fmt.Printf("%s", html)
}
