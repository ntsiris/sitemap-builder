package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

func main() {
	urlFlag := flag.String("url", "", "Specify the URL to be processed")

	flag.Parse()

	if *urlFlag == "" {
		fmt.Println("No URL provided. Please specify a URL with the -url flag.")
	} else {
		fmt.Printf("URL provided: %s\n", *urlFlag)
	}

	/*
		1. Get the HTML page
		2. Parse all the links on the page
		3. Build proper urls with our links
		4. Filter out any links with a different domain
		5. Find all pages (bfs)
		6. Generate XML
	*/

	// Get the HTML Page
	resp, err := http.Get(*urlFlag)
	if err != nil {
		log.Fatalf("Failed to get site \"%s\": %v", *urlFlag, err)
	}
	// Not closing the response body can cause a memory leak
	defer resp.Body.Close()

	io.Copy(os.Stdout, resp.Body)
}
