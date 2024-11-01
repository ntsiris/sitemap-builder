package utils

import (
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"

	"github.com/ntsiris/sitemap-builder/pkg/link"
)

func ParsePageLinks(urlStr string) (string, []string) {
	resp, err := http.Get(urlStr)
	if err != nil {
		log.Fatalf("Failed to get site \"%s\": %v", urlStr, err)
	}
	// Not closing the response body can cause a memory leak
	defer resp.Body.Close()

	base := baseURL(resp.Request.URL)

	return base, filterAndNormalizeLinks(resp.Body, base)
}

func FilterExternal(links []string, keepFn func(string) bool) []string {
	var retLinks []string

	for _, l := range links {
		if keepFn(l) {
			retLinks = append(retLinks, l)
		}
	}

	return retLinks
}

func filterAndNormalizeLinks(r io.Reader, base string) []string {
	links, _ := link.Parse(r)
	var hrefs []string
	for _, l := range links {
		switch {
		case strings.HasPrefix(l.Href, "/"):
			hrefs = append(hrefs, base+l.Href)
		case strings.HasPrefix(l.Href, "http"):
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
