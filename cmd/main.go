package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strings"

	"github.com/ntsiris/sitemap-builder/pkg/crawler"
)

func main() {
	urlFlag := flag.String("url", "", "Specify the URL to be processed")
	maxDepth := flag.Int("depth", 3, "The maximum number of links deep to traverse")

	flag.Parse()

	if *urlFlag == "" {
		fmt.Println("No URL provided. Please specify a URL with the -url flag.")
	}

	resp, err := http.Get(*urlFlag)
	if err != nil {
		log.Fatalf("Failed to get site \"%s\": %v", *urlFlag, err)
	}
	// Not closing the response body can cause a memory leak
	defer resp.Body.Close()

	base := &url.URL{
		Scheme: resp.Request.URL.Scheme,
		Host:   resp.Request.URL.Host,
	}

	pages, err := crawler.Crawl(*urlFlag, *maxDepth, withPrefix(base.String()))

	for _, page := range pages {
		fmt.Println(page)
	}
}

func withPrefix(prefix string) func(string) bool {
	return func(link string) bool {
		return strings.HasPrefix(link, prefix)
	}
}
