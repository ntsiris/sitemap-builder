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

package link

import (
	"io"
	"strings"

	stack "github.com/ntsiris/sitemap-builder/pkg/collections"
	"golang.org/x/net/html"
)

// HTMLLink represents a hyperlink in an HTML document, containing the link's Href and the Text displayed.
type HTMLLink struct {
	Href string // URL destination of the link
	Text string // Text content within the link
}

// nodeProcessor is an interface for processing HTML nodes.
// It defines methods to check if a node qualifies for processing and to perform processing on the node.
type nodeProcessor interface {
	qualifies(node *html.Node) bool
	processNode(node *html.Node)
}

// textNodeProcessor is used to process text nodes within HTML elements,
// collecting and building the text content.
type textNodeProcessor struct {
	textBuilder *strings.Builder
}

func (tp *textNodeProcessor) qualifies(node *html.Node) bool {
	return node.Type == html.TextNode
}

func (tp *textNodeProcessor) processNode(node *html.Node) {
	tp.textBuilder.WriteString(strings.TrimSpace(node.Data))
}

// linkNodeProcessor is a processor for identifying and processing anchor (<a>) nodes.
type linkNodeProcessor struct{}

func (lp *linkNodeProcessor) qualifies(node *html.Node) bool {
	return node.Type == html.ElementNode && node.Data == "a"
}

func (lp *linkNodeProcessor) processNode(node *html.Node) {}

// Parse reads an HTML document from an io.Reader, extracts anchor (<a>) nodes,
// and returns a slice of HTMLLink objects representing the links found.
func Parse(r io.Reader) ([]HTMLLink, error) {
	htmlDocument, err := html.Parse(r)

	if err != nil {
		return nil, err
	}

	// Extract anchor nodes from the document
	nodes := dfsHTMLNodes(htmlDocument, &linkNodeProcessor{})

	var links []HTMLLink

	for _, node := range nodes {
		links = append(links, buildHTMLLink(node))
	}

	return links, nil
}

// dfsHTMLNodes performs a depth-first search (DFS) on an HTML node tree,
// starting from the seed node, and applies a nodeProcessor to each node.
// It returns a slice of nodes that qualify based on the nodeProcessor's criteria.
func dfsHTMLNodes(seed *html.Node, pr nodeProcessor) []*html.Node {
	var ret []*html.Node
	visited := make(map[*html.Node]bool)
	nodeStack := stack.Stack{}

	nodeStack.Push(seed)

	for nodeStack.Len() > 0 {
		node := nodeStack.Pop().(*html.Node)

		if !visited[node] {
			visited[node] = true

			for child := node.FirstChild; child != nil; child = child.NextSibling {
				nodeStack.Push(child)
			}

			if pr.qualifies(node) {
				pr.processNode(node)
				ret = append(ret, node)
			}
		}
	}

	return ret
}

func buildHTMLLink(node *html.Node) HTMLLink {
	var ret HTMLLink

	for _, attr := range node.Attr {
		if attr.Key == "href" {
			ret.Href = attr.Val
			break
		}
	}

	ret.Text = getLinkText(node)

	return ret
}

// getLinkText traverses the given node's subtree, collects text content from all
// text nodes within it, and returns the concatenated result as a single string.
func getLinkText(node *html.Node) string {

	tp := textNodeProcessor{}
	tp.textBuilder = &strings.Builder{}

	_ = dfsHTMLNodes(node, &tp)

	return tp.textBuilder.String()
}
