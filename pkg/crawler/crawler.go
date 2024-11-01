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

package crawler

import (
	"net/http"
	"net/url"
	"strings"

	"github.com/ntsiris/sitemap-builder/pkg/link"
)

// Crawl performs a web crawl starting from the provided URL (urlStr) up to a specified depth (maxDepth).
// It uses a keepFn function to determine which URLs should be kept. The function returns a slice
// of all unique URLs found within the depth limit, or an error if something goes wrong.
func Crawl(urlStr string, maxDepth int, keepFn func(string) bool) ([]string, error) {
	seen := make(map[string]struct{})

	// queue holds URLs to process at the current depth level
	// nextQueue holds URLs to process at the next depth level
	var queue map[string]struct{}
	nextQueue := map[string]struct{}{
		urlStr: {}, // Start with the initial URL in the nextQueue
	}

	for i := 0; i <= maxDepth; i++ {
		queue, nextQueue = nextQueue, make(map[string]struct{})

		for url := range queue {
			if _, ok := seen[url]; ok {
				continue
			}
			seen[url] = struct{}{}

			for _, link := range parsePageLinks(url, keepFn) {
				nextQueue[link] = struct{}{}
			}
		}
	}

	ret := make([]string, 0, len(seen))
	for url := range seen {
		ret = append(ret, url)
	}

	return ret, nil
}

// parsePageLinks fetches and parses all links on the page specified by urlStr.
// It filters links based on the keepFn function and normalizes relative links using the base URL.
func parsePageLinks(urlStr string, keepFn func(string) bool) []string {
	resp, err := http.Get(urlStr)
	if err != nil {
		return []string{}
	}
	// Ensure the response body is closed to avoid memory leaks
	defer resp.Body.Close()

	base := baseURL(resp.Request.URL)

	links, _ := link.Parse(resp.Body)

	return filterAndNormalizeLinks(links, base, keepFn)
}

// filterAndNormalizeLinks takes parsed HTML links, a base URL for normalization, and a keepFn
// function to filter external links. It returns a slice of strings containing the URLs to keep.
func filterAndNormalizeLinks(links []link.HTMLLink, base string, keepFn func(string) bool) []string {
	var hrefs []string
	for _, l := range links {
		switch {
		case strings.HasPrefix(l.Href, "/"):
			hrefs = append(hrefs, base+l.Href)
		case keepFn(l.Href):
			hrefs = append(hrefs, l.Href)
		}
	}

	return hrefs
}

func baseURL(u *url.URL) string {
	base := &url.URL{
		Scheme: u.Scheme,
		Host:   u.Host,
	}

	return base.String()
}
