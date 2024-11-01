package main

import (
	"flag"
	"fmt"
	"strings"

	"github.com/ntsiris/sitemap-builder/internal/utils"
)

func main() {
	urlFlag := flag.String("url", "", "Specify the URL to be processed")

	flag.Parse()

	if *urlFlag == "" {
		fmt.Println("No URL provided. Please specify a URL with the -url flag.")
	}

	/*

		Process:
			1. Get the HTML page
			2. Parse all the links on the page
			3. Build proper urls with our links
			4. Filter out any links with a different domain
			5. Find all pages (bfs)
			6. Generate XML
	*/

	/*
		Link cases:
			Handle:
				-> /some-path [add domain]
				-> https://example.com/some-path

			Do not Handle:
				-> #fragment [Don't handle]
				-> mailto:someone@example.com [Don't handle]
	*/

	base, pages := utils.ParsePageLinks(*urlFlag)
	pages = utils.FilterExternal(pages, withPrefix(base))

	for _, page := range pages {
		fmt.Println(page)
	}
}

func withPrefix(prefix string) func(string) bool {
	return func(link string) bool {
		return strings.HasPrefix(link, prefix)
	}
}
