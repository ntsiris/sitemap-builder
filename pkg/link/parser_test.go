package link

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"golang.org/x/net/html"
)

// TestParse_SingleLink tests that Parse identifies a single link with correct href and text content.
func TestParse_SingleLink(t *testing.T) {
	htmlContent := `<html><body><a href="http://example.com">Example</a></body></html>`
	r := strings.NewReader(htmlContent)

	links, err := Parse(r)
	assert.NoError(t, err, "Parse should not return an error for valid HTML")
	assert.Equal(t, 1, len(links), "Expected one link in the parsed HTML")
	assert.Equal(t, "http://example.com", links[0].Href, "Expected href to match the link's href attribute")
	assert.Equal(t, "Example", links[0].Text, "Expected link text to match the anchor's text content")
}

// TestParse_MultipleLinks tests that Parse identifies multiple links with correct hrefs and text contents.
func TestParse_MultipleLinks(t *testing.T) {
	htmlContent := `<html><body>
		<a href="http://example1.com">Example 1</a>
		<a href="http://example2.com">Example 2</a>
		<a href="http://example3.com">Example 3</a>
	</body></html>`
	r := strings.NewReader(htmlContent)

	links, err := Parse(r)
	assert.NoError(t, err, "Parse should not return an error for valid HTML")
	assert.Equal(t, 3, len(links), "Expected three links in the parsed HTML")

	expectedLinks := []HTMLLink{
		{Href: "http://example3.com", Text: "Example 3"},
		{Href: "http://example2.com", Text: "Example 2"},
		{Href: "http://example1.com", Text: "Example 1"},
	}
	assert.Equal(t, expectedLinks, links, "Parsed links should match expected links")
}

// TestParse_Empty tests that Parse handles HTML with no links.
func TestParse_Empty(t *testing.T) {
	htmlContent := `<html><body><p>No links here!</p></body></html>`
	r := strings.NewReader(htmlContent)

	links, err := Parse(r)
	assert.NoError(t, err, "Parse should not return an error for valid HTML with no links")
	assert.Equal(t, 0, len(links), "Expected no links in the parsed HTML")
}

// TestDFSHTMLNodes_AnchorNode checks that dfsHTMLNodes identifies anchor nodes.
func TestDFSHTMLNodes_AnchorNode(t *testing.T) {
	htmlContent := `<html><body><a href="http://example.com">Example</a></body></html>`
	doc, _ := html.Parse(strings.NewReader(htmlContent))

	processor := &linkNodeProcessor{}
	nodes := dfsHTMLNodes(doc, processor)
	assert.Equal(t, 1, len(nodes), "Expected one anchor node")
	assert.Equal(t, "a", nodes[0].Data, "Expected the node to be an anchor element")
}

// TestBuildHTMLLink tests the buildHTMLLink function for correct href and text extraction.
func TestBuildHTMLLink(t *testing.T) {
	htmlContent := `<a href="http://example.com">Example Link</a>`
	// Parse the content as a complete HTML document to ensure it matches the expected structure.
	doc, _ := html.Parse(strings.NewReader(`<html><body>` + htmlContent + `</body></html>`))

	// Traverse to the <a> node within the parsed document structure.
	body := doc.FirstChild.LastChild // Assumes structure is <html><body>...</body></html>
	aNode := body.FirstChild         // Now, aNode should point to the <a> tag.

	// Now call buildHTMLLink with the actual <a> node.
	link := buildHTMLLink(aNode)
	assert.Equal(t, "http://example.com", link.Href, "Expected href to match the link's href attribute")
	assert.Equal(t, "Example Link", link.Text, "Expected link text to match the anchor's text content")
}

// TestGetLinkText tests the getLinkText function with nested text content.
func TestGetLinkText(t *testing.T) {
	htmlContent := `<a href="http://example.com"><span>Nested</span> Text</a>`
	// Parse as a full document to ensure the correct structure.
	doc, _ := html.Parse(strings.NewReader(`<html><body>` + htmlContent + `</body></html>`))

	// Locate the <a> node within the document structure.
	body := doc.FirstChild.LastChild // Navigate to the <body> tag.
	aNode := body.FirstChild         // <a> should be the first child of <body>.

	// Call getLinkText on the <a> node.
	text := getLinkText(aNode)
	assert.Equal(t, "TextNested", text, "Expected getLinkText to return concatenated text content")
}

// TestTextNodeProcessor tests correctly qualifies and processes text nodes.
func TestTextNodeProcessor(t *testing.T) {
	htmlContent := `<p>This is a <span>test</span> text.</p>`
	doc, _ := html.ParseFragment(strings.NewReader(htmlContent), nil)

	textProcessor := &textNodeProcessor{textBuilder: &strings.Builder{}}
	_ = dfsHTMLNodes(doc[0], textProcessor)
	result := textProcessor.textBuilder.String()

	expectedText := "text.testThis is a"
	assert.Equal(t, expectedText, result, "Expected concatenated text content to match")
}
