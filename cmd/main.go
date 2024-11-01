package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/ntsiris/sitemap-builder/pkg/crawler"
)

func main() {
	urlFlag := flag.String("url", "", "Specify the URL to be processed")
	maxDepth := flag.Int("depth", 3, "The maximum number of links deep to traverse")
	outputFile := flag.String("out", "./map.xml", "Specify the output file for the sitemap")

	flag.Parse()

	if *urlFlag == "" {
		fmt.Println("No URL provided. Please specify a URL with the -url flag.")
	}

	resp, err := http.Get(*urlFlag)
	if err != nil {
		log.Fatalf("Failed to get site \"%s\": %v", *urlFlag, err)
	}
	// Not closing the response body can cause a memory leak
	defer func() {
		err := resp.Body.Close()
		if err != nil {
			log.Print(err)
		}
	}()

	base := &url.URL{
		Scheme: resp.Request.URL.Scheme,
		Host:   resp.Request.URL.Host,
	}

	pages, err := crawler.Crawl(*urlFlag, *maxDepth, withPrefix(base.String()))

	file, err := os.OpenFile(*outputFile, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0664)
	if err != nil {
		log.Fatalf("Failed to open specified output file %s: %v", *outputFile, err)
	}
	defer func() {
		err = file.Close()
		if err != nil {
			log.Print(err)
		}
	}()

	GenerateSitemap(pages, file)
}

func withPrefix(prefix string) func(string) bool {
	return func(link string) bool {
		return strings.HasPrefix(link, prefix)
	}
}
