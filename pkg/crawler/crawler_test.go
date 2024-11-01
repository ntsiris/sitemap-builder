package crawler

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/ntsiris/sitemap-builder/pkg/link"
)

// TestCrawl verifies that the Crawl function follows links up to the specified depth,
// and filters them using the provided keep function.
func TestCrawl(t *testing.T) {
	// Set up a mock HTTP server with some linked pages
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/" {
			w.Write([]byte(`<a href="/page1">Page 1</a><a href="https://external.com">External</a>`))
		} else if r.URL.Path == "/page1" {
			w.Write([]byte(`<a href="/page2">Page 2</a>`))
		} else if r.URL.Path == "/page2" {
			w.Write([]byte(`<a href="/">Home</a>`))
		}
	}))
	defer mockServer.Close()

	// keepFn filters out external links
	keepFn := func(link string) bool {
		return strings.HasPrefix(link, mockServer.URL)
	}

	// Call Crawl with maxDepth 2 to test link discovery across multiple pages
	links, err := Crawl(mockServer.URL, 2, keepFn)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	// Add diagnostic log for debugging
	t.Logf("Links found: %v", links)

	// Expected links based on the pages and depth
	expectedLinks := []string{
		mockServer.URL,            // Base URL
		mockServer.URL + "/page1", // Linked from "/"
		mockServer.URL + "/page2", // Linked from "/page1"
	}

	// Check that each expected link is in the result
	for _, expected := range expectedLinks {
		found := false
		for _, link := range links {
			if link == expected {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Expected link %s to be in result, but it was not found", expected)
		}
	}
}

// TestParsePageLinks tests that parsePageLinks correctly extracts and normalizes links,
// and filters them based on the keep function.
func TestParsePageLinks(t *testing.T) {
	// HTML content simulating a page with internal and external links
	htmlContent := `<a href="/about">About</a><a href="https://external.com">External</a>`
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(htmlContent))
	}))
	defer server.Close()

	// Define a keep function to filter for external links only
	keepFn := func(link string) bool {
		return strings.Contains(link, "external")
	}

	// Run parsePageLinks and verify that it correctly applies base URL and filtering
	links := parsePageLinks(server.URL, keepFn)

	// Check that the results contain the expected links
	expectedLinks := []string{server.URL + "/about", "https://external.com"}
	for _, expected := range expectedLinks {
		found := false
		for _, link := range links {
			if link == expected {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Expected link %s to be in result, but it was not found", expected)
		}
	}
}

// TestFilterAndNormalizeLinks tests that filterAndNormalizeLinks normalizes relative URLs and applies keepFn to filter links.
func TestFilterAndNormalizeLinks(t *testing.T) {
	// Simulate links parsed from an HTML document
	links := []link.HTMLLink{
		{Href: "/contact", Text: "Contact"},
		{Href: "https://external.com", Text: "External"},
	}

	// Base URL
	base := "https://example.com"
	// Keep function to filter for external links only
	keepFn := func(link string) bool {
		return strings.HasPrefix(link, base)
	}

	result := filterAndNormalizeLinks(links, base, keepFn)

	// Expected result after filtering and normalizing
	expectedLinks := []string{"https://example.com/contact"}
	if len(result) != len(expectedLinks) {
		t.Errorf("Expected %d links, got %d", len(expectedLinks), len(result))
	}
	for i, expected := range expectedLinks {
		if result[i] != expected {
			t.Errorf("Expected link %s, got %s", expected, result[i])
		}
	}
}

// TestBaseURL tests that baseURL constructs the base part of a URL, preserving scheme and host.
func TestBaseURL(t *testing.T) {
	// Parse a sample URL
	parsedURL, _ := url.Parse("https://example.com/path/to/page")

	// Generate the base URL
	base := baseURL(parsedURL)

	// Expected base URL with only scheme and host
	expected := "https://example.com"
	if base != expected {
		t.Errorf("Expected base URL %s, got %s", expected, base)
	}
}
